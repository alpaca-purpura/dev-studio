import { useEffect, useMemo, useRef, useState } from "react";
import type { Repo, Session } from "../../../shared/api/types";
import { useRepos } from "../../../shared/store/repos-store";
import { useSessions } from "../../../shared/store/sessions-store";
import { useUi } from "../../../shared/store/ui-store";
import { cn } from "../../../shared/lib/cn";
import { Button, DiffStat, Input, Kbd, StatusDot, statusLabel, type WorkspaceState } from "../../../shared/ui";
import { UpdateFooter } from "./update-footer";

/** Pseudo-repo para sesiones F1 sin repo_id (migración, spec shell §6). */
const SIN_REPO: Repo = { id: "", nombre: "(sin repositorio)", ruta: "" };

function workspaceState(s: Session, dirty: boolean, conflict: boolean): WorkspaceState {
  if (s.status === "streaming") return "streaming";
  if (conflict) return "conflict";
  return dirty ? "ready" : "idle";
}

function WorkspaceItem({ session, index }: { session: Session; index: number }) {
  const activeId = useSessions((s) => s.activeId);
  const switchTo = useSessions((s) => s.switchTo);
  const requestClose = useUi((s) => s.requestClose);
  // PB-02: branch/±N del worktree DE ESTA SESIÓN (ya no el status agregado del repo)
  const st = useRepos((s) => s.sessionStatus[session.id]);
  const refreshSessionStatus = useRepos((s) => s.refreshSessionStatus);

  useEffect(() => {
    void refreshSessionStatus(session.id);
  }, [refreshSessionStatus, session.id]);

  const dirty = (st?.files.length ?? 0) > 0;
  const conflict = st?.files.some((f) => f.state === "U") ?? false;
  const state = workspaceState(session, dirty, conflict);
  const active = activeId === session.id;

  return (
    <div
      data-session={session.id}
      className={cn(
        "group relative w-full rounded-md border border-transparent transition-colors hover:bg-sidebar-accent",
        active && "border-border bg-sidebar-accent shadow-[inset_2px_0_0_var(--primary)]",
      )}
    >
      <button onClick={() => switchTo(session.id)} className="w-full cursor-pointer px-2 py-1.5 text-left">
        <div className="flex items-center gap-1.5">
          <span aria-hidden className="text-muted-foreground">⎇</span>
          <span className="min-w-0 flex-1 truncate font-mono text-xs font-semibold text-foreground">
            {session.nombre}
          </span>
          <DiffStat add={st?.add ?? 0} del={st?.del ?? 0} />
        </div>
        <div className="mt-0.5 flex items-center gap-1.5 pl-5 text-[10px] text-muted-foreground">
          <span className="truncate font-mono">{session.branch ?? st?.branch ?? "—"}</span>
          <span aria-hidden>·</span>
          <StatusDot state={state} />
          <span className="truncate">{statusLabel[state]}</span>
          {index < 9 && <Kbd className="ml-auto">⌘{index + 1}</Kbd>}
        </div>
      </button>
      <button
        aria-label={`Cerrar ${session.nombre}`}
        title="Cerrar sesión"
        onClick={(e) => {
          e.stopPropagation();
          requestClose(session.id);
        }}
        className="absolute right-1 top-1 hidden cursor-pointer rounded p-0.5 text-muted-foreground hover:bg-secondary hover:text-foreground group-hover:block"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="size-3">
          <path d="M6 6l12 12M18 6L6 18" />
        </svg>
      </button>
    </div>
  );
}

/** clave de grupo estable — el pseudo-repo «(sin repositorio)» también colapsa (UX estándar
 *  de tree views: todo grupo colapsa, el contador siempre visible; bug dogfooding DH-16.1) */
export const repoKey = (r: Repo) => r.id || "sin-repo";

function RepoBlock({ repo, sessions, startIndex }: { repo: Repo; sessions: Session[]; startIndex: number }) {
  const expanded = useRepos((s) => s.expanded[repoKey(repo)] ?? true);
  const toggleExpanded = useRepos((s) => s.toggleExpanded);
  const startPicker = useUi((s) => s.startPicker);
  const [menuOpen, setMenuOpen] = useState(false);
  const isReal = repo.id !== "";

  return (
    <section className="py-1">
      <button
        onClick={() => toggleExpanded(repoKey(repo))}
        aria-expanded={expanded}
        className="flex w-full cursor-pointer items-center gap-1.5 rounded-md px-2 py-1 text-left hover:bg-sidebar-accent"
      >
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          className={cn("size-3 text-muted-foreground transition-transform", !expanded && "-rotate-90")}
        >
          <polyline points="6 9 12 15 18 9" />
        </svg>
        <span className="min-w-0 flex-1 truncate text-sm font-bold text-foreground">{repo.nombre}</span>
        <span className="rounded-full bg-secondary px-1.5 font-mono text-[10px] text-muted-foreground">
          {sessions.length}
        </span>
      </button>

      {expanded && (
        <div className="space-y-0.5 pl-2">
          {isReal && (
            <div className="relative">
              <button
                onClick={() => startPicker(repo.id)}
                className="flex w-full cursor-pointer items-center gap-1.5 rounded-md px-2 py-1.5 text-left text-xs text-muted-foreground hover:bg-sidebar-accent hover:text-foreground"
              >
                <span className="text-primary">+</span>
                <span className="flex-1">Nuevo Workspace</span>
                <span
                  role="button"
                  tabIndex={0}
                  aria-label="Opciones de Nuevo Workspace"
                  onClick={(e) => {
                    e.stopPropagation();
                    setMenuOpen((v) => !v);
                  }}
                  className="rounded px-1 hover:bg-secondary"
                >
                  ⋯
                </span>
              </button>
              {menuOpen && (
                <div
                  className="absolute right-0 top-8 z-40 w-52 rounded-md border border-border bg-popover p-1 shadow-lg"
                  onMouseLeave={() => setMenuOpen(false)}
                >
                  {["Duplicar último workspace", "Nuevo desde branch…", "Nuevo worktree…"].map((op) => (
                    <button
                      key={op}
                      disabled
                      title="Se define con el workspace aislado (PB-02)"
                      className="block w-full cursor-not-allowed rounded px-2 py-1.5 text-left text-xs text-muted-foreground/50"
                    >
                      {op}
                    </button>
                  ))}
                  <div className="border-t border-border px-2 py-1 text-[10px] text-muted-foreground">
                    Se define con PB-02 (workspace aislado)
                  </div>
                </div>
              )}
            </div>
          )}
          {sessions.map((s, i) => (
            <WorkspaceItem key={s.id} session={s} index={startIndex + i} />
          ))}
        </div>
      )}
    </section>
  );
}

export function ReposRail() {
  const repos = useRepos((s) => s.repos);
  const addOpen = useRepos((s) => s.addOpen);
  const addError = useRepos((s) => s.addError);
  const setAddOpen = useRepos((s) => s.setAddOpen);
  const register = useRepos((s) => s.register);
  const sessions = useSessions((s) => s.sessions);
  const switchTo = useSessions((s) => s.switchTo);
  const inputRef = useRef<HTMLInputElement>(null);
  const [ruta, setRuta] = useState("");

  // orden plano de workspaces para ⌘1..⌘9 — mismo orden visual del rail
  const ordered = useMemo(() => {
    const byRepo = [...repos, SIN_REPO].map((r) => ({
      repo: r,
      sessions: sessions.filter((s) => (r.id === "" ? !s.repo_id : s.repo_id === r.id)),
    }));
    return byRepo.filter((b) => b.repo.id !== "" || b.sessions.length > 0);
  }, [repos, sessions]);

  const flat = useMemo(() => ordered.flatMap((b) => b.sessions), [ordered]);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (!(e.metaKey || e.ctrlKey)) return;
      const n = Number(e.key);
      if (n >= 1 && n <= 9 && flat[n - 1]) {
        e.preventDefault();
        switchTo(flat[n - 1].id);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [flat, switchTo]);

  // reveal: activar una sesión (clic/⌘N) expande su grupo si estaba colapsado
  const activeId = useSessions((s) => s.activeId);
  const ensureExpanded = useRepos((s) => s.ensureExpanded);
  useEffect(() => {
    if (!activeId) return;
    const sess = sessions.find((s) => s.id === activeId);
    if (sess) ensureExpanded(sess.repo_id || "sin-repo");
  }, [activeId, sessions, ensureExpanded]);

  const submitAdd = async () => {
    const r = ruta.trim();
    if (!r) return;
    const repo = await register(r);
    if (repo) setRuta("");
  };

  let acc = 0;

  return (
    <nav className="flex w-56 shrink-0 flex-col border-r border-sidebar-border bg-sidebar" aria-label="Repositorios y workspaces">
      <div className="p-2">
        <button
          title="Inicio — placeholder, sin destino aún (RN-7)"
          className="flex w-full cursor-default items-center gap-2 rounded-md px-2 py-1.5 text-sm font-bold text-sidebar-foreground"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className="size-4">
            <path d="M3 11l9-8 9 8" />
            <path d="M5 10v10a1 1 0 0 0 1 1h4v-6h4v6h4a1 1 0 0 0 1-1V10" />
          </svg>
          Inicio
        </button>
      </div>
      <div className="mx-2 border-t border-sidebar-border" />

      <div className="flex items-center justify-between px-3 py-2">
        <span className="font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">
          Repositorios
        </span>
        <button
          aria-label="Agregar repositorio"
          aria-expanded={addOpen}
          onClick={() => {
            setAddOpen(!addOpen);
            setTimeout(() => inputRef.current?.focus(), 0);
          }}
          className="cursor-pointer rounded p-1 text-muted-foreground hover:bg-sidebar-accent hover:text-foreground"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className="size-4">
            <path d="M13 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8z" />
            <path d="M13 3v5h5" />
            <path d="M12 13v5M9.5 15.5h5" />
          </svg>
        </button>
      </div>

      {addOpen && (
        <div className="space-y-1 px-2 pb-2">
          <div className="flex gap-1">
            <Input
              ref={inputRef}
              value={ruta}
              onChange={(e) => setRuta(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && void submitAdd()}
              placeholder="Ruta local del repositorio…"
              className="h-7 px-2 text-xs"
            />
            <Button onClick={() => void submitAdd()} className="h-7 px-2 text-xs">
              Agregar
            </Button>
          </div>
          {addError && <p className="px-1 text-[10px] leading-snug text-crit">{addError}</p>}
        </div>
      )}

      <div className="min-h-0 flex-1 overflow-y-auto px-1">
        {ordered.length === 0 && (
          <p className="px-3 py-4 text-xs leading-relaxed text-muted-foreground">
            Sin repositorios. Registrá el primero con el botón de arriba (ruta local con `.git`).
          </p>
        )}
        {ordered.map((b, i) => {
          const start = acc;
          acc += b.sessions.length;
          return (
            <div key={b.repo.id || "sin-repo"}>
              {i > 0 && <div className="mx-1 my-1 border-t border-sidebar-border" />}
              <RepoBlock repo={b.repo} sessions={b.sessions} startIndex={start} />
            </div>
          );
        })}
      </div>

      <UpdateFooter />
    </nav>
  );
}
