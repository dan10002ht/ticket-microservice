import { Shield } from "lucide-react";
import { PageHeader } from "@/components/molecules/page-header";
import { EmptyState } from "@/components/molecules/empty-state";

export default function AdminSecurityPage() {
  return (
    <>
      <PageHeader
        title="Security"
        description="Security settings and audit logs"
      />
      <EmptyState
        icon={Shield}
        title="Coming Soon"
        description="Security dashboard with audit logs, access policies, and threat monitoring will be available here."
        className="mt-12"
      />
    </>
  );
}
