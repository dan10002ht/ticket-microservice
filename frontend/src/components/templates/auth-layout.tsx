import Link from "next/link";
import { Ticket } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

interface AuthLayoutProps {
  children: React.ReactNode;
}

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="flex min-h-screen">
      {/* Decorative left panel (desktop only) */}
      <div className="hidden w-1/2 gradient-brand lg:flex lg:flex-col lg:items-center lg:justify-center">
        <div className="px-12 text-center text-white">
          <Ticket className="mx-auto h-16 w-16 mb-6" />
          <h2 className="text-3xl font-bold">Welcome to TicketBox</h2>
          <p className="mt-4 text-lg text-white/80">
            Your one-stop platform for discovering and booking amazing events
            across Vietnam.
          </p>
        </div>
      </div>

      {/* Form panel */}
      <div className="flex w-full items-center justify-center gradient-brand-subtle px-4 lg:w-1/2">
        <div className="w-full max-w-md">
          <div className="mb-8 flex flex-col items-center gap-2 lg:hidden">
            <Link
              href="/"
              className="flex items-center gap-2 font-bold text-2xl"
            >
              <Ticket className="h-8 w-8 text-primary" />
              <span>TicketBox</span>
            </Link>
          </div>
          <Card>
            <CardContent className="p-6">{children}</CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
