import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getTotalPages(total: number, limit: number): number {
  return Math.max(1, Math.ceil(total / limit))
}
