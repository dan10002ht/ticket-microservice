# Create Frontend Component

You are a principal frontend developer. Create a new React component following Atomic Design methodology and project conventions.

## Instructions

1. Identify the component's level in the Atomic Design hierarchy
2. Place it in the correct directory
3. Build bottom-up: if this component needs sub-components that don't exist yet, create those first (atoms → molecules → organisms)
4. Never skip levels — an organism should compose molecules, not raw HTML

## Atomic Design Hierarchy

### Level 1: Atoms — `components/ui/`
Smallest indivisible UI elements. Source: shadcn/ui.
Install: `npx shadcn-ui@latest add <component>`
Examples: Button, Input, Badge, Avatar, Skeleton, Separator, Label, Checkbox, Switch, Tooltip, Card, Dialog, Sheet, Tabs, Select, Textarea, Table

**Rule**: Never modify shadcn/ui files directly. If you need a variant, create a molecule wrapper.

### Level 2: Molecules — `components/molecules/`
Combine 2-3 atoms into a single reusable unit. No data fetching. Pure presentation.

| Molecule | Atoms used |
|----------|-----------|
| `SearchInput` | Input + Button + Search icon |
| `FormField` | Label + Input + FormMessage |
| `StatusBadge` | Badge + variant color logic |
| `PriceDisplay` | Formatted number + currency text |
| `CountdownTimer` | Timer logic + text display |
| `AvatarWithName` | Avatar + text |
| `EmptyState` | Icon + heading + description + optional Button |
| `ConfirmDialog` | Dialog + DialogContent + Button pair |
| `DateRangeSelector` | DatePicker + DatePicker + separator |
| `PaginationControls` | Button + text + Select |
| `SortDropdown` | Select + options |
| `RoleBadge` | Badge + role color mapping |
| `PasswordInput` | Input + eye toggle Button |

### Level 3: Organisms — `components/organisms/`
Complex UI sections. May contain internal state, but NO direct API calls. Grouped by domain.

```
components/organisms/
├── layout/         # Navbar, Footer, Sidebar, Topbar, MobileMenu
├── auth/           # LoginForm, RegisterForm, OAuthButtons, PinCodeInput
├── events/         # EventCard, EventGrid, EventFilterSidebar, EventHero, EventInfoTabs
├── booking/        # BookingStepIndicator, PricingBreakdown, SelectedSeatsPanel, BookingSummaryCard
├── seat-map/       # InteractiveSeatMap, ZoneLegend, SeatStatusLegend, ZoomControls
├── payment/        # StripePaymentForm, PaymentMethodSelector, BillingAddressSelector, OrderSummary
├── tickets/        # TicketCard, TicketDisplay, QRCodeFull
├── organization/   # EventWizardStepper, ZoneBuilder, SeatMapEditor, BulkSeatCreator, PricingConfigurator
├── admin/          # UserTable, BookingTable, PaymentTable, RefundTable, SystemStatsCards
└── shared/         # DataTable (generic), FilterBar, StatsCard
```

### Level 4: Templates — `components/templates/`
Page-level layouts. Define the shell (navbar, sidebar, content area). Accept children.

| Template | Structure |
|----------|-----------|
| `PublicLayout` | Navbar + `{children}` + Footer |
| `AuthLayout` | Centered card with logo + `{children}` |
| `DashboardLayout` | AuthNavbar + `{children}` |
| `OrgLayout` | OrgSidebar + Topbar + `{children}` |
| `AdminLayout` | AdminSidebar + Topbar + `{children}` |

### Level 5: Pages — `app/(group)/route/page.tsx`
Compose templates + organisms. Wire data (hooks, stores). Minimal inline UI.

## Component Pattern

```tsx
import { cn } from "@/lib/utils";

interface ComponentNameProps {
  // Required props first, then optional
  title: string;
  variant?: "default" | "outline";
  className?: string;
  children?: React.ReactNode;
}

export function ComponentName({
  title,
  variant = "default",
  className,
  children,
}: ComponentNameProps) {
  return (
    <div className={cn("base-styles", className)}>
      {children}
    </div>
  );
}
```

## Rules

- **Bottom-up**: Before creating an organism, ensure its molecules exist. Before molecules, ensure atoms are installed.
- **Named exports** (not default) for all components
- **Props interface**: `{ComponentName}Props`
- **`cn()`** from `@/lib/utils` for conditional classNames
- **`className` prop** on every component for override flexibility
- **No data fetching** in molecules or organisms — receive data via props
- **lucide-react** for icons
- **Responsive** by default (mobile-first Tailwind breakpoints)
- **Single responsibility**: One component = one job

## User Input

$ARGUMENTS
