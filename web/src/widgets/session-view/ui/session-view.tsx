import { useEffect, useRef, useState } from "react";
import { useSessions } from "../../../shared/store/sessions-store";
import { Avatar, Chip, ModalShell } from "../../../shared/ui";
import { cn } from "../../../shared/lib/cn";

/** Modales stub del composer (spec §3.3, RN-7): presentes, honestos, con destino. */
type StubKey = "boceto" | "captura" | "html" | "movil" | "prompts" | "voz" | "dictar" | "comandos" | "ssh";

const STUBS: Record<StubKey, { title: string; sub: string; body: string }> = {
  boceto: { title: "Boceto", sub: "Un dibujo dice lo que a las palabras les cuesta un párrafo.", body: "El lienzo de dibujo llega con la rebanada de composer completo (PB-09). Todavía no disponible." },
  captura: { title: "Captura de pantalla", sub: "Capturá cualquier área y colocala directo en el mensaje.", body: "Depende de permisos del escritorio — todavía no disponible (PB-09)." },
  html: { title: "Importar página web", sub: "Capturá una página como HTML.", body: "Requiere una extensión de navegador propia — fuera del alcance por ahora (spec clon §10.1)." },
  movil: { title: "Enviar desde el móvil", sub: "Archivos y capturas desde tu teléfono al composer.", body: "Depende de la app móvil (roadmap FN). Todavía no disponible." },
  prompts: { title: "Biblioteca de prompts", sub: "Fragmentos reutilizables para tus roles.", body: "Llega como rebanada propia (PB-16). Todavía no disponible." },
  voz: { title: "Conversación por voz", sub: "Conversación bidireccional hablada con el rol.", body: "Requiere backend propio — no incluida en esta fase. El botón queda visible a propósito: mensaje claro de «todavía no», no una feature a medias." },
  dictar: { title: "Dictar mensaje", sub: "Grabá tu voz y se transcribe en el borrador.", body: "Distinto de «Conversación por voz»: esto es solo transcripción. Todavía no disponible." },
  comandos: { title: "Comandos", sub: "Procesos que corren junto a tus roles: dev server, watchers.", body: "Los comandos persistentes llegan con la rebanada DevOps (PB-17). Todavía no disponible." },
  ssh: { title: "Conexiones remotas (SSH)", sub: "Un terminal en una máquina remota sin salir de la sesión.", body: "Llega con la rebanada DevOps (PB-17). Todavía no disponible." },
};

const COMPOSER_ICONS: { key: StubKey; title: string; path: string }[] = [
  { key: "boceto", title: "Boceto", path: "M4 20l4-1 11-11-3-3L5 16l-1 4Z" },
  { key: "captura", title: "Captura de pantalla", path: "M4 8V5a1 1 0 0 1 1-1h3M4 16v3a1 1 0 0 0 1 1h3M20 8V5a1 1 0 0 0-1-1h-3M20 16v3a1 1 0 0 1-1 1h-3" },
  { key: "html", title: "Importar página web", path: "M8 6 3 12l5 6M16 6l5 6-5 6" },
  { key: "movil", title: "Enviar desde el móvil", path: "M7 4a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v16a2 2 0 0 1-2 2H9a2 2 0 0 1-2-2zM11 18h2" },
  { key: "prompts", title: "Biblioteca de prompts", path: "M4 5a2 2 0 0 1 2-2h11v18H6a2 2 0 0 1-2-2ZM17 3v18M8 7h5M8 11h5" },
  { key: "voz", title: "Conversación por voz", path: "M12 2l1.6 4.9L18 8l-4.4 1.6L12 14l-1.6-4.4L6 8l4.4-1.1Z" },
  { key: "dictar", title: "Dictar mensaje", path: "M9 5a3 3 0 0 1 6 0v5a3 3 0 0 1-6 0zM5 11a7 7 0 0 0 14 0M12 18v3" },
];

const SIN_COLA: string[] = []; // referencia estable — un selector que fabrica [] nuevo por snapshot loopea React (#185)

export function SessionView() {
  const session = useSessions((s) => s.sessions.find((x) => x.id === s.activeId));
  const streamBuffer = useSessions((s) => (s.activeId ? s.streamBuffer[s.activeId] : undefined));
  const queued = useSessions((s) => (s.activeId ? (s.queue[s.activeId] ?? SIN_COLA) : SIN_COLA));
  const sendTurn = useSessions((s) => s.sendTurn);
  const enqueue = useSessions((s) => s.enqueue);
  const [text, setText] = useState("");
  const [stub, setStub] = useState<StubKey | null>(null);
  const [composerH, setComposerH] = useState(88);
  const scrollRef = useRef<HTMLDivElement>(null);
  const dragRef = useRef<{ y: number; h: number } | null>(null);

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight });
  }, [session?.conv.length, streamBuffer]);

  if (!session) {
    return (
      <div className="flex flex-1 items-center justify-center text-sm text-muted-foreground">
        Sin sesiones. Registrá un repositorio y creá un workspace desde el rail.
      </div>
    );
  }

  const streaming = session.status === "streaming";
  const canSend = text.trim().length > 0;

  const doSend = () => {
    const t = text.trim();
    if (!t || streaming) return;
    setText("");
    void sendTurn(session.id, t);
  };
  const doQueue = () => {
    const t = text.trim();
    if (!t) return;
    setText("");
    enqueue(session.id, t);
  };

  const startDrag = (e: React.PointerEvent) => {
    dragRef.current = { y: e.clientY, h: composerH };
    const move = (ev: PointerEvent) => {
      if (!dragRef.current) return;
      const dh = dragRef.current.y - ev.clientY;
      setComposerH(Math.min(320, Math.max(64, dragRef.current.h + dh)));
    };
    const up = () => {
      dragRef.current = null;
      window.removeEventListener("pointermove", move);
      window.removeEventListener("pointerup", up);
    };
    window.addEventListener("pointermove", move);
    window.addEventListener("pointerup", up);
  };

  return (
    <div className="flex min-w-0 flex-1 flex-col">
      {/* header de sesión */}
      <div className="flex items-center gap-3 border-b border-border px-5 py-3">
        <Avatar label={session.rol ? siglas(session.rol) : "··"} />
        <div className="min-w-0">
          <div className="truncate text-sm font-semibold">
            {session.rol ? `${session.rol} · ` : ""}
            {session.nombre}
          </div>
          <div className="truncate font-mono text-[10px] text-muted-foreground">
            {session.cwd} · Claude Code{session.model ? ` (${session.model})` : ""}
          </div>
        </div>
        <span className="flex-1" />
        {session.modo === "exploracion" && (
          <Chip
            className="border-primary/40 bg-accent-soft text-primary"
            title="Sesión de exploración: Claude Code en modo plan — lee y analiza, sin permisos de edición"
          >
            🔍 Exploración · solo lectura
          </Chip>
        )}
        {session.branch && <Chip title={`Workspace aislado: ${session.workspace ?? ""}`}>⎇ {session.branch}</Chip>}
        {session.historia && <Chip title="Paquete de trabajo ligado (RN-1)">📋 {session.historia.titulo}</Chip>}
      </div>

      {/* conversación */}
      <div ref={scrollRef} className="min-h-0 flex-1 space-y-4 overflow-y-auto p-5">
        {(session.conv ?? []).map((t, i) => (
          <div key={i} className={cn("max-w-[75%]", t.role === "user" ? "ml-auto" : "mr-auto")}>
            <div className="mb-1 font-mono text-[9px] uppercase tracking-[0.14em] text-muted-foreground">
              {t.role === "user" ? "Tú" : (session.rol ?? "Asistente")}
            </div>
            <div
              className={cn(
                "whitespace-pre-wrap rounded-md border px-3 py-2 text-sm leading-relaxed",
                t.role === "user"
                  ? "border-primary/30 bg-accent-soft text-foreground"
                  : "border-border bg-card text-card-foreground",
              )}
            >
              {t.text}
            </div>
          </div>
        ))}
        {streamBuffer !== undefined && (
          <div className="mr-auto max-w-[75%]">
            <div className="mb-1 font-mono text-[9px] uppercase tracking-[0.14em] text-muted-foreground">
              {session.rol ?? "Asistente"}
            </div>
            <div className="whitespace-pre-wrap rounded-md border border-border bg-card px-3 py-2 text-sm leading-relaxed">
              {streamBuffer || "…"}
            </div>
          </div>
        )}
      </div>

      {/* composer */}
      <div className="border-t border-border px-5 py-3">
        <div className="rounded-lg border border-border bg-card">
          <div
            onPointerDown={startDrag}
            title="Arrastrá para redimensionar el composer"
            className="mx-auto mt-1 h-1 w-10 cursor-ns-resize rounded-full bg-border"
          />
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                if (streaming) doQueue();
                else doSend();
              }
            }}
            style={{ height: composerH }}
            placeholder="Redactá un mensaje: Enter para enviar, Shift+Enter para una nueva línea"
            className="w-full resize-none bg-transparent px-4 py-2 text-sm text-foreground outline-none placeholder:text-muted-foreground/60"
          />
          <div className="flex items-center gap-1 px-3 pb-3">
            {COMPOSER_ICONS.map((b, i) => (
              <span key={b.key} className="flex items-center">
                {(i === 2 || i === 5) && <span className="mx-1 h-4 w-px bg-border" />}
                <button
                  title={`${b.title} — todavía no (RN-7)`}
                  onClick={() => setStub(b.key)}
                  className="cursor-pointer rounded-full p-1.5 text-muted-foreground/70 hover:bg-secondary hover:text-foreground"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="size-4">
                    <path d={b.path} />
                  </svg>
                </button>
              </span>
            ))}
            <span className="flex-1" />
            {queued.length > 0 && (
              <Chip title="Mensajes en cola — se despachan al terminar el turno (RN-8)">
                ⇒ {queued.length} en cola
              </Chip>
            )}
            <button
              onClick={doQueue}
              disabled={!canSend}
              title="Añadir a la cola — se envía cuando el rol termine el turno en curso"
              className="cursor-pointer rounded-full border border-border px-3 py-1 text-xs font-semibold text-muted-foreground hover:text-foreground disabled:cursor-not-allowed disabled:opacity-40"
            >
              ⇒ Encolar
            </button>
            <button
              onClick={doSend}
              disabled={!canSend || streaming}
              title={streaming ? "Turno en curso — encolá el mensaje" : "Enviar (Enter)"}
              aria-label="Enviar"
              className="cursor-pointer rounded-full bg-primary p-2 text-primary-foreground shadow-[var(--shadow-glow)] transition-transform hover:scale-105 disabled:cursor-not-allowed disabled:opacity-40 disabled:shadow-none"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="size-4">
                <path d="M3 12l18-8-8 18-2-8-8-2Z" />
              </svg>
            </button>
          </div>
          {/* fila terminal — stub honesto (terminal PTY = PB-09; Comandos/SSH = PB-17) */}
          <div className="flex items-center gap-2 border-t border-border px-3 py-2">
            <span
              title="Terminal embebida — llega con PB-09"
              className="rounded-t-md border border-b-0 border-border px-2 py-0.5 font-mono text-[10px] text-muted-foreground/50"
            >
              ⌨ Terminal
            </span>
            <button
              onClick={() => setStub("comandos")}
              className="cursor-pointer rounded-full border border-border px-2 py-0.5 text-[10px] text-muted-foreground hover:text-foreground"
            >
              ⚡ Comandos de desarrollo
            </button>
            <button
              onClick={() => setStub("ssh")}
              className="cursor-pointer rounded-full border border-border px-2 py-0.5 text-[10px] text-muted-foreground hover:text-foreground"
            >
              🖧 SSH remoto
            </button>
          </div>
        </div>
      </div>

      <ModalShell
        open={stub !== null}
        onClose={() => setStub(null)}
        title={stub ? STUBS[stub].title : ""}
        subtitle={stub ? STUBS[stub].sub : undefined}
        icon="◔"
      >
        <p className="text-sm leading-relaxed text-muted-foreground">{stub ? STUBS[stub].body : null}</p>
      </ModalShell>
    </div>
  );
}

function siglas(rol: string): string {
  const parts = rol.split(/\s+/);
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase();
  return rol.slice(0, 2).toUpperCase();
}
