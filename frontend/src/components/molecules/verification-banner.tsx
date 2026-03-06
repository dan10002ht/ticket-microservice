"use client";

import Link from "next/link";
import { AlertTriangle, X } from "lucide-react";
import { useState } from "react";
import { useAuthStore } from "@/stores";

export function VerificationBanner() {
  const { user, isHydrated } = useAuthStore();
  const [dismissed, setDismissed] = useState(false);

  if (!isHydrated || !user || user.isVerified || dismissed) return null;

  return (
    <div className="border-b bg-yellow-50 dark:bg-yellow-950/30">
      <div className="container mx-auto flex items-center justify-between gap-3 px-4 py-2.5">
        <div className="flex items-center gap-2 text-sm">
          <AlertTriangle className="h-4 w-4 shrink-0 text-yellow-600 dark:text-yellow-500" />
          <p className="text-yellow-800 dark:text-yellow-200">
            Your email is not verified.{" "}
            <Link
              href="/verify"
              className="font-medium underline underline-offset-2 hover:text-yellow-900 dark:hover:text-yellow-100"
            >
              Verify now
            </Link>
          </p>
        </div>
        <button
          onClick={() => setDismissed(true)}
          className="shrink-0 rounded p-1 text-yellow-600 hover:bg-yellow-100 dark:text-yellow-500 dark:hover:bg-yellow-900/50"
          aria-label="Dismiss"
        >
          <X className="h-4 w-4" />
        </button>
      </div>
    </div>
  );
}
