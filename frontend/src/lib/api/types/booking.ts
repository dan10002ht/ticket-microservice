// ── Domain types ──

export type BookingStatus =
  | "pending"
  | "confirmed"
  | "cancelled"
  | "completed"
  | "expired";

export interface Booking {
  id: string;
  user_id: string;
  event_id: string;
  ticket_quantity: number;
  seat_numbers?: string[];
  total_amount: number;
  currency?: string;
  status: BookingStatus;
  special_requests?: string;
  payment_reference?: string;
  metadata?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
  expires_at?: string;
}

export interface SeatReservation {
  reservation_id: string;
  event_id: string;
  seat_numbers: string[];
  user_id: string;
  expires_at: string;
}

// ── Request types ──

export interface BookingCreateRequest {
  event_id: string;
  ticket_quantity: number;
  seat_numbers?: string[];
  special_requests?: string;
  idempotency_key?: string;
  metadata?: Record<string, unknown>;
}

export interface BookingUpdateRequest {
  ticket_quantity?: number;
  seat_numbers?: string[];
  special_requests?: string;
  metadata?: Record<string, unknown>;
}

export interface BookingConfirmRequest {
  payment_reference?: string;
}

export interface BookingCancelRequest {
  reason?: string;
}

export interface SeatReservationRequest {
  event_id: string;
  seat_numbers: string[];
  timeout_seconds?: number;
}

export interface SeatReleaseRequest {
  reservation_id: string;
  seat_numbers?: string[];
}

export interface BookingFilters {
  status?: BookingStatus;
  event_id?: string;
}
