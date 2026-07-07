import type { GitDiff, GitLogEntry, GitStatus, Historia, Repo, Session } from "./types";

async function json<T>(res: Response): Promise<T> {
  if (!res.ok) throw new Error(`${res.status} ${await res.text()}`);
  return (await res.json()) as T;
}

const post = (url: string, body: unknown) =>
  fetch(url, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });

export const api = {
  list: (): Promise<Session[]> => fetch("/api/sessions").then((r) => json<Session[]>(r)),

  create: (p: { nombre: string; repo_id?: string; historia?: Historia; rol?: string; cwd?: string }): Promise<Session> =>
    post("/api/sessions", p).then((r) => json<Session>(r)),

  rename: (id: string, nombre: string): Promise<Session> =>
    fetch(`/api/sessions/${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ nombre }),
    }).then((r) => json<Session>(r)),

  close: (id: string): Promise<void> =>
    fetch(`/api/sessions/${id}`, { method: "DELETE" }).then((r) => {
      if (!r.ok) throw new Error(`${r.status}`);
    }),

  turn: (id: string, text: string): Promise<void> =>
    post(`/api/sessions/${id}/turn`, { text }).then((r) => {
      if (!r.ok) throw new Error(`${r.status}`);
    }),

  // --- repos ---
  repos: (): Promise<Repo[]> => fetch("/api/repos").then((r) => json<Repo[]>(r)),

  registerRepo: (ruta: string): Promise<Repo> => post("/api/repos", { ruta }).then((r) => json<Repo>(r)),

  removeRepo: (id: string): Promise<void> =>
    fetch(`/api/repos/${id}`, { method: "DELETE" }).then((r) => {
      if (!r.ok) throw new Error(`${r.status}`);
    }),

  repoGitStatus: (id: string): Promise<GitStatus> =>
    fetch(`/api/repos/${id}/git/status`).then((r) => json<GitStatus>(r)),

  // --- git de la sesión ---
  gitStatus: (sessionId: string): Promise<GitStatus> =>
    fetch(`/api/sessions/${sessionId}/git/status`).then((r) => json<GitStatus>(r)),

  gitDiff: (sessionId: string, path: string): Promise<GitDiff> =>
    fetch(`/api/sessions/${sessionId}/git/diff?path=${encodeURIComponent(path)}`).then((r) => json<GitDiff>(r)),

  gitLog: (sessionId: string, n = 30): Promise<GitLogEntry[]> =>
    fetch(`/api/sessions/${sessionId}/git/log?n=${n}`).then((r) => json<GitLogEntry[]>(r)),

  gitCommit: (sessionId: string, paths: string[], mensaje: string): Promise<{ sha: string }> =>
    post(`/api/sessions/${sessionId}/git/commit`, { paths, mensaje }).then((r) => json<{ sha: string }>(r)),
};
