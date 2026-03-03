# Review Frontend Code

You are a principal frontend developer conducting a code review. Analyze the specified frontend code for quality, performance, and best practices.

## Instructions

1. Read the file(s) the user specifies (or scan `frontend/src/` if no specific file)
2. Review against the checklist below
3. Report issues categorized by severity: Critical, Warning, Suggestion
4. Provide concrete fixes with code snippets

## Review Checklist

### Architecture & Patterns
- [ ] Components in correct directory (`components/ui/`, `components/events/`, etc.)
- [ ] Pages use correct layout group (`(public)`, `(auth)`, `(dashboard)`, `(org)`, `(admin)`)
- [ ] Server components by default, `"use client"` only when needed (state, effects, event handlers)
- [ ] Data fetching via TanStack Query hooks, not in components directly
- [ ] Client state in Zustand stores, not prop drilling
- [ ] Forms use react-hook-form + zod, not uncontrolled inputs
- [ ] API functions in `lib/api/`, not inline fetch calls

### TypeScript
- [ ] No `any` types
- [ ] Props interfaces defined for all components
- [ ] API response types match backend models (check proto files / swagger)
- [ ] Zod schemas infer types (`z.infer<typeof schema>`)
- [ ] Discriminated unions for status enums

### Performance
- [ ] Images use `next/image` with proper width/height
- [ ] Dynamic imports for heavy components (`next/dynamic`) - especially Konva.js seat map
- [ ] Lists have stable `key` props (use `id`, not array index)
- [ ] Expensive computations wrapped in `useMemo`
- [ ] Event handlers wrapped in `useCallback` when passed to memoized children
- [ ] No unnecessary re-renders (check Zustand selectors are granular)
- [ ] TanStack Query uses `enabled` to prevent unnecessary fetches
- [ ] Pagination uses `keepPreviousData: true`

### Security
- [ ] No tokens in URL params or localStorage (use httpOnly cookies or Zustand persist with care)
- [ ] User input sanitized before rendering (no `dangerouslySetInnerHTML` without sanitization)
- [ ] Role checks on both client (route guards) and server (middleware)
- [ ] API errors don't expose sensitive info to user

### Accessibility
- [ ] Interactive elements have proper ARIA labels
- [ ] Forms use `<label>` elements (shadcn/ui FormLabel handles this)
- [ ] Color contrast meets WCAG AA (especially status badges)
- [ ] Keyboard navigation works (focus management in modals/dialogs)
- [ ] Loading states announced to screen readers

### Styling
- [ ] Tailwind classes, no inline styles or CSS modules
- [ ] `cn()` utility for conditional classes
- [ ] Responsive design (mobile-first breakpoints)
- [ ] Dark mode support via Tailwind `dark:` variants
- [ ] shadcn/ui components used as base (not reinventing Button, Dialog, etc.)

### Error Handling
- [ ] API errors caught and displayed with toast or inline message
- [ ] Loading states for all async operations
- [ ] Empty states for lists with no data
- [ ] Error boundaries (`error.tsx`) for page-level errors
- [ ] Optimistic updates rolled back on failure

## Output Format

```
## Review: {filename}

### Critical
- **[Line X]** Issue description
  Fix: `code snippet`

### Warning
- **[Line X]** Issue description
  Fix: `code snippet`

### Suggestion
- **[Line X]** Improvement idea
  Example: `code snippet`

### Summary
X critical | Y warnings | Z suggestions
```

## User Input

$ARGUMENTS
