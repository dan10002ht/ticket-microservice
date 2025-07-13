# Venue, Event & Ticket Management Flow

## 🎯 Overview

Hệ thống quản lý venue, event và ticket được thiết kế để hỗ trợ:
- **Venue Management**: Quản lý địa điểm với sơ đồ chỗ ngồi linh hoạt
- **Event Management**: Tạo và quản lý sự kiện với pricing động
- **Ticket Management**: Hệ thống đặt vé theo thời gian thực
- **Seating Layout**: Canvas-based UI cho sơ đồ chỗ ngồi

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Venue Service │    │   Event Service │    │  Ticket Service │
│                 │    │                 │    │                 │
│ • Venue CRUD    │    │ • Event CRUD    │    │ • Ticket CRUD   │
│ • Layout Mgmt   │    │ • Pricing Mgmt  │    │ • Booking Flow  │
│ • Zone Mgmt     │    │ • Availability  │    │ • Seat Lock     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Gateway API   │
                    │                 │
                    │ • REST Endpoints│
                    │ • Auth/Authorize│
                    │ • Rate Limiting │
                    └─────────────────┘
```

## 📊 Database Schema

### Core Tables

#### 1. Venues
```sql
venues (
  id, public_id, name, description, address, city, state, country,
  latitude, longitude, phone, email, website, capacity,
  facilities, images, is_active, created_by, created_at, updated_at
)
```

#### 2. Venue Layouts
```sql
venue_layouts (
  id, venue_id, name, description, layout_config, total_seats,
  total_rows, total_columns, is_active, created_by, created_at, updated_at
)
```

#### 3. Seating Zones
```sql
seating_zones (
  id, venue_layout_id, name, description, color, zone_coordinates,
  total_seats, base_price, is_active, display_order, created_at, updated_at
)
```

#### 4. Seats
```sql
seats (
  id, venue_layout_id, seating_zone_id, seat_number, row_number, column_number,
  seat_coordinates, seat_type, is_available, is_active, created_at, updated_at
)
```

#### 5. Events
```sql
events (
  id, public_id, title, description, venue_id, venue_layout_id, organizer_id,
  event_type, status, start_time, end_time, doors_open_time,
  max_capacity, current_bookings, event_images, event_details,
  is_featured, is_active, created_at, updated_at
)
```

#### 6. Event Pricing
```sql
event_pricing (
  id, event_id, seating_zone_id, price, original_price, available_seats,
  total_seats, pricing_type, pricing_rules, early_bird_end, early_bird_price,
  is_active, created_at, updated_at
)
```

#### 7. Event Seat Availability
```sql
event_seat_availability (
  id, event_id, seat_id, status, reserved_by, reserved_until,
  booked_by, booked_at, final_price, created_at, updated_at
)
```

## 🔄 Business Flows

### 1. Venue Creation Flow

```
1. Create Venue
   ├── Basic Info (name, address, contact)
   ├── Location (coordinates, city, country)
   ├── Capacity & Facilities
   └── Images & Description

2. Create Venue Layout
   ├── Layout Name & Description
   ├── Canvas Configuration
   ├── Total Dimensions (rows, columns)
   └── Save Layout

3. Define Seating Zones
   ├── Zone Name & Color
   ├── Canvas Coordinates
   ├── Base Pricing
   └── Seat Count

4. Generate Seats
   ├── Auto-generate based on dimensions
   ├── Assign to zones
   ├── Set seat numbers
   └── Canvas coordinates
```

### 2. Event Creation Flow

```
1. Create Event
   ├── Basic Info (title, description, type)
   ├── Select Venue & Layout
   ├── Set Date & Time
   └── Organizer Details

2. Configure Pricing
   ├── Set zone-based pricing
   ├── Configure dynamic pricing rules
   ├── Set early bird pricing
   └── Define seat availability

3. Publish Event
   ├── Review all details
   ├── Set status to published
   ├── Initialize seat availability
   └── Start accepting bookings
```

### 3. Ticket Booking Flow

```
1. User Browse Events
   ├── Search & Filter events
   ├── View event details
   ├── Check availability
   └── Select event

2. Seat Selection
   ├── View seating layout
   ├── Interactive seat map
   ├── Real-time availability
   └── Select seats

3. Pricing & Payment
   ├── Calculate total price
   ├── Apply discounts
   ├── Process payment
   └── Confirm booking

4. Ticket Generation
   ├── Generate ticket
   ├── Send confirmation
   ├── Update availability
   └── Notify user
```

## 🎨 Canvas Layout System

### Layout Configuration Structure

```json
{
  "canvas": {
    "width": 1200,
    "height": 800,
    "background": "#f5f5f5",
    "grid": {
      "enabled": true,
      "size": 20,
      "color": "#e0e0e0"
    }
  },
  "stage": {
    "x": 100,
    "y": 50,
    "width": 1000,
    "height": 100,
    "type": "rectangle",
    "color": "#333333",
    "label": "Stage"
  },
  "zones": [
    {
      "id": "zone_1",
      "name": "VIP",
      "x": 150,
      "y": 200,
      "width": 200,
      "height": 150,
      "color": "#ffd700",
      "price": 150.00,
      "seats": [
        {
          "row": 1,
          "col": 1,
          "x": 160,
          "y": 210,
          "seat_number": "A1",
          "type": "standard"
        }
      ]
    }
  ],
  "legend": {
    "enabled": true,
    "position": "top-right",
    "items": [
      {
        "label": "Available",
        "color": "#4CAF50"
      },
      {
        "label": "Reserved",
        "color": "#FF9800"
      },
      {
        "label": "Booked",
        "color": "#F44336"
      }
    ]
  }
}
```

### Interactive Features

1. **Drag & Drop**: Kéo thả zones và seats
2. **Real-time Updates**: Cập nhật availability ngay lập tức
3. **Zoom & Pan**: Phóng to, thu nhỏ và di chuyển
4. **Multi-select**: Chọn nhiều chỗ ngồi cùng lúc
5. **Price Calculation**: Tính giá tự động
6. **Validation**: Kiểm tra tính hợp lệ

## 💰 Pricing System

### Pricing Types

1. **Fixed Pricing**
   - Giá cố định cho mỗi zone
   - Không thay đổi theo thời gian

2. **Dynamic Pricing**
   - Giá thay đổi theo demand
   - Algorithm-based pricing
   - Real-time adjustments

3. **Tiered Pricing**
   - Early bird discounts
   - Group discounts
   - VIP pricing
   - Last-minute pricing

### Pricing Rules

```json
{
  "base_price": 100.00,
  "early_bird": {
    "enabled": true,
    "end_date": "2024-02-01T00:00:00Z",
    "discount_percent": 20
  },
  "dynamic_pricing": {
    "enabled": true,
    "rules": [
      {
        "condition": "occupancy > 80%",
        "action": "increase_price",
        "percentage": 10
      },
      {
        "condition": "days_until_event < 7",
        "action": "decrease_price",
        "percentage": 15
      }
    ]
  },
  "group_discounts": [
    {
      "min_tickets": 5,
      "discount_percent": 10
    },
    {
      "min_tickets": 10,
      "discount_percent": 20
    }
  ]
}
```

## 🔒 Seat Management

### Seat Status Flow

```
Available → Reserved → Booked
    ↑           ↓         ↓
    └─── Released ←─── Cancelled
```

### Reservation System

1. **Temporary Hold**: 15 phút để user hoàn tất booking
2. **Auto-release**: Tự động giải phóng sau timeout
3. **Conflict Prevention**: Không cho phép double booking
4. **Real-time Sync**: Cập nhật ngay lập tức cho tất cả users

### Seat Types

- **Standard**: Chỗ ngồi thường
- **Wheelchair**: Chỗ ngồi cho xe lăn
- **Companion**: Chỗ ngồi đi kèm
- **VIP**: Chỗ ngồi VIP
- **Reserved**: Chỗ ngồi được đặt trước

## 📱 API Endpoints

### Venue Management

```
GET    /api/venues                    # List venues
POST   /api/venues                    # Create venue
GET    /api/venues/:id                # Get venue details
PUT    /api/venues/:id                # Update venue
DELETE /api/venues/:id                # Delete venue

GET    /api/venues/:id/layouts        # Get venue layouts
POST   /api/venues/:id/layouts        # Create layout
PUT    /api/venues/:id/layouts/:lid   # Update layout
DELETE /api/venues/:id/layouts/:lid   # Delete layout

GET    /api/layouts/:id/zones         # Get seating zones
POST   /api/layouts/:id/zones         # Create zone
PUT    /api/layouts/:id/zones/:zid    # Update zone
DELETE /api/layouts/:id/zones/:zid    # Delete zone

GET    /api/layouts/:id/seats         # Get seats
POST   /api/layouts/:id/seats         # Create seats
PUT    /api/layouts/:id/seats/:sid    # Update seat
DELETE /api/layouts/:id/seats/:sid    # Delete seat
```

### Event Management

```
GET    /api/events                    # List events
POST   /api/events                    # Create event
GET    /api/events/:id                # Get event details
PUT    /api/events/:id                # Update event
DELETE /api/events/:id                # Delete event

GET    /api/events/:id/pricing        # Get event pricing
POST   /api/events/:id/pricing        # Set pricing
PUT    /api/events/:id/pricing/:pid   # Update pricing

GET    /api/events/:id/availability   # Get seat availability
POST   /api/events/:id/reserve        # Reserve seats
PUT    /api/events/:id/book           # Book seats
DELETE /api/events/:id/release        # Release seats
```

### Ticket Management

```
GET    /api/tickets                   # Get user tickets
POST   /api/tickets                   # Create ticket
GET    /api/tickets/:id               # Get ticket details
PUT    /api/tickets/:id/cancel        # Cancel ticket

GET    /api/tickets/:id/qr            # Get QR code
POST   /api/tickets/:id/validate      # Validate ticket
```

## 🔐 Authorization & Permissions

### Venue Management
- `venue.create` - Tạo venue mới
- `venue.edit` - Chỉnh sửa venue
- `venue.delete` - Xóa venue
- `venue.view` - Xem venue

### Event Management
- `event.create` - Tạo event mới
- `event.edit` - Chỉnh sửa event
- `event.delete` - Xóa event
- `event.publish` - Publish event
- `event.view` - Xem event

### Ticket Management
- `ticket.book` - Đặt vé
- `ticket.cancel` - Hủy vé
- `ticket.view` - Xem vé
- `ticket.validate` - Validate vé

## 📊 Monitoring & Analytics

### Key Metrics

1. **Venue Performance**
   - Occupancy rate
   - Revenue per venue
   - Popular zones
   - Peak hours

2. **Event Analytics**
   - Ticket sales
   - Revenue trends
   - Popular events
   - Cancellation rate

3. **User Behavior**
   - Booking patterns
   - Seat preferences
   - Payment methods
   - User journey

### Real-time Monitoring

- **Seat Availability**: Real-time updates
- **Booking Rate**: TPS (transactions per second)
- **Error Rate**: Failed bookings
- **Performance**: API response times

## 🚀 Implementation Phases

### Phase 1: Core Infrastructure
- [ ] Database schema setup
- [ ] Basic CRUD operations
- [ ] Authentication & authorization
- [ ] Basic API endpoints

### Phase 2: Venue Management
- [ ] Venue CRUD
- [ ] Layout management
- [ ] Zone configuration
- [ ] Seat generation

### Phase 3: Event Management
- [ ] Event CRUD
- [ ] Pricing system
- [ ] Availability tracking
- [ ] Event publishing

### Phase 4: Ticket System
- [ ] Booking flow
- [ ] Payment integration
- [ ] Ticket generation
- [ ] QR code system

### Phase 5: Advanced Features
- [ ] Canvas UI
- [ ] Dynamic pricing
- [ ] Real-time updates
- [ ] Analytics dashboard

## 🔧 Technical Considerations

### Performance
- **Caching**: Redis cho availability data
- **Database**: Optimized queries với indexes
- **Real-time**: WebSocket cho live updates
- **CDN**: Static assets (images, layouts)

### Scalability
- **Microservices**: Separate services cho venue, event, ticket
- **Load Balancing**: Distribute traffic
- **Database Sharding**: Scale horizontally
- **Queue System**: Async processing

### Security
- **Rate Limiting**: Prevent abuse
- **Input Validation**: Sanitize all inputs
- **SQL Injection**: Parameterized queries
- **XSS Protection**: Sanitize outputs

### Reliability
- **Circuit Breaker**: Handle service failures
- **Retry Logic**: Automatic retries
- **Monitoring**: Health checks
- **Backup**: Regular data backups
