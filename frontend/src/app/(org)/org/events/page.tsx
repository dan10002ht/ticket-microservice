"use client";

import Link from "next/link";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/molecules/page-header";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { useEvents } from "@/lib/api/queries";
import type { Event } from "@/lib/api/types/event";

const columns: Column<Event>[] = [
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

export default function OrgEventsPage() {
  const { data, isLoading } = useEvents();
  const events = data?.items ?? [];

  return (
    <>
      <PageHeader
        title="My Events"
        description="Create and manage your events"
      >
        <Button asChild>
          <Link href="/org/events/create">
            <Plus className="mr-2 h-4 w-4" />
            Create Event
          </Link>
        </Button>
      </PageHeader>

      <div className="mt-6">
        <DataTable
          columns={columns}
          data={events}
          isLoading={isLoading}
          emptyMessage="No events yet. Create your first event!"
          keyExtractor={(e) => e.id}
        />
      </div>
    </>
  );
}
