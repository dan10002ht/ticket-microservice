"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, CalendarDays, LogIn, Ticket, User } from "lucide-react";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/stores";

const publicTabs = [
  { href: "/", label: "Home", icon: Home },
  { href: "/events", label: "Events", icon: CalendarDays },
  { href: "/login", label: "Login", icon: LogIn },
];

const authTabs = [
  { href: "/", label: "Home", icon: Home },
  { href: "/events", label: "Events", icon: CalendarDays },
  { href: "/my-bookings", label: "Bookings", icon: Ticket },
  { href: "/profile", label: "Profile", icon: User },
];

interface MobileBottomTabsProps {
  className?: string;
}

export function MobileBottomTabs({ className }: MobileBottomTabsProps) {
  const pathname = usePathname();
  const { isAuthenticated, isHydrated } = useAuthStore();

  const tabs = isHydrated && isAuthenticated ? authTabs : publicTabs;

  return (
    <nav
      className={cn(
        "fixed bottom-0 left-0 right-0 z-50 border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 md:hidden",
        className
      )}
      aria-label="Mobile navigation"
    >
      <div className="flex h-16 items-center justify-around px-2">
        {tabs.map((tab) => {
          const isActive =
            tab.href === "/"
              ? pathname === "/"
              : pathname.startsWith(tab.href);

          return (
            <Link
              key={tab.href}
              href={tab.href}
              className={cn(
                "flex flex-col items-center gap-1 px-3 py-2 text-xs transition-colors",
                isActive
                  ? "text-primary"
                  : "text-muted-foreground hover:text-foreground"
              )}
            >
              <tab.icon className="h-5 w-5" />
              <span>{tab.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
