package com.coworking.reservations.service;

import com.coworking.reservations.client.SpaceServiceClient;
import com.coworking.reservations.dto.ReservationRequest;
import com.coworking.reservations.dto.ReservationResponse;
import com.coworking.reservations.exception.ResourceNotFoundException;
import com.coworking.reservations.model.Reservation;
import com.coworking.reservations.model.ReservationStatus;
import com.coworking.reservations.repository.ReservationRepository;
import feign.FeignException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.time.Duration;
import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class ReservationService {

    private final ReservationRepository reservationRepository;
    private final SpaceServiceClient spaceServiceClient;

    // -------------------------------------------------------------------------
    // CREATE
    // -------------------------------------------------------------------------

    @Transactional
    public ReservationResponse createReservation(
            ReservationRequest request,
            String userId,
            String authHeader
    ) {
        // 1. Validate temporal ordering
        if (!request.isEndTimeAfterStartTime()) {
            throw new IllegalArgumentException("End time must be after start time");
        }

        // 2. Validate space exists and is active via space-service
        SpaceServiceClient.SpaceDto space = fetchSpace(request.getSpaceId(), authHeader);

        if (Boolean.FALSE.equals(space.getIsActive())) {
            throw new IllegalArgumentException("Space '" + space.getName() + "' is not available for booking");
        }

        // 3. Check for overlapping PENDING or CONFIRMED reservations
        List<Reservation> overlaps = reservationRepository.findOverlappingReservations(
                request.getSpaceId(),
                request.getStartTime(),
                request.getEndTime()
        );

        if (!overlaps.isEmpty()) {
            throw new IllegalArgumentException(
                    "Space is already reserved during the requested time window");
        }

        // 4. Calculate total cost: pricePerHour * hours (rounded up to nearest quarter-hour)
        BigDecimal totalCost = calculateCost(
                space.getPricePerHour(),
                request.getStartTime(),
                request.getEndTime()
        );

        // 5. Create reservation — start PENDING then immediately CONFIRMED
        Reservation reservation = Reservation.builder()
                .userId(userId)
                .spaceId(request.getSpaceId())
                .startTime(request.getStartTime())
                .endTime(request.getEndTime())
                .status(ReservationStatus.PENDING)
                .totalCost(totalCost)
                .notes(request.getNotes())
                .build();

        reservation = reservationRepository.save(reservation);

        // Immediately confirm
        reservation.setStatus(ReservationStatus.CONFIRMED);
        reservation = reservationRepository.save(reservation);

        log.info("Reservation {} created and confirmed for user {} at space {}",
                reservation.getId(), userId, space.getName());

        return toResponse(reservation, space.getName());
    }

    // -------------------------------------------------------------------------
    // READ
    // -------------------------------------------------------------------------

    @Transactional(readOnly = true)
    public List<ReservationResponse> getMyReservations(String userId, String authHeader) {
        return reservationRepository
                .findByUserIdOrderByCreatedAtDesc(userId)
                .stream()
                .map(r -> toResponseWithSpaceName(r, authHeader))
                .collect(Collectors.toList());
    }

    @Transactional(readOnly = true)
    public List<ReservationResponse> getAllReservations(String authHeader) {
        return reservationRepository
                .findAllWithFilters(null, null)
                .stream()
                .map(r -> toResponseWithSpaceName(r, authHeader))
                .collect(Collectors.toList());
    }

    @Transactional(readOnly = true)
    public ReservationResponse getReservation(UUID id, String userId, String role, String authHeader) {
        Reservation reservation = findByIdOrThrow(id);
        validateOwnershipOrAdmin(reservation, userId, role);
        String spaceName = fetchSpaceNameSafe(reservation.getSpaceId(), authHeader);
        return toResponse(reservation, spaceName);
    }

    // -------------------------------------------------------------------------
    // CANCEL
    // -------------------------------------------------------------------------

    @Transactional
    public ReservationResponse cancelReservation(UUID id, String userId, String role, String authHeader) {
        Reservation reservation = findByIdOrThrow(id);
        validateOwnershipOrAdmin(reservation, userId, role);

        if (reservation.getStatus() == ReservationStatus.CANCELLED) {
            throw new IllegalArgumentException("Reservation is already cancelled");
        }
        if (reservation.getStatus() == ReservationStatus.COMPLETED) {
            throw new IllegalArgumentException("Cannot cancel a completed reservation");
        }

        reservation.setStatus(ReservationStatus.CANCELLED);
        reservation = reservationRepository.save(reservation);

        log.info("Reservation {} cancelled by user {}", id, userId);

        String spaceName = fetchSpaceNameSafe(reservation.getSpaceId(), authHeader);
        return toResponse(reservation, spaceName);
    }

    // -------------------------------------------------------------------------
    // COMPLETE (admin only)
    // -------------------------------------------------------------------------

    @Transactional
    public ReservationResponse completeReservation(UUID id, String authHeader) {
        Reservation reservation = findByIdOrThrow(id);

        if (reservation.getStatus() != ReservationStatus.CONFIRMED) {
            throw new IllegalArgumentException(
                    "Only CONFIRMED reservations can be marked as COMPLETED. Current status: "
                    + reservation.getStatus());
        }

        reservation.setStatus(ReservationStatus.COMPLETED);
        reservation = reservationRepository.save(reservation);

        log.info("Reservation {} marked as COMPLETED", id);

        String spaceName = fetchSpaceNameSafe(reservation.getSpaceId(), authHeader);
        return toResponse(reservation, spaceName);
    }

    // -------------------------------------------------------------------------
    // AVAILABILITY CHECK
    // -------------------------------------------------------------------------

    @Transactional(readOnly = true)
    public boolean checkAvailability(String spaceId, LocalDateTime startTime, LocalDateTime endTime) {
        if (endTime.isBefore(startTime) || endTime.isEqual(startTime)) {
            throw new IllegalArgumentException("End time must be after start time");
        }
        List<Reservation> overlaps = reservationRepository.findOverlappingReservations(
                spaceId, startTime, endTime);
        return overlaps.isEmpty();
    }

    // -------------------------------------------------------------------------
    // Private helpers
    // -------------------------------------------------------------------------

    private Reservation findByIdOrThrow(UUID id) {
        return reservationRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Reservation", "id", id));
    }

    private void validateOwnershipOrAdmin(Reservation reservation, String userId, String role) {
        boolean isAdmin = "ADMIN".equalsIgnoreCase(role);
        if (!isAdmin && !reservation.getUserId().equals(userId)) {
            throw new AccessDeniedException("You do not have permission to access this reservation");
        }
    }

    private SpaceServiceClient.SpaceDto fetchSpace(String spaceId, String authHeader) {
        try {
            return spaceServiceClient.getSpace(spaceId, authHeader);
        } catch (FeignException.NotFound e) {
            throw new ResourceNotFoundException("Space", "id", spaceId);
        } catch (FeignException e) {
            log.error("Error fetching space {}: {}", spaceId, e.getMessage());
            throw e;
        }
    }

    private String fetchSpaceNameSafe(String spaceId, String authHeader) {
        try {
            SpaceServiceClient.SpaceDto space = spaceServiceClient.getSpace(spaceId, authHeader);
            return space.getName();
        } catch (Exception e) {
            log.warn("Could not fetch space name for {}: {}", spaceId, e.getMessage());
            return spaceId; // Fallback to ID if service is unavailable
        }
    }

    private ReservationResponse toResponseWithSpaceName(Reservation reservation, String authHeader) {
        String spaceName = fetchSpaceNameSafe(reservation.getSpaceId(), authHeader);
        return toResponse(reservation, spaceName);
    }

    private ReservationResponse toResponse(Reservation reservation, String spaceName) {
        return ReservationResponse.builder()
                .id(reservation.getId())
                .userId(reservation.getUserId())
                .spaceId(reservation.getSpaceId())
                .spaceName(spaceName)
                .startTime(reservation.getStartTime())
                .endTime(reservation.getEndTime())
                .status(reservation.getStatus())
                .totalCost(reservation.getTotalCost())
                .notes(reservation.getNotes())
                .createdAt(reservation.getCreatedAt())
                .updatedAt(reservation.getUpdatedAt())
                .build();
    }

    /**
     * Calculate cost = pricePerHour × hours elapsed (ceiling, 2 decimal places).
     */
    private BigDecimal calculateCost(BigDecimal pricePerHour, LocalDateTime start, LocalDateTime end) {
        if (pricePerHour == null || pricePerHour.compareTo(BigDecimal.ZERO) <= 0) {
            return BigDecimal.ZERO;
        }
        long minutes = Duration.between(start, end).toMinutes();
        // Convert minutes to fractional hours, rounded up to 2 decimal places
        BigDecimal hours = BigDecimal.valueOf(minutes)
                .divide(BigDecimal.valueOf(60), 4, RoundingMode.HALF_UP);
        return pricePerHour.multiply(hours).setScale(2, RoundingMode.HALF_UP);
    }
}
