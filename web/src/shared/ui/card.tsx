import type { HTMLAttributes } from "react";
import { cn } from "../lib/cn";

/** Tarjeta base — superficie card + hairline (jerarquía por luz, no por cajas pesadas). */
export function Card({ className, ...rest }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("rounded-lg border border-border bg-card text-card-foreground", className)}
      {...rest}
    />
  );
}
