"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { CalendarDays, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { cn } from "@/lib/utils";

const sampleResults = [
  { title: "Summer Music Festival 2026", href: "/events/1" },
  { title: "Tech Conference Vietnam", href: "/events/2" },
  { title: "Stand-up Comedy Night", href: "/events/3" },
  { title: "Jazz Night at the Rooftop", href: "/events/4" },
];

interface SearchDialogProps {
  className?: string;
}

export function SearchDialog({ className }: SearchDialogProps) {
  const [open, setOpen] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((prev) => !prev);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  return (
    <>
      <Button
        variant="outline"
        onClick={() => setOpen(true)}
        className={cn(
          "relative h-9 w-9 p-0 md:h-9 md:w-60 md:justify-start md:px-3 md:py-2",
          className
        )}
        aria-label="Search events"
      >
        <Search className="h-4 w-4 md:mr-2" />
        <span className="hidden md:inline-flex text-sm text-muted-foreground">
          Search events...
        </span>
        <kbd className="pointer-events-none absolute right-1.5 top-1.5 hidden h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 md:flex">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Button>

      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput placeholder="Search events, bookings..." />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup heading="Events">
            {sampleResults.map((result) => (
              <CommandItem
                key={result.href}
                onSelect={() => {
                  router.push(result.href);
                  setOpen(false);
                }}
              >
                <CalendarDays className="mr-2 h-4 w-4" />
                <span>{result.title}</span>
              </CommandItem>
            ))}
          </CommandGroup>
        </CommandList>
      </CommandDialog>
    </>
  );
}
