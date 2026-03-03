"use client";

import { useState, useMemo } from "react";
import { SearchInput } from "@/components/molecules/search-input";
import { EventCard } from "@/components/organisms/events/event-card";
import { EmptyState } from "@/components/molecules/empty-state";
import { CalendarX } from "lucide-react";
import type { Event } from "@/types";

interface EventListingContentProps {
  events: Event[];
}

export function EventListingContent({ events }: EventListingContentProps) {
  const [searchQuery, setSearchQuery] = useState("");

  const filteredEvents = useMemo(() => {
    if (!searchQuery.trim()) return events;
    const q = searchQuery.toLowerCase();
    return events.filter(
      (event) =>
        event.title.toLowerCase().includes(q) ||
        event.venue.toLowerCase().includes(q) ||
        event.address?.toLowerCase().includes(q) ||
        event.category?.toLowerCase().includes(q)
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

      {filteredEvents.length === 0 ? (
        <div className="mt-8">
          <EmptyState
            icon={CalendarX}
            title="No events found"
            description={
              searchQuery
                ? `No events matching "${searchQuery}"`
                : "There are no events available at the moment."
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
