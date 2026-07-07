import type { ButtonHTMLAttributes } from "react";
import { cn } from "../lib/cn";

type Variant = "primary" | "outline" | "ghost";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
}

/** Botón PRENTER. `primary` lleva el glow de marca (único elemento con glow — RN-6). */
export function Button({ variant = "primary", className, ...rest }: ButtonProps) {
  return (
    <button
      className={cn(
        "inline-flex cursor-pointer items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-45",
        variant === "primary" &&
          "bg-primary text-primary-foreground shadow-[var(--shadow-glow)] hover:brightness-110",
        variant === "outline" &&
          "border border-border bg-transparent text-foreground hover:bg-secondary",
        variant === "ghost" && "bg-transparent text-muted-foreground hover:bg-secondary hover:text-foreground",
        className,
      )}
      {...rest}
    />
  );
}
