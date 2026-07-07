import type { StoryState } from "../ui/state-pill";
import type { TipoItem } from "../lib/tipos";

/** Ítems de EJEMPLO del overlay Backlog (spec shell §3.4, RN-7): el board real llega con
 *  el port de Historia/Capability (PB-07). El modo picker SÍ es real — elegir uno liga la
 *  sesión nueva a su paquete. `tipoItem` = taxonomía estándar (spec nuevo-workspace §1). */
export interface MockStory {
  id: string;
  titulo: string;
  estado: StoryState;
  rol: string | null;
  tipoItem: TipoItem;
  release: "F1" | "F2";
  capability: string;
  prioridad: "crítica" | "alta" | "media" | "baja";
}

export const MOCK_BACKLOG: MockStory[] = [
  { id: "hx-panorama", titulo: "Panorama transversal de sesiones entre repositorios", estado: "idea", rol: null, tipoItem: "historia", release: "F2", capability: "multi-repositorio", prioridad: "media" },
  { id: "hx-changelog", titulo: "Publicar changelog automático por capability", estado: "idea", rol: null, tipoItem: "tarea", release: "F2", capability: "backlog-historias", prioridad: "baja" },
  { id: "hx-board", titulo: "Tablero backlog con los estados del proceso", estado: "refining", rol: "PO", tipoItem: "historia", release: "F1", capability: "backlog-historias", prioridad: "alta" },
  { id: "hx-checkpoints", titulo: "Checkpoints y revert por turno de conversación", estado: "refined", rol: null, tipoItem: "historia", release: "F2", capability: "sesiones-multiples", prioridad: "media" },
  { id: "hx-diff-conv", titulo: "Diff conversacional: comentar una línea y que el rol la resuelva", estado: "ready", rol: null, tipoItem: "historia", release: "F2", capability: "revision-cambios", prioridad: "media" },
  { id: "hx-spike-pty", titulo: "Spike: evaluar librerías PTY para el terminal embebido", estado: "ready", rol: null, tipoItem: "spike", release: "F2", capability: "comandos-dev", prioridad: "media" },
  { id: "hx-worktree", titulo: "Aislar sesión en workspace propio (PB-02)", estado: "developing", rol: "FS", tipoItem: "tarea", release: "F1", capability: "sesiones-multiples", prioridad: "crítica" },
  { id: "hx-sse-retry", titulo: "Resolver reconexión SSE con retry", estado: "developing", rol: "BE", tipoItem: "bug", release: "F1", capability: "sesiones-multiples", prioridad: "alta" },
  { id: "hx-watcher", titulo: "Watcher de build para Vite (mantenimiento)", estado: "developed", rol: "DO", tipoItem: "tarea", release: "F1", capability: "comandos-dev", prioridad: "baja" },
  { id: "hx-rutas", titulo: "Validar rutas protegidas por sesión", estado: "reviewing", rol: "QA", tipoItem: "bug", release: "F1", capability: "sesiones-multiples", prioridad: "crítica" },
  { id: "hx-roster", titulo: "Vista Roles conectada al registry (PB-06)", estado: "reviewing", rol: "AR", tipoItem: "historia", release: "F2", capability: "roles-roster", prioridad: "media" },
  { id: "hx-rail", titulo: "Rail de sesiones colapsable (F1)", estado: "done", rol: "FS", tipoItem: "historia", release: "F1", capability: "sesiones-multiples", prioridad: "media" },
  { id: "hx-pausada", titulo: "Compañero móvil (fuera de alcance por ahora)", estado: "parked", rol: null, tipoItem: "historia", release: "F2", capability: "multi-repositorio", prioridad: "baja" },
  { id: "hx-dropped", titulo: "Torre de control transversal (eliminada del shell, §0.5)", estado: "dropped", rol: null, tipoItem: "historia", release: "F2", capability: "multi-repositorio", prioridad: "baja" },
];

/** Roster estático — vista previa del registry (PB-25). Refleja los grupos del mockup. */
export interface MockRol {
  sigla: string;
  nombre: string;
  grupo: "Builder" | "Auditor" | "Humano-complementario";
  descripcion: string;
  prompt: string;
}

export const MOCK_ROSTER: MockRol[] = [
  {
    sigla: "FS", nombre: "Full-Stack", grupo: "Builder",
    descripcion: "Dueño de dominio, casos de uso, adaptadores y la SPA embebida. Escribe el test primero, siempre.",
    prompt: "# Identidad\nSos el rol Full-Stack de este proyecto. Trabajás en Go (arquitectura hexagonal:\ndomain/ports/usecase/adapters) y en la SPA embebida (React + Zustand + Tailwind).\n\n# Reglas\n- TDD obligatorio: test en rojo antes que implementación.\n- Un turno = un cambio revisable. No mezclás refactors con features.\n- Si el diff toca un boundary de arch/, avisás antes de escribir.\n\n# Entregables\n- Diff acotado + test nuevo o migrado + una línea de qué RN cubre.",
  },
  { sigla: "FE", nombre: "Frontend", grupo: "Builder", descripcion: "Especialista en la SPA embebida y el design system PRENTER.", prompt: "# Identidad\nSos el rol Frontend. Componés UI solo con átomos catalogados en Storybook (RN-9)." },
  { sigla: "BE", nombre: "Backend", grupo: "Builder", descripcion: "Especialista en el binario Go: dominio, puertos, adaptadores.", prompt: "# Identidad\nSos el rol Backend. Hexagonal estricta; TDD obligatorio." },
  { sigla: "DO", nombre: "DevOps", grupo: "Builder", descripcion: "Comandos persistentes, túneles, SSH — dueño de gates de deploy.", prompt: "# Identidad\nSos el rol DevOps del proyecto." },
  { sigla: "QA", nombre: "QA", grupo: "Auditor", descripcion: "Verifica escenarios reales contra la app viva, nunca «GET 200».", prompt: "# Identidad\nSos el rol QA. Verificación real: acción ejercida + efecto observado." },
  { sigla: "SG", nombre: "Seguridad", grupo: "Auditor", descripcion: "Rutas protegidas, aislamiento entre sesiones, superficie de la API local.", prompt: "# Identidad\nSos el rol Seguridad." },
  { sigla: "AR", nombre: "Arquitecto de Software", grupo: "Auditor", descripcion: "Custodia los boundaries de arch/ y el fitness.", prompt: "# Identidad\nSos el rol Arquitecto. Los boundaries de arch/ son tu contrato." },
  { sigla: "PO", nombre: "Product Owner", grupo: "Humano-complementario", descripcion: "Refina historias ligadas a capabilities; dueño del gate refined→ready.", prompt: "# Identidad\nSos el rol Product Owner. Historias con criterios verificables, siempre ligadas a una capability." },
  { sigla: "ID", nombre: "Ideación", grupo: "Humano-complementario", descripcion: "Convierte ideas crudas en historias con forma.", prompt: "# Identidad\nSos el rol Ideación." },
];
