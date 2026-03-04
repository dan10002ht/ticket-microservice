import { z } from "zod";

export const ticketSelectionSchema = z.object({
  ticket_type_id: z.string().min(1, "Please select a ticket type"),
  quantity: z.coerce
    .number()
    .int("Quantity must be a whole number")
    .min(1, "Minimum 1 ticket required"),
});

export const specialRequestsSchema = z.object({
  special_requests: z
    .string()
    .max(500, "Maximum 500 characters")
    .optional()
    .default(""),
});

export type TicketSelectionInput = z.infer<typeof ticketSelectionSchema>;
export type SpecialRequestsInput = z.infer<typeof specialRequestsSchema>;
