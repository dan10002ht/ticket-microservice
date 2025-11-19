# 💳 Payment Service - Documentation

Complete documentation for Payment Service implementation.

---

## 📚 Documentation Index

### Phase 1: Core Setup

1. **[01_SETUP_COMPLETE.md](./01_SETUP_COMPLETE.md)** ✅

   - Project structure setup
   - Maven dependencies configuration
   - Spring Boot application setup
   - Docker integration
   - Development environment setup

2. **[02_DATABASE_SETUP_COMPLETE.md](./02_DATABASE_SETUP_COMPLETE.md)** ✅

   - Flyway migrations (V1-V5)
   - Database schema design
   - Tables, views, functions
   - Indexes and constraints
   - Docker database integration

3. **[03_DATABASE_SCHEMA.md](./03_DATABASE_SCHEMA.md)** ✅
   - Complete database schema reference
   - Table structures and columns
   - Views and functions documentation
   - Performance optimization notes
   - Maintenance procedures

---

## 🎯 Implementation Progress

| Phase                          | Section             | Status         | Documents                                                                                                                                  |
| ------------------------------ | ------------------- | -------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| **Phase 1: Core Setup**        |                     |                |                                                                                                                                            |
|                                | Project Setup       | ✅ Complete    | [01_SETUP_COMPLETE.md](./01_SETUP_COMPLETE.md)                                                                                             |
|                                | Database Setup      | ✅ Complete    | [02_DATABASE_SETUP_COMPLETE.md](./02_DATABASE_SETUP_COMPLETE.md)                                                                           |
|                                | Entity Models       | ✅ Complete    | Source code (`src/main/java/com/ticketing/payment/entity`)                                                                                 |
|                                | Repository Layer    | ✅ Complete    | Source code (`src/main/java/com/ticketing/payment/repository`)                                                                             |
| **Phase 2: Payment Flow**      |                     |                |                                                                                                                                            |
|                                | Service Layer       | ✅ Complete    | Source code (`src/main/java/com/ticketing/payment/service`)                                                                                |
|                                | Stripe Adapter      | 🚧 In Progress | `StripeGatewayAdapter` (`src/main/java/com/ticketing/payment/adapter/stripe`) - authorize + provider refs wired, capture/refund hooking up |
|                                | REST API            | ⏳ Not Started | -                                                                                                                                          |
|                                | gRPC API            | ✅ Complete    | [`shared-lib/protos/payment.proto`](../../shared-lib/protos/payment.proto) + gRPC adapter (`src/main/java/com/ticketing/payment/grpc`)     |
| **Phase 3: Webhook & Refunds** |                     |                |                                                                                                                                            |
|                                | Webhook Handling    | 🚧 In Progress | gRPC `ProcessWebhook` + Stripe handler + gateway `/webhooks/payment/:gateway`                                                              |
|                                | Refund Service      | ✅ Complete    | `PaymentService` / `RefundService` implementations                                                                                         |
|                                | Transaction Logging | 🚧 In Progress | `TransactionLogService` + logging hooks                                                                                                    |

---

## 📖 Quick Links

### Setup & Getting Started

- [Project Setup](./01_SETUP_COMPLETE.md#-completed-tasks)
- [Running Locally](../../payment-service/README.md#-getting-started)
- [Docker Setup](./01_SETUP_COMPLETE.md#-docker-integration)

### Database

- [Schema Overview](./03_DATABASE_SCHEMA.md#-tables-overview)
- [Migrations](./02_DATABASE_SETUP_COMPLETE.md#-flyway-migrations-created)
- [Test Migrations](./02_DATABASE_SETUP_COMPLETE.md#-testing)

### Development

- [Project Structure](./01_SETUP_COMPLETE.md#-completed-tasks)
- [Configuration](../../payment-service/src/main/resources/application.yml)
- [Environment Variables](./01_SETUP_COMPLETE.md#-environment-variables)

---

## 🛠️ Technical Stack

- **Language**: Java 17
- **Framework**: Spring Boot 3.2.0
- **Database**: PostgreSQL 15
- **Migration**: Flyway
- **Build Tool**: Maven
- **Container**: Docker
- **API**: REST + gRPC

---

## 📊 Database Objects

| Type       | Count | Status      |
| ---------- | ----- | ----------- |
| Tables     | 4     | ✅ Complete |
| Views      | 4     | ✅ Complete |
| Functions  | 5     | ✅ Complete |
| Indexes    | 50+   | ✅ Complete |
| Migrations | 5     | ✅ Complete |

See [Database Schema](./03_DATABASE_SCHEMA.md) for details.

---

## 🚀 Next Steps

### Phase 1 - Entity Models

1. Create JPA entity classes
2. Add validation annotations
3. Define entity relationships
4. Create enums for status values
5. Add business logic methods
6. Write entity tests

### Phase 1 - Repository Layer

1. Create repository interfaces
2. Add custom query methods
3. Document repository patterns
4. ✅ Integration tests deferred (see TODO section)

---

## ✅ TODO / Deferred Items

- [ ] Viết integration tests cho `PaymentRepository`, `RefundRepository`, `TransactionLogRepository`, `IdempotencyKeyRepository` (ưu tiên Testcontainers + Flyway để bám sát schema thật)
- [ ] Bổ sung unit tests cho entity logic (validation + business helper methods)
- [ ] Thiết lập test data fixtures và test strategy documentation

---

## 🔜 Immediate Next Steps

1. Hoàn thiện Stripe capture/refund workflow (sử dụng `providerReference`, fallback lookup trong webhook)
2. Bổ sung REST controller hoặc gateway client để truyền `IdempotencyKeyContext`/headers vào service layer
3. Hoàn thiện logging/webhook flow (verify signature tất cả provider, mapping refundId đầy đủ)
4. Sau khi adapter hoạt động ổn định, thêm integration/unit tests cho service + adapter layer

---

## 📝 Related Documentation

- [Main Project Documentation](../../README.md)
- [Payment Service README](../../payment-service/README.md)
- [Payment Service Strategy](../../PAYMENT_SERVICE.md)
- [API Event Creation Checklist](../../gateway/API_EVENT_CREATION_CHECKLIST.md)

---

**Last Updated**: 2024  
**Version**: 1.0.0  
**Status**: Phase 1 - Database Setup Complete ✅
