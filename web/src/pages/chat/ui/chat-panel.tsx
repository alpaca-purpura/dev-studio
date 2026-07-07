import { useState } from "react";
import { useSessions } from "../../../shared/store/sessions-store";

// Panel mínimo para verificar el driver CLI-nativo end-to-end. Todo lo demás (Mapa, Portafolio,
// otras vistas del shell) queda en blanco a propósito — fuera de scope de esta épica.
export function ChatPanel() {
  const activeId = useSessions((s) => s.activeId);
  const session = useSessions((s) => s.sessions.find((x) => x.id === s.activeId));
  const streamBuffer = useSessions((s) => (s.activeId ? s.streamBuffer[s.activeId] : undefined));
  const sendTurn = useSessions((s) => s.sendTurn);
  const [text, setText] = useState("");

  if (!activeId || !session) {
    return (
      <div className="flex flex-1 items-center justify-center text-sm text-muted-foreground">
        Sin sesiones. Creá una desde el rail.
      </div>
    );
  }

  const submit = () => {
    const t = text.trim();
    if (!t || session.status === "streaming") return;
    setText("");
    void sendTurn(activeId, t);
  };

  return (
    <div className="flex flex-1 flex-col">
      <header className="flex items-center gap-2 border-b border-border px-4 py-2">
        <span className="text-sm font-semibold">{session.nombre}</span>
        <span className="font-mono text-[10px] text-muted-foreground">{session.cwd}</span>
        {session.model && (
          <span className="ml-auto font-mono text-[10px] text-muted-foreground">
            {session.model}
          </span>
        )}
      </header>

      <div className="flex-1 space-y-3 overflow-y-auto p-4">
        {(session.conv ?? []).map((t, i) => (
          <div
            key={i}
            className={
              t.role === "user"
                ? "ml-auto max-w-[70%] rounded-lg bg-primary px-3 py-2 text-sm text-primary-foreground"
                : "mr-auto max-w-[70%] whitespace-pre-wrap rounded-lg bg-secondary px-3 py-2 text-sm text-secondary-foreground"
            }
          >
            {t.text}
          </div>
        ))}
        {streamBuffer !== undefined && (
          <div className="mr-auto max-w-[70%] whitespace-pre-wrap rounded-lg bg-secondary px-3 py-2 text-sm text-secondary-foreground">
            {streamBuffer || "…"}
          </div>
        )}
      </div>

      <div className="flex items-center gap-2 border-t border-border p-3">
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") submit();
          }}
          placeholder={session.status === "streaming" ? "streaming…" : "Escribí un mensaje"}
          disabled={session.status === "streaming"}
          className="flex-1 rounded-md border border-input bg-card px-3 py-2 text-sm outline-none focus:border-primary disabled:opacity-60"
        />
        <button
          type="button"
          onClick={submit}
          disabled={session.status === "streaming"}
          className="rounded-md bg-primary px-3 py-2 text-sm font-semibold text-primary-foreground disabled:opacity-60"
        >
          Enviar
        </button>
      </div>
    </div>
  );
}
