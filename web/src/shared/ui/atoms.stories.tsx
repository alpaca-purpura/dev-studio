import type { Meta, StoryObj } from "@storybook/react-vite";
import {
  Avatar,
  Button,
  Card,
  Chip,
  DiffStat,
  FieldLabel,
  FilterChip,
  HonestBanner,
  IconButton,
  Input,
  Kbd,
  Pill,
  ReleaseChip,
  Select,
  StatePill,
  StatusDot,
  statusLabel,
  Tabs,
  Textarea,
  Toggle,
  storyStateLabel,
  type StoryState,
  type WorkspaceState,
} from "./index";
import { useState } from "react";

const meta: Meta = { title: "Átomos PRENTER" };
export default meta;

export const Botones: StoryObj = {
  render: () => (
    <div className="flex flex-wrap items-center gap-3">
      <Button>Crear sesión aislada</Button>
      <Button variant="outline">Confirmar cambios</Button>
      <Button variant="ghost">Cancelar</Button>
      <Button disabled>Deshabilitado</Button>
      <IconButton aria-label="Cerrar">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <path d="M6 6l12 12M18 6L6 18" />
        </svg>
      </IconButton>
    </div>
  ),
};

export const ChipsYPills: StoryObj = {
  render: () => (
    <div className="flex flex-wrap items-center gap-3">
      <Chip>wt/aislar-sesion-worktree</Chip>
      <Chip>3 cambios</Chip>
      <Pill tone="ok">🟢 live</Pill>
      <Pill tone="warn">⏳ wip</Pill>
      <Pill tone="crit">drift</Pill>
      <Pill tone="muted">📋 planned</Pill>
      <Pill tone="primary">Rol actual</Pill>
      <ReleaseChip release="F2" />
      <DiffStat add={41} del={9} />
      <Kbd>⌘1</Kbd>
    </div>
  ),
};

export const EstadosDeWorkspace: StoryObj = {
  render: () => (
    <div className="flex flex-col gap-2">
      {(["idle", "ready", "conflict", "archived", "streaming"] as WorkspaceState[]).map((s) => (
        <div key={s} className="flex items-center gap-2 text-sm">
          <StatusDot state={s} />
          <span>{statusLabel[s]}</span>
          <code className="font-mono text-xs text-muted-foreground">{s}</code>
        </div>
      ))}
    </div>
  ),
};

export const EstadosDeHistoria: StoryObj = {
  render: () => (
    <div className="flex max-w-md flex-wrap gap-2">
      {(Object.keys(storyStateLabel) as StoryState[]).map((s) => (
        <StatePill key={s} state={s} />
      ))}
    </div>
  ),
};

export const Formulario: StoryObj = {
  render: () => (
    <Card className="max-w-sm space-y-3 p-4">
      <div>
        <FieldLabel>Ruta del repositorio</FieldLabel>
        <Input placeholder="/home/usuaria/proyectos/mi-repo" />
      </div>
      <div>
        <FieldLabel>Proveedor</FieldLabel>
        <Select>
          <option>Claude Code</option>
        </Select>
      </div>
      <div>
        <FieldLabel>Mensaje de commit</FieldLabel>
        <Textarea rows={2} placeholder="Mensaje de commit…" />
      </div>
    </Card>
  ),
};

export const TogglesYTabs: StoryObj = {
  render: function Render() {
    const [on, setOn] = useState(true);
    const [tab, setTab] = useState("cambios");
    const [chip, setChip] = useState("all");
    return (
      <div className="flex max-w-md flex-col gap-4">
        <div className="flex items-center gap-2 text-sm">
          <Toggle on={on} onChange={setOn} aria-label="Restaurar comandos al iniciar" />
          <span>Restaurar los comandos al iniciar</span>
        </div>
        <Tabs
          items={[
            { id: "archivos", label: "Archivos" },
            { id: "cambios", label: "Cambios" },
            { id: "pruebas", label: "Pruebas" },
          ]}
          active={tab}
          onSelect={setTab}
        />
        <div className="flex gap-2">
          {["all", "fs", "qa"].map((k) => (
            <FilterChip key={k} active={chip === k} onClick={() => setChip(k)}>
              {k === "all" ? "Todos 3" : k === "fs" ? "Full-Stack 2" : "QA 1"}
            </FilterChip>
          ))}
        </div>
      </div>
    );
  },
};

export const AvataresYHonestidad: StoryObj = {
  render: () => (
    <div className="flex max-w-md flex-col gap-4">
      <div className="flex gap-2">
        {["FS", "FE", "BE", "QA", "SG", "AR", "PO", "ID"].map((r) => (
          <Avatar key={r} label={r} />
        ))}
      </div>
      <HonestBanner>
        Los roles vienen del registry de arneses de tu organización — registry propio en
        construcción (PB-25). Este roster es una vista previa estática.
      </HonestBanner>
    </div>
  ),
};
