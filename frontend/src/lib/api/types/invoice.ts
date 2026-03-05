// ── Domain types ──

export type InvoiceStatus = "pending" | "generated" | "cancelled";

export interface InvoiceItem {
  id: string;
  description: string;
  quantity: number;
  unit_price: number;
  total_price: number;
}

export interface Invoice {
  id: string;
  invoice_number: string;
  booking_id: string;
  payment_id?: string;
  user_id: string;
  event_id?: string;
  items: InvoiceItem[];
  subtotal: number;
  tax_amount: number;
  total_amount: number;
  currency: string;
  status: InvoiceStatus;
  issued_at?: string;
  created_at: string;
  updated_at: string;
}

// ── Filter types ──

export interface InvoiceFilters {
  status?: InvoiceStatus;
  booking_id?: string;
}
