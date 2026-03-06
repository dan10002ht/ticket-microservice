export { queryKeys } from "./query-keys";

// Auth
export {
  useMe,
  useLogin,
  useRegister,
  useLogout,
  useForgotPassword,
  useResetPassword,
  useSendVerificationEmail,
  useVerifyUser,
} from "./auth.queries";

// Events
export {
  useEvents,
  useEvent,
  useEventZones,
  useEventPricing,
  useEventAvailability,
  useCreateEvent,
  useUpdateEvent,
  useDeleteEvent,
  usePublishEvent,
  useCreateZone,
  useCreatePricing,
  useEventSeats,
  useBulkCreateSeats,
} from "./event.queries";

// Bookings
export {
  useBookings,
  useBooking,
  useAdminBookings,
  useCreateBooking,
  useUpdateBooking,
  useCancelBooking,
  useConfirmBooking,
  useReserveSeats,
  useReleaseSeats,
} from "./booking.queries";

// Tickets
export {
  useTickets,
  useTicketTypes,
  useTicketAvailability,
  useCreateTicketType,
  useUpdateTicketType,
  useDeleteTicketType,
  useReserveTickets,
  useReleaseTickets,
} from "./ticket.queries";

// Payments
export {
  usePayments,
  usePayment,
  usePaymentMethods,
  usePaymentRefunds,
  useCreatePayment,
  useCapturePayment,
  useCancelPayment,
  useRefundPayment,
} from "./payment.queries";

// Users
export {
  useProfile,
  useAddresses,
  useUpdateProfile,
  useCreateAddress,
  useUpdateAddress,
  useDeleteAddress,
} from "./user.queries";

// Admin
export {
  useAdminUsers,
  useAdminUser,
  useAdminUpdateUser,
  useAdminDeleteUser,
  useAdminPayments,
} from "./admin.queries";

// Invoices
export {
  useInvoices,
  useInvoice,
  useDownloadInvoicePdf,
} from "./invoice.queries";

// Check-ins
export {
  useCheckins,
  useCheckInStats,
  useCheckInDetail,
  useProcessCheckIn,
} from "./checkin.queries";
