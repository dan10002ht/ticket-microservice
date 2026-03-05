"use client";

import { PageHeader } from "@/components/molecules/page-header";
import { EventListingContent } from "@/components/organisms/events/event-listing-content";
import { useEvents } from "@/lib/api/queries";

export default function EventsPage() {
  const { data, isLoading, error } = useEvents();

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

      <EventListingContent
        events={data?.items ?? []}
        isLoading={isLoading}
      />
    </section>
  );
}
