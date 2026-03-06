---
paths:
  - "frontend/src/**/*.tsx"
  - "frontend/src/**/*.ts"
---

# Frontend Rules

## Component Hierarchy (Atomic Design)
- `components/ui/` — shadcn/ui atoms. NEVER modify directly.
- `components/molecules/` — 2-3 atoms combined. Pure presentation. No data fetching.
- `components/organisms/<domain>/` — Complex sections. May have state. No direct API calls in molecules.
- `app/(group)/route/page.tsx` — Pages compose organisms + wire data hooks.

## Imports
- Use `@/` alias for all imports (maps to `frontend/src/`)
- Import hooks from `@/lib/api/queries` (barrel export)
- Import types from `@/lib/api/types/<domain>`
- Import UI atoms from `@/components/ui/<component>`
- Use `cn()` from `@/lib/utils` for conditional classNames

## Patterns
- Named exports for components, `export default` only for `page.tsx`
- Props interface: `{ComponentName}Props`
- Forms: `react-hook-form` + `zod` + `zodResolver`
- Toasts: `import { toast } from "sonner"` — `toast.success()`, `toast.error()`
- Icons: `lucide-react` only

## After Changes
- Run `npx tsc --noEmit` to verify no type errors
