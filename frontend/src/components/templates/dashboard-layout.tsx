import { AuthNavbar } from "@/components/organisms/layout/auth-navbar";
import { MobileBottomTabs } from "@/components/organisms/layout/mobile-bottom-tabs";
import { VerificationBanner } from "@/components/molecules/verification-banner";

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <div className="flex min-h-screen flex-col">
      <AuthNavbar />
      <VerificationBanner />
      <main className="flex-1 container mx-auto px-4 py-6 pb-20 md:pb-6">
        {children}
      </main>
      <MobileBottomTabs />
    </div>
  );
}
