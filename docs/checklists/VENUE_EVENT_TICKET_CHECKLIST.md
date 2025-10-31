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

### 🔄 Hệ thống Kho hàng Thời gian thực & Kiểm soát Đồng thời

- [x] **Hệ thống Kho hàng Thời gian thực**

  - [x] Thêm tracking trạng thái vào bảng event_seats (status, reserved_by, reserved_until, booked_by)
  - [x] Xây dựng API kiểm tra ghế có sẵn theo thời gian thực
  - [ ] Thêm WebSocket hỗ trợ cập nhật thời gian thực
  - [ ] Xây dựng hệ thống đặt chỗ với timeout

- [x] **Kiểm soát Đồng thời**

  - [ ] Triển khai Redis distributed locks cho việc đặt ghế
  - [x] Thêm optimistic locking với version fields
  - [ ] Xây dựng cơ chế timeout cho đặt chỗ
  - [ ] Ngăn chặn race condition

- [ ] **Tính năng Đặt vé Nâng cao**
  - [ ] Đặt chỗ với thời hạn hết hạn
  - [ ] Quản lý phiên đặt vé
  - [ ] Giải phóng ghế khi hết timeout
  - [ ] Ngăn chặn đặt vé đồng thời

## 🚀 Phase 5: Hiệu suất & Khả năng Mở rộng

- [ ] **Chiến lược Caching**

  - [ ] Redis caching cho tình trạng ghế có sẵn
  - [ ] Vô hiệu hóa cache khi trạng thái ghế thay đổi
  - [ ] Cache phân tán cho đồng thời cao

- [x] **Tối ưu hóa Database**

  - [x] Thêm indexes cho truy vấn trạng thái ghế
  - [x] Triển khai connection pooling (PgPool-II)
  - [ ] Thêm giám sát database

- [ ] **Cập nhật Thời gian thực**
  - [ ] Tích hợp WebSocket cho cập nhật ghế
  - [ ] Redis pub/sub cho thay đổi trạng thái ghế
  - [ ] Phát sóng kho hàng thời gian thực

### ⏳ Backlog (sẽ thực hiện sau)

- [ ] Xây dựng unit test, integration test cho các luồng booking, seat reservation

---

## 📝 Ghi chú

- **Không còn venue-service, không còn bảng venue/layout/zone/seat dùng chung.**
- **Mọi logic venue/layout/zone/seat đều nằm trong event-service, gắn với từng event.**
- Nếu có tài liệu cũ, hãy tham khảo README_EVENT_MODEL.md và EVENT_NEW.md để hiểu mô hình mới.
