import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "My Bookings",
  description: "View and manage your event bookings.",
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return children;
}
