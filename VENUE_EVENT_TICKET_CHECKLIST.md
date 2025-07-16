# 🎫 Event-centric Venue, Layout, Zone, Seat Checklist (Mô hình mới)

## 📋 Tổng quan

- Từ nay, **venue/layout/zone/seat đều là thuộc tính của event** (không còn venue-service, không còn bảng venue/layout/zone/seat dùng chung).
- Mỗi event thuộc organization, có venue info nhúng, layout riêng, zone/seat riêng.
- Tất cả logic venue/layout/zone/seat chỉ nằm trong event-service.

## 🏗️ Phase 1: Database Schema & Infrastructure

- [x] Tạo Database Migrations cho event-service (venue/layout/zone/seat đều thuộc event)
- [x] Tạo Protobuf Definitions (event.proto, ticket.proto)
- [x] Service Infrastructure (Go server, gRPC, DB, Redis, Prometheus)

---

## 🏢 Phase 2: Event Service Implementation (Go)

- [x] Core Models & Repositories (event, event_seating_zone, event_seat)
- [x] Business Logic Services (CRUD event, layout, zone, seat)
- [x] Canvas Layout System (canvas_config lưu ở event)
- [x] gRPC Controllers (CRUD event, layout, zone, seat)
- [x] gRPC Services (Internal Communication)
- [x] Caching & Performance (nếu cần)

---

## 🎭 Phase 3: Event Service Advanced Features

- [x] Pricing, discount, seat availability, schedule, ...
- [x] Integration với organization, ticket, payment, ...

---

## 🎫 Phase 4: Ticket Service Implementation (Go)

- [ ] Thiết kế schema/migration cho ticket (gắn với event_id, seat_id, zone_id)
- [ ] Xây dựng models & repository cho ticket, seat reservation
- [ ] Xây dựng service cho booking, seat reservation, ticket issuance
- [ ] Tích hợp gRPC client để lấy dữ liệu seat/zone/layout từ event-service
- [ ] Xây dựng gRPC controllers cho các luồng booking, reserve, issue ticket
- [ ] Xây dựng logic kiểm tra seat availability (lấy từ event-service hoặc cache)
- [ ] Tích hợp với payment-service, notification-service, booking-worker
- [ ] Expose Prometheus metrics cho booking/ticket
- [ ] Viết tài liệu hướng dẫn API, flow booking, seat reservation

### ⏳ Backlog (sẽ thực hiện sau)

- [ ] Xây dựng unit test, integration test cho các luồng booking, seat reservation

---

## 📝 Ghi chú

- **Không còn venue-service, không còn bảng venue/layout/zone/seat dùng chung.**
- **Mọi logic venue/layout/zone/seat đều nằm trong event-service, gắn với từng event.**
- Nếu có tài liệu cũ, hãy tham khảo README_EVENT_MODEL.md và EVENT_NEW.md để hiểu mô hình mới.
