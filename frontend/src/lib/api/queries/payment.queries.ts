import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type {
  Payment,
  PaymentCreateRequest,
  PaymentCancelRequest,
  PaymentFilters,
  PaymentMethodInfo,
  Refund,
  RefundCreateRequest,
} from "../types/payment";
import type { PaginatedResponse, PaginationParams, ApiError } from "../types/common";

// ── Queries ──

export function usePayments(filters?: PaymentFilters & PaginationParams) {
  return useQuery<PaginatedResponse<Payment>, ApiError>({
    queryKey: queryKeys.payments.list(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.payments.list, {
        params: { ...filters, page: filters?.page, size: filters?.limit },
      });
      return {
        items: data.payments ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        limit: data.size ?? filters?.limit ?? 10,
      };
    },
  });
}

export function usePayment(id: string) {
  return useQuery<Payment, ApiError>({
    queryKey: queryKeys.payments.detail(id),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.payments.detail(id));
      return data.payment;
    },
    enabled: !!id,
  });
}

export function usePaymentMethods() {
  return useQuery<PaymentMethodInfo[], ApiError>({
    queryKey: queryKeys.payments.methods(),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.payments.methods);
      return data.methods ?? [];
    },
  });
}

export function usePaymentRefunds(paymentId: string) {
  return useQuery<Refund[], ApiError>({
    queryKey: queryKeys.payments.refunds(paymentId),
    queryFn: async () => {
      const { data } = await apiClient.get(
        API_ENDPOINTS.payments.refunds(paymentId)
      );
      return data.refunds ?? [];
    },
    enabled: !!paymentId,
  });
}

// ── Mutations ──

export function useCreatePayment() {
  const queryClient = useQueryClient();
  return useMutation<Payment, ApiError, PaymentCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(API_ENDPOINTS.payments.create, input);
      return data.payment;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.lists() });
      queryClient.invalidateQueries({ queryKey: queryKeys.bookings.all });
    },
  });
}

export function useCapturePayment() {
  const queryClient = useQueryClient();
  return useMutation<Payment, ApiError, string>({
    mutationFn: async (id) => {
      const { data } = await apiClient.post(API_ENDPOINTS.payments.capture(id));
      return data.payment;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.detail(id) });
    },
  });
}

export function useCancelPayment() {
  const queryClient = useQueryClient();
  return useMutation<Payment, ApiError, { id: string } & PaymentCancelRequest>({
    mutationFn: async ({ id, ...input }) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.payments.cancel(id),
        input
      );
      return data.payment;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.payments.all });
    },
  });
}

export function useRefundPayment(paymentId: string) {
  const queryClient = useQueryClient();
  return useMutation<Refund, ApiError, RefundCreateRequest>({
    mutationFn: async (input) => {
      const { data } = await apiClient.post(
        API_ENDPOINTS.payments.refund(paymentId),
        input
      );
      return data.refund;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.payments.refunds(paymentId),
      });
      queryClient.invalidateQueries({
        queryKey: queryKeys.payments.detail(paymentId),
      });
    },
  });
}
