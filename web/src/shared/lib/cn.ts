import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/** Une clases Tailwind resolviendo conflictos (patrón shadcn, copiado de harness-studio). */
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}
