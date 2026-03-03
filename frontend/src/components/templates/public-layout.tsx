import { Navbar } from "@/components/organisms/layout/navbar";
import { Footer } from "@/components/organisms/layout/footer";
import { MobileBottomTabs } from "@/components/organisms/layout/mobile-bottom-tabs";

interface PublicLayoutProps {
  children: React.ReactNode;
}

export function PublicLayout({ children }: PublicLayoutProps) {
  return (
    <div className="flex min-h-screen flex-col">
      <Navbar />
      <main className="flex-1 pb-16 md:pb-0">{children}</main>
      <Footer className="hidden md:block" />
      <MobileBottomTabs />
    </div>
  );
}
