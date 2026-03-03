import type { Metadata } from "next";
import { PageHeader } from "@/components/molecules/page-header";
import { StatusBadge } from "@/components/molecules/status-badge";
import { PriceDisplay } from "@/components/molecules/price-display";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import type { Booking } from "@/types";

export const metadata: Metadata = {
  title: "My Bookings",
  description: "View and manage your event bookings.",
};

const sampleBookings: Booking[] = [
  {
    id: "BK-001",
    userId: "user-1",
    eventId: "1",
    status: "confirmed",
    totalAmount: 1000000,
    ticketCount: 2,
    createdAt: "2026-03-01T10:00:00Z",
    event: {
      id: "1",
      title: "Summer Music Festival 2026",
      description: "",
      venue: "Sân vận động Mỹ Đình",
      address: "Hà Nội",
      startDate: "2026-07-15T18:00:00Z",
      endDate: "2026-07-15T23:00:00Z",
      status: "published",
      organizerId: "org-1",
      totalCapacity: 5000,
      availableCapacity: 1200,
      minPrice: 500000,
      maxPrice: 2000000,
      createdAt: "2026-01-01T00:00:00Z",
    },
  },
  {
    id: "BK-002",
    userId: "user-1",
    eventId: "2",
    status: "pending",
    totalAmount: 600000,
    ticketCount: 1,
    createdAt: "2026-03-02T14:00:00Z",
    event: {
      id: "2",
      title: "Tech Conference Vietnam",
      description: "",
      venue: "GEM Center",
      address: "TP. Hồ Chí Minh",
      startDate: "2026-08-20T09:00:00Z",
      endDate: "2026-08-21T17:00:00Z",
      status: "published",
      organizerId: "org-2",
      totalCapacity: 800,
      availableCapacity: 350,
      minPrice: 300000,
      maxPrice: 1500000,
      createdAt: "2026-02-01T00:00:00Z",
    },
  },
  {
    id: "BK-003",
    userId: "user-1",
    eventId: "3",
    status: "cancelled",
    totalAmount: 400000,
    ticketCount: 1,
    createdAt: "2026-02-20T09:00:00Z",
    event: {
      id: "3",
      title: "Stand-up Comedy Night",
      description: "",
      venue: "Nhà hát Lớn",
      address: "Hà Nội",
      startDate: "2026-06-10T19:30:00Z",
      endDate: "2026-06-10T22:00:00Z",
      status: "published",
      organizerId: "org-3",
      totalCapacity: 300,
      availableCapacity: 45,
      minPrice: 200000,
      maxPrice: 800000,
      createdAt: "2026-03-01T00:00:00Z",
    },
  },
];

const columns: Column<Booking>[] = [
  { key: "id", header: "Booking ID" },
  {
    key: "event",
    header: "Event",
    render: (booking) => booking.event?.title ?? "-",
  },
  {
    key: "ticketCount",
    header: "Tickets",
    className: "text-center",
  },
  {
    key: "totalAmount",
    header: "Total",
    render: (booking) => <PriceDisplay amount={booking.totalAmount} size="sm" />,
  },
  {
    key: "status",
    header: "Status",
    render: (booking) => <StatusBadge status={booking.status} />,
  },
  {
    key: "createdAt",
    header: "Date",
    render: (booking) =>
      new Date(booking.createdAt).toLocaleDateString("vi-VN"),
  },
];

export default function MyBookingsPage() {
  return (
    <>
      <PageHeader
        title="My Bookings"
        description="View and manage your event bookings"
      />

      <div className="mt-6">
        <DataTable
          columns={columns}
          data={sampleBookings}
          keyExtractor={(b) => b.id}
        />
      </div>
    </>
  );
}
