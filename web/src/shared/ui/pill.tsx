import type { HTMLAttributes } from "react";
import { cn } from "../lib/cn";

type Tone = "ok" | "warn" | "crit" | "muted" | "primary";

export interface PillProps extends HTMLAttributes<HTMLSpanElement> {
  tone?: Tone;
}

const tones: Record<Tone, string> = {
  ok: "bg-ok-soft text-ok",
  warn: "bg-warn-soft text-warn",
  crit: "bg-crit-soft text-crit",
  muted: "bg-muted-soft text-muted-foreground",
  primary: "bg-accent-soft text-primary",
};

/** Pill de estado (live/planned/wip, severidades) — fondo soft, texto del tono. */
export function Pill({ tone = "muted", className, ...rest }: PillProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-semibold",
        tones[tone],
        className,
      )}
      {...rest}
    />
  );
}
