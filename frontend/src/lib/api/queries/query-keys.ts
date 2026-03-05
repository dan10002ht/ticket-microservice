import type { EventFilters } from "../types/event";
import type { BookingFilters } from "../types/booking";
import type { TicketFilters } from "../types/ticket";
import type { PaymentFilters } from "../types/payment";
import type { InvoiceFilters } from "../types/invoice";
import type { PaginationParams } from "../types/common";

export const queryKeys = {
  auth: {
    all: ["auth"] as const,
    user: () => [...queryKeys.auth.all, "user"] as const,
    permissions: () => [...queryKeys.auth.all, "permissions"] as const,
    roles: () => [...queryKeys.auth.all, "roles"] as const,
  },

  events: {
    all: ["events"] as const,
    lists: () => [...queryKeys.events.all, "list"] as const,
    list: (filters?: EventFilters & PaginationParams) =>
      [...queryKeys.events.lists(), filters] as const,
    details: () => [...queryKeys.events.all, "detail"] as const,
    detail: (id: string) => [...queryKeys.events.details(), id] as const,
    zones: (eventId: string) =>
      [...queryKeys.events.detail(eventId), "zones"] as const,
    seats: (eventId: string, filters?: Record<string, string>) =>
      [...queryKeys.events.detail(eventId), "seats", filters] as const,
    pricing: (eventId: string) =>
      [...queryKeys.events.detail(eventId), "pricing"] as const,
    availability: (eventId: string) =>
      [...queryKeys.events.detail(eventId), "availability"] as const,
    templates: () => [...queryKeys.events.all, "templates"] as const,
  },

  bookings: {
    all: ["bookings"] as const,
    lists: () => [...queryKeys.bookings.all, "list"] as const,
    list: (filters?: BookingFilters & PaginationParams) =>
      [...queryKeys.bookings.lists(), filters] as const,
    details: () => [...queryKeys.bookings.all, "detail"] as const,
    detail: (id: string) => [...queryKeys.bookings.details(), id] as const,
    adminList: (filters?: PaginationParams) =>
      [...queryKeys.bookings.all, "admin", filters] as const,
  },

  tickets: {
    all: ["tickets"] as const,
    lists: () => [...queryKeys.tickets.all, "list"] as const,
    list: (filters?: TicketFilters & PaginationParams) =>
      [...queryKeys.tickets.lists(), filters] as const,
    detail: (id: string) =>
      [...queryKeys.tickets.all, "detail", id] as const,
    types: (eventId: string) =>
      [...queryKeys.tickets.all, "types", eventId] as const,
    availability: (eventId: string) =>
      [...queryKeys.tickets.all, "availability", eventId] as const,
  },

  payments: {
    all: ["payments"] as const,
    lists: () => [...queryKeys.payments.all, "list"] as const,
    list: (filters?: PaymentFilters & PaginationParams) =>
      [...queryKeys.payments.lists(), filters] as const,
    detail: (id: string) =>
      [...queryKeys.payments.all, "detail", id] as const,
    methods: () => [...queryKeys.payments.all, "methods"] as const,
    refunds: (paymentId: string) =>
      [...queryKeys.payments.detail(paymentId), "refunds"] as const,
    adminList: (filters?: PaginationParams) =>
      [...queryKeys.payments.all, "admin", filters] as const,
  },

  users: {
    all: ["users"] as const,
    profile: () => [...queryKeys.users.all, "profile"] as const,
    addresses: () => [...queryKeys.users.all, "addresses"] as const,
    adminList: (filters?: PaginationParams) =>
      [...queryKeys.users.all, "admin", filters] as const,
    adminDetail: (id: string) =>
      [...queryKeys.users.all, "admin", id] as const,
  },
  invoices: {
    all: ["invoices"] as const,
    lists: () => [...queryKeys.invoices.all, "list"] as const,
    list: (filters?: InvoiceFilters & PaginationParams) =>
      [...queryKeys.invoices.lists(), filters] as const,
    detail: (id: string) =>
      [...queryKeys.invoices.all, "detail", id] as const,
  },

  checkins: {
    all: ["checkins"] as const,
    event: (eventId: string) =>
      [...queryKeys.checkins.all, "event", eventId] as const,
    stats: (eventId: string) =>
      [...queryKeys.checkins.all, "stats", eventId] as const,
    detail: (id: string) =>
      [...queryKeys.checkins.all, "detail", id] as const,
  },
} as const;
