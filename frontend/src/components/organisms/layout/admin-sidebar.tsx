"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  Users,
  CalendarDays,
  CreditCard,
  Shield,
  Settings,
  Ticket,
} from "lucide-react";
import { NavItem } from "@/components/molecules/nav-item";
import { cn } from "@/lib/utils";

const adminNavItems = [
  { href: "/admin/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/admin/users", label: "Users", icon: Users },
  { href: "/admin/events", label: "Events", icon: CalendarDays },
  { href: "/admin/payments", label: "Payments", icon: CreditCard },
  { href: "/admin/security", label: "Security", icon: Shield },
  { href: "/admin/settings", label: "Settings", icon: Settings },
];

interface AdminSidebarProps {
  className?: string;
}

export function AdminSidebar({ className }: AdminSidebarProps) {
  const pathname = usePathname();

  return (
    <aside
      className={cn(
        "hidden w-64 shrink-0 border-r bg-muted/30 lg:block",
        className
      )}
    >
      <div className="flex h-16 items-center border-b px-6">
        <Link href="/admin/dashboard" className="flex items-center gap-2 font-bold text-lg">
          <Ticket className="h-5 w-5 text-primary" />
          <span>Admin Panel</span>
        </Link>
      </div>

      <nav className="flex flex-col gap-1 p-4" aria-label="Admin navigation">
        {adminNavItems.map((item) => (
          <NavItem
            key={item.href}
            href={item.href}
            label={item.label}
            icon={item.icon}
            isActive={pathname.startsWith(item.href)}
          />
        ))}
      </nav>
    </aside>
  );
}
