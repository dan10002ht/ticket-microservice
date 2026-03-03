import type { Metadata } from "next";
import {
  Activity,
  CalendarDays,
  DollarSign,
  Users,
} from "lucide-react";
import { PageHeader } from "@/components/molecules/page-header";
import { StatsCard } from "@/components/molecules/stats-card";

export const metadata: Metadata = {
  title: "Admin Dashboard",
  description: "System overview and statistics.",
};

export default function AdminDashboardPage() {
  return (
    <>
      <PageHeader
        title="Admin Dashboard"
        description="System overview and statistics"
      />

      <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          title="Total Users"
          value="2,543"
          icon={Users}
          trend={{ value: 15, isPositive: true }}
        />
        <StatsCard
          title="Active Events"
          value={48}
          icon={CalendarDays}
          trend={{ value: 8, isPositive: true }}
        />
        <StatsCard
          title="Revenue (MTD)"
          value="125.8M VND"
          icon={DollarSign}
          trend={{ value: 22, isPositive: true }}
        />
        <StatsCard
          title="System Health"
          value="99.9%"
          icon={Activity}
          description="All services operational"
        />
      </div>

      <div className="mt-8 grid gap-6 lg:grid-cols-2">
        <div className="rounded-lg border bg-card p-6">
          <h3 className="font-semibold">Revenue Chart</h3>
          <div className="mt-4 flex h-48 items-center justify-center text-sm text-muted-foreground">
            Chart placeholder — will integrate with recharts
          </div>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <h3 className="font-semibold">Recent Activity</h3>
          <div className="mt-4 flex h-48 items-center justify-center text-sm text-muted-foreground">
            Activity feed placeholder
          </div>
        </div>
      </div>
    </>
  );
}
