"use client";

import { useState, useMemo } from "react";
import { SearchInput } from "@/components/molecules/search-input";
import { EventCard } from "@/components/organisms/events/event-card";
import { EmptyState } from "@/components/molecules/empty-state";
import { Skeleton } from "@/components/ui/skeleton";
import { CalendarX } from "lucide-react";
import type { Event } from "@/lib/api/types/event";

interface EventListingContentProps {
  events: Event[];
  isLoading?: boolean;
  emptyMessage?: string;
}

function EventCardSkeleton() {
  return (
    <div className="rounded-lg border overflow-hidden">
      <Skeleton className="aspect-[16/9] w-full" />
      <div className="p-4 space-y-2">
        <Skeleton className="h-5 w-3/4" />
        <Skeleton className="h-4 w-1/2" />
        <Skeleton className="h-4 w-2/3" />
      </div>
      <div className="border-t px-4 py-3 flex justify-between">
        <Skeleton className="h-4 w-20" />
        <Skeleton className="h-4 w-16" />
      </div>
    </div>
  );
}

export function EventListingContent({
  events,
  isLoading = false,
  emptyMessage,
}: EventListingContentProps) {
  const [searchQuery, setSearchQuery] = useState("");

  const filteredEvents = useMemo(() => {
    if (!searchQuery.trim()) return events;
    const q = searchQuery.toLowerCase();
    return events.filter(
      (event) =>
        event.name.toLowerCase().includes(q) ||
        event.venue_name.toLowerCase().includes(q) ||
        event.venue_address?.toLowerCase().includes(q) ||
        event.venue_city?.toLowerCase().includes(q)
    );
  }, [events, searchQuery]);

  return (
    <>
      <div className="mt-6">
        <SearchInput
          placeholder="Search events by name, venue..."
          onSearch={setSearchQuery}
          className="max-w-md"
        />
      </div>

      {isLoading ? (
        <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <EventCardSkeleton key={i} />
          ))}
        </div>
      ) : filteredEvents.length === 0 ? (
        <div className="mt-8">
          <EmptyState
            icon={CalendarX}
            title="No events found"
            description={
              searchQuery
                ? `No events matching "${searchQuery}"`
                : emptyMessage || "There are no events available at the moment."
            }
          />
        </div>
      ) : (
        <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {filteredEvents.map((event) => (
            <EventCard key={event.id} event={event} />
          ))}
        </div>
      )}
    </>
  );
}
