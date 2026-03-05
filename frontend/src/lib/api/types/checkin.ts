// ── Domain types ──

export type CheckInStatus = "success" | "invalid" | "already_used" | "cancelled";

export interface CheckIn {
  id: string;
  ticket_id: string;
  event_id: string;
  user_id?: string;
  qr_code?: string;
  status: CheckInStatus;
  staff_id?: string;
  device_id?: string;
  gate?: string;
  notes?: string;
  check_in_time: string;
  created_at: string;
}

export interface EventCheckInStats {
  event_id: string;
  total_checkins: number;
  unique_tickets: number;
  by_gate: Record<string, number>;
  last_checkin_at?: string;
}

// ── Request types ──

export interface CheckInRequest {
  ticket_id: string;
  qr_code?: string;
  event_id: string;
  staff_id?: string;
  device_id?: string;
  gate?: string;
}
