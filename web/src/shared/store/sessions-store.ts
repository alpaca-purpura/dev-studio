import { create } from "zustand";
import { api } from "../api/client";
import { connectDock } from "../api/sse";
import type { DockFrame, Historia, Session } from "../api/types";

interface SessionsState {
  sessions: Session[];
  activeId: string | null;
  streamBuffer: Record<string, string>;
  /** Cola de mensajes por sesión (Encolar ⇒). Client-side, respeta un-turno-a-la-vez (RN-8):
   *  el siguiente mensaje se despacha recién cuando llega result/error del turno en curso. */
  queue: Record<string, string[]>;

  init: () => Promise<void>;
  create: (p: { nombre: string; repo_id?: string; historia?: Historia; rol?: string }) => Promise<Session>;
  closeSession: (id: string) => Promise<void>;
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

  closeSession: async (id: string) => {
    await api.close(id);
    set((st) => {
      const sessions = st.sessions.filter((s) => s.id !== id);
      const activeId = st.activeId === id ? (sessions[0]?.id ?? null) : st.activeId;
      return { sessions, activeId };
    });
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
      sessions: st.sessions.map((s) =>
        s.id === id
          ? { ...s, status: "streaming", conv: [...s.conv, { role: "user", text }] }
          : s,
      ),
    }));
    await api.turn(id, text);
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
