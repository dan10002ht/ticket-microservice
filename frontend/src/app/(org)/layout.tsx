import { OrgLayout } from "@/components/templates/org-layout";

export default function Layout({ children }: { children: React.ReactNode }) {
  return <OrgLayout>{children}</OrgLayout>;
}
