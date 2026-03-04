// ── Domain types ──

export type TicketStatus =
  | "available"
  | "reserved"
  | "sold"
  | "used"
  | "cancelled";

export interface Ticket {
  id: string;
  event_id: string;
  ticket_type_id: string;
  user_id?: string;
  booking_id?: string;
  ticket_number: string;
  seat_number?: string;
  price: number;
  currency?: string;
  status: TicketStatus;
  qr_code?: string;
  valid_from?: string;
  valid_until?: string;
  created_at: string;
  updated_at: string;
}

export type TicketTypeStatus = "active" | "inactive" | "sold_out";

export interface TicketType {
  id: string;
  event_id: string;
  name: string;
  description?: string;
  price: number;
  currency?: string;
  quantity: number;
  available_quantity: number;
  max_per_purchase: number;
  min_per_purchase: number;
  valid_from?: string;
  valid_until?: string;
  status: TicketTypeStatus;
}

// ── Request types ──

export interface TicketTypeCreateRequest {
  event_id: string;
  name: string;
  description?: string;
  price: number;
  currency?: string;
  quantity: number;
  max_per_purchase?: number;
  min_per_purchase?: number;
}

export interface TicketReserveRequest {
  ticket_type_id: string;
  quantity: number;
  seat_numbers?: string[];
  timeout_seconds?: number;
}

export interface TicketReleaseRequest {
  reservation_id: string;
  ticket_ids?: string[];
}

export interface TicketFilters {
  event_id?: string;
  ticket_type_id?: string;
  status?: TicketStatus;
}

// ── Response types ──

export interface ReservationResponse {
  success: boolean;
  tickets: Ticket[];
  reservation_id: string;
  expires_at: string;
  message?: string;
}

export interface AvailabilityResponse {
  success: boolean;
  available: boolean;
  available_quantity: number;
  available_seats?: string[];
  message?: string;
}
