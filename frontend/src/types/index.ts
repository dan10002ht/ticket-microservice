export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: "individual" | "organization" | "admin";
  phone?: string;
  createdAt: string;
}

export interface Event {
  id: string;
  title: string;
  description: string;
  venue: string;
  address: string;
  startDate: string;
  endDate: string;
  imageUrl?: string;
  status: "draft" | "published" | "cancelled" | "completed";
  organizerId: string;
  totalCapacity: number;
  availableCapacity: number;
  minPrice: number;
  maxPrice: number;
  category?: string;
  createdAt: string;
}

export interface Booking {
  id: string;
  userId: string;
  eventId: string;
  status: "pending" | "confirmed" | "cancelled" | "completed";
  totalAmount: number;
  ticketCount: number;
  createdAt: string;
  event?: Event;
}

export interface Ticket {
  id: string;
  bookingId: string;
  eventId: string;
  seatId?: string;
  zoneName: string;
  seatLabel?: string;
  status: "active" | "used" | "cancelled" | "expired";
  qrCode?: string;
  price: number;
}

export interface Payment {
  id: string;
  bookingId: string;
  amount: number;
  currency: string;
  status: "pending" | "completed" | "failed" | "refunded";
  method: string;
  transactionId?: string;
  createdAt: string;
}

export interface Zone {
  id: string;
  eventId: string;
  name: string;
  capacity: number;
  price: number;
  color: string;
}

export interface Seat {
  id: string;
  zoneId: string;
  label: string;
  row: string;
  number: number;
  status: "available" | "reserved" | "sold";
  price: number;
}
