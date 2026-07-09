import { useState } from "react";
import type { ToolCall } from "../api/types";
import { cn } from "../lib/cn";
import { DiffStat } from "./diff-stat";

/** Tarjeta de una herramienta que el agente invocó (R1 · hueco #1 de conductor.go).
 *  Deriva el arg legible y el diff del input crudo; el cuerpo colapsa el output largo.
 *  Estilo making-of: `⎿ Nombre arg · estado`, teal único, JetBrains Mono. */
export function ToolCard({ call }: { call: ToolCall }) {
  const view = deriveView(call);
  const collapsible = !!view.body && view.body.split("\n").length > 6;
  const [open, setOpen] = useState(!collapsible);

  return (
    <div className="rounded-md border border-border bg-card/60 font-mono text-xs">
      <button
        type="button"
        onClick={() => collapsible && setOpen((o) => !o)}
        className={cn(
          "flex w-full items-center gap-2 px-3 py-1.5 text-left",
          collapsible && "cursor-pointer hover:bg-secondary/40",
        )}
      >
        <span className="text-primary">⎿</span>
        <span className="font-semibold text-primary">{call.name}</span>
        {view.arg && <span className="min-w-0 flex-1 truncate text-muted-foreground">{view.arg}</span>}
        {!view.arg && <span className="flex-1" />}
        {view.diff && <DiffStat add={view.diff.added} del={view.diff.removed} />}
        <StatusMark status={call.status} />
        {collapsible && (
          <span className="text-muted-foreground/60">{open ? "▾" : "▸"}</span>
        )}
      </button>

      {open && view.diff && (
        <pre className="overflow-x-auto border-t border-border px-3 py-2 leading-relaxed">
          {view.diff.lines.map((l, i) => (
            <div
              key={i}
              className={cn(
                "whitespace-pre",
                l.sign === "+" && "text-ok",
                l.sign === "-" && "text-crit",
                l.sign === " " && "text-muted-foreground/70",
              )}
            >
              {l.sign}
              {l.text}
            </div>
          ))}
          {view.diff.truncated > 0 && (
            <div className="text-muted-foreground/50">… +{view.diff.truncated} líneas más</div>
          )}
        </pre>
      )}

      {open && !view.diff && view.body && (
        <pre
          className={cn(
            "overflow-x-auto whitespace-pre-wrap border-t border-border px-3 py-2 leading-relaxed",
            call.status === "error" ? "text-crit" : "text-muted-foreground",
          )}
        >
          {view.body}
        </pre>
      )}
    </div>
  );
}

function StatusMark({ status }: { status: ToolCall["status"] }) {
  if (status === "running")
    return (
      <span className="inline-block size-2 rounded-full bg-primary motion-safe:animate-pulse" title="corriendo…" />
    );
  if (status === "error") return <span className="text-crit" title="error">✕</span>;
  return <span className="text-ok" title="ok">✓</span>;
}

interface DiffView {
  added: number;
  removed: number;
  lines: { sign: "+" | "-" | " "; text: string }[];
  truncated: number;
}
interface ToolView {
  arg: string;
  body?: string;
  diff?: DiffView;
}

const MAX_DIFF_LINES = 14;

/** deriveView traduce el input crudo (JSON) + output a lo que la card muestra, por herramienta. */
function deriveView(call: ToolCall): ToolView {
  const input = safeParse(call.input);
  const name = call.name;

  if (name === "Edit" || name === "MultiEdit") {
    const arg = shortPath(str(input.file_path));
    return { arg, diff: buildDiff(str(input.old_string), str(input.new_string)) };
  }
  if (name === "Write") {
    const arg = shortPath(str(input.file_path));
    return { arg, diff: buildDiff("", str(input.content)) };
  }
  if (name === "Bash") {
    return { arg: str(input.command), body: call.output };
  }
  if (name === "Read" || name === "Glob") {
    return { arg: shortPath(str(input.file_path) || str(input.path) || str(input.pattern)), body: call.output };
  }
  if (name === "Grep") {
    return { arg: str(input.pattern), body: call.output };
  }
  // desconocida: mostrar el input compacto como arg
  const compact = call.input && call.input !== "{}" ? call.input : "";
  return { arg: compact.length > 80 ? compact.slice(0, 80) + "…" : compact, body: call.output };
}

function buildDiff(oldStr: string, newStr: string): DiffView {
  const oldLines = oldStr ? oldStr.split("\n") : [];
  const newLines = newStr ? newStr.split("\n") : [];
  const lines: DiffView["lines"] = [];
  for (const t of oldLines) lines.push({ sign: "-", text: t });
  for (const t of newLines) lines.push({ sign: "+", text: t });
  const truncated = Math.max(0, lines.length - MAX_DIFF_LINES);
  return { added: newLines.length, removed: oldLines.length, lines: lines.slice(0, MAX_DIFF_LINES), truncated };
}

function safeParse(s: string): Record<string, unknown> {
  try {
    const v = JSON.parse(s || "{}");
    return v && typeof v === "object" ? (v as Record<string, unknown>) : {};
  } catch {
    return {};
  }
}
function str(v: unknown): string {
  return typeof v === "string" ? v : "";
}
/** deja ~2 segmentos finales de una ruta larga para que quepa en el header. */
function shortPath(p: string): string {
  if (!p) return "";
  const parts = p.split("/");
  return parts.length > 3 ? "…/" + parts.slice(-2).join("/") : p;
}
