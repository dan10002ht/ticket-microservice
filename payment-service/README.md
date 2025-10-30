# 💳 Payment Service

Payment Service là một microservice được phát triển bằng **Java Spring Boot**, quản lý toàn bộ nghiệp vụ thanh toán cho hệ thống ticket booking.

## 📋 Features

- ✅ Multiple payment gateway support (Stripe, PayPal, VNPay, Momo)
- ✅ Payment processing và tracking
- ✅ Refund management
- ✅ Webhook handling
- ✅ Idempotency guarantee
- ✅ Transaction logging
- ✅ Security với JWT
- ✅ Monitoring với Prometheus
- ✅ Comprehensive testing

## 🛠️ Tech Stack

- **Java 17**
- **Spring Boot 3.2.0**
- **Spring Data JPA**
- **PostgreSQL**
- **Flyway** (Database migrations)
- **gRPC** (Inter-service communication)
- **Stripe Java SDK**
- **PayPal REST SDK**
- **Micrometer** (Metrics)
- **TestContainers** (Integration testing)

## 🚀 Getting Started

### Prerequisites

- Java 17 or higher
- Maven 3.8+
- PostgreSQL 15+
- Docker (optional, for local development)

### Environment Variables

Create a `.env` file trong root directory:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=payment_db
DB_USER=postgres
DB_PASSWORD=postgres_password

# Server
SERVER_PORT=8080
GRPC_PORT=9090

# JWT
JWT_SECRET=your-secret-key-change-in-production

# Stripe
STRIPE_API_KEY=sk_test_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
STRIPE_PUBLIC_KEY=pk_test_xxx

# PayPal
PAYPAL_MODE=sandbox
PAYPAL_CLIENT_ID=xxx
PAYPAL_CLIENT_SECRET=xxx

# VNPay
VNPAY_MERCHANT_ID=xxx
VNPAY_HASH_SECRET=xxx

# Momo
MOMO_PARTNER_CODE=xxx
MOMO_ACCESS_KEY=xxx
MOMO_SECRET_KEY=xxx
```

### Database Setup

1. Create database:

```bash
createdb payment_db
```

2. Run migrations (automatic on startup):

```bash
mvn flyway:migrate
```

### Running the Application

#### Development mode:

```bash
mvn spring-boot:run -Dspring-boot.run.profiles=dev
```

#### Production mode:

```bash
mvn clean package
java -jar target/payment-service-1.0.0.jar --spring.profiles.active=prod
```

#### With Docker (via centralized docker-compose):

```bash
# Run from project root
cd scripts
./dev-all.sh
```

## 📚 API Documentation

### REST API

- **Base URL**: `http://localhost:8080`
- **Swagger UI**: `http://localhost:8080/swagger-ui.html`

### gRPC API

- **Port**: `9090`
- **Proto files**: `src/main/proto/payment.proto`

## 🧪 Testing

### Run all tests:

```bash
mvn test
```

### Run integration tests:

```bash
mvn verify -P integration-tests
```

### Run with coverage:

```bash
mvn clean test jacoco:report
```

## 📊 Monitoring

### Health Check

```bash
curl http://localhost:8080/actuator/health
```

### Metrics (Prometheus)

```bash
curl http://localhost:8080/actuator/prometheus
```

### Endpoints

- Health: `/actuator/health`
- Metrics: `/actuator/metrics`
- Prometheus: `/actuator/prometheus`
- Info: `/actuator/info`

## 🏗️ Project Structure

```
payment-service/
├── src/
│   ├── main/
│   │   ├── java/com/ticketing/payment/
│   │   │   ├── controller/        # REST controllers
│   │   │   ├── grpc/             # gRPC services
│   │   │   ├── service/          # Business logic
│   │   │   ├── adapter/          # Payment gateway adapters
│   │   │   ├── repository/       # Data access
│   │   │   ├── entity/           # JPA entities
│   │   │   ├── dto/              # Data transfer objects
│   │   │   ├── config/           # Configuration
│   │   │   ├── security/         # Security components
│   │   │   ├── exception/        # Custom exceptions
│   │   │   └── util/             # Utilities
│   │   └── resources/
│   │       ├── application.yml
│   │       └── db/migration/     # Flyway migrations
│   └── test/                     # Tests
├── pom.xml
├── Dockerfile
└── README.md
```

## 🔐 Security

- JWT authentication
- Webhook signature validation
- Idempotency keys
- Rate limiting
- Input validation
- SQL injection prevention

## 📝 Development Guide

### Adding a New Payment Gateway

1. Create adapter class implementing `PaymentGatewayAdapter`
2. Add configuration properties
3. Implement payment processing logic
4. Add webhook handler
5. Write tests

### Database Migrations

1. Create new migration file: `V{version}__{description}.sql`
2. Place in `src/main/resources/db/migration/`
3. Run: `mvn flyway:migrate`

## 🤝 Contributing

1. Create feature branch
2. Make changes
3. Write tests
4. Submit pull request

## 📄 License

Copyright © 2024 Ticket Booking System

---

**Built with ❤️ using Java Spring Boot**
