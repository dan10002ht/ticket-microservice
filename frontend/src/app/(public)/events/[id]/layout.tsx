import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Event Details",
  description: "View event details, pricing, and availability.",
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return children;
}
