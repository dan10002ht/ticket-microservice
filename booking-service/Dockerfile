FROM maven:3.9.5-eclipse-temurin-17 AS build

WORKDIR /app

COPY pom.xml .
RUN mvn dependency:go-offline -B

COPY src ./src
COPY ../shared-lib/protos ./protos
RUN mvn clean package -DskipTests

FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

RUN addgroup -S spring && adduser -S spring -G spring
USER spring:spring

COPY --from=build /app/target/booking-service-1.0.0.jar app.jar

EXPOSE 8084 50058

ENTRYPOINT ["java","-jar","/app/app.jar"]

