import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "My Events",
  description: "Manage your events.",
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return children;
}
