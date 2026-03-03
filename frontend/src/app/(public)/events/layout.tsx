import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Events",
  description: "Browse all upcoming events and find your next experience.",
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return children;
}
