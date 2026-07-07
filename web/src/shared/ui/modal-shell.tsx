import { useEffect, type ReactNode } from "react";
import { cn } from "../lib/cn";
import { IconButton } from "./icon-button";

/** Modal centrado sobre backdrop (composer, quick-look, revisión). Esc y clic-afuera cierran. */
export function ModalShell({
  open,
  onClose,
  title,
  subtitle,
  icon,
  children,
  footer,
  wide,
}: {
  open: boolean;
  onClose: () => void;
  title: string;
  subtitle?: string;
  icon?: ReactNode;
  children: ReactNode;
  footer?: ReactNode;
  wide?: boolean;
}) {
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => e.key === "Escape" && onClose();
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  if (!open) return null;
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-6"
      onClick={(e) => e.target === e.currentTarget && onClose()}
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-label={title}
        className={cn(
          "flex max-h-[85vh] w-full flex-col overflow-hidden rounded-lg border border-border bg-popover shadow-2xl",
          wide ? "max-w-4xl" : "max-w-lg",
        )}
      >
        <div className="flex items-start gap-3 border-b border-border p-4">
          {icon && (
            <span className="flex size-9 shrink-0 items-center justify-center rounded-md bg-secondary text-base">
              {icon}
            </span>
          )}
          <div className="min-w-0 flex-1">
            <h4 className="font-display text-base font-semibold">{title}</h4>
            {subtitle && <p className="mt-0.5 text-xs text-muted-foreground">{subtitle}</p>}
          </div>
          <IconButton aria-label="Cerrar" onClick={onClose}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M6 6l12 12M18 6L6 18" />
            </svg>
          </IconButton>
        </div>
        <div className="min-h-0 flex-1 overflow-auto p-4">{children}</div>
        {footer && <div className="flex justify-end gap-2 border-t border-border p-3">{footer}</div>}
      </div>
    </div>
  );
}
