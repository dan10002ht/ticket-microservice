"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useCreateTicketType, useUpdateTicketType } from "@/lib/api/queries";
import { ticketTypeSchema, type TicketTypeInput } from "@/lib/validators/event";
import { showToast } from "@/lib/toast";
import type { TicketType } from "@/lib/api/types/ticket";
import type { ApiError } from "@/lib/api/types/common";

interface TicketTypeFormProps {
  eventId: string;
  ticketType?: TicketType;
  onSuccess: () => void;
  onCancel: () => void;
}

export function TicketTypeForm({
  eventId,
  ticketType,
  onSuccess,
  onCancel,
}: TicketTypeFormProps) {
  const isEdit = !!ticketType;
  const createMutation = useCreateTicketType();
  const updateMutation = useUpdateTicketType(ticketType?.id ?? "");

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<TicketTypeInput>({
    resolver: zodResolver(ticketTypeSchema),
    defaultValues: {
      name: ticketType?.name ?? "",
      description: ticketType?.description ?? "",
      price: ticketType?.price ?? 0,
      currency: ticketType?.currency ?? "VND",
      quantity: ticketType?.quantity ?? 100,
      max_per_purchase: ticketType?.max_per_purchase ?? 10,
      min_per_purchase: ticketType?.min_per_purchase ?? 1,
    },
  });

  const isPending = createMutation.isPending || updateMutation.isPending;

  const onSubmit = async (data: TicketTypeInput) => {
    try {
      if (isEdit) {
        await updateMutation.mutateAsync({
          ...data,
          event_id: eventId,
        });
        showToast.success("Ticket type updated successfully.");
      } else {
        await createMutation.mutateAsync({
          ...data,
          event_id: eventId,
        });
        showToast.success("Ticket type created successfully.");
      }
      onSuccess();
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="space-y-4 rounded-lg border bg-card p-4"
    >
      <h4 className="font-semibold">
        {isEdit ? "Edit Ticket Type" : "Add Ticket Type"}
      </h4>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="tt-name">Name *</Label>
          <Input
            id="tt-name"
            placeholder="e.g. VIP, Standard, Early Bird"
            {...register("name")}
          />
          {errors.name && (
            <p className="text-xs text-destructive">{errors.name.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="tt-price">Price *</Label>
          <Input
            id="tt-price"
            type="number"
            min={0}
            step="1000"
            {...register("price", { valueAsNumber: true })}
          />
          {errors.price && (
            <p className="text-xs text-destructive">{errors.price.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="tt-quantity">Total Quantity *</Label>
          <Input
            id="tt-quantity"
            type="number"
            min={1}
            {...register("quantity", { valueAsNumber: true })}
          />
          {errors.quantity && (
            <p className="text-xs text-destructive">
              {errors.quantity.message}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="tt-currency">Currency</Label>
          <Input id="tt-currency" placeholder="VND" {...register("currency")} />
        </div>

        <div className="space-y-2">
          <Label htmlFor="tt-max">Max per Purchase</Label>
          <Input
            id="tt-max"
            type="number"
            min={1}
            {...register("max_per_purchase", { valueAsNumber: true })}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="tt-min">Min per Purchase</Label>
          <Input
            id="tt-min"
            type="number"
            min={1}
            {...register("min_per_purchase", { valueAsNumber: true })}
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="tt-desc">Description</Label>
        <Textarea
          id="tt-desc"
          placeholder="Optional description..."
          rows={2}
          {...register("description")}
        />
      </div>

      <div className="flex items-center justify-end gap-2">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isPending}>
          {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          {isEdit ? "Save Changes" : "Add Ticket Type"}
        </Button>
      </div>
    </form>
  );
}
