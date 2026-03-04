// ── Domain types ──

export type PaymentMethod =
  | "credit_card"
  | "debit_card"
  | "bank_transfer"
  | "digital_wallet";

export type GatewayProvider = "stripe" | "paypal" | "square";

export type PaymentStatus =
  | "pending"
  | "authorized"
  | "captured"
  | "completed"
  | "failed"
  | "refunded"
  | "cancelled";

export interface Payment {
  id: string;
  booking_id: string;
  ticket_id?: string;
  user_id: string;
  amount: number;
  currency?: string;
  payment_method: PaymentMethod;
  gateway_provider: GatewayProvider;
  status: PaymentStatus;
  transaction_id?: string;
  idempotency_key?: string;
  metadata?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export type RefundType = "full" | "partial";

export type RefundStatus = "pending" | "processing" | "completed" | "failed";

export interface Refund {
  id: string;
  payment_id: string;
  amount: number;
  currency?: string;
  reason?: string;
  description?: string;
  refund_type: RefundType;
  status: RefundStatus;
  external_reference?: string;
  failure_reason?: string;
  metadata?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface PaymentMethodInfo {
  id: string;
  name: string;
  enabled: boolean;
}

// ── Request types ──

export interface PaymentCreateRequest {
  booking_id: string;
  ticket_id?: string;
  amount: number;
  currency?: string;
  payment_method: PaymentMethod;
  gateway_provider?: GatewayProvider;
  idempotency_key?: string;
  metadata?: Record<string, unknown>;
}

export interface PaymentCancelRequest {
  reason?: string;
}

export interface RefundCreateRequest {
  amount?: number;
  reason?: string;
  description?: string;
  refund_type?: RefundType;
  idempotency_key?: string;
  metadata?: Record<string, unknown>;
}

export interface RefundUpdateRequest {
  status: RefundStatus;
  external_reference?: string;
  failure_reason?: string;
}

export interface PaymentFilters {
  status?: PaymentStatus;
  user_id?: string;
}
