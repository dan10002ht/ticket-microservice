"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { Loader2, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useEventZones, useEventSeats, useBulkCreateSeats } from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import type { ApiError } from "@/lib/api/types/common";
import { cn } from "@/lib/utils";

interface SeatManagementStepProps {
  eventId: string;
  className?: string;
}

interface BulkSeatForm {
  zone_id: string;
  row_prefix: string;
  start_number: number;
  count: number;
}

export function SeatManagementStep({
  eventId,
  className,
}: SeatManagementStepProps) {
  const { data: zones = [] } = useEventZones(eventId);
  const [selectedZoneId, setSelectedZoneId] = useState<string>("");
  const { data: seats = [], isLoading: seatsLoading } = useEventSeats(eventId, {
    zone_id: selectedZoneId || undefined,
  });
  const bulkCreateMutation = useBulkCreateSeats(eventId);

  const { register, handleSubmit, setValue, reset, watch } =
    useForm<BulkSeatForm>({
      defaultValues: {
        zone_id: "",
        row_prefix: "A",
        start_number: 1,
        count: 10,
      },
    });

  const watchZoneId = watch("zone_id");

  const onBulkCreate = async (data: BulkSeatForm) => {
    if (!data.zone_id) {
      showToast.error("Please select a zone");
      return;
    }

    const seatInputs = Array.from({ length: data.count }, (_, i) => ({
      seat_number: `${data.row_prefix}${data.start_number + i}`,
      row_number: data.row_prefix,
    }));

    try {
      await bulkCreateMutation.mutateAsync({
        zone_id: data.zone_id,
        seats: seatInputs,
      });
      showToast.success(`Created ${data.count} seats successfully!`);
      reset({ zone_id: data.zone_id, row_prefix: "A", start_number: 1, count: 10 });
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  return (
    <div className={cn("space-y-6", className)}>
      {/* Bulk Create Form */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Bulk Create Seats</CardTitle>
          <CardDescription>
            Quickly generate multiple seats for a zone
          </CardDescription>
        </CardHeader>
        <CardContent>
          {zones.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              No zones created yet. Add zones in the &quot;Zones & Pricing&quot; tab first.
            </p>
          ) : (
            <form onSubmit={handleSubmit(onBulkCreate)} className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-4">
                <div className="space-y-2">
                  <Label>Zone *</Label>
                  <Select
                    value={watchZoneId}
                    onValueChange={(val) => {
                      setValue("zone_id", val);
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select zone" />
                    </SelectTrigger>
                    <SelectContent>
                      {zones.map((z) => (
                        <SelectItem key={z.id} value={z.id}>
                          {z.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="row_prefix">Row</Label>
                  <Input
                    id="row_prefix"
                    placeholder="A"
                    {...register("row_prefix")}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="start_number">Start #</Label>
                  <Input
                    id="start_number"
                    type="number"
                    min={1}
                    {...register("start_number", { valueAsNumber: true })}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="count">Count</Label>
                  <Input
                    id="count"
                    type="number"
                    min={1}
                    max={100}
                    {...register("count", { valueAsNumber: true })}
                  />
                </div>
              </div>

              <Button type="submit" disabled={bulkCreateMutation.isPending}>
                {bulkCreateMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Plus className="mr-2 h-4 w-4" />
                )}
                Create Seats
              </Button>
            </form>
          )}
        </CardContent>
      </Card>

      {/* Seat List */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-lg">Seats</CardTitle>
              <CardDescription>
                {seats.length} seat{seats.length !== 1 ? "s" : ""} total
              </CardDescription>
            </div>
            {zones.length > 0 && (
              <Select
                value={selectedZoneId}
                onValueChange={setSelectedZoneId}
              >
                <SelectTrigger className="w-[180px]">
                  <SelectValue placeholder="All zones" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All zones</SelectItem>
                  {zones.map((z) => (
                    <SelectItem key={z.id} value={z.id}>
                      {z.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {seatsLoading ? (
            <p className="text-sm text-muted-foreground">Loading seats...</p>
          ) : seats.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              No seats created yet. Use the bulk create form above.
            </p>
          ) : (
            <div className="max-h-[400px] overflow-auto rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Seat #</TableHead>
                    <TableHead>Row</TableHead>
                    <TableHead>Zone</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Price</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {seats.map((seat) => {
                    const zone = zones.find((z) => z.id === seat.zone_id);
                    return (
                      <TableRow key={seat.id}>
                        <TableCell className="font-medium">
                          {seat.seat_number}
                        </TableCell>
                        <TableCell>{seat.row_number || "—"}</TableCell>
                        <TableCell>{zone?.name || seat.zone_id}</TableCell>
                        <TableCell>
                          <span
                            className={cn(
                              "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium",
                              seat.status === "available"
                                ? "bg-green-100 text-green-800"
                                : seat.status === "reserved"
                                  ? "bg-yellow-100 text-yellow-800"
                                  : "bg-gray-100 text-gray-800"
                            )}
                          >
                            {seat.status || "available"}
                          </span>
                        </TableCell>
                        <TableCell className="text-right">
                          {seat.final_price
                            ? `${seat.final_price} ${seat.currency || "VND"}`
                            : "—"}
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
