import { useUi, type NavKey } from "../../../shared/store/ui-store";
import { cn } from "../../../shared/lib/cn";

interface NavItem {
  key: NavKey;
  label: string;
  soon?: boolean;
  path: string;
}

/** studio-nav (spec shell §3.2): microcopy en español (RN-10) — «Producto», no "Product". */
const ITEMS: NavItem[] = [
  { key: "studio", label: "Studio", path: "M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" },
  { key: "backlog", label: "Backlog", path: "M9 6h11M9 12h11M9 18h11M4 6h.01M4 12h.01M4 18h.01" },
  { key: "producto", label: "Producto", path: "M3 3h7v7H3zM14 3h7v7h-7zM3 14h7v7H3zM14 14h7v7h-7z" },
  { key: "roles", label: "Roles", soon: true, path: "M9 11a3 3 0 1 0 0-6 3 3 0 0 0 0 6zM2.5 20c0-3.3 2.9-6 6.5-6s6.5 2.7 6.5 6M18 11.3a2.3 2.3 0 1 0 0-4.6M15.8 14.3c2.6.5 4.7 2.7 4.7 5.7" },
  { key: "config", label: "Config", path: "M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" },
];

export function StudioNav() {
  const nav = useUi((s) => s.nav);
  const setNav = useUi((s) => s.setNav);

  return (
    <nav className="flex w-16 shrink-0 flex-col items-stretch gap-1 border-r border-sidebar-border bg-sidebar px-1.5 py-3" aria-label="Navegación del estudio">
      {ITEMS.map((it) => (
        <button
          key={it.key}
          title={it.soon ? `${it.label} — pronto (llega con el registry, PB-25→PB-06)` : it.label}
          aria-current={nav === it.key ? "page" : undefined}
          onClick={() => setNav(nav === it.key && it.key !== "studio" ? "studio" : it.key)}
          className={cn(
            "relative flex cursor-pointer flex-col items-center gap-1 rounded-md px-1 py-2 text-[10px] font-semibold transition-colors",
            nav === it.key
              ? "bg-accent-soft text-primary"
              : "text-muted-foreground hover:bg-sidebar-accent hover:text-foreground",
            it.soon && "opacity-60",
          )}
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className="size-4.5">
            <path d={it.path} />
          </svg>
          {it.label}
          {it.soon && (
            <span className="absolute right-0.5 top-0.5 rounded-full bg-secondary px-1 font-mono text-[8px] uppercase tracking-wide text-muted-foreground">
              pronto
            </span>
          )}
        </button>
      ))}
    </nav>
  );
}
