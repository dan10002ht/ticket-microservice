"use client";

import { useEffect, useCallback } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { BookingStepIndicator } from "@/components/molecules/booking-step-indicator";
import { TicketSelectionStep } from "./ticket-selection-step";
import { ReviewStep } from "./review-step";
import { PaymentStep } from "./payment-step";
import { ConfirmationStep } from "./confirmation-step";
import { useBookingStore } from "@/stores";
import { useReleaseTickets, useTicketTypes } from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import type { Event } from "@/lib/api/types/event";
import type { ApiError } from "@/lib/api/types/common";

interface BookingFlowContentProps {
  event?: Event;
  isLoading: boolean;
  error: ApiError | null;
}

export function BookingFlowContent({
  event,
  isLoading,
  error,
}: BookingFlowContentProps) {
  const { step, reservation, reset } = useBookingStore();
  const releaseMutation = useReleaseTickets();
  const { data: ticketTypes } = useTicketTypes(event?.id ?? "");

  // Cleanup on unmount: release reservation + reset store
  useEffect(() => {
    return () => {
      const { reservation: currentRes, step: currentStep } =
        useBookingStore.getState();
      if (currentRes && currentStep !== "confirmation") {
        releaseMutation.mutate({
          reservation_id: currentRes.reservation_id,
        });
      }
      // Always reset store on unmount
      useBookingStore.getState().reset();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Warn user before closing tab mid-flow
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      const { step: currentStep } = useBookingStore.getState();
      if (currentStep !== "select" && currentStep !== "confirmation") {
        e.preventDefault();
      }
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => window.removeEventListener("beforeunload", handleBeforeUnload);
  }, []);

  const handleReservationExpired = useCallback(() => {
    const { reservation: currentRes } = useBookingStore.getState();
    if (currentRes) {
      // Best-effort release (server expires it anyway)
      releaseMutation.mutate({
        reservation_id: currentRes.reservation_id,
      });
    }
    reset();
    showToast.warning("Your reservation has expired. Please try again.");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [reset]);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="mx-auto h-12 w-96" />
        <Skeleton className="h-24 w-full" />
        <Skeleton className="h-24 w-full" />
      </div>
    );
  }

  if (error || !event) {
    return (
      <div className="rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
        {error?.error?.message || "Failed to load event details. Please try again."}
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <BookingStepIndicator currentStep={step} />

      {step === "select" && <TicketSelectionStep eventId={event.id} />}

      {step === "review" && (
        <ReviewStep event={event} ticketTypes={ticketTypes ?? []} />
      )}

      {step === "payment" && (
        <PaymentStep
          event={event}
          ticketTypes={ticketTypes ?? []}
          onExpired={handleReservationExpired}
        />
      )}

      {step === "confirmation" && (
        <ConfirmationStep
          eventId={event.id}
          eventName={event.name}
          ticketTypes={ticketTypes ?? []}
        />
      )}
    </div>
  );
}
