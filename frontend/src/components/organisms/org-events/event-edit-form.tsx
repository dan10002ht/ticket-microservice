"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, ChevronDown, ChevronUp, Loader2 } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Skeleton } from "@/components/ui/skeleton";
import { TabsContent } from "@/components/ui/tabs";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { StatusBadge } from "@/components/molecules/status-badge";
import { useUpdateEvent } from "@/lib/api/queries";
import { createEventSchema, type CreateEventInput } from "@/lib/validators/event";
import { showToast } from "@/lib/toast";
import type { Event } from "@/lib/api/types/event";
import type { ApiError } from "@/lib/api/types/common";

const EVENT_TYPES = [
  { value: "concert", label: "Concert" },
  { value: "theater", label: "Theater" },
  { value: "sports", label: "Sports" },
  { value: "conference", label: "Conference" },
  { value: "festival", label: "Festival" },
  { value: "other", label: "Other" },
];

const CATEGORIES = [
  { value: "music", label: "Music" },
  { value: "arts", label: "Arts" },
  { value: "sports", label: "Sports" },
  { value: "business", label: "Business" },
  { value: "education", label: "Education" },
  { value: "entertainment", label: "Entertainment" },
  { value: "other", label: "Other" },
];

interface EventEditFormProps {
  event?: Event;
  isLoading: boolean;
  error: ApiError | null;
  children?: React.ReactNode;
}

function toDatetimeLocal(iso: string | undefined): string {
  if (!iso) return "";
  const d = new Date(iso);
  if (isNaN(d.getTime())) return "";
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function EventEditForm({ event, isLoading, error, children }: EventEditFormProps) {
  const router = useRouter();
  const updateMutation = useUpdateEvent(event?.id ?? "");
  const [showAdvanced, setShowAdvanced] = useState(false);

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
          venue_country: event.venue_country ?? "",
          venue_capacity: event.venue_capacity?.toString() ?? "",
          event_type: event.event_type ?? "",
          category: event.category ?? "",
          sale_start_date: toDatetimeLocal(event.sale_start_date),
          sale_end_date: toDatetimeLocal(event.sale_end_date),
          min_age: event.min_age?.toString() ?? "",
          is_featured: event.is_featured ?? false,
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
        venue_country: data.venue_country || undefined,
        venue_capacity: data.venue_capacity
          ? parseInt(data.venue_capacity, 10)
          : undefined,
        event_type: data.event_type || undefined,
        category: data.category || undefined,
        sale_start_date: data.sale_start_date
          ? new Date(data.sale_start_date).toISOString()
          : undefined,
        sale_end_date: data.sale_end_date
          ? new Date(data.sale_end_date).toISOString()
          : undefined,
        min_age: data.min_age ? parseInt(data.min_age, 10) : undefined,
        is_featured: data.is_featured || undefined,
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
        <StatusBadge status={event.status ?? "draft"} />
      </div>

      {children}

      <TabsContent value="basic">
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
              <Label htmlFor="venue_country">Country</Label>
              <Input
                id="venue_country"
                placeholder="e.g. Vietnam"
                {...register("venue_country")}
              />
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

          {/* Advanced Options */}
          <div>
            <Button
              type="button"
              variant="ghost"
              className="gap-2 px-0 text-muted-foreground"
              onClick={() => setShowAdvanced(!showAdvanced)}
            >
              {showAdvanced ? (
                <ChevronUp className="h-4 w-4" />
              ) : (
                <ChevronDown className="h-4 w-4" />
              )}
              Advanced Options
            </Button>

            {showAdvanced && (
              <div className="mt-4 rounded-lg border bg-muted/30 p-5">
                <div className="grid gap-x-6 gap-y-5 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="event_type">Event Type</Label>
                    <Select
                      defaultValue={event.event_type ?? ""}
                      onValueChange={(val) =>
                        register("event_type").onChange({
                          target: { name: "event_type", value: val },
                        })
                      }
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select type" />
                      </SelectTrigger>
                      <SelectContent>
                        {EVENT_TYPES.map((t) => (
                          <SelectItem key={t.value} value={t.value}>
                            {t.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="category">Category</Label>
                    <Select
                      defaultValue={event.category ?? ""}
                      onValueChange={(val) =>
                        register("category").onChange({
                          target: { name: "category", value: val },
                        })
                      }
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select category" />
                      </SelectTrigger>
                      <SelectContent>
                        {CATEGORIES.map((c) => (
                          <SelectItem key={c.value} value={c.value}>
                            {c.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="sale_start_date">Sale Start Date</Label>
                    <Input
                      id="sale_start_date"
                      type="datetime-local"
                      {...register("sale_start_date")}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="sale_end_date">Sale End Date</Label>
                    <Input
                      id="sale_end_date"
                      type="datetime-local"
                      {...register("sale_end_date")}
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="min_age">Minimum Age</Label>
                    <Input
                      id="min_age"
                      type="number"
                      min={0}
                      placeholder="e.g. 18"
                      {...register("min_age")}
                    />
                  </div>

                  <div className="sm:col-span-2">
                    <label
                      htmlFor="is_featured"
                      className="flex cursor-pointer items-center gap-3 rounded-md border px-4 py-3 transition-colors hover:bg-accent/50"
                    >
                      <input
                        type="checkbox"
                        id="is_featured"
                        className="h-4 w-4 rounded border-input"
                        {...register("is_featured")}
                      />
                      <div>
                        <span className="text-sm font-medium">Featured Event</span>
                        <p className="text-xs text-muted-foreground">
                          Highlight this event on the homepage
                        </p>
                      </div>
                    </label>
                  </div>
                </div>
              </div>
            )}
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
      </TabsContent>
    </div>
  );
}
