"use client";

import { m } from "framer-motion";
import { cn } from "@/lib/utils";

interface PageTransitionProps {
  className?: string;
  children: React.ReactNode;
}

export function PageTransition({ className, children }: PageTransitionProps) {
  return (
    <m.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, ease: "easeOut" }}
      className={cn(className)}
    >
      {children}
    </m.div>
  );
}
