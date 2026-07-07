import { create } from "zustand";
import type { Historia } from "../api/types";

export type NavKey = "studio" | "backlog" | "producto" | "roles" | "config";

interface UiState {
  nav: NavKey;
  /** picker: el Backlog abierto para ELEGIR el paquete de una sesión nueva (RN-1) */
  pickerMode: boolean;
  /** historia elegida en el picker, esperando en Config a `Crear sesión aislada` */
  pendingHistoria: Historia | null;
  /** repo para el que se está creando la sesión (desde su `+ Nuevo Workspace`) */
  pendingRepoId: string | null;
  rolElegido: string;

  setNav: (k: NavKey) => void;
  /** `+ Nuevo Workspace` de un repo → Backlog en modo picker */
  startPicker: (repoId: string) => void;
  /** clic en una historia estando en picker → Config con el paquete */
  pickHistoria: (h: Historia) => void;
  setRol: (rol: string) => void;
  /** vuelta a Studio limpiando estado transitorio del flujo */
  resetToStudio: () => void;
  clearPending: () => void;
}

export const useUi = create<UiState>((set) => ({
  nav: "studio",
  pickerMode: false,
  pendingHistoria: null,
  pendingRepoId: null,
  rolElegido: "Full-Stack",

  setNav: (k) => set({ nav: k, pickerMode: false }),
  startPicker: (repoId) => set({ nav: "backlog", pickerMode: true, pendingRepoId: repoId }),
  pickHistoria: (h) => set({ nav: "config", pickerMode: false, pendingHistoria: h }),
  setRol: (rol) => set({ rolElegido: rol }),
  resetToStudio: () => set({ nav: "studio", pickerMode: false }),
  clearPending: () => set({ pendingHistoria: null, pendingRepoId: null }),
}));
