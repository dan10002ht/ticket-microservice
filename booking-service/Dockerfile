# ============================================
# Stage 1: Dependencies
# ============================================
# Pin exact image digest for security & reproducibility
FROM maven:3.9-eclipse-temurin-17@sha256:f7ae6cafd12a38da76a1a105744522387b32caacba5e23674a652133607d007c AS deps

WORKDIR /app

# Copy pom.xml and download dependencies (layer caching)
COPY pom.xml .
RUN mvn dependency:go-offline -B

# ============================================
# Stage 2: Builder
# ============================================
FROM deps AS builder

# Copy source code
COPY src ./src

# Copy protos if needed
COPY ../shared-lib/protos ./protos 2>/dev/null || true

# Build the application
RUN mvn clean package -DskipTests -B

# ============================================
# Stage 3: Security Scan (CI/CD - use with --target security)
# ============================================
FROM deps AS security

# Run OWASP dependency check (non-blocking)
RUN mvn org.owasp:dependency-check-maven:check || true

# ============================================
# Stage 4: Production (Distroless Java)
# ============================================
# Distroless benefits:
# - No shell, package manager, or unnecessary binaries
# - Smaller attack surface (fewer CVEs)
# - Optimized for Java containers
FROM gcr.io/distroless/java17-debian12:nonroot@sha256:032c7e120a1f5e1667a9647f4b129b313fe0f095aeae46e39d7425f23bc48912 AS production

WORKDIR /app

# Copy jar from builder
COPY --from=builder /app/target/booking-service-*.jar /app/app.jar

# Distroless runs as nonroot by default
USER nonroot:nonroot

# Expose HTTP and gRPC ports
EXPOSE 8084 50058

# JVM tuning for containers
# -XX:+UseContainerSupport: Use container memory limits
# -XX:MaxRAMPercentage: Use 75% of container memory for heap
# -XX:InitialRAMPercentage: Start with 50% of max heap
ENV JAVA_TOOL_OPTIONS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0 -XX:InitialRAMPercentage=50.0"

# Note: Distroless has no shell, health checks must be external
# Example in docker-compose.yml:
#   healthcheck:
#     test: ["CMD-SHELL", "curl -f http://localhost:8084/actuator/health || exit 1"]

CMD ["app.jar"]
