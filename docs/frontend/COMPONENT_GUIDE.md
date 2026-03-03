# Frontend Component & Pattern Guide

## 1. Atomic Design (Build Order)

Luôn build từ nhỏ đến lớn. Không bao giờ skip level.

```
Atoms → Molecules → Organisms → Templates → Pages
```

| Level | Folder | Ví dụ | Data fetching? |
|-------|--------|-------|----------------|
| **Atom** | `components/ui/` | Button, Input, Badge, Card, Dialog | No |
| **Molecule** | `components/molecules/` | SearchInput, StatusBadge, PriceDisplay, CountdownTimer | No |
| **Organism** | `components/organisms/` | EventCard, Navbar, SeatMap, PaymentForm, DataTable | No (nhận qua props) |
| **Template** | `components/templates/` | PublicLayout, DashboardLayout, OrgLayout, AdminLayout | No |
| **Page** | `app/(group)/page.tsx` | EventListingPage, BookingPage | Yes (qua hooks) |

**Atoms** = shadcn/ui primitives. Install, không sửa.
**Molecules** = ghép 2-3 atoms, pure presentation.
**Organisms** = UI sections phức tạp, có thể có internal state, nhưng KHÔNG fetch data.
**Templates** = shell layout, nhận children.
**Pages** = wire data vào organisms.

---

## 2. Container/Presentational Pattern

Tách logic khỏi UI.

```
Page (Container) ──── fetch data, handle state
  └── Organism (Presentational) ──── nhận props, render UI
        └── Molecule (Presentational) ──── pure, stateless
```

**Rule**: Organisms và Molecules KHÔNG gọi API trực tiếp. Data luôn được truyền từ Page qua props.

```tsx
// Page (Container) - fetches data
export default function EventListingPage() {
  const { data, isLoading } = useEvents(params);
  return <EventGrid events={data} isLoading={isLoading} />;
}

// Organism (Presentational) - receives data via props
export function EventGrid({ events, isLoading }: EventGridProps) {
  if (isLoading) return <EventGridSkeleton />;
  return events.map((e) => <EventCard key={e.id} event={e} />);
}
```

---

## 3. Composition over Inheritance

Dùng composition (children, render props, slots) thay vì inheritance.

```tsx
// Good: Composition via children
<Card>
  <CardHeader><CardTitle>{title}</CardTitle></CardHeader>
  <CardContent>{children}</CardContent>
</Card>

// Good: Composition via slots
<DataTable
  columns={columns}
  data={data}
  filterBar={<EventFilterSidebar />}
  emptyState={<EmptyState message="No events found" />}
/>
```

---

## 4. Compound Component Pattern

Cho các component phức tạp có nhiều phần liên kết (stepper, tabs, form wizard).

```tsx
// Usage
<BookingStepper currentStep={2}>
  <BookingStepper.Step label="Select Seats" />
  <BookingStepper.Step label="Checkout" />
  <BookingStepper.Step label="Payment" />
  <BookingStepper.Step label="Confirmation" />
</BookingStepper>

// Implementation
function BookingStepper({ currentStep, children }) { ... }
BookingStepper.Step = function Step({ label }) { ... };
```

---

## 5. Custom Hook Pattern (Separation of Concerns)

Mỗi hook 1 concern. Không trộn lẫn.

```
lib/hooks/
├── use-events.ts         # TanStack Query: events CRUD
├── use-bookings.ts       # TanStack Query: bookings CRUD
├── use-auth.ts           # Auth logic: login, logout, refresh
├── use-websocket.ts      # WebSocket connection + message handling
├── use-countdown.ts      # Generic countdown timer logic
├── use-debounce.ts       # Generic debounce
├── use-media-query.ts    # Responsive breakpoint detection
└── use-local-storage.ts  # localStorage wrapper
```

**Rule**: Hooks trong `lib/hooks/` là domain-specific (events, bookings, etc.). Generic hooks (debounce, media query) cũng ở đây.

---

## 6. Error Boundary Pattern

3 levels of error handling:

```
Level 1: Field-level    → FormMessage (zod validation)
Level 2: Component-level → toast notification (API errors in mutations)
Level 3: Page-level     → error.tsx (uncaught errors, crashed components)
```

```tsx
// error.tsx - Page-level error boundary
"use client";
export default function Error({ error, reset }: { error: Error; reset: () => void }) {
  return (
    <div className="flex flex-col items-center gap-4 py-20">
      <h2>Something went wrong</h2>
      <Button onClick={reset}>Try again</Button>
    </div>
  );
}
```

---

## 7. Loading State Pattern

3 types, always implement:

| Type | Usage | Component |
|------|-------|-----------|
| **Skeleton** | Initial page load, data fetch | `loading.tsx` + Skeleton atoms |
| **Spinner** | Button click, form submit | Inline spinner in Button |
| **Optimistic** | Mutation (update, delete) | TanStack Query optimistic update |

```tsx
// loading.tsx - Next.js automatic loading state
export default function Loading() {
  return <EventGridSkeleton />;  // Skeleton matching the page layout
}
```

---

## 8. Route Guard Pattern

Protect routes by role. Implement at layout level.

```tsx
// components/templates/DashboardLayout.tsx
export function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, user } = useAuthStore();

  if (!isAuthenticated) redirect("/login");
  if (user?.role !== "individual") redirect("/");

  return (
    <>
      <AuthNavbar />
      <main>{children}</main>
    </>
  );
}
```

---

## 9. Optimistic Update Pattern

Cho UX mượt khi mutate data.

```tsx
export function useToggleSeat() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: toggleSeatApi,
    onMutate: async (seatId) => {
      await queryClient.cancelQueries({ queryKey: ["seats"] });
      const prev = queryClient.getQueryData(["seats"]);
      queryClient.setQueryData(["seats"], (old) =>
        old.map((s) => s.id === seatId ? { ...s, selected: !s.selected } : s)
      );
      return { prev };
    },
    onError: (_err, _vars, context) => {
      queryClient.setQueryData(["seats"], context.prev); // rollback
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["seats"] });
    },
  });
}
```

---

## 10. Barrel Export Pattern

Mỗi folder có `index.ts` export tất cả.

```typescript
// components/molecules/index.ts
export { SearchInput } from "./search-input";
export { StatusBadge } from "./status-badge";
export { PriceDisplay } from "./price-display";
export { CountdownTimer } from "./countdown-timer";
export { EmptyState } from "./empty-state";
```

```tsx
// Clean imports in pages
import { SearchInput, StatusBadge, PriceDisplay } from "@/components/molecules";
```

---

## Component Inventory

### Atoms (shadcn/ui) — cần install
```
Button, Input, Label, Textarea, Select, Checkbox, Switch, RadioGroup,
Badge, Avatar, Skeleton, Separator, Tooltip, Card, Dialog, Sheet,
Tabs, Table, DropdownMenu, Command, Popover, Calendar, Toast
```

### Molecules — cần build
```
SearchInput, FormField, StatusBadge, PriceDisplay, CountdownTimer,
AvatarWithName, EmptyState, ConfirmDialog, DateRangeSelector,
PaginationControls, SortDropdown, RoleBadge, PasswordInput,
FileUploader, StepIndicator
```

### Organisms — cần build (by domain)
```
layout/    Navbar, AuthNavbar, Footer, OrgSidebar, AdminSidebar, Topbar, MobileMenu
auth/      LoginForm, RegisterForm, OAuthButtons, PinCodeInput, ForgotPasswordForm
events/    EventCard, EventGrid, EventFilterSidebar, EventHero, EventInfoTabs, VenueMapPreview
booking/   BookingStepIndicator, PricingBreakdown, SelectedSeatsPanel, BookingSummaryCard
seat-map/  InteractiveSeatMap, ZoneLegend, SeatStatusLegend, ZoomControls
payment/   StripePaymentForm, PaymentMethodSelector, BillingAddressSelector, OrderSummary
tickets/   TicketCard, TicketDisplay, QRCodeFull
org/       EventWizardStepper, ZoneBuilder, SeatMapEditor, BulkSeatCreator, PricingConfigurator, TicketTypeCreator
admin/     UserTable, BookingTable, PaymentTable, RefundTable, SystemStatsCards
shared/    DataTable, FilterBar, StatsCard, RevenueChart, SalesChart
```
