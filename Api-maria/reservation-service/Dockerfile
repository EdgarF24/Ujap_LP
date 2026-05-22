# ---- Builder Stage ----
FROM maven:3.9.6-eclipse-temurin-17-alpine AS builder

WORKDIR /app

# Copy pom first for dependency caching
COPY pom.xml .

# Download dependencies (cached layer)
RUN mvn dependency:go-offline -B || true

# Copy source and build
COPY src ./src
RUN mvn package -DskipTests -B

# ---- Runtime Stage ----
FROM eclipse-temurin:17-jre-alpine AS runtime

WORKDIR /app

# Create non-root user for security and install curl for healthchecks
RUN apk --no-cache add curl && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the built jar from builder stage
COPY --from=builder /app/target/*.jar app.jar

RUN chown appuser:appgroup app.jar

USER appuser

EXPOSE 8003

ENTRYPOINT ["java", "-XX:+UseContainerSupport", "-XX:MaxRAMPercentage=75.0", "-jar", "app.jar"]
