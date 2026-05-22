package com.coworking.reservations.dto;

import com.coworking.reservations.model.ReservationStatus;
import lombok.*;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ReservationResponse {

    private UUID id;
    private String userId;
    private String spaceId;
    private String spaceName;
    private LocalDateTime startTime;
    private LocalDateTime endTime;
    private ReservationStatus status;
    private BigDecimal totalCost;
    private String notes;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}
