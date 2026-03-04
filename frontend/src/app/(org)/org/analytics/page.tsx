import { BarChart3 } from "lucide-react";
import { PageHeader } from "@/components/molecules/page-header";
import { EmptyState } from "@/components/molecules/empty-state";

export default function OrgAnalyticsPage() {
  return (
    <>
      <PageHeader
        title="Analytics"
        description="Event performance and insights"
      />
      <EmptyState
        icon={BarChart3}
        title="Coming Soon"
        description="Analytics dashboard with event performance metrics, ticket sales trends, and audience insights will be available here."
        className="mt-12"
      />
    </>
  );
}
