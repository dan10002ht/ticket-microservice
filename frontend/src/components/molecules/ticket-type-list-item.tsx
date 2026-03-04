"use client";

import { useState } from "react";
import { Pencil, Trash2, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import type { TicketType } from "@/lib/api/types/ticket";

interface TicketTypeListItemProps {
  ticketType: TicketType;
  onEdit: () => void;
  onDelete: () => void;
  isDeleting?: boolean;
}

function formatPrice(price: number, currency = "VND") {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 0,
  }).format(price);
}

const statusVariant: Record<string, "default" | "secondary" | "destructive"> = {
  active: "default",
  inactive: "secondary",
  sold_out: "destructive",
};

export function TicketTypeListItem({
  ticketType,
  onEdit,
  onDelete,
  isDeleting = false,
}: TicketTypeListItemProps) {
  const [confirmDelete, setConfirmDelete] = useState(false);

  const handleDelete = () => {
    if (!confirmDelete) {
      setConfirmDelete(true);
      return;
    }
    onDelete();
  };

  return (
    <div className="flex items-center justify-between rounded-lg border px-4 py-3">
      <div className="flex-1 space-y-0.5">
        <div className="flex items-center gap-2">
          <span className="font-medium">{ticketType.name}</span>
          <Badge variant={statusVariant[ticketType.status] ?? "secondary"} className="capitalize">
            {ticketType.status.replace("_", " ")}
          </Badge>
        </div>
        <div className="flex items-center gap-3 text-sm text-muted-foreground">
          <span>{formatPrice(ticketType.price, ticketType.currency)}</span>
          <span>
            {ticketType.available_quantity}/{ticketType.quantity} available
          </span>
        </div>
      </div>

      <div className="flex items-center gap-1">
        <Button variant="ghost" size="icon" className="h-8 w-8" onClick={onEdit}>
          <Pencil className="h-3.5 w-3.5" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className={`h-8 w-8 ${confirmDelete ? "text-destructive" : ""}`}
          onClick={handleDelete}
          onBlur={() => setConfirmDelete(false)}
          disabled={isDeleting}
        >
          {isDeleting ? (
            <Loader2 className="h-3.5 w-3.5 animate-spin" />
          ) : (
            <Trash2 className="h-3.5 w-3.5" />
          )}
        </Button>
      </div>
    </div>
  );
}
