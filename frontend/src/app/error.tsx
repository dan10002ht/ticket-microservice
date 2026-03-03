"use client";

import { useEffect } from "react";
import { AlertTriangle, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <div className="text-center">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-destructive/10">
          <AlertTriangle className="h-8 w-8 text-destructive" />
        </div>
        <h2 className="mt-4 text-2xl font-semibold">Something went wrong</h2>
        <p className="mt-2 max-w-md text-sm text-muted-foreground">
          An unexpected error occurred. Please try again.
        </p>
        <Button onClick={reset} className="mt-6">
          <RotateCcw className="mr-2 h-4 w-4" />
          Try again
        </Button>
      </div>
    </div>
  );
}
