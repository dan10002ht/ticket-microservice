"use client";

import Link from "next/link";
import {
  CalendarDays,
  CheckCircle,
  MapPin,
  Ticket,
} from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { PageHeader } from "@/components/molecules/page-header";
import { StatsCard } from "@/components/molecules/stats-card";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { useEvents } from "@/lib/api/queries";
import type { Event } from "@/lib/api/types/event";

const recentColumns: Column<Event>[] = [
  {
    key: "name",
    header: "Event",
    render: (e) => (
      <Link
        href={`/events/${e.id}`}
        className="font-medium hover:underline"
      >
        {e.name}
      </Link>
    ),
  },
  {
    key: "start_date",
    header: "Date",
    render: (e) =>
      new Date(e.start_date).toLocaleDateString("vi-VN", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      }),
  },
  { key: "venue_name", header: "Venue" },
  {
    key: "venue_capacity",
    header: "Capacity",
    render: (e) => (
      <span className="tabular-nums">{e.venue_capacity ?? "—"}</span>
    ),
    className: "text-center",
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

export default function OrgDashboardPage() {
  const { data, isLoading } = useEvents();
  const events = data?.items ?? [];

  const totalEvents = events.length;
  const totalCapacity = events.reduce(
    (sum, e) => sum + (e.venue_capacity ?? 0),
    0
  );
  const uniqueVenues = new Set(events.map((e) => e.venue_name)).size;

  return (
    <>
      <PageHeader
        title="Dashboard"
        description="Overview of your events and sales"
      />

      {isLoading ? (
        <StatsSkeleton />
      ) : (
        <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard
            title="Total Events"
            value={totalEvents}
            icon={CalendarDays}
          />
          <StatsCard
            title="Published"
            value={totalEvents}
            icon={CheckCircle}
          />
          <StatsCard
            title="Total Capacity"
            value={totalCapacity.toLocaleString("vi-VN")}
            icon={Ticket}
          />
          <StatsCard
            title="Venues"
            value={uniqueVenues}
            icon={MapPin}
          />
        </div>
      )}

      <div className="mt-8">
        <h2 className="mb-4 text-lg font-semibold">Recent Events</h2>
        <DataTable
          columns={recentColumns}
          data={events.slice(0, 5)}
          isLoading={isLoading}
          emptyMessage="No events yet. Create your first event!"
          keyExtractor={(e) => e.id}
        />
      </div>
    </>
  );
}
