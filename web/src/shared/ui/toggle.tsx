import { cn } from "../lib/cn";

/** Switch on/off (modales Comandos, ajustes). */
export function Toggle({
  on,
  onChange,
  disabled,
  "aria-label": ariaLabel,
}: {
  on: boolean;
  onChange?: (v: boolean) => void;
  disabled?: boolean;
  "aria-label": string;
}) {
  return (
    <button
      role="switch"
      aria-checked={on}
      aria-label={ariaLabel}
      disabled={disabled}
      onClick={() => onChange?.(!on)}
      className={cn(
        "relative h-5 w-9 shrink-0 cursor-pointer rounded-full border border-border transition-colors disabled:cursor-not-allowed disabled:opacity-40",
        on ? "bg-primary" : "bg-secondary",
      )}
    >
      <span
        className={cn(
          "absolute top-0.5 size-3.5 rounded-full bg-foreground transition-transform",
          on ? "translate-x-4.5 bg-primary-foreground" : "translate-x-0.5",
        )}
      />
    </button>
  );
}
