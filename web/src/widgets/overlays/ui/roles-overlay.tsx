import { useUi } from "../../../shared/store/ui-store";
import { PanelOverlay, Pill } from "../../../shared/ui";

export function RolesOverlay() {
  const nav = useUi((s) => s.nav);
  const resetToStudio = useUi((s) => s.resetToStudio);

  return (
    <PanelOverlay
      open={nav === "roles"}
      onClose={resetToStudio}
      title="Roles"
      subtitle="El catálogo de arneses de tu organización, más allá de asignarlos a una sesión nueva."
    >
      <div className="flex h-full flex-col items-center justify-center gap-3 p-8 text-center">
        <div className="text-3xl">🧩</div>
        <Pill tone="muted">Próximamente</Pill>
        <p className="max-w-sm text-sm leading-relaxed text-muted-foreground">
          Esta vista será el navegador del <b className="text-foreground">registry de arneses</b>{" "}
          (PB-25 → PB-06): explorar los roles publicados por tu organización e instalarlos por
          proyecto. Hoy el roster solo se ve al crear una sesión, desde <b>Configuración</b>.
        </p>
      </div>
    </PanelOverlay>
  );
}
