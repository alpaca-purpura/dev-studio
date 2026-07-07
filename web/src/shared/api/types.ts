export type SessionStatus = "streaming" | "idle";

export interface Turn {
  role: "user" | "assistant";
  text: string;
}

export interface Session {
  id: string;
  nombre: string;
  cwd: string;
  status: SessionStatus;
  claude_session_id?: string;
  model?: string;
  conv: Turn[];
}

export interface DockFrame {
  session_id: string;
  kind: "status" | "init" | "delta" | "result" | "error";
  text?: string;
  status?: SessionStatus;
  claude_session_id?: string;
  model?: string;
}
