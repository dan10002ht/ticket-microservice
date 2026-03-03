import { PublicLayout } from "@/components/templates/public-layout";

export default function Layout({ children }: { children: React.ReactNode }) {
  return <PublicLayout>{children}</PublicLayout>;
}
