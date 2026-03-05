"use client";

import {
  LayoutDashboard,
  Users,
  CalendarDays,
  CreditCard,
  Shield,
  Settings,
} from "lucide-react";
import { AdminSidebar } from "@/components/organisms/layout/admin-sidebar";
import { Topbar } from "@/components/organisms/layout/topbar";

const adminNavItems = [
  { href: "/admin/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/admin/users", label: "Users", icon: Users },
  { href: "/admin/events", label: "Events", icon: CalendarDays },
  { href: "/admin/payments", label: "Payments", icon: CreditCard },
  { href: "/admin/security", label: "Security", icon: Shield },
  { href: "/admin/settings", label: "Settings", icon: Settings },
];

interface AdminLayoutProps {
  children: React.ReactNode;
}

export function AdminLayout({ children }: AdminLayoutProps) {
  return (
    <div className="flex h-screen overflow-hidden">
      <AdminSidebar />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Topbar navItems={adminNavItems} brandLabel="Admin Panel" brandHref="/admin/dashboard" />
        <main className="flex-1 overflow-y-auto p-6">{children}</main>
      </div>
    </div>
  );
}
