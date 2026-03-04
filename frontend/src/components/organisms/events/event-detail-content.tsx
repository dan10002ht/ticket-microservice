"use client";

import Link from "next/link";
import {
  CalendarDays,
  MapPin,
  Users,
  ArrowLeft,
  Ticket,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { PriceDisplay } from "@/components/molecules/price-display";
import type { Event, Pricing, Availability, EventSeatingZone } from "@/lib/api/types/event";
import type { ApiError } from "@/lib/api/types/common";

interface EventDetailContentProps {
  event?: Event;
  pricing?: Pricing[];
  availability?: Availability;
  zones?: EventSeatingZone[];
  isLoading: boolean;
  error: ApiError | null;
}

function DetailSkeleton() {
  return (
    <div className="container mx-auto px-4 py-8 space-y-8">
      <Skeleton className="h-8 w-32" />
      <Skeleton className="aspect-[21/9] w-full rounded-xl" />
      <div className="grid gap-8 lg:grid-cols-3">
        <div className="lg:col-span-2 space-y-4">
          <Skeleton className="h-10 w-3/4" />
          <Skeleton className="h-5 w-1/3" />
          <Skeleton className="h-5 w-1/4" />
          <Skeleton className="h-32 w-full" />
        </div>
        <div className="space-y-4">
          <Skeleton className="h-48 w-full rounded-lg" />
          <Skeleton className="h-12 w-full rounded-lg" />
        </div>
      </div>
    </div>
  );
}

export function EventDetailContent({
  event,
  pricing,
  availability,
  zones,
  isLoading,
  error,
}: EventDetailContentProps) {
  if (isLoading) return <DetailSkeleton />;

  if (error || !event) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="flex flex-col items-center justify-center text-center">
          <div className="rounded-full bg-muted p-4">
            <CalendarDays className="h-8 w-8 text-muted-foreground" />
          </div>
          <h3 className="mt-4 text-lg font-semibold">Event not found</h3>
          <p className="mt-2 max-w-sm text-sm text-muted-foreground">
            {error?.error || "The event you're looking for doesn't exist or has been removed."}
          </p>
          <Button asChild variant="outline" className="mt-4">
            <Link href="/events">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Events
            </Link>
          </Button>
        </div>
      </div>
    );
  }

  const formattedStartDate = new Date(event.start_date).toLocaleDateString(
    "vi-VN",
    { weekday: "long", day: "numeric", month: "long", year: "numeric" }
  );
  const formattedStartTime = new Date(event.start_date).toLocaleTimeString(
    "vi-VN",
    { hour: "2-digit", minute: "2-digit" }
  );
  const formattedEndTime = event.end_date
    ? new Date(event.end_date).toLocaleTimeString("vi-VN", {
        hour: "2-digit",
        minute: "2-digit",
      })
    : null;

  const minPrice = pricing?.length
    ? Math.min(...pricing.map((p) => p.price))
    : null;
  const maxPrice = pricing?.length
    ? Math.max(...pricing.map((p) => p.price))
    : null;

  const venueDisplay = [event.venue_name, event.venue_city, event.venue_country]
    .filter(Boolean)
    .join(", ");

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Back button */}
      <Button asChild variant="ghost" className="mb-6 -ml-2">
        <Link href="/events">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Events
        </Link>
      </Button>

      {/* Hero */}
      <div className="relative aspect-[21/9] overflow-hidden rounded-xl bg-muted">
        <div className="flex h-full items-center justify-center gradient-brand-subtle">
          <CalendarDays className="h-24 w-24 text-primary/20" />
        </div>
      </div>

      {/* Content grid */}
      <div className="mt-8 grid gap-8 lg:grid-cols-3">
        {/* Left: Event info */}
        <div className="lg:col-span-2 space-y-6">
          <div>
            <h1 className="text-3xl font-bold">{event.name}</h1>

            <div className="mt-4 flex flex-col gap-2">
              <div className="flex items-center gap-2 text-muted-foreground">
                <CalendarDays className="h-4 w-4" />
                <span>
                  {formattedStartDate} &middot; {formattedStartTime}
                  {formattedEndTime && ` – ${formattedEndTime}`}
                </span>
              </div>
              <div className="flex items-center gap-2 text-muted-foreground">
                <MapPin className="h-4 w-4" />
                <span>{venueDisplay}</span>
              </div>
              {event.venue_capacity && (
                <div className="flex items-center gap-2 text-muted-foreground">
                  <Users className="h-4 w-4" />
                  <span>{event.venue_capacity.toLocaleString()} capacity</span>
                </div>
              )}
            </div>
          </div>

          <Separator />

          {/* Description */}
          {event.description && (
            <div>
              <h2 className="text-lg font-semibold mb-2">About this event</h2>
              <p className="text-muted-foreground whitespace-pre-line">
                {event.description}
              </p>
            </div>
          )}

          {/* Zones & Pricing */}
          {zones && zones.length > 0 && (
            <div>
              <h2 className="text-lg font-semibold mb-4">Zones & Pricing</h2>
              <div className="grid gap-4 sm:grid-cols-2">
                {zones.map((zone) => {
                  const zonePricing = pricing?.filter(
                    (p) => p.zone_id === zone.id
                  );
                  return (
                    <Card key={zone.id}>
                      <CardHeader className="pb-2">
                        <CardTitle className="flex items-center gap-2 text-base">
                          {zone.color && (
                            <span
                              className="h-3 w-3 rounded-full"
                              style={{ backgroundColor: zone.color }}
                            />
                          )}
                          {zone.name}
                          <Badge variant="secondary" className="ml-auto capitalize">
                            {zone.zone_type}
                          </Badge>
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="pt-0">
                        {zone.seat_count != null && (
                          <p className="text-sm text-muted-foreground">
                            {zone.seat_count} seats
                          </p>
                        )}
                        {zonePricing && zonePricing.length > 0 && (
                          <div className="mt-2 space-y-1">
                            {zonePricing.map((p) => (
                              <div
                                key={p.id}
                                className="flex items-center justify-between text-sm"
                              >
                                <span>{p.name}</span>
                                <PriceDisplay
                                  amount={p.price}
                                  currency={p.currency}
                                  size="sm"
                                />
                              </div>
                            ))}
                          </div>
                        )}
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            </div>
          )}
        </div>

        {/* Right: Booking sidebar */}
        <div className="space-y-4">
          <Card>
            <CardContent className="py-6 space-y-4">
              {/* Price range */}
              {minPrice != null && (
                <div>
                  <p className="text-sm text-muted-foreground">
                    {minPrice === maxPrice ? "Price" : "Starting from"}
                  </p>
                  <PriceDisplay amount={minPrice} size="lg" />
                  {maxPrice != null && maxPrice !== minPrice && (
                    <span className="ml-1 text-sm text-muted-foreground">
                      – <PriceDisplay amount={maxPrice} size="sm" />
                    </span>
                  )}
                </div>
              )}

              {/* Availability */}
              {availability && (
                <div className="space-y-2">
                  <Separator />
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <div>
                      <p className="text-muted-foreground">Available</p>
                      <p className="font-semibold text-green-600">
                        {availability.available_seats}
                      </p>
                    </div>
                    <div>
                      <p className="text-muted-foreground">Total</p>
                      <p className="font-semibold">
                        {availability.total_seats}
                      </p>
                    </div>
                    <div>
                      <p className="text-muted-foreground">Reserved</p>
                      <p className="font-semibold text-yellow-600">
                        {availability.reserved_seats}
                      </p>
                    </div>
                    <div>
                      <p className="text-muted-foreground">Sold</p>
                      <p className="font-semibold text-blue-600">
                        {availability.sold_seats}
                      </p>
                    </div>
                  </div>
                </div>
              )}

              <Separator />

              <Button className="w-full" size="lg" asChild>
                <Link href={`/events/${event.id}/book`}>
                  <Ticket className="mr-2 h-4 w-4" />
                  Book Now
                </Link>
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
