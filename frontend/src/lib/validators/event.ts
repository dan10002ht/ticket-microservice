import { z } from "zod";

export const createEventSchema = z.object({
  name: z.string().min(2, "Event name must be at least 2 characters"),
  description: z.string().optional(),
  start_date: z.string().min(1, "Please select a start date"),
  end_date: z.string().min(1, "Please select an end date"),
  venue_name: z.string().min(2, "Venue name must be at least 2 characters"),
  venue_address: z.string().optional(),
  venue_city: z.string().optional(),
  venue_country: z.string().optional(),
  venue_capacity: z.string().optional(),
  event_type: z.string().optional(),
  category: z.string().optional(),
  sale_start_date: z.string().optional(),
  sale_end_date: z.string().optional(),
  min_age: z.string().optional(),
  is_featured: z.boolean().optional(),
});

export type CreateEventInput = z.infer<typeof createEventSchema>;

export const ticketTypeSchema = z.object({
  name: z.string().min(1, "Ticket type name is required"),
  description: z.string().optional(),
  price: z.number().min(0, "Price must be 0 or more"),
  currency: z.string().optional(),
  quantity: z.number().int().min(1, "At least 1 ticket required"),
  max_per_purchase: z.number().int().min(1).optional(),
  min_per_purchase: z.number().int().min(1).optional(),
});

export type TicketTypeInput = z.infer<typeof ticketTypeSchema>;
