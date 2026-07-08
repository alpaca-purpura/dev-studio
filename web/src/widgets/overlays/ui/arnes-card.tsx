import type { Arnes } from "../../../shared/api/types";
import { Avatar, Pill } from "../../../shared/ui";
import { cn } from "../../../shared/lib/cn";

/** Sigla derivada del id/nombre del arnés (2 letras) — el formato ArnesIA no trae avatar. */
export function siglaDeArnes(a: Arnes): string {
  const base = (a.nombre || a.id).trim();
  const partes = base.split(/[\s·\-_/]+/).filter(Boolean);
  const s = partes.length >= 2 ? partes[0][0] + partes[1][0] : base.slice(0, 2);
  return s.toUpperCase();
}

/** Tarjeta de arnés — se usa en el roster instalado y en el catálogo del registry (PB-25). */
export function ArnesCard({
  arnes,
  selected,
  onSelect,
  accion,
}: {
  arnes: Arnes;
  selected?: boolean;
  onSelect?: () => void;
  accion?: React.ReactNode;
}) {
  return (
    <div
      className={cn(
        "flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm transition-colors",
        onSelect && "cursor-pointer",
        selected
          ? "bg-accent-soft font-semibold text-primary"
          : "text-muted-foreground hover:bg-secondary hover:text-foreground",
      )}
      onClick={onSelect}
      role={onSelect ? "button" : undefined}
    >
      <Avatar label={siglaDeArnes(arnes)} className="size-6 shrink-0 text-[8px]" />
      <div className="min-w-0 flex-1">
        <div className="truncate">{arnes.nombre || arnes.id}</div>
        <div className="truncate font-mono text-[10px] font-normal text-muted-foreground">
          v{arnes.version}
          {arnes.canal ? ` · ${arnes.canal}` : ""}
        </div>
      </div>
      {accion}
    </div>
  );
}

/** Detalle del arnés elegido: el meta rol×proceso×fases del manifiesto (arnes.l0.json). */
export function ArnesDetalle({ arnes }: { arnes: Arnes }) {
  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2">
        <h3 className="font-display text-lg font-semibold">{arnes.nombre || arnes.id}</h3>
        <Pill tone="primary">Arnés instalado</Pill>
        <span className="font-mono text-[10px] text-muted-foreground">
          v{arnes.version}
          {arnes.canal ? ` · ${arnes.canal}` : ""}
        </span>
      </div>
      {arnes.descripcion && (
        <p className="max-w-xl text-sm leading-relaxed text-muted-foreground">{arnes.descripcion}</p>
      )}
      <dl className="max-w-xl space-y-1 rounded-md border border-border bg-card p-3 text-xs leading-relaxed">
        <div>
          <dt className="inline font-mono font-bold uppercase tracking-[0.14em] text-muted-foreground">rol · </dt>
          <dd className="inline text-foreground">{arnes.rol}</dd>
        </div>
        {arnes.proceso && (
          <div>
            <dt className="inline font-mono font-bold uppercase tracking-[0.14em] text-muted-foreground">proceso · </dt>
            <dd className="inline text-foreground">{arnes.proceso}</dd>
          </div>
        )}
        {arnes.fases && arnes.fases.length > 0 && (
          <div>
            <dt className="inline font-mono font-bold uppercase tracking-[0.14em] text-muted-foreground">fases · </dt>
            <dd className="inline text-foreground">{arnes.fases.join(" → ")}</dd>
          </div>
        )}
      </dl>
      <p className="max-w-xl text-[11px] leading-relaxed text-muted-foreground">
        La sesión con este rol carga su forma-plugin (skills del arnés) vía{" "}
        <code className="font-mono">--plugin-dir</code> + su banda Base + un preámbulo compuesto de
        este meta — el payload no se escribe en tu repo, solo el lock{" "}
        <code className="font-mono">.devstudio/arneses.yaml</code>.
      </p>
    </div>
  );
}
