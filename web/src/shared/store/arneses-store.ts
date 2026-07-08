import { create } from "zustand";
import { api } from "../api/client";
import type { Arnes, RegistryEstado } from "../api/types";

/** Registry de arneses (PB-25): conexión a nivel app + roster instalado POR repo.
 *  El registry es la ÚNICA fuente de roles (DH-14) — acá no se crea nada, se consume. */
interface ArnesesState {
  registry: RegistryEstado | null;
  /** roster instalado por repoId (lock .devstudio/arneses.yaml del repo) */
  instalados: Record<string, Arnes[]>;
  error: string | null;
  busy: boolean;

  init: () => Promise<void>;
  conectar: (source: string) => Promise<boolean>;
  sync: () => Promise<void>;
  refreshInstalados: (repoId: string) => Promise<void>;
  instalar: (repoId: string, arnesId: string) => Promise<boolean>;
  desinstalar: (repoId: string, arnesId: string) => Promise<boolean>;
}

const msg = (e: unknown, fallback: string) => (e instanceof Error ? e.message : fallback);

export const useArneses = create<ArnesesState>((set, get) => ({
  registry: null,
  instalados: {},
  error: null,
  busy: false,

  init: async () => {
    try {
      set({ registry: await api.registry() });
    } catch {
      // API vieja sin registry (pre-update): la sección se muestra vacía honesta
    }
  },

  conectar: async (source: string) => {
    set({ busy: true, error: null });
    try {
      set({ registry: await api.conectarRegistry(source), busy: false });
      return true;
    } catch (e) {
      set({ error: msg(e, "No se pudo conectar el registry"), busy: false });
      return false;
    }
  },

  sync: async () => {
    set({ busy: true, error: null });
    try {
      set({ registry: await api.syncRegistry(), busy: false });
    } catch (e) {
      set({ error: msg(e, "No se pudo actualizar el registry"), busy: false });
    }
  },

  refreshInstalados: async (repoId: string) => {
    try {
      const lista = await api.arnesesInstalados(repoId);
      set((st) => ({ instalados: { ...st.instalados, [repoId]: lista } }));
    } catch {
      // repo sin lock todavía: roster vacío
      set((st) => ({ instalados: { ...st.instalados, [repoId]: [] } }));
    }
  },

  instalar: async (repoId: string, arnesId: string) => {
    set({ busy: true, error: null });
    try {
      await api.instalarArnes(repoId, arnesId);
      await get().refreshInstalados(repoId);
      set({ busy: false });
      return true;
    } catch (e) {
      set({ error: msg(e, "No se pudo instalar el arnés"), busy: false });
      return false;
    }
  },

  desinstalar: async (repoId: string, arnesId: string) => {
    set({ busy: true, error: null });
    try {
      await api.desinstalarArnes(repoId, arnesId);
      await get().refreshInstalados(repoId);
      set({ busy: false });
      return true;
    } catch (e) {
      set({ error: msg(e, "No se pudo desinstalar el arnés"), busy: false });
      return false;
    }
  },
}));
