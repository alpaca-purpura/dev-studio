import { useEffect, useState } from "react";
import { api } from "../../../shared/api/client";
import { Button, ModalShell } from "../../../shared/ui";

/** Footer del rail (PB-26): versión del binario corriendo + «Actualizar» (rebuild local →
 *  restart → esta SPA pollea la versión y recarga sola cuando cambia — RN-5). */
export function UpdateFooter() {
  const [version, setVersion] = useState<string>("…");
  const [stamp, setStamp] = useState<string>(""); // version@build_date — cambia en CADA rebuild
  const [source, setSource] = useState<string>("");
  const [estado, setEstado] = useState<"idle" | "updating" | "waiting">("idle");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api
      .version()
      .then((v) => {
        setVersion(v.version);
        setStamp(`${v.version}@${v.build_date}`);
        setSource(v.source);
      })
      .catch(() => setVersion("?"));
  }, []);

  const actualizar = async () => {
    setEstado("updating");
    setError(null);
    try {
      await api.update();
      setEstado("waiting");
      // el proceso se reinicia: pollear hasta ver OTRO build (version@fecha) y recargar
      const original = stamp;
      const t0 = Date.now();
      const poll = async () => {
        if (Date.now() - t0 > 120_000) {
          setEstado("idle");
          setError("El reinicio tardó demasiado — recargá a mano.");
          return;
        }
        try {
          const v = await api.version();
          if (`${v.version}@${v.build_date}` !== original) {
            location.reload();
            return;
          }
        } catch {
          // reiniciando — seguir intentando
        }
        setTimeout(() => void poll(), 1000);
      };
      setTimeout(() => void poll(), 1000);
    } catch (e) {
      setEstado("idle");
      setError(e instanceof Error ? e.message : "No se pudo actualizar");
    }
  };

  return (
    <div className="border-t border-sidebar-border p-2">
      <div className="flex items-center gap-2">
        <span
          className="min-w-0 flex-1 truncate font-mono text-[9px] uppercase tracking-wider text-muted-foreground"
          title={source ? `Binario ${version} — rebuild desde ${source}` : `Versión ${version} (sin app.json: corré scripts/install.sh)`}
        >
          v {version}
        </span>
        <button
          onClick={() => void actualizar()}
          disabled={estado !== "idle"}
          title="Rebuild del repo local + reinicio (PB-26)"
          className="cursor-pointer rounded-full border border-border px-2 py-0.5 text-[10px] font-semibold text-muted-foreground transition-colors hover:border-primary hover:text-primary disabled:cursor-wait disabled:opacity-60"
        >
          {estado === "idle" ? "↻ Actualizar" : estado === "updating" ? "Compilando…" : "Reiniciando…"}
        </button>
      </div>

      <ModalShell
        open={error !== null}
        onClose={() => setError(null)}
        title="La actualización falló"
        subtitle="Salida real del build — nada se ocultó (RN de honestidad)."
        icon="⚠"
        wide
        footer={<Button onClick={() => setError(null)}>Entendido</Button>}
      >
        <pre className="overflow-x-auto whitespace-pre-wrap font-mono text-xs leading-relaxed text-crit">
          {error}
        </pre>
      </ModalShell>
    </div>
  );
}
