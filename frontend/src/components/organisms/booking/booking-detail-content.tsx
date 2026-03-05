"use client";

import Link from "next/link";
import {
  CalendarDays,
  Ticket,
  CreditCard,
  ArrowLeft,
  Loader2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { useCancelBooking } from "@/lib/api/queries";
import { InvoiceSection } from "@/components/molecules/invoice-section";
import { showToast } from "@/lib/toast";
import type { Booking } from "@/lib/api/types/booking";
import type { Payment } from "@/lib/api/types/payment";
import type { ApiError } from "@/lib/api/types/common";

interface BookingDetailContentProps {
  booking?: Booking;
  payment?: Payment;
  isLoading: boolean;
  error: ApiError | null;
}

function formatPrice(price: number, currency = "VND") {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 0,
  }).format(price);
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

const statusVariant: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  confirmed: "default",
  completed: "default",
  pending: "secondary",
  cancelled: "destructive",
  expired: "destructive",
};

export function BookingDetailContent({
  booking,
  payment,
  isLoading,
  error,
}: BookingDetailContentProps) {
  const cancelMutation = useCancelBooking();

  const handleCancel = async () => {
    if (!booking) return;
    try {
      await cancelMutation.mutateAsync({ id: booking.id });
      showToast.success("Booking cancelled successfully.");
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full rounded-lg" />
      </div>
    );
  }

  if (error || !booking) {
    return (
      <div className="rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
        {error?.error?.message || "Booking not found."}
      </div>
    );
  }

  const canCancel = booking.status === "pending";

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/my-bookings">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-xl font-bold">Booking Details</h1>
          <p className="font-mono text-sm text-muted-foreground">
            {booking.id.slice(0, 8).toUpperCase()}
          </p>
        </div>
        <Badge variant={statusVariant[booking.status] ?? "outline"} className="ml-auto capitalize">
          {booking.status}
        </Badge>
      </div>

      <div className="rounded-lg border bg-card">
        <div className="p-4 space-y-4">
          {/* Booking Info */}
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="flex items-start gap-3">
              <Ticket className="mt-0.5 h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Tickets</p>
                <p className="font-medium">{booking.ticket_quantity} ticket(s)</p>
                {booking.seat_numbers && booking.seat_numbers.length > 0 && (
                  <p className="text-xs text-muted-foreground">
                    Seats: {booking.seat_numbers.join(", ")}
                  </p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <CreditCard className="mt-0.5 h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Total Amount</p>
                <p className="font-medium">
                  {formatPrice(booking.total_amount, booking.currency)}
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <CalendarDays className="mt-0.5 h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Booked On</p>
                <p className="font-medium">{formatDate(booking.created_at)}</p>
              </div>
            </div>

            {booking.expires_at && (
              <div className="flex items-start gap-3">
                <CalendarDays className="mt-0.5 h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Expires</p>
                  <p className="font-medium">{formatDate(booking.expires_at)}</p>
                </div>
              </div>
            )}
          </div>

          {booking.special_requests && (
            <>
              <Separator />
              <div>
                <p className="text-sm text-muted-foreground">Special Requests</p>
                <p className="mt-1 text-sm">{booking.special_requests}</p>
              </div>
            </>
          )}
        </div>

        {/* Payment Info */}
        {payment && (
          <>
            <Separator />
            <div className="p-4 space-y-2">
              <h3 className="text-sm font-semibold">Payment</h3>
              <div className="grid gap-2 text-sm sm:grid-cols-2">
                <div>
                  <span className="text-muted-foreground">Method: </span>
                  <span className="capitalize">
                    {payment.payment_method.replace("_", " ")}
                  </span>
                </div>
                <div>
                  <span className="text-muted-foreground">Status: </span>
                  <Badge
                    variant={payment.status === "completed" ? "default" : "secondary"}
                    className="capitalize"
                  >
                    {payment.status}
                  </Badge>
                </div>
                {payment.transaction_id && (
                  <div>
                    <span className="text-muted-foreground">Transaction: </span>
                    <span className="font-mono text-xs">
                      {payment.transaction_id}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </>
        )}

        {/* Invoice */}
        <InvoiceSection bookingId={booking.id} />
      </div>

      {/* Actions */}
      <div className="flex gap-3">
        <Button variant="outline" asChild>
          <Link href="/my-bookings">Back to Bookings</Link>
        </Button>
        {canCancel && (
          <Button
            variant="destructive"
            onClick={handleCancel}
            disabled={cancelMutation.isPending}
          >
            {cancelMutation.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            Cancel Booking
          </Button>
        )}
      </div>
    </div>
  );
}
