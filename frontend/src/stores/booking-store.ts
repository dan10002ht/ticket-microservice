import { create } from "zustand";
import type { Booking } from "@/lib/api/types/booking";
import type { ReservationResponse } from "@/lib/api/types/ticket";
import type { Payment } from "@/lib/api/types/payment";

export type BookingStep = "select" | "review" | "payment" | "confirmation";

interface BookingFlowState {
  // Current step
  step: BookingStep;

  // Ticket selection
  selectedTicketTypeId: string | null;
  quantity: number;
  specialRequests: string;

  // Server data
  reservation: ReservationResponse | null;
  booking: Booking | null;
  payment: Payment | null;

  // Actions
  setTicketSelection: (ticketTypeId: string, quantity: number) => void;
  setSpecialRequests: (text: string) => void;
  setReservation: (reservation: ReservationResponse) => void;
  setBooking: (booking: Booking) => void;
  setPayment: (payment: Payment) => void;
  goToStep: (step: BookingStep) => void;
  reset: () => void;
}

const initialState = {
  step: "select" as BookingStep,
  selectedTicketTypeId: null,
  quantity: 1,
  specialRequests: "",
  reservation: null,
  booking: null,
  payment: null,
};

export const useBookingStore = create<BookingFlowState>((set) => ({
  ...initialState,

  setTicketSelection: (ticketTypeId, quantity) =>
    set({ selectedTicketTypeId: ticketTypeId, quantity }),

  setSpecialRequests: (text) => set({ specialRequests: text }),

  setReservation: (reservation) => set({ reservation }),

  setBooking: (booking) => set({ booking }),

  setPayment: (payment) => set({ payment }),

  goToStep: (step) => set({ step }),

  reset: () => set(initialState),
}));
