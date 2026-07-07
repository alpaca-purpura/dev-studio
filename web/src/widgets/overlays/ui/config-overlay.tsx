import { useState } from "react";
import { MOCK_ROSTER } from "../../../shared/mock/backlog";
import { useSessions } from "../../../shared/store/sessions-store";
import { useRepos } from "../../../shared/store/repos-store";
import { useUi } from "../../../shared/store/ui-store";
import { Avatar, Button, Chip, FieldLabel, HonestBanner, PanelOverlay, Pill, Select } from "../../../shared/ui";
import { cn } from "../../../shared/lib/cn";

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
  const [busy, setBusy] = useState(false);

  const rol = MOCK_ROSTER.find((r) => r.nombre === rolElegido) ?? MOCK_ROSTER[0];
  const grupos = ["Builder", "Auditor", "Humano-complementario"] as const;
  const repoDestino = repos.find((r) => r.id === pendingRepoId) ?? repos[0];

  const crear = async () => {
    // RN-1: no hay sesión sin paquete — sin historia elegida, al Backlog en modo picker.
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
        rol: rol.nombre,
      });
      clearPending();
      resetToStudio();
    } finally {
      setBusy(false);
    }
  };

  return (
    <PanelOverlay
      open={nav === "config"}
      onClose={() => {
        clearPending();
        resetToStudio();
      }}
      title="Configuración"
      subtitle="Ajustes del proyecto. De momento: el roster de roles disponibles para nuevas sesiones."
    >
      <div className="space-y-4 p-4">
        <HonestBanner>
          Los roles vienen del registry de arneses de tu organización — registry propio en
          construcción (PB-25). Este roster es una vista previa estática (DH-14: rol = arnés
          instalado; cero roles locales).
        </HonestBanner>

        <div className="flex gap-4">
          {/* roster */}
          <nav className="w-56 shrink-0 space-y-3">
            {grupos.map((g) => (
              <div key={g}>
                <div className="mb-1 font-mono text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">
                  {g}
                </div>
                {MOCK_ROSTER.filter((r) => r.grupo === g).map((r) => (
                  <button
                    key={r.sigla}
                    onClick={() => setRol(r.nombre)}
                    className={cn(
                      "flex w-full cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm transition-colors",
                      rol.nombre === r.nombre
                        ? "bg-accent-soft font-semibold text-primary"
                        : "text-muted-foreground hover:bg-secondary hover:text-foreground",
                    )}
                  >
                    <Avatar label={r.sigla} className="size-6 text-[8px]" />
                    {r.nombre}
                  </button>
                ))}
              </div>
            ))}
            {/* «+ Rol personalizado» ELIMINADO a propósito (spec §3.5): rol custom = publicarlo al registry */}
          </nav>

          {/* detalle del rol */}
          <div className="min-w-0 flex-1 space-y-3">
            <div className="flex items-center gap-2">
              <h3 className="font-display text-lg font-semibold">{rol.nombre}</h3>
              <Pill tone="primary">Rol del roster</Pill>
            </div>
            <p className="max-w-xl text-sm leading-relaxed text-muted-foreground">{rol.descripcion}</p>

            {pendingHistoria ? (
              <Chip className="border-primary/40 bg-accent-soft text-primary">
                📋 Nueva sesión para: {pendingHistoria.titulo}
              </Chip>
            ) : (
              <div className="max-w-xl rounded-md border border-dashed border-border p-3 text-xs leading-relaxed text-muted-foreground">
                Todavía no elegiste un paquete de trabajo. Toda sesión debe ligarse a una historia
                del Backlog (RN-1) — <b className="text-foreground">Crear sesión aislada</b> te
                llevará ahí a elegirla.
              </div>
            )}

            <pre className="max-w-xl overflow-x-auto whitespace-pre-wrap rounded-md border border-border bg-card p-3 font-mono text-xs leading-relaxed text-muted-foreground">
              {rol.prompt}
            </pre>

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
                  (v1: la sesión trabaja sobre la raíz del repo — el aislamiento por worktree llega con PB-02)
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
      </div>
    </PanelOverlay>
  );
}
