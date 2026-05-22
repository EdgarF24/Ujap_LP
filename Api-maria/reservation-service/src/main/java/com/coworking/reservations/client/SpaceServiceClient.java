package com.coworking.reservations.client;

import lombok.*;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestParam;

import java.math.BigDecimal;

@FeignClient(name = "space-service", url = "${space.service.url}")
public interface SpaceServiceClient {

    /**
     * Retrieve a space by its ID.
     * Passes the JWT Authorization header from the current request context.
     */
    @GetMapping("/spaces/{id}")
    SpaceDto getSpace(
            @PathVariable("id") String id,
            @RequestHeader("Authorization") String authorizationHeader
    );

    /**
     * Check availability for a space during a time window.
     */
    @GetMapping("/spaces/{id}/availability")
    Boolean checkSpaceAvailability(
            @PathVariable("id") String id,
            @RequestParam("startTime") String startTime,
            @RequestParam("endTime") String endTime,
            @RequestHeader("Authorization") String authorizationHeader
    );

    // -------------------------------------------------------------------------
    // Inner DTO representing a Space as returned by space-service
    // -------------------------------------------------------------------------

    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    class SpaceDto {
        private String id;
        private String name;
        private BigDecimal pricePerHour;
        private Integer capacity;
        private Boolean isActive;
        private String type;
        private String description;
    }
}
