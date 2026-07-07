import { useMemo, useState } from "react";
import { MOCK_BACKLOG, type MockStory } from "../../../shared/mock/backlog";
import { useUi } from "../../../shared/store/ui-store";
import { Chip, HonestBanner, PanelOverlay, StatePill, ReleaseChip, storyStateLabel, type StoryState } from "../../../shared/ui";
import { cn } from "../../../shared/lib/cn";

/** Los 8 estados del stream visible; parked/dropped se revelan con el toggle (F0: 10 completos). */
const STREAM: StoryState[] = ["idea", "refining", "refined", "ready", "developing", "developed", "reviewing", "done"];
const EXCEPCIONALES: StoryState[] = ["parked", "dropped"];

const prioColor: Record<MockStory["prioridad"], string> = {
  crítica: "bg-crit",
  alta: "bg-warn",
  media: "bg-[#3b82f6]",
  baja: "bg-muted-foreground",
};

function StoryCard({ story, picking, onPick }: { story: MockStory; picking: boolean; onPick: () => void }) {
  return (
    <button
      onClick={onPick}
      disabled={!picking}
      className={cn(
        "w-full rounded-md border border-border bg-card p-2.5 text-left transition-colors",
        picking ? "cursor-pointer hover:border-primary hover:shadow-[var(--shadow-glow)]" : "cursor-default",
        story.estado === "done" && "opacity-60",
      )}
    >
      <div className="flex items-center gap-1.5">
        <span className={cn("font-mono text-[10px] font-bold", story.rol ? "text-primary" : "text-muted-foreground")}>
          {story.rol ?? "— sin asignar"}
        </span>
        <span className="flex-1" />
        <ReleaseChip release={story.release} />
      </div>
      <div className="mt-1.5 text-xs font-semibold leading-snug">{story.titulo}</div>
      <div className="mt-2 flex items-center gap-1.5">
        <Chip className="max-w-36 truncate px-1.5 text-[9px]">{story.capability}</Chip>
        <span className="flex-1" />
        <span className={cn("size-1.5 rounded-full", prioColor[story.prioridad])} />
        <span className="text-[9px] text-muted-foreground">{story.prioridad}</span>
      </div>
    </button>
  );
}

export function BacklogOverlay() {
  const nav = useUi((s) => s.nav);
  const pickerMode = useUi((s) => s.pickerMode);
  const pickHistoria = useUi((s) => s.pickHistoria);
  const resetToStudio = useUi((s) => s.resetToStudio);
  const [showExcepcionales, setShowExcepcionales] = useState(false);
  const [search, setSearch] = useState("");

  const lanes = useMemo(() => {
    const visible = STREAM.concat(showExcepcionales ? EXCEPCIONALES : []);
    const q = search.toLowerCase();
    return visible.map((estado) => ({
      estado,
      stories: MOCK_BACKLOG.filter(
        (s) => s.estado === estado && (!q || s.titulo.toLowerCase().includes(q) || s.capability.includes(q)),
      ),
    }));
  }, [showExcepcionales, search]);

  // WIP advisory (F0: se visualizan, no restringen) — umbrales visuales 2/1
  const developing = MOCK_BACKLOG.filter((s) => s.estado === "developing").length;
  const reviewing = MOCK_BACKLOG.filter((s) => s.estado === "reviewing").length;

  return (
    <PanelOverlay
      open={nav === "backlog"}
      onClose={resetToStudio}
      title="Backlog"
      subtitle="Historias del proyecto — cada una liga a la capability que incrementa. Estados del descriptor de proceso (I-77), no un kanban genérico."
    >
      <div className="space-y-3 p-4">
        {pickerMode && (
          <div className="rounded-md border border-primary/40 bg-accent-soft px-3 py-2 text-sm font-semibold text-primary">
            Elegí un paquete de trabajo para la nueva sesión — hacé clic en una historia.
          </div>
        )}
        <HonestBanner>
          Datos de ejemplo — el board real llega con el port de Historia/Capability (PB-07).
          El picker sí es real: la historia elegida queda ligada a la sesión (RN-1).
        </HonestBanner>

        <div className="flex flex-wrap items-center gap-2">
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="🔍 Buscar…"
            className="w-48 rounded-md border border-input bg-card px-2 py-1 text-xs outline-none focus:border-primary"
          />
          <label className="flex cursor-pointer items-center gap-1.5 text-xs text-muted-foreground">
            <input
              type="checkbox"
              checked={showExcepcionales}
              onChange={(e) => setShowExcepcionales(e.target.checked)}
              className="accent-[var(--primary)]"
            />
            mostrar pausadas + descartadas
          </label>
          <span className="flex-1" />
          <Chip className={cn(developing > 2 && "border-crit text-crit")} title="WIP advisory (F0): se visualiza, no restringe">
            desarrollo {developing}/2
          </Chip>
          <Chip className={cn(reviewing > 1 && "border-crit text-crit")} title="WIP advisory (F0): se visualiza, no restringe">
            revisión {reviewing}/1
          </Chip>
        </div>

        <div className="overflow-x-auto pb-2">
          <div className="flex min-w-max gap-3">
            {lanes.map(({ estado, stories }) => (
              <div key={estado} className="w-60 shrink-0">
                <div className="mb-2 flex items-center gap-2">
                  <StatePill state={estado} />
                  <span className="font-mono text-[10px] text-muted-foreground">{stories.length}</span>
                </div>
                <div className="space-y-2">
                  {stories.map((s) => (
                    <StoryCard
                      key={s.id}
                      story={s}
                      picking={pickerMode}
                      onPick={() => pickerMode && pickHistoria({ id: s.id, titulo: s.titulo })}
                    />
                  ))}
                  {stories.length === 0 && (
                    <p className="rounded-md border border-dashed border-border p-3 text-center text-[10px] text-muted-foreground">
                      Sin historias en {storyStateLabel[estado].toLowerCase()}
                    </p>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </PanelOverlay>
  );
}
