"use client";

import { Separator } from "@/components/ui/separator";

interface OrderSummaryCardProps {
  eventName: string;
  ticketTypeName: string;
  quantity: number;
  unitPrice: number;
  currency?: string;
}

function formatPrice(price: number, currency = "VND") {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 0,
  }).format(price);
}

export function OrderSummaryCard({
  eventName,
  ticketTypeName,
  quantity,
  unitPrice,
  currency = "VND",
}: OrderSummaryCardProps) {
  const total = unitPrice * quantity;

  return (
    <div className="rounded-lg border bg-card p-4">
      <h3 className="font-semibold">Order Summary</h3>
      <Separator className="my-3" />

      <div className="space-y-2 text-sm">
        <div>
          <p className="font-medium">{eventName}</p>
          <p className="text-muted-foreground">{ticketTypeName}</p>
        </div>

        <div className="flex justify-between">
          <span className="text-muted-foreground">
            {formatPrice(unitPrice, currency)} &times; {quantity}
          </span>
          <span>{formatPrice(total, currency)}</span>
        </div>
      </div>

      <Separator className="my-3" />

      <div className="flex justify-between font-semibold">
        <span>Total</span>
        <span>{formatPrice(total, currency)}</span>
      </div>
    </div>
  );
}
