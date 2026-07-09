import { create } from "zustand";
import { api } from "../api/client";
import { connectDock } from "../api/sse";
import type { CloseSessionResp, DockFrame, Historia, Session, ToolCall } from "../api/types";
import { useRepos } from "./repos-store";

interface SessionsState {
  sessions: Session[];
  activeId: string | null;
  streamBuffer: Record<string, string>;
  /** Tarjetas de herramienta del turno en curso por sesión (R1). Se reinician al enviar un
   *  turno nuevo; la llamada y su resultado se parean por tool_id. */
  toolCalls: Record<string, ToolCall[]>;
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
  queue: {},

  init: async () => {
    const sessions = await api.list();
    set({ sessions, activeId: sessions[0]?.id ?? null });
    disconnect?.();
    disconnect = connectDock((f) => get().onDock(f));
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

  switchTo: (id: string) => set({ activeId: id }),

  sendTurn: async (id: string, text: string) => {
    set((st) => ({
      // turno nuevo ⇒ las tarjetas de herramienta del turno anterior se limpian
      toolCalls: { ...st.toolCalls, [id]: [] },
      sessions: st.sessions.map((s) =>
        s.id === id
          ? { ...s, status: "streaming", conv: [...s.conv, { role: "user", text }] }
          : s,
      ),
    }));
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
