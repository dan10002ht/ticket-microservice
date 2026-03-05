"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  CalendarDays,
  MapPin,
  Users,
  Pencil,
  Trash2,
  Send,
  Plus,
  Loader2,
  ArrowLeft,
  ScanLine,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { TicketTypeForm } from "@/components/molecules/ticket-type-form";
import { TicketTypeListItem } from "@/components/molecules/ticket-type-list-item";
import {
  usePublishEvent,
  useDeleteEvent,
  useDeleteTicketType,
} from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import type { Event } from "@/lib/api/types/event";
import type { TicketType } from "@/lib/api/types/ticket";
import type { ApiError } from "@/lib/api/types/common";

interface EventDetailOrgContentProps {
  event?: Event;
  ticketTypes?: TicketType[];
  isLoading: boolean;
  error: ApiError | null;
}

const statusVariant: Record<
  string,
  "default" | "secondary" | "destructive" | "outline"
> = {
  published: "default",
  draft: "secondary",
  cancelled: "destructive",
  completed: "outline",
};

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function EventDetailOrgContent({
  event,
  ticketTypes,
  isLoading,
  error,
}: EventDetailOrgContentProps) {
  const router = useRouter();
  const publishMutation = usePublishEvent();
  const deleteMutation = useDeleteEvent();
  const deleteTicketTypeMutation = useDeleteTicketType();

  const [showTicketForm, setShowTicketForm] = useState(false);
  const [editingTicketType, setEditingTicketType] = useState<TicketType | null>(
    null
  );
  const [confirmDelete, setConfirmDelete] = useState(false);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-48 w-full rounded-lg" />
        <Skeleton className="h-32 w-full rounded-lg" />
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

  const handlePublish = async () => {
    try {
      await publishMutation.mutateAsync(event.id);
      showToast.success("Event published successfully!");
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  const handleDelete = async () => {
    if (!confirmDelete) {
      setConfirmDelete(true);
      return;
    }
    try {
      await deleteMutation.mutateAsync(event.id);
      showToast.success("Event deleted.");
      router.push("/org/events");
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  const handleDeleteTicketType = async (typeId: string) => {
    try {
      await deleteTicketTypeMutation.mutateAsync(typeId);
      showToast.success("Ticket type deleted.");
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  const isDraft = !event.created_at || event.description !== undefined; // fallback: always show publish if no status field
  // Event type doesn't have explicit status field in the current type definition
  // We'll show publish button always and let the API handle the validation

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/org/events">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div>
            <h1 className="text-2xl font-bold">{event.name}</h1>
            <p className="text-sm text-muted-foreground">
              Created {formatDate(event.created_at)}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" asChild>
            <Link href={`/org/events/${event.id}/edit`}>
              <Pencil className="mr-1.5 h-3.5 w-3.5" />
              Edit
            </Link>
          </Button>
          <Button variant="outline" size="sm" asChild>
            <Link href={`/org/events/${event.id}/checkin`}>
              <ScanLine className="mr-1.5 h-3.5 w-3.5" />
              Check-in
            </Link>
          </Button>
          <Button
            variant="default"
            size="sm"
            onClick={handlePublish}
            disabled={publishMutation.isPending}
          >
            {publishMutation.isPending ? (
              <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />
            ) : (
              <Send className="mr-1.5 h-3.5 w-3.5" />
            )}
            Publish
          </Button>
          <Button
            variant={confirmDelete ? "destructive" : "outline"}
            size="sm"
            onClick={handleDelete}
            onBlur={() => setConfirmDelete(false)}
            disabled={deleteMutation.isPending}
          >
            {deleteMutation.isPending ? (
              <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />
            ) : (
              <Trash2 className="mr-1.5 h-3.5 w-3.5" />
            )}
            {confirmDelete ? "Confirm Delete" : "Delete"}
          </Button>
        </div>
      </div>

      {/* Event Info */}
      <div className="rounded-lg border bg-card p-4 space-y-3">
        <h3 className="font-semibold">Event Details</h3>
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="flex items-start gap-3">
            <CalendarDays className="mt-0.5 h-4 w-4 text-muted-foreground" />
            <div>
              <p className="text-sm text-muted-foreground">Date</p>
              <p className="font-medium">{formatDate(event.start_date)}</p>
              {event.end_date && (
                <p className="text-sm text-muted-foreground">
                  to {formatDate(event.end_date)}
                </p>
              )}
            </div>
          </div>

          <div className="flex items-start gap-3">
            <MapPin className="mt-0.5 h-4 w-4 text-muted-foreground" />
            <div>
              <p className="text-sm text-muted-foreground">Venue</p>
              <p className="font-medium">{event.venue_name}</p>
              {(event.venue_address || event.venue_city) && (
                <p className="text-sm text-muted-foreground">
                  {[event.venue_address, event.venue_city]
                    .filter(Boolean)
                    .join(", ")}
                </p>
              )}
            </div>
          </div>

          {event.venue_capacity && (
            <div className="flex items-start gap-3">
              <Users className="mt-0.5 h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Capacity</p>
                <p className="font-medium">
                  {event.venue_capacity.toLocaleString()}
                </p>
              </div>
            </div>
          )}
        </div>

        {event.description && (
          <>
            <Separator />
            <div>
              <p className="text-sm text-muted-foreground">Description</p>
              <p className="mt-1 text-sm whitespace-pre-line">
                {event.description}
              </p>
            </div>
          </>
        )}
      </div>

      {/* Ticket Types */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold">Ticket Types</h3>
          {!showTicketForm && !editingTicketType && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowTicketForm(true)}
            >
              <Plus className="mr-1.5 h-3.5 w-3.5" />
              Add Ticket Type
            </Button>
          )}
        </div>

        {/* Inline form for adding */}
        {showTicketForm && (
          <TicketTypeForm
            eventId={event.id}
            onSuccess={() => setShowTicketForm(false)}
            onCancel={() => setShowTicketForm(false)}
          />
        )}

        {/* Inline form for editing */}
        {editingTicketType && (
          <TicketTypeForm
            eventId={event.id}
            ticketType={editingTicketType}
            onSuccess={() => setEditingTicketType(null)}
            onCancel={() => setEditingTicketType(null)}
          />
        )}

        {/* Ticket type list */}
        {ticketTypes && ticketTypes.length > 0 ? (
          <div className="space-y-2">
            {ticketTypes.map((tt) => (
              <TicketTypeListItem
                key={tt.id}
                ticketType={tt}
                onEdit={() => {
                  setShowTicketForm(false);
                  setEditingTicketType(tt);
                }}
                onDelete={() => handleDeleteTicketType(tt.id)}
                isDeleting={deleteTicketTypeMutation.isPending}
              />
            ))}
          </div>
        ) : (
          !showTicketForm && (
            <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
              No ticket types yet. Add one to allow bookings.
            </div>
          )
        )}
      </div>
    </div>
  );
}
