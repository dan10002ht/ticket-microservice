import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

interface AvatarWithNameProps {
  name: string;
  email?: string;
  imageUrl?: string;
  size?: "sm" | "md" | "lg";
  className?: string;
}

const sizeClasses = {
  sm: { avatar: "h-8 w-8", name: "text-sm", email: "text-xs" },
  md: { avatar: "h-10 w-10", name: "text-sm", email: "text-xs" },
  lg: { avatar: "h-12 w-12", name: "text-base", email: "text-sm" },
};

function getInitials(name: string): string {
  return name
    .split(" ")
    .map((part) => part[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);
}

export function AvatarWithName({
  name,
  email,
  imageUrl,
  size = "md",
  className,
}: AvatarWithNameProps) {
  const s = sizeClasses[size];

  return (
    <div className={cn("flex items-center gap-3", className)}>
      <Avatar className={s.avatar}>
        <AvatarImage src={imageUrl} alt={name} />
        <AvatarFallback>{getInitials(name)}</AvatarFallback>
      </Avatar>
      <div className="flex flex-col">
        <span className={cn("font-medium leading-none", s.name)}>{name}</span>
        {email && (
          <span className={cn("text-muted-foreground", s.email)}>{email}</span>
        )}
      </div>
    </div>
  );
}
