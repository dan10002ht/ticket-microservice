"use client";

import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/molecules/page-header";
import { StatusBadge } from "@/components/molecules/status-badge";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";

interface OrgEvent {
  id: string;
  title: string;
  date: string;
  venue: string;
  capacity: number;
  sold: number;
  status: string;
}

const events: OrgEvent[] = [
  { id: "1", title: "Summer Music Festival 2026", date: "Jul 15, 2026", venue: "Mỹ Đình Stadium", capacity: 5000, sold: 3800, status: "published" },
  { id: "2", title: "Tech Conference Vietnam", date: "Aug 20, 2026", venue: "GEM Center", capacity: 800, sold: 450, status: "published" },
  { id: "3", title: "Stand-up Comedy Night", date: "Jun 10, 2026", venue: "Nhà hát Lớn", capacity: 300, sold: 255, status: "published" },
  { id: "4", title: "New Year Countdown 2027", date: "Dec 31, 2026", venue: "Phú Mỹ Hưng", capacity: 10000, sold: 0, status: "draft" },
];

const columns: Column<OrgEvent>[] = [
  { key: "title", header: "Event" },
  { key: "date", header: "Date" },
  { key: "venue", header: "Venue" },
  {
    key: "sold",
    header: "Sold / Capacity",
    render: (e) => (
      <span className="tabular-nums">
        {e.sold} / {e.capacity}
      </span>
    ),
  },
  {
    key: "status",
    header: "Status",
    render: (e) => <StatusBadge status={e.status} />,
  },
];

export default function OrgEventsPage() {
  return (
    <>
      <PageHeader
        title="My Events"
        description="Create and manage your events"
      >
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Create Event
        </Button>
      </PageHeader>

      <div className="mt-6">
        <DataTable
          columns={columns}
          data={events}
          keyExtractor={(e) => e.id}
        />
      </div>
    </>
  );
}
