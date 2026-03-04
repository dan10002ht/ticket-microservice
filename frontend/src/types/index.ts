// Canonical API types — re-exported for convenience
export type { AuthUser, AuthResponse } from "@/lib/api/types/auth";
export type {
  Event,
  EventSeatingZone,
  EventSeat,
  EventStatus,
  Pricing,
  Availability,
} from "@/lib/api/types/event";
export type { Booking, BookingStatus, SeatReservation } from "@/lib/api/types/booking";
export type {
  Ticket,
  TicketType,
  TicketStatus,
  TicketTypeStatus,
} from "@/lib/api/types/ticket";
export type {
  Payment,
  PaymentStatus,
  Refund,
  PaymentMethodInfo,
} from "@/lib/api/types/payment";
export type { UserProfile, UserAddress } from "@/lib/api/types/user";
