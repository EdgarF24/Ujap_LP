package com.coworking.reservations.security;

import io.jsonwebtoken.*;
import io.jsonwebtoken.security.Keys;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;

@Slf4j
@Component
public class JwtUtil {

    private final SecretKey signingKey;

    public JwtUtil(@Value("${jwt.secret}") String jwtSecret) {
        byte[] keyBytes = jwtSecret.getBytes(StandardCharsets.UTF_8);
        this.signingKey = Keys.hmacShaKeyFor(keyBytes);
    }

    /**
     * Validate a JWT token — returns true only if the token is well-formed,
     * properly signed and not expired.
     */
    public boolean validateToken(String token) {
        try {
            Jwts.parser()
                    .verifyWith(signingKey)
                    .build()
                    .parseSignedClaims(token);
            return true;
        } catch (ExpiredJwtException e) {
            log.warn("JWT token is expired: {}", e.getMessage());
        } catch (UnsupportedJwtException e) {
            log.warn("JWT token is unsupported: {}", e.getMessage());
        } catch (MalformedJwtException e) {
            log.warn("JWT token is malformed: {}", e.getMessage());
        } catch (SecurityException e) {
            log.warn("JWT signature validation failed: {}", e.getMessage());
        } catch (IllegalArgumentException e) {
            log.warn("JWT claims string is empty: {}", e.getMessage());
        }
        return false;
    }

    /**
     * Parse and return all claims from a valid JWT token.
     */
    public Claims extractClaims(String token) {
        return Jwts.parser()
                .verifyWith(signingKey)
                .build()
                .parseSignedClaims(token)
                .getPayload();
    }

    /**
     * Extract the user_id claim.
     */
    public String extractUserId(String token) {
        Claims claims = extractClaims(token);
        // Support both "user_id" custom claim and standard "sub"
        Object userId = claims.get("user_id");
        if (userId != null) {
            return userId.toString();
        }
        return claims.getSubject();
    }

    /**
     * Extract the email claim.
     */
    public String extractEmail(String token) {
        return (String) extractClaims(token).get("email");
    }

    /**
     * Extract the role claim.
     */
    public String extractRole(String token) {
        return (String) extractClaims(token).get("role");
    }
}
