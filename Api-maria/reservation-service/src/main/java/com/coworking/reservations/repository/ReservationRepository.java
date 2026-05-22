package com.coworking.reservations.repository;

import com.coworking.reservations.model.Reservation;
import com.coworking.reservations.model.ReservationStatus;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Repository
public interface ReservationRepository extends JpaRepository<Reservation, UUID> {

    List<Reservation> findByUserId(String userId);

    List<Reservation> findByUserIdOrderByCreatedAtDesc(String userId);

    /**
     * Find overlapping reservations for a space, excluding a given status.
     * Used to detect booking conflicts.
     */
    List<Reservation> findBySpaceIdAndStatusNotAndStartTimeLessThanAndEndTimeGreaterThan(
            String spaceId,
            ReservationStatus excludedStatus,
            LocalDateTime endTime,
            LocalDateTime startTime
    );

    /**
     * Find active (PENDING or CONFIRMED) overlapping reservations for a space.
     */
    @Query("""
            SELECT r FROM Reservation r
            WHERE r.spaceId = :spaceId
              AND r.status IN ('PENDING', 'CONFIRMED')
              AND r.startTime < :endTime
              AND r.endTime > :startTime
            """)
    List<Reservation> findOverlappingReservations(
            @Param("spaceId") String spaceId,
            @Param("startTime") LocalDateTime startTime,
            @Param("endTime") LocalDateTime endTime
    );

    /**
     * Admin: list all reservations with optional status and spaceId filters.
     */
    @Query("""
            SELECT r FROM Reservation r
            WHERE (:status IS NULL OR r.status = :status)
              AND (:spaceId IS NULL OR r.spaceId = :spaceId)
            ORDER BY r.createdAt DESC
            """)
    List<Reservation> findAllWithFilters(
            @Param("status") ReservationStatus status,
            @Param("spaceId") String spaceId
    );

    /**
     * Count reservations for a user.
     */
    long countByUserId(String userId);

    /**
     * Find reservations by status.
     */
    List<Reservation> findByStatus(ReservationStatus status);
}
