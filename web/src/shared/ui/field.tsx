import type { ComponentProps } from "react";
import { cn } from "../lib/cn";

const base =
  "w-full rounded-md border border-input bg-card px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground/60 focus-visible:outline-2 focus-visible:outline-ring disabled:cursor-not-allowed disabled:opacity-45";

export function Input({ className, ...rest }: ComponentProps<"input">) {
  return <input className={cn(base, className)} {...rest} />;
}

export function Textarea({ className, ...rest }: ComponentProps<"textarea">) {
  return <textarea className={cn(base, "resize-none", className)} {...rest} />;
}

export function Select({ className, children, ...rest }: ComponentProps<"select">) {
  return (
    <select className={cn(base, "cursor-pointer", className)} {...rest}>
      {children}
    </select>
  );
}

/** Label de campo — mono uppercase tracking ancho (firma técnica PRENTER). */
export function FieldLabel({ className, ...rest }: React.HTMLAttributes<HTMLLabelElement>) {
  return (
    <label
      className={cn("mb-1 block font-mono text-[10px] uppercase tracking-[0.14em] text-muted-foreground", className)}
      {...rest}
    />
  );
}
