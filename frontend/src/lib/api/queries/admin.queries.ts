import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type { AuthUser } from "../types/auth";
import type { Payment, PaymentFilters } from "../types/payment";
import type { PaginatedResponse, PaginationParams, ApiError } from "../types/common";

// ── Admin User Queries ──

export function useAdminUsers(filters?: PaginationParams) {
  return useQuery<PaginatedResponse<AuthUser>, ApiError>({
    queryKey: queryKeys.users.adminList(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.users.adminList, {
        params: filters,
      });
      return {
        items: data.users ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        limit: data.limit ?? 10,
      };
    },
  });
}

export function useAdminUser(userId: string) {
  return useQuery<AuthUser, ApiError>({
    queryKey: queryKeys.users.adminDetail(userId),
    queryFn: async () => {
      const { data } = await apiClient.get(
        API_ENDPOINTS.users.adminDetail(userId)
      );
      return data.user;
    },
    enabled: !!userId,
  });
}

export function useAdminUpdateUser(userId: string) {
  const queryClient = useQueryClient();
  return useMutation<AuthUser, ApiError, Partial<AuthUser>>({
    mutationFn: async (input) => {
      const { data } = await apiClient.put(
        API_ENDPOINTS.auth.user(userId),
        input
      );
      return data.user;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.users.adminDetail(userId),
      });
      queryClient.invalidateQueries({
        queryKey: queryKeys.users.adminList(),
      });
    },
  });
}

export function useAdminDeleteUser() {
  const queryClient = useQueryClient();
  return useMutation<void, ApiError, string>({
    mutationFn: async (userId) => {
      await apiClient.delete(API_ENDPOINTS.auth.user(userId));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.users.adminList(),
      });
    },
  });
}

// ── Admin Payment Queries ──

export function useAdminPayments(filters?: PaymentFilters & PaginationParams) {
  return useQuery<PaginatedResponse<Payment>, ApiError>({
    queryKey: queryKeys.payments.adminList(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.payments.adminList, {
        params: filters,
      });
      return {
        items: data.payments ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        limit: data.limit ?? 10,
      };
    },
  });
}
