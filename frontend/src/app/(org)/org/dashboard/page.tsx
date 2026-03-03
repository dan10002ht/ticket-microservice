import type { Metadata } from "next";
import {
  CalendarDays,
  DollarSign,
  Ticket,
  Users,
} from "lucide-react";
import { PageHeader } from "@/components/molecules/page-header";
import { StatsCard } from "@/components/molecules/stats-card";
import { StatusBadge } from "@/components/molecules/status-badge";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";

export const metadata: Metadata = {
  title: "Organizer Dashboard",
  description: "Overview of your events and sales.",
};

interface RecentBooking {
  id: string;
  customer: string;
  event: string;
  tickets: number;
  total: string;
  status: string;
}

const recentBookings: RecentBooking[] = [
  { id: "BK-101", customer: "Nguyen Van A", event: "Summer Music Festival", tickets: 3, total: "1,500,000 VND", status: "confirmed" },
  { id: "BK-102", customer: "Tran Thi B", event: "Tech Conference", tickets: 1, total: "300,000 VND", status: "pending" },
  { id: "BK-103", customer: "Le Van C", event: "Summer Music Festival", tickets: 2, total: "1,000,000 VND", status: "confirmed" },
  { id: "BK-104", customer: "Pham Thi D", event: "Comedy Night", tickets: 4, total: "3,200,000 VND", status: "confirmed" },
  { id: "BK-105", customer: "Hoang Van E", event: "Tech Conference", tickets: 1, total: "1,500,000 VND", status: "cancelled" },
];

const columns: Column<RecentBooking>[] = [
  { key: "id", header: "ID" },
  { key: "customer", header: "Customer" },
  { key: "event", header: "Event" },
  { key: "tickets", header: "Tickets", className: "text-center" },
  { key: "total", header: "Total" },
  {
    key: "status",
    header: "Status",
    render: (booking) => <StatusBadge status={booking.status} />,
  },
];

export default function OrgDashboardPage() {
  return (
    <>
      <PageHeader
        title="Dashboard"
        description="Overview of your events and sales"
      />

      <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          title="Total Events"
          value={12}
          icon={CalendarDays}
          trend={{ value: 8, isPositive: true }}
        />
        <StatsCard
          title="Tickets Sold"
          value="1,234"
          icon={Ticket}
          trend={{ value: 12, isPositive: true }}
        />
        <StatsCard
          title="Total Revenue"
          value="45.2M VND"
          icon={DollarSign}
          trend={{ value: 5, isPositive: true }}
        />
        <StatsCard
          title="Attendees"
          value="856"
          icon={Users}
          trend={{ value: 3, isPositive: false }}
        />
      </div>

      <div className="mt-8">
        <h2 className="mb-4 text-lg font-semibold">Recent Bookings</h2>
        <DataTable
          columns={columns}
          data={recentBookings}
          keyExtractor={(b) => b.id}
        />
      </div>
    </>
  );
}
