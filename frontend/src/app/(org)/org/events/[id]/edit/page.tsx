"use client";

import { useParams } from "next/navigation";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { EventEditForm } from "@/components/organisms/org-events/event-edit-form";
import { ZonePricingStep } from "@/components/organisms/org-events/zone-pricing-step";
import { SeatManagementStep } from "@/components/organisms/org-events/seat-management-step";
import { useEvent } from "@/lib/api/queries";

export default function OrgEventEditPage() {
  const params = useParams();
  const eventId = params.id as string;
  const { data: event, isLoading, error } = useEvent(eventId);

  return (
    <Tabs defaultValue="basic" className="space-y-6">
      <EventEditForm event={event} isLoading={isLoading} error={error}>
        <TabsList>
          <TabsTrigger value="basic">Basic Info</TabsTrigger>
          <TabsTrigger value="zones" disabled={!event}>
            Zones & Pricing
          </TabsTrigger>
          <TabsTrigger value="seats" disabled={!event}>
            Seats
          </TabsTrigger>
        </TabsList>
      </EventEditForm>

      <TabsContent value="zones">
        {event && (
          <ZonePricingStep
            eventId={event.id}
            onNext={() => {}}
            onSkip={() => {}}
          />
        )}
      </TabsContent>

      <TabsContent value="seats">
        {event && <SeatManagementStep eventId={event.id} />}
      </TabsContent>
    </Tabs>
  );
}
