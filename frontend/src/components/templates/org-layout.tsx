"use client";

import {
  LayoutDashboard,
  CalendarDays,
  BarChart3,
  Settings,
} from "lucide-react";
import { OrgSidebar } from "@/components/organisms/layout/org-sidebar";
import { Topbar } from "@/components/organisms/layout/topbar";

const orgNavItems = [
  { href: "/org/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/org/events", label: "Events", icon: CalendarDays },
  { href: "/org/analytics", label: "Analytics", icon: BarChart3 },
  { href: "/org/settings", label: "Settings", icon: Settings },
];

interface OrgLayoutProps {
  children: React.ReactNode;
}

export function OrgLayout({ children }: OrgLayoutProps) {
  return (
    <div className="flex h-screen overflow-hidden">
      <OrgSidebar />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Topbar navItems={orgNavItems} brandLabel="Organizer" brandHref="/org/dashboard" />
        <main className="flex-1 overflow-y-auto p-6">{children}</main>
      </div>
    </div>
  );
}
