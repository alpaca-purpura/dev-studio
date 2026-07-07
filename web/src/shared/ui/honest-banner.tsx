import type { ReactNode } from "react";
import { cn } from "../lib/cn";

/** Banner de honestidad de superficie (RN-7): lo que no existe se dice, con destino. */
export function HonestBanner({ children, className }: { children: ReactNode; className?: string }) {
  return (
    <div
      role="note"
      className={cn(
        "flex items-start gap-2 rounded-md border border-warn/30 bg-warn-soft px-3 py-2 text-xs text-warn",
        className,
      )}
    >
      <span aria-hidden>◔</span>
      <span className="text-foreground/80">{children}</span>
    </div>
  );
}
