"use client";

import { useState } from "react";
import { toast } from "sonner";
import { Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useEventZones,
  useCreateZone,
  useEventPricing,
  useCreatePricing,
} from "@/lib/api/queries";
import type { ApiError } from "@/lib/api/types/common";
import { cn } from "@/lib/utils";

interface ZonePricingStepProps {
  eventId: string;
  onNext: () => void;
  onSkip: () => void;
  className?: string;
}

const ZONE_TYPES = [
  { value: "seated", label: "Seated" },
  { value: "standing", label: "Standing" },
  { value: "vip", label: "VIP" },
] as const;

const ZONE_COLORS = [
  "#3B82F6", // blue
  "#EF4444", // red
  "#10B981", // green
  "#F59E0B", // amber
  "#8B5CF6", // purple
  "#EC4899", // pink
  "#06B6D4", // cyan
  "#F97316", // orange
];

export function ZonePricingStep({
  eventId,
  onNext,
  onSkip,
  className,
}: ZonePricingStepProps) {
  const { data: zones = [], isLoading: zonesLoading } = useEventZones(eventId);
  const { data: pricing = [], isLoading: pricingLoading } =
    useEventPricing(eventId);
  const createZoneMutation = useCreateZone(eventId);
  const createPricingMutation = useCreatePricing(eventId);

  // Zone form state
  const [zoneName, setZoneName] = useState("");
  const [zoneType, setZoneType] = useState<string>("seated");
  const [seatCount, setSeatCount] = useState("");
  const [zoneColor, setZoneColor] = useState(ZONE_COLORS[0]);

  // Pricing form state
  const [pricingZoneId, setPricingZoneId] = useState("");
  const [pricingName, setPricingName] = useState("");
  const [pricingPrice, setPricingPrice] = useState("");
  const [pricingCurrency, setPricingCurrency] = useState("USD");

  const handleAddZone = async () => {
    if (!zoneName.trim()) {
      toast.error("Zone name is required");
      return;
    }

    try {
      await createZoneMutation.mutateAsync({
        name: zoneName.trim(),
        zone_type: zoneType as "seated" | "standing" | "vip",
        seat_count: seatCount ? parseInt(seatCount, 10) : undefined,
        color: zoneColor,
      });
      toast.success("Zone added!");
      setZoneName("");
      setSeatCount("");
    } catch (err: unknown) {
      const message =
        (err as ApiError)?.error?.message || "Failed to add zone.";
      toast.error(message);
    }
  };

  const handleAddPricing = async () => {
    if (!pricingZoneId) {
      toast.error("Please select a zone");
      return;
    }
    if (!pricingName.trim()) {
      toast.error("Pricing name is required");
      return;
    }
    if (!pricingPrice || parseFloat(pricingPrice) < 0) {
      toast.error("Please enter a valid price");
      return;
    }

    try {
      await createPricingMutation.mutateAsync({
        zone_id: pricingZoneId,
        name: pricingName.trim(),
        price: parseFloat(pricingPrice),
        currency: pricingCurrency || undefined,
      });
      toast.success("Pricing added!");
      setPricingName("");
      setPricingPrice("");
    } catch (err: unknown) {
      const message =
        (err as ApiError)?.error?.message || "Failed to add pricing.";
      toast.error(message);
    }
  };

  const getPricingForZone = (zoneId: string) =>
    pricing.filter((p) => p.zone_id === zoneId);

  return (
    <div className={cn("space-y-6", className)}>
      {/* Add Zone Form */}
      <Card>
        <CardHeader>
          <CardTitle>Add Seating Zone</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div className="space-y-2">
              <Label>Zone name *</Label>
              <Input
                placeholder="e.g. VIP Section A"
                value={zoneName}
                onChange={(e) => setZoneName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label>Type</Label>
              <Select value={zoneType} onValueChange={setZoneType}>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {ZONE_TYPES.map((t) => (
                    <SelectItem key={t.value} value={t.value}>
                      {t.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Seat count</Label>
              <Input
                type="number"
                placeholder="e.g. 500"
                value={seatCount}
                onChange={(e) => setSeatCount(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label>Color</Label>
              <div className="flex gap-1.5 pt-1">
                {ZONE_COLORS.map((c) => (
                  <button
                    key={c}
                    type="button"
                    className={cn(
                      "h-7 w-7 rounded-full border-2 transition-transform",
                      zoneColor === c
                        ? "border-foreground scale-110"
                        : "border-transparent"
                    )}
                    style={{ backgroundColor: c }}
                    onClick={() => setZoneColor(c)}
                  />
                ))}
              </div>
            </div>
          </div>
          <Button
            className="mt-4"
            size="sm"
            onClick={handleAddZone}
            disabled={createZoneMutation.isPending}
          >
            <Plus className="mr-1.5 h-4 w-4" />
            {createZoneMutation.isPending ? "Adding..." : "Add Zone"}
          </Button>
        </CardContent>
      </Card>

      {/* Zones List */}
      {zones.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Zones ({zones.length})</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {zones.map((zone) => (
                <div
                  key={zone.id}
                  className="flex items-center gap-3 rounded-lg border p-3"
                >
                  <div
                    className="h-4 w-4 rounded-full shrink-0"
                    style={{ backgroundColor: zone.color || "#3B82F6" }}
                  />
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm">{zone.name}</p>
                    <p className="text-xs text-muted-foreground">
                      {zone.zone_type} &middot;{" "}
                      {zone.seat_count ?? "N/A"} seats
                    </p>
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {getPricingForZone(zone.id).length > 0 ? (
                      getPricingForZone(zone.id).map((p) => (
                        <span key={p.id} className="block">
                          {p.name}: {p.currency} {p.price}
                        </span>
                      ))
                    ) : (
                      <span className="italic">No pricing</span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Add Pricing Form */}
      {zones.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Add Pricing</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div className="space-y-2">
                <Label>Zone *</Label>
                <Select value={pricingZoneId} onValueChange={setPricingZoneId}>
                  <SelectTrigger className="w-full">
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
                <Label>Pricing name *</Label>
                <Input
                  placeholder="e.g. Early Bird"
                  value={pricingName}
                  onChange={(e) => setPricingName(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label>Price *</Label>
                <Input
                  type="number"
                  step="0.01"
                  placeholder="e.g. 99.00"
                  value={pricingPrice}
                  onChange={(e) => setPricingPrice(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label>Currency</Label>
                <Select value={pricingCurrency} onValueChange={setPricingCurrency}>
                  <SelectTrigger className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="USD">USD</SelectItem>
                    <SelectItem value="VND">VND</SelectItem>
                    <SelectItem value="EUR">EUR</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <Button
              className="mt-4"
              size="sm"
              onClick={handleAddPricing}
              disabled={createPricingMutation.isPending}
            >
              <Plus className="mr-1.5 h-4 w-4" />
              {createPricingMutation.isPending ? "Adding..." : "Add Pricing"}
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Navigation */}
      <div className="flex justify-end gap-3">
        <Button variant="outline" onClick={onSkip}>
          Skip
        </Button>
        <Button onClick={onNext}>Next</Button>
      </div>
    </div>
  );
}
