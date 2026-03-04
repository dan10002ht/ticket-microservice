# TicketBox Frontend — Flow & Screen Map

## Route Overview

```
                            ┌──────────────────────────┐
                            │     Root Layout          │
                            │  (Providers, Themes,     │
                            │   AuthInitializer)       │
                            └──────────┬───────────────┘
                                       │
          ┌────────────┬───────────────┼────────────────┬──────────────┐
          ▼            ▼               ▼                ▼              ▼
     (public)       (auth)        (dashboard)        (org)         (admin)
     Navbar +      AuthLayout     AuthNavbar +     OrgLayout      AdminLayout
     Footer        (centered)     DashboardLayout  Sidebar +      Sidebar +
                                                   Topbar         Topbar
```

---

## 1. Public Pages (no auth required)

| Route | Screen | Mô tả |
|-------|--------|-------|
| `/` | **Homepage** | Hero banner, featured events (3 items từ API), categories, features, social proof stats |
| `/events` | **Event Listing** | Danh sách events từ API, search filter by name/venue/city, skeleton loading, empty state |
| `/events/[id]` | **Event Detail** | Hero image, event info (date, venue, capacity), zones & pricing cards, availability stats, "Book Now" CTA |

### Flow: Browse → Detail

```
Homepage ──[Browse Events]──→ Event Listing ──[Click card]──→ Event Detail
                                    ▲                              │
                                    └──────[Back to Events]────────┘
```

---

## 2. Auth Pages (redirect nếu đã login)

| Route | Screen | Mô tả |
|-------|--------|-------|
| `/login` | **Login** | Email/password form, "Forgot password?" link, error banner, redirect theo `callbackUrl` |
| `/register` | **Register** | First/last name, email, password + confirm, password strength meter, error banner |

### Flow: Auth

```
Navbar [Login] ──→ Login Page ──[success]──→ callbackUrl hoặc /
                       │
                  [Create one] ──→ Register Page ──[success]──→ /
                       │
                [Forgot password?] ──→ /forgot-password (chưa implement)
```

---

## 3. Dashboard Pages (yêu cầu auth)

| Route | Screen | Mô tả |
|-------|--------|-------|
| `/my-bookings` | **My Bookings** | DataTable: booking ID, event link, tickets, total (VND), status badge, date. Cancel button cho pending bookings |
| `/profile` | **Profile** | Profile card (avatar, role, member since, status) + Edit form (firstName, lastName, email disabled, phone) |

### Flow: User Dashboard

```
AuthNavbar [My Bookings] ──→ My Bookings ──[View Event]──→ /events/[id]
                                   │
                              [Cancel] ──→ confirm cancel ──→ refresh list

AuthNavbar [Profile] ──→ Profile Page ──[Save changes]──→ API update
```

---

## 4. Organizer Pages (yêu cầu auth + org role)

| Route | Screen | Mô tả |
|-------|--------|-------|
| `/org/dashboard` | **Org Dashboard** | Stats cards (events, tickets sold, revenue, attendees), recent bookings table |
| `/org/events` | **Org Events** | Events management table (name, sold/capacity, status, date, venue), "Create Event" button |

> **Status**: Hardcoded data — chưa wire API

### Flow: Organizer

```
Sidebar [Dashboard] ──→ Org Dashboard (overview stats)
Sidebar [Events] ──→ Org Events ──[Create Event]──→ (chưa implement)
                         │
                    [Edit/Manage] ──→ (chưa implement)
```

---

## 5. Admin Pages (yêu cầu auth + admin role)

| Route | Screen | Mô tả |
|-------|--------|-------|
| `/admin/dashboard` | **Admin Dashboard** | System stats (users, events, revenue, health), chart placeholder, activity log placeholder |
| `/admin/users` | **Admin Users** | User management table (name, email, role, status, joined date), search input |

> **Status**: Hardcoded data — chưa wire API

### Flow: Admin

```
Sidebar [Dashboard] ──→ Admin Dashboard (system overview)
Sidebar [Users] ──→ Admin Users ──[Search]──→ filter users
                        │
                   [Edit/Ban] ──→ (chưa implement)
```

---

## Route Protection (Middleware)

```
Unauthenticated user
  │
  ├── /events, /events/[id], / ──→ ✅ Cho phép
  ├── /login, /register ──→ ✅ Cho phép
  ├── /my-bookings, /profile ──→ ❌ Redirect → /login?callbackUrl=/my-bookings
  ├── /org/* ──→ ❌ Redirect → /login?callbackUrl=/org/dashboard
  └── /admin/* ──→ ❌ Redirect → /login?callbackUrl=/admin/dashboard

Authenticated user
  │
  ├── /login, /register ──→ ❌ Redirect → /
  └── Tất cả routes khác ──→ ✅ Cho phép
```

---

## Auth State Flow

```
App Mount
  │
  ▼
AuthInitializer (trong Providers)
  │
  ├── Có cookie tb_access_token?
  │     ├── Có → useMe() validate token → ✅ setUser(user) vào Zustand
  │     │                                 → ❌ clearUser() (token invalid)
  │     └── Không → clearUser() + setHydrated(true)
  │
  ▼
Navbar/AuthNavbar render
  │
  ├── isHydrated=false → không render auth section (tránh flash)
  ├── isAuthenticated=true → User dropdown (avatar + name + logout)
  └── isAuthenticated=false → Login/Register buttons
```

---

## Token Refresh Flow

```
API Request (401 Unauthorized)
  │
  ▼
Auth Interceptor
  │
  ├── Là request /auth/* ? → Skip (tránh infinite loop)
  ├── Đang refresh? → Đưa vào failedQueue, chờ
  └── Chưa refresh → Set isRefreshing=true
         │
         ▼
    POST /auth/refresh { refresh_token }
         │
         ├── ✅ Success → setTokens() → replay tất cả failedQueue
         └── ❌ Fail → clearTokens() + clearUser() → redirect /login
```

---

## Data Flow Summary

| Page | Query Hook | Mutation Hook | Store |
|------|-----------|---------------|-------|
| `/` | `useEvents({ limit: 3 })` | — | — |
| `/events` | `useEvents()` | — | — |
| `/events/[id]` | `useEvent`, `useEventPricing`, `useEventAvailability`, `useEventZones` | — | — |
| `/login` | — | `useLogin()` | — |
| `/register` | — | `useRegister()` | — |
| `/my-bookings` | `useBookings()` | `useCancelBooking()` | — |
| `/profile` | `useProfile()` | `useUpdateProfile()` | `useAuthStore` (fallback) |
| Navbar/Topbar | — | `useLogout()` | `useAuthStore` |
| AuthInitializer | `useMe()` | — | `useAuthStore` |

---

## Chưa Implement (Next Phases)

### Phase 4: Dashboard & Admin
- [ ] Wire `/org/dashboard` với real stats API
- [ ] Wire `/org/events` với `useEvents()` + `useCreateEvent()`
- [ ] Wire `/admin/users` với admin user API
- [ ] Wire `/admin/dashboard` với system stats API

### Phase 5: Booking Flow
- [ ] Tạo `/events/[id]/book` — select zone/seats → create booking
- [ ] Tạo `/events/[id]/book/payment` — payment processing
- [ ] Tạo `/my-bookings/[id]` — booking detail page

### Phase 6: Advanced
- [ ] `/forgot-password` — send reset email
- [ ] `/reset-password` — reset with token
- [ ] Real-time seat availability (WebSocket)
- [ ] `/org/events/[id]/edit` — event management
- [ ] `/my-tickets` — user's tickets list
