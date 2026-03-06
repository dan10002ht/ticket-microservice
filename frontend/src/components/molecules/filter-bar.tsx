"use client";

import { X } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";

interface FilterSlot {
  key: string;
  placeholder: string;
  options: { value: string; label: string }[];
  value: string;
  onChange: (value: string) => void;
}

interface FilterBarProps {
  slots: FilterSlot[];
  className?: string;
}

export function FilterBar({ slots, className }: FilterBarProps) {
  const hasActive = slots.some((s) => s.value !== "");

  return (
    <div className={cn("flex flex-wrap items-center gap-2", className)}>
      {slots.map((slot) => (
        <Select
          key={slot.key}
          value={slot.value || undefined}
          onValueChange={slot.onChange}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder={slot.placeholder} />
          </SelectTrigger>
          <SelectContent>
            {slot.options.map((opt) => (
              <SelectItem key={opt.value} value={opt.value}>
                {opt.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ))}

      {hasActive && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => slots.forEach((s) => s.onChange(""))}
          className="text-muted-foreground"
        >
          <X className="mr-1 h-3 w-3" />
          Clear
        </Button>
      )}
    </div>
  );
}

export type { FilterSlot, FilterBarProps };
