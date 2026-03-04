"use client";

import { useEffect, useRef } from "react";
import { Clock, AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCountdown } from "@/lib/hooks/use-countdown";

interface ReservationTimerProps {
  expiresAt: string | null;
  onExpired: () => void;
}

export function ReservationTimer({ expiresAt, onExpired }: ReservationTimerProps) {
  const { minutes, seconds, isExpired } = useCountdown(expiresAt);
  const expiredCalled = useRef(false);

  useEffect(() => {
    if (isExpired && !expiredCalled.current) {
      expiredCalled.current = true;
      onExpired();
    }
  }, [isExpired, onExpired]);

  // Reset ref when expiresAt changes (new reservation)
  useEffect(() => {
    expiredCalled.current = false;
  }, [expiresAt]);

  if (!expiresAt) return null;

  if (isExpired) {
    return (
      <div className="flex items-center gap-2 rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
        <AlertTriangle className="h-4 w-4" />
        <span>Reservation expired</span>
        <Button
          variant="ghost"
          size="sm"
          className="ml-auto h-auto px-2 py-1 text-xs"
          onClick={onExpired}
        >
          Start Over
        </Button>
      </div>
    );
  }

  const totalSeconds = minutes * 60 + seconds;
  const colorClass =
    totalSeconds > 300
      ? "text-green-600"
      : totalSeconds > 120
        ? "text-yellow-600"
        : "text-red-600";

  return (
    <div className={`flex items-center gap-2 text-sm font-medium ${colorClass}`}>
      <Clock className="h-4 w-4" />
      <span>
        Time remaining: {String(minutes).padStart(2, "0")}:
        {String(seconds).padStart(2, "0")}
      </span>
    </div>
  );
}
