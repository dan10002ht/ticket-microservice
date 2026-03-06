"use client";

import { useState } from "react";
import Link from "next/link";
import { PageHeader } from "@/components/molecules/page-header";
import { StatusBadge } from "@/components/molecules/status-badge";
import { PriceDisplay } from "@/components/molecules/price-display";
import { Pagination } from "@/components/molecules/pagination";
import { FilterBar, type FilterSlot } from "@/components/molecules/filter-bar";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { Button } from "@/components/ui/button";
import { useBookings, useCancelBooking } from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import { getTotalPages } from "@/lib/utils";
import type { Booking, BookingStatus } from "@/lib/api/types/booking";
import type { ApiError } from "@/lib/api/types/common";

const LIMIT = 10;

const statusOptions = [
  { value: "pending", label: "Pending" },
  { value: "confirmed", label: "Confirmed" },
  { value: "cancelled", label: "Cancelled" },
  { value: "completed", label: "Completed" },
  { value: "expired", label: "Expired" },
];

const columns: Column<Booking>[] = [
  {
    key: "id",
    header: "Booking ID",
    render: (booking) => (
      <Link
        href={`/my-bookings/${booking.id}`}
        className="font-mono text-xs text-primary hover:underline"
      >
        {booking.id.slice(0, 8)}...
      </Link>
    ),
  },
  {
    key: "event_id",
    header: "Event",
    render: (booking) => (
      <Link
        href={`/events/${booking.event_id}`}
        className="text-primary hover:underline"
      >
        View Event
      </Link>
    ),
  },
  {
    key: "ticket_quantity",
    header: "Tickets",
    className: "text-center",
  },
  {
    key: "total_amount",
    header: "Total",
    render: (booking) => (
      <PriceDisplay
        amount={booking.total_amount}
        currency={booking.currency || "VND"}
        size="sm"
      />
    ),
  },
  {
    key: "status",
    header: "Status",
    render: (booking) => <StatusBadge status={booking.status} />,
  },
  {
    key: "created_at",
    header: "Date",
    render: (booking) =>
      new Date(booking.created_at).toLocaleDateString("en-US"),
  },
];

export default function MyBookingsPage() {
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("");

  const { data, isLoading, error } = useBookings({
    status: (statusFilter || undefined) as BookingStatus | undefined,
    page,
    limit: LIMIT,
  });
  const cancelMutation = useCancelBooking();

  const totalPages = getTotalPages(data?.total ?? 0, LIMIT);

  const handleCancel = async (id: string) => {
    try {
      await cancelMutation.mutateAsync({ id });
      showToast.success("Booking cancelled successfully");
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  const columnsWithActions: Column<Booking>[] = [
    ...columns,
    {
      key: "actions",
      header: "",
      render: (booking) => (
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="sm" asChild>
            <Link href={`/my-bookings/${booking.id}`}>View</Link>
          </Button>
          {booking.status === "pending" && (
            <Button
              variant="ghost"
              size="sm"
              className="text-destructive"
              onClick={() => handleCancel(booking.id)}
              disabled={cancelMutation.isPending}
            >
              Cancel
            </Button>
          )}
        </div>
      ),
    },
  ];

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
        title="My Bookings"
        description="View and manage your event bookings"
      />

      {error && (
        <div className="mt-4 rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {error.error?.message || "Failed to load bookings."}
        </div>
      )}

      <div className="mt-6">
        <FilterBar slots={filterSlots} />
      </div>

      <div className="mt-6">
        <DataTable
          columns={columnsWithActions}
          data={data?.items ?? []}
          isLoading={isLoading}
          emptyMessage="No bookings yet"
          keyExtractor={(b) => b.id}
        />
      </div>

      <div className="mt-6">
        <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
      </div>
    </>
  );
}
