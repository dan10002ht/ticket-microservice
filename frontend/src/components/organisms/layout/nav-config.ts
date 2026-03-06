import {
  LayoutDashboard,
  CalendarDays,
  BarChart3,
  Settings,
  Users,
  CreditCard,
  Shield,
} from "lucide-react";
import type { LucideIcon } from "lucide-react";

export interface NavItemConfig {
  href: string;
  label: string;
  icon: LucideIcon;
}

export const orgNavItems: NavItemConfig[] = [
  { href: "/org/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/org/events", label: "Events", icon: CalendarDays },
  { href: "/org/analytics", label: "Analytics", icon: BarChart3 },
  { href: "/org/settings", label: "Settings", icon: Settings },
];

export const adminNavItems: NavItemConfig[] = [
  { href: "/admin/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/admin/users", label: "Users", icon: Users },
  { href: "/admin/events", label: "Events", icon: CalendarDays },
  { href: "/admin/payments", label: "Payments", icon: CreditCard },
  { href: "/admin/security", label: "Security", icon: Shield },
  { href: "/admin/settings", label: "Settings", icon: Settings },
];
