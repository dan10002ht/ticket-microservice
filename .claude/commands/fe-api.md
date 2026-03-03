# Create API Client + Query Hook

You are a principal frontend developer. Generate an API client function and TanStack Query hook for connecting to the backend.

## Instructions

1. Check the backend API endpoint in `gateway/src/routes/` and `gateway/src/swagger/` to understand the exact request/response shape
2. Create/update the API client function in `frontend/src/lib/api/`
3. Create/update the TanStack Query hook in `frontend/src/lib/hooks/`
4. Add/update TypeScript types in `frontend/src/types/`

## Backend Gateway

- Base URL: `http://localhost:53000` (configured via env `NEXT_PUBLIC_API_URL`)
- Auth routes: `/api/auth/*` (no auth header needed)
- All other routes: `/api/*` (require `Authorization: Bearer <token>` header)
- Swagger docs: `http://localhost:53000/api/docs`

## API Client Pattern

```typescript
// lib/api/client.ts - Base axios instance
import axios from "axios";

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:53000",
  headers: { "Content-Type": "application/json" },
});

// Interceptor adds auth token automatically
apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().tokens?.accessToken;
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});
```

```typescript
// lib/api/events.ts - Per-domain API functions
import { apiClient } from "./client";
import type { Event, EventListParams } from "@/types/event";

export const eventsApi = {
  list: (params?: EventListParams) =>
    apiClient.get<Event[]>("/api/events", { params }).then((r) => r.data),

  getById: (eventId: string) =>
    apiClient.get<Event>(`/api/events/${eventId}`).then((r) => r.data),

  create: (data: CreateEventInput) =>
    apiClient.post<Event>("/api/events", data).then((r) => r.data),
};
```

## TanStack Query Hook Pattern

```typescript
// lib/hooks/use-events.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { eventsApi } from "@/lib/api/events";

// Query keys - consistent structure
export const eventKeys = {
  all: ["events"] as const,
  lists: () => [...eventKeys.all, "list"] as const,
  list: (params?: EventListParams) => [...eventKeys.lists(), params] as const,
  details: () => [...eventKeys.all, "detail"] as const,
  detail: (id: string) => [...eventKeys.details(), id] as const,
};

// Query hook
export function useEvents(params?: EventListParams) {
  return useQuery({
    queryKey: eventKeys.list(params),
    queryFn: () => eventsApi.list(params),
  });
}

export function useEvent(eventId: string) {
  return useQuery({
    queryKey: eventKeys.detail(eventId),
    queryFn: () => eventsApi.getById(eventId),
    enabled: !!eventId,
  });
}

// Mutation hook
export function useCreateEvent() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: eventsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: eventKeys.lists() });
    },
  });
}
```

## Query Key Structure

```
["auth", "permissions"]
["auth", "roles"]
["users", "profile"]
["users", "addresses"]
["events", "list", { status, page, limit, category }]
["events", "detail", eventId]
["events", "detail", eventId, "zones"]
["events", "detail", eventId, "seats", { zone_id, status }]
["events", "detail", eventId, "pricing"]
["events", "detail", eventId, "availability"]
["tickets", "types", eventId]
["tickets", "availability", eventId]
["tickets", "list", { status }]
["tickets", "detail", ticketId]
["bookings", "list", { status, page, limit }]
["bookings", "detail", bookingId]
["bookings", "admin", { status, event_id, page, limit }]
["payments", "list", { status, page, limit }]
["payments", "detail", paymentId]
["payments", "detail", paymentId, "refunds"]
["payments", "admin", { page, limit }]
```

## API Domains

| Domain | File | Base path |
|--------|------|-----------|
| Auth | `lib/api/auth.ts` | `/api/auth` |
| Users | `lib/api/users.ts` | `/api/users` |
| Events | `lib/api/events.ts` | `/api/events` |
| Tickets | `lib/api/tickets.ts` | `/api/tickets` |
| Bookings | `lib/api/bookings.ts` | `/api/bookings` |
| Payments | `lib/api/payments.ts` | `/api/payments` |
| Organizations | `lib/api/organizations.ts` | `/api/organizations` |

## Rules

- Always check the actual backend route/swagger before generating types
- Use `enabled` option in queries that depend on a param (e.g., `enabled: !!eventId`)
- Invalidate related queries after mutations
- Handle pagination with `keepPreviousData: true`
- Error types should match backend error response format

## User Input

$ARGUMENTS
