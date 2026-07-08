import type { Meta, StoryObj } from "@storybook/react-vite";
import { Button } from "../../../shared/ui";
import { ArnesCard, ArnesDetalle } from "./arnes-card";

const dfc = {
  id: "dev-full-cycle",
  nombre: "dev-full-cycle",
  descripcion: "Arnés dogfood: desarrollo de software end-to-end (idea → released).",
  rol: "Ingeniería · Desarrollo full-cycle",
  proceso: "desarrollo de software end-to-end (idea → released)",
  version: "0.1.0",
  canal: "beta",
  fases: ["spec", "build", "review", "release"],
};

const meta: Meta = {
  title: "Overlays/ArnesCard",
  parameters: { layout: "padded" },
};
export default meta;

export const Tarjeta: StoryObj = {
  render: () => (
    <div className="w-64 space-y-1">
      <ArnesCard arnes={dfc} selected onSelect={() => {}} />
      <ArnesCard
        arnes={{ id: "demo-auditor", nombre: "demo-auditor", rol: "Auditoría", version: "0.2.0" }}
        onSelect={() => {}}
        accion={
          <Button variant="ghost" className="px-2 py-1 text-xs">
            Instalar
          </Button>
        }
      />
    </div>
  ),
};

export const Detalle: StoryObj = {
  render: () => (
    <div className="max-w-xl">
      <ArnesDetalle arnes={dfc} />
    </div>
  ),
};
