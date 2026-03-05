"use client";

import { Users, Ticket, DoorOpen, Clock } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { useCheckInStats } from "@/lib/api/queries";

interface EventCheckinStatsProps {
  eventId: string;
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function EventCheckinStats({ eventId }: EventCheckinStatsProps) {
  const { data: stats, isLoading } = useCheckInStats(eventId);

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-3">
        <Skeleton className="h-24 rounded-lg" />
        <Skeleton className="h-24 rounded-lg" />
        <Skeleton className="h-24 rounded-lg" />
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
        No check-in data available yet.
      </div>
    );
  }

  const gates = Object.entries(stats.by_gate ?? {});

  return (
    <div className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-3">
        <div className="rounded-lg border bg-card p-4">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Users className="h-4 w-4" />
            Total Check-ins
          </div>
          <p className="mt-2 text-2xl font-bold tabular-nums">
            {stats.total_checkins}
          </p>
        </div>

        <div className="rounded-lg border bg-card p-4">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Ticket className="h-4 w-4" />
            Unique Tickets
          </div>
          <p className="mt-2 text-2xl font-bold tabular-nums">
            {stats.unique_tickets}
          </p>
        </div>

        <div className="rounded-lg border bg-card p-4">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Clock className="h-4 w-4" />
            Last Check-in
          </div>
          <p className="mt-2 text-sm font-medium">
            {stats.last_checkin_at
              ? formatDate(stats.last_checkin_at)
              : "No check-ins yet"}
          </p>
        </div>
      </div>

      {gates.length > 0 && (
        <div className="rounded-lg border bg-card p-4">
          <h4 className="flex items-center gap-2 text-sm font-semibold text-muted-foreground">
            <DoorOpen className="h-4 w-4" />
            By Gate
          </h4>
          <div className="mt-3 space-y-2">
            {gates.map(([gate, count]) => (
              <div key={gate} className="flex items-center justify-between text-sm">
                <span>{gate}</span>
                <span className="font-medium tabular-nums">{count}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
