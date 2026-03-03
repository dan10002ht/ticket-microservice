# Create Frontend Page

You are a principal frontend developer. Create a new Next.js page following Atomic Design and project conventions.

## Instructions

1. Read `docs/frontend/PAGES.md` to understand the page specs (route, layout, components, API endpoints)
2. Identify which page to create from the docs
3. **Check dependencies first**: Does this page need organisms/molecules that don't exist yet? If yes, create those first (bottom-up)
4. Generate the page file at the correct path under `frontend/src/app/`

## Atomic Design Build Order

A page is Level 5 in Atomic Design. Before creating a page, verify:

```
Page (Level 5) — app/(group)/route/page.tsx
  └── composes Templates (Level 4) — components/templates/
        └── composes Organisms (Level 3) — components/organisms/
              └── composes Molecules (Level 2) — components/molecules/
                    └── composes Atoms (Level 1) — components/ui/
```

**Pre-flight checklist:**
1. Are all needed atoms installed? (shadcn/ui: `npx shadcn-ui@latest add button card ...`)
2. Are all needed molecules created? (`components/molecules/`)
3. Are all needed organisms created? (`components/organisms/`)
4. Is the layout template created? (`components/templates/`)
5. Are API hooks ready? (`lib/hooks/`)
6. Are Zustand stores ready? (`stores/`)

If any dependency is missing, create it first before the page.

## Route Groups & Layouts

| Group | Path prefix | Template | Auth | Role |
|-------|------------|----------|------|------|
| `(public)` | `/` | PublicLayout | No | All |
| `(auth)` | `/login`, `/register` | AuthLayout | No (guest only) | Guest |
| `(dashboard)` | `/booking/*`, `/my-*`, `/profile`, `/payments` | DashboardLayout | Yes | individual |
| `(org)` | `/org/*` | OrgLayout | Yes | organization |
| `(admin)` | `/admin/*` | AdminLayout | Yes | admin |

## Page Pattern

```tsx
// "use client" only if the page needs client-side interactivity
"use client";

// 1. Hooks (API, stores)
// 2. Organisms (compose the page UI)
// 3. Molecules (if directly used)
import { useEvents } from "@/lib/hooks/use-events";
import { EventGrid } from "@/components/organisms/events/event-grid";
import { EventFilterSidebar } from "@/components/organisms/events/event-filter-sidebar";

export default function EventListingPage() {
  // Data hooks
  const { data: events, isLoading } = useEvents(params);

  // Page composes organisms — minimal inline UI
  return (
    <div className="container mx-auto py-6">
      <div className="flex gap-6">
        <EventFilterSidebar onFilterChange={setParams} />
        <EventGrid events={events} isLoading={isLoading} />
      </div>
    </div>
  );
}
```

## Checklist

For each page, generate:
- [ ] Verify all atom/molecule/organism dependencies exist
- [ ] Page component at correct route path (`page.tsx`)
- [ ] Loading state (`loading.tsx`) if data fetching is needed
- [ ] Error boundary (`error.tsx`) for pages with API calls
- [ ] TypeScript types if not already in `types/`
- [ ] API hooks if not already in `lib/hooks/`

## User Input

$ARGUMENTS
