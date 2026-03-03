# Frontend Pages & Routes

Tech stack: **Next.js 14+ (App Router)**, TypeScript, shadcn/ui + Tailwind CSS, TanStack Query v5, Zustand, Konva.js (seat map), Stripe.js, WebSocket.

Nằm tại `frontend/` trong monorepo.

---

## Table of Contents

- [Project Structure](#project-structure)
- [Layouts](#layouts)
- [A. Public Pages](#a-public-pages)
- [B. Auth Pages](#b-auth-pages)
- [C. User Pages (Individual)](#c-user-pages-individual)
- [D. Organization Pages](#d-organization-pages)
- [E. Admin Pages](#e-admin-pages)
- [Key Flows](#key-flows)

---

## Project Structure

```
frontend/src/app/
├── (public)/                    # Public layout (Navbar + Footer)
│   ├── page.tsx                 # Landing page          /
│   ├── events/
│   │   ├── page.tsx             # Event listing          /events
│   │   └── [eventId]/
│   │       └── page.tsx         # Event detail           /events/:eventId
│   └── layout.tsx
│
├── (auth)/                      # Auth layout (centered card, no navbar)
│   ├── login/page.tsx           #                        /login
│   ├── register/page.tsx        #                        /register
│   ├── forgot-password/page.tsx #                        /forgot-password
│   ├── reset-password/page.tsx  #                        /reset-password
│   ├── verify-email/page.tsx    #                        /verify-email
│   └── layout.tsx
│
├── (dashboard)/                 # User dashboard layout (Auth Navbar + content)
│   ├── booking/
│   │   └── [eventId]/
│   │       ├── seats/page.tsx       #                    /booking/:eventId/seats
│   │       ├── checkout/page.tsx    #                    /booking/:eventId/checkout
│   │       ├── payment/page.tsx     #                    /booking/:eventId/payment
│   │       └── confirmation/page.tsx#                    /booking/:eventId/confirmation
│   ├── my-bookings/
│   │   ├── page.tsx                 #                    /my-bookings
│   │   └── [bookingId]/page.tsx     #                    /my-bookings/:bookingId
│   ├── my-tickets/
│   │   ├── page.tsx                 #                    /my-tickets
│   │   └── [ticketId]/page.tsx      #                    /my-tickets/:ticketId
│   ├── payments/
│   │   ├── page.tsx                 #                    /payments
│   │   └── [paymentId]/page.tsx     #                    /payments/:paymentId
│   ├── profile/
│   │   ├── page.tsx                 #                    /profile
│   │   └── addresses/page.tsx       #                    /profile/addresses
│   └── layout.tsx
│
├── (org)/                       # Organization layout (Sidebar + Topbar)
│   ├── dashboard/page.tsx           #                    /org/dashboard
│   ├── events/
│   │   ├── page.tsx                 #                    /org/events
│   │   ├── create/page.tsx          #                    /org/events/create
│   │   └── [eventId]/
│   │       ├── edit/page.tsx        #                    /org/events/:eventId/edit
│   │       ├── zones/page.tsx       #                    /org/events/:eventId/zones
│   │       ├── seats/page.tsx       #                    /org/events/:eventId/seats
│   │       ├── pricing/page.tsx     #                    /org/events/:eventId/pricing
│   │       ├── availability/page.tsx#                    /org/events/:eventId/availability
│   │       └── ticket-types/page.tsx#                    /org/events/:eventId/ticket-types
│   └── layout.tsx
│
├── (admin)/                     # Admin layout (Sidebar + Topbar)
│   ├── dashboard/page.tsx           #                    /admin/dashboard
│   ├── users/
│   │   ├── page.tsx                 #                    /admin/users
│   │   └── [userId]/page.tsx        #                    /admin/users/:userId
│   ├── bookings/page.tsx            #                    /admin/bookings
│   ├── payments/
│   │   ├── page.tsx                 #                    /admin/payments
│   │   └── refunds/page.tsx         #                    /admin/payments/refunds
│   └── layout.tsx
│
└── layout.tsx                   # Root layout
```

---

## Layouts

### 1. Public Layout

Dùng cho: Landing, Event Listing, Event Detail.

```
+------------------------------------------------------------------+
|  [Logo]   Events                     [Login]  [Register]          |
+------------------------------------------------------------------+
|                                                                    |
|                        PAGE CONTENT                                |
|                                                                    |
+------------------------------------------------------------------+
|  Footer: About | Contact | Terms | Privacy          © 2026        |
+------------------------------------------------------------------+
```

- Navbar: Logo, link Events, Login/Register buttons (hoặc avatar dropdown nếu đã login)
- Footer: Links, copyright
- Responsive: hamburger menu trên mobile

### 2. Dashboard Layout (Individual User)

Dùng cho: Booking flow, My Bookings, My Tickets, Payments, Profile. Yêu cầu đăng nhập.

```
+------------------------------------------------------------------+
|  [Logo]  Events  My Bookings  My Tickets     [Avatar v] [Bell]   |
+------------------------------------------------------------------+
|                                                                    |
|                        PAGE CONTENT                                |
|                                                                    |
+------------------------------------------------------------------+
```

- Navbar: Logo, navigation links (Events, My Bookings, My Tickets), notification bell (WebSocket), avatar dropdown (Profile, Payments, Settings, Logout)
- Không có footer (app-like)
- Route guard: redirect về `/login` nếu chưa auth

### 3. Organization Layout

Dùng cho: Org Dashboard, Event Management. Yêu cầu role `organization`.

```
+------------------+-----------------------------------------------+
|                  |  Organization Name          [Avatar v] [Bell]  |
|  [Logo]          +-----------------------------------------------+
|                  |                                                 |
|  Dashboard       |                                                 |
|  Events          |              PAGE CONTENT                       |
|  Create Event    |                                                 |
|  Analytics       |                                                 |
|  Settings        |                                                 |
|                  |                                                 |
+------------------+-----------------------------------------------+
```

- Sidebar: Logo, navigation items, collapsible trên mobile
- Topbar: Org name, avatar dropdown, notifications
- Route guard: redirect về `/login` nếu chưa auth, về `/` nếu không phải role `organization`

### 4. Admin Layout

Dùng cho: Admin Dashboard, User/Booking/Payment Management. Yêu cầu role `admin`.

```
+------------------+-----------------------------------------------+
|                  |  Admin Panel                [Avatar v] [Bell]  |
|  [Logo]          +-----------------------------------------------+
|                  |                                                 |
|  Dashboard       |                                                 |
|  Users           |              PAGE CONTENT                       |
|  Bookings        |                                                 |
|  Payments        |                                                 |
|    Refunds       |                                                 |
|                  |                                                 |
+------------------+-----------------------------------------------+
```

- Tương tự Organization Layout nhưng navigation items khác
- Route guard: chỉ cho role `admin` hoặc `super_admin`

---

## A. Public Pages

### A1. Landing Page

| | |
|---|---|
| **Route** | `/` |
| **Layout** | Public |
| **Mô tả** | Trang chủ. Hero section với search bar, featured events carousel, upcoming events grid, browse by category. |

**Components:**
- `HeroSection` - Banner lớn với search bar (keyword, location, date)
- `FeaturedEventsCarousel` - Horizontal scroll các event nổi bật
- `UpcomingEventsGrid` - Grid 3-4 cột các event sắp diễn ra
- `CategoryGrid` - Chips/cards: Music, Sports, Theater, Comedy, Conferences, Festivals

**API:**
- `GET /api/events?status=published&is_featured=true&limit=8` - Featured events
- `GET /api/events?status=published&sort=start_date&limit=12` - Upcoming events

```
+------------------------------------------------------------------+
|  NAVBAR                                                           |
+------------------------------------------------------------------+
|                                                                    |
|              Find Your Next Experience                             |
|   [______Search events______] [City v] [Date] [Search]            |
|                                                                    |
+------------------------------------------------------------------+
|  Featured Events                                    [See all ->]   |
|  +----------+  +----------+  +----------+  +----------+           |
|  |  Image   |  |  Image   |  |  Image   |  |  Image   |          |
|  |  Name    |  |  Name    |  |  Name    |  |  Name    |          |
|  | Date/Loc |  | Date/Loc |  | Date/Loc |  | Date/Loc |          |
|  | From $XX |  | From $XX |  | From $XX |  | From $XX |          |
|  +----------+  +----------+  +----------+  +----------+           |
|                                                                    |
|  Browse by Category                                                |
|  [Music] [Sports] [Theater] [Comedy] [Conferences] [Festivals]    |
|                                                                    |
+------------------------------------------------------------------+
|  FOOTER                                                            |
+------------------------------------------------------------------+
```

---

### A2. Event Listing

| | |
|---|---|
| **Route** | `/events` |
| **Layout** | Public |
| **Mô tả** | Danh sách events có filter, search, sort, pagination. Hỗ trợ grid/list view toggle. |

**Components:**
- `EventFilterSidebar` - Filters: category, date range, location, price range
- `EventGrid` / `EventList` - Toggle grid/list view
- `EventCard` - Card: image, name, date, venue, starting price
- `SortDropdown` - Sort by: date, price, popularity
- `PaginationControls`

**API:**
- `GET /api/events?status=published&category=X&location=X&date_from=X&date_to=X&price_min=X&price_max=X&sort=X&page=X&limit=12`

```
+------------------------------------------------------------------+
|  Events                                        [Grid] [List]      |
+----------+-------------------------------------------------------+
| FILTERS  |  Sort: [Date v]     Showing 1-12 of 156               |
|          |                                                         |
| Category |  +----------+  +----------+  +----------+              |
| [x]Music |  | Image    |  | Image    |  | Image    |             |
| [ ]Sport |  | Name     |  | Name     |  | Name     |             |
|          |  | Mar 15   |  | Mar 22   |  | Apr 1    |             |
| Date     |  | Venue    |  | Venue    |  | Venue    |             |
| [From]   |  | $45+     |  | $30+     |  | $55+     |             |
| [To  ]   |  +----------+  +----------+  +----------+              |
|          |                                                         |
| Location |  +----------+  +----------+  +----------+              |
| [City v] |  | ...      |  | ...      |  | ...      |             |
|          |  +----------+  +----------+  +----------+              |
| Price    |                                                         |
| [--o---] |       [<< 1  2  3  4  5 ... 13 >>]                    |
+----------+-------------------------------------------------------+
```

---

### A3. Event Detail

| | |
|---|---|
| **Route** | `/events/:eventId` |
| **Layout** | Public |
| **Mô tả** | Chi tiết event: thông tin, venue map preview (read-only canvas), pricing theo zone, ticket types, availability summary. CTA "Book Now" (redirect login nếu chưa auth). |

**Components:**
- `EventHero` - Image lớn, event name, date/time, venue
- `EventInfoTabs` - Tabs: Description, Venue Map, Pricing, Schedule
- `VenueMapPreview` - Konva.js canvas read-only hiển thị zones với màu sắc
- `AvailabilitySummary` - Cards: Total seats, Available, Reserved, Sold
- `PricingTable` - Bảng giá theo zone/category
- `TicketTypeSelector` - Danh sách ticket types với price và available quantity
- `BookNowButton` - CTA chính, redirect `/booking/:eventId/seats`

**API:**
- `GET /api/events/:eventId`
- `GET /api/events/:eventId/zones`
- `GET /api/events/:eventId/availability`
- `GET /api/events/:eventId/pricing`
- `GET /api/tickets/types/:eventId`
- `GET /api/tickets/availability/:eventId`

```
+------------------------------------------------------------------+
|  [< Back to Events]                                                |
|                                                                    |
|  +---------------------------+  +-------------------------------+  |
|  |                           |  |  CONCERT ABC                  |  |
|  |      Event Image          |  |  Mar 15, 2026 - 8:00 PM      |  |
|  |                           |  |  Madison Square Garden        |  |
|  |                           |  |  New York, US                 |  |
|  +---------------------------+  |                               |  |
|                                 |  Capacity: 20,000             |  |
|                                 |  Available: 12,450            |  |
|                                 |  Starting from $45.00         |  |
|                                 |                               |  |
|                                 |  [      BOOK NOW      ]       |  |
|                                 +-------------------------------+  |
|                                                                    |
|  [Description] [Venue Map] [Pricing] [Schedule]                    |
|  +--------------------------------------------------------------+  |
|  |  Tab content...                                               |  |
|  +--------------------------------------------------------------+  |
|                                                                    |
|  Ticket Types                                                      |
|  +------------------------+----------+----------+---------+        |
|  | Type                   | Price    | Avail.   | Action  |        |
|  | Standard               | $45.00   | 8,000    | [Select]|        |
|  | VIP                    | $150.00  | 500      | [Select]|        |
|  | Wheelchair Accessible  | $45.00   | 50       | [Select]|        |
|  +------------------------+----------+----------+---------+        |
+------------------------------------------------------------------+
```

---

## B. Auth Pages

Tất cả auth pages dùng layout centered card (không navbar, không footer).

```
+------------------------------------------------------------------+
|                                                                    |
|                        [Logo]                                      |
|                                                                    |
|                 +------------------------+                         |
|                 |                        |                         |
|                 |      AUTH FORM         |                         |
|                 |                        |                         |
|                 +------------------------+                         |
|                                                                    |
+------------------------------------------------------------------+
```

### B1. Login

| | |
|---|---|
| **Route** | `/login` |
| **Layout** | Auth (centered card) |
| **Mô tả** | Đăng nhập bằng email/password hoặc OAuth (Google, Facebook). |

**Components:**
- `LoginForm` - Email, password fields, remember me checkbox
- `OAuthButtons` - Google, Facebook buttons
- `ForgotPasswordLink` - Link đến `/forgot-password`
- `RegisterLink` - Link đến `/register`

**API:**
- `POST /api/auth/login` - Email/password login
- `POST /api/auth/oauth/login` - OAuth login

---

### B2. Register

| | |
|---|---|
| **Route** | `/register` |
| **Layout** | Auth (centered card) |
| **Mô tả** | Đăng ký tài khoản. Chọn role: "Mua vé" (individual) hoặc "Tổ chức sự kiện" (organization). |

**Components:**
- `RoleSelector` - 2 cards chọn role trước khi fill form
- `RegistrationForm` - First name, last name, email, password, confirm password
- `OrganizationFields` - Hiện thêm nếu chọn organization: org name, description, website, tax ID
- `OAuthButtons` - Google, Facebook
- `TermsCheckbox`
- `LoginLink` - Link đến `/login`

**API:**
- `POST /api/auth/register/email` - Email registration
- `POST /api/auth/register/oauth` - OAuth registration

---

### B3. Verify Email

| | |
|---|---|
| **Route** | `/verify-email` |
| **Layout** | Auth (centered card) |
| **Mô tả** | Nhập PIN 6 số từ email xác thực. |

**Components:**
- `PinCodeInput` - 6 ô input cho PIN
- `ResendButton` - Gửi lại email (có cooldown 60s)
- `EmailDisplay` - Hiển thị email đã đăng ký

**API:**
- `POST /api/auth/verify-user` - Verify email with PIN
- `POST /api/auth/resend-verification-email` - Resend PIN

---

### B4. Forgot Password

| | |
|---|---|
| **Route** | `/forgot-password` |
| **Layout** | Auth (centered card) |
| **Mô tả** | Nhập email để nhận link reset password. |

**Components:**
- `EmailInput`
- `SubmitButton`
- `BackToLoginLink`

**API:**
- `POST /api/auth/forgot-password`

---

### B5. Reset Password

| | |
|---|---|
| **Route** | `/reset-password?token=xxx` |
| **Layout** | Auth (centered card) |
| **Mô tả** | Đặt password mới bằng token từ email. |

**Components:**
- `NewPasswordForm` - New password, confirm password
- `PasswordStrengthIndicator`

**API:**
- `POST /api/auth/reset-password`

---

## C. User Pages (Individual)

Tất cả yêu cầu đăng nhập. Dùng Dashboard Layout.

### C1. Seat Selection

| | |
|---|---|
| **Route** | `/booking/:eventId/seats` |
| **Layout** | Dashboard |
| **Mô tả** | Interactive seat map (Konva.js canvas). User chọn ghế, xem giá realtime. Khi chọn xong → reserve seats → countdown 10 phút bắt đầu. Đây là page phức tạp nhất. |

**Components:**
- `BookingStepIndicator` - Steps: 1.Select Seats > 2.Checkout > 3.Payment > 4.Confirmation
- `InteractiveSeatMap` - Konva.js canvas render zones + seats từ `canvas_config` và `coordinates`. Color-coded theo zone/status. Click để chọn/bỏ ghế
- `ZoneLegend` - Chú thích màu cho từng zone
- `SeatStatusLegend` - Available (xanh), Reserved (vàng), Booked (đỏ), Blocked (xám), Selected (xanh dương)
- `SelectedSeatsPanel` - Panel bên phải: danh sách ghế đã chọn, giá từng ghế, tổng tiền
- `CountdownTimer` - Hiện sau khi reserve, đếm ngược 10 phút
- `ContinueButton` - Reserve seats → chuyển sang checkout
- `ZoomControls` - Zoom in/out/reset cho seat map

**API:**
- `GET /api/events/:eventId` - Event details + canvas_config
- `GET /api/events/:eventId/zones` - Zones với coordinates/colors
- `GET /api/events/:eventId/seats` - Tất cả seats với coordinates/status
- `GET /api/events/:eventId/pricing` - Pricing rules
- `GET /api/events/:eventId/availability` - Availability overview
- `POST /api/events/:eventId/pricing/calculate` - Tính giá cho selection
- `POST /api/bookings/seats/reserve` - Reserve ghế đã chọn

**WebSocket:**
- `room:join { event_id }` - Join room nhận updates realtime
- Listen `ticket:availability` - Update status ghế realtime
- Listen `ticket:reserved` / `ticket:released` - Ghế bị người khác reserve/release

```
+------------------------------------------------------------------+
|  Book: EVENT NAME                       Timer: 09:45 remaining    |
+------------------------------------------------------------------+
|  [1.Select Seats] > 2.Checkout > 3.Payment > 4.Confirmation      |
+------------------------------------------------------------------+
|                                         |  Your Selection          |
|  +-----------------------------------+ |                          |
|  |                                   | |  Seat A12 (VIP)   $150  |
|  |          STAGE                    | |  Seat A13 (VIP)   $150  |
|  |                                   | |  Seat B5 (Std)     $45  |
|  |     [Zone A - VIP]               | |                          |
|  |    o o o X o o o o o             | |  ----------------------   |
|  |    o o o o o * * o o             | |  Subtotal:       $345.00 |
|  |                                   | |  Service Fee:     $17.25 |
|  |     [Zone B - Standard]          | |  Total:          $362.25 |
|  |    o o o o o o o o o             | |                          |
|  |    o o o X o o * o o             | |  [Continue to Checkout]  |
|  |    o o o o o o o o o             | |                          |
|  |                                   | |  Legend:                 |
|  +-----------------------------------+ |  o Available  * Selected |
|                                         |  X Booked     - Blocked |
|  Zoom: [+] [-] [Reset]                 |                          |
+------------------------------------------------------------------+
```

---

### C2. Checkout

| | |
|---|---|
| **Route** | `/booking/:eventId/checkout` |
| **Layout** | Dashboard |
| **Mô tả** | Review booking: danh sách ghế, pricing breakdown, optional discount code, special requests. Tạo booking object. |

**Components:**
- `BookingStepIndicator` - Step 2 active
- `BookingSummaryCard` - Event info, danh sách ghế đã reserve
- `PricingBreakdown` - Base price, service fee, discount, total
- `DiscountCodeInput` - Nhập mã giảm giá
- `SpecialRequestsTextarea` - Ghi chú đặc biệt
- `CountdownTimer` - Countdown reservation
- `ProceedToPaymentButton`

**API:**
- `POST /api/bookings` - Tạo booking (event_id, ticket_quantity, seat_numbers, special_requests)
- `POST /api/events/:eventId/pricing/calculate` - Tính lại giá nếu có discount
- `POST /api/events/:eventId/pricing/discount` - Validate discount code

**WebSocket:**
- Listen `booking:queue_position` - Vị trí trong queue
- Listen `booking:processing` - Booking đang được xử lý

---

### C3. Payment

| | |
|---|---|
| **Route** | `/booking/:eventId/payment` |
| **Layout** | Dashboard |
| **Mô tả** | Form thanh toán. Chọn payment method, nhập thẻ qua Stripe Elements, chọn billing address. |

**Components:**
- `BookingStepIndicator` - Step 3 active
- `OrderSummary` - Compact: event name, seat count, total
- `PaymentMethodSelector` - Radio: Credit Card, Debit Card, E-Wallet, Bank Transfer
- `StripePaymentForm` - Stripe Elements: card number, expiry, CVC
- `SavedPaymentMethods` - Danh sách thẻ đã lưu
- `BillingAddressSelector` - Dropdown chọn address đã lưu hoặc nhập mới
- `CountdownTimer`
- `PayButton` - "Pay $362.25"
- `PaymentProcessingOverlay` - Loading overlay trong lúc xử lý

**API:**
- `GET /api/payments/methods` - Payment methods đã lưu
- `GET /api/users/addresses` - Billing addresses
- `POST /api/payments` - Process payment (booking_id, amount, currency, payment_method)
- `POST /api/payments/:paymentId/capture` - Capture sau khi Stripe confirm
- `POST /api/bookings/:bookingId/confirm` - Confirm booking sau payment success

**WebSocket:**
- Listen `payment:processing` - Payment đang xử lý
- Listen `payment:success` - Thanh toán thành công → redirect confirmation
- Listen `payment:failed` - Thanh toán thất bại → hiện error, cho retry

```
+------------------------------------------------------------------+
|  1.Select Seats > 2.Checkout > [3.Payment] > 4.Confirmation      |
+------------------------------------------------------------------+
|                                         |  Order Summary           |
|  Payment Method                         |  EVENT NAME              |
|  (*) Credit Card                        |  Mar 15, 2026 8:00 PM   |
|  ( ) Debit Card                         |                          |
|  ( ) E-Wallet                           |  3x Seats: $345.00      |
|                                         |  Service Fee: $17.25    |
|  Card Details                           |  ----------------------  |
|  +-----------------------------------+ |  Total: $362.25         |
|  | Card Number                       | |                          |
|  | [4242 4242 4242 4242            ] | |  Timer: 07:22           |
|  | Expiry        CVC                 | |                          |
|  | [MM/YY]       [123]               | |                          |
|  +-----------------------------------+ |                          |
|                                         |                          |
|  Billing Address                        |                          |
|  [v] Home - 123 Main St, NYC           |                          |
|                                         |                          |
|  [          PAY $362.25          ]      |                          |
+------------------------------------------------------------------+
```

---

### C4. Booking Confirmation

| | |
|---|---|
| **Route** | `/booking/:eventId/confirmation` |
| **Layout** | Dashboard |
| **Mô tả** | Trang success. Hiển thị booking reference, ticket cards với QR code, download/share options. |

**Components:**
- `SuccessAnimation` - Checkmark animation
- `BookingReference` - Hiển thị lớn: "BK-20260315-A7X9F2"
- `TicketCards` - Mỗi ticket 1 card: seat info, QR code (dùng `qrcode.react`)
- `DownloadTicketsButton` - Download PDF
- `AddToCalendarButton` - iCal/Google Calendar
- `ShareButton`
- `ViewBookingsLink` - Link đến `/my-bookings`

**API:**
- `GET /api/bookings/:bookingId` - Booking details
- `GET /api/tickets` - Tickets của booking này

```
+------------------------------------------------------------------+
|  1.Select Seats > 2.Checkout > 3.Payment > [4.Confirmation]      |
+------------------------------------------------------------------+
|                                                                    |
|                     [checkmark]                                    |
|                Booking Confirmed!                                  |
|                                                                    |
|         Reference: BK-20260315-A7X9F2                             |
|                                                                    |
|   +-------------------+  +-------------------+  +--------------+  |
|   | TICKET #TK-001    |  | TICKET #TK-002    |  | TICKET #003  |  |
|   | Seat: A12 (VIP)   |  | Seat: A13 (VIP)   |  | Seat: B5     |  |
|   | $150.00           |  | $150.00           |  | $45.00       |  |
|   |  +------+         |  |  +------+         |  |  +------+    |  |
|   |  | QR   |         |  |  | QR   |         |  |  | QR   |    |  |
|   |  | Code |         |  |  | Code |         |  |  | Code |    |  |
|   |  +------+         |  |  +------+         |  |  +------+    |  |
|   +-------------------+  +-------------------+  +--------------+  |
|                                                                    |
|   [Download Tickets]  [Add to Calendar]  [Share]                   |
|   [View All My Bookings]                                           |
+------------------------------------------------------------------+
```

---

### C5. My Bookings

| | |
|---|---|
| **Route** | `/my-bookings` |
| **Layout** | Dashboard |
| **Mô tả** | Danh sách tất cả bookings của user. Filter theo status, phân trang. |

**Components:**
- `BookingStatusTabs` - All, Pending, Confirmed, Cancelled, Expired
- `BookingCard` - Event name, date, status badge, seat count, total amount, actions (View, Cancel)
- `EmptyState` - "No bookings yet" với CTA browse events
- `PaginationControls`

**API:**
- `GET /api/bookings?status=X&page=X&limit=10`

---

### C6. Booking Detail

| | |
|---|---|
| **Route** | `/my-bookings/:bookingId` |
| **Layout** | Dashboard |
| **Mô tả** | Chi tiết 1 booking: timeline trạng thái, thông tin event, danh sách ghế, payment info, tickets. Actions: cancel, request refund. |

**Components:**
- `BookingStatusBadge` - Badge màu theo status
- `BookingTimeline` - Visual timeline: Created → Confirmed → (Cancelled?)
- `EventInfoCard` - Event name, date, venue
- `SeatList` - Danh sách ghế: zone, row, seat number, price
- `PaymentInfoCard` - Amount, method, status, transaction date
- `TicketList` - Ticket cards với QR
- `CancelBookingButton` - Với confirmation dialog
- `RequestRefundButton` - Nếu đã thanh toán

**API:**
- `GET /api/bookings/:bookingId`
- `GET /api/tickets` - Tickets cho booking này
- `GET /api/payments` - Payment info
- `POST /api/bookings/:bookingId/cancel`
- `POST /api/payments/:paymentId/refund`

---

### C7. My Tickets

| | |
|---|---|
| **Route** | `/my-tickets` |
| **Layout** | Dashboard |
| **Mô tả** | Danh sách tất cả tickets. Filter theo status. |

**Components:**
- `TicketStatusTabs` - All, Active, Used, Cancelled, Refunded
- `TicketCard` - Event name, seat info, QR code thumbnail, status badge, validity dates
- `PaginationControls`

**API:**
- `GET /api/tickets?status=X&page=X&limit=10`

---

### C8. Ticket Detail

| | |
|---|---|
| **Route** | `/my-tickets/:ticketId` |
| **Layout** | Dashboard |
| **Mô tả** | Chi tiết 1 ticket. QR code lớn để scan tại venue. |

**Components:**
- `TicketDisplay` - Full ticket card layout
- `QRCodeFull` - QR code lớn (dùng `qrcode.react`)
- `TicketInfo` - Event, seat, zone, pricing, validity, status
- `DownloadButton` - Download ticket PDF

**API:**
- `GET /api/tickets/:ticketId`

---

### C9. Profile

| | |
|---|---|
| **Route** | `/profile` |
| **Layout** | Dashboard |
| **Mô tả** | Xem/sửa thông tin cá nhân. |

**Components:**
- `AvatarUploader` - Upload/change avatar
- `ProfileForm` - First name, last name, phone, date of birth
- `PreferencesSection` - Notification preferences
- `ManageAddressesLink` - Link đến `/profile/addresses`

**API:**
- `GET /api/users/profile`
- `PUT /api/users/profile`
- `POST /api/users/profile` (create nếu chưa có)

---

### C10. Address Management

| | |
|---|---|
| **Route** | `/profile/addresses` |
| **Layout** | Dashboard |
| **Mô tả** | CRUD địa chỉ (dùng cho billing). |

**Components:**
- `AddressList` - Danh sách addresses, đánh dấu default
- `AddressCard` - Label (Home/Work/Other), street, city, state, country, postal code
- `AddressFormDialog` - Modal add/edit address
- `SetDefaultButton`
- `DeleteButton` - Với confirmation

**API:**
- `GET /api/users/addresses`
- `POST /api/users/addresses`
- `PUT /api/users/addresses/:addressId`
- `DELETE /api/users/addresses/:addressId`

---

### C11. Payment History

| | |
|---|---|
| **Route** | `/payments` |
| **Layout** | Dashboard |
| **Mô tả** | Lịch sử thanh toán. |

**Components:**
- `PaymentStatusFilter` - All, Success, Pending, Failed, Refunded
- `PaymentTable` - Date, amount, status badge, booking reference, method
- `PaginationControls`

**API:**
- `GET /api/payments?status=X&page=X&limit=10`

---

### C12. Payment Detail

| | |
|---|---|
| **Route** | `/payments/:paymentId` |
| **Layout** | Dashboard |
| **Mô tả** | Chi tiết 1 payment. Lịch sử refund nếu có. |

**Components:**
- `PaymentSummary` - Amount, method, status, date, booking reference
- `RefundHistory` - Danh sách refunds (partial/full) với status
- `RequestRefundButton` - Nếu payment success và chưa fully refunded
- `RefundDialog` - Nhập amount, reason

**API:**
- `GET /api/payments/:paymentId`
- `GET /api/payments/:paymentId/refunds`
- `POST /api/payments/:paymentId/refund`

---

## D. Organization Pages

Tất cả yêu cầu role `organization`. Dùng Organization Layout.

### D1. Organization Dashboard

| | |
|---|---|
| **Route** | `/org/dashboard` |
| **Layout** | Organization |
| **Mô tả** | Tổng quan: số events, tickets bán được, revenue, upcoming events. Charts. |

**Components:**
- `StatsCards` - 4 cards: Total Events, Tickets Sold, Revenue, Upcoming
- `RevenueChart` - Line chart revenue theo thời gian (dùng `recharts`)
- `SalesChart` - Bar chart ticket sales theo event
- `RecentEventsTable` - Events gần nhất: name, status, sold, revenue, action
- `QuickActions` - Buttons: Create Event, View All Events

**API:**
- `GET /api/organizations/:orgId/dashboard`
- `GET /api/events?organization_id=X`

```
+------------------------------------------------------------------+
|  SIDEBAR  |  Dashboard                                            |
|           |                                                        |
|           |  +--------+ +--------+ +--------+ +--------+          |
|           |  |  12    | | 4,520  | | $98.5K | |   3    |          |
|           |  | Events | | Sold   | | Revenue| |Upcoming|          |
|           |  +--------+ +--------+ +--------+ +--------+          |
|           |                                                        |
|           |  Revenue (Last 30 Days)                                |
|           |  +--------------------------------------------------+ |
|           |  |       /\         /\                               | |
|           |  |      /  \       /  \      /\                      | |
|           |  |     /    \     /    \    /  \                     | |
|           |  |    /      \   /      \  /    \                    | |
|           |  +--------------------------------------------------+ |
|           |                                                        |
|           |  Recent Events                                         |
|           |  +--------+--------+------+--------+--------+         |
|           |  | Name   | Status | Sold | Revenue| Action |         |
|           |  | Evt 1  | Live   | 1200 | $54K   | [View] |         |
|           |  | Evt 2  | Draft  | --   | --     | [Edit] |         |
|           |  +--------+--------+------+--------+--------+         |
+------------------------------------------------------------------+
```

---

### D2. Events List

| | |
|---|---|
| **Route** | `/org/events` |
| **Layout** | Organization |
| **Mô tả** | Danh sách events của organization. Filter theo status, search, actions. |

**Components:**
- `EventStatusTabs` - All, Draft, Published, Cancelled, Completed
- `EventTable` - Name, date, venue, status badge, capacity, sold count, actions
- `CreateEventButton` - Link đến `/org/events/create`
- `EventActionsMenu` - Per-row: Edit, Duplicate, Publish/Unpublish, Delete
- `SearchBar`

**API:**
- `GET /api/events` (filtered by organization)
- `POST /api/events/:eventId/publish`
- `POST /api/events/:eventId/duplicate`
- `DELETE /api/events/:eventId`

---

### D3. Create Event (Multi-step Wizard)

| | |
|---|---|
| **Route** | `/org/events/create` |
| **Layout** | Organization |
| **Mô tả** | Wizard tạo event mới. 6 bước: Basic Info → Venue → Zones → Seats → Pricing → Ticket Types. Có thể save draft bất kỳ lúc nào. |

**Steps & Components:**

**Step 1 - Basic Info:**
- `BasicInfoForm` - Name, description (rich text), start/end date+time, event type, category, tags, images upload
- `EventTemplateSelector` - Chọn template có sẵn để bắt đầu

**Step 2 - Venue:**
- `VenueForm` - Venue name, address, city, country, capacity

**Step 3 - Zones:**
- `ZoneBuilder` - Tạo zones: name, type (seated/standing/VIP), color picker, capacity
- `ZoneCanvasPreview` - Preview layout zones trên canvas

**Step 4 - Seats:**
- `SeatMapEditor` - Konva.js canvas editor: drag-drop đặt ghế vào zone
- `BulkSeatCreator` - Tạo hàng loạt: chọn zone, số hàng, số ghế/hàng → auto generate
- `SeatTable` - Bảng tất cả seats, inline editing

**Step 5 - Pricing:**
- `PricingConfigurator` - Set giá theo zone/category: base_price, currency, discount rules, valid from/until

**Step 6 - Ticket Types:**
- `TicketTypeCreator` - Tạo ticket types: name, price, quantity, max/min per purchase
- `ReviewSummary` - Tổng hợp tất cả settings

**Navigation:**
- `WizardStepper` - Step indicator, có thể quay lại bước trước
- `SaveDraftButton` - Lưu nháp bất kỳ lúc nào
- `PublishButton` - Publish event (ở bước cuối)

**API:**
- `GET /api/events/templates` - Load templates
- `POST /api/events` - Create event (draft)
- `PUT /api/events/:eventId/draft` - Save draft
- `POST /api/events/:eventId/zones` - Create zones
- `POST /api/events/:eventId/seats/bulk` - Bulk create seats
- `POST /api/events/:eventId/pricing` - Create pricing
- `POST /api/tickets/types` - Create ticket types
- `POST /api/events/:eventId/publish` - Publish

---

### D4. Edit Event

| | |
|---|---|
| **Route** | `/org/events/:eventId/edit` |
| **Layout** | Organization |
| **Mô tả** | Edit event info (same form như Step 1-2 của Create, nhưng pre-filled). |

**Components:** Reuse `BasicInfoForm` + `VenueForm` từ Create Event.

**API:**
- `GET /api/events/:eventId`
- `PUT /api/events/:eventId`
- `PUT /api/events/:eventId/draft`

---

### D5. Zone Management

| | |
|---|---|
| **Route** | `/org/events/:eventId/zones` |
| **Layout** | Organization |
| **Mô tả** | CRUD zones cho event. Bảng + canvas preview. |

**Components:**
- `ZoneTable` - Name, type, seat count, color swatch, actions
- `ZoneFormDialog` - Modal create/edit: name, type (seated/standing/VIP), color picker
- `ZoneCanvasPreview` - Konva.js preview zones
- `DeleteZoneConfirmation`

**API:**
- `GET /api/events/:eventId/zones`
- `POST /api/events/:eventId/zones`
- `PUT /api/events/:eventId/zones/:zoneId`
- `DELETE /api/events/:eventId/zones/:zoneId`

---

### D6. Seat Map Editor

| | |
|---|---|
| **Route** | `/org/events/:eventId/seats` |
| **Layout** | Organization |
| **Mô tả** | Interactive canvas editor để đặt/quản lý ghế trong zones. Hỗ trợ bulk creation. |

**Components:**
- `SeatMapCanvas` - Konva.js full editor: drag-drop, select, multi-select
- `BulkSeatCreator` - Form: zone, start row, end row, seats per row → generate
- `SeatPropertyPanel` - Edit selected seat: number, row, pricing category, status
- `SeatTable` - Bảng tất cả seats, filter by zone
- `ZoomControls` + `SnapToGrid`

**API:**
- `GET /api/events/:eventId/seats`
- `GET /api/events/:eventId/zones`
- `POST /api/events/:eventId/seats`
- `POST /api/events/:eventId/seats/bulk`
- `PUT /api/events/:eventId/seats/:seatId`
- `DELETE /api/events/:eventId/seats/:seatId`

---

### D7. Pricing Management

| | |
|---|---|
| **Route** | `/org/events/:eventId/pricing` |
| **Layout** | Organization |
| **Mô tả** | CRUD pricing rules theo zone/category. Test price calculation. |

**Components:**
- `PricingTable` - Zone, category, base price, currency, valid dates, active status
- `PricingFormDialog` - Create/edit: zone selector, category selector, price, discount rules, validity
- `DiscountRuleBuilder` - Visual builder cho discount rules (JSON)
- `PriceCalculatorPreview` - Test tính giá cho 1 selection

**API:**
- `GET /api/events/:eventId/pricing`
- `POST /api/events/:eventId/pricing`
- `PUT /api/events/:eventId/pricing/:pricingId`
- `DELETE /api/events/:eventId/pricing/:pricingId`
- `GET /api/events/:eventId/pricing/zone/:zoneId`
- `POST /api/events/:eventId/pricing/calculate`
- `POST /api/events/:eventId/pricing/discount`

---

### D8. Availability Management

| | |
|---|---|
| **Route** | `/org/events/:eventId/availability` |
| **Layout** | Organization |
| **Mô tả** | Xem tổng quan availability. Block/release seats. |

**Components:**
- `AvailabilitySummaryCards` - Total, Available, Reserved, Booked, Blocked
- `AvailabilityMap` - Konva.js canvas read-only, color-coded theo status, hỗ trợ multi-select để block
- `ZoneAvailabilityTable` - Breakdown theo zone
- `BlockSeatsDialog` - Chọn seats + reason → block
- `ReleaseSeatsDialog` - Chọn blocked seats → release

**API:**
- `GET /api/events/:eventId/availability`
- `GET /api/events/:eventId/availability/zones/:zoneId`
- `GET /api/events/:eventId/availability/seats/:seatId`
- `PUT /api/events/:eventId/availability/seats/:seatId`
- `POST /api/events/:eventId/availability/block`
- `POST /api/events/:eventId/availability/release`

---

### D9. Ticket Types Management

| | |
|---|---|
| **Route** | `/org/events/:eventId/ticket-types` |
| **Layout** | Organization |
| **Mô tả** | CRUD ticket types cho event. |

**Components:**
- `TicketTypeTable` - Name, price, total quantity, available quantity, min/max per purchase, status
- `TicketTypeFormDialog` - Create/edit: name, description, price, currency, quantity, max/min per purchase
- `AvailabilityBadge` - Visual indicator remaining/total

**API:**
- `GET /api/tickets/types/:eventId`
- `GET /api/tickets/availability/:eventId`
- `POST /api/tickets/types`
- `PUT /api/tickets/types/:typeId`
- `DELETE /api/tickets/types/:typeId`

---

## E. Admin Pages

Tất cả yêu cầu role `admin` hoặc `super_admin`. Dùng Admin Layout.

### E1. Admin Dashboard

| | |
|---|---|
| **Route** | `/admin/dashboard` |
| **Layout** | Admin |
| **Mô tả** | System-wide overview. Stats cards, charts, recent activity. |

**Components:**
- `SystemStatsCards` - Total Users, Total Events, Total Bookings, Total Revenue
- `BookingVolumeChart` - Line chart bookings theo thời gian
- `RevenueChart` - Revenue by day/week/month
- `PaymentStatusBreakdown` - Pie chart payment statuses
- `RecentBookingsTable` - 10 bookings gần nhất
- `RecentPaymentsTable` - 10 payments gần nhất

**API:**
- `GET /api/bookings/admin/list?limit=10`
- `GET /api/payments/admin/list?limit=10`
- `GET /api/users/admin/list`

---

### E2. User Management

| | |
|---|---|
| **Route** | `/admin/users` |
| **Layout** | Admin |
| **Mô tả** | Danh sách tất cả users. Search, filter by role, actions. |

**Components:**
- `UserTable` - Email, name, role badge, status (active/inactive), registration date, actions
- `SearchBar` - Search by name/email
- `RoleFilter` - All, Individual, Organization, Admin
- `CreateUserDialog`
- `UserActionsMenu` - Per-row: View, Edit, Deactivate, Reset Password

**API:**
- `GET /api/users/admin/list`
- `POST /api/users/admin` - Create user
- `GET /api/auth/users/:userId`
- `PUT /api/auth/users/:userId`
- `DELETE /api/auth/users/:userId`

---

### E3. User Detail (Admin)

| | |
|---|---|
| **Route** | `/admin/users/:userId` |
| **Layout** | Admin |
| **Mô tả** | Chi tiết 1 user: profile, bookings, payments. Admin actions. |

**Components:**
- `UserProfileCard` - Name, email, role, registration date, status
- `UserBookingsTable` - Bookings của user này
- `UserPaymentsTable` - Payments của user này
- `UserAddressList` - Addresses
- `AdminActionsPanel` - Change role, activate/deactivate, reset password

**API:**
- `GET /api/users/admin/:userId`
- `GET /api/auth/users/:userId`
- `GET /api/bookings/admin/list?user_id=X`
- `GET /api/payments/admin/list?user_id=X`
- `PUT /api/auth/users/:userId`

---

### E4. All Bookings

| | |
|---|---|
| **Route** | `/admin/bookings` |
| **Layout** | Admin |
| **Mô tả** | Tất cả bookings trong hệ thống. Filter, search, export. |

**Components:**
- `BookingTable` - Reference, user email, event name, status badge, amount, date, actions
- `StatusFilter` - Pending, Confirmed, Cancelled, Failed, Expired
- `DateRangeFilter`
- `SearchBar` - By reference or user email
- `BookingDetailDrawer` - Side drawer xem chi tiết
- `ExportButton` - Export CSV

**API:**
- `GET /api/bookings/admin/list?status=X&event_id=X&page=X&limit=20`

---

### E5. All Payments

| | |
|---|---|
| **Route** | `/admin/payments` |
| **Layout** | Admin |
| **Mô tả** | Tất cả payments. Filter, export. |

**Components:**
- `PaymentTable` - ID, user, booking ref, amount, method, status badge, date
- `StatusFilter` - Pending, Processing, Success, Failed, Cancelled, Refunded
- `DateRangeFilter`
- `PaymentDetailDrawer`
- `ExportButton`

**API:**
- `GET /api/payments/admin/list?status=X&page=X&limit=20`

---

### E6. Refund Management

| | |
|---|---|
| **Route** | `/admin/payments/refunds` |
| **Layout** | Admin |
| **Mô tả** | Quản lý refund requests. Approve/reject refunds. |

**Components:**
- `RefundTable` - Refund ID, payment ID, amount, type (full/partial), status badge, reason, date
- `RefundStatusFilter` - Pending, Processing, Success, Failed, Cancelled
- `UpdateRefundStatusDialog` - Approve/reject với reason
- `RefundDetailDrawer` - Chi tiết refund + payment gốc

**API:**
- `GET /api/payments/admin/list` (filter refunded)
- `GET /api/payments/:paymentId/refunds`
- `PUT /api/payments/refunds/:refundId`

---

## Key Flows

### Booking Flow (User)

```
Landing (/) ──> Event Listing (/events) ──> Event Detail (/events/:id)
                                                    │
                                              [Book Now]
                                                    │
                                                    ▼
                                    ┌─── Seat Selection ───┐
                                    │  /booking/:id/seats   │
                                    │  Select seats on map  │
                                    │  Reserve seats (10min)│
                                    └───────────┬───────────┘
                                                │
                                                ▼
                                    ┌──── Checkout ─────────┐
                                    │  /booking/:id/checkout │
                                    │  Review + discount     │
                                    │  Create booking        │
                                    └───────────┬───────────┘
                                                │
                                                ▼
                                    ┌──── Payment ──────────┐
                                    │  /booking/:id/payment  │
                                    │  Stripe payment form   │
                                    │  Process payment       │
                                    └───────────┬───────────┘
                                                │
                                         [payment:success]
                                                │
                                                ▼
                                    ┌── Confirmation ───────┐
                                    │  /booking/:id/confirm  │
                                    │  Tickets + QR codes    │
                                    │  Download / Share      │
                                    └───────────────────────┘
```

### Event Creation Flow (Organization)

```
Org Dashboard (/org/dashboard) ──> [Create Event]
                                        │
                                        ▼
                        ┌── Create Event Wizard ──────────────────┐
                        │  /org/events/create                     │
                        │                                         │
                        │  Step 1: Basic Info (name, dates, tags) │
                        │           │                             │
                        │  Step 2: Venue (name, address, capacity)│
                        │           │                             │
                        │  Step 3: Zones (create zones + colors)  │
                        │           │                             │
                        │  Step 4: Seats (bulk create, map editor)│
                        │           │                             │
                        │  Step 5: Pricing (per zone/category)    │
                        │           │                             │
                        │  Step 6: Ticket Types (Standard, VIP..) │
                        │           │                             │
                        │  [Save Draft] or [Publish]              │
                        └─────────────────────────────────────────┘
                                        │
                                   [Published]
                                        │
                                        ▼
                        ┌── Manage Event ─────────────────────────┐
                        │  /org/events/:eventId/...               │
                        │                                         │
                        │  /edit          - Edit basic info       │
                        │  /zones         - Manage zones          │
                        │  /seats         - Seat map editor       │
                        │  /pricing       - Manage pricing        │
                        │  /availability  - View/block seats      │
                        │  /ticket-types  - Manage ticket types   │
                        └─────────────────────────────────────────┘
```

### WebSocket Connection Lifecycle

```
1. User login ──> Store JWT token
2. Connect WebSocket: ws://localhost:3003?token=JWT
3. Receive: system:connected { connection_id }
4. Auto-subscribe to user channel: user:{userId}
5. On event page:
   Send: room:join { event_id }
   ──> Receive realtime seat updates
6. During booking:
   ──> booking:queue_position (queue updates)
   ──> booking:processing
   ──> booking:confirmed / booking:failed
   ──> payment:processing / payment:success / payment:failed
7. On page leave:
   Send: room:leave { event_id }
8. On logout: disconnect WebSocket
```

---

## Route Summary

| Route | Page | Layout | Role |
|-------|------|--------|------|
| `/` | Landing | Public | All |
| `/events` | Event Listing | Public | All |
| `/events/:eventId` | Event Detail | Public | All |
| `/login` | Login | Auth | Guest |
| `/register` | Register | Auth | Guest |
| `/forgot-password` | Forgot Password | Auth | Guest |
| `/reset-password` | Reset Password | Auth | Guest |
| `/verify-email` | Verify Email | Auth | Guest |
| `/booking/:eventId/seats` | Seat Selection | Dashboard | individual |
| `/booking/:eventId/checkout` | Checkout | Dashboard | individual |
| `/booking/:eventId/payment` | Payment | Dashboard | individual |
| `/booking/:eventId/confirmation` | Confirmation | Dashboard | individual |
| `/my-bookings` | My Bookings | Dashboard | individual |
| `/my-bookings/:bookingId` | Booking Detail | Dashboard | individual |
| `/my-tickets` | My Tickets | Dashboard | individual |
| `/my-tickets/:ticketId` | Ticket Detail | Dashboard | individual |
| `/payments` | Payment History | Dashboard | individual |
| `/payments/:paymentId` | Payment Detail | Dashboard | individual |
| `/profile` | Profile | Dashboard | All auth |
| `/profile/addresses` | Addresses | Dashboard | All auth |
| `/org/dashboard` | Org Dashboard | Organization | organization |
| `/org/events` | Org Events | Organization | organization |
| `/org/events/create` | Create Event | Organization | organization |
| `/org/events/:eventId/edit` | Edit Event | Organization | organization |
| `/org/events/:eventId/zones` | Zone Management | Organization | organization |
| `/org/events/:eventId/seats` | Seat Map Editor | Organization | organization |
| `/org/events/:eventId/pricing` | Pricing Management | Organization | organization |
| `/org/events/:eventId/availability` | Availability | Organization | organization |
| `/org/events/:eventId/ticket-types` | Ticket Types | Organization | organization |
| `/admin/dashboard` | Admin Dashboard | Admin | admin |
| `/admin/users` | User Management | Admin | admin |
| `/admin/users/:userId` | User Detail | Admin | admin |
| `/admin/bookings` | All Bookings | Admin | admin |
| `/admin/payments` | All Payments | Admin | admin |
| `/admin/payments/refunds` | Refund Management | Admin | admin |

**Total: 34 pages across 4 layouts**
