import { Settings } from "lucide-react";
import { PageHeader } from "@/components/molecules/page-header";
import { EmptyState } from "@/components/molecules/empty-state";

export default function OrgSettingsPage() {
  return (
    <>
      <PageHeader
        title="Settings"
        description="Organization preferences and configuration"
      />
      <EmptyState
        icon={Settings}
        title="Coming Soon"
        description="Organization settings including profile, team management, and notification preferences will be available here."
        className="mt-12"
      />
    </>
  );
}
