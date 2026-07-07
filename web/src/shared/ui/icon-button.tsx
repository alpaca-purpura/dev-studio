import type { ButtonHTMLAttributes, ReactNode } from "react";
import { cn } from "../lib/cn";

export interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
  /** accesibilidad obligatoria: los icon-button no tienen texto visible */
  "aria-label": string;
}

/** Botón circular de solo icono (composer, rails, headers). */
export function IconButton({ className, ...rest }: IconButtonProps) {
  return (
    <button
      className={cn(
        "inline-flex size-8 shrink-0 cursor-pointer items-center justify-center rounded-full border border-transparent text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground disabled:cursor-not-allowed disabled:opacity-40 [&>svg]:size-4",
        className,
      )}
      {...rest}
    />
  );
}
