"use client";

import { useState } from "react";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { TicketTypeCard } from "@/components/molecules/ticket-type-card";
import { useTicketTypes } from "@/lib/api/queries";
import { useBookingStore } from "@/stores";

interface TicketSelectionStepProps {
  eventId: string;
}

export function TicketSelectionStep({ eventId }: TicketSelectionStepProps) {
  const { data: ticketTypes, isLoading, error } = useTicketTypes(eventId);
  const { selectedTicketTypeId, quantity, setTicketSelection, goToStep } =
    useBookingStore();

  // Local selection state synced to store on "Continue"
  const [localTypeId, setLocalTypeId] = useState(selectedTicketTypeId ?? "");
  const [localQty, setLocalQty] = useState(quantity);

  const handleQuantityChange = (typeId: string, qty: number) => {
    if (qty === 0) {
      // Deselect
      setLocalTypeId("");
      setLocalQty(1);
    } else {
      setLocalTypeId(typeId);
      setLocalQty(qty);
    }
  };

  const handleContinue = () => {
    if (!localTypeId || localQty < 1) return;
    setTicketSelection(localTypeId, localQty);
    goToStep("review");
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-24 w-full rounded-lg" />
        <Skeleton className="h-24 w-full rounded-lg" />
        <Skeleton className="h-24 w-full rounded-lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
        Failed to load ticket types. Please try again.
      </div>
    );
  }

  if (!ticketTypes || ticketTypes.length === 0) {
    return (
      <div className="py-8 text-center text-muted-foreground">
        No ticket types available for this event.
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-lg font-semibold">Select Tickets</h2>
        <p className="text-sm text-muted-foreground">
          Choose your ticket type and quantity
        </p>
      </div>

      <div className="space-y-3">
        {ticketTypes.map((tt) => (
          <TicketTypeCard
            key={tt.id}
            ticketType={tt}
            selectedQuantity={localTypeId === tt.id ? localQty : 0}
            onQuantityChange={(qty) => handleQuantityChange(tt.id, qty)}
          />
        ))}
      </div>

      <div className="flex justify-end pt-2">
        <Button
          onClick={handleContinue}
          disabled={!localTypeId || localQty < 1}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
