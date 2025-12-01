# 🔄 Saga Pattern - Giải Thích Chi Tiết

## 📚 Saga Pattern là gì?

**Saga Pattern** là một pattern để quản lý **distributed transactions** trong kiến trúc microservices. Thay vì dùng traditional ACID transactions (không khả thi trong distributed systems), Saga chia transaction thành nhiều **local transactions** nhỏ, mỗi transaction có **compensation action** để rollback.

---

## ❓ Tại sao cần Saga Pattern?

### Vấn đề với Traditional Transactions (ACID)

Trong **monolithic application** với single database:
```java
@Transactional
public void createBooking() {
    // Tất cả trong 1 transaction
    bookingRepository.save(booking);      // Step 1
    ticketService.reserveSeats();         // Step 2
    paymentService.processPayment();       // Step 3
    // Nếu bất kỳ step nào fail → TẤT CẢ rollback tự động
}
```

**Vấn đề**: Trong **microservices**, mỗi service có database riêng:
- `booking-service` → PostgreSQL
- `ticket-service` → PostgreSQL (khác database)
- `payment-service` → PostgreSQL (khác database)

**Không thể dùng 2PC (Two-Phase Commit)** vì:
- Performance issues (locks kéo dài)
- Availability issues (nếu 1 service down → tất cả block)
- Complexity cao

### Giải pháp: Saga Pattern

Saga chia transaction thành các **local transactions** độc lập, mỗi transaction có **compensation** để undo.

---

## 🏗️ Saga Orchestrator Pattern

Có 2 loại Saga:
1. **Choreography** - Mỗi service tự quyết định next step (event-driven)
2. **Orchestration** - Có 1 orchestrator điều phối tất cả steps ⭐ (chúng ta dùng)

### Booking Saga Orchestrator

```java
@Transactional  // Chỉ cho booking-service database
public BookingResult executeBookingSaga(BookingCreateCommand command) {
    // Step 1: Local transaction trong booking-service
    Booking booking = createBookingRecord(command);  // ✅ Commit ngay
    
    // Step 2: Call external service (ticket-service)
    String reservationId = reserveSeats(booking);    // ✅ Commit trong ticket-service
    
    // Step 3: Call external service (payment-service)
    String paymentId = processPayment(booking);      // ✅ Commit trong payment-service
    
    // Step 4: Update local state
    booking.confirm(paymentId);                      // ✅ Commit trong booking-service
    
    return result;
}
```

**Quan trọng**: Mỗi step là **independent transaction**, không có global rollback!

---

## 🔄 Compensation Pattern (Rollback)

Khi một step fail, chúng ta **không thể rollback** các steps trước đó (đã commit rồi). Thay vào đó, chúng ta **compensate** (undo) bằng cách gọi **reverse operations**.

### Ví dụ Compensation

```java
try {
    // Step 1: Create booking ✅ (committed)
    Booking booking = createBookingRecord(command);
    
    // Step 2: Reserve seats ✅ (committed trong ticket-service)
    String reservationId = reserveSeats(booking);
    
    // Step 3: Process payment ❌ FAILED!
    String paymentId = processPayment(booking);  // Throws exception
    
} catch (Exception e) {
    // Compensation: Undo các steps đã commit
    compensate(booking, reservationId, paymentId);
}

private void compensate(Booking booking, String reservationId, String paymentId) {
    // Compensate Step 3: Cancel payment (nếu đã tạo)
    if (paymentId != null) {
        paymentServiceClient.cancelPayment(paymentId);  // Reverse Step 3
    }
    
    // Compensate Step 2: Release seats (nếu đã reserve)
    if (reservationId != null) {
        ticketServiceClient.releaseTickets(reservationId);  // Reverse Step 2
    }
    
    // Compensate Step 1: Mark booking as FAILED
    booking.setStatus(BookingStatus.FAILED);  // Update local state
}
```

---

## 📊 Flow Diagram Chi Tiết

### Happy Path (Success)

```
┌─────────────────────────────────────────────────────────┐
│ Step 1: Create Booking (Local Transaction)              │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ bookingRepository.save(booking)                     │ │
│ │ Status: PENDING                                      │ │
│ │ ✅ COMMIT trong booking-service database             │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│ Step 2: Reserve Seats (External Service Call)           │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ticketServiceClient.reserveTickets()                │ │
│ │ ✅ COMMIT trong ticket-service database             │ │
│ │ Returns: reservationId                               │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│ Step 3: Process Payment (External Service Call)         │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ paymentServiceClient.createPayment()                │ │
│ │ paymentServiceClient.capturePayment()               │ │
│ │ ✅ COMMIT trong payment-service database            │ │
│ │ Returns: paymentId                                   │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│ Step 4: Confirm Booking (Local Transaction)            │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ booking.setStatus(CONFIRMED)                        │ │
│ │ booking.setPaymentReference(paymentId)              │ │
│ │ ✅ COMMIT trong booking-service database            │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                        ↓
                    ✅ SUCCESS
```

### Failure Path (Compensation)

```
┌─────────────────────────────────────────────────────────┐
│ Step 1: Create Booking ✅                                │
│ ✅ COMMITTED                                             │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│ Step 2: Reserve Seats ✅                                 │
│ ✅ COMMITTED (reservationId = "res-123")                │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│ Step 3: Process Payment ❌ FAILED!                      │
│ Exception: "Payment gateway timeout"                    │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│ COMPENSATION: Undo Steps                                 │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Step 3 Reverse: Cancel Payment (nếu đã tạo)         │ │
│ │ ❌ paymentId = null → Skip                           │ │
│ └─────────────────────────────────────────────────────┘ │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Step 2 Reverse: Release Seats                       │ │
│ │ ✅ ticketServiceClient.releaseTickets("res-123")    │ │
│ └─────────────────────────────────────────────────────┘ │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Step 1 Reverse: Mark as FAILED                      │ │
│ │ ✅ booking.setStatus(FAILED)                        │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                        ↓
                    ❌ FAILED
```

---

## 🔑 Key Concepts

### 1. **Local Transactions Only**

Mỗi step là **independent transaction** trong service của nó:

```java
// Step 1: Local transaction trong booking-service
@Transactional  // Chỉ cho booking-service DB
Booking booking = bookingRepository.save(newBooking);  // ✅ Commit ngay

// Step 2: External call → transaction trong ticket-service
String reservationId = ticketServiceClient.reserveTickets(...);  // ✅ Commit trong ticket-service

// Step 3: External call → transaction trong payment-service  
String paymentId = paymentServiceClient.createPayment(...);  // ✅ Commit trong payment-service
```

**Không có global transaction!**

### 2. **Compensation Actions**

Mỗi step cần có **reverse operation**:

| Step | Forward Action | Compensation Action |
|------|---------------|---------------------|
| Create Booking | `save(booking)` | `setStatus(FAILED)` |
| Reserve Seats | `reserveTickets()` | `releaseTickets()` |
| Process Payment | `createPayment()` | `cancelPayment()` |

### 3. **Idempotency**

Compensation actions phải **idempotent** (có thể gọi nhiều lần mà kết quả giống nhau):

```java
// ✅ Idempotent
public void releaseTickets(String reservationId) {
    if (isAlreadyReleased(reservationId)) {
        return;  // Already released, do nothing
    }
    // Release tickets...
}

// ❌ Not idempotent (sẽ fail nếu gọi 2 lần)
public void releaseTickets(String reservationId) {
    deleteReservation(reservationId);  // Fail nếu đã xóa
}
```

### 4. **Eventual Consistency**

Saga không đảm bảo **immediate consistency** như ACID transactions. Có thể có **temporary inconsistency**:

```
Time T1: Booking created (PENDING)
Time T2: Seats reserved (RESERVING)
Time T3: Payment processing (PROCESSING_PAYMENT)
Time T4: Payment failed
Time T5: Compensation starts
Time T6: Seats released
Time T7: Booking marked as FAILED
```

Trong khoảng T3-T7, system có thể **inconsistent** (seats reserved nhưng payment failed).

---

## 💡 So Sánh với ACID Transactions

| Aspect | ACID Transaction | Saga Pattern |
|--------|------------------|--------------|
| **Scope** | Single database | Multiple services/databases |
| **Consistency** | Immediate (strong) | Eventual (weak) |
| **Rollback** | Automatic (database) | Manual (compensation) |
| **Performance** | Fast (local) | Slower (network calls) |
| **Availability** | All or nothing | Partial success possible |
| **Complexity** | Low | High (need compensation logic) |

---

## 🎯 Trong Booking Service

### Code Structure

```java
@Transactional  // Chỉ cho booking-service DB
public BookingResult executeBookingSaga(BookingCreateCommand command) {
    RLock lock = null;
    Booking booking = null;
    String reservationId = null;
    String paymentId = null;

    try {
        // Step 1: Lock + Create booking (local transaction)
        lock = acquireLock(command.getEventId());
        booking = createBookingRecord(command);  // ✅ Committed
        
        // Step 2: Reserve seats (external service)
        reservationId = reserveSeats(booking);    // ✅ Committed in ticket-service
        
        // Step 3: Process payment (external service)
        paymentId = processPayment(booking);      // ✅ Committed in payment-service
        
        // Step 4: Confirm (local transaction)
        booking.confirm(paymentId);               // ✅ Committed
        
        return result;
        
    } catch (Exception e) {
        // Compensation: Undo all committed steps
        compensate(booking, reservationId, paymentId, e.getMessage());
        throw e;
    } finally {
        releaseLock(lock);
    }
}
```

### Compensation Logic

```java
private void compensate(Booking booking, String reservationId, 
                        String paymentId, String reason) {
    // Compensate Step 3: Cancel payment
    if (paymentId != null) {
        paymentServiceClient.cancelPayment(paymentId);  // Reverse payment
    }
    
    // Compensate Step 2: Release seats
    if (reservationId != null && booking != null) {
        ticketServiceClient.releaseTickets(reservationId, null);  // Reverse reservation
    }
    
    // Compensate Step 1: Mark as failed
    if (booking != null) {
        booking.setStatus(BookingStatus.FAILED);
        bookingRepository.save(booking);  // Update local state
    }
}
```

---

## ⚠️ Challenges & Best Practices

### 1. **Compensation Failures**

Nếu compensation cũng fail thì sao?

```java
// Compensation Step 2 fails
try {
    ticketServiceClient.releaseTickets(reservationId);  // ❌ Network error
} catch (Exception e) {
    log.error("Compensation failed - seats may remain reserved", e);
    // Options:
    // 1. Retry compensation later (background job)
    // 2. Alert operations team
    // 3. Manual intervention required
}
```

**Solution**: Implement **retry mechanism** cho compensation actions.

### 2. **Partial Compensation**

Cần track **compensation status** để biết đã compensate step nào:

```java
private void compensate(Booking booking, String reservationId, 
                        String paymentId, String reason) {
    boolean paymentCompensated = false;
    boolean seatsCompensated = false;
    
    try {
        if (paymentId != null) {
            paymentServiceClient.cancelPayment(paymentId);
            paymentCompensated = true;
        }
    } catch (Exception e) {
        log.error("Payment compensation failed", e);
    }
    
    try {
        if (reservationId != null) {
            ticketServiceClient.releaseTickets(reservationId);
            seatsCompensated = true;
        }
    } catch (Exception e) {
        log.error("Seat compensation failed", e);
    }
    
    // Log compensation status for monitoring
    log.info("Compensation status: payment={}, seats={}", 
             paymentCompensated, seatsCompensated);
}
```

### 3. **Idempotency Keys**

Dùng **idempotency keys** để tránh duplicate operations:

```java
// Use booking reference as idempotency key
String idempotencyKey = booking.getBookingReference();

PaymentProto.Payment payment = paymentServiceClient.createPayment(
    bookingId,
    userId,
    amount,
    currency,
    paymentMethod,
    gatewayProvider,
    idempotencyKey,  // ✅ Prevent duplicate payments
    metadata
);
```

---

## 📝 Tóm Tắt

1. **Saga Pattern** = Chia transaction thành nhiều local transactions
2. **Mỗi step** commit độc lập (không có global rollback)
3. **Compensation** = Reverse operations để undo khi fail
4. **Eventual Consistency** = Có thể có temporary inconsistency
5. **Idempotency** = Compensation actions phải idempotent
6. **Orchestrator** = Central coordinator điều phối tất cả steps

**Trade-off**: 
- ✅ Availability cao (không block toàn bộ system)
- ✅ Performance tốt (không cần global locks)
- ❌ Complexity cao (cần implement compensation)
- ❌ Eventual consistency (không phải immediate)

---

**Last Updated**: 2024


