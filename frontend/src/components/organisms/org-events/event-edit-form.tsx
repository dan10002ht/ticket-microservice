"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, Loader2 } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Skeleton } from "@/components/ui/skeleton";
import { useUpdateEvent } from "@/lib/api/queries";
import { createEventSchema, type CreateEventInput } from "@/lib/validators/event";
import { showToast } from "@/lib/toast";
import type { Event } from "@/lib/api/types/event";
import type { ApiError } from "@/lib/api/types/common";

interface EventEditFormProps {
  event?: Event;
  isLoading: boolean;
  error: ApiError | null;
}

function toDatetimeLocal(iso: string | undefined): string {
  if (!iso) return "";
  const d = new Date(iso);
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function EventEditForm({ event, isLoading, error }: EventEditFormProps) {
  const router = useRouter();
  const updateMutation = useUpdateEvent(event?.id ?? "");

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateEventInput>({
    resolver: zodResolver(createEventSchema),
    values: event
      ? {
          name: event.name,
          description: event.description ?? "",
          start_date: toDatetimeLocal(event.start_date),
          end_date: toDatetimeLocal(event.end_date),
          venue_name: event.venue_name,
          venue_address: event.venue_address ?? "",
          venue_city: event.venue_city ?? "",
          venue_capacity: event.venue_capacity?.toString() ?? "",
        }
      : undefined,
  });

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full rounded-lg" />
      </div>
    );
  }

  if (error || !event) {
    return (
      <div className="rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
        {error?.error?.message || "Event not found."}
      </div>
    );
  }

  const onSubmit = async (data: CreateEventInput) => {
    try {
      await updateMutation.mutateAsync({
        name: data.name,
        description: data.description || undefined,
        start_date: data.start_date
          ? new Date(data.start_date).toISOString()
          : undefined,
        end_date: data.end_date
          ? new Date(data.end_date).toISOString()
          : undefined,
        venue_name: data.venue_name,
        venue_address: data.venue_address || undefined,
        venue_city: data.venue_city || undefined,
        venue_capacity: data.venue_capacity
          ? parseInt(data.venue_capacity, 10)
          : undefined,
      });
      showToast.success("Event updated successfully.");
      router.push(`/org/events/${event.id}`);
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" asChild>
          <Link href={`/org/events/${event.id}`}>
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold">Edit Event</h1>
      </div>

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="space-y-4 rounded-lg border bg-card p-6"
      >
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2 sm:col-span-2">
            <Label htmlFor="name">Event Name *</Label>
            <Input id="name" {...register("name")} />
            {errors.name && (
              <p className="text-xs text-destructive">{errors.name.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="start_date">Start Date *</Label>
            <Input
              id="start_date"
              type="datetime-local"
              {...register("start_date")}
            />
            {errors.start_date && (
              <p className="text-xs text-destructive">
                {errors.start_date.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="end_date">End Date</Label>
            <Input
              id="end_date"
              type="datetime-local"
              {...register("end_date")}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="venue_name">Venue Name *</Label>
            <Input id="venue_name" {...register("venue_name")} />
            {errors.venue_name && (
              <p className="text-xs text-destructive">
                {errors.venue_name.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="venue_city">City</Label>
            <Input id="venue_city" {...register("venue_city")} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="venue_address">Address</Label>
            <Input id="venue_address" {...register("venue_address")} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="venue_capacity">Capacity</Label>
            <Input
              id="venue_capacity"
              type="number"
              min={0}
              {...register("venue_capacity")}
            />
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="description">Description</Label>
          <Textarea
            id="description"
            rows={4}
            {...register("description")}
          />
        </div>

        <div className="flex items-center justify-end gap-3 pt-2">
          <Button variant="outline" type="button" asChild>
            <Link href={`/org/events/${event.id}`}>Cancel</Link>
          </Button>
          <Button type="submit" disabled={updateMutation.isPending}>
            {updateMutation.isPending && (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            )}
            Save Changes
          </Button>
        </div>
      </form>
    </div>
  );
}
