import { cn } from "../lib/cn";

/** Los 10 estados del value_stream firmado en F0 (parked/dropped = excepcionales). */
export type StoryState =
  | "idea"
  | "refining"
  | "refined"
  | "ready"
  | "developing"
  | "developed"
  | "reviewing"
  | "done"
  | "parked"
  | "dropped";

export const storyStateLabel: Record<StoryState, string> = {
  idea: "Idea",
  refining: "Refinando",
  refined: "Refinada",
  ready: "Lista",
  developing: "En desarrollo",
  developed: "Desarrollada",
  reviewing: "En revisión",
  done: "Hecha",
  parked: "Pausada",
  dropped: "Descartada",
};

const styles: Record<StoryState, string> = {
  idea: "bg-muted-soft text-muted-foreground",
  refining: "bg-[hsl(294_30%_16%)] text-[hsl(294_55%_78%)]",
  refined: "bg-[hsl(258_30%_18%)] text-[hsl(258_60%_80%)]",
  ready: "bg-accent-soft text-primary",
  developing: "bg-warn-soft text-warn",
  developed: "bg-[hsl(200_40%_16%)] text-[hsl(200_70%_78%)]",
  reviewing: "bg-crit-soft text-crit",
  done: "bg-ok-soft text-ok",
  parked: "bg-muted-soft text-muted-foreground",
  dropped: "bg-muted-soft text-muted-foreground line-through",
};

export function StatePill({ state, className }: { state: StoryState; className?: string }) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2 py-0.5 font-mono text-[10px] font-bold uppercase tracking-wider",
        styles[state],
        className,
      )}
    >
      {storyStateLabel[state]}
    </span>
  );
}

/** Chip de release (F1/F2/FN) de una tarjeta de historia. */
export function ReleaseChip({ release }: { release: string }) {
  return (
    <span className="inline-flex items-center rounded-sm bg-[hsl(200_40%_16%)] px-1.5 py-0.5 font-mono text-[10px] font-bold text-[hsl(200_70%_78%)]">
      {release}
    </span>
  );
}
