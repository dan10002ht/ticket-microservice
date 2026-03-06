"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { ArrowLeft, ChevronDown, ChevronUp } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/molecules/page-header";
import { StepIndicator } from "@/components/molecules/step-indicator";
import { useCreateEvent } from "@/lib/api/queries";
import type { ApiError } from "@/lib/api/types/common";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  createEventSchema,
  type CreateEventInput,
} from "@/lib/validators/event";

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
import { ZonePricingStep } from "@/components/organisms/org-events/zone-pricing-step";
import { EventReviewStep } from "@/components/organisms/org-events/event-review-step";

const STEPS = ["Basic Info", "Zones & Pricing", "Review"];

export default function CreateEventPage() {
  const router = useRouter();
  const createMutation = useCreateEvent();
  const [step, setStep] = useState(1);
  const [eventId, setEventId] = useState<string | null>(null);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const {
    register,
    handleSubmit,
    trigger,
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
      venue_country: "",
      venue_capacity: undefined,
      event_type: "",
      category: "",
      sale_start_date: "",
      sale_end_date: "",
      min_age: "",
      is_featured: false,
    },
  });

  const onSubmitStep1 = async (data: CreateEventInput) => {
    try {
      const event = await createMutation.mutateAsync({
        name: data.name,
        description: data.description || undefined,
        start_date: new Date(data.start_date).toISOString(),
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
      setEventId(event.id);
      setStep(2);
      toast.success("Event created as draft!");
    } catch (err: unknown) {
      const message =
        (err as ApiError)?.error?.message || "Failed to create event.";
      toast.error(message);
    }
  };

  const handleNext = async () => {
    if (step === 1) {
      const valid = await trigger();
      if (valid) {
        handleSubmit(onSubmitStep1)();
      }
    } else if (step === 2) {
      setStep(3);
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

      <StepIndicator steps={STEPS} currentStep={step} className="mt-6" />

      {step === 1 && (
        <Card className="mt-6">
          <CardHeader>
            <CardTitle>Event Details</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit(onSubmitStep1)} className="space-y-6">
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
                  <Label htmlFor="end_date">End date *</Label>
                  <Input
                    id="end_date"
                    type="datetime-local"
                    {...register("end_date")}
                  />
                  {errors.end_date && (
                    <p className="text-xs text-destructive">
                      {errors.end_date.message}
                    </p>
                  )}
                </div>
              </div>

              {/* Venue fields */}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="venue_name">Venue name *</Label>
                  <Input
                    id="venue_name"
                    placeholder="e.g. My Dinh Stadium"
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
                    placeholder="e.g. Ha Noi"
                    {...register("venue_city")}
                  />
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <Label htmlFor="venue_address">Address</Label>
                  <Input
                    id="venue_address"
                    placeholder="Full address"
                    {...register("venue_address")}
                  />
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
                            <span className="text-sm font-medium">
                              Featured Event
                            </span>
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

              {/* Submit */}
              <div className="flex justify-end gap-3">
                <Button type="button" variant="outline" asChild>
                  <Link href="/org/events">Cancel</Link>
                </Button>
                <Button type="submit" disabled={createMutation.isPending}>
                  {createMutation.isPending ? "Creating..." : "Next"}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {step === 2 && eventId && (
        <ZonePricingStep
          eventId={eventId}
          onNext={() => setStep(3)}
          onSkip={() => setStep(3)}
          className="mt-6"
        />
      )}

      {step === 3 && eventId && (
        <EventReviewStep
          eventId={eventId}
          onDone={() => router.push("/org/events")}
          className="mt-6"
        />
      )}
    </>
  );
}
