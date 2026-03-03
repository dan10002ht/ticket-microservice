"use client";

import Link from "next/link";
import Image from "next/image";
import { CalendarDays, MapPin } from "lucide-react";
import { m } from "framer-motion";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { PriceDisplay } from "@/components/molecules/price-display";
import { cn } from "@/lib/utils";
import type { Event } from "@/types";

interface EventCardProps {
  event: Event;
  className?: string;
}

export function EventCard({ event, className }: EventCardProps) {
  const formattedDate = new Date(event.startDate).toLocaleDateString("vi-VN", {
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
              {event.imageUrl ? (
                <Image
                  src={event.imageUrl}
                  alt={event.title}
                  fill
                  className="object-cover transition-transform group-hover:scale-105"
                  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                />
              ) : (
                <div className="flex h-full items-center justify-center gradient-brand-subtle">
                  <CalendarDays className="h-12 w-12 text-primary/30" />
                </div>
              )}
              {event.category && (
                <Badge className="absolute left-3 top-3">{event.category}</Badge>
              )}
            </div>
          </CardHeader>

          <CardContent className="p-4">
            <h3 className="line-clamp-1 font-semibold text-lg group-hover:text-primary transition-colors">
              {event.title}
            </h3>

            <div className="mt-2 flex flex-col gap-1">
              <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <CalendarDays className="h-3.5 w-3.5" />
                <time dateTime={event.startDate}>{formattedDate}</time>
              </div>
              <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <MapPin className="h-3.5 w-3.5" />
                <span className="line-clamp-1">{event.venue}</span>
              </div>
            </div>
          </CardContent>

          <CardFooter className="flex items-center justify-between border-t px-4 py-3">
            <PriceDisplay amount={event.minPrice} size="sm" />
            <span className="text-xs text-muted-foreground">
              {event.availableCapacity} seats left
            </span>
          </CardFooter>
        </Card>
      </Link>
    </m.div>
  );
}
