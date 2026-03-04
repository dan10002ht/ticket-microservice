"use client";

import {
  Activity,
  CalendarDays,
  DollarSign,
  Users,
} from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { PageHeader } from "@/components/molecules/page-header";
import { StatsCard } from "@/components/molecules/stats-card";
import { StatusBadge } from "@/components/molecules/status-badge";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { useAdminUsers, useAdminBookings, useEvents } from "@/lib/api/queries";
import type { Booking } from "@/lib/api/types/booking";

const bookingColumns: Column<Booking>[] = [
  { key: "id", header: "Booking ID", render: (b) => b.id.slice(0, 8) + "…" },
  {
    key: "status",
    header: "Status",
    render: (b) => <StatusBadge status={b.status} />,
  },
  {
    key: "ticket_quantity",
    header: "Tickets",
    className: "text-center",
  },
  {
    key: "total_amount",
    header: "Amount",
    render: (b) => `${(b.total_amount ?? 0).toLocaleString("vi-VN")} VND`,
  },
  {
    key: "created_at",
    header: "Date",
    render: (b) =>
      new Date(b.created_at).toLocaleDateString("vi-VN", {
        day: "2-digit",
        month: "short",
      }),
  },
];

function StatsSkeleton() {
  return (
    <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {Array.from({ length: 4 }).map((_, i) => (
        <Skeleton key={i} className="h-[106px] rounded-lg" />
      ))}
    </div>
  );
}

export default function AdminDashboardPage() {
  const { data: usersData, isLoading: usersLoading } = useAdminUsers({ limit: 1 });
  const { data: eventsData, isLoading: eventsLoading } = useEvents({ limit: 1 });
  const { data: bookingsData, isLoading: bookingsLoading } = useAdminBookings({ limit: 5 });

  const isLoading = usersLoading || eventsLoading;
  const totalUsers = usersData?.total ?? 0;
  const totalEvents = eventsData?.total ?? 0;
  const recentBookings = bookingsData?.items ?? [];

  return (
    <>
      <PageHeader
        title="Admin Dashboard"
        description="System overview and statistics"
      />

      {isLoading ? (
        <StatsSkeleton />
      ) : (
        <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard
            title="Total Users"
            value={totalUsers.toLocaleString("vi-VN")}
            icon={Users}
          />
          <StatsCard
            title="Total Events"
            value={totalEvents}
            icon={CalendarDays}
          />
          <StatsCard
            title="Revenue (MTD)"
            value="—"
            icon={DollarSign}
            description="Payment stats coming soon"
          />
          <StatsCard
            title="System Health"
            value="99.9%"
            icon={Activity}
            description="All services operational"
          />
        </div>
      )}

      <div className="mt-8 grid gap-6 lg:grid-cols-2">
        <div className="rounded-lg border bg-card p-6">
          <h3 className="font-semibold">Revenue Chart</h3>
          <div className="mt-4 flex h-48 items-center justify-center text-sm text-muted-foreground">
            Chart placeholder — will integrate with recharts
          </div>
        </div>
        <div>
          <h3 className="mb-4 font-semibold">Recent Bookings</h3>
          <DataTable
            columns={bookingColumns}
            data={recentBookings}
            isLoading={bookingsLoading}
            emptyMessage="No bookings yet"
            keyExtractor={(b) => b.id}
          />
        </div>
      </div>
    </>
  );
}
