import type { Meta, StoryObj } from "@storybook/react-vite";
import { useEffect, useState } from "react";
import { StudioNav } from "./studio-nav";
import { ReposRail } from "./repos-rail";
import { SessionView } from "./session-view";
import { useSessions } from "../shared/store/sessions-store";
import { useRepos } from "../shared/store/repos-store";
import { Button, ModalShell, PanelOverlay } from "../shared/ui";

const meta: Meta = { title: "Organismos del shell" };
export default meta;

/** Seed de stores para render standalone (sin backend): el estado es demo, el componente es el real. */
function seed() {
  useRepos.setState({
    repos: [{ id: "r1", nombre: "dev-studio", ruta: "/home/demo/dev-studio" }],
    sessionStatus: {
      s1: { branch: "wt/aislar-sesion", files: [{ path: "internal/adapters/git/cli/cli.go", state: "M" }], add: 41, del: 9 },
      s2: { branch: "wt/fix-sse", files: [], add: 0, del: 0 },
    },
  });
  useSessions.setState({
    sessions: [
      {
        id: "s1",
        nombre: "aislar-sesion-worktree",
        cwd: "/home/demo/dev-studio",
        status: "idle",
        repo_id: "r1",
        rol: "Full-Stack",
        historia: { id: "hx-worktree", titulo: "Aislar sesión en workspace propio (PB-02)" },
        model: "opus",
        conv: [
          { role: "user", text: "PB-02: cerrá sesion-aislada-por-cwd creando un workspace propio por sesión." },
          { role: "assistant", text: "Entendido. Escribo el test de aislamiento primero (TDD) y después el manager." },
        ],
      },
      { id: "s2", nombre: "fix-sse-reconexion", cwd: "/home/demo/dev-studio", status: "streaming", repo_id: "r1", rol: "Backend", conv: [] },
    ],
    activeId: "s1",
  });
}

export const NavDelEstudio: StoryObj = {
  render: () => (
    <div className="h-96 overflow-hidden rounded-lg border border-border">
      <StudioNav />
    </div>
  ),
};

export const RailDeRepositorios: StoryObj = {
  render: function Render() {
    useEffect(seed, []);
    return (
      <div className="h-[480px] overflow-hidden rounded-lg border border-border">
        <ReposRail />
      </div>
    );
  },
};

export const VistaDeSesion: StoryObj = {
  render: function Render() {
    useEffect(seed, []);
    return (
      <div className="flex h-[520px] overflow-hidden rounded-lg border border-border">
        <SessionView />
      </div>
    );
  },
};

export const OverlayDePanel: StoryObj = {
  render: function Render() {
    const [open, setOpen] = useState(true);
    return (
      <div className="relative h-96 overflow-hidden rounded-lg border border-border bg-card p-4">
        <Button onClick={() => setOpen(true)}>Abrir overlay</Button>
        <PanelOverlay
          open={open}
          onClose={() => setOpen(false)}
          title="Backlog"
          subtitle="El overlay tapa SOLO el área de contenido — los rails quedan visibles (RN-3)."
        >
          <p className="p-4 text-sm text-muted-foreground">Contenido del panel.</p>
        </PanelOverlay>
      </div>
    );
  },
};

export const Modal: StoryObj = {
  render: function Render() {
    const [open, setOpen] = useState(false);
    return (
      <>
        <Button onClick={() => setOpen(true)}>Abrir modal</Button>
        <ModalShell
          open={open}
          onClose={() => setOpen(false)}
          title="Conversación por voz"
          subtitle="Todavía no incluida — requiere backend propio."
          icon="◔"
          footer={<Button onClick={() => setOpen(false)}>Entendido</Button>}
        >
          <p className="text-sm text-muted-foreground">
            El botón queda visible a propósito: mensaje claro de «todavía no», no una feature a
            medias (RN-7).
          </p>
        </ModalShell>
      </>
    );
  },
};
