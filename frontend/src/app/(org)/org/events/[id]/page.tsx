"use client";

import { useParams } from "next/navigation";
import { EventDetailOrgContent } from "@/components/organisms/org-events/event-detail-org-content";
import { useEvent, useTicketTypes } from "@/lib/api/queries";

export default function OrgEventDetailPage() {
  const params = useParams();
  const eventId = params.id as string;
  const { data: event, isLoading, error } = useEvent(eventId);
  const { data: ticketTypes } = useTicketTypes(eventId);

  return (
    <EventDetailOrgContent
      event={event}
      ticketTypes={ticketTypes}
      isLoading={isLoading}
      error={error}
    />
  );
}
