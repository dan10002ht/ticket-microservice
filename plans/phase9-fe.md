# Plan: Phase 9 — Pagination & Filter

## Role & Standards

Bạn là **Senior Frontend Engineer** chuyên React/Next.js, chịu trách nhiệm implement theo các tiêu chuẩn sau:

- **Atomic Design Pattern**: Luôn build bottom-up — Atoms (shadcn/ui) → Molecules → Organisms → Pages. Không skip layer.
- **Tech Stack**: Next.js 14+ App Router, TypeScript strict, shadcn/ui + Tailwind CSS, TanStack Query v5, Zustand
- **Code Quality**:
  - TypeScript strict — type mọi props, không dùng `any`
  - Components phải có interface Props rõ ràng, export named (không default export cho components)
  - Dùng `"use client"` directive cho client components
  - Tailwind classes gọn gàng, dùng `cn()` utility khi cần merge classes
- **UI/UX Best Practices**:
  - Loading states (Skeleton) cho mọi data fetch
  - Empty states có icon + message rõ ràng
  - Responsive design: mobile-first, dùng `sm:`, `md:`, `lg:` breakpoints
  - Accessible: proper `aria-label`, keyboard navigation cho interactive elements
  - Consistent spacing: dùng `space-y-*`, `gap-*` pattern
- **State Management**:
  - Server state → TanStack Query (queries + mutations)
  - Local UI state → `useState` trong page component, truyền xuống qua props
  - Shared client state → Zustand stores
- **Pattern Conventions**:
  - Query hooks: `use[Resource]s()` cho list, `use[Resource]()` cho detail
  - Mutations: `use[Action][Resource]()` (e.g. `useCreateEvent`, `useCancelBooking`)
  - Error handling: `showToast.apiError(err as ApiError)` cho mutations
  - File naming: kebab-case cho files, PascalCase cho components
- **Codebase Location**: Tất cả files trong `frontend/src/`
- **Language**: UI text bằng English, code comments bằng English

---

## Context

Phase 8 hoàn thành (Invoice & Check-in). Tất cả list pages hiện tại load **toàn bộ data** không có pagination. Infrastructure đã ready:
- `PaginatedResponse<T>`: `{ items, total, page, limit, totalPages? }`
- `PaginationParams`: `{ page?, limit? }`
- Tất cả query hooks đã accept `PaginationParams` + domain filters
- shadcn `Select` component có sẵn
- Một số pages có client-side search (admin users, admin events, public events)

**Mục tiêu:** Thêm Pagination molecule + FilterBar molecule, wire vào 6 list pages.

---

## Implementation (Bottom-Up)

### Step 1: Utility — `getTotalPages` helper

**Modify** `src/lib/utils.ts`
- Thêm `getTotalPages(total: number, limit: number): number` — `Math.max(1, Math.ceil(total / limit))`
- Dùng ở tất cả pages thay vì rely vào `totalPages?` optional field

### Step 2: Molecule — `Pagination`

**Create** `src/components/molecules/pagination.tsx`
- Props: `{ page, totalPages, onPageChange, className? }`
- Previous/Next buttons (`ChevronLeft`/`ChevronRight` từ lucide-react)
- Hiện tối đa 5 page numbers với ellipsis: `1 ... 4 5 6 ... 20`
- Return `null` khi `totalPages <= 1`
- Dùng shadcn `Button` variant="outline"/"ghost"

### Step 3: Molecule — `FilterBar`

**Create** `src/components/molecules/filter-bar.tsx`
- Generic, domain-agnostic filter row
- Props:
  ```ts
  interface FilterSlot {
    key: string;
    placeholder: string;
    options: { value: string; label: string }[];
    value: string;
    onChange: (value: string) => void;
  }
  interface FilterBarProps { slots: FilterSlot[]; className?: string; }
  ```
- Render mỗi slot = shadcn `Select`
- "Clear filters" button khi bất kỳ slot nào có value

### Step 4: Export molecules

**Modify** `src/components/molecules/index.ts`
- Thêm export `Pagination` + `FilterBar`

### Step 5: Public Events page — pagination + status filter

**Modify** `src/app/(public)/events/page.tsx`
- Add state: `page`, `statusFilter`
- Call `useEvents({ status: statusFilter || undefined, page, limit: 12 })`
- Render `<FilterBar>` (1 slot: status = published/cancelled/completed)
- Render `<Pagination>` dưới event grid
- Reset page=1 khi filter thay đổi

**Modify** `src/components/organisms/events/event-listing-content.tsx`
- Giữ nguyên client-side search (search trên data current page)
- Thêm optional `footer` ReactNode prop cho Pagination render slot
- Nhận `events` từ parent (đã filtered server-side)

### Step 6: My Bookings page — pagination + status filter

**Modify** `src/app/(dashboard)/my-bookings/page.tsx`
- Add state: `page`, `statusFilter`
- Call `useBookings({ status: statusFilter || undefined, page, limit: 10 })`
- `<FilterBar>` (1 slot: status = pending/confirmed/cancelled/completed/expired)
- `<Pagination>` dưới DataTable

### Step 7: Org Events page — pagination + search

**Modify** `src/app/(org)/org/events/page.tsx`
- Add state: `page`, `search`
- Call `useEvents({ page, limit: 20 })`
- Client-side `useMemo` name filter trên current page data
- Thêm `SearchInput` trên table
- `<Pagination>` dưới DataTable

### Step 8: Admin Events page — pagination

**Modify** `src/app/(admin)/admin/events/page.tsx`
- Add `page` state, pass `{ page, limit: 20 }` to `useEvents()`
- `<Pagination>` dưới DataTable
- Giữ existing client-side search

### Step 9: Admin Users page — pagination

**Modify** `src/app/(admin)/admin/users/page.tsx`
- Add `page` state, pass `{ page, limit: 20 }` to `useAdminUsers()`
- `<Pagination>` dưới DataTable
- Giữ existing client-side search

### Step 10: Admin Payments page — pagination + status filter

**Modify** `src/app/(admin)/admin/payments/page.tsx`
- Add state: `page`, `statusFilter`
- Call `useAdminPayments({ status: statusFilter || undefined, page, limit: 20 })`
- `<FilterBar>` (1 slot: status = pending/authorized/captured/completed/failed/refunded/cancelled)
- `<Pagination>` dưới DataTable

---

## Files Summary

| # | File | Action | Step |
|---|------|--------|------|
| 1 | `src/lib/utils.ts` | **Modify** | 1 |
| 2 | `src/components/molecules/pagination.tsx` | **Create** | 2 |
| 3 | `src/components/molecules/filter-bar.tsx` | **Create** | 3 |
| 4 | `src/components/molecules/index.ts` | **Modify** | 4 |
| 5 | `src/app/(public)/events/page.tsx` | **Modify** | 5 |
| 6 | `src/components/organisms/events/event-listing-content.tsx` | **Modify** | 5 |
| 7 | `src/app/(dashboard)/my-bookings/page.tsx` | **Modify** | 6 |
| 8 | `src/app/(org)/org/events/page.tsx` | **Modify** | 7 |
| 9 | `src/app/(admin)/admin/events/page.tsx` | **Modify** | 8 |
| 10 | `src/app/(admin)/admin/users/page.tsx` | **Modify** | 9 |
| 11 | `src/app/(admin)/admin/payments/page.tsx` | **Modify** | 10 |

**2 create, 9 modify = 11 files total**

---

## Key Decisions

1. **Client-side search giữ nguyên** — Backend không có text search param, search chỉ filter trên current page data.
2. **FilterBar generic** — Nhận `FilterSlot[]`, không import domain types. Reusable cho mọi page.
3. **Reset page=1 khi filter thay đổi** — Tránh page 5 với 0 results.
4. **`getTotalPages` utility** — Không rely vào `totalPages?` optional, tự tính từ `total` và `limit`.

---

## Verification

1. `/events` → grid 12 items/page, pagination controls, status filter dropdown
2. `/my-bookings` → table 10 items/page, pagination, status filter
3. `/org/events` → table 20 items/page, pagination, search input
4. `/admin/events` → table 20 items/page, pagination, existing search works
5. `/admin/users` → table 20 items/page, pagination, existing search works
6. `/admin/payments` → table 20 items/page, pagination, status filter
7. Pagination returns `null` khi totalPages <= 1
8. `npx tsc --noEmit` + `npx next build` pass
