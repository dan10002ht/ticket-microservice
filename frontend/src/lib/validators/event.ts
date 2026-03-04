import { z } from "zod";

export const createEventSchema = z.object({
  name: z.string().min(2, "Event name must be at least 2 characters"),
  description: z.string().optional(),
  start_date: z.string().min(1, "Please select a start date"),
  end_date: z.string().optional(),
  venue_name: z.string().min(2, "Venue name must be at least 2 characters"),
  venue_address: z.string().optional(),
  venue_city: z.string().optional(),
  venue_capacity: z.string().optional(),
});

export type CreateEventInput = z.infer<typeof createEventSchema>;
