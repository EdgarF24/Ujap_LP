package com.coworking.reservations.controller;

import com.coworking.reservations.dto.ApiResponse;
import com.coworking.reservations.dto.ReservationRequest;
import com.coworking.reservations.dto.ReservationResponse;
import com.coworking.reservations.service.ReservationService;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;
import java.util.UUID;

@Slf4j
@RestController
@RequestMapping
@RequiredArgsConstructor
public class ReservationController {

    private final ReservationService reservationService;

    // -------------------------------------------------------------------------
    // GET /reservations — Admin: list all reservations
    // -------------------------------------------------------------------------

    @GetMapping("/reservations")
    public ResponseEntity<ApiResponse<List<ReservationResponse>>> getAllReservations(
            HttpServletRequest httpRequest) {

        requireRole("ADMIN");
        String authHeader = httpRequest.getHeader("Authorization");
        List<ReservationResponse> reservations = reservationService.getAllReservations(authHeader);
        return ResponseEntity.ok(ApiResponse.success(reservations));
    }

    // -------------------------------------------------------------------------
    // GET /reservations/my — Auth: my reservations
    // -------------------------------------------------------------------------

    @GetMapping("/reservations/my")
    public ResponseEntity<ApiResponse<List<ReservationResponse>>> getMyReservations(
            HttpServletRequest httpRequest) {

        String userId  = extractUserId();
        String authHeader = httpRequest.getHeader("Authorization");
        List<ReservationResponse> reservations = reservationService.getMyReservations(userId, authHeader);
        return ResponseEntity.ok(ApiResponse.success(reservations));
    }

    // -------------------------------------------------------------------------
    // POST /reservations — Auth: create reservation
    // -------------------------------------------------------------------------

    @PostMapping("/reservations")
    public ResponseEntity<ApiResponse<ReservationResponse>> createReservation(
            @Valid @RequestBody ReservationRequest request,
            HttpServletRequest httpRequest) {

        String userId     = extractUserId();
        String authHeader = httpRequest.getHeader("Authorization");
        ReservationResponse reservation = reservationService.createReservation(request, userId, authHeader);
        return ResponseEntity
                .status(HttpStatus.CREATED)
                .body(ApiResponse.success(reservation, "Reservation created and confirmed successfully"));
    }

    // -------------------------------------------------------------------------
    // GET /reservations/availability — Auth: check availability
    // -------------------------------------------------------------------------

    @GetMapping("/reservations/availability")
    public ResponseEntity<ApiResponse<Boolean>> checkAvailability(
            @RequestParam String spaceId,
            @RequestParam @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime startTime,
            @RequestParam @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime endTime) {

        boolean available = reservationService.checkAvailability(spaceId, startTime, endTime);
        String message = available ? "Space is available" : "Space is not available for the requested time";
        return ResponseEntity.ok(ApiResponse.success(available, message));
    }

    // -------------------------------------------------------------------------
    // GET /reservations/{id} — Auth: get specific reservation
    // -------------------------------------------------------------------------

    @GetMapping("/reservations/{id}")
    public ResponseEntity<ApiResponse<ReservationResponse>> getReservation(
            @PathVariable UUID id,
            HttpServletRequest httpRequest) {

        String userId     = extractUserId();
        String role       = extractRole();
        String authHeader = httpRequest.getHeader("Authorization");
        ReservationResponse reservation = reservationService.getReservation(id, userId, role, authHeader);
        return ResponseEntity.ok(ApiResponse.success(reservation));
    }

    // -------------------------------------------------------------------------
    // PUT /reservations/{id}/cancel — Auth: cancel reservation
    // -------------------------------------------------------------------------

    @PutMapping("/reservations/{id}/cancel")
    public ResponseEntity<ApiResponse<ReservationResponse>> cancelReservation(
            @PathVariable UUID id,
            HttpServletRequest httpRequest) {

        String userId     = extractUserId();
        String role       = extractRole();
        String authHeader = httpRequest.getHeader("Authorization");
        ReservationResponse reservation = reservationService.cancelReservation(id, userId, role, authHeader);
        return ResponseEntity.ok(ApiResponse.success(reservation, "Reservation cancelled successfully"));
    }

    // -------------------------------------------------------------------------
    // PUT /reservations/{id}/complete — Admin: complete reservation
    // -------------------------------------------------------------------------

    @PutMapping("/reservations/{id}/complete")
    public ResponseEntity<ApiResponse<ReservationResponse>> completeReservation(
            @PathVariable UUID id,
            HttpServletRequest httpRequest) {

        requireRole("ADMIN");
        String authHeader = httpRequest.getHeader("Authorization");
        ReservationResponse reservation = reservationService.completeReservation(id, authHeader);
        return ResponseEntity.ok(ApiResponse.success(reservation, "Reservation marked as completed"));
    }

    // -------------------------------------------------------------------------
    // GET /health — Public
    // -------------------------------------------------------------------------

    @GetMapping("/health")
    public ResponseEntity<ApiResponse<Map<String, String>>> health() {
        return ResponseEntity.ok(ApiResponse.success(
                Map.of(
                        "status", "UP",
                        "service", "reservation-service",
                        "version", "1.0.0"
                )
        ));
    }

    // -------------------------------------------------------------------------
    // Private helpers
    // -------------------------------------------------------------------------

    private Authentication getAuthentication() {
        return SecurityContextHolder.getContext().getAuthentication();
    }

    private String extractUserId() {
        Authentication auth = getAuthentication();
        if (auth == null || auth.getPrincipal() == null) {
            throw new AccessDeniedException("Not authenticated");
        }
        return auth.getPrincipal().toString();
    }

    private String extractRole() {
        Authentication auth = getAuthentication();
        if (auth == null || auth.getAuthorities() == null || auth.getAuthorities().isEmpty()) {
            return "USER";
        }
        String authority = auth.getAuthorities().iterator().next().getAuthority();
        // Strip "ROLE_" prefix if present
        return authority.startsWith("ROLE_") ? authority.substring(5) : authority;
    }

    @SuppressWarnings("unchecked")
    private String extractEmail() {
        Authentication auth = getAuthentication();
        if (auth instanceof UsernamePasswordAuthenticationToken token) {
            Object details = token.getDetails();
            if (details instanceof Map<?, ?> map) {
                Object email = ((Map<String, Object>) map).get("email");
                return email != null ? email.toString() : "";
            }
        }
        return "";
    }

    /**
     * Enforce that the current authenticated user has the given role.
     * Throws AccessDeniedException otherwise.
     */
    private void requireRole(String requiredRole) {
        String currentRole = extractRole();
        if (!requiredRole.equalsIgnoreCase(currentRole)) {
            throw new AccessDeniedException(
                    "Access denied: role '" + requiredRole + "' required, current role is '" + currentRole + "'");
        }
    }
}
