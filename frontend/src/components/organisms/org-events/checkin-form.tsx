"use client";

import { useState } from "react";
import { ScanLine, Loader2, CheckCircle2, XCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { useProcessCheckIn } from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import type { CheckIn } from "@/lib/api/types/checkin";
import type { ApiError } from "@/lib/api/types/common";

interface CheckinFormProps {
  eventId: string;
}

const statusConfig: Record<
  string,
  { label: string; variant: "default" | "destructive" | "secondary"; icon: typeof CheckCircle2 }
> = {
  success: { label: "Check-in Successful", variant: "default", icon: CheckCircle2 },
  invalid: { label: "Invalid Ticket", variant: "destructive", icon: XCircle },
  already_used: { label: "Already Used", variant: "secondary", icon: XCircle },
  cancelled: { label: "Ticket Cancelled", variant: "destructive", icon: XCircle },
};

export function CheckinForm({ eventId }: CheckinFormProps) {
  const processCheckIn = useProcessCheckIn(eventId);
  const [ticketId, setTicketId] = useState("");
  const [gate, setGate] = useState("");
  const [lastResult, setLastResult] = useState<CheckIn | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!ticketId.trim()) return;

    try {
      const result = await processCheckIn.mutateAsync({
        ticket_id: ticketId.trim(),
        event_id: eventId,
        gate: gate.trim() || undefined,
      });
      setLastResult(result);
      if (result.status === "success") {
        showToast.success("Check-in successful!");
        setTicketId("");
      }
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  const resultConfig = lastResult ? statusConfig[lastResult.status] : null;

  return (
    <div className="space-y-4">
      <form onSubmit={handleSubmit} className="space-y-4 rounded-lg border bg-card p-4">
        <h3 className="flex items-center gap-2 font-semibold">
          <ScanLine className="h-4 w-4" />
          Manual Check-in
        </h3>

        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="ticket-id">Ticket ID *</Label>
            <Input
              id="ticket-id"
              placeholder="Enter ticket ID or scan QR code"
              value={ticketId}
              onChange={(e) => setTicketId(e.target.value)}
              autoFocus
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="gate">Gate</Label>
            <Input
              id="gate"
              placeholder="e.g. Gate A, Main Entrance"
              value={gate}
              onChange={(e) => setGate(e.target.value)}
            />
          </div>
        </div>

        <Button
          type="submit"
          disabled={processCheckIn.isPending || !ticketId.trim()}
        >
          {processCheckIn.isPending && (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          )}
          Process Check-in
        </Button>
      </form>

      {/* Last result feedback */}
      {lastResult && resultConfig && (
        <div
          className={`flex items-center gap-3 rounded-lg border p-4 ${
            lastResult.status === "success"
              ? "border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-950"
              : "border-red-200 bg-red-50 dark:border-red-800 dark:bg-red-950"
          }`}
        >
          <resultConfig.icon
            className={`h-5 w-5 ${
              lastResult.status === "success"
                ? "text-green-600"
                : "text-red-600"
            }`}
          />
          <div className="flex-1">
            <p className="font-medium">{resultConfig.label}</p>
            <p className="text-xs text-muted-foreground">
              Ticket: {lastResult.ticket_id.slice(0, 8)}...
            </p>
          </div>
          <Badge variant={resultConfig.variant} className="capitalize">
            {lastResult.status.replace("_", " ")}
          </Badge>
        </div>
      )}
    </div>
  );
}
