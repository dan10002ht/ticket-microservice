"use client";

import { useState } from "react";
import { PageHeader } from "@/components/molecules/page-header";
import { StatusBadge } from "@/components/molecules/status-badge";
import { Pagination } from "@/components/molecules/pagination";
import { FilterBar, type FilterSlot } from "@/components/molecules/filter-bar";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { useAdminPayments } from "@/lib/api/queries";
import { getTotalPages } from "@/lib/utils";
import type { Payment, PaymentStatus } from "@/lib/api/types/payment";

const LIMIT = 20;

const statusOptions = [
  { value: "pending", label: "Pending" },
  { value: "authorized", label: "Authorized" },
  { value: "captured", label: "Captured" },
  { value: "completed", label: "Completed" },
  { value: "failed", label: "Failed" },
  { value: "refunded", label: "Refunded" },
  { value: "cancelled", label: "Cancelled" },
];

const columns: Column<Payment>[] = [
  {
    key: "id",
    header: "Payment ID",
    render: (p) => (
      <span className="font-mono text-xs">{p.id.slice(0, 8)}…</span>
    ),
  },
  {
    key: "amount",
    header: "Amount",
    render: (p) =>
      `${p.amount.toLocaleString("vi-VN")} ${(p.currency ?? "VND").toUpperCase()}`,
  },
  {
    key: "payment_method",
    header: "Method",
    render: (p) => (
      <span className="capitalize">
        {p.payment_method.replace("_", " ")}
      </span>
    ),
  },
  {
    key: "status",
    header: "Status",
    render: (p) => <StatusBadge status={p.status} />,
  },
  {
    key: "created_at",
    header: "Date",
    render: (p) =>
      new Date(p.created_at).toLocaleDateString("vi-VN", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      }),
  },
];

export default function AdminPaymentsPage() {
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("");

  const { data, isLoading } = useAdminPayments({
    status: (statusFilter || undefined) as PaymentStatus | undefined,
    page,
    limit: LIMIT,
  });
  const totalPages = getTotalPages(data?.total ?? 0, LIMIT);

  const filterSlots: FilterSlot[] = [
    {
      key: "status",
      placeholder: "Filter by status",
      options: statusOptions,
      value: statusFilter,
      onChange: (v) => {
        setStatusFilter(v);
        setPage(1);
      },
    },
  ];

  return (
    <>
      <PageHeader
        title="Payments"
        description="Overview of all platform payments"
      />

      <div className="mt-6">
        <FilterBar slots={filterSlots} />
      </div>

      <div className="mt-6">
        <DataTable
          columns={columns}
          data={data?.items ?? []}
          isLoading={isLoading}
          emptyMessage="No payments found"
          keyExtractor={(p) => p.id}
        />
      </div>

      <div className="mt-6">
        <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
      </div>
    </>
  );
}
