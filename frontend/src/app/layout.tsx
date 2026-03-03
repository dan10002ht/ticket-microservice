import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { Providers } from "./providers";
import "./globals.css";

const inter = Inter({
  subsets: ["latin", "vietnamese"],
  display: "swap",
  variable: "--font-inter",
});

export const metadata: Metadata = {
  title: {
    default: "TicketBox — Event Ticketing Platform",
    template: "%s | TicketBox",
  },
  description:
    "Browse events, book tickets, and manage your bookings all in one place.",
  keywords: ["events", "tickets", "booking", "concerts", "entertainment"],
  openGraph: {
    type: "website",
    siteName: "TicketBox",
    title: "TicketBox — Event Ticketing Platform",
    description:
      "Browse events, book tickets, and manage your bookings all in one place.",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="vi" suppressHydrationWarning>
      <body className={`${inter.variable} font-sans antialiased`}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
