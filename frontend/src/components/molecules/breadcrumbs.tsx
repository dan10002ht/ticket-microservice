"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Fragment } from "react";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { cn } from "@/lib/utils";

const labelMap: Record<string, string> = {
  admin: "Admin",
  org: "Organizer",
  dashboard: "Dashboard",
  users: "Users",
  events: "Events",
  payments: "Payments",
  security: "Security",
  settings: "Settings",
  analytics: "Analytics",
  "my-bookings": "My Bookings",
  "my-tickets": "My Tickets",
  profile: "Profile",
  create: "Create",
  edit: "Edit",
};

// Root sections that should link to their dashboard instead of themselves
const rootRedirects: Record<string, string> = {
  org: "/org/dashboard",
  admin: "/admin/dashboard",
};

interface BreadcrumbsProps {
  className?: string;
}

export function Breadcrumbs({ className }: BreadcrumbsProps) {
  const pathname = usePathname();

  const segments = pathname.split("/").filter(Boolean);

  if (segments.length === 0) return null;

  return (
    <Breadcrumb className={cn(className)}>
      <BreadcrumbList>
        {segments.map((segment, index) => {
          const rawHref = "/" + segments.slice(0, index + 1).join("/");
          const href = rootRedirects[segment] && index === 0 ? rootRedirects[segment] : rawHref;
          const label =
            labelMap[segment] ??
            segment.charAt(0).toUpperCase() + segment.slice(1);
          const isLast = index === segments.length - 1;

          return (
            <Fragment key={rawHref}>
              {index > 0 && <BreadcrumbSeparator />}
              <BreadcrumbItem>
                {isLast ? (
                  <BreadcrumbPage>{label}</BreadcrumbPage>
                ) : (
                  <BreadcrumbLink asChild>
                    <Link href={href}>{label}</Link>
                  </BreadcrumbLink>
                )}
              </BreadcrumbItem>
            </Fragment>
          );
        })}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
