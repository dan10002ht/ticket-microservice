"use client";

import { PageHeader } from "@/components/molecules/page-header";
import { StatusBadge } from "@/components/molecules/status-badge";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { useAdminPayments } from "@/lib/api/queries";
import type { Payment } from "@/lib/api/types/payment";

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
  const { data, isLoading } = useAdminPayments();
  const payments = data?.items ?? [];

  return (
    <>
      <PageHeader
        title="Payments"
        description="Overview of all platform payments"
      />

      <div className="mt-6">
        <DataTable
          columns={columns}
          data={payments}
          isLoading={isLoading}
          emptyMessage="No payments found"
          keyExtractor={(p) => p.id}
        />
      </div>
    </>
  );
}
