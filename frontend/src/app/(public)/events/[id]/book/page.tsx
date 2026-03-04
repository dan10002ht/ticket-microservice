"use client";

import { useParams } from "next/navigation";
import { BookingFlowContent } from "@/components/organisms/booking/booking-flow-content";
import { useEvent } from "@/lib/api/queries";

export default function BookEventPage() {
  const params = useParams();
  const eventId = params.id as string;
  const { data: event, isLoading, error } = useEvent(eventId);

  return (
    <section className="container mx-auto max-w-2xl px-4 py-8">
      <BookingFlowContent event={event} isLoading={isLoading} error={error} />
    </section>
  );
}
