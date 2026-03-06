"use client";

import { toast } from "sonner";
import { Calendar, MapPin, Users, Tag, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { StatusBadge } from "@/components/molecules/status-badge";
import {
  useEvent,
  useEventZones,
  useEventPricing,
  usePublishEvent,
} from "@/lib/api/queries";
import type { ApiError } from "@/lib/api/types/common";
import { cn } from "@/lib/utils";

interface EventReviewStepProps {
  eventId: string;
  onDone: () => void;
  className?: string;
}

export function EventReviewStep({
  eventId,
  onDone,
  className,
}: EventReviewStepProps) {
  const { data: event, isLoading: eventLoading } = useEvent(eventId);
  const { data: zones = [] } = useEventZones(eventId);
  const { data: pricing = [] } = useEventPricing(eventId);
  const publishMutation = usePublishEvent();

  const handlePublish = async () => {
    try {
      await publishMutation.mutateAsync(eventId);
      toast.success("Event published successfully!");
      onDone();
    } catch (err: unknown) {
      const message =
        (err as ApiError)?.error?.message || "Failed to publish event.";
      toast.error(message);
    }
  };

  const handleSaveDraft = () => {
    toast.success("Event saved as draft.");
    onDone();
  };

  if (eventLoading) {
    return (
      <div className={cn("flex items-center justify-center py-12", className)}>
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!event) return null;

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleString();
    } catch {
      return dateStr;
    }
  };

  return (
    <div className={cn("space-y-6", className)}>
      {/* Event Summary */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Event Summary</CardTitle>
            <StatusBadge status={event.status ?? "draft"} />
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <h3 className="text-lg font-semibold">{event.name}</h3>
            {event.description && (
              <p className="mt-1 text-sm text-muted-foreground">
                {event.description}
              </p>
            )}
          </div>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="flex items-start gap-2">
              <Calendar className="mt-0.5 h-4 w-4 text-muted-foreground shrink-0" />
              <div className="text-sm">
                <p>{formatDate(event.start_date)}</p>
                {event.end_date && (
                  <p className="text-muted-foreground">
                    to {formatDate(event.end_date)}
                  </p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-2">
              <MapPin className="mt-0.5 h-4 w-4 text-muted-foreground shrink-0" />
              <div className="text-sm">
                <p>{event.venue_name}</p>
                <p className="text-muted-foreground">
                  {[event.venue_address, event.venue_city, event.venue_country]
                    .filter(Boolean)
                    .join(", ") || "No address"}
                </p>
              </div>
            </div>

            {event.venue_capacity && (
              <div className="flex items-center gap-2">
                <Users className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm">
                  Capacity: {event.venue_capacity.toLocaleString()}
                </span>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Zones & Pricing Summary */}
      <Card>
        <CardHeader>
          <CardTitle>
            Zones & Pricing ({zones.length} zone{zones.length !== 1 ? "s" : ""}
            )
          </CardTitle>
        </CardHeader>
        <CardContent>
          {zones.length === 0 ? (
            <p className="text-sm text-muted-foreground italic">
              No zones configured. You can add zones later from the event detail
              page.
            </p>
          ) : (
            <div className="space-y-3">
              {zones.map((zone) => {
                const zonePricing = pricing.filter(
                  (p) => p.zone_id === zone.id
                );
                return (
                  <div
                    key={zone.id}
                    className="flex items-center gap-3 rounded-lg border p-3"
                  >
                    <div
                      className="h-4 w-4 rounded-full shrink-0"
                      style={{ backgroundColor: zone.color || "#3B82F6" }}
                    />
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-sm">{zone.name}</p>
                      <p className="text-xs text-muted-foreground">
                        {zone.zone_type} &middot;{" "}
                        {zone.seat_count ?? "N/A"} seats
                      </p>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <Tag className="h-3.5 w-3.5 text-muted-foreground" />
                      {zonePricing.length > 0 ? (
                        <div className="text-xs text-right">
                          {zonePricing.map((p) => (
                            <span key={p.id} className="block">
                              {p.name}: {p.currency} {p.price}
                            </span>
                          ))}
                        </div>
                      ) : (
                        <span className="text-xs text-muted-foreground italic">
                          No pricing
                        </span>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Actions */}
      <div className="flex justify-end gap-3">
        <Button variant="outline" onClick={handleSaveDraft}>
          Save as Draft
        </Button>
        <Button onClick={handlePublish} disabled={publishMutation.isPending}>
          {publishMutation.isPending ? "Publishing..." : "Publish Event"}
        </Button>
      </div>
    </div>
  );
}
