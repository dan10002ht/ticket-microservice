"use client";

import { useParams } from "next/navigation";
import { BookingDetailContent } from "@/components/organisms/booking/booking-detail-content";
import { useBooking, usePayment } from "@/lib/api/queries";

export default function BookingDetailPage() {
  const params = useParams();
  const bookingId = params.id as string;
  const { data: booking, isLoading, error } = useBooking(bookingId);
  const { data: payment } = usePayment(booking?.payment_reference ?? "");

  return (
    <section className="container mx-auto max-w-2xl px-4 py-8">
      <BookingDetailContent
        booking={booking}
        payment={payment}
        isLoading={isLoading}
        error={error}
      />
    </section>
  );
}
