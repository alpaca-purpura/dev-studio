export type SessionStatus = "streaming" | "idle";

export interface Turn {
  role: "user" | "assistant";
  text: string;
}

export interface Historia {
  id: string;
  titulo: string;
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
  workspace?: string; // ruta del worktree propio (PB-02); vacío = legacy
  branch?: string; // wt/{slug}
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
  kind: "status" | "init" | "delta" | "result" | "error";
  text?: string;
  status?: SessionStatus;
  claude_session_id?: string;
  model?: string;
}
