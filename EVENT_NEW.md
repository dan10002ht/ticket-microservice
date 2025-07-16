# Booking System (Event-centric Model)

## Thay đổi lớn:

- Không còn venue-service.
- Mọi logic venue/layout/zone/seat đều thuộc event-service, gắn với từng event.
- Mỗi event thuộc organization, có venue info nhúng, layout riêng, zone/seat riêng.

> **Lưu ý:** Mô hình này đã thay thế hoàn toàn mô hình venue/layout/zone/seat dùng chung trước đây. Nếu bạn từng tham khảo tài liệu cũ, hãy chuyển sang đọc README_EVENT_MODEL.md và tài liệu này.

## Các service còn lại:

- event-service: Quản lý event, venue info, layout, zone, seat, pricing, schedule, ...
- ticket-service: Quản lý booking, seat reservation, ticket issuance
- gateway: API Gateway (Node.js)
- ...

## Luồng nghiệp vụ:

1. Organization tạo event (nhập venue info, vẽ layout, tạo zone/seat)
2. Event-service lưu toàn bộ venue/layout/zone/seat cho event
3. Booking, bán vé, ... dựa trên layout/zone/seat của event

## Xem thêm chi tiết trong README_EVENT_MODEL.md
