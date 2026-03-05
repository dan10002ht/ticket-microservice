"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { EventCheckinStats } from "@/components/organisms/org-events/event-checkin-stats";
import { CheckinForm } from "@/components/organisms/org-events/checkin-form";
import { useEvent } from "@/lib/api/queries";

export default function OrgEventCheckinPage() {
  const params = useParams();
  const eventId = params.id as string;
  const { data: event } = useEvent(eventId);

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" asChild>
          <Link href={`/org/events/${eventId}`}>
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-xl font-bold">Check-in</h1>
          {event && (
            <p className="text-sm text-muted-foreground">{event.name}</p>
          )}
        </div>
      </div>

      <EventCheckinStats eventId={eventId} />

      <Separator />

      <CheckinForm eventId={eventId} />
    </div>
  );
}
