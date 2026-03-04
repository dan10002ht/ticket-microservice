"use client";

import { Minus, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import type { TicketType } from "@/lib/api/types/ticket";

interface TicketTypeCardProps {
  ticketType: TicketType;
  selectedQuantity: number;
  onQuantityChange: (quantity: number) => void;
  disabled?: boolean;
}

function formatPrice(price: number, currency = "VND") {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 0,
  }).format(price);
}

export function TicketTypeCard({
  ticketType,
  selectedQuantity,
  onQuantityChange,
  disabled = false,
}: TicketTypeCardProps) {
  const isSoldOut =
    ticketType.status === "sold_out" || ticketType.available_quantity === 0;
  const isInactive = ticketType.status === "inactive";
  const isUnavailable = isSoldOut || isInactive || disabled;

  const min = ticketType.min_per_purchase || 1;
  const max = Math.min(
    ticketType.max_per_purchase || 10,
    ticketType.available_quantity
  );

  const handleDecrement = () => {
    if (selectedQuantity <= min) {
      onQuantityChange(0); // deselect
    } else {
      onQuantityChange(selectedQuantity - 1);
    }
  };

  const handleIncrement = () => {
    if (selectedQuantity === 0) {
      onQuantityChange(min);
    } else if (selectedQuantity < max) {
      onQuantityChange(selectedQuantity + 1);
    }
  };

  const isSelected = selectedQuantity > 0;

  return (
    <div
      className={`rounded-lg border p-4 transition-colors ${
        isSelected
          ? "border-primary bg-primary/5"
          : isUnavailable
            ? "border-muted bg-muted/30 opacity-60"
            : "border-border hover:border-primary/50"
      }`}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 space-y-1">
          <div className="flex items-center gap-2">
            <h3 className="font-semibold">{ticketType.name}</h3>
            {isSoldOut && <Badge variant="destructive">Sold Out</Badge>}
            {isInactive && <Badge variant="secondary">Unavailable</Badge>}
          </div>
          {ticketType.description && (
            <p className="text-sm text-muted-foreground">
              {ticketType.description}
            </p>
          )}
          <p className="text-sm text-muted-foreground">
            {ticketType.available_quantity} of {ticketType.quantity} available
          </p>
        </div>

        <div className="text-right">
          <p className="text-lg font-bold">
            {formatPrice(ticketType.price, ticketType.currency)}
          </p>
        </div>
      </div>

      {!isUnavailable && (
        <div className="mt-3 flex items-center justify-end gap-2">
          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8"
            onClick={handleDecrement}
            disabled={selectedQuantity === 0}
          >
            <Minus className="h-3 w-3" />
          </Button>
          <span className="w-8 text-center text-sm font-medium tabular-nums">
            {selectedQuantity}
          </span>
          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8"
            onClick={handleIncrement}
            disabled={selectedQuantity >= max}
          >
            <Plus className="h-3 w-3" />
          </Button>
        </div>
      )}
    </div>
  );
}
