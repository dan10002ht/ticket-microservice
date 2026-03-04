"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { PageHeader } from "@/components/molecules/page-header";
import { SearchInput } from "@/components/molecules/search-input";
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

export default function AdminEventsPage() {
  const { data, isLoading } = useEvents();
  const events = data?.items ?? [];
  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    if (!search) return events;
    const q = search.toLowerCase();
    return events.filter(
      (e) =>
        e.name.toLowerCase().includes(q) ||
        e.venue_name.toLowerCase().includes(q) ||
        (e.venue_city ?? "").toLowerCase().includes(q)
    );
  }, [events, search]);

  return (
    <>
      <PageHeader
        title="Events"
        description="Overview of all platform events"
      />

      <div className="mt-6">
        <SearchInput
          placeholder="Search by event name, venue..."
          onSearch={setSearch}
          className="max-w-md"
        />
      </div>

      <div className="mt-6">
        <DataTable
          columns={columns}
          data={filtered}
          isLoading={isLoading}
          emptyMessage="No events found"
          keyExtractor={(e) => e.id}
        />
      </div>
    </>
  );
}
