"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  CalendarDays,
  BarChart3,
  Settings,
  Ticket,
} from "lucide-react";
import { NavItem } from "@/components/molecules/nav-item";
import { cn } from "@/lib/utils";

const orgNavItems = [
  { href: "/org/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/org/events", label: "Events", icon: CalendarDays },
  { href: "/org/analytics", label: "Analytics", icon: BarChart3 },
  { href: "/org/settings", label: "Settings", icon: Settings },
];

interface OrgSidebarProps {
  className?: string;
}

export function OrgSidebar({ className }: OrgSidebarProps) {
  const pathname = usePathname();

  return (
    <aside
      className={cn(
        "hidden w-64 shrink-0 border-r bg-muted/30 lg:block",
        className
      )}
    >
      <div className="flex h-16 items-center border-b px-6">
        <Link href="/org/dashboard" className="flex items-center gap-2 font-bold text-lg">
          <Ticket className="h-5 w-5 text-primary" />
          <span>Organizer</span>
        </Link>
      </div>

      <nav className="flex flex-col gap-1 p-4" aria-label="Organization navigation">
        {orgNavItems.map((item) => (
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
