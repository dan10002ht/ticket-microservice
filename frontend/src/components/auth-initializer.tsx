"use client";

import { useEffect } from "react";
import { useMe } from "@/lib/api/queries";
import { useAuthStore } from "@/stores";

/**
 * Invisible component that runs useMe() on mount to restore auth state
 * from cookies. Does not block render — UI shows immediately while
 * the auth check runs in the background.
 */
export function AuthInitializer() {
  const { data, isSuccess, isError } = useMe();
  const { setUser, clearUser, setHydrated } = useAuthStore();

  useEffect(() => {
    if (isSuccess && data) {
      setUser(data);
      setHydrated(true);
      return;
    }
    if (isError) {
      clearUser();
      setHydrated(true);
    }
  }, [isSuccess, isError, data, setUser, clearUser, setHydrated]);

  return null;
}
