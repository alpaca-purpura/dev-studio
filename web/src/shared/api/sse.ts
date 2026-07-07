import type { DockFrame } from "./types";

/** Un único EventSource para todas las sesiones — el filtrado por session_id vive en el store. */
export function connectDock(onFrame: (f: DockFrame) => void): () => void {
  const es = new EventSource("/events");
  es.addEventListener("dock", (e) => {
    onFrame(JSON.parse((e as MessageEvent).data) as DockFrame);
  });
  return () => es.close();
}
