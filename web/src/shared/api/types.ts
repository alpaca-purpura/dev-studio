export type SessionStatus = "streaming" | "idle";

export interface Turn {
  role: "user" | "assistant";
  text: string;
}

import type { TipoItem } from "../lib/tipos";

export interface Historia {
  id: string;
  titulo: string;
  tipo?: TipoItem; // vacío = historia (legacy)
}

export interface Session {
  id: string;
  nombre: string;
  cwd: string;
  status: SessionStatus;
  claude_session_id?: string;
  model?: string;
  repo_id?: string;
  historia?: Historia;
  rol?: string;
  workspace?: string; // ruta del worktree propio (PB-02); vacío = sin aislamiento
  branch?: string; // {prefijo-tipo}/{slug} (PB-27)
  modo?: string; // trabajo (default) | exploracion (read-only)
  conv: Turn[];
}

export interface CloseSessionResp {
  closed: boolean;
  workspace_removed: boolean;
  detalle?: string;
}

export interface VersionInfo {
  version: string;
  build_date: string;
  source: string;
}

export interface Repo {
  id: string;
  nombre: string;
  ruta: string;
}

/** Arnés del registry (PB-25) — formato ArnesIA (nomenclatura-arnes v1): meta rol×proceso. */
export interface Arnes {
  id: string;
  nombre: string;
  descripcion?: string;
  rol: string;
  proceso?: string;
  version: string;
  canal?: string;
  fases?: string[];
}

export interface RegistryEstado {
  source: string;
  arneses: Arnes[] | null;
}

export interface GitFile {
  path: string;
  state: string; // M | A | D | R
}

export interface GitStatus {
  branch: string;
  files: GitFile[];
  add: number;
  del: number;
}

export interface GitDiff {
  path: string;
  unified: string;
  original: string;
  modified: string;
  binary: boolean;
}

export interface GitLogEntry {
  sha: string;
  mensaje: string;
  autor: string;
  fecha: string;
}

export interface DockFrame {
  session_id: string;
  kind: "status" | "init" | "delta" | "result" | "error" | "tool.call" | "tool.result";
  text?: string;
  status?: SessionStatus;
  claude_session_id?: string;
  model?: string;
  // herramientas (R1 tool-cards) — tool_id parea la llamada con su resultado
  tool_id?: string;
  tool_name?: string;
  tool_input?: string; // JSON crudo del input de la herramienta
  tool_is_error?: boolean;
}

/** Una invocación de herramienta del agente + su resultado (R1 tool-cards). */
export interface ToolCall {
  tool_id: string;
  name: string; // "Bash" | "Read" | "Edit" | …
  input: string; // JSON crudo del input
  status: "running" | "ok" | "error";
  output?: string; // llega con el tool.result
}

/** Ítem del transcript reconstruido (R1.5): historial ordenado que sobrevive al reentrar. */
export interface TranscriptItem {
  kind: "user" | "assistant" | "tool.call" | "tool.result";
  text?: string;
  tool_id?: string;
  tool_name?: string;
  tool_input?: string;
  tool_is_error?: boolean;
}
