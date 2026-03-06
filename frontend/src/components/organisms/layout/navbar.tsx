"use client";

import Link from "next/link";
import { LogOut, Menu, Ticket } from "lucide-react";
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

const navLinks = [
  { href: "/", label: "Home" },
  { href: "/events", label: "Events" },
];

interface NavbarProps {
  className?: string;
}

export function Navbar({ className }: NavbarProps) {
  const { user, isAuthenticated, isHydrated } = useAuthStore();
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
                <UserAvatarMenu className="hidden md:inline-flex" />
              ) : (
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
