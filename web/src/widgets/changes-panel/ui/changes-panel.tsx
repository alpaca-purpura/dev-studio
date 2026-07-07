import { useCallback, useEffect, useState } from "react";
import { api } from "../../../shared/api/client";
import type { GitDiff, GitLogEntry, GitStatus } from "../../../shared/api/types";
import { useSessions } from "../../../shared/store/sessions-store";
import { useRepos } from "../../../shared/store/repos-store";
import { useUi } from "../../../shared/store/ui-store";
import {
  Avatar,
  Button,
  Chip,
  DiffStat,
  FilterChip,
  ModalShell,
  Tabs,
  Textarea,
} from "../../../shared/ui";
import { cn } from "../../../shared/lib/cn";

const stateLabel: Record<string, string> = { M: "modificado", A: "nuevo", D: "eliminado", R: "renombrado" };

function DiffLines({ text, sideBySide }: { text: string; sideBySide?: boolean }) {
  return (
    <pre className={cn("overflow-x-auto p-3 font-mono text-xs leading-relaxed", sideBySide && "h-full")}>
      {text.split("\n").map((l, i) => (
        <div
          key={i}
          className={cn(
            "px-1",
            l.startsWith("+") && !l.startsWith("+++") && "bg-ok-soft text-ok",
            l.startsWith("-") && !l.startsWith("---") && "bg-crit-soft text-crit",
          )}
        >
          {l || " "}
        </div>
      ))}
    </pre>
  );
}

export function ChangesPanel() {
  const session = useSessions((s) => s.sessions.find((x) => x.id === s.activeId));
  const refreshSessionStatus = useRepos((s) => s.refreshSessionStatus);
  const requestClose = useUi((s) => s.requestClose);
  const [tab, setTab] = useState("cambios");
  const [pane, setPane] = useState<"pendientes" | "historial">("pendientes");
  const [status, setStatus] = useState<GitStatus | null>(null);
  const [statusErr, setStatusErr] = useState<string | null>(null);
  const [log, setLog] = useState<GitLogEntry[]>([]);
  const [checked, setChecked] = useState<Record<string, boolean>>({});
  const [mensaje, setMensaje] = useState("");
  const [busy, setBusy] = useState(false);
  const [quick, setQuick] = useState<GitDiff | null>(null);
  const [review, setReview] = useState<{ diffs: GitDiff[]; i: number; done: boolean } | null>(null);
  const [logSearch, setLogSearch] = useState("");

  const sessionId = session?.id;
  const streaming = session?.status === "streaming";

  const refresh = useCallback(async () => {
    if (!sessionId) return;
    try {
      const st = await api.gitStatus(sessionId);
      setStatus(st);
      setStatusErr(null);
      setChecked((prev) => {
        const next: Record<string, boolean> = {};
        for (const f of st.files) next[f.path] = prev[f.path] ?? true;
        return next;
      });
    } catch (e) {
      setStatus(null);
      setStatusErr(e instanceof Error ? e.message : "No se pudo leer el estado git");
    }
    try {
      setLog(await api.gitLog(sessionId));
    } catch {
      setLog([]);
    }
  }, [sessionId]);

  // refresh al cambiar de sesión y al terminar cada turno del agente (spec §5: on-demand, sin watcher)
  useEffect(() => {
    void refresh();
  }, [refresh, streaming]);

  if (!session) return null;

  const files = status?.files ?? [];
  const selected = files.filter((f) => checked[f.path]).map((f) => f.path);

  const openQuick = async (path: string) => setQuick(await api.gitDiff(session.id, path));

  const openReview = async () => {
    const diffs = await Promise.all(files.map((f) => api.gitDiff(session.id, f.path)));
    setReview({ diffs, i: 0, done: false });
  };

  const doCommit = async (cerrar: boolean) => {
    if (selected.length === 0 || busy) return;
    setBusy(true);
    try {
      await api.gitCommit(session.id, selected, mensaje.trim());
      setMensaje("");
      await refresh();
      void refreshSessionStatus(session.id);
      // cerrar pasa por el modal conservar/borrar del workspace (spec workspace-aislado §2.2)
      if (cerrar) requestClose(session.id);
    } catch (e) {
      setStatusErr(e instanceof Error ? e.message : "No se pudo commitear");
    } finally {
      setBusy(false);
    }
  };

  return (
    <aside className="flex w-80 shrink-0 flex-col border-l border-border bg-sidebar">
      <Tabs
        items={[
          { id: "archivos", label: "Archivos" },
          { id: "cambios", label: "Cambios" },
          { id: "pruebas", label: "Pruebas" },
        ]}
        active={tab}
        onSelect={setTab}
        className="px-2"
      />

      {tab === "archivos" && (
        <div className="min-h-0 flex-1 overflow-y-auto p-3">
          {files.length === 0 && <p className="text-xs text-muted-foreground">Sin archivos tocados en este workspace.</p>}
          {files.map((f) => (
            <button
              key={f.path}
              onClick={() => void openQuick(f.path)}
              className="flex w-full cursor-pointer items-center gap-2 rounded px-1.5 py-1 text-left hover:bg-sidebar-accent"
            >
              <span className="font-mono text-[10px] text-muted-foreground">{f.state}</span>
              <span className="truncate font-mono text-xs">{f.path}</span>
            </button>
          ))}
        </div>
      )}

      {tab === "cambios" && (
        <div className="flex min-h-0 flex-1 flex-col">
          <div className="flex gap-1 px-3 pt-2">
            {(["pendientes", "historial"] as const).map((p) => (
              <FilterChip key={p} active={pane === p} onClick={() => setPane(p)}>
                {p === "pendientes" ? "Cambios pendientes" : "Historial"}
              </FilterChip>
            ))}
          </div>

          {pane === "pendientes" && (
            <div className="flex min-h-0 flex-1 flex-col">
              <div className="flex items-center gap-2 px-3 py-2">
                <Chip className="max-w-40 truncate">{status?.branch || "—"}</Chip>
                <DiffStat add={status?.add ?? 0} del={status?.del ?? 0} />
                <span className="flex-1" />
                <Chip>{files.length} cambios</Chip>
              </div>

              {/* sync: deshabilitados por RN-4 — la app no toca el remoto */}
              <div className="flex items-center gap-1 px-3 pb-2">
                {["Fetch", "Pull", "Push"].map((op) => (
                  <button
                    key={op}
                    disabled
                    title={`${op} — la app no toca el remoto (boundary git-solo-lectura-y-commit); se decidirá con el PR-flow (PB-14)`}
                    className="cursor-not-allowed rounded border border-border px-2 py-0.5 text-[10px] text-muted-foreground/40"
                  >
                    {op}
                  </button>
                ))}
                <span className="flex-1" />
                <button
                  onClick={() => void refresh()}
                  title="Actualizar estado"
                  className="cursor-pointer rounded border border-border px-2 py-0.5 text-[10px] text-muted-foreground hover:text-foreground"
                >
                  ↻ Actualizar
                </button>
              </div>

              <div className="px-3 pb-2">
                <Button onClick={() => void openReview()} disabled={files.length === 0} className="w-full">
                  ⚡ Revisar cambios
                </Button>
              </div>

              {/* chips por rol: v1 solo «Todos» — la atribución por rol necesita PB-02 (RN-7) */}
              <div className="flex gap-1 px-3 pb-2">
                <FilterChip active>Todos {files.length}</FilterChip>
                <span
                  className="self-center text-[9px] text-muted-foreground/60"
                  title="Filtro por rol: necesita workspace aislado por sesión (PB-02)"
                >
                  por rol → PB-02
                </span>
              </div>

              <div className="min-h-0 flex-1 overflow-y-auto px-3">
                {statusErr && <p className="py-2 text-xs text-crit">{statusErr}</p>}
                {!statusErr && files.length === 0 && (
                  <p className="py-2 text-xs text-muted-foreground">Working tree limpio.</p>
                )}
                {files.map((f) => (
                  <div key={f.path} className="flex items-center gap-2 rounded px-1 py-1 hover:bg-sidebar-accent">
                    <input
                      type="checkbox"
                      checked={checked[f.path] ?? false}
                      onChange={(e) => setChecked((c) => ({ ...c, [f.path]: e.target.checked }))}
                      aria-label={`Incluir ${f.path} en el commit`}
                      className="size-3.5 cursor-pointer accent-[var(--primary)]"
                    />
                    <button
                      onClick={() => void openQuick(f.path)}
                      title={f.path}
                      className="min-w-0 flex-1 cursor-pointer truncate text-left font-mono text-xs hover:text-primary"
                    >
                      {f.path}
                    </button>
                    <span className="font-mono text-[9px] uppercase text-muted-foreground">{stateLabel[f.state] ?? f.state}</span>
                  </div>
                ))}
              </div>

              <div className="space-y-2 border-t border-border p-3">
                <Textarea
                  rows={2}
                  value={mensaje}
                  onChange={(e) => setMensaje(e.target.value)}
                  placeholder="Mensaje de commit…"
                />
                <label
                  className="flex cursor-not-allowed items-start gap-2 text-[10px] text-muted-foreground/50"
                  title="Adjuntar la conversación llega con PB-10"
                >
                  <input type="checkbox" disabled className="mt-0.5" />
                  Adjuntar la conversación del rol al commit — todavía no (PB-10).
                </label>
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    disabled={selected.length === 0 || busy}
                    onClick={() => void doCommit(false)}
                    className="flex-1"
                  >
                    Confirmar cambios
                  </Button>
                  <Button
                    disabled={selected.length === 0 || busy}
                    onClick={() => void doCommit(true)}
                    title="Commit + cerrar la sesión"
                    className="flex-1"
                  >
                    Confirmar y cerrar
                  </Button>
                </div>
                <p className="text-[9px] leading-snug text-muted-foreground/70">
                  Se hará commit de {selected.length}/{files.length} archivos — solo los tildados (RN-5).
                </p>
              </div>
            </div>
          )}

          {pane === "historial" && (
            <div className="flex min-h-0 flex-1 flex-col">
              <div className="p-3 pb-1">
                <input
                  value={logSearch}
                  onChange={(e) => setLogSearch(e.target.value)}
                  placeholder="🔍 Buscar mensajes…"
                  className="w-full rounded-md border border-input bg-card px-2 py-1 text-xs outline-none focus:border-primary"
                />
              </div>
              <div className="min-h-0 flex-1 overflow-y-auto px-3 pb-3">
                {log
                  .filter(
                    (e) =>
                      !logSearch ||
                      e.mensaje.toLowerCase().includes(logSearch.toLowerCase()) ||
                      e.autor.toLowerCase().includes(logSearch.toLowerCase()),
                  )
                  .map((e) => (
                    <div key={e.sha} className="flex gap-2 border-b border-border/50 py-2">
                      <Avatar label={e.autor} className="size-6 text-[8px]" />
                      <div className="min-w-0 flex-1">
                        <div className="truncate text-xs font-semibold">{e.mensaje}</div>
                        <div className="mt-0.5 flex items-center gap-2 text-[9px] text-muted-foreground">
                          <span className="truncate">{e.autor}</span>
                          <code className="font-mono">{e.sha.slice(0, 7)}</code>
                          <span className="ml-auto shrink-0">{new Date(e.fecha).toLocaleDateString("es", { day: "numeric", month: "short" })}</span>
                        </div>
                      </div>
                    </div>
                  ))}
                {log.length === 0 && <p className="py-2 text-xs text-muted-foreground">Sin historial (repo sin commits).</p>}
              </div>
            </div>
          )}
        </div>
      )}

      {tab === "pruebas" && (
        <div className="p-4 text-xs leading-relaxed text-muted-foreground">
          Superficies de prueba con navegador y MCP — todavía no disponible en esta rebanada (RN-7).
        </div>
      )}

      {/* quick-look (diff unificado) */}
      <ModalShell
        open={quick !== null}
        onClose={() => setQuick(null)}
        title={quick?.path.split("/").pop() ?? ""}
        subtitle={quick?.path}
        icon="📄"
        wide
      >
        {quick?.binary ? (
          <p className="text-sm text-muted-foreground">Archivo binario — sin vista previa.</p>
        ) : (
          <DiffLines text={quick?.unified || quick?.modified || ""} />
        )}
      </ModalShell>

      {/* revisión side-by-side */}
      <ModalShell
        open={review !== null}
        onClose={() => setReview(null)}
        title={review?.done ? "Listo para hacer commit" : (review ? `${review.diffs[review.i]?.path} · ${review.i + 1}/${review.diffs.length}` : "")}
        icon="⚡"
        wide
        footer={
          review && !review.done ? (
            <>
              <Button
                variant="outline"
                disabled={review.i === 0}
                onClick={() => setReview((r) => (r ? { ...r, i: r.i - 1 } : r))}
              >
                ← Anterior
              </Button>
              <Button
                variant="outline"
                disabled
                title="Generar el mensaje con IA — todavía no en esta rebanada (RN-7)"
              >
                ✨ Generar con IA
              </Button>
              <Button
                onClick={() =>
                  setReview((r) =>
                    r ? (r.i === r.diffs.length - 1 ? { ...r, done: true } : { ...r, i: r.i + 1 }) : r,
                  )
                }
              >
                {review.i === review.diffs.length - 1 ? "Finalizar" : "Siguiente →"}
              </Button>
            </>
          ) : undefined
        }
      >
        {review && !review.done && review.diffs[review.i] && (
          <div className="grid h-[50vh] grid-cols-2 gap-2">
            <div className="flex min-h-0 flex-col overflow-hidden rounded-md border border-border">
              <div className="border-b border-border bg-secondary px-2 py-1 font-mono text-[10px] uppercase tracking-wider text-crit">
                Original
              </div>
              <div className="min-h-0 flex-1 overflow-auto">
                {review.diffs[review.i].original ? (
                  <DiffLines text={review.diffs[review.i].original} sideBySide />
                ) : (
                  <p className="p-3 text-xs italic text-muted-foreground">— archivo nuevo —</p>
                )}
              </div>
            </div>
            <div className="flex min-h-0 flex-col overflow-hidden rounded-md border border-border">
              <div className="border-b border-border bg-secondary px-2 py-1 font-mono text-[10px] uppercase tracking-wider text-ok">
                Modificado
              </div>
              <div className="min-h-0 flex-1 overflow-auto">
                <DiffLines text={review.diffs[review.i].modified} sideBySide />
              </div>
            </div>
          </div>
        )}
        {review?.done && (
          <div className="flex flex-col items-center gap-3 py-6 text-center">
            <div className="flex size-12 items-center justify-center rounded-full bg-accent-soft text-xl text-primary shadow-[var(--shadow-glow)]">
              ✓
            </div>
            <p className="text-sm font-semibold">Revisión completa</p>
            <p className="max-w-sm text-xs text-muted-foreground">
              Todos los archivos fueron revisados. Cerrá este panel y escribí el mensaje de commit
              en «Cambios pendientes» para confirmar los tildados.
            </p>
            <Button onClick={() => setReview(null)}>Ir a confirmar</Button>
          </div>
        )}
      </ModalShell>
    </aside>
  );
}
