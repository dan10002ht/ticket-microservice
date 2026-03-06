# Create API Client + Query Hook

You are a principal frontend developer. Generate an API client function and TanStack Query hook for connecting to the backend.

## Instructions

1. Check the backend API endpoint in `gateway/src/routes/` to understand the request/response shape
2. Add endpoint URL to `frontend/src/lib/api/endpoints.ts`
3. Add TypeScript types in `frontend/src/lib/api/types/<domain>.ts`
4. Add query keys in `frontend/src/lib/api/queries/query-keys.ts`
5. Create TanStack Query hook in `frontend/src/lib/api/queries/<domain>.queries.ts`
6. Export from `frontend/src/lib/api/queries/index.ts`

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

**Note:** This project does NOT use a separate `eventsApi` object. API calls are inline in query hooks using `apiClient` + `API_ENDPOINTS`.

## Actual Project Pattern

```typescript
// lib/api/queries/event.queries.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type { Event } from "../types/event";
import type { PaginatedResponse, ApiError } from "../types/common";

export function useEvents(filters?: EventFilters & PaginationParams) {
  return useQuery<PaginatedResponse<Event>, ApiError>({
    queryKey: queryKeys.events.list(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.events.list, {
        params: filters,
      });
      return {
        items: data.events ?? [],
        total: data.pagination?.total ?? 0,
        page: data.pagination?.page ?? 1,
        limit: data.pagination?.limit ?? 10,
        totalPages: data.pagination?.totalPages,
      };
    },
  });
}

export function useCreateEvent() {
  const queryClient = useQueryClient();
  return useMutation<Event, ApiError, EventCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(API_ENDPOINTS.events.list, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.events.lists() });
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

## Key Files

| File                                  | Purpose                              |
| ------------------------------------- | ------------------------------------ |
| `lib/api/client.ts`                   | Axios instance with auth interceptor |
| `lib/api/endpoints.ts`                | All API endpoint URLs (centralized)  |
| `lib/api/types/<domain>.ts`           | TypeScript interfaces per domain     |
| `lib/api/queries/query-keys.ts`       | Centralized query key factory        |
| `lib/api/queries/<domain>.queries.ts` | TanStack Query hooks per domain      |
| `lib/api/queries/index.ts`            | Re-exports all hooks                 |

## API Domains

| Domain   | Types file         | Queries file                 | Base path   |
| -------- | ------------------ | ---------------------------- | ----------- |
| Auth     | `types/auth.ts`    | `queries/auth.queries.ts`    | `/auth`     |
| Events   | `types/event.ts`   | `queries/event.queries.ts`   | `/events`   |
| Tickets  | `types/ticket.ts`  | `queries/ticket.queries.ts`  | `/tickets`  |
| Bookings | `types/booking.ts` | `queries/booking.queries.ts` | `/bookings` |
| Payments | `types/payment.ts` | `queries/payment.queries.ts` | `/payments` |
| Users    | `types/user.ts`    | `queries/user.queries.ts`    | `/users`    |
| Invoices | `types/invoice.ts` | `queries/invoice.queries.ts` | `/invoices` |
| Checkins | `types/checkin.ts` | `queries/checkin.queries.ts` | `/checkins` |

## Rules

- Always check the actual backend route/swagger before generating types
- Use `enabled` option in queries that depend on a param (e.g., `enabled: !!eventId`)
- Invalidate related queries after mutations
- Handle pagination with `keepPreviousData: true`
- Error types should match backend error response format

## User Input

$ARGUMENTS
