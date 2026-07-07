import { useEffect } from "react";
import { ReposRail } from "./widgets/repos-rail";
import { StudioNav } from "./widgets/studio-nav";
import { SessionView } from "./widgets/session-view";
import { ChangesPanel } from "./widgets/changes-panel";
import { BacklogOverlay, ConfigOverlay, ProductoOverlay, RolesOverlay } from "./widgets/overlays";
import { CloseSessionModal } from "./widgets/session-close/ui/close-session-modal";
import { NewWorkspaceWizard } from "./widgets/new-workspace/ui/new-workspace-wizard";
import { useSessions } from "./shared/store/sessions-store";
import { useRepos } from "./shared/store/repos-store";
import { useUi } from "./shared/store/ui-store";

export function App() {
  const initSessions = useSessions((s) => s.init);
  const initRepos = useRepos((s) => s.init);
  const nav = useUi((s) => s.nav);
  const resetToStudio = useUi((s) => s.resetToStudio);

  useEffect(() => {
    void initSessions();
    void initRepos();
  }, [initSessions, initRepos]);

  // Esc cierra el overlay activo y vuelve a Studio (spec §4.1 paso 5)
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape" && nav !== "studio") resetToStudio();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [nav, resetToStudio]);

  return (
    <div className="flex h-full">
      <ReposRail />
      <StudioNav />
      {/* área de contenido: los overlays viven ACÁ (absolute) — RN-3: jamás tapan los rails */}
      <div className="relative flex min-w-0 flex-1">
        <SessionView />
        <ChangesPanel />
        <BacklogOverlay />
        <ConfigOverlay />
        <ProductoOverlay />
        <RolesOverlay />
        <CloseSessionModal />
        <NewWorkspaceWizard />
      </div>
    </div>
  );
}
