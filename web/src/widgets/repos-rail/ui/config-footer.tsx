import { useUi } from "../../../shared/store/ui-store";
import { cn } from "../../../shared/lib/cn";

/** Footer del rail: la puerta a Configuración (feedback dogfooding DH-18.1 — Config es
 *  de TODA la app, no de la zona por sesión del studio-nav; patrón ArnesIA: abajo del
 *  primer sidebar). Versión + «Actualizar» viven DENTRO de Config (§ Aplicación). */
export function ConfigFooter() {
  const nav = useUi((s) => s.nav);
  const setNav = useUi((s) => s.setNav);
  const activo = nav === "config";

  return (
    <div className="border-t border-sidebar-border p-2">
      <button
        onClick={() => setNav(activo ? "studio" : "config")}
        aria-current={activo ? "page" : undefined}
        title="Configuración de la aplicación: registry de arneses · versión y actualización"
        className={cn(
          "flex w-full cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-xs font-semibold transition-colors",
          activo
            ? "bg-accent-soft text-primary"
            : "text-muted-foreground hover:bg-sidebar-accent hover:text-foreground",
        )}
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className="size-4 shrink-0">
          <path d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" />
          <path d="M19.4 15a1.7 1.7 0 0 0 .34 1.87l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.7 1.7 0 0 0-1.87-.34 1.7 1.7 0 0 0-1.03 1.56V21a2 2 0 1 1-4 0v-.09a1.7 1.7 0 0 0-1.11-1.56 1.7 1.7 0 0 0-1.87.34l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.7 1.7 0 0 0 .34-1.87 1.7 1.7 0 0 0-1.56-1.03H3a2 2 0 1 1 0-4h.09a1.7 1.7 0 0 0 1.56-1.11 1.7 1.7 0 0 0-.34-1.87l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.7 1.7 0 0 0 1.87.34h.01A1.7 1.7 0 0 0 10 4.09V4a2 2 0 1 1 4 0v.09a1.7 1.7 0 0 0 1.03 1.56 1.7 1.7 0 0 0 1.87-.34l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.7 1.7 0 0 0-.34 1.87v.01c.26.63.87 1.04 1.56 1.04H21a2 2 0 1 1 0 4h-.09a1.7 1.7 0 0 0-1.51 1z" />
        </svg>
        Configuración
      </button>
    </div>
  );
}
