import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type {
  Booking,
  BookingCreateRequest,
  BookingUpdateRequest,
  BookingConfirmRequest,
  BookingCancelRequest,
  BookingFilters,
  SeatReservation,
  SeatReservationRequest,
  SeatReleaseRequest,
} from "../types/booking";
import type { PaginatedResponse, PaginationParams, ApiError } from "../types/common";

// ── Queries ──

export function useBookings(filters?: BookingFilters & PaginationParams) {
  return useQuery<PaginatedResponse<Booking>, ApiError>({
    queryKey: queryKeys.bookings.list(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.bookings.list, {
        params: filters,
      });
      // Normalize: { success, bookings: [], total, page, limit }
      return {
        items: data.bookings ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        limit: data.limit ?? 10,
      };
    },
  });
}

export function useBooking(id: string) {
  return useQuery<Booking, ApiError>({
    queryKey: queryKeys.bookings.detail(id),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.bookings.detail(id));
      return data.booking;
    },
    enabled: !!id,
  });
}

export function useAdminBookings(filters?: BookingFilters & PaginationParams) {
  return useQuery<PaginatedResponse<Booking>, ApiError>({
    queryKey: queryKeys.bookings.adminList(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.bookings.adminList, {
        params: filters,
      });
      return {
        items: data.bookings ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        limit: data.limit ?? 10,
      };
    },
  });
}

// ── Mutations ──

export function useCreateBooking() {
  const queryClient = useQueryClient();
  return useMutation<Booking, ApiError, BookingCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(API_ENDPOINTS.bookings.list, input);
      return data.booking;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.bookings.lists() });
    },
  });
}

export function useUpdateBooking(id: string) {
  const queryClient = useQueryClient();
  return useMutation<Booking, ApiError, BookingUpdateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.put(
        API_ENDPOINTS.bookings.detail(id),
        input
      );
      return data.booking;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.bookings.detail(id) });
      queryClient.invalidateQueries({ queryKey: queryKeys.bookings.lists() });
    },
  });
}

export function useCancelBooking() {
  const queryClient = useQueryClient();
  return useMutation<Booking, ApiError, { id: string } & BookingCancelRequest>({
    mutationFn: async ({ id, ...input }) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.bookings.cancel(id),
        input
      );
      return data.booking;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.bookings.all });
    },
  });
}

export function useConfirmBooking() {
  const queryClient = useQueryClient();
  return useMutation<Booking, ApiError, { id: string } & BookingConfirmRequest>({
    mutationFn: async ({ id, ...input }) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.bookings.confirm(id),
        input
      );
      return data.booking;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.bookings.all });
    },
  });
}

export function useReserveSeats() {
  return useMutation<SeatReservation, ApiError, SeatReservationRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.bookings.seatsReserve,
        input
      );
      return data.reservation;
    },
  });
}

export function useReleaseSeats() {
  return useMutation<void, ApiError, SeatReleaseRequest>({
    mutationFn: async (input) => {
      await apiClient.post(API_ENDPOINTS.bookings.seatsRelease, input);
    },
  });
}
