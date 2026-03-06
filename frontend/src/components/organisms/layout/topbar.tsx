"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Bell,
  LogOut,
  Menu,
  Ticket,
  User,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetTitle,
} from "@/components/ui/sheet";
import { Breadcrumbs } from "@/components/molecules/breadcrumbs";
import { ThemeToggle } from "@/components/molecules/theme-toggle";
import { UserAvatarMenu } from "@/components/molecules/user-avatar-menu";
import { useLogout } from "@/lib/api/queries";
import { cn } from "@/lib/utils";
import type { NavItemConfig } from "./nav-config";

interface TopbarProps {
  className?: string;
  navItems?: NavItemConfig[];
  brandLabel?: string;
  brandHref?: string;
}

export function Topbar({ className, navItems = [], brandLabel = "Organizer", brandHref = "/org/dashboard" }: TopbarProps) {
  const logoutMutation = useLogout();
  const pathname = usePathname();

  return (
    <header
      className={cn(
        "flex h-16 items-center justify-between border-b bg-background px-4 lg:px-6",
        className
      )}
    >
      <div className="flex items-center gap-3">
        {/* Mobile hamburger menu — visible below lg where sidebar is hidden */}
        {navItems.length > 0 && <Sheet>
          <SheetTrigger asChild className="lg:hidden">
            <Button variant="ghost" size="icon" aria-label="Open menu">
              <Menu className="h-5 w-5" />
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-72 px-0 pt-0">
            <SheetTitle className="sr-only">Navigation menu</SheetTitle>
            <div className="flex h-16 items-center border-b px-6">
              <Link href={brandHref} className="flex items-center gap-2 font-bold text-lg">
                <Ticket className="h-5 w-5 text-primary" />
                <span>{brandLabel}</span>
              </Link>
            </div>
            <nav className="flex flex-col gap-1 p-4">
              {navItems.map((item) => {
                const Icon = item.icon;
                const isActive = pathname.startsWith(item.href);
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={cn(
                      "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground",
                      isActive && "bg-accent text-accent-foreground"
                    )}
                  >
                    <Icon className="h-4 w-4" />
                    {item.label}
                  </Link>
                );
              })}
            </nav>
            <div className="mt-auto border-t p-4">
              <Link
                href="/profile"
                className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent"
              >
                <User className="h-4 w-4" />
                Profile
              </Link>
              <Button
                variant="ghost"
                className="w-full justify-start gap-3 px-3 text-destructive hover:text-destructive"
                onClick={() => logoutMutation.mutate()}
                disabled={logoutMutation.isPending}
              >
                <LogOut className="h-4 w-4" />
                {logoutMutation.isPending ? "Logging out..." : "Logout"}
              </Button>
            </div>
          </SheetContent>
        </Sheet>}

        <div className="hidden lg:block">
          <Breadcrumbs />
        </div>
      </div>

      <div className="flex items-center gap-2">
        <ThemeToggle />
        <Button variant="ghost" size="icon" aria-label="Notifications">
          <Bell className="h-5 w-5" />
        </Button>

        <UserAvatarMenu showHomeLink className="hidden lg:inline-flex" />
      </div>
    </header>
  );
}
