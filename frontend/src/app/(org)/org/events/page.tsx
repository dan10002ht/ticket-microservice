"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/molecules/page-header";
import { SearchInput } from "@/components/molecules/search-input";
import { Pagination } from "@/components/molecules/pagination";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { useEvents, useDeleteEvent } from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import { getTotalPages } from "@/lib/utils";
import type { Event } from "@/lib/api/types/event";
import type { ApiError } from "@/lib/api/types/common";

const LIMIT = 20;

const columns: Column<Event>[] = [
  {
    key: "name",
    header: "Event",
    render: (e) => (
      <Link
        href={`/org/events/${e.id}`}
        className="font-medium text-primary hover:underline"
      >
        {e.name}
      </Link>
    ),
  },
  {
    key: "start_date",
    header: "Date",
    render: (e) =>
      new Date(e.start_date).toLocaleDateString("en-US", {
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
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");

  const { data, isLoading } = useEvents({ page, limit: LIMIT });
  const deleteMutation = useDeleteEvent();

  const events = data?.items ?? [];
  const totalPages = getTotalPages(data?.total ?? 0, LIMIT);

  const filtered = useMemo(() => {
    if (!search) return events;
    const q = search.toLowerCase();
    return events.filter(
      (e) =>
        e.name.toLowerCase().includes(q) ||
        e.venue_name.toLowerCase().includes(q)
    );
  }, [events, search]);

  const handleDelete = async (id: string) => {
    try {
      await deleteMutation.mutateAsync(id);
      showToast.success("Event deleted.");
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  const columnsWithActions: Column<Event>[] = [
    ...columns,
    {
      key: "actions",
      header: "",
      render: (e) => (
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="sm" asChild>
            <Link href={`/org/events/${e.id}`}>View</Link>
          </Button>
          <Button variant="ghost" size="sm" asChild>
            <Link href={`/org/events/${e.id}/edit`}>Edit</Link>
          </Button>
        </div>
      ),
    },
  ];

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
        <SearchInput
          placeholder="Search by event name, venue..."
          onSearch={setSearch}
          className="max-w-md"
        />
      </div>

      <div className="mt-6">
        <DataTable
          columns={columnsWithActions}
          data={filtered}
          isLoading={isLoading}
          emptyMessage="No events yet. Create your first event!"
          keyExtractor={(e) => e.id}
        />
      </div>

      <div className="mt-6">
        <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
      </div>
    </>
  );
}
