import { useState } from "react";
import { TIPO_META, TIPOS, type TipoItem } from "../../../shared/lib/tipos";
import { useRepos } from "../../../shared/store/repos-store";
import { useSessions } from "../../../shared/store/sessions-store";
import { useUi } from "../../../shared/store/ui-store";
import { Button, FieldLabel, Input, ModalShell, Select } from "../../../shared/ui";
import { cn } from "../../../shared/lib/cn";

type Ubicacion = "worktree" | "checkout";
type Proposito = "explorar" | "existente" | "nuevo";

function OptionCard({
  selected,
  disabled,
  title,
  icon,
  desc,
  note,
  onClick,
}: {
  selected?: boolean;
  disabled?: boolean;
  title: string;
  icon: string;
  desc: string;
  note?: string;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      title={disabled ? note : undefined}
      className={cn(
        "flex w-full items-start gap-3 rounded-lg border p-3 text-left transition-colors",
        disabled
          ? "cursor-not-allowed border-border opacity-40"
          : selected
            ? "cursor-pointer border-primary bg-accent-soft"
            : "cursor-pointer border-border hover:border-primary/50",
      )}
    >
      <span className="text-lg" aria-hidden>{icon}</span>
      <span className="min-w-0">
        <span className="block text-sm font-semibold">{title}</span>
        <span className="mt-0.5 block text-xs leading-relaxed text-muted-foreground">{desc}</span>
        {disabled && note && <span className="mt-1 block text-[10px] text-warn">{note}</span>}
      </span>
    </button>
  );
}

/** Wizard «Nuevo Workspace» (PB-27, spec nuevo-workspace §2): ¿dónde? → ¿para qué? */
export function NewWorkspaceWizard() {
  const wizardRepoId = useUi((s) => s.wizardRepoId);
  const closeWizard = useUi((s) => s.closeWizard);
  const startPicker = useUi((s) => s.startPicker);
  const pickHistoria = useUi((s) => s.pickHistoria);
  const repo = useRepos((s) => s.repos.find((r) => r.id === wizardRepoId));
  const create = useSessions((s) => s.create);

  const [ubicacion, setUbicacion] = useState<Ubicacion>("worktree");
  const [proposito, setProposito] = useState<Proposito | null>(null);
  const [nombreExplorar, setNombreExplorar] = useState("explorar");
  const [tipoNuevo, setTipoNuevo] = useState<TipoItem>("historia");
  const [tituloNuevo, setTituloNuevo] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!wizardRepoId || !repo) return null;

  const reset = () => {
    setUbicacion("worktree");
    setProposito(null);
    setNombreExplorar("explorar");
    setTituloNuevo("");
    setError(null);
  };

  const cerrar = () => {
    reset();
    closeWizard();
  };

  const crearExploracion = async () => {
    if (busy) return;
    setBusy(true);
    setError(null);
    try {
      await create({
        nombre: nombreExplorar.trim() || "explorar",
        repo_id: repo.id,
        modo: "exploracion",
        ubicacion,
      });
      cerrar();
    } catch (e) {
      setError(e instanceof Error ? e.message : "No se pudo crear la sesión");
    } finally {
      setBusy(false);
    }
  };

  const crearItemNuevo = () => {
    const titulo = tituloNuevo.trim();
    if (!titulo) return;
    // el ítem viaja con la sesión (spec §2: board global = PB-07); id local estable
    pickHistoria({ id: `local-${Date.now().toString(36)}`, titulo, tipo: tipoNuevo });
    reset();
  };

  return (
    <ModalShell
      open
      onClose={cerrar}
      title="Nuevo Workspace"
      subtitle={`Repositorio: ${repo.nombre}`}
      icon="⎇"
    >
      <div className="space-y-4">
        <div>
          <FieldLabel>1 · ¿Dónde va a trabajar esta sesión?</FieldLabel>
          <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
            <OptionCard
              icon="⎇"
              title="Workspace nuevo"
              desc="Worktree aislado con su propia branch — para todo lo que toque código."
              selected={ubicacion === "worktree"}
              onClick={() => setUbicacion("worktree")}
            />
            <OptionCard
              icon="📂"
              title="Checkout actual"
              desc="La raíz del repo, sin aislamiento — solo para navegar y comprender."
              selected={ubicacion === "checkout"}
              onClick={() => {
                setUbicacion("checkout");
                setProposito("explorar");
              }}
            />
          </div>
        </div>

        <div>
          <FieldLabel>2 · ¿Para qué?</FieldLabel>
          <div className="space-y-2">
            <OptionCard
              icon="🔍"
              title="Navegar y comprender"
              desc="Sesión de exploración: Claude Code lee y analiza, sin permisos de edición (modo plan). No requiere ítem del backlog."
              selected={proposito === "explorar"}
              onClick={() => setProposito("explorar")}
            />
            <OptionCard
              icon="📋"
              title="Trabajar un ítem existente"
              desc="Elegir una historia, bug o tarea del Backlog — la sesión nace ligada (RN-1)."
              selected={proposito === "existente"}
              disabled={ubicacion === "checkout"}
              note="Trabajo con edición siempre en workspace aislado"
              onClick={() => startPicker(repo.id)}
            />
            <OptionCard
              icon="➕"
              title="Crear un ítem nuevo"
              desc="Historia, bug, hotfix, tarea o spike — nace acá y la sesión queda ligada."
              selected={proposito === "nuevo"}
              disabled={ubicacion === "checkout"}
              note="Trabajo con edición siempre en workspace aislado"
              onClick={() => setProposito("nuevo")}
            />
          </div>
        </div>

        {proposito === "explorar" && (
          <div className="space-y-2 rounded-md border border-border p-3">
            <FieldLabel>Nombre de la sesión</FieldLabel>
            <Input value={nombreExplorar} onChange={(e) => setNombreExplorar(e.target.value)} />
            <p className="text-[10px] leading-relaxed text-muted-foreground">
              {ubicacion === "checkout"
                ? "Corre sobre la raíz del repo, en solo-lectura — no crea worktree ni branch."
                : "Crea un worktree con branch explore/… — útil para revisar sin tocar tu checkout."}
            </p>
            <Button onClick={() => void crearExploracion()} disabled={busy}>
              🔍 Crear sesión de exploración
            </Button>
          </div>
        )}

        {proposito === "nuevo" && (
          <div className="space-y-2 rounded-md border border-border p-3">
            <div className="flex gap-2">
              <div className="w-40 shrink-0">
                <FieldLabel>Tipo</FieldLabel>
                <Select value={tipoNuevo} onChange={(e) => setTipoNuevo(e.target.value as TipoItem)}>
                  {TIPOS.map((t) => (
                    <option key={t} value={t}>
                      {TIPO_META[t].icon} {TIPO_META[t].label}
                    </option>
                  ))}
                </Select>
              </div>
              <div className="min-w-0 flex-1">
                <FieldLabel>Título</FieldLabel>
                <Input
                  value={tituloNuevo}
                  onChange={(e) => setTituloNuevo(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && crearItemNuevo()}
                  placeholder="Qué hay que hacer…"
                />
              </div>
            </div>
            <p className="text-[10px] text-muted-foreground">
              Branch: <code className="font-mono">{TIPO_META[tipoNuevo].branch}</code> · commit{" "}
              <code className="font-mono">{TIPO_META[tipoNuevo].commit}</code>
            </p>
            <Button onClick={crearItemNuevo} disabled={!tituloNuevo.trim()}>
              Continuar → elegir rol
            </Button>
          </div>
        )}

        {error && <p className="text-xs text-crit">{error}</p>}
      </div>
    </ModalShell>
  );
}
