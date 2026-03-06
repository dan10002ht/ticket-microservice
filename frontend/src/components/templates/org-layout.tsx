"use client";

import { Sidebar } from "@/components/organisms/layout/sidebar";
import { Topbar } from "@/components/organisms/layout/topbar";
import { VerificationBanner } from "@/components/molecules/verification-banner";
import { orgNavItems } from "@/components/organisms/layout/nav-config";

interface OrgLayoutProps {
  children: React.ReactNode;
}

export function OrgLayout({ children }: OrgLayoutProps) {
  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar navItems={orgNavItems} brandLabel="Organizer" brandHref="/org/dashboard" />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Topbar navItems={orgNavItems} brandLabel="Organizer" brandHref="/org/dashboard" />
        <VerificationBanner />
        <main className="flex-1 overflow-y-auto p-6">{children}</main>
      </div>
    </div>
  );
}
