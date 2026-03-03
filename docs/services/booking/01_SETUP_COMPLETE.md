# ✅ Phase 1 - Project Setup (Planned)

This document will track booking-service setup milestones once implementation starts.

## 🧱 Planned Deliverables
- Spring Boot 3.2 baseline (Java 17, Maven wrapper)
- Standardized folders (config, controller, entity, repository, service, saga)
- Application entrypoint + layered logging config
- Dockerfile + docker-compose instructions
- Profiles (`dev`, `prod`) with shared config overrides

## 📋 Checklist (to be updated)
- [ ] Initialize `booking-service` module (copy structure from payment-service)
- [ ] Configure Maven dependencies (Spring Web, Data JPA, Validation, Kafka, Redis, gRPC)
- [ ] Add application.yml with placeholders for Postgres/Redis/Kafka
- [ ] Provide local env instructions (`mvn spring-boot:run -Dspring-boot.run.profiles=dev`)
- [ ] Confirm bootstrapping via health endpoint (TODO)

## 🧰 Tools & Dependencies
- Java 17, Maven 3.9+
- Spring Boot 3.2.x
- Flyway, Kafka, Redis, gRPC stubs

## 📝 Notes
- Follow payment-service conventions for package naming (`com.ticketing.booking`)
- Reuse shared base classes (logging, exception handling, `BaseEntity`)
- Add TODO markers where implementation is pending

_Last updated: Planning stage (2024)_
