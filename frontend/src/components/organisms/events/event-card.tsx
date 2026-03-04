"use client";

import Link from "next/link";
import { CalendarDays, MapPin } from "lucide-react";
import { m } from "framer-motion";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { PriceDisplay } from "@/components/molecules/price-display";
import { cn } from "@/lib/utils";
import type { Event } from "@/lib/api/types/event";

interface EventCardProps {
  event: Event;
  minPrice?: number;
  availableSeats?: number;
  className?: string;
}

export function EventCard({
  event,
  minPrice,
  availableSeats,
  className,
}: EventCardProps) {
  const formattedDate = new Date(event.start_date).toLocaleDateString("vi-VN", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });

  return (
    <m.div
      whileHover={{ y: -4 }}
      whileTap={{ scale: 0.98 }}
      transition={{ type: "spring", stiffness: 300, damping: 20 }}
    >
      <Link href={`/events/${event.id}`}>
        <Card
          className={cn(
            "group overflow-hidden transition-shadow hover:shadow-lg",
            className
          )}
        >
          <CardHeader className="p-0">
            <div className="relative aspect-[16/9] overflow-hidden bg-muted">
              <div className="flex h-full items-center justify-center gradient-brand-subtle">
                <CalendarDays className="h-12 w-12 text-primary/30" />
              </div>
            </div>
          </CardHeader>

          <CardContent className="p-4">
            <h3 className="line-clamp-1 font-semibold text-lg group-hover:text-primary transition-colors">
              {event.name}
            </h3>

            <div className="mt-2 flex flex-col gap-1">
              <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <CalendarDays className="h-3.5 w-3.5" />
                <time dateTime={event.start_date}>{formattedDate}</time>
              </div>
              <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <MapPin className="h-3.5 w-3.5" />
                <span className="line-clamp-1">
                  {event.venue_name}
                  {event.venue_city && `, ${event.venue_city}`}
                </span>
              </div>
            </div>
          </CardContent>

          <CardFooter className="flex items-center justify-between border-t px-4 py-3">
            {minPrice != null ? (
              <PriceDisplay amount={minPrice} size="sm" />
            ) : (
              <span className="text-sm text-muted-foreground">--</span>
            )}
            {availableSeats != null ? (
              <span className="text-xs text-muted-foreground">
                {availableSeats} seats left
              </span>
            ) : event.venue_capacity != null ? (
              <span className="text-xs text-muted-foreground">
                {event.venue_capacity} capacity
              </span>
            ) : null}
          </CardFooter>
        </Card>
      </Link>
    </m.div>
  );
}
