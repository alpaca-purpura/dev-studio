import { useState } from "react";
import { cn } from "../../../shared/lib/cn";
import { Pip } from "../../../shared/ui/indicators";
import { useSessions } from "../../../shared/store/sessions-store";
import type { Session } from "../../../shared/api/types";

// SessionRail: rail izquierdo estilo WARP-tabs, colapsable a un gutter de solo iconos.
// Estilo + comportamiento copiados de harness-studio `widgets/session-rail` (mockup it.14);
// el dominio de la card se simplificó a lo que dev-studio ya tiene (sin arnés/salud/vistas —
// eso vive en ArnesIA, no acá; "todo lo demás dejalo en blanco").
export function SessionRail() {
  const sessions = useSessions((s) => s.sessions);
  const activeId = useSessions((s) => s.activeId);
  const collapsed = useSessions((s) => s.railCollapsed);
  const switchTo = useSessions((s) => s.switchTo);
  const toggleRail = useSessions((s) => s.toggleRail);
  const create = useSessions((s) => s.create);

  return (
    <aside
      className={cn(
        "flex flex-none flex-col border-r border-border bg-card transition-[width] duration-150",
        collapsed ? "w-[52px]" : "w-[224px]",
      )}
    >
      <div
        className={cn(
          "flex items-center gap-2 px-2.5 pb-2 pt-2.5",
          collapsed && "flex-col gap-1.5 px-0",
        )}
      >
        <div className="grid size-7 flex-none place-items-center rounded-lg bg-primary text-sm font-extrabold text-primary-foreground">
          D
        </div>
        {!collapsed && (
          <span className="text-[13px] font-bold tracking-tight text-foreground">DevStudio</span>
        )}
        <button
          type="button"
          onClick={toggleRail}
          title={collapsed ? "Expandir rail" : "Colapsar a gutter"}
          className={cn(
            "rounded-md px-1.5 py-0.5 text-muted-foreground hover:bg-secondary hover:text-foreground",
            !collapsed && "ml-auto",
          )}
        >
          {collapsed ? "»" : "«"}
        </button>
      </div>

      {!collapsed && (
        <div className="flex items-center gap-1.5 px-3 pb-1 pt-1.5 font-mono text-[9px] uppercase tracking-wider text-muted-foreground">
          Sesiones
        </div>
      )}

      <div
        className={cn(
          "flex flex-1 flex-col gap-1 overflow-y-auto p-2",
          collapsed && "items-center p-1",
        )}
      >
        {sessions.map((s) => (
          <SessionCard
            key={s.id}
            session={s}
            active={s.id === activeId}
            collapsed={collapsed}
            onClick={() => switchTo(s.id)}
          />
        ))}
      </div>

      <div className="flex flex-none flex-col gap-1 border-t border-border p-2">
        <NewSessionButton collapsed={collapsed} onCreate={create} />
      </div>
    </aside>
  );
}

function SessionCard({
  session: s,
  active,
  collapsed,
  onClick,
}: {
  session: Session;
  active: boolean;
  collapsed: boolean;
  onClick: () => void;
}) {
  const rename = useSessions((st) => st.rename);
  const closeSession = useSessions((st) => st.closeSession);
  const canClose = useSessions((st) => st.sessions.length > 1);
  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState(s.nombre);

  const commit = () => {
    setEditing(false);
    if (value.trim() && value.trim() !== s.nombre) void rename(s.id, value.trim());
    else setValue(s.nombre);
  };

  if (collapsed) {
    return (
      <button
        type="button"
        onClick={onClick}
        title={s.nombre}
        aria-current={active}
        className={cn(
          "flex w-[38px] flex-col items-center gap-1 rounded-[10px] border py-1.5",
          active ? "border-primary bg-accent-soft" : "border-border hover:bg-secondary",
        )}
      >
        <Pip status={s.status} />
      </button>
    );
  }

  return (
    <div
      role="button"
      tabIndex={0}
      onClick={onClick}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          onClick();
        }
      }}
      aria-current={active}
      title={s.nombre}
      className={cn(
        "group relative flex cursor-pointer flex-col gap-1 rounded-[9px] border p-2 text-left transition-colors",
        active
          ? "border-primary bg-accent-soft shadow-[inset_3px_0_0_var(--primary)]"
          : "border-border bg-card hover:border-input hover:bg-secondary",
      )}
    >
      {canClose && (
        <button
          type="button"
          title="Cerrar sesión"
          onClick={(e) => {
            e.stopPropagation();
            void closeSession(s.id);
          }}
          className="absolute right-1.5 top-1.5 hidden size-4 items-center justify-center rounded text-muted-foreground hover:bg-border hover:text-foreground group-hover:flex"
        >
          ✕
        </button>
      )}

      <div className="flex items-center gap-1.5">
        <Pip status={s.status} />
        {editing ? (
          <input
            autoFocus
            value={value}
            onClick={(e) => e.stopPropagation()}
            onChange={(e) => setValue(e.target.value)}
            onBlur={commit}
            onKeyDown={(e) => {
              e.stopPropagation();
              if (e.key === "Enter") commit();
              if (e.key === "Escape") {
                setValue(s.nombre);
                setEditing(false);
              }
            }}
            className="min-w-0 flex-1 rounded border border-primary bg-card px-1.5 py-px text-xs font-semibold text-foreground"
          />
        ) : (
          <span className="flex-1 truncate text-xs font-semibold">{s.nombre}</span>
        )}
        <button
          type="button"
          title="Renombrar sesión"
          onClick={(e) => {
            e.stopPropagation();
            setValue(s.nombre);
            setEditing(true);
          }}
          className="hidden size-[15px] items-center justify-center rounded text-[10px] text-muted-foreground hover:bg-border hover:text-foreground group-hover:flex"
        >
          ✎
        </button>
      </div>

      <div className="truncate font-mono text-[9.5px] text-muted-foreground">{s.cwd}</div>
    </div>
  );
}

function NewSessionButton({
  collapsed,
  onCreate,
}: {
  collapsed: boolean;
  onCreate: (cwd: string) => Promise<void>;
}) {
  return (
    <button
      type="button"
      onClick={() => {
        const dir = window.prompt("Directorio de trabajo de la sesión (opcional):")?.trim() ?? "";
        void onCreate(dir);
      }}
      title="Nueva sesión — spawnea un proceso claude propio"
      className={cn(
        "flex items-center justify-center gap-1.5 rounded-lg border border-dashed border-input py-2 text-xs font-semibold text-muted-foreground hover:border-primary hover:bg-accent-soft hover:text-primary",
        collapsed && "px-0",
      )}
    >
      ＋{!collapsed && <span>Nueva sesión</span>}
    </button>
  );
}
