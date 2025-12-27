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
FROM gcr.io/distroless/java17-debian12:nonroot@sha256:032c7e120a1f5e1667a9647f4b129b313fe0f095aeae46e39d7425f23bc48912 AS production

WORKDIR /app

# Copy jar from builder
COPY --from=builder /app/target/payment-service-*.jar /app/app.jar

# Distroless runs as nonroot by default
USER nonroot:nonroot

# Expose HTTP and gRPC ports
EXPOSE 8080 9090

# JVM tuning for containers
ENV JAVA_TOOL_OPTIONS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0 -XX:InitialRAMPercentage=50.0"

# Note: Distroless has no shell, health checks must be external
CMD ["app.jar"]
