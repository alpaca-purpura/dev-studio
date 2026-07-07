import type { Session } from "./types";

async function json<T>(res: Response): Promise<T> {
  if (!res.ok) throw new Error(`${res.status} ${await res.text()}`);
  return (await res.json()) as T;
}

export const api = {
  list: (): Promise<Session[]> => fetch("/api/sessions").then((r) => json<Session[]>(r)),

  create: (nombre: string, cwd: string): Promise<Session> =>
    fetch("/api/sessions", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ nombre, cwd }),
    }).then((r) => json<Session>(r)),

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
    fetch(`/api/sessions/${id}/turn`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text }),
    }).then((r) => {
      if (!r.ok) throw new Error(`${r.status}`);
    }),
};
