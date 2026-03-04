import { OrgLayout } from "@/components/templates/org-layout";
import { RoleGuard } from "@/components/role-guard";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={["organization"]}>
      <OrgLayout>{children}</OrgLayout>
    </RoleGuard>
  );
}
