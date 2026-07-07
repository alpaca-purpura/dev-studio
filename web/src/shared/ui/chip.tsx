import type { HTMLAttributes } from "react";
import { cn } from "../lib/cn";

/** Chip informativo (branch, contadores, metadatos) — mono, hairline, sin acento. */
export function Chip({ className, ...rest }: HTMLAttributes<HTMLSpanElement>) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 rounded-full border border-border bg-secondary px-2 py-0.5 font-mono text-xs text-muted-foreground",
        className,
      )}
      {...rest}
    />
  );
}
