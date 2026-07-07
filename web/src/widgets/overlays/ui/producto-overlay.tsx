import { useUi } from "../../../shared/store/ui-store";
import { Card, HonestBanner, PanelOverlay, Pill } from "../../../shared/ui";

/** Mapa ESTÁTICO (spec §3.6): refleja el INCREMENTO real (CAP-01..06 + esta rebanada) a mano.
 *  Se deriva del dato real cuando el SYSTEM-MAP as code entre (PB-08+). */
interface Area {
  nombre: string;
  estado: "live" | "wip" | "planned";
  destino?: string;
}
interface Caja {
  emoji: string;
  nombre: string;
  sub: string;
  areas: Area[];
}
interface Zona {
  nombre: string;
  tier: string;
  desc: string;
  cajas: Caja[];
}

const ZONAS: Zona[] = [
  {
    nombre: "Estudio",
    tier: "valor",
    desc: "Lo que el usuario opera directo: conversación con Claude Code + navegación multisesión.",
    cajas: [
      {
        emoji: "💬",
        nombre: "Dev Studio",
        sub: "conversación + composer + panel lateral",
        areas: [
          { nombre: "Chat conversacional", estado: "live" },
          { nombre: "Composer con cola (Encolar)", estado: "live" },
          { nombre: "Cambios (diff + commit por selección)", estado: "live" },
          { nombre: "Terminal integrada", estado: "planned", destino: "PB-09" },
          { nombre: "Pruebas (navegador + MCP)", estado: "planned", destino: "PB-09+" },
        ],
      },
      {
        emoji: "🗂",
        nombre: "Workspaces multisesión",
        sub: "workspace = sesión (1:1)",
        areas: [
          { nombre: "Crear/cerrar sesión desde la app", estado: "live" },
          { nombre: "Sesión ligada a historia (picker)", estado: "live" },
          { nombre: "Aislar sesión en workspace propio", estado: "planned", destino: "PB-02 — la siguiente" },
        ],
      },
    ],
  },
  {
    nombre: "Plataforma",
    tier: "transversal",
    desc: "El usuario entra, gestiona su espacio y su backlog — se opera a mano, no vía agente.",
    cajas: [
      {
        emoji: "📁",
        nombre: "Repositorios y Workspaces",
        sub: "rail izquierdo",
        areas: [
          { nombre: "Registro de repos (ruta local)", estado: "live" },
          { nombre: "Branch + diff ±N por workspace", estado: "live" },
          { nombre: "Clonar desde URL", estado: "planned" },
        ],
      },
      {
        emoji: "📋",
        nombre: "Backlog",
        sub: "board de estados del proceso",
        areas: [
          { nombre: "Picker de paquete de trabajo", estado: "live" },
          { nombre: "Board real (Historia/Capability as code)", estado: "planned", destino: "PB-07" },
        ],
      },
      {
        emoji: "⚙",
        nombre: "Configuración",
        sub: "roster de roles",
        areas: [
          { nombre: "Roster (vista previa del registry)", estado: "wip", destino: "PB-25→PB-06" },
          { nombre: "Registry de arneses propio", estado: "planned", destino: "PB-25" },
        ],
      },
    ],
  },
  {
    nombre: "Infraestructura",
    tier: "no-funcional",
    desc: "Boundaries as code + driver CLI-nativo + API local — arch/INDEX.md manda.",
    cajas: [
      {
        emoji: "🔧",
        nombre: "Motor",
        sub: "driver + API + persistencia",
        areas: [
          { nombre: "Driver CLI-nativo (subproceso claude)", estado: "live" },
          { nombre: "API REST + SSE local-only", estado: "live" },
          { nombre: "Git solo-lectura + commit (boundary)", estado: "live" },
          { nombre: "SSE replay Last-Event-ID", estado: "planned", destino: "PB-18" },
        ],
      },
    ],
  },
];

const pillTone = { live: "ok", wip: "warn", planned: "muted" } as const;
const pillLabel = { live: "🟢 live", wip: "⏳ wip", planned: "📋 planned" } as const;

export function ProductoOverlay() {
  const nav = useUi((s) => s.nav);
  const resetToStudio = useUi((s) => s.resetToStudio);

  return (
    <PanelOverlay
      open={nav === "producto"}
      onClose={resetToStudio}
      title="Producto"
      subtitle="Mapa del producto — la estructura prevista y lo que ya está construido sobre ella."
    >
      <div className="space-y-5 p-4">
        <HonestBanner>
          Mapa estático, reflejo manual del INCREMENTO — se deriva del dato real cuando el
          SYSTEM-MAP as code entre (PB-08+).
        </HonestBanner>
        {ZONAS.map((z) => (
          <section key={z.nombre}>
            <div className="mb-2 flex items-baseline gap-2">
              <h3 className="font-display text-base font-semibold">{z.nombre}</h3>
              <span className="rounded-full bg-secondary px-2 py-0.5 font-mono text-[9px] uppercase tracking-wider text-muted-foreground">
                {z.tier}
              </span>
              <p className="hidden text-xs text-muted-foreground md:block">{z.desc}</p>
            </div>
            <div className="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-3">
              {z.cajas.map((c) => {
                const live = c.areas.filter((a) => a.estado === "live").length;
                return (
                  <Card key={c.nombre} className="p-3">
                    <div className="flex items-center gap-2">
                      <span>{c.emoji}</span>
                      <div className="min-w-0 flex-1">
                        <div className="truncate text-sm font-semibold">{c.nombre}</div>
                        <div className="truncate text-[10px] text-muted-foreground">{c.sub}</div>
                      </div>
                      <span className="font-mono text-[10px] text-muted-foreground">
                        {live}/{c.areas.length}
                      </span>
                    </div>
                    <div className="mt-2 h-1 overflow-hidden rounded-full bg-secondary">
                      <div className="h-full bg-primary" style={{ width: `${(live / c.areas.length) * 100}%` }} />
                    </div>
                    <div className="mt-2 space-y-1">
                      {c.areas.map((a) => (
                        <div key={a.nombre} className="flex items-center gap-2 text-xs">
                          <span className="min-w-0 flex-1 truncate">{a.nombre}</span>
                          <Pill tone={pillTone[a.estado]} className="text-[9px]">
                            {pillLabel[a.estado]}
                          </Pill>
                          {a.destino && (
                            <span className="shrink-0 font-mono text-[9px] text-muted-foreground">→ {a.destino}</span>
                          )}
                        </div>
                      ))}
                    </div>
                  </Card>
                );
              })}
            </div>
          </section>
        ))}
      </div>
    </PanelOverlay>
  );
}
