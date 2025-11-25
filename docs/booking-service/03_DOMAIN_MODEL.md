# 📐 Booking Domain Model (Planning)

Reference document for booking-service entities, state machine, and saga orchestration.

## 🧬 Core Entities
- `Booking`: aggregate root, tracks lifecycle status, user, session, totals, provider references.
- `BookingItem`: per-seat / per-ticket detail (eventId, seatId, pricing category, amount).
- `BookingAction` / `BookingTask`: saga steps, retries, compensation commands.
- `BookingLock`: Redis-based locks (document key format, TTL).

## 🔄 Lifecycle States (draft)
1. `PENDING` – request accepted, waiting for queue
2. `RESERVING` – ticket-service seat hold in progress
3. `AWAITING_PAYMENT` – seats held, waiting for payment authorization
4. `PROCESSING_PAYMENT` – payment-service interaction
5. `CONFIRMED` – booking finalized, tickets issued
6. `FAILED` / `CANCELLED` – compensation triggered (release seats, refund if needed)

## 🔁 Saga Flow (high level)
1. Receive booking command from gateway/worker
2. Lock booking session → call Ticket Service to reserve seats
3. On success, call Payment Service (authorize/capture depending on policy)
4. Persist saga actions (Kafka events, compensation tasks)
5. Release locks + publish final booking events

## 🧱 TODOs
- [ ] Finalize state diagram + transitions
- [ ] Document compensation steps for each failure point
- [ ] Define Redis key strategy (`booking:lock:{sessionId}` etc.)
- [ ] Align with ticket-service reservation expiry to avoid deadlocks
- [ ] Add sequence diagrams (queue → booking-service → ticket-service/payment-service)

_Last updated: Planning stage (2024)_
