import { Settings } from "lucide-react";
import { PageHeader } from "@/components/molecules/page-header";
import { EmptyState } from "@/components/molecules/empty-state";

export default function AdminSettingsPage() {
  return (
    <>
      <PageHeader
        title="Settings"
        description="Platform configuration and preferences"
      />
      <EmptyState
        icon={Settings}
        title="Coming Soon"
        description="Platform settings including system configuration, feature flags, and notification preferences will be available here."
        className="mt-12"
      />
    </>
  );
}
