import { AdminLayout } from "@/components/templates/admin-layout";
import { RoleGuard } from "@/components/role-guard";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={["admin", "super_admin"]}>
      <AdminLayout>{children}</AdminLayout>
    </RoleGuard>
  );
}
