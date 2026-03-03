# Create Zustand Store

You are a principal frontend developer. Generate a Zustand store for client-side state management.

## Instructions

1. Create the store in `frontend/src/stores/`
2. Follow the project's Zustand patterns
3. Only use Zustand for client state that doesn't belong in TanStack Query (server state)

## When to Use Zustand vs TanStack Query

| Use Case | Tool |
|----------|------|
| API response data (events, bookings, etc.) | TanStack Query |
| Auth state (tokens, current user) | Zustand |
| Booking flow state (selected seats, current step) | Zustand |
| WebSocket connection state | Zustand |
| UI state (modals, sidebar open/close) | Zustand |
| Form wizard multi-step data | Zustand |

## Store Pattern

```typescript
// stores/auth-store.ts
import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthState {
  // State
  user: User | null;
  tokens: { accessToken: string; refreshToken: string } | null;
  isAuthenticated: boolean;

  // Actions
  setAuth: (user: User, tokens: { accessToken: string; refreshToken: string }) => void;
  logout: () => void;
  updateTokens: (tokens: { accessToken: string; refreshToken: string }) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      // Initial state
      user: null,
      tokens: null,
      isAuthenticated: false,

      // Actions
      setAuth: (user, tokens) =>
        set({ user, tokens, isAuthenticated: true }),

      logout: () =>
        set({ user: null, tokens: null, isAuthenticated: false }),

      updateTokens: (tokens) =>
        set({ tokens }),
    }),
    {
      name: "auth-storage", // localStorage key
      partialize: (state) => ({
        user: state.user,
        tokens: state.tokens,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
```

## Existing Stores (planned)

### auth-store.ts
- `user` - Current user (id, email, role, firstName, lastName)
- `tokens` - { accessToken, refreshToken }
- `isAuthenticated` - boolean
- Actions: `setAuth()`, `logout()`, `updateTokens()`
- Middleware: `persist` (localStorage)

### booking-store.ts
- `currentEvent` - Event being booked
- `selectedSeats` - Array of selected seat objects
- `reservationId` - Active reservation ID
- `reservationExpiresAt` - Countdown reference
- `bookingId` - Created booking ID
- `paymentId` - Created payment ID
- `step` - Current step (seats → checkout → payment → confirmation)
- Actions: `selectSeat()`, `deselectSeat()`, `clearSelection()`, `setReservation()`, `setBooking()`, `setPayment()`, `nextStep()`, `reset()`
- No persist (session-only)

### websocket-store.ts
- `connected` - boolean
- `connectionId` - From system:connected
- `joinedRooms` - Set of joined room IDs
- `lastMessage` - Last received message
- Actions: `setConnected()`, `setDisconnected()`, `joinRoom()`, `leaveRoom()`, `setLastMessage()`
- No persist

## Rules

- Use `persist` middleware only for data that should survive page refresh (auth tokens)
- Use `partialize` to control what gets persisted
- Actions are methods on the store interface, not separate functions
- Use `useAuthStore.getState()` for accessing state outside React (e.g., in API interceptors)
- Keep stores focused - one concern per store
- Derive computed values with selectors: `const isOrg = useAuthStore((s) => s.user?.role === "organization")`

## User Input

$ARGUMENTS
