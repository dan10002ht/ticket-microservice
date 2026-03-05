import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type {
  UserProfile,
  UserProfileUpdateRequest,
  UserAddress,
  UserAddressCreateRequest,
} from "../types/user";
import type { ApiError } from "../types/common";

// ── Queries ──

export function useProfile() {
  return useQuery<UserProfile, ApiError>({
    queryKey: queryKeys.users.profile(),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.users.profile);
      // Gateway returns { profile: { ... } } with camelCase keys
      const p = data?.profile ?? data;
      return {
        id: p.userId ?? p.id ?? "",
        email: p.email ?? "",
        firstName: p.firstName ?? "",
        lastName: p.lastName ?? "",
        phone: p.phone ?? "",
        dateOfBirth: p.dateOfBirth ?? "",
        avatarUrl: p.avatarUrl ?? "",
        preferences: p.preferences ?? {},
        createdAt: p.createdAt ?? "",
        updatedAt: p.updatedAt ?? "",
      };
    },
  });
}

export function useAddresses() {
  return useQuery<UserAddress[], ApiError>({
    queryKey: queryKeys.users.addresses(),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.users.addresses);
      return data;
    },
  });
}

// ── Mutations ──

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  return useMutation<UserProfile, ApiError, UserProfileUpdateRequest>({
    mutationFn: async (input) => {
      // Send camelCase — gateway auto-transforms to snake_case for gRPC
      const { data } = await apiClient.put(API_ENDPOINTS.users.profile, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.profile() });
    },
  });
}

export function useCreateAddress() {
  const queryClient = useQueryClient();
  return useMutation<UserAddress, ApiError, UserAddressCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(API_ENDPOINTS.users.addresses, input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.addresses() });
    },
  });
}

export function useUpdateAddress(id: string) {
  const queryClient = useQueryClient();
  return useMutation<UserAddress, ApiError, UserAddressCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.put(
        API_ENDPOINTS.users.address(id),
        input
      );
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.addresses() });
    },
  });
}

export function useDeleteAddress() {
  const queryClient = useQueryClient();
  return useMutation<void, ApiError, string>({
    mutationFn: async (id) => {
      await apiClient.delete(API_ENDPOINTS.users.address(id));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.addresses() });
    },
  });
}
