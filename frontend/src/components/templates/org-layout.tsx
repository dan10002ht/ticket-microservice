import { OrgSidebar } from "@/components/organisms/layout/org-sidebar";
import { Topbar } from "@/components/organisms/layout/topbar";

interface OrgLayoutProps {
  children: React.ReactNode;
}

export function OrgLayout({ children }: OrgLayoutProps) {
  return (
    <div className="flex h-screen overflow-hidden">
      <OrgSidebar />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Topbar />
        <main className="flex-1 overflow-y-auto p-6">{children}</main>
      </div>
    </div>
  );
}
