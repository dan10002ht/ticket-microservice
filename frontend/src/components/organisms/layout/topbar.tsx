"use client";

import { Bell, LogOut, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { AvatarWithName } from "@/components/molecules/avatar-with-name";
import { Breadcrumbs } from "@/components/molecules/breadcrumbs";
import { ThemeToggle } from "@/components/molecules/theme-toggle";
import { cn } from "@/lib/utils";

interface TopbarProps {
  className?: string;
}

export function Topbar({ className }: TopbarProps) {
  return (
    <header
      className={cn(
        "flex h-16 items-center justify-between border-b bg-background px-6",
        className
      )}
    >
      <Breadcrumbs />

      <div className="flex items-center gap-2">
        <ThemeToggle />
        <Button variant="ghost" size="icon" aria-label="Notifications">
          <Bell className="h-5 w-5" />
        </Button>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="gap-2 px-2">
              <AvatarWithName name="Admin User" size="sm" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-48">
            <DropdownMenuItem className="flex items-center gap-2">
              <User className="h-4 w-4" />
              Profile
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="flex items-center gap-2 text-destructive">
              <LogOut className="h-4 w-4" />
              Logout
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
