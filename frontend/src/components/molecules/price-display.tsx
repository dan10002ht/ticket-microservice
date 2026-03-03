import { cn } from "@/lib/utils";

interface PriceDisplayProps {
  amount: number;
  currency?: string;
  locale?: string;
  size?: "sm" | "md" | "lg";
  className?: string;
}

const sizeClasses = {
  sm: "text-sm",
  md: "text-base",
  lg: "text-2xl font-bold",
};

export function PriceDisplay({
  amount,
  currency = "VND",
  locale = "vi-VN",
  size = "md",
  className,
}: PriceDisplayProps) {
  const formatted = new Intl.NumberFormat(locale, {
    style: "currency",
    currency,
    minimumFractionDigits: 0,
  }).format(amount);

  return (
    <span className={cn(sizeClasses[size], "tabular-nums", className)}>
      {formatted}
    </span>
  );
}
