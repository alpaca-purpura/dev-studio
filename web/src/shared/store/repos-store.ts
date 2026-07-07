import { create } from "zustand";
import { api } from "../api/client";
import type { GitStatus, Repo } from "../api/types";

interface ReposState {
  repos: Repo[];
  /** status git por repo (branch + ±N para el rail) — refrescado on-demand */
  status: Record<string, GitStatus>;
  /** status git POR SESIÓN (su worktree propio, PB-02) — la verdad del workspace-item */
  sessionStatus: Record<string, GitStatus>;
  expanded: Record<string, boolean>;
  addOpen: boolean;
  addError: string | null;

  init: () => Promise<void>;
  register: (ruta: string) => Promise<Repo | null>;
  remove: (id: string) => Promise<void>;
  refreshStatus: (repoId: string) => Promise<void>;
  refreshSessionStatus: (sessionId: string) => Promise<void>;
  toggleExpanded: (repoId: string) => void;
  setAddOpen: (v: boolean) => void;
}

export const useRepos = create<ReposState>((set, get) => ({
  repos: [],
  status: {},
  sessionStatus: {},
  expanded: {},
  addOpen: false,
  addError: null,

  init: async () => {
    const repos = await api.repos();
    set({ repos });
    for (const r of repos) void get().refreshStatus(r.id);
  },

  register: async (ruta: string) => {
    try {
      const repo = await api.registerRepo(ruta);
      set((st) => ({
        repos: st.repos.some((r) => r.id === repo.id) ? st.repos : [...st.repos, repo],
        addOpen: false,
        addError: null,
      }));
      void get().refreshStatus(repo.id);
      return repo;
    } catch (e) {
      set({ addError: e instanceof Error ? e.message : "No se pudo registrar el repositorio" });
      return null;
    }
  },

  remove: async (id: string) => {
    await api.removeRepo(id);
    set((st) => ({ repos: st.repos.filter((r) => r.id !== id) }));
  },

  refreshStatus: async (repoId: string) => {
    try {
      const st = await api.repoGitStatus(repoId);
      set((s) => ({ status: { ...s.status, [repoId]: st } }));
    } catch {
      // repo movido/borrado del disco: el rail lo muestra sin status
    }
  },

  refreshSessionStatus: async (sessionId: string) => {
    try {
      const st = await api.gitStatus(sessionId);
      set((s) => ({ sessionStatus: { ...s.sessionStatus, [sessionId]: st } }));
    } catch {
      // sesión sin cwd válido (legacy con directorio borrado): sin status, sin drama
    }
  },

  toggleExpanded: (repoId: string) =>
    set((st) => ({ expanded: { ...st.expanded, [repoId]: !(st.expanded[repoId] ?? true) } })),

  setAddOpen: (v: boolean) => set({ addOpen: v, addError: null }),
}));
