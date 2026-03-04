"use client";

import { useState } from "react";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { OrderSummaryCard } from "@/components/molecules/order-summary-card";
import { useReserveTickets } from "@/lib/api/queries";
import { useCreateBooking } from "@/lib/api/queries";
import { useBookingStore } from "@/stores";
import { showToast } from "@/lib/toast";
import type { ApiError } from "@/lib/api/types/common";
import type { TicketType } from "@/lib/api/types/ticket";
import type { Event } from "@/lib/api/types/event";

interface ReviewStepProps {
  event: Event;
  ticketTypes: TicketType[];
}

export function ReviewStep({ event, ticketTypes }: ReviewStepProps) {
  const {
    selectedTicketTypeId,
    quantity,
    specialRequests,
    setSpecialRequests,
    setReservation,
    setBooking,
    goToStep,
  } = useBookingStore();

  const reserveMutation = useReserveTickets(event.id);
  const createBookingMutation = useCreateBooking();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const selectedType = ticketTypes.find((t) => t.id === selectedTicketTypeId);

  if (!selectedType) {
    goToStep("select");
    return null;
  }

  const handleConfirm = async () => {
    if (isSubmitting) return;
    setIsSubmitting(true);

    try {
      // 1. Reserve tickets
      const reservation = await reserveMutation.mutateAsync({
        ticket_type_id: selectedType.id,
        quantity,
        timeout_seconds: 600,
      });
      setReservation(reservation);

      // 2. Create booking
      const booking = await createBookingMutation.mutateAsync({
        event_id: event.id,
        ticket_quantity: quantity,
        special_requests: specialRequests || undefined,
        idempotency_key: crypto.randomUUID(),
      });
      setBooking(booking);

      // 3. Advance to payment
      goToStep("payment");
    } catch (err) {
      showToast.apiError(err as ApiError);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold">Review Your Order</h2>
        <p className="text-sm text-muted-foreground">
          Confirm your selection before reserving
        </p>
      </div>

      <OrderSummaryCard
        eventName={event.name}
        ticketTypeName={selectedType.name}
        quantity={quantity}
        unitPrice={selectedType.price}
        currency={selectedType.currency}
      />

      <div className="space-y-2">
        <Label htmlFor="special-requests">
          Special Requests{" "}
          <span className="text-muted-foreground">(optional)</span>
        </Label>
        <Textarea
          id="special-requests"
          placeholder="Any special requirements or notes..."
          value={specialRequests}
          onChange={(e) => setSpecialRequests(e.target.value)}
          maxLength={500}
          rows={3}
        />
        <p className="text-xs text-muted-foreground">
          {specialRequests.length}/500
        </p>
      </div>

      <div className="flex items-center justify-between pt-2">
        <Button variant="outline" onClick={() => goToStep("select")}>
          Back
        </Button>
        <Button onClick={handleConfirm} disabled={isSubmitting}>
          {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Confirm & Reserve
        </Button>
      </div>
    </div>
  );
}
