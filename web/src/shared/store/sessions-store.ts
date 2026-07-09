import { create } from "zustand";
import { api } from "../api/client";
import { connectDock } from "../api/sse";
import type { CloseSessionResp, DockFrame, Historia, Session, ToolCall, TranscriptItem, Turn } from "../api/types";

/** conv (solo texto) → ítems del transcript, para sembrar el optimista antes de que cargue el rico. */
const convToItems = (conv: Turn[]): TranscriptItem[] => conv.map((t) => ({ kind: t.role, text: t.text }));
import { useRepos } from "./repos-store";

interface SessionsState {
  sessions: Session[];
  activeId: string | null;
  streamBuffer: Record<string, string>;
  /** Tarjetas de herramienta del turno en curso por sesión (R1). Se reinician al enviar un
   *  turno nuevo; la llamada y su resultado se parean por tool_id. */
  toolCalls: Record<string, ToolCall[]>;
  /** Transcript reconstruido por sesión (R1.5): historial ordenado (texto + tool-cards en su
   *  lugar) que sobrevive al reentrar. Fuente única de render; el turno vivo va como overlay. */
  transcript: Record<string, TranscriptItem[]>;
  /** Cola de mensajes por sesión (Encolar ⇒). Client-side, respeta un-turno-a-la-vez (RN-8):
   *  el siguiente mensaje se despacha recién cuando llega result/error del turno en curso. */
  queue: Record<string, string[]>;

  init: () => Promise<void>;
  create: (p: {
    nombre: string;
    repo_id?: string;
    historia?: Historia;
    rol?: string;
    modo?: "trabajo" | "exploracion";
    ubicacion?: "worktree" | "checkout";
  }) => Promise<Session>;
  closeSession: (id: string, workspace?: "keep" | "remove") => Promise<CloseSessionResp>;
  rename: (id: string, nombre: string) => Promise<void>;
  switchTo: (id: string) => void;
  loadTranscript: (id: string) => Promise<void>;
  sendTurn: (id: string, text: string) => Promise<void>;
  enqueue: (id: string, text: string) => void;
  onDock: (f: DockFrame) => void;
}

let disconnect: (() => void) | null = null;

export const useSessions = create<SessionsState>((set, get) => ({
  sessions: [],
  activeId: null,
  streamBuffer: {},
  toolCalls: {},
  transcript: {},
  queue: {},

  init: async () => {
    const sessions = await api.list();
    const activeId = sessions[0]?.id ?? null;
    set({ sessions, activeId });
    disconnect?.();
    disconnect = connectDock((f) => get().onDock(f));
    if (activeId) void get().loadTranscript(activeId);
  },

  create: async (p) => {
    const sess = await api.create(p);
    set((st) => ({ sessions: [...st.sessions, sess], activeId: sess.id }));
    return sess;
  },

  closeSession: async (id: string, workspace: "keep" | "remove" = "keep") => {
    const resp = await api.close(id, workspace);
    set((st) => {
      const sessions = st.sessions.filter((s) => s.id !== id);
      const activeId = st.activeId === id ? (sessions[0]?.id ?? null) : st.activeId;
      return { sessions, activeId };
    });
    return resp;
  },

  rename: async (id: string, nombre: string) => {
    set((st) => ({
      sessions: st.sessions.map((s) => (s.id === id ? { ...s, nombre } : s)),
    }));
    await api.rename(id, nombre);
  },

  switchTo: (id: string) => {
    set({ activeId: id });
    void get().loadTranscript(id);
  },

  // loadTranscript trae el historial reconstruido (R1.5) y limpia el overlay vivo — las cards
  // que estaban en vivo ahora viven en el transcript, en su lugar.
  loadTranscript: async (id: string) => {
    try {
      const items = await api.transcript(id);
      set((st) => ({
        transcript: { ...st.transcript, [id]: items },
        toolCalls: { ...st.toolCalls, [id]: [] },
      }));
    } catch {
      /* sin transcript (sesión nueva / offline): el render cae al conv */
    }
  },

  sendTurn: async (id: string, text: string) => {
    set((st) => {
      const base = st.transcript[id] ?? convToItems(st.sessions.find((s) => s.id === id)?.conv ?? []);
      return {
        // turno nuevo ⇒ overlay de tool-cards limpio; el prompt entra optimista al transcript
        toolCalls: { ...st.toolCalls, [id]: [] },
        transcript: { ...st.transcript, [id]: [...base, { kind: "user", text }] },
        sessions: st.sessions.map((s) =>
          s.id === id ? { ...s, status: "streaming", conv: [...s.conv, { role: "user", text }] } : s,
        ),
      };
    });
    try {
      await api.turn(id, text);
    } catch (e) {
      // el POST falló (p. ej. spawn de claude imposible): jamás quedar «streaming» mudo —
      // se revierte el estado y el error REAL entra como burbuja (bug dogfooding DH-16.1)
      const msg = e instanceof Error ? e.message : String(e);
      set((st) => ({
        sessions: st.sessions.map((s) =>
          s.id === id
            ? { ...s, status: "idle", conv: [...s.conv, { role: "assistant", text: `⚠ No se pudo enviar el turno: ${msg}` }] }
            : s,
        ),
      }));
    }
  },

  enqueue: (id: string, text: string) => {
    const sess = get().sessions.find((s) => s.id === id);
    if (sess && sess.status !== "streaming") {
      void get().sendTurn(id, text);
      return;
    }
    set((st) => ({ queue: { ...st.queue, [id]: [...(st.queue[id] ?? []), text] } }));
  },

  onDock: (f: DockFrame) => {
    set((st) => {
      switch (f.kind) {
        case "status":
          return {
            sessions: st.sessions.map((s) =>
              s.id === f.session_id ? { ...s, status: f.status ?? s.status } : s,
            ),
          };
        case "init":
          return {
            sessions: st.sessions.map((s) =>
              s.id === f.session_id
                ? { ...s, claude_session_id: f.claude_session_id, model: f.model }
                : s,
            ),
          };
        case "delta":
          return {
            streamBuffer: {
              ...st.streamBuffer,
              [f.session_id]: (st.streamBuffer[f.session_id] ?? "") + (f.text ?? ""),
            },
          };
        case "tool.call": {
          if (!f.tool_id) return {};
          const cur = st.toolCalls[f.session_id] ?? [];
          const nc: ToolCall = {
            tool_id: f.tool_id,
            name: f.tool_name ?? "",
            input: f.tool_input ?? "",
            status: "running",
          };
          return { toolCalls: { ...st.toolCalls, [f.session_id]: [...cur, nc] } };
        }
        case "tool.result": {
          if (!f.tool_id) return {};
          const cur = st.toolCalls[f.session_id] ?? [];
          return {
            toolCalls: {
              ...st.toolCalls,
              [f.session_id]: cur.map((c): ToolCall =>
                c.tool_id === f.tool_id
                  ? { ...c, status: f.tool_is_error ? "error" : "ok", output: f.text }
                  : c,
              ),
            },
          };
        }
        case "result": {
          const buf = { ...st.streamBuffer };
          delete buf[f.session_id];
          return {
            streamBuffer: buf,
            sessions: st.sessions.map((s) =>
              s.id === f.session_id
                ? { ...s, status: "idle", conv: [...s.conv, { role: "assistant", text: f.text ?? "" }] }
                : s,
            ),
          };
        }
        case "error": {
          const buf = { ...st.streamBuffer };
          delete buf[f.session_id];
          return {
            streamBuffer: buf,
            sessions: st.sessions.map((s) =>
              s.id === f.session_id
                ? {
                    ...s,
                    status: "idle",
                    conv: [...s.conv, { role: "assistant", text: `⚠ error: ${f.text ?? ""}` }],
                  }
                : s,
            ),
          };
        }
        default:
          return {};
      }
    });

    // turno terminado → refrescar el status git del workspace (branch/±N del rail y Cambios)
    if (f.kind === "result" || f.kind === "error") {
      void useRepos.getState().refreshSessionStatus(f.session_id);
      // …y recargar el transcript (R1.5): el turno recién cerrado ya está en el JSONL → las
      // tool-cards vivas se pliegan a su lugar en el historial. Pequeño respiro para que asiente.
      const sid = f.session_id;
      setTimeout(() => void get().loadTranscript(sid), 300);
    }
    // turno terminado → despachar el siguiente de la cola (RN-8: nunca en paralelo)
    if (f.kind === "result" || f.kind === "error") {
      const { queue, sendTurn } = get();
      const pending = queue[f.session_id];
      if (pending && pending.length > 0) {
        const [next, ...rest] = pending;
        set((st) => ({ queue: { ...st.queue, [f.session_id]: rest } }));
        void sendTurn(f.session_id, next);
      }
    }
  },
}));
