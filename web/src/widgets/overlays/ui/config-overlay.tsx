import { useEffect, useMemo, useState } from "react";
import { useArneses } from "../../../shared/store/arneses-store";
import { useSessions } from "../../../shared/store/sessions-store";
import { useRepos } from "../../../shared/store/repos-store";
import { useUi } from "../../../shared/store/ui-store";
import { Button, Chip, FieldLabel, HonestBanner, Input, PanelOverlay, Select } from "../../../shared/ui";
import { ArnesCard, ArnesDetalle } from "./arnes-card";
import { AppUpdateSection } from "./app-update-section";

function slugify(s: string): string {
  return (
    s
      .toLowerCase()
      .normalize("NFD")
      .replace(/[̀-ͯ]/g, "")
      .replace(/\(.*?\)/g, "")
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/^-+|-+$/g, "")
      .slice(0, 26) || "sesion"
  );
}

export function ConfigOverlay() {
  const nav = useUi((s) => s.nav);
  const pendingHistoria = useUi((s) => s.pendingHistoria);
  const pendingRepoId = useUi((s) => s.pendingRepoId);
  const rolElegido = useUi((s) => s.rolElegido);
  const setRol = useUi((s) => s.setRol);
  const startPicker = useUi((s) => s.startPicker);
  const resetToStudio = useUi((s) => s.resetToStudio);
  const clearPending = useUi((s) => s.clearPending);
  const repos = useRepos((s) => s.repos);
  const create = useSessions((s) => s.create);

  const registry = useArneses((s) => s.registry);
  const instaladosPorRepo = useArneses((s) => s.instalados);
  const arnesesError = useArneses((s) => s.error);
  const arnesesBusy = useArneses((s) => s.busy);
  const initArneses = useArneses((s) => s.init);
  const conectar = useArneses((s) => s.conectar);
  const syncRegistry = useArneses((s) => s.sync);
  const refreshInstalados = useArneses((s) => s.refreshInstalados);
  const instalar = useArneses((s) => s.instalar);
  const desinstalar = useArneses((s) => s.desinstalar);

  const [busy, setBusy] = useState(false);
  const [source, setSource] = useState("");

  const abierto = nav === "config";
  const repoDestino = repos.find((r) => r.id === pendingRepoId) ?? repos[0];
  const instalados = useMemo(
    () => (repoDestino ? (instaladosPorRepo[repoDestino.id] ?? []) : []),
    [instaladosPorRepo, repoDestino],
  );
  const catalogo = registry?.arneses ?? [];
  const conectado = Boolean(registry?.source);

  useEffect(() => {
    if (!abierto) return;
    void initArneses();
    if (repoDestino) void refreshInstalados(repoDestino.id);
  }, [abierto, repoDestino, initArneses, refreshInstalados]);

  useEffect(() => {
    if (registry?.source) setSource(registry.source);
  }, [registry?.source]);

  // el rol elegido debe ser un arnés INSTALADO (RN-5); si no, cae al primero (o ninguno)
  const arnesElegido = instalados.find((a) => a.id === rolElegido) ?? instalados[0] ?? null;

  // roster agrupado por proceso (el meta rol×proceso de la fábrica — los grupos
  // Builder/Auditor/Humano del mock murieron con el mock)
  const grupos = useMemo(() => {
    const g = new Map<string, typeof instalados>();
    for (const a of instalados) {
      const k = a.proceso || "sin proceso declarado";
      g.set(k, [...(g.get(k) ?? []), a]);
    }
    return [...g.entries()];
  }, [instalados]);

  const crear = async () => {
    // RN-1: no hay sesión de trabajo sin paquete — sin historia elegida, al picker.
    if (!pendingHistoria) {
      if (repoDestino) startPicker(repoDestino.id);
      return;
    }
    if (!repoDestino || busy) return;
    setBusy(true);
    try {
      await create({
        nombre: slugify(pendingHistoria.titulo),
        repo_id: repoDestino.id,
        historia: pendingHistoria,
        rol: arnesElegido?.id ?? "",
        modo: "trabajo",
        ubicacion: "worktree",
      });
      clearPending();
      resetToStudio();
    } finally {
      setBusy(false);
    }
  };

  return (
    <PanelOverlay
      open={abierto}
      onClose={() => {
        clearPending();
        resetToStudio();
      }}
      title="Configuración"
      subtitle="Registry de arneses de tu organización + roster del proyecto (rol = arnés instalado, DH-14)."
    >
      <div className="space-y-4 p-4">
        {/* --- registry (nivel app) --- */}
        <section className="space-y-2 rounded-md border border-border bg-card p-3">
          <div className="font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">
            Registry de arneses
          </div>
          <div className="flex max-w-2xl gap-2">
            <Input
              placeholder="URL git del marketplace o ruta local (p. ej. git@github.com:org/marketplace)"
              value={source}
              onChange={(e) => setSource(e.target.value)}
              className="flex-1 font-mono text-xs"
            />
            <Button
              variant="outline"
              disabled={arnesesBusy || !source.trim()}
              onClick={() => void conectar(source.trim())}
            >
              {conectado ? "Reconectar" : "Conectar"}
            </Button>
            {conectado && (
              <Button variant="ghost" disabled={arnesesBusy} onClick={() => void syncRegistry()}>
                Actualizar
              </Button>
            )}
          </div>
          {arnesesError && <p className="text-xs text-warn">{arnesesError}</p>}
          {!conectado && !arnesesError && (
            <HonestBanner>
              Sin registry conectado no hay roles: DevStudio no crea roles locales — la curaduría
              vive en el marketplace git de tu organización (formato ArnesIA).
            </HonestBanner>
          )}
        </section>

        <div className="flex gap-4">
          {/* --- roster instalado + catálogo --- */}
          <nav className="w-64 shrink-0 space-y-4">
            <div>
              <div className="mb-1 font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">
                Roster del proyecto{repoDestino ? ` · ${repoDestino.nombre}` : ""}
              </div>
              {grupos.length === 0 && (
                <p className="rounded-md border border-dashed border-border p-2 text-xs leading-relaxed text-muted-foreground">
                  Ningún arnés instalado en este proyecto todavía — instalá uno del catálogo de
                  abajo. Una sesión también puede crearse sin rol.
                </p>
              )}
              {grupos.map(([proceso, lista]) => (
                <div key={proceso} className="mb-2">
                  <div className="mb-1 truncate text-[10px] text-muted-foreground" title={proceso}>
                    {proceso}
                  </div>
                  {lista.map((a) => (
                    <ArnesCard
                      key={a.id}
                      arnes={a}
                      selected={arnesElegido?.id === a.id}
                      onSelect={() => setRol(a.id)}
                    />
                  ))}
                </div>
              ))}
            </div>

            {conectado && (
              <div>
                <div className="mb-1 font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">
                  Catálogo del registry
                </div>
                {catalogo.length === 0 && (
                  <p className="text-xs text-muted-foreground">El marketplace no publica arneses aún.</p>
                )}
                {catalogo.map((a) => {
                  const instalado = instalados.some((i) => i.id === a.id);
                  return (
                    <ArnesCard
                      key={a.id}
                      arnes={a}
                      accion={
                        repoDestino && (
                          <Button
                            variant="ghost"
                            className="shrink-0 px-2 py-1 text-xs"
                            disabled={arnesesBusy}
                            onClick={() =>
                              void (instalado
                                ? desinstalar(repoDestino.id, a.id)
                                : instalar(repoDestino.id, a.id))
                            }
                          >
                            {instalado ? "Desinstalar" : "Instalar"}
                          </Button>
                        )
                      }
                    />
                  );
                })}
              </div>
            )}
          </nav>

          {/* --- detalle + crear sesión --- */}
          <div className="min-w-0 flex-1 space-y-3">
            {arnesElegido ? (
              <ArnesDetalle arnes={arnesElegido} />
            ) : (
              <div className="max-w-xl rounded-md border border-dashed border-border p-3 text-xs leading-relaxed text-muted-foreground">
                Sin rol elegido: la sesión nace <b className="text-foreground">sin arnés</b> (Claude
                Code pelado). Instalá un arnés del registry para trabajar con rol.
              </div>
            )}

            {pendingHistoria ? (
              <Chip className="border-primary/40 bg-accent-soft text-primary">
                📋 Nueva sesión para: {pendingHistoria.titulo}
              </Chip>
            ) : (
              <div className="max-w-xl rounded-md border border-dashed border-border p-3 text-xs leading-relaxed text-muted-foreground">
                Todavía no elegiste un paquete de trabajo. Toda sesión con edición liga a un ítem
                (RN-1) — <b className="text-foreground">Crear sesión aislada</b> te llevará ahí a
                elegirlo.
              </div>
            )}

            <div className="flex max-w-xl gap-3">
              <div className="flex-1">
                <FieldLabel>Proveedor</FieldLabel>
                <Select disabled title="Segundo proveedor — PB-20">
                  <option>Claude Code</option>
                </Select>
              </div>
              <div className="flex-1">
                <FieldLabel>Modelo</FieldLabel>
                <Select disabled title="Elección de modelo — sin efecto en esta rebanada (RN-7)">
                  <option>El de tu CLI</option>
                </Select>
              </div>
              <div className="flex-1">
                <FieldLabel>Cuenta</FieldLabel>
                <Select disabled title="BYO licencia (DH-10) — la app no toca credenciales">
                  <option>Personal (BYO CLI)</option>
                </Select>
              </div>
            </div>

            {repoDestino ? (
              <p className="text-xs text-muted-foreground">
                Repositorio destino: <b className="font-mono text-foreground">{repoDestino.nombre}</b>
                <span className="ml-2 text-[10px]">
                  — la sesión nace en su workspace aislado; el arnés se inyecta al claude por flags
                  (PB-25)
                </span>
              </p>
            ) : (
              <p className="text-xs text-warn">Registrá un repositorio primero (rail izquierdo).</p>
            )}

            <Button onClick={() => void crear()} disabled={!repoDestino || busy}>
              Crear sesión aislada
            </Button>
          </div>
        </div>

        {/* app-level: versión + updater (movidos del footer del rail — DH-18.1) */}
        <AppUpdateSection />
      </div>
    </PanelOverlay>
  );
}
