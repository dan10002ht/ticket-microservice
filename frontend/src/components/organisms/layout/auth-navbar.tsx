"use client";

import Link from "next/link";
import { Bell, Building2, LogOut, Menu, Shield, Ticket } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetTitle,
} from "@/components/ui/sheet";
import { SearchDialog } from "@/components/molecules/search-dialog";
import { ThemeToggle } from "@/components/molecules/theme-toggle";
import { UserAvatarMenu } from "@/components/molecules/user-avatar-menu";
import { useAuthStore } from "@/stores";
import { useLogout } from "@/lib/api/queries";
import { cn } from "@/lib/utils";

const dashboardLinks = [
  { href: "/my-bookings", label: "Bookings" },
  { href: "/my-tickets", label: "Tickets" },
];

interface AuthNavbarProps {
  className?: string;
}

export function AuthNavbar({ className }: AuthNavbarProps) {
  const { user } = useAuthStore();
  const logoutMutation = useLogout();

  return (
    <header
      className={cn(
        "sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60",
        className,
      )}
    >
      <nav className="container mx-auto flex h-16 items-center justify-between px-4">
        {/* Logo + Nav links (left-aligned) */}
        <div className="flex items-center gap-8">
          <Link href="/" className="flex items-center gap-2 font-bold text-xl">
            <Ticket className="h-6 w-6 text-primary" />
            <span>TicketBox</span>
          </Link>

          {/* Desktop nav */}
          <div className="hidden items-center gap-6 md:flex">
            {dashboardLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
              >
                {link.label}
              </Link>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-2">
          <SearchDialog />
          <ThemeToggle />
          <Button
            variant="ghost"
            size="icon"
            aria-label="Notifications"
            className="hidden md:inline-flex"
          >
            <Bell className="h-5 w-5" />
          </Button>

          {/* Desktop: reuse shared avatar dropdown */}
          <UserAvatarMenu className="hidden md:inline-flex" />

          {/* Mobile menu */}
          <Sheet>
            <SheetTrigger asChild className="md:hidden">
              <Button variant="ghost" size="icon" aria-label="Open menu">
                <Menu className="h-5 w-5" />
              </Button>
            </SheetTrigger>
            <SheetContent side="right" className="w-72 px-6 pt-12">
              <SheetTitle className="sr-only">Navigation menu</SheetTitle>
              <div className="flex flex-col gap-4">
                {dashboardLinks.map((link) => (
                  <Link
                    key={link.href}
                    href={link.href}
                    className="text-sm font-medium transition-colors hover:text-primary"
                  >
                    {link.label}
                  </Link>
                ))}
                <Link
                  href="/notifications"
                  className="flex items-center gap-2 text-sm font-medium transition-colors hover:text-primary"
                >
                  <Bell className="h-4 w-4" />
                  Notifications
                </Link>

                {user?.role === "organization" && (
                  <Link
                    href="/org/dashboard"
                    className="flex items-center gap-2 text-sm font-medium text-primary transition-colors hover:text-primary/80"
                  >
                    <Building2 className="h-4 w-4" />
                    Organizer Dashboard
                  </Link>
                )}
                {(user?.role === "admin" || user?.role === "super_admin") && (
                  <Link
                    href="/admin/dashboard"
                    className="flex items-center gap-2 text-sm font-medium text-primary transition-colors hover:text-primary/80"
                  >
                    <Shield className="h-4 w-4" />
                    Admin Panel
                  </Link>
                )}

                <Link
                  href="/profile"
                  className="text-sm font-medium transition-colors hover:text-primary"
                >
                  Profile
                </Link>
                <div className="mt-4">
                  <Button
                    variant="outline"
                    className="w-full text-destructive"
                    onClick={() => logoutMutation.mutate()}
                    disabled={logoutMutation.isPending}
                  >
                    <LogOut className="mr-2 h-4 w-4" />
                    {logoutMutation.isPending ? "Logging out..." : "Logout"}
                  </Button>
                </div>
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </nav>
    </header>
  );
}
