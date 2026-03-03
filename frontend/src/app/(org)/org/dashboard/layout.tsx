import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Organizer Dashboard",
  description: "Overview of your events and sales.",
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return children;
}
