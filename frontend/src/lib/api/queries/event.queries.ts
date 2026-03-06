import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type {
  Event,
  EventCreateRequest,
  EventUpdateRequest,
  EventFilters,
  EventSeatingZone,
  EventSeatFull,
  ZoneCreateRequest,
  BulkSeatCreateRequest,
  Pricing,
  PricingCreateRequest,
  Availability,
} from "../types/event";
import type { PaginatedResponse, PaginationParams, ApiError } from "../types/common";

// ── Queries ──

export function useEvents(filters?: EventFilters & PaginationParams) {
  return useQuery<PaginatedResponse<Event>, ApiError>({
    queryKey: queryKeys.events.list(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.events.list, {
        params: filters,
      });
      // Normalize: gateway returns { events: [], pagination: {} }
      return {
        items: data.events ?? [],
        total: data.pagination?.total ?? 0,
        page: data.pagination?.page ?? 1,
        limit: data.pagination?.limit ?? 10,
        totalPages: data.pagination?.totalPages,
      };
    },
  });
}

export function useEvent(id: string) {
  return useQuery<Event, ApiError>({
    queryKey: queryKeys.events.detail(id),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.events.detail(id));
      return data;
    },
    enabled: !!id,
  });
}

export function useEventZones(eventId: string) {
  return useQuery<EventSeatingZone[], ApiError>({
    queryKey: queryKeys.events.zones(eventId),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.events.zones(eventId));
      return data;
    },
    enabled: !!eventId,
  });
}

export function useEventPricing(eventId: string) {
  return useQuery<Pricing[], ApiError>({
    queryKey: queryKeys.events.pricing(eventId),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.events.pricing(eventId));
      return data;
    },
    enabled: !!eventId,
  });
}

export function useEventAvailability(eventId: string) {
  return useQuery<Availability, ApiError>({
    queryKey: queryKeys.events.availability(eventId),
    queryFn: async () => {
      const { data } = await apiClient.get(
        API_ENDPOINTS.events.availability(eventId)
      );
      return data;
    },
    enabled: !!eventId,
  });
}

// ── Mutations ──

export function useCreateEvent() {
  const queryClient = useQueryClient();
  return useMutation<Event, ApiError, EventCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(API_ENDPOINTS.events.list, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.events.lists() });
    },
  });
}

export function useUpdateEvent(id: string) {
  const queryClient = useQueryClient();
  return useMutation<Event, ApiError, EventUpdateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.put(
        API_ENDPOINTS.events.detail(id),
        input
      );
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.events.detail(id) });
      queryClient.invalidateQueries({ queryKey: queryKeys.events.lists() });
    },
  });
}

export function useDeleteEvent() {
  const queryClient = useQueryClient();
  return useMutation<void, ApiError, string>({
    mutationFn: async (id) => {
      await apiClient.delete(API_ENDPOINTS.events.detail(id));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.events.all });
    },
  });
}

export function usePublishEvent() {
  const queryClient = useQueryClient();
  return useMutation<Event, ApiError, string>({
    mutationFn: async (id) => {
      const { data } = await apiClient.post(API_ENDPOINTS.events.publish(id));
      return data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.events.detail(id) });
      queryClient.invalidateQueries({ queryKey: queryKeys.events.lists() });
    },
  });
}

export function useCreateZone(eventId: string) {
  const queryClient = useQueryClient();
  return useMutation<EventSeatingZone, ApiError, ZoneCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.events.zones(eventId),
        input
      );
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.events.zones(eventId),
      });
    },
  });
}

export function useCreatePricing(eventId: string) {
  const queryClient = useQueryClient();
  return useMutation<Pricing, ApiError, PricingCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.events.pricing(eventId),
        input
      );
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.events.pricing(eventId),
      });
    },
  });
}

export function useEventSeats(eventId: string, filters?: { zone_id?: string }) {
  return useQuery<EventSeatFull[], ApiError>({
    queryKey: queryKeys.events.seats(eventId, filters as Record<string, string>),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.events.seats(eventId), {
        params: { ...filters, limit: 500 },
      });
      return data.seats ?? [];
    },
    enabled: !!eventId,
  });
}

export function useBulkCreateSeats(eventId: string) {
  const queryClient = useQueryClient();
  return useMutation<{ createdCount: number }, ApiError, BulkSeatCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.events.seatsBulk(eventId),
        input
      );
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.events.seats(eventId),
      });
    },
  });
}
