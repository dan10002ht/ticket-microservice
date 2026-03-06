"use client";

import { Sidebar } from "@/components/organisms/layout/sidebar";
import { Topbar } from "@/components/organisms/layout/topbar";
import { VerificationBanner } from "@/components/molecules/verification-banner";
import { adminNavItems } from "@/components/organisms/layout/nav-config";

interface AdminLayoutProps {
  children: React.ReactNode;
}

export function AdminLayout({ children }: AdminLayoutProps) {
  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar navItems={adminNavItems} brandLabel="Admin Panel" brandHref="/admin/dashboard" />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Topbar navItems={adminNavItems} brandLabel="Admin Panel" brandHref="/admin/dashboard" />
        <VerificationBanner />
        <main className="flex-1 overflow-y-auto p-6">{children}</main>
      </div>
    </div>
  );
}
