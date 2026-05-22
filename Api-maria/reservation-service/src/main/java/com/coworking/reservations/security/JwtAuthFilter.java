package com.coworking.reservations.security;

import io.jsonwebtoken.Claims;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.lang.NonNull;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.List;
import java.util.Map;

@Slf4j
@Component
@RequiredArgsConstructor
public class JwtAuthFilter extends OncePerRequestFilter {

    private final JwtUtil jwtUtil;

    @Override
    protected void doFilterInternal(
            @NonNull HttpServletRequest request,
            @NonNull HttpServletResponse response,
            @NonNull FilterChain filterChain
    ) throws ServletException, IOException {

        try {
            String token = extractToken(request);

            if (StringUtils.hasText(token) && jwtUtil.validateToken(token)) {
                Claims claims = jwtUtil.extractClaims(token);

                String userId = jwtUtil.extractUserId(token);
                String email  = jwtUtil.extractEmail(token);
                String role   = jwtUtil.extractRole(token);

                // Build Spring Security authority from role claim
                String authority = (role != null && !role.startsWith("ROLE_"))
                        ? "ROLE_" + role.toUpperCase()
                        : (role != null ? role.toUpperCase() : "ROLE_USER");

                UsernamePasswordAuthenticationToken authentication =
                        new UsernamePasswordAuthenticationToken(
                                userId,
                                null,
                                List.of(new SimpleGrantedAuthority(authority))
                        );

                // Store additional claims as details so controllers can access them
                authentication.setDetails(Map.of(
                        "userId", userId != null ? userId : "",
                        "email",  email  != null ? email  : "",
                        "role",   role   != null ? role   : ""
                ));

                SecurityContextHolder.getContext().setAuthentication(authentication);
                log.debug("Authenticated user '{}' with role '{}'", email, role);
            }
        } catch (Exception ex) {
            log.error("Cannot set user authentication: {}", ex.getMessage());
        }

        filterChain.doFilter(request, response);
    }

    /**
     * Pull the raw JWT from the Authorization: Bearer <token> header.
     */
    private String extractToken(HttpServletRequest request) {
        String bearerToken = request.getHeader("Authorization");
        if (StringUtils.hasText(bearerToken) && bearerToken.startsWith("Bearer ")) {
            return bearerToken.substring(7);
        }
        return null;
    }
}
