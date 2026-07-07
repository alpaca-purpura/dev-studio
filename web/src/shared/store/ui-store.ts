import { create } from "zustand";
import type { Historia } from "../api/types";

export type NavKey = "studio" | "backlog" | "producto" | "roles" | "config";

interface UiState {
  nav: NavKey;
  /** sesión con cierre pendiente — dispara el modal conservar/borrar workspace (PB-02) */
  closeRequestId: string | null;
  /** wizard «Nuevo Workspace» abierto para este repo (PB-27) */
  wizardRepoId: string | null;
  /** picker: el Backlog abierto para ELEGIR el paquete de una sesión nueva (RN-1) */
  pickerMode: boolean;
  /** historia elegida en el picker o creada en el wizard, esperando en Config */
  pendingHistoria: Historia | null;
  /** repo para el que se está creando la sesión (desde su `+ Nuevo Workspace`) */
  pendingRepoId: string | null;
  rolElegido: string;

  setNav: (k: NavKey) => void;
  requestClose: (sessionId: string | null) => void;
  /** `+ Nuevo Workspace` de un repo → abre el wizard (PB-27) */
  openWizard: (repoId: string) => void;
  closeWizard: () => void;
  /** desde el wizard: «trabajar un ítem existente» → Backlog en modo picker */
  startPicker: (repoId: string) => void;
  /** clic en una historia estando en picker (o ítem creado en el wizard) → Config con el paquete */
  pickHistoria: (h: Historia) => void;
  setRol: (rol: string) => void;
  /** vuelta a Studio limpiando estado transitorio del flujo */
  resetToStudio: () => void;
  clearPending: () => void;
}

export const useUi = create<UiState>((set) => ({
  nav: "studio",
  closeRequestId: null,
  wizardRepoId: null,
  pickerMode: false,
  pendingHistoria: null,
  pendingRepoId: null,
  rolElegido: "Full-Stack",

  setNav: (k) => set({ nav: k, pickerMode: false }),
  requestClose: (sessionId) => set({ closeRequestId: sessionId }),
  openWizard: (repoId) => set({ wizardRepoId: repoId, pendingRepoId: repoId }),
  closeWizard: () => set({ wizardRepoId: null }),
  startPicker: (repoId) => set({ wizardRepoId: null, nav: "backlog", pickerMode: true, pendingRepoId: repoId }),
  pickHistoria: (h) => set({ wizardRepoId: null, nav: "config", pickerMode: false, pendingHistoria: h }),
  setRol: (rol) => set({ rolElegido: rol }),
  resetToStudio: () => set({ nav: "studio", pickerMode: false }),
  clearPending: () => set({ pendingHistoria: null, pendingRepoId: null }),
}));
