import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type {
  Ticket,
  TicketType,
  TicketTypeCreateRequest,
  TicketReserveRequest,
  TicketReleaseRequest,
  TicketFilters,
  ReservationResponse,
  AvailabilityResponse,
} from "../types/ticket";
import type { PaginatedResponse, PaginationParams, ApiError } from "../types/common";

// ── Queries ──

export function useTickets(filters?: TicketFilters & PaginationParams) {
  return useQuery<PaginatedResponse<Ticket>, ApiError>({
    queryKey: queryKeys.tickets.list(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.tickets.list, {
        params: filters,
      });
      return {
        items: data.tickets ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        limit: data.limit ?? 10,
      };
    },
  });
}

export function useTicketTypes(eventId: string) {
  return useQuery<TicketType[], ApiError>({
    queryKey: queryKeys.tickets.types(eventId),
    queryFn: async () => {
      const { data } = await apiClient.get(
        API_ENDPOINTS.tickets.types(eventId),
        { params: { include_availability: true } }
      );
      return data.ticket_types ?? [];
    },
    enabled: !!eventId,
  });
}

export function useTicketAvailability(eventId: string) {
  return useQuery<AvailabilityResponse, ApiError>({
    queryKey: queryKeys.tickets.availability(eventId),
    queryFn: async () => {
      const { data } = await apiClient.get(
        API_ENDPOINTS.tickets.availability(eventId)
      );
      return data;
    },
    enabled: !!eventId,
  });
}

// ── Mutations ──

export function useCreateTicketType() {
  const queryClient = useQueryClient();
  return useMutation<TicketType, ApiError, TicketTypeCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.tickets.typesCreate,
        input
      );
      return data;
    },
    onSuccess: (_, input) => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.tickets.types(input.event_id),
      });
    },
  });
}

export function useUpdateTicketType(typeId: string) {
  const queryClient = useQueryClient();
  return useMutation<TicketType, ApiError, Partial<TicketTypeCreateRequest>>({
    mutationFn: async (input) => {
      const { data } = await apiClient.put(
        API_ENDPOINTS.tickets.typeDetail(typeId),
        input
      );
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets.all });
    },
  });
}

export function useReserveTickets(eventId: string) {
  return useMutation<ReservationResponse, ApiError, TicketReserveRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.tickets.reserve(eventId),
        input
      );
      return data;
    },
  });
}

export function useDeleteTicketType() {
  const queryClient = useQueryClient();
  return useMutation<void, ApiError, string>({
    mutationFn: async (typeId) => {
      await apiClient.delete(API_ENDPOINTS.tickets.typeDetail(typeId));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tickets.all });
    },
  });
}

export function useReleaseTickets() {
  return useMutation<void, ApiError, TicketReleaseRequest>({
    mutationFn: async (input) => {
      await apiClient.post(API_ENDPOINTS.tickets.release, input);
    },
  });
}
