"use client";

import { useState, useEffect, useCallback } from "react";

interface CountdownResult {
  minutes: number;
  seconds: number;
  totalSeconds: number;
  isExpired: boolean;
}

function getRemainingSeconds(expiresAt: string): number {
  const diff = new Date(expiresAt).getTime() - Date.now();
  return Math.max(0, Math.floor(diff / 1000));
}

export function useCountdown(expiresAt: string | null): CountdownResult {
  const [totalSeconds, setTotalSeconds] = useState(() =>
    expiresAt ? getRemainingSeconds(expiresAt) : 0
  );

  const isExpired = totalSeconds <= 0 && expiresAt !== null;

  useEffect(() => {
    if (!expiresAt) return;

    setTotalSeconds(getRemainingSeconds(expiresAt));

    const interval = setInterval(() => {
      const remaining = getRemainingSeconds(expiresAt);
      setTotalSeconds(remaining);
      if (remaining <= 0) clearInterval(interval);
    }, 1000);

    return () => clearInterval(interval);
  }, [expiresAt]);

  return {
    minutes: Math.floor(totalSeconds / 60),
    seconds: totalSeconds % 60,
    totalSeconds,
    isExpired,
  };
}
