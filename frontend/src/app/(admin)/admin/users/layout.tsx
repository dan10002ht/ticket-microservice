import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Manage Users",
  description: "View and manage all platform users.",
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return children;
}
