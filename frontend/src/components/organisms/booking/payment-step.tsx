"use client";

import { useState } from "react";
import { CreditCard, Building, Wallet, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { OrderSummaryCard } from "@/components/molecules/order-summary-card";
import { ReservationTimer } from "@/components/molecules/reservation-timer";
import { useCreatePayment, useConfirmBooking } from "@/lib/api/queries";
import { useBookingStore } from "@/stores";
import { showToast } from "@/lib/toast";
import type { ApiError } from "@/lib/api/types/common";
import type { PaymentMethod } from "@/lib/api/types/payment";
import type { TicketType } from "@/lib/api/types/ticket";
import type { Event } from "@/lib/api/types/event";

interface PaymentStepProps {
  event: Event;
  ticketTypes: TicketType[];
  onExpired: () => void;
}

const PAYMENT_METHODS: {
  value: PaymentMethod;
  label: string;
  icon: typeof CreditCard;
}[] = [
  { value: "credit_card", label: "Credit Card", icon: CreditCard },
  { value: "debit_card", label: "Debit Card", icon: CreditCard },
  { value: "bank_transfer", label: "Bank Transfer", icon: Building },
  { value: "digital_wallet", label: "Digital Wallet", icon: Wallet },
];

export function PaymentStep({ event, ticketTypes, onExpired }: PaymentStepProps) {
  const {
    selectedTicketTypeId,
    quantity,
    reservation,
    booking,
    setPayment,
    goToStep,
  } = useBookingStore();

  const createPaymentMutation = useCreatePayment();
  const confirmBookingMutation = useConfirmBooking();
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>("credit_card");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const selectedType = ticketTypes.find((t) => t.id === selectedTicketTypeId);

  if (!booking || !selectedType) {
    goToStep("select");
    return null;
  }

  const handlePay = async () => {
    if (isSubmitting) return;
    setIsSubmitting(true);

    try {
      // 1. Create payment
      const payment = await createPaymentMutation.mutateAsync({
        booking_id: booking.id,
        amount: booking.total_amount,
        currency: booking.currency || selectedType.currency || "VND",
        payment_method: paymentMethod,
        idempotency_key: crypto.randomUUID(),
      });

      // 2. Confirm booking
      await confirmBookingMutation.mutateAsync({
        id: booking.id,
        payment_reference: payment.id,
      });

      setPayment(payment);
      goToStep("confirmation");
    } catch (err) {
      showToast.apiError(err as ApiError);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold">Payment</h2>
        <p className="text-sm text-muted-foreground">
          Choose your payment method to complete the booking
        </p>
      </div>

      <ReservationTimer
        expiresAt={reservation?.expires_at ?? null}
        onExpired={onExpired}
      />

      <OrderSummaryCard
        eventName={event.name}
        ticketTypeName={selectedType.name}
        quantity={quantity}
        unitPrice={selectedType.price}
        currency={selectedType.currency}
      />

      <div className="space-y-3">
        <Label>Payment Method</Label>
        <RadioGroup
          value={paymentMethod}
          onValueChange={(val) => setPaymentMethod(val as PaymentMethod)}
          className="space-y-2"
        >
          {PAYMENT_METHODS.map((method) => {
            const Icon = method.icon;
            return (
              <label
                key={method.value}
                className={`flex cursor-pointer items-center gap-3 rounded-lg border p-3 transition-colors ${
                  paymentMethod === method.value
                    ? "border-primary bg-primary/5"
                    : "border-border hover:border-primary/50"
                }`}
              >
                <RadioGroupItem value={method.value} />
                <Icon className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm font-medium">{method.label}</span>
              </label>
            );
          })}
        </RadioGroup>
      </div>

      <div className="flex justify-end pt-2">
        <Button onClick={handlePay} disabled={isSubmitting}>
          {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Pay Now
        </Button>
      </div>
    </div>
  );
}
