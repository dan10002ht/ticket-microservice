"use client";

import Link from "next/link";
import { CheckCircle2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { useBookingStore } from "@/stores";
import type { TicketType } from "@/lib/api/types/ticket";

interface ConfirmationStepProps {
  eventId: string;
  eventName: string;
  ticketTypes: TicketType[];
}

function formatPrice(price: number, currency = "VND") {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 0,
  }).format(price);
}

export function ConfirmationStep({
  eventId,
  eventName,
  ticketTypes,
}: ConfirmationStepProps) {
  const { booking, payment, selectedTicketTypeId, quantity } = useBookingStore();
  const selectedType = ticketTypes.find((t) => t.id === selectedTicketTypeId);

  return (
    <div className="space-y-6 text-center">
      <div className="flex flex-col items-center gap-3">
        <div className="flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
          <CheckCircle2 className="h-10 w-10 text-green-600" />
        </div>
        <h2 className="text-2xl font-bold">Booking Confirmed!</h2>
        <p className="text-muted-foreground">
          Your tickets have been reserved successfully
        </p>
      </div>

      {booking && (
        <div className="mx-auto max-w-sm space-y-3 rounded-lg border bg-card p-4 text-left">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Booking ID</span>
            <span className="font-mono text-xs">
              {booking.id.slice(0, 8).toUpperCase()}
            </span>
          </div>
          <Separator />
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Event</span>
            <span className="font-medium">{eventName}</span>
          </div>
          {selectedType && (
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Ticket</span>
              <span>
                {selectedType.name} &times; {quantity}
              </span>
            </div>
          )}
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Total</span>
            <span className="font-semibold">
              {formatPrice(
                booking.total_amount,
                booking.currency
              )}
            </span>
          </div>
          <Separator />
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Status</span>
            <Badge variant="default">Confirmed</Badge>
          </div>
          {payment && (
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Payment</span>
              <Badge variant="secondary" className="capitalize">
                {payment.payment_method.replace("_", " ")}
              </Badge>
            </div>
          )}
        </div>
      )}

      <div className="flex flex-col items-center gap-3 pt-2 sm:flex-row sm:justify-center">
        <Button asChild>
          <Link href="/my-bookings">View My Bookings</Link>
        </Button>
        <Button variant="outline" asChild>
          <Link href={`/events/${eventId}`}>Back to Event</Link>
        </Button>
      </div>
    </div>
  );
}
