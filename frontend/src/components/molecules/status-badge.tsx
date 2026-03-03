import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

const statusColorMap: Record<string, string> = {
  active: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
  confirmed: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
  completed: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
  pending: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
  reserved: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
  cancelled: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  failed: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  refunded: "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300",
  draft: "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300",
};

interface StatusBadgeProps {
  status: string;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const colorClass = statusColorMap[status.toLowerCase()] ?? statusColorMap.draft;

  return (
    <Badge
      variant="outline"
      className={cn("capitalize border-0", colorClass, className)}
    >
      {status}
    </Badge>
  );
}
