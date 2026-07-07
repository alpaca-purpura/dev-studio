import type { ReactNode } from "react";
import { cn } from "../lib/cn";
import { IconButton } from "./icon-button";

/** Overlay de panel (§0.5, RN-3): tapa SOLO el área de contenido — los rails quedan
 *  visibles y clickeables porque este overlay se posiciona `absolute` DENTRO del área
 *  de sesión, nunca `fixed` sobre la app entera. */
export function PanelOverlay({
  open,
  onClose,
  title,
  subtitle,
  children,
  className,
}: {
  open: boolean;
  onClose: () => void;
  title: string;
  subtitle?: ReactNode;
  children: ReactNode;
  className?: string;
}) {
  if (!open) return null;
  return (
    <div
      className={cn("absolute inset-0 z-30 flex flex-col overflow-hidden bg-background", className)}
      role="region"
      aria-label={title}
    >
      <div className="flex items-start gap-3 border-b border-border px-6 py-4">
        <div className="min-w-0 flex-1">
          <h2 className="font-display text-xl font-semibold tracking-tight">{title}</h2>
          {subtitle && <p className="mt-1 max-w-3xl text-xs leading-relaxed text-muted-foreground">{subtitle}</p>}
        </div>
        <IconButton aria-label={`Cerrar ${title}`} onClick={onClose}>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M6 6l12 12M18 6L6 18" />
          </svg>
        </IconButton>
      </div>
      <div className="min-h-0 flex-1 overflow-auto">{children}</div>
    </div>
  );
}
