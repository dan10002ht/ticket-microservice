"use client";

import { useParams } from "next/navigation";
import { EventDetailContent } from "@/components/organisms/events/event-detail-content";
import {
  useEvent,
  useEventPricing,
  useEventAvailability,
  useEventZones,
} from "@/lib/api/queries";

export default function EventDetailPage() {
  const { id } = useParams<{ id: string }>();

  const { data: event, isLoading, error } = useEvent(id);
  const { data: pricing } = useEventPricing(id);
  const { data: availability } = useEventAvailability(id);
  const { data: zones } = useEventZones(id);

  return (
    <EventDetailContent
      event={event}
      pricing={pricing}
      availability={availability}
      zones={zones}
      isLoading={isLoading}
      error={error}
    />
  );
}
