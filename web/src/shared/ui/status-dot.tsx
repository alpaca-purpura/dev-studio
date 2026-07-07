import { cn } from "../lib/cn";

export type WorkspaceState = "idle" | "ready" | "conflict" | "archived" | "streaming";

const colors: Record<WorkspaceState, string> = {
  idle: "bg-muted-foreground/50",
  ready: "bg-ok",
  conflict: "bg-crit",
  archived: "bg-muted-foreground/30",
  streaming: "bg-primary animate-pulse",
};

/** Punto de estado de un workspace/sesión (rail izquierdo). */
export function StatusDot({ state, className }: { state: WorkspaceState; className?: string }) {
  return <span className={cn("inline-block size-2 shrink-0 rounded-full", colors[state], className)} />;
}

/** Etiqueta humana del estado — español neutro (RN-10). */
export const statusLabel: Record<WorkspaceState, string> = {
  idle: "Sin cambios",
  ready: "Cambios listos",
  conflict: "Conflictos",
  archived: "Archivado",
  streaming: "Trabajando…",
};
