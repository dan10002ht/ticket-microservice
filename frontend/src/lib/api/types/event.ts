// ── Domain types ──

export interface Event {
  id: string;
  organization_id: string;
  name: string;
  description?: string;
  start_date: string;
  end_date?: string;
  venue_name: string;
  venue_address?: string;
  venue_city?: string;
  venue_country?: string;
  venue_capacity?: number;
  canvas_config?: string;
  zones?: EventSeatingZone[];
  seats?: EventSeat[];
  created_at: string;
  updated_at: string;
}

export interface EventSeatingZone {
  id: string;
  event_id: string;
  name: string;
  zone_type: "seated" | "standing" | "vip";
  coordinates?: string;
  seat_count?: number;
  color?: string;
  created_at: string;
  updated_at: string;
}

export interface EventSeat {
  id: string;
  event_id: string;
  zone_id: string;
  seat_number: string;
  row_number?: string;
  coordinates?: string;
  created_at: string;
  updated_at: string;
}

export interface Pricing {
  id: string;
  event_id: string;
  zone_id: string;
  name: string;
  price: number;
  currency: string;
}

export interface Availability {
  total_seats: number;
  available_seats: number;
  reserved_seats: number;
  sold_seats: number;
  blocked_seats: number;
}

// ── Request types ──

export type EventStatus = "draft" | "published" | "cancelled" | "completed";

export interface EventCreateRequest {
  name: string;
  description?: string;
  start_date: string;
  end_date?: string;
  venue_name: string;
  venue_address?: string;
  venue_city?: string;
  venue_country?: string;
  venue_capacity?: number;
  canvas_config?: string;
}

export interface EventUpdateRequest extends Partial<EventCreateRequest> {
  status?: EventStatus;
}

export interface EventFilters {
  category?: string;
  location?: string;
  date?: string;
  status?: EventStatus;
}

export interface ZoneCreateRequest {
  name: string;
  zone_type: "seated" | "standing" | "vip";
  coordinates?: string;
  seat_count?: number;
  color?: string;
}

export interface SeatCreateRequest {
  zone_id: string;
  seat_number: string;
  row_number?: string;
  coordinates?: string;
}

export interface BulkSeatCreateRequest {
  zone_id: string;
  seats: Omit<SeatCreateRequest, "zone_id">[];
}

export interface PricingCreateRequest {
  zone_id: string;
  name: string;
  price: number;
  currency?: string;
}

export interface PriceCalculateRequest {
  seat_ids?: string[];
  quantity?: number;
  zone_id?: string;
  discount_code?: string;
}

export interface BlockSeatsRequest {
  seat_ids: string[];
  reason?: string;
}
