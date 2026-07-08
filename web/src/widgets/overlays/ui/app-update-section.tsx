import { useEffect, useState } from "react";
import { api } from "../../../shared/api/client";
import { Button, ModalShell } from "../../../shared/ui";

/** Sección «Aplicación» del overlay Configuración (movida del footer del rail —
 *  feedback dogfooding DH-18.1): versión del binario corriendo + «Actualizar» (PB-26:
 *  rebuild local → restart → esta SPA pollea la versión y recarga sola cuando cambia). */
export function AppUpdateSection() {
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
    <section className="space-y-2 rounded-md border border-border bg-card p-3">
      <div className="font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">
        Aplicación
      </div>
      <div className="flex max-w-2xl items-center gap-2">
        <span
          className="min-w-0 flex-1 truncate font-mono text-xs text-muted-foreground"
          title={source ? `Binario ${version} — rebuild desde ${source}` : `Versión ${version} (sin app.json: corré scripts/install.sh)`}
        >
          v {version}
          {source && <span className="ml-2 text-[10px]">— fuente: {source}</span>}
        </span>
        <Button
          variant="outline"
          onClick={() => void actualizar()}
          disabled={estado !== "idle"}
          title="Rebuild del repo local + reinicio (PB-26)"
        >
          {estado === "idle" ? "↻ Actualizar" : estado === "updating" ? "Compilando…" : "Reiniciando…"}
        </Button>
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
    </section>
  );
}
