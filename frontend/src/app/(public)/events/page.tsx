"use client";

import { useState } from "react";
import { PageHeader } from "@/components/molecules/page-header";
import { Pagination } from "@/components/molecules/pagination";
import { FilterBar, type FilterSlot } from "@/components/molecules/filter-bar";
import { EventListingContent } from "@/components/organisms/events/event-listing-content";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useEvents } from "@/lib/api/queries";
import { getTotalPages } from "@/lib/utils";
import type { EventStatus } from "@/lib/api/types/event";

const LIMIT = 12;

const categoryOptions = [
  { value: "music", label: "Music" },
  { value: "arts", label: "Arts" },
  { value: "sports", label: "Sports" },
  { value: "business", label: "Business" },
  { value: "education", label: "Education" },
  { value: "entertainment", label: "Entertainment" },
  { value: "other", label: "Other" },
];

const eventTypeOptions = [
  { value: "concert", label: "Concert" },
  { value: "theater", label: "Theater" },
  { value: "sports", label: "Sports" },
  { value: "conference", label: "Conference" },
  { value: "festival", label: "Festival" },
  { value: "other", label: "Other" },
];

const statusOptions = [
  { value: "published", label: "Published" },
  { value: "cancelled", label: "Cancelled" },
  { value: "completed", label: "Completed" },
];

export default function EventsPage() {
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("");
  const [eventTypeFilter, setEventTypeFilter] = useState("");
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");

  const { data, isLoading, error } = useEvents({
    status: (statusFilter || undefined) as EventStatus | undefined,
    category: categoryFilter || undefined,
    event_type: eventTypeFilter || undefined,
    start_date_from: dateFrom ? new Date(dateFrom).toISOString() : undefined,
    start_date_to: dateTo ? new Date(dateTo).toISOString() : undefined,
    page,
    limit: LIMIT,
  });

  const totalPages = getTotalPages(data?.total ?? 0, LIMIT);

  const resetPage = () => setPage(1);

  const filterSlots: FilterSlot[] = [
    {
      key: "category",
      placeholder: "Category",
      options: categoryOptions,
      value: categoryFilter,
      onChange: (v) => {
        setCategoryFilter(v);
        resetPage();
      },
    },
    {
      key: "event_type",
      placeholder: "Event Type",
      options: eventTypeOptions,
      value: eventTypeFilter,
      onChange: (v) => {
        setEventTypeFilter(v);
        resetPage();
      },
    },
    {
      key: "status",
      placeholder: "Status",
      options: statusOptions,
      value: statusFilter,
      onChange: (v) => {
        setStatusFilter(v);
        resetPage();
      },
    },
  ];

  return (
    <section className="container mx-auto px-4 py-8">
      <PageHeader
        title="All Events"
        description="Browse and discover upcoming events"
      />

      {error && (
        <div className="mt-6 rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {error.error?.message || "Failed to load events. Please try again later."}
        </div>
      )}

      <div className="mt-6 space-y-4">
        <FilterBar slots={filterSlots} />

        <div className="flex flex-wrap items-end gap-4">
          <div className="space-y-1">
            <Label className="text-xs text-muted-foreground">From</Label>
            <Input
              type="date"
              className="w-[160px]"
              value={dateFrom}
              onChange={(e) => {
                setDateFrom(e.target.value);
                resetPage();
              }}
            />
          </div>
          <div className="space-y-1">
            <Label className="text-xs text-muted-foreground">To</Label>
            <Input
              type="date"
              className="w-[160px]"
              value={dateTo}
              onChange={(e) => {
                setDateTo(e.target.value);
                resetPage();
              }}
            />
          </div>
        </div>
      </div>

      <EventListingContent
        events={data?.items ?? []}
        isLoading={isLoading}
      />

      <div className="mt-8">
        <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
      </div>
    </section>
  );
}
