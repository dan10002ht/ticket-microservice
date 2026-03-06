# Plan: Phase 12 — Org Dashboard Analytics & Event Publishing UX

## Role & Standards

Bạn là **Senior Full-Stack Engineer**, chịu trách nhiệm implement theo tiêu chuẩn project (xem CLAUDE.md).

---

## Context

Phase 11 hoàn thành (Advanced Event Management). Audit hiện tại cho thấy:
- Org Dashboard chỉ có 4 stats cards đơn giản (total events, capacity) — thiếu revenue, bookings, status breakdown
- Event publish flow hoạt động nhưng không có validation (có thể publish event chưa có zones/pricing)
- ~15 gateway endpoints tồn tại nhưng chưa được frontend sử dụng
- Event detail page thiếu event images (chỉ có gradient placeholder)
- Duplicate event endpoint có nhưng chưa wire vào UI

**Mục tiêu:** Nâng cấp Org Dashboard thành analytics dashboard thực sự, cải thiện publishing UX với validation, và wire các gateway endpoints đang bỏ phí.

---

## Implementation

### Step 1: Backend — Org Dashboard Stats Endpoint

**Modify** `event-service/` — Thêm gRPC method `GetEventStats`

Stats cần trả về:
```
total_events, published_events, draft_events, cancelled_events
total_capacity, total_zones, total_seats
```

**Modify** `gateway/src/handlers/eventHandlers.js` — Implement `getOrganizationDashboardHandler` (hiện trả 501)

Aggregated response:
```json
{
  "event_stats": { "total": 10, "published": 5, "draft": 3, "cancelled": 2 },
  "capacity_stats": { "total_capacity": 50000, "total_zones": 25, "total_seats": 40000 },
  "recent_events": [...]
}
```

### Step 2: Frontend — Dashboard Stats API Hook

**Modify** `frontend/src/lib/api/types/event.ts` — Add `OrgDashboardStats` interface
**Modify** `frontend/src/lib/api/endpoints.ts` — Add org dashboard endpoint
**Modify** `frontend/src/lib/api/queries/event.queries.ts` — Add `useOrgDashboardStats()` hook
**Modify** `frontend/src/lib/api/queries/index.ts` — Export new hook

### Step 3: Frontend — Org Dashboard Redesign

**Modify** `frontend/src/app/(org)/org/dashboard/page.tsx`

New layout:
```
[ Total Events ] [ Published ] [ Draft ] [ Cancelled ]    ← Stats cards with icons + colors
[ Recent Events table with status badges + quick actions ] ← Existing table enhanced
[ Quick Actions: Create Event | View All Events ]          ← CTA buttons
```

- Replace client-side calculation with `useOrgDashboardStats()` hook
- Color-coded stats cards (green=published, yellow=draft, red=cancelled)
- Recent events table with StatusBadge, quick edit/view links
- Empty state khi chưa có events

### Step 4: Frontend — Publish Validation

**Modify** `frontend/src/app/(org)/org/events/[id]/page.tsx`

Before publish, check:
1. Event có ít nhất 1 zone? (fetch zones)
2. Event có ít nhất 1 pricing tier? (fetch pricing)
3. Start date > now?

Show validation dialog nếu thiếu:
```
"Cannot publish: Missing zones/pricing"
- [ ] At least 1 seating zone required
- [x] At least 1 pricing tier required
- [x] Start date is in the future
[Set up Zones] [Cancel]
```

### Step 5: Frontend — Duplicate Event

**Modify** `frontend/src/app/(org)/org/events/[id]/page.tsx`

- Add "Duplicate" button trong header actions (bên cạnh Edit/Delete)
- Call `POST /events/:id/duplicate` endpoint
- Redirect to new event edit page sau khi duplicate
- Toast success: "Event duplicated as draft"

**Modify** `frontend/src/lib/api/endpoints.ts` — Add duplicate endpoint
**Modify** `frontend/src/lib/api/queries/event.queries.ts` — Add `useDuplicateEvent()` mutation

### Step 6: Frontend — Org Events Table Enhancement

**Modify** `frontend/src/app/(org)/org/events/page.tsx`

- Add status filter dropdown (FilterBar slot: draft/published/cancelled/completed)
- Add event_type filter dropdown
- Quick action dropdown per row: View | Edit | Duplicate | Delete
- Status column with StatusBadge

---

## Files Summary

| # | File | Action | Step |
|---|------|--------|------|
| 1 | `event-service/` (proto + controller + service + repo) | **Modify** | 1 |
| 2 | `gateway/src/handlers/eventHandlers.js` | **Modify** | 1 |
| 3 | `frontend/src/lib/api/types/event.ts` | **Modify** | 2 |
| 4 | `frontend/src/lib/api/endpoints.ts` | **Modify** | 2, 5 |
| 5 | `frontend/src/lib/api/queries/event.queries.ts` | **Modify** | 2, 5 |
| 6 | `frontend/src/lib/api/queries/index.ts` | **Modify** | 2, 5 |
| 7 | `frontend/src/app/(org)/org/dashboard/page.tsx` | **Modify** | 3 |
| 8 | `frontend/src/app/(org)/org/events/[id]/page.tsx` | **Modify** | 4, 5 |
| 9 | `frontend/src/app/(org)/org/events/page.tsx` | **Modify** | 6 |

**0 create, 9 modify = 9 files total** (+ backend proto/Go changes in Step 1)

---

## Key Decisions

1. **Dashboard stats from backend** — Không tính client-side nữa, dùng proper aggregation query
2. **Publish validation client-side** — Fetch zones/pricing trước khi cho publish, không cần backend validation endpoint riêng
3. **Duplicate as draft** — Event duplicate luôn tạo dưới dạng draft, user phải review + publish lại
4. **Reuse existing FilterBar** — Org events page dùng FilterBar molecule đã có từ Phase 9

---

## Verification

1. `/org/dashboard` → 4 color-coded stats cards, recent events table, loading skeleton
2. `/org/events` → Status + event_type filters, quick action dropdown per row
3. `/org/events/:id` → Publish button shows validation dialog nếu thiếu zones/pricing
4. `/org/events/:id` → Duplicate button creates new draft event + redirect
5. `npx tsc --noEmit` pass
6. `go build ./...` pass cho event-service

---

---

# Future Phases Roadmap

## Phase 13: Event Image Upload & Media Management
- Cloud storage integration (S3/Cloudflare R2)
- Image upload component (drag & drop, preview, crop)
- Event hero banner + gallery images
- Image optimization (resize, compress, WebP)
- Gateway multipart/form-data endpoint
- Public event detail page hero section

## Phase 14: Real-time Notifications (WebSocket)
- Frontend WebSocket client (connect to realtime-service port 50070)
- Booking status notifications (confirmed, cancelled, expired)
- Payment status notifications
- Reconnection logic with exponential backoff
- Notification center UI (bell icon dropdown)
- Unread count badge

## Phase 15: Booking Flow Enhancement
- Seat selection UI (interactive seat map from canvas_config)
- Real-time seat availability (WebSocket updates)
- Reservation countdown timer
- Queue position display
- Booking confirmation page redesign

## Phase 16: Admin Panel Enhancement
- Wire unused admin endpoints (payments admin list, refund management, permission CRUD)
- User management (activate/deactivate, role change)
- System-wide analytics dashboard
- Audit log viewer

## Phase 17: Notification Service & Email
- Notification service implementation (Go)
- Email templates (booking confirmation, payment receipt, event reminder)
- In-app notification history
- Notification preferences (email, in-app, push)

## Phase 18: Search & Discovery
- Full-text search (PostgreSQL tsvector or Elasticsearch)
- Event recommendations (category-based, location-based)
- Featured events carousel on homepage
- Event tags filtering
- Nearby events (geolocation)

## Phase 19: Performance & Production Readiness
- API response caching (Redis)
- Image CDN integration
- Bundle optimization (code splitting, lazy loading)
- Error tracking (Sentry)
- Health check dashboard
- Rate limiting
- CORS production config
