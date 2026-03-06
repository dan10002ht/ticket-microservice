"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Ticket } from "lucide-react";
import { NavItem } from "@/components/molecules/nav-item";
import { cn } from "@/lib/utils";
import type { NavItemConfig } from "./nav-config";

interface SidebarProps {
  navItems: NavItemConfig[];
  brandLabel: string;
  brandHref: string;
  className?: string;
}

export function Sidebar({ navItems, brandLabel, brandHref, className }: SidebarProps) {
  const pathname = usePathname();

  return (
    <aside
      className={cn(
        "hidden w-64 shrink-0 border-r bg-muted/30 lg:block",
        className
      )}
    >
      <div className="flex h-16 items-center border-b px-6">
        <Link href={brandHref} className="flex items-center gap-2 font-bold text-lg">
          <Ticket className="h-5 w-5 text-primary" />
          <span>{brandLabel}</span>
        </Link>
      </div>

      <nav className="flex flex-col gap-1 p-4" aria-label={`${brandLabel} navigation`}>
        {navItems.map((item) => (
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
