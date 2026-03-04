# ============================================
# Stage 1: Build
# ============================================
FROM maven:3.9-eclipse-temurin-21 AS builder

WORKDIR /app

# Cache Maven dependencies first
COPY pom.xml .
RUN mvn dependency:go-offline -B --no-transfer-progress

# Build the application
COPY src ./src
COPY ../shared-lib/protos /shared-lib/protos 2>/dev/null || true
RUN mvn package -DskipTests -B --no-transfer-progress

# ============================================
# Stage 2: Production
# ============================================
FROM eclipse-temurin:21-jre-alpine AS production

WORKDIR /app

# Non-root user for security
RUN addgroup -S invoicesvc && adduser -S invoicesvc -G invoicesvc

COPY --from=builder /app/target/*.jar /app/app.jar

USER invoicesvc

EXPOSE 8083
EXPOSE 50060

ENV SPRING_PROFILES_ACTIVE=prod
ENV JAVA_OPTS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0 -XX:+UseG1GC"

ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar /app/app.jar"]
