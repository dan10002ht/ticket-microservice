"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "next-themes";
import { LazyMotion, domAnimation } from "framer-motion";
import { useState } from "react";
import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000,
            refetchOnWindowFocus: false,
          },
        },
      })
  );

  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <LazyMotion features={domAnimation} strict>
        <QueryClientProvider client={queryClient}>
          <TooltipProvider>
            {children}
          </TooltipProvider>
          <Toaster richColors position="top-right" />
        </QueryClientProvider>
      </LazyMotion>
    </ThemeProvider>
  );
}
