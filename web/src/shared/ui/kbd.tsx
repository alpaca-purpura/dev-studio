import type { HTMLAttributes } from "react";
import { cn } from "../lib/cn";

/** Tecla/atajo — mono uppercase, hairline (firma técnica PRENTER). */
export function Kbd({ className, ...rest }: HTMLAttributes<HTMLElement>) {
  return (
    <kbd
      className={cn(
        "inline-flex items-center rounded-sm border border-border bg-secondary px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-wider text-muted-foreground",
        className,
      )}
      {...rest}
    />
  );
}
