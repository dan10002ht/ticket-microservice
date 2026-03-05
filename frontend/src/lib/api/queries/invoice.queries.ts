import { useQuery, useMutation } from "@tanstack/react-query";
import { apiClient } from "../client";
import { API_ENDPOINTS } from "../endpoints";
import { queryKeys } from "./query-keys";
import type { Invoice, InvoiceFilters } from "../types/invoice";
import type { PaginatedResponse, PaginationParams, ApiError } from "../types/common";

// ── Queries ──

export function useInvoices(filters?: InvoiceFilters & PaginationParams) {
  return useQuery<PaginatedResponse<Invoice>, ApiError>({
    queryKey: queryKeys.invoices.list(filters),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.invoices.list, {
        params: filters,
      });
      return {
        items: data.invoices ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        limit: data.limit ?? 10,
      };
    },
  });
}

export function useInvoice(id: string) {
  return useQuery<Invoice, ApiError>({
    queryKey: queryKeys.invoices.detail(id),
    queryFn: async () => {
      const { data } = await apiClient.get(API_ENDPOINTS.invoices.detail(id));
      return data.invoice ?? data;
    },
    enabled: !!id,
  });
}

// ── Mutations ──

export function useDownloadInvoicePdf() {
  return useMutation<void, ApiError, string>({
    mutationFn: async (id) => {
      const response = await apiClient.get(API_ENDPOINTS.invoices.pdf(id), {
        responseType: "blob",
      });
      // Trigger browser download
      const blob = new Blob([response.data], { type: "application/pdf" });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `invoice-${id}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    },
  });
}
