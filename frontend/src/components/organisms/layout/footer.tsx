import Link from "next/link";
import { Ticket } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";

const footerLinks = [
  {
    title: "Product",
    links: [
      { href: "/events", label: "Browse Events" },
      { href: "/register", label: "Create Account" },
    ],
  },
  {
    title: "Support",
    links: [
      { href: "#", label: "Help Center" },
      { href: "#", label: "Contact Us" },
    ],
  },
  {
    title: "Legal",
    links: [
      { href: "#", label: "Privacy Policy" },
      { href: "#", label: "Terms of Service" },
    ],
  },
];

interface FooterProps {
  className?: string;
}

export function Footer({ className }: FooterProps) {
  return (
    <footer className={cn("border-t bg-muted/30", className)}>
      <div className="container mx-auto px-4 py-12">
        <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
          {/* Brand */}
          <div>
            <Link href="/" className="flex items-center gap-2 font-bold text-lg">
              <Ticket className="h-5 w-5 text-primary" />
              <span>TicketBox</span>
            </Link>
            <p className="mt-2 text-sm text-muted-foreground">
              Your one-stop platform for event tickets.
            </p>
          </div>

          {/* Link groups */}
          {footerLinks.map((group) => (
            <div key={group.title}>
              <h3 className="font-semibold text-sm">{group.title}</h3>
              <ul className="mt-3 space-y-2">
                {group.links.map((link) => (
                  <li key={link.label}>
                    <Link
                      href={link.href}
                      className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <Separator className="my-8" />

        <p className="text-center text-xs text-muted-foreground">
          &copy; {new Date().getFullYear()} TicketBox. All rights reserved.
        </p>
      </div>
    </footer>
  );
}
