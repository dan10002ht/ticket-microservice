"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import type { LucideIcon } from "lucide-react";
import {
  Bell,
  Building2,
  Home,
  LogOut,
  Shield,
  User,
  Menu,
  Ticket,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetTitle,
} from "@/components/ui/sheet";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Breadcrumbs } from "@/components/molecules/breadcrumbs";
import { ThemeToggle } from "@/components/molecules/theme-toggle";
import { useAuthStore } from "@/stores";
import { useLogout } from "@/lib/api/queries";
import { cn } from "@/lib/utils";

interface NavItemConfig {
  href: string;
  label: string;
  icon: LucideIcon;
}

interface TopbarProps {
  className?: string;
  navItems?: NavItemConfig[];
  brandLabel?: string;
  brandHref?: string;
}

export function Topbar({ className, navItems = [], brandLabel = "Organizer", brandHref = "/org/dashboard" }: TopbarProps) {
  const { user } = useAuthStore();
  const logoutMutation = useLogout();
  const pathname = usePathname();

  const displayName = user
    ? `${user.first_name ?? ""} ${user.last_name ?? ""}`.trim() || user.email || "User"
    : "User";

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

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="hidden lg:inline-flex">
              <Avatar className="h-8 w-8">
                <AvatarFallback className="text-xs">
                  {displayName
                    .split(" ")
                    .map((p) => p[0])
                    .join("")
                    .toUpperCase()
                    .slice(0, 2) || "U"}
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-48">
            <div className="px-2 py-1.5">
              <p className="text-sm font-medium">{displayName}</p>
              {user?.email && (
                <p className="text-xs text-muted-foreground">{user.email}</p>
              )}
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem asChild>
              <Link href="/" className="flex items-center gap-2">
                <Home className="h-4 w-4" />
                Home
              </Link>
            </DropdownMenuItem>
            {user?.role === "organization" && (
              <DropdownMenuItem asChild>
                <Link href="/org/dashboard" className="flex items-center gap-2">
                  <Building2 className="h-4 w-4" />
                  Organizer Dashboard
                </Link>
              </DropdownMenuItem>
            )}
            {(user?.role === "admin" || user?.role === "super_admin") && (
              <DropdownMenuItem asChild>
                <Link href="/admin/dashboard" className="flex items-center gap-2">
                  <Shield className="h-4 w-4" />
                  Admin Panel
                </Link>
              </DropdownMenuItem>
            )}
            <DropdownMenuItem asChild>
              <Link href="/my-bookings" className="flex items-center gap-2">
                <Ticket className="h-4 w-4" />
                Bookings
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link href="/profile" className="flex items-center gap-2">
                <User className="h-4 w-4" />
                Profile
              </Link>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="flex items-center gap-2 text-destructive"
              onClick={() => logoutMutation.mutate()}
              disabled={logoutMutation.isPending}
            >
              <LogOut className="h-4 w-4" />
              {logoutMutation.isPending ? "Logging out..." : "Logout"}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
