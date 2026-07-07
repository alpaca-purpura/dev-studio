import { cn } from "../lib/cn";
import type { SessionStatus } from "../api/types";

/** Punto de estado de una sesión (streaming/idle) — copiado de harness-studio `Pip`. */
export function Pip({ status, className }: { status: SessionStatus; className?: string }) {
  return <span className={cn("pip", `pip-${status}`, className)} aria-hidden />;
}
