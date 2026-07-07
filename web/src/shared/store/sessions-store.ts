import { create } from "zustand";
import { api } from "../api/client";
import { connectDock } from "../api/sse";
import type { DockFrame, Session } from "../api/types";

interface SessionsState {
  sessions: Session[];
  activeId: string | null;
  railCollapsed: boolean;
  streamBuffer: Record<string, string>;

  init: () => Promise<void>;
  create: (cwd: string) => Promise<void>;
  closeSession: (id: string) => Promise<void>;
  rename: (id: string, nombre: string) => Promise<void>;
  switchTo: (id: string) => void;
  toggleRail: () => void;
  sendTurn: (id: string, text: string) => Promise<void>;
  onDock: (f: DockFrame) => void;
}

let disconnect: (() => void) | null = null;

export const useSessions = create<SessionsState>((set, get) => ({
  sessions: [],
  activeId: null,
  railCollapsed: false,
  streamBuffer: {},

  init: async () => {
    const sessions = await api.list();
    set({ sessions, activeId: sessions[0]?.id ?? null });
    disconnect?.();
    disconnect = connectDock((f) => get().onDock(f));
  },

  create: async (cwd: string) => {
    const sess = await api.create("Nueva sesión", cwd);
    set((st) => ({ sessions: [...st.sessions, sess], activeId: sess.id }));
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
  toggleRail: () => set((st) => ({ railCollapsed: !st.railCollapsed })),

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
  },
}));
