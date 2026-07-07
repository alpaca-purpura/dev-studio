import { cn } from "../lib/cn";

export interface TabItem {
  id: string;
  label: string;
}

/** Fila de tabs subrayadas (panel lateral, sub-tabs de Cambios). */
export function Tabs({
  items,
  active,
  onSelect,
  className,
}: {
  items: TabItem[];
  active: string;
  onSelect: (id: string) => void;
  className?: string;
}) {
  return (
    <div role="tablist" className={cn("flex items-center gap-1 border-b border-border", className)}>
      {items.map((t) => (
        <button
          key={t.id}
          role="tab"
          aria-selected={active === t.id}
          onClick={() => onSelect(t.id)}
          className={cn(
            "cursor-pointer border-b-2 px-3 py-2 text-sm transition-colors",
            active === t.id
              ? "border-primary font-semibold text-foreground"
              : "border-transparent text-muted-foreground hover:text-foreground",
          )}
        >
          {t.label}
        </button>
      ))}
    </div>
  );
}

/** Chip de filtro seleccionable (chips por rol del panel Cambios). */
export function FilterChip({
  active,
  onClick,
  children,
}: {
  active?: boolean;
  onClick?: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      aria-pressed={active}
      className={cn(
        "cursor-pointer rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors",
        active
          ? "border-primary bg-accent-soft text-primary"
          : "border-border bg-transparent text-muted-foreground hover:text-foreground",
      )}
    >
      {children}
    </button>
  );
}
