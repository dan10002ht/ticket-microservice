"use client";

import { useParams } from "next/navigation";
import { EventEditForm } from "@/components/organisms/org-events/event-edit-form";
import { useEvent } from "@/lib/api/queries";

export default function OrgEventEditPage() {
  const params = useParams();
  const eventId = params.id as string;
  const { data: event, isLoading, error } = useEvent(eventId);

  return <EventEditForm event={event} isLoading={isLoading} error={error} />;
}
