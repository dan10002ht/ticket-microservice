"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { Building2, LogOut, Menu, Shield, Ticket, User } from "lucide-react";
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
import { SearchDialog } from "@/components/molecules/search-dialog";
import { ThemeToggle } from "@/components/molecules/theme-toggle";
import { useAuthStore } from "@/stores";
import { useLogout } from "@/lib/api/queries";
import { cn } from "@/lib/utils";

const navLinks = [
  { href: "/", label: "Home" },
  { href: "/events", label: "Events" },
];

interface NavbarProps {
  className?: string;
}

export function Navbar({ className }: NavbarProps) {
  const router = useRouter();
  const { user, isAuthenticated, isHydrated } = useAuthStore();
  const logoutMutation = useLogout();

  const displayName = user
    ? `${user.firstName ?? ""} ${user.lastName ?? ""}`.trim() || user.email || "User"
    : "User";

  const initials = displayName
    .split(" ")
    .map((p) => p[0])
    .join("")
    .toUpperCase()
    .slice(0, 2) || "U";

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
            {navLinks.map((link) => (
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

        {/* Search + Theme + Auth */}
        <div className="flex items-center gap-2">
          <SearchDialog />
          <ThemeToggle />

          {/* Auth section — only render after hydration to avoid flash */}
          {isHydrated && (
            <>
              {isAuthenticated ? (
                /* Authenticated: user dropdown */
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="hidden md:inline-flex"
                    >
                      <Avatar className="h-8 w-8">
                        <AvatarFallback className="text-xs">{initials}</AvatarFallback>
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
                    {user?.role === "organization" && (
                      <DropdownMenuItem asChild>
                        <Link
                          href="/org/dashboard"
                          className="flex items-center gap-2"
                        >
                          <Building2 className="h-4 w-4" />
                          Organizer Dashboard
                        </Link>
                      </DropdownMenuItem>
                    )}
                    {(user?.role === "admin" ||
                      user?.role === "super_admin") && (
                      <DropdownMenuItem asChild>
                        <Link
                          href="/admin/dashboard"
                          className="flex items-center gap-2"
                        >
                          <Shield className="h-4 w-4" />
                          Admin Panel
                        </Link>
                      </DropdownMenuItem>
                    )}
                    <DropdownMenuItem asChild>
                      <Link
                        href="/my-bookings"
                        className="flex items-center gap-2"
                      >
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
              ) : (
                /* Unauthenticated: Login/Register buttons */
                <div className="hidden items-center gap-2 md:flex">
                  <Button variant="ghost" asChild>
                    <Link href="/login">Login</Link>
                  </Button>
                  <Button asChild>
                    <Link href="/register">Register</Link>
                  </Button>
                </div>
              )}
            </>
          )}

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
                {navLinks.map((link) => (
                  <Link
                    key={link.href}
                    href={link.href}
                    className="text-sm font-medium transition-colors hover:text-primary"
                  >
                    {link.label}
                  </Link>
                ))}

                {isHydrated && isAuthenticated ? (
                  <>
                    {user?.role === "organization" && (
                      <Link
                        href="/org/dashboard"
                        className="text-sm font-medium text-primary transition-colors hover:text-primary/80"
                      >
                        Organizer Dashboard
                      </Link>
                    )}
                    {(user?.role === "admin" ||
                      user?.role === "super_admin") && (
                      <Link
                        href="/admin/dashboard"
                        className="text-sm font-medium text-primary transition-colors hover:text-primary/80"
                      >
                        Admin Panel
                      </Link>
                    )}
                    <Link
                      href="/my-bookings"
                      className="text-sm font-medium transition-colors hover:text-primary"
                    >
                      Bookings
                    </Link>
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
                  </>
                ) : (
                  <div className="mt-4 flex flex-col gap-2">
                    <Button variant="outline" asChild>
                      <Link href="/login">Login</Link>
                    </Button>
                    <Button asChild>
                      <Link href="/register">Register</Link>
                    </Button>
                  </div>
                )}
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </nav>
    </header>
  );
}
