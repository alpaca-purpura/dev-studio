import { cn } from "../lib/cn";

/** Avatar de rol: iniciales sobre acento suave (FS, QA, PO…). */
export function Avatar({ label, className }: { label: string; className?: string }) {
  return (
    <span
      className={cn(
        "inline-flex size-8 shrink-0 items-center justify-center rounded-full bg-accent-soft font-mono text-xs font-bold uppercase text-primary",
        className,
      )}
    >
      {label.slice(0, 2)}
    </span>
  );
}
