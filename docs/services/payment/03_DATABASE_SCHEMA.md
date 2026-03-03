## 🗄️ Payment Service - Database Schema

Complete database schema documentation for Payment Service.

---

## 📊 Tables Overview

| Table              | Purpose                          | Records                  |
| ------------------ | -------------------------------- | ------------------------ |
| `payments`         | Core payment transactions        | High volume              |
| `refunds`          | Refund transactions              | Medium volume            |
| `transaction_logs` | Audit trail for all transactions | Very high volume         |
| `idempotency_keys` | Prevent duplicate requests       | Medium volume (with TTL) |

---

## 📋 Table: `payments`

**Purpose**: Store all payment transactions

### Columns

| Column                | Type          | Constraints                 | Description              |
| --------------------- | ------------- | --------------------------- | ------------------------ |
| `id`                  | BIGSERIAL     | PRIMARY KEY                 | Internal ID              |
| `payment_id`          | UUID          | UNIQUE, NOT NULL            | Public-facing UUID       |
| `booking_id`          | VARCHAR(255)  | NOT NULL, INDEXED           | Reference to booking     |
| `ticket_id`           | VARCHAR(255)  | NULLABLE                    | Reference to ticket      |
| `user_id`             | VARCHAR(255)  | NOT NULL, INDEXED           | User who made payment    |
| `amount`              | DECIMAL(10,2) | NOT NULL, CHECK >= 0        | Payment amount           |
| `currency`            | VARCHAR(3)    | NOT NULL, DEFAULT 'USD'     | Currency code (ISO 4217) |
| `payment_method`      | VARCHAR(50)   | NOT NULL                    | Payment method used      |
| `status`              | VARCHAR(50)   | NOT NULL, DEFAULT 'pending' | Payment status           |
| `external_reference`  | VARCHAR(255)  | NULLABLE, INDEXED           | Gateway transaction ID   |
| `gateway_provider`    | VARCHAR(50)   | NULLABLE, INDEXED           | Gateway used             |
| `gateway_response`    | JSONB         | NULLABLE                    | Full gateway response    |
| `failure_reason`      | TEXT          | NULLABLE                    | Failure description      |
| `description`         | TEXT          | NULLABLE                    | Payment description      |
| `metadata`            | JSONB         | NULLABLE                    | Additional data          |
| `idempotency_key`     | VARCHAR(255)  | UNIQUE                      | Idempotency key          |
| `created_at`          | TIMESTAMP     | NOT NULL, DEFAULT NOW()     | Creation timestamp       |
| `updated_at`          | TIMESTAMP     | NOT NULL, AUTO-UPDATE       | Update timestamp         |
| `created_by`          | VARCHAR(255)  | NULLABLE                    | Creator user ID          |
| `updated_by`          | VARCHAR(255)  | NULLABLE                    | Last updater user ID     |
| `paid_at`             | TIMESTAMP     | NULLABLE                    | Payment completion time  |
| `cancelled_at`        | TIMESTAMP     | NULLABLE                    | Cancellation time        |
| `cancelled_by`        | VARCHAR(255)  | NULLABLE                    | Canceller user ID        |
| `cancellation_reason` | TEXT          | NULLABLE                    | Cancellation reason      |

### Status Values

- `pending`: Payment initiated, waiting for processing
- `processing`: Payment being processed by gateway
- `success`: Payment completed successfully
- `failed`: Payment failed
- `cancelled`: Payment cancelled by user/system
- `refunded`: Payment fully refunded
- `partially_refunded`: Payment partially refunded

### Payment Method Values

- `credit_card`: Credit card payment
- `debit_card`: Debit card payment
- `bank_transfer`: Bank transfer
- `e_wallet`: E-wallet (PayPal, Momo, etc.)
- `cash`: Cash payment

### Indexes

- `idx_payments_booking_id`: Fast lookup by booking
- `idx_payments_user_id`: Fast lookup by user
- `idx_payments_status`: Filter by status
- `idx_payments_payment_id`: Lookup by public ID
- `idx_payments_external_reference`: Gateway transaction lookup
- `idx_payments_gateway_provider`: Filter by gateway
- `idx_payments_idempotency_key`: Idempotency check
- `idx_payments_user_status`: Composite for user + status queries
- `idx_payments_booking_status`: Composite for booking + status queries

---

## 📋 Table: `refunds`

**Purpose**: Store refund transactions

### Columns

| Column                | Type          | Constraints                 | Description               |
| --------------------- | ------------- | --------------------------- | ------------------------- |
| `id`                  | BIGSERIAL     | PRIMARY KEY                 | Internal ID               |
| `refund_id`           | UUID          | UNIQUE, NOT NULL            | Public-facing UUID        |
| `payment_id`          | BIGINT        | FK to payments, NOT NULL    | Parent payment            |
| `payment_uuid`        | UUID          | NOT NULL                    | Denormalized payment UUID |
| `amount`              | DECIMAL(10,2) | NOT NULL, CHECK > 0         | Refund amount             |
| `currency`            | VARCHAR(3)    | NOT NULL, DEFAULT 'USD'     | Currency code             |
| `refund_type`         | VARCHAR(50)   | NOT NULL                    | full or partial           |
| `status`              | VARCHAR(50)   | NOT NULL, DEFAULT 'pending' | Refund status             |
| `external_reference`  | VARCHAR(255)  | NULLABLE, INDEXED           | Gateway refund ID         |
| `gateway_provider`    | VARCHAR(50)   | NULLABLE                    | Gateway used              |
| `gateway_response`    | JSONB         | NULLABLE                    | Full gateway response     |
| `failure_reason`      | TEXT          | NULLABLE                    | Failure description       |
| `reason`              | VARCHAR(255)  | NOT NULL                    | Short reason              |
| `description`         | TEXT          | NULLABLE                    | Detailed description      |
| `metadata`            | JSONB         | NULLABLE                    | Additional data           |
| `created_at`          | TIMESTAMP     | NOT NULL, DEFAULT NOW()     | Creation timestamp        |
| `updated_at`          | TIMESTAMP     | NOT NULL, AUTO-UPDATE       | Update timestamp          |
| `created_by`          | VARCHAR(255)  | NULLABLE                    | Creator user ID           |
| `updated_by`          | VARCHAR(255)  | NULLABLE                    | Last updater user ID      |
| `refunded_at`         | TIMESTAMP     | NULLABLE                    | Refund completion time    |
| `cancelled_at`        | TIMESTAMP     | NULLABLE                    | Cancellation time         |
| `cancelled_by`        | VARCHAR(255)  | NULLABLE                    | Canceller user ID         |
| `cancellation_reason` | TEXT          | NULLABLE                    | Cancellation reason       |

### Status Values

- `pending`: Refund initiated
- `processing`: Refund being processed
- `success`: Refund completed
- `failed`: Refund failed
- `cancelled`: Refund cancelled

### Refund Type Values

- `full`: Full refund (entire payment amount)
- `partial`: Partial refund (less than payment amount)

---

## 📋 Table: `transaction_logs`

**Purpose**: Audit trail for all payment/refund transactions

### Columns

| Column               | Type         | Constraints              | Description               |
| -------------------- | ------------ | ------------------------ | ------------------------- |
| `id`                 | BIGSERIAL    | PRIMARY KEY              | Internal ID               |
| `log_id`             | UUID         | UNIQUE, NOT NULL         | Public-facing UUID        |
| `payment_id`         | BIGINT       | FK to payments, NULLABLE | Related payment           |
| `payment_uuid`       | UUID         | NULLABLE                 | Denormalized payment UUID |
| `refund_id`          | BIGINT       | FK to refunds, NULLABLE  | Related refund            |
| `refund_uuid`        | UUID         | NULLABLE                 | Denormalized refund UUID  |
| `transaction_type`   | VARCHAR(50)  | NOT NULL                 | Event type                |
| `event_name`         | VARCHAR(100) | NOT NULL                 | Specific event name       |
| `gateway_provider`   | VARCHAR(50)  | NULLABLE                 | Gateway used              |
| `external_reference` | VARCHAR(255) | NULLABLE                 | Gateway transaction ID    |
| `request_data`       | JSONB        | NULLABLE, INDEXED        | Request payload           |
| `response_data`      | JSONB        | NULLABLE, INDEXED        | Response payload          |
| `headers`            | JSONB        | NULLABLE                 | HTTP headers              |
| `status`             | VARCHAR(50)  | NOT NULL                 | Log status                |
| `error_code`         | VARCHAR(100) | NULLABLE                 | Error code if failed      |
| `error_message`      | TEXT         | NULLABLE                 | Error message             |
| `duration_ms`        | INTEGER      | NULLABLE                 | Duration in ms            |
| `user_id`            | VARCHAR(255) | NULLABLE                 | User context              |
| `ip_address`         | VARCHAR(45)  | NULLABLE                 | Client IP                 |
| `user_agent`         | TEXT         | NULLABLE                 | Client user agent         |
| `correlation_id`     | UUID         | NULLABLE, INDEXED        | Distributed tracing ID    |
| `metadata`           | JSONB        | NULLABLE, INDEXED        | Additional data           |
| `created_at`         | TIMESTAMP    | NOT NULL, DEFAULT NOW()  | Timestamp                 |
| `created_by`         | VARCHAR(255) | NULLABLE                 | Creator                   |

### Transaction Type Values

- `payment_initiated`
- `payment_success`
- `payment_failed`
- `refund_initiated`
- `refund_success`
- `refund_failed`
- `webhook_received`

---

## 📋 Table: `idempotency_keys`

**Purpose**: Prevent duplicate payment requests

### Columns

| Column             | Type         | Constraints                    | Description         |
| ------------------ | ------------ | ------------------------------ | ------------------- |
| `id`               | BIGSERIAL    | PRIMARY KEY                    | Internal ID         |
| `idempotency_key`  | VARCHAR(255) | UNIQUE, NOT NULL               | Client-provided key |
| `request_path`     | VARCHAR(500) | NOT NULL                       | API endpoint        |
| `request_method`   | VARCHAR(10)  | NOT NULL                       | HTTP method         |
| `request_body`     | JSONB        | NULLABLE                       | Original request    |
| `request_headers`  | JSONB        | NULLABLE                       | Original headers    |
| `response_status`  | INTEGER      | NULLABLE                       | HTTP status code    |
| `response_body`    | JSONB        | NULLABLE                       | Response data       |
| `response_headers` | JSONB        | NULLABLE                       | Response headers    |
| `payment_id`       | BIGINT       | FK to payments, NULLABLE       | Related payment     |
| `payment_uuid`     | UUID         | NULLABLE                       | Denormalized UUID   |
| `refund_id`        | BIGINT       | FK to refunds, NULLABLE        | Related refund      |
| `refund_uuid`      | UUID         | NULLABLE                       | Denormalized UUID   |
| `user_id`          | VARCHAR(255) | NULLABLE                       | User context        |
| `ip_address`       | VARCHAR(45)  | NULLABLE                       | Client IP           |
| `user_agent`       | TEXT         | NULLABLE                       | Client user agent   |
| `status`           | VARCHAR(50)  | NOT NULL, DEFAULT 'processing' | Status              |
| `expires_at`       | TIMESTAMP    | NOT NULL, INDEXED              | Expiration time     |
| `metadata`         | JSONB        | NULLABLE                       | Additional data     |
| `created_at`       | TIMESTAMP    | NOT NULL, DEFAULT NOW()        | Creation time       |
| `updated_at`       | TIMESTAMP    | NOT NULL, AUTO-UPDATE          | Update time         |
| `completed_at`     | TIMESTAMP    | NULLABLE                       | Completion time     |

### Status Values

- `processing`: Request in progress
- `completed`: Request completed
- `failed`: Request failed

---

## 📊 Views

### `v_user_payment_summary`

User payment statistics

**Columns**: `user_id`, `total_payments`, `successful_payments`, `failed_payments`, `refunded_payments`, `total_amount_paid`, `total_amount_refunded`, `average_payment_amount`, `first_payment_date`, `last_payment_date`

### `v_gateway_payment_summary`

Gateway performance statistics

**Columns**: `gateway_provider`, `total_payments`, `successful_payments`, `failed_payments`, `total_amount_processed`, `average_payment_amount`, `success_rate_percentage`

### `v_daily_payment_statistics`

Daily payment metrics

**Columns**: `payment_date`, `total_payments`, `successful_payments`, `failed_payments`, `daily_revenue`, `average_payment_amount`, `unique_users`, `unique_bookings`

### `v_refund_statistics`

Daily refund metrics

**Columns**: `refund_date`, `total_refunds`, `successful_refunds`, `full_refunds`, `partial_refunds`, `total_refunded_amount`, `average_refund_amount`

---

## 🔧 Functions

### `get_total_refunded_amount(payment_id BIGINT)`

Calculate total refunded amount for a payment

**Returns**: `DECIMAL(10,2)`

### `can_payment_be_refunded(payment_id BIGINT)`

Check if payment can be refunded

**Returns**: `BOOLEAN`

### `get_remaining_refundable_amount(payment_id BIGINT)`

Calculate remaining refundable amount

**Returns**: `DECIMAL(10,2)`

### `cleanup_expired_idempotency_keys()`

Delete expired idempotency keys

**Returns**: `void`

---

## 🔐 Security Considerations

1. **Sensitive Data**: Gateway responses stored in JSONB (ensure encryption at rest)
2. **PII**: User IDs stored as references, not full user data
3. **Access Control**: Application-level permissions required
4. **Audit Trail**: All changes tracked via updated_at triggers
5. **Idempotency**: Prevents duplicate charges

---

## 📈 Performance Optimizations

1. **Indexes**: Comprehensive indexing on frequently queried columns
2. **Composite Indexes**: For common multi-column queries
3. **GIN Indexes**: On JSONB columns for flexible querying
4. **Partial Indexes**: On active records only
5. **Views**: Pre-aggregated statistics for reporting
6. **Functions**: Database-level business logic for performance

---

## 🧹 Maintenance

### Cleanup Expired Idempotency Keys

```sql
SELECT cleanup_expired_idempotency_keys();
```

Schedule this via cron or Spring Scheduler.

### Archive Old Transaction Logs

Transaction logs grow quickly. Consider archiving logs older than 90 days:

```sql
-- Create archive table
CREATE TABLE transaction_logs_archive (LIKE transaction_logs INCLUDING ALL);

-- Move old logs
INSERT INTO transaction_logs_archive
SELECT * FROM transaction_logs
WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '90 days';

-- Delete from main table
DELETE FROM transaction_logs
WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '90 days';
```

---

## 📝 Migration Files

- `V1__create_payments_table.sql` - Payments table
- `V2__create_refunds_table.sql` - Refunds table
- `V3__create_transaction_logs_table.sql` - Transaction logs
- `V4__create_idempotency_keys_table.sql` - Idempotency keys
- `V5__add_constraints_and_views.sql` - Views, functions, constraints

---

**Database schema version**: 5  
**Last updated**: 2024  
**Total tables**: 4  
**Total views**: 4  
**Total functions**: 4
