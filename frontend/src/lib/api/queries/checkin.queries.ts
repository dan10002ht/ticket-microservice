import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type {
  CheckIn,
  CheckInRequest,
  EventCheckInStats,
} from "../types/checkin";
import type { PaginatedResponse, PaginationParams, ApiError } from "../types/common";

// ── Queries ──

export function useCheckins(
  eventId: string,
  params?: PaginationParams & { gate?: string }
) {
  return useQuery<PaginatedResponse<CheckIn>, ApiError>({
    queryKey: queryKeys.checkins.event(eventId),
    queryFn: async () => {
      const { data } = await apiClient.get(
        API_ENDPOINTS.checkins.event(eventId),
        { params }
      );
      return {
        items: data.checkins ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        limit: data.limit ?? 10,
      };
    },
    enabled: !!eventId,
  });
}

export function useCheckInStats(eventId: string) {
  return useQuery<EventCheckInStats, ApiError>({
    queryKey: queryKeys.checkins.stats(eventId),
    queryFn: async () => {
      const { data } = await apiClient.get(
        API_ENDPOINTS.checkins.stats(eventId)
      );
      return data;
    },
    enabled: !!eventId,
  });
}

export function useCheckInDetail(id: string) {
  return useQuery<CheckIn, ApiError>({
    queryKey: queryKeys.checkins.detail(id),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.checkins.detail(id));
      return data;
    },
    enabled: !!id,
  });
}

// ── Mutations ──

export function useProcessCheckIn(eventId: string) {
  const queryClient = useQueryClient();
  return useMutation<CheckIn, ApiError, CheckInRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.checkins.checkin(eventId),
        input
      );
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.checkins.event(eventId),
      });
      queryClient.invalidateQueries({
        queryKey: queryKeys.checkins.stats(eventId),
      });
    },
  });
}
