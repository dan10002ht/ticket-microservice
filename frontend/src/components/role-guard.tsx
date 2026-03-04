"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/stores";

interface RoleGuardProps {
  allowedRoles: string[];
  fallbackUrl?: string;
  children: React.ReactNode;
}

export function RoleGuard({
  allowedRoles,
  fallbackUrl = "/",
  children,
}: RoleGuardProps) {
  const { user, isHydrated, isAuthenticated } = useAuthStore();
  const router = useRouter();

  const userRole = user?.role ?? "";
  const isAllowed = isAuthenticated && allowedRoles.includes(userRole);

  useEffect(() => {
    if (!isHydrated) return;

    if (!isAuthenticated) {
      router.replace("/login");
      return;
    }

    if (!allowedRoles.includes(userRole)) {
      router.replace(fallbackUrl);
    }
  }, [isHydrated, isAuthenticated, userRole, allowedRoles, fallbackUrl, router]);

  // Don't render anything until hydrated and authorized
  if (!isHydrated || !isAllowed) {
    return null;
  }

  return <>{children}</>;
}
