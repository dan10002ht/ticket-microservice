# 💳 PAYMENT SERVICE - IMPLEMENTATION STRATEGY

## 📋 OVERVIEW

Payment Service là một microservice được implement bằng **Java Spring Boot**, quản lý toàn bộ nghiệp vụ thanh toán cho hệ thống ticket booking. Service này được thiết kế với độ chặt chẽ cao, đảm bảo tính bảo mật, transaction integrity, và khả năng scale.

---

## 🎯 WHY JAVA?

### **✅ Advantages**

- **Strong typing & compile-time safety**: Phát hiện lỗi sớm, giảm runtime errors
- **Enterprise-grade**: Spring Boot ecosystem mạnh mẽ cho financial services
- **Transaction management**: Spring @Transactional cho ACID transactions
- **Security**: Spring Security cho authentication/authorization
- **Thread safety**: Built-in concurrency support cho high-throughput
- **Testing**: JUnit, Mockito, TestContainers cho comprehensive testing
- **Monitoring**: Micrometer, Actuator cho production monitoring
- **Industry standard**: Financial services thường dùng Java

---

## 🏗️ ARCHITECTURE OVERVIEW

```
┌─────────────────────────────────────────────────────────────┐
│                    API GATEWAY                              │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ↓
┌─────────────────────────────────────────────────────────────┐
│                 PAYMENT SERVICE (Java)                      │
├─────────────────────────────────────────────────────────────┤
│  Controller Layer (REST/gRPC)                               │
│  ├── PaymentController                                      │
│  ├── RefundController                                       │
│  └── WebhookController                                      │
├─────────────────────────────────────────────────────────────┤
│  Service Layer (Business Logic)                             │
│  ├── PaymentService                                         │
│  ├── RefundService                                          │
│  ├── PaymentGatewayService                                  │
│  └── ReconciliationService                                  │
├─────────────────────────────────────────────────────────────┤
│  Gateway Adapter Layer                                      │
│  ├── StripeAdapter                                          │
│  ├── PayPalAdapter                                          │
│  ├── VNPayAdapter                                           │
│  └── MomoAdapter                                            │
├─────────────────────────────────────────────────────────────┤
│  Repository Layer (Data Access)                             │
│  ├── PaymentRepository                                      │
│  ├── RefundRepository                                       │
│  └── TransactionLogRepository                               │
├─────────────────────────────────────────────────────────────┤
│  Security Layer                                             │
│  ├── JWT Authentication                                     │
│  ├── Webhook Signature Validation                           │
│  └── Idempotency Guard                                      │
└─────────────────────────────────────────────────────────────┘
                      │
                      ↓
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL Database                            │
│  ├── payments                                               │
│  ├── refunds                                                │
│  ├── transaction_logs                                       │
│  └── idempotency_keys                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 📊 DATABASE DESIGN

### **1. Payments Table**

```sql
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    payment_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

    -- Booking reference
    booking_id VARCHAR(255) NOT NULL,
    ticket_id VARCHAR(255),
    user_id VARCHAR(255) NOT NULL,

    -- Payment details
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    payment_method VARCHAR(50) NOT NULL, -- credit_card, bank_transfer, e_wallet

    -- Status tracking
    status VARCHAR(50) NOT NULL, -- pending, processing, success, failed, cancelled

    -- External references
    external_reference VARCHAR(255), -- Gateway transaction ID
    gateway_provider VARCHAR(50), -- stripe, paypal, vnpay, momo

    -- Idempotency
    idempotency_key VARCHAR(255) UNIQUE,

    -- Correlation
    correlation_id VARCHAR(255) NOT NULL,

    -- Payment metadata
    metadata JSONB,

    -- Timestamps
    initiated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    failed_at TIMESTAMP,

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),

    -- Constraints
    CONSTRAINT check_amount CHECK (amount > 0),
    CONSTRAINT check_status CHECK (status IN ('pending', 'processing', 'success', 'failed', 'cancelled'))
);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_correlation_id ON payments(correlation_id);
CREATE INDEX idx_payments_external_reference ON payments(external_reference);
```

### **2. Refunds Table**

```sql
CREATE TABLE refunds (
    id BIGSERIAL PRIMARY KEY,
    refund_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

    -- Payment reference
    payment_id BIGINT NOT NULL REFERENCES payments(id),

    -- Refund details
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    reason VARCHAR(255),

    -- Status tracking
    status VARCHAR(50) NOT NULL, -- pending, processing, success, failed

    -- External references
    external_reference VARCHAR(255),

    -- Timestamps
    initiated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),

    -- Constraints
    CONSTRAINT check_refund_amount CHECK (amount > 0),
    CONSTRAINT check_refund_status CHECK (status IN ('pending', 'processing', 'success', 'failed'))
);

CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_status ON refunds(status);
```

### **3. Transaction Logs Table**

```sql
CREATE TABLE transaction_logs (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL REFERENCES payments(id),

    -- Event details
    event_type VARCHAR(50) NOT NULL, -- initiated, processing, webhook_received, completed, failed
    event_data JSONB,

    -- External data
    gateway_request JSONB,
    gateway_response JSONB,

    -- Timestamp
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transaction_logs_payment_id ON transaction_logs(payment_id);
CREATE INDEX idx_transaction_logs_event_type ON transaction_logs(event_type);
```

### **4. Idempotency Keys Table**

```sql
CREATE TABLE idempotency_keys (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    payment_id BIGINT REFERENCES payments(id),

    -- Request hash
    request_hash VARCHAR(255) NOT NULL,

    -- Response cache
    response_body TEXT,
    response_status INT,

    -- TTL
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_idempotency_keys_expires_at ON idempotency_keys(expires_at);
```

---

## 🛠️ TECH STACK

### **Core Framework**

- **Spring Boot 3.x**: Application framework
- **Spring Data JPA**: Database access
- **Spring Security**: Authentication/Authorization
- **Spring Cloud**: Microservices patterns

### **Database**

- **PostgreSQL**: Primary database
- **HikariCP**: Connection pooling
- **Flyway**: Database migrations

### **Communication**

- **gRPC (Spring Boot Starter gRPC)**: Inter-service communication
- **REST API**: External API endpoints
- **Spring WebFlux**: Reactive programming (for webhooks)

### **Payment Gateways**

- **Stripe Java SDK**: Credit card processing
- **PayPal REST SDK**: PayPal integration
- **Custom adapters**: VNPay, Momo, ZaloPay

### **Monitoring & Observability**

- **Micrometer**: Metrics collection
- **Spring Boot Actuator**: Health checks, metrics
- **Prometheus**: Metrics storage
- **Grafana**: Metrics visualization
- **SLF4J + Logback**: Logging

### **Testing**

- **JUnit 5**: Unit testing
- **Mockito**: Mocking framework
- **TestContainers**: Integration testing with PostgreSQL
- **WireMock**: API mocking for payment gateways

### **Security**

- **Spring Security**: Authentication/Authorization
- **JWT**: Token-based authentication
- **HMAC**: Webhook signature validation
- **BCrypt**: Password hashing

---

## 📦 PROJECT STRUCTURE

```
payment-service/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/
│   │   │       └── ticketing/
│   │   │           └── payment/
│   │   │               ├── PaymentServiceApplication.java
│   │   │               ├── controller/
│   │   │               │   ├── PaymentController.java
│   │   │               │   ├── RefundController.java
│   │   │               │   └── WebhookController.java
│   │   │               ├── grpc/
│   │   │               │   ├── PaymentGrpcService.java
│   │   │               │   └── RefundGrpcService.java
│   │   │               ├── service/
│   │   │               │   ├── PaymentService.java
│   │   │               │   ├── RefundService.java
│   │   │               │   ├── PaymentGatewayService.java
│   │   │               │   └── ReconciliationService.java
│   │   │               ├── adapter/
│   │   │               │   ├── PaymentGatewayAdapter.java (interface)
│   │   │               │   ├── StripeAdapter.java
│   │   │               │   ├── PayPalAdapter.java
│   │   │               │   ├── VNPayAdapter.java
│   │   │               │   └── MomoAdapter.java
│   │   │               ├── repository/
│   │   │               │   ├── PaymentRepository.java
│   │   │               │   ├── RefundRepository.java
│   │   │               │   ├── TransactionLogRepository.java
│   │   │               │   └── IdempotencyKeyRepository.java
│   │   │               ├── entity/
│   │   │               │   ├── Payment.java
│   │   │               │   ├── Refund.java
│   │   │               │   ├── TransactionLog.java
│   │   │               │   └── IdempotencyKey.java
│   │   │               ├── dto/
│   │   │               │   ├── request/
│   │   │               │   │   ├── InitiatePaymentRequest.java
│   │   │               │   │   ├── InitiateRefundRequest.java
│   │   │               │   │   └── WebhookRequest.java
│   │   │               │   └── response/
│   │   │               │       ├── PaymentResponse.java
│   │   │               │       ├── RefundResponse.java
│   │   │               │       └── PaymentStatusResponse.java
│   │   │               ├── config/
│   │   │               │   ├── DatabaseConfig.java
│   │   │               │   ├── SecurityConfig.java
│   │   │               │   ├── GrpcConfig.java
│   │   │               │   └── PaymentGatewayConfig.java
│   │   │               ├── security/
│   │   │               │   ├── JwtAuthenticationFilter.java
│   │   │               │   ├── WebhookSignatureValidator.java
│   │   │               │   └── IdempotencyGuard.java
│   │   │               ├── exception/
│   │   │               │   ├── PaymentException.java
│   │   │               │   ├── DuplicatePaymentException.java
│   │   │               │   ├── InvalidAmountException.java
│   │   │               │   └── GatewayException.java
│   │   │               └── util/
│   │   │                   ├── CorrelationIdGenerator.java
│   │   │                   └── PaymentValidator.java
│   │   └── resources/
│   │       ├── application.yml
│   │       ├── application-dev.yml
│   │       ├── application-prod.yml
│   │       └── db/
│   │           └── migration/
│   │               ├── V1__create_payments_table.sql
│   │               ├── V2__create_refunds_table.sql
│   │               ├── V3__create_transaction_logs_table.sql
│   │               └── V4__create_idempotency_keys_table.sql
│   └── test/
│       └── java/
│           └── com/
│               └── ticketing/
│                   └── payment/
│                       ├── service/
│                       │   ├── PaymentServiceTest.java
│                       │   └── RefundServiceTest.java
│                       ├── controller/
│                       │   └── PaymentControllerTest.java
│                       └── integration/
│                           └── PaymentIntegrationTest.java
├── pom.xml
├── Dockerfile
└── README.md
```

---

## 🔧 CORE COMPONENTS

### **1. Payment Entity**

```java
@Entity
@Table(name = "payments")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Payment {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(unique = true, nullable = false)
    private UUID paymentId = UUID.randomUUID();

    @Column(nullable = false)
    private String bookingId;

    private String ticketId;

    @Column(nullable = false)
    private String userId;

    @Column(nullable = false, precision = 10, scale = 2)
    private BigDecimal amount;

    @Column(nullable = false, length = 3)
    private String currency = "USD";

    @Column(nullable = false, length = 50)
    @Enumerated(EnumType.STRING)
    private PaymentMethod paymentMethod;

    @Column(nullable = false, length = 50)
    @Enumerated(EnumType.STRING)
    private PaymentStatus status = PaymentStatus.PENDING;

    private String externalReference;

    @Column(length = 50)
    @Enumerated(EnumType.STRING)
    private PaymentGateway gatewayProvider;

    @Column(unique = true)
    private String idempotencyKey;

    @Column(nullable = false)
    private String correlationId;

    @Type(JsonBinaryType.class)
    @Column(columnDefinition = "jsonb")
    private Map<String, Object> metadata;

    @Column(nullable = false)
    private LocalDateTime initiatedAt = LocalDateTime.now();

    private LocalDateTime completedAt;
    private LocalDateTime failedAt;

    @CreationTimestamp
    private LocalDateTime createdAt;

    @UpdateTimestamp
    private LocalDateTime updatedAt;

    private String createdBy;
    private String updatedBy;

    // Business logic methods
    public boolean canBeRefunded() {
        return status == PaymentStatus.SUCCESS && completedAt != null;
    }

    public boolean isPending() {
        return status == PaymentStatus.PENDING;
    }

    public boolean isCompleted() {
        return status == PaymentStatus.SUCCESS;
    }
}
```

### **2. Payment Service**

```java
@Service
@Transactional
@Slf4j
public class PaymentService {

    private final PaymentRepository paymentRepository;
    private final TransactionLogRepository transactionLogRepository;
    private final IdempotencyKeyRepository idempotencyKeyRepository;
    private final PaymentGatewayService gatewayService;

    @Transactional
    public PaymentResponse initiatePayment(InitiatePaymentRequest request) {
        // 1. Validate idempotency
        if (request.getIdempotencyKey() != null) {
            Optional<PaymentResponse> cached = checkIdempotency(request);
            if (cached.isPresent()) {
                return cached.get();
            }
        }

        // 2. Validate payment request
        validatePaymentRequest(request);

        // 3. Create payment record
        Payment payment = createPaymentRecord(request);
        payment = paymentRepository.save(payment);

        // 4. Log transaction
        logTransaction(payment, "INITIATED", request);

        // 5. Call payment gateway
        try {
            PaymentGatewayResponse gatewayResponse =
                gatewayService.processPayment(payment, request);

            payment.setExternalReference(gatewayResponse.getTransactionId());
            payment.setStatus(PaymentStatus.PROCESSING);
            payment = paymentRepository.save(payment);

            logTransaction(payment, "GATEWAY_CALLED", gatewayResponse);

        } catch (GatewayException e) {
            payment.setStatus(PaymentStatus.FAILED);
            payment.setFailedAt(LocalDateTime.now());
            paymentRepository.save(payment);

            logTransaction(payment, "GATEWAY_FAILED", e.getMessage());
            throw e;
        }

        // 6. Cache idempotency response
        PaymentResponse response = toResponse(payment);
        if (request.getIdempotencyKey() != null) {
            cacheIdempotency(request.getIdempotencyKey(), response);
        }

        return response;
    }

    @Transactional
    public void handleWebhookCallback(WebhookRequest webhook) {
        // 1. Validate webhook signature
        validateWebhookSignature(webhook);

        // 2. Find payment by external reference
        Payment payment = paymentRepository
            .findByExternalReference(webhook.getTransactionId())
            .orElseThrow(() -> new PaymentNotFoundException("Payment not found"));

        // 3. Update payment status based on webhook
        switch (webhook.getStatus()) {
            case "SUCCESS":
                payment.setStatus(PaymentStatus.SUCCESS);
                payment.setCompletedAt(LocalDateTime.now());
                break;
            case "FAILED":
                payment.setStatus(PaymentStatus.FAILED);
                payment.setFailedAt(LocalDateTime.now());
                break;
            default:
                payment.setStatus(PaymentStatus.PROCESSING);
        }

        paymentRepository.save(payment);
        logTransaction(payment, "WEBHOOK_RECEIVED", webhook);

        // 4. Publish event (ready for Kafka later)
        publishPaymentCompletedEvent(payment);
    }

    private void validateWebhookSignature(WebhookRequest webhook) {
        String signature = webhook.getSignature();
        String payload = webhook.getPayload();
        String secret = getGatewaySecret(webhook.getGateway());

        String expectedSignature = HmacUtils.hmacSha256Hex(secret, payload);

        if (!signature.equals(expectedSignature)) {
            throw new InvalidWebhookSignatureException("Invalid webhook signature");
        }
    }

    // More methods...
}
```

### **3. Payment Gateway Adapter Interface**

```java
public interface PaymentGatewayAdapter {
    PaymentGatewayResponse processPayment(Payment payment, InitiatePaymentRequest request);
    RefundGatewayResponse processRefund(Payment payment, BigDecimal amount, String reason);
    PaymentStatus getPaymentStatus(String externalReference);
    boolean validateWebhook(String signature, String payload);
}
```

### **4. Stripe Adapter Implementation**

```java
@Component
@Slf4j
public class StripeAdapter implements PaymentGatewayAdapter {

    private final Stripe stripeClient;

    @Value("${payment.stripe.secret-key}")
    private String secretKey;

    @Override
    public PaymentGatewayResponse processPayment(Payment payment, InitiatePaymentRequest request) {
        try {
            Stripe.apiKey = secretKey;

            PaymentIntentCreateParams params = PaymentIntentCreateParams.builder()
                .setAmount(payment.getAmount().multiply(new BigDecimal(100)).longValue())
                .setCurrency(payment.getCurrency().toLowerCase())
                .setPaymentMethod(request.getPaymentMethodToken())
                .setConfirm(true)
                .putMetadata("booking_id", payment.getBookingId())
                .putMetadata("user_id", payment.getUserId())
                .build();

            PaymentIntent intent = PaymentIntent.create(params);

            return PaymentGatewayResponse.builder()
                .transactionId(intent.getId())
                .status(mapStripeStatus(intent.getStatus()))
                .gatewayResponse(intent.toJson())
                .build();

        } catch (StripeException e) {
            log.error("Stripe payment failed", e);
            throw new GatewayException("Stripe payment failed: " + e.getMessage(), e);
        }
    }

    @Override
    public boolean validateWebhook(String signature, String payload) {
        try {
            Event event = Webhook.constructEvent(
                payload, signature, webhookSecret
            );
            return true;
        } catch (SignatureVerificationException e) {
            return false;
        }
    }

    private PaymentStatus mapStripeStatus(String stripeStatus) {
        return switch (stripeStatus) {
            case "succeeded" -> PaymentStatus.SUCCESS;
            case "processing" -> PaymentStatus.PROCESSING;
            case "requires_payment_method", "canceled" -> PaymentStatus.FAILED;
            default -> PaymentStatus.PENDING;
        };
    }
}
```

---

## 🔐 SECURITY FEATURES

### **1. Idempotency Guard**

```java
@Component
public class IdempotencyGuard {

    private final IdempotencyKeyRepository repository;

    @Transactional
    public Optional<PaymentResponse> checkIdempotency(String key, String requestHash) {
        return repository.findByKeyAndNotExpired(key, LocalDateTime.now())
            .filter(ik -> ik.getRequestHash().equals(requestHash))
            .map(this::buildResponseFromCache);
    }

    @Transactional
    public void cacheResponse(String key, String requestHash, PaymentResponse response) {
        IdempotencyKey idempotencyKey = IdempotencyKey.builder()
            .key(key)
            .requestHash(requestHash)
            .responseBody(serializeResponse(response))
            .responseStatus(200)
            .expiresAt(LocalDateTime.now().plusHours(24))
            .build();

        repository.save(idempotencyKey);
    }
}
```

### **2. Webhook Signature Validator**

```java
@Component
public class WebhookSignatureValidator {

    public boolean validate(String gateway, String signature, String payload, String secret) {
        return switch (gateway.toLowerCase()) {
            case "stripe" -> validateStripeSignature(signature, payload, secret);
            case "paypal" -> validatePayPalSignature(signature, payload, secret);
            case "vnpay" -> validateVNPaySignature(signature, payload, secret);
            default -> throw new UnsupportedGatewayException("Unknown gateway: " + gateway);
        };
    }

    private boolean validateStripeSignature(String signature, String payload, String secret) {
        String expectedSignature = HmacUtils.hmacSha256Hex(secret, payload);
        return MessageDigest.isEqual(
            signature.getBytes(StandardCharsets.UTF_8),
            expectedSignature.getBytes(StandardCharsets.UTF_8)
        );
    }
}
```

---

## 🚀 API ENDPOINTS

### **gRPC Service**

```protobuf
service PaymentService {
  rpc InitiatePayment(InitiatePaymentRequest) returns (PaymentResponse);
  rpc GetPaymentStatus(GetPaymentStatusRequest) returns (PaymentStatusResponse);
  rpc InitiateRefund(InitiateRefundRequest) returns (RefundResponse);
  rpc GetPaymentHistory(GetPaymentHistoryRequest) returns (PaymentHistoryResponse);
}

message InitiatePaymentRequest {
  string booking_id = 1;
  string user_id = 2;
  double amount = 3;
  string currency = 4;
  string payment_method = 5;
  string payment_method_token = 6;
  string idempotency_key = 7;
  string correlation_id = 8;
}

message PaymentResponse {
  string payment_id = 1;
  string status = 2;
  string external_reference = 3;
  string message = 4;
}
```

### **REST API**

```
POST   /api/payments                    # Initiate payment
GET    /api/payments/{paymentId}        # Get payment status
POST   /api/payments/{paymentId}/refund # Initiate refund
GET    /api/payments/user/{userId}      # Get user payment history
POST   /api/webhooks/{gateway}          # Webhook callback
GET    /health                          # Health check
```

---

## 📊 CONFIGURATION

### **application.yml**

```yaml
spring:
  application:
    name: payment-service

  datasource:
    url: jdbc:postgresql://localhost:5432/payment_db
    username: ${DB_USER}
    password: ${DB_PASSWORD}
    hikari:
      maximum-pool-size: 10
      minimum-idle: 5

  jpa:
    hibernate:
      ddl-auto: validate
    show-sql: false
    properties:
      hibernate:
        dialect: org.hibernate.dialect.PostgreSQLDialect
        format_sql: true

  flyway:
    enabled: true
    locations: classpath:db/migration

grpc:
  server:
    port: 9090

server:
  port: 8080

payment:
  stripe:
    api-key: ${STRIPE_API_KEY}
    webhook-secret: ${STRIPE_WEBHOOK_SECRET}
  paypal:
    client-id: ${PAYPAL_CLIENT_ID}
    client-secret: ${PAYPAL_CLIENT_SECRET}
  vnpay:
    merchant-id: ${VNPAY_MERCHANT_ID}
    hash-secret: ${VNPAY_HASH_SECRET}

management:
  endpoints:
    web:
      exposure:
        include: health,metrics,prometheus
  metrics:
    export:
      prometheus:
        enabled: true
```

---

## 🧪 TESTING STRATEGY

### **1. Unit Tests**

```java
@ExtendWith(MockitoExtension.class)
class PaymentServiceTest {

    @Mock
    private PaymentRepository paymentRepository;

    @Mock
    private PaymentGatewayService gatewayService;

    @InjectMocks
    private PaymentService paymentService;

    @Test
    void shouldInitiatePaymentSuccessfully() {
        // Given
        InitiatePaymentRequest request = createTestRequest();
        Payment payment = createTestPayment();

        when(paymentRepository.save(any())).thenReturn(payment);
        when(gatewayService.processPayment(any(), any()))
            .thenReturn(createSuccessGatewayResponse());

        // When
        PaymentResponse response = paymentService.initiatePayment(request);

        // Then
        assertNotNull(response);
        assertEquals(PaymentStatus.PROCESSING, response.getStatus());
        verify(paymentRepository, times(2)).save(any());
    }
}
```

### **2. Integration Tests**

```java
@SpringBootTest
@Testcontainers
class PaymentIntegrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15")
        .withDatabaseName("payment_test")
        .withUsername("test")
        .withPassword("test");

    @Autowired
    private PaymentService paymentService;

    @Test
    void shouldProcessPaymentEndToEnd() {
        // Test full payment flow with real database
    }
}
```

---

## 📈 MONITORING & OBSERVABILITY

### **Metrics**

```java
@Component
public class PaymentMetrics {

    private final Counter paymentInitiated;
    private final Counter paymentSuccess;
    private final Counter paymentFailed;
    private final Timer paymentDuration;

    public PaymentMetrics(MeterRegistry registry) {
        this.paymentInitiated = Counter.builder("payment.initiated")
            .description("Total payments initiated")
            .register(registry);

        this.paymentSuccess = Counter.builder("payment.success")
            .description("Total successful payments")
            .register(registry);

        this.paymentFailed = Counter.builder("payment.failed")
            .description("Total failed payments")
            .register(registry);

        this.paymentDuration = Timer.builder("payment.duration")
            .description("Payment processing duration")
            .register(registry);
    }
}
```

---

## 🎯 IMPLEMENTATION ROADMAP

### **Phase 1: Core Setup (Week 1)**

1. Project setup với Spring Boot
2. Database schema design và migration
3. Basic entity models
4. Repository layer
5. Basic service layer

### **Phase 2: Payment Flow (Week 2)**

6. Payment service implementation
7. Stripe adapter (primary gateway)
8. REST API endpoints
9. gRPC endpoints
10. Idempotency guard

### **Phase 3: Webhook & Refunds (Week 3)**

11. Webhook handling
12. Signature validation
13. Refund service
14. Transaction logging
15. Error handling

### **Phase 4: Testing & Security (Week 4)**

16. Unit tests
17. Integration tests
18. Security hardening
19. Performance testing
20. Documentation

### **Phase 5: Additional Gateways (Week 5)**

21. PayPal adapter
22. VNPay adapter
23. Momo adapter
24. Gateway abstraction refinement

### **Phase 6: Production Ready (Week 6)**

25. Monitoring & metrics
26. Logging enhancement
27. Load testing
28. Production deployment
29. Reconciliation service
30. Documentation completion

---

## 🔮 FUTURE ENHANCEMENTS

### **Phase 7: Kafka Integration**

- Event publishing for payment events
- Async payment processing
- Event sourcing for audit trail

### **Phase 8: Advanced Features**

- Partial payments
- Split payments
- Recurring payments
- Payment methods management
- 3D Secure support
- Multi-currency support
- Payment analytics dashboard

---

## 🎯 SUCCESS CRITERIA

### **Functional**

- ✅ Process payments reliably
- ✅ Handle refunds correctly
- ✅ Webhook processing
- ✅ Multiple gateway support
- ✅ Idempotency guarantee

### **Non-Functional**

- ✅ 99.9% uptime
- ✅ < 2s payment initiation
- ✅ < 5s webhook processing
- ✅ Zero duplicate charges
- ✅ Complete audit trail

---

## 🚀 CONCLUSION

Payment Service với Java Spring Boot sẽ cung cấp:

- **Strong type safety** cho financial transactions
- **Enterprise-grade reliability** với Spring ecosystem
- **Comprehensive security** với Spring Security
- **Scalable architecture** sẵn sàng cho Kafka
- **Production-ready** với monitoring và testing

**Java là lựa chọn hoàn hảo cho Payment Service với yêu cầu cao về độ chặt chẽ và bảo mật!** 💳

---

## ✅ IMPLEMENTATION CHECKLIST

### **📦 Phase 1: Core Setup (Week 1)**

#### **Project Setup**

- [x] Initialize Spring Boot 3.x project với Maven/Gradle
- [x] Add dependencies (Spring Data JPA, PostgreSQL, Flyway, gRPC)
- [x] Configure project structure theo package convention
- [x] Setup Git repository và .gitignore
- [x] Create README.md và documentation

#### **Database Setup**

- [x] Create PostgreSQL database `payment_db`
- [x] Write Flyway migration V1: `payments` table
- [x] Write Flyway migration V2: `refunds` table
- [x] Write Flyway migration V3: `transaction_logs` table
- [x] Write Flyway migration V4: `idempotency_keys` table
- [x] Test migrations locally
- [x] Create database indexes

#### **Entity Models**

- [x] Create `Payment` entity với JPA annotations
- [x] Create `Refund` entity với relationships
- [x] Create `TransactionLog` entity
- [x] Create `IdempotencyKey` entity
- [x] Add validation annotations (@NotNull, @Size, etc.)
- [x] Add business logic methods (canBeRefunded, isPending, etc.)
- [x] Create enums (PaymentStatus, PaymentMethod, PaymentGateway)

#### **Repository Layer**

- [ ] Create `PaymentRepository` extends JpaRepository
- [ ] Add custom query methods (findByBookingId, findByUserId, etc.)
- [ ] Create `RefundRepository`
- [ ] Create `TransactionLogRepository`
- [ ] Create `IdempotencyKeyRepository`
- [ ] Write repository unit tests

---

### **💳 Phase 2: Payment Flow (Week 2)**

#### **DTO Layer**

- [ ] Create `InitiatePaymentRequest` DTO
- [ ] Create `PaymentResponse` DTO
- [ ] Create `InitiateRefundRequest` DTO
- [ ] Create `RefundResponse` DTO
- [ ] Create `PaymentStatusResponse` DTO
- [ ] Create `WebhookRequest` DTO
- [ ] Add validation annotations

#### **Service Layer**

- [ ] Create `PaymentService` với @Service annotation
- [ ] Implement `initiatePayment()` method
- [ ] Implement `getPaymentStatus()` method
- [ ] Implement `getPaymentHistory()` method
- [ ] Add transaction management (@Transactional)
- [ ] Implement validation logic
- [ ] Add correlation ID generation
- [ ] Implement transaction logging

#### **Gateway Adapter**

- [ ] Create `PaymentGatewayAdapter` interface
- [ ] Create `StripeAdapter` implementation
- [ ] Implement Stripe payment processing
- [ ] Add Stripe error handling
- [ ] Test Stripe integration với sandbox
- [ ] Create `PaymentGatewayService` orchestrator
- [ ] Add gateway selection logic

#### **REST API**

- [ ] Create `PaymentController` với @RestController
- [ ] Add `POST /api/payments` endpoint
- [ ] Add `GET /api/payments/{id}` endpoint
- [ ] Add `GET /api/payments/user/{userId}` endpoint
- [ ] Add request validation
- [ ] Add error handling (@ExceptionHandler)
- [ ] Add API documentation (Swagger/OpenAPI)

#### **gRPC API**

- [ ] Define payment.proto service definition
- [ ] Generate Java classes từ proto
- [ ] Create `PaymentGrpcService` implementation
- [ ] Implement gRPC endpoints
- [ ] Add gRPC interceptors (logging, auth)
- [ ] Test gRPC với gRPC client

---

### **🔐 Phase 3: Webhook & Refunds (Week 3)**

#### **Webhook Handling**

- [ ] Create `WebhookController`
- [ ] Add `POST /api/webhooks/{gateway}` endpoint
- [ ] Implement webhook signature validation
- [ ] Create `WebhookSignatureValidator` component
- [ ] Implement Stripe webhook validation
- [ ] Handle webhook events (payment.succeeded, payment.failed)
- [ ] Update payment status từ webhook
- [ ] Add webhook retry mechanism
- [ ] Test webhook với mock gateway

#### **Refund Service**

- [ ] Create `RefundService` với @Service annotation
- [ ] Implement `initiateRefund()` method
- [ ] Validate refund eligibility
- [ ] Call gateway refund API
- [ ] Create refund record
- [ ] Update payment status
- [ ] Add refund transaction logging
- [ ] Implement partial refund support

#### **Security Layer**

- [ ] Create `IdempotencyGuard` component
- [ ] Implement idempotency checking
- [ ] Implement response caching
- [ ] Add idempotency cleanup job
- [ ] Create `JwtAuthenticationFilter`
- [ ] Add Spring Security configuration
- [ ] Implement webhook signature validation
- [ ] Add rate limiting for API endpoints

#### **Transaction Logging**

- [ ] Implement `logTransaction()` method
- [ ] Log payment initiation
- [ ] Log gateway requests/responses
- [ ] Log webhook events
- [ ] Log status changes
- [ ] Add structured logging với correlation IDs
- [ ] Implement log rotation

---

### **🧪 Phase 4: Testing & Security (Week 4)**

#### **Unit Tests**

- [ ] Write `PaymentServiceTest`
- [ ] Test payment initiation flow
- [ ] Test idempotency guard
- [ ] Test validation logic
- [ ] Write `RefundServiceTest`
- [ ] Test refund flow
- [ ] Write `StripeAdapterTest`
- [ ] Test gateway error handling
- [ ] Achieve 80%+ code coverage

#### **Integration Tests**

- [ ] Setup TestContainers for PostgreSQL
- [ ] Write `PaymentIntegrationTest`
- [ ] Test end-to-end payment flow
- [ ] Test webhook processing
- [ ] Test refund flow
- [ ] Test concurrent requests
- [ ] Test idempotency with real DB
- [ ] Test transaction rollback

#### **Security Testing**

- [ ] Test webhook signature validation
- [ ] Test idempotency key expiration
- [ ] Test duplicate payment prevention
- [ ] Test authorization checks
- [ ] Perform security audit
- [ ] Test SQL injection prevention
- [ ] Test XSS prevention

#### **Performance Testing**

- [ ] Load test payment endpoints
- [ ] Test concurrent payment processing
- [ ] Test database connection pooling
- [ ] Optimize slow queries
- [ ] Test gateway timeout handling
- [ ] Benchmark response times
- [ ] Identify bottlenecks

---

### **🌐 Phase 5: Additional Gateways (Week 5)**

#### **PayPal Integration**

- [ ] Create `PayPalAdapter` implementation
- [ ] Integrate PayPal REST SDK
- [ ] Implement PayPal payment processing
- [ ] Implement PayPal webhook validation
- [ ] Test PayPal integration với sandbox
- [ ] Add PayPal error handling

#### **VNPay Integration**

- [ ] Create `VNPayAdapter` implementation
- [ ] Implement VNPay payment URL generation
- [ ] Implement VNPay callback handling
- [ ] Implement VNPay signature validation
- [ ] Test VNPay integration
- [ ] Add VNPay error handling

#### **Momo Integration**

- [ ] Create `MomoAdapter` implementation
- [ ] Integrate Momo SDK/API
- [ ] Implement Momo payment processing
- [ ] Implement Momo IPN handling
- [ ] Test Momo integration
- [ ] Add Momo error handling

#### **Gateway Abstraction**

- [ ] Refine `PaymentGatewayAdapter` interface
- [ ] Implement gateway factory pattern
- [ ] Add gateway configuration management
- [ ] Implement gateway fallback logic
- [ ] Add gateway health checks
- [ ] Implement gateway selection strategy

---

### **🚀 Phase 6: Production Ready (Week 6)**

#### **Monitoring & Observability**

- [ ] Add Micrometer metrics
- [ ] Create custom metrics (payment.initiated, payment.success, etc.)
- [ ] Configure Prometheus endpoint
- [ ] Create Grafana dashboards
- [ ] Add Spring Boot Actuator
- [ ] Configure health checks
- [ ] Add liveness probe
- [ ] Add readiness probe

#### **Logging Enhancement**

- [ ] Configure Logback
- [ ] Add structured logging (JSON format)
- [ ] Add correlation ID to all logs
- [ ] Configure log levels per environment
- [ ] Add sensitive data masking
- [ ] Setup log aggregation (ELK/Loki)
- [ ] Add error alerting

#### **Configuration Management**

- [ ] Create application-dev.yml
- [ ] Create application-staging.yml
- [ ] Create application-prod.yml
- [ ] Externalize sensitive configs (Vault/Secrets Manager)
- [ ] Add configuration validation
- [ ] Document all configuration properties

#### **Documentation**

- [ ] Write API documentation (Swagger/OpenAPI)
- [ ] Write gRPC service documentation
- [ ] Create architecture diagram
- [ ] Write deployment guide
- [ ] Write operational runbook
- [ ] Document error codes
- [ ] Create troubleshooting guide

#### **Deployment**

- [ ] Create Dockerfile
- [ ] Create docker-compose.yml for local development
- [ ] Write Kubernetes manifests (if applicable)
- [ ] Configure CI/CD pipeline
- [ ] Setup staging environment
- [ ] Setup production environment
- [ ] Perform smoke tests
- [ ] Load testing in staging

#### **Reconciliation Service**

- [ ] Create `ReconciliationService`
- [ ] Implement daily reconciliation job
- [ ] Compare internal records với gateway
- [ ] Generate reconciliation reports
- [ ] Alert on discrepancies
- [ ] Implement dispute resolution workflow

---

### **🔮 Phase 7: Kafka Integration (Future)**

#### **Event Publishing**

- [ ] Add Kafka dependencies
- [ ] Create Kafka configuration
- [ ] Define payment event schemas
- [ ] Create `PaymentEventPublisher`
- [ ] Publish `payment.initiated` event
- [ ] Publish `payment.completed` event
- [ ] Publish `payment.failed` event
- [ ] Publish `refund.completed` event

#### **Outbox Pattern**

- [ ] Create `outbox` table
- [ ] Implement outbox write on payment completion
- [ ] Create outbox dispatcher worker
- [ ] Implement at-least-once delivery
- [ ] Add retry mechanism
- [ ] Implement dead letter queue
- [ ] Monitor outbox lag

#### **Event Sourcing**

- [ ] Design event store schema
- [ ] Implement event append
- [ ] Implement event replay
- [ ] Create event projections
- [ ] Add snapshots for performance
- [ ] Implement event versioning

---

### **🎨 Phase 8: Advanced Features (Future)**

#### **Payment Methods Management**

- [ ] Create `payment_methods` table
- [ ] Implement save payment method
- [ ] Implement tokenization
- [ ] Add default payment method
- [ ] Implement payment method deletion
- [ ] Add PCI compliance measures

#### **Partial & Split Payments**

- [ ] Design partial payment schema
- [ ] Implement partial payment processing
- [ ] Track payment installments
- [ ] Implement split payment (multiple sources)
- [ ] Add payment plan support

#### **Advanced Security**

- [ ] Implement 3D Secure
- [ ] Add fraud detection rules
- [ ] Implement risk scoring
- [ ] Add velocity checks
- [ ] Implement geolocation checks
- [ ] Add device fingerprinting

#### **Analytics Dashboard**

- [ ] Create payment analytics schema
- [ ] Implement revenue tracking
- [ ] Create conversion funnel
- [ ] Add payment method statistics
- [ ] Implement cohort analysis
- [ ] Create executive dashboard

---

## 📊 IMPLEMENTATION PROGRESS TRACKING

### **Overall Progress**

| Phase                        | Tasks   | Completed | Progress |
| ---------------------------- | ------- | --------- | -------- |
| Phase 1: Core Setup          | 21      | 0         | 0%       |
| Phase 2: Payment Flow        | 27      | 0         | 0%       |
| Phase 3: Webhook & Refunds   | 22      | 0         | 0%       |
| Phase 4: Testing & Security  | 21      | 0         | 0%       |
| Phase 5: Additional Gateways | 18      | 0         | 0%       |
| Phase 6: Production Ready    | 31      | 0         | 0%       |
| Phase 7: Kafka Integration   | 15      | 0         | 0%       |
| Phase 8: Advanced Features   | 20      | 0         | 0%       |
| **TOTAL**                    | **175** | **0**     | **0%**   |

### **Priority Levels**

- 🔴 **Critical (Phase 1-3)**: Must have cho MVP - 70 tasks
- 🟡 **Important (Phase 4-6)**: Production ready - 73 tasks
- 🟢 **Nice to have (Phase 7-8)**: Advanced features - 32 tasks

---

## 🎯 SUCCESS METRICS

### **Development Milestones**

- [ ] MVP ready (Phase 1-3 complete) - Week 3
- [ ] Production ready (Phase 1-6 complete) - Week 6
- [ ] Kafka integration complete (Phase 7) - Week 8
- [ ] Advanced features complete (Phase 8) - Week 12

### **Quality Metrics**

- [ ] Code coverage > 80%
- [ ] Zero critical security vulnerabilities
- [ ] All integration tests passing
- [ ] Load test: 1000 TPS without errors
- [ ] API response time < 2s (p95)

### **Production Metrics**

- [ ] 99.9% uptime
- [ ] Zero duplicate charges
- [ ] < 1% failed payments (excluding user errors)
- [ ] Webhook processing < 5s
- [ ] Complete audit trail for all transactions

---

**Sử dụng checklist này để track progress và đảm bảo không bỏ sót bất kỳ component nào!** ✅
