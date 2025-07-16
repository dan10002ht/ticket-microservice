# Event-centric Venue & Layout Model

> **Lưu ý:** Không còn venue-service, không còn bảng venue/layout/zone/seat dùng chung. Tất cả venue/layout/zone/seat là property của event, chỉ quản lý trong event-service.

## 1. Tổng quan nghiệp vụ

- **Venue** không còn là bảng dùng chung, mà là thông tin nhúng trong từng **event**.
- **Mỗi event** thuộc về một **organization** (tổ chức sở hữu event).
- **Mỗi event** có venue riêng (có thể trùng tên, không cần unique toàn hệ thống).
- **Layout (canvas_config)**, **zone**, **seat** đều gắn trực tiếp với event.
- Khi tạo event, user nhập venue info, vẽ layout, tạo zone/seat cho event đó.

## 2. Sơ đồ quan hệ

```
Organization
  |
  | 1 - N
  v
Event (venue info nhúng, layout riêng)
  |
  | 1 - N
  v
EventSeatingZone
  |
  | 1 - N
  v
EventSeat
```

## 3. Ví dụ schema (PostgreSQL)

```sql
-- Organization
CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ...
);

-- Event (venue info nhúng, layout riêng)
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    -- Venue info nhúng
    venue_name VARCHAR(255) NOT NULL,
    venue_address TEXT,
    venue_city VARCHAR(100),
    venue_country VARCHAR(100),
    venue_capacity INTEGER,
    -- Layout
    canvas_config JSONB NOT NULL, -- Layout cho event này
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Zone cho event
CREATE TABLE event_seating_zones (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    zone_type VARCHAR(50),
    coordinates JSONB NOT NULL, -- Boundary trên canvas
    seat_count INTEGER NOT NULL,
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seat cho event
CREATE TABLE event_seats (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    zone_id BIGINT NOT NULL REFERENCES event_seating_zones(id) ON DELETE CASCADE,
    seat_number VARCHAR(20) NOT NULL,
    row_number VARCHAR(10),
    coordinates JSONB NOT NULL, -- Vị trí trên canvas
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 4. Luồng nghiệp vụ

1. **Organization** tạo event → nhập venue info (tên, địa chỉ, ...), nhập thông tin event.
2. **Tạo layout cho event** (canvas_config, vẽ zone, seat).
3. **Tạo các zone và seat cho event** (gắn với event_id).
4. **Booking, bán vé, ...**: Dựa trên layout/zone/seat của event.

## 5. Ưu điểm

- Linh hoạt tối đa: mỗi event có thể có venue, layout, zone, seat riêng biệt.
- Không bị ràng buộc venue/layout cứng.
- Dễ mở rộng, dễ quản lý theo từng event/organization.

> **Mô hình này thay thế hoàn toàn mô hình venue/layout/zone/seat dùng chung trước đây.**
