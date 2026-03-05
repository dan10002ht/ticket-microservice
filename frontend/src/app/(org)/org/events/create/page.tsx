"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/molecules/page-header";
import { useCreateEvent } from "@/lib/api/queries";
import type { ApiError } from "@/lib/api/types/common";
import {
  createEventSchema,
  type CreateEventInput,
} from "@/lib/validators/event";

export default function CreateEventPage() {
  const router = useRouter();
  const createMutation = useCreateEvent();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateEventInput>({
    resolver: zodResolver(createEventSchema),
    defaultValues: {
      name: "",
      description: "",
      start_date: "",
      end_date: "",
      venue_name: "",
      venue_address: "",
      venue_city: "",
      venue_capacity: undefined,
    },
  });

  const onSubmit = async (data: CreateEventInput) => {
    try {
      await createMutation.mutateAsync({
        name: data.name,
        description: data.description || undefined,
        start_date: new Date(data.start_date).toISOString(),
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
      toast.success("Event created successfully!");
      router.push("/org/events");
    } catch (err: unknown) {
      const message =
        (err as ApiError)?.error?.message || "Failed to create event.";
      toast.error(message);
    }
  };

  return (
    <>
      <PageHeader title="Create Event" description="Set up a new event">
        <Button variant="ghost" asChild>
          <Link href="/org/events">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Events
          </Link>
        </Button>
      </PageHeader>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Event Details</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Event name */}
            <div className="space-y-2">
              <Label htmlFor="name">Event name *</Label>
              <Input
                id="name"
                placeholder="e.g. Summer Music Festival 2026"
                {...register("name")}
              />
              {errors.name && (
                <p className="text-xs text-destructive">
                  {errors.name.message}
                </p>
              )}
            </div>

            {/* Description */}
            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                placeholder="Tell attendees about your event..."
                rows={4}
                {...register("description")}
              />
            </div>

            {/* Date fields */}
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="start_date">Start date *</Label>
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
                <Label htmlFor="end_date">End date</Label>
                <Input
                  id="end_date"
                  type="datetime-local"
                  {...register("end_date")}
                />
              </div>
            </div>

            {/* Venue fields */}
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="venue_name">Venue name *</Label>
                <Input
                  id="venue_name"
                  placeholder="e.g. Mỹ Đình Stadium"
                  {...register("venue_name")}
                />
                {errors.venue_name && (
                  <p className="text-xs text-destructive">
                    {errors.venue_name.message}
                  </p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="venue_city">City</Label>
                <Input
                  id="venue_city"
                  placeholder="e.g. Hà Nội"
                  {...register("venue_city")}
                />
              </div>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="venue_address">Address</Label>
                <Input
                  id="venue_address"
                  placeholder="Full address"
                  {...register("venue_address")}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="venue_capacity">Capacity</Label>
                <Input
                  id="venue_capacity"
                  type="number"
                  placeholder="e.g. 5000"
                  {...register("venue_capacity")}
                />
                {errors.venue_capacity && (
                  <p className="text-xs text-destructive">
                    {errors.venue_capacity.message}
                  </p>
                )}
              </div>
            </div>

            {/* Submit */}
            <div className="flex justify-end gap-3">
              <Button type="button" variant="outline" asChild>
                <Link href="/org/events">Cancel</Link>
              </Button>
              <Button type="submit" disabled={createMutation.isPending}>
                {createMutation.isPending
                  ? "Creating..."
                  : "Create Event"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </>
  );
}
