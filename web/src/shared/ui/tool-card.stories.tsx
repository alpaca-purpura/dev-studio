import type { Meta, StoryObj } from "@storybook/react-vite";
import { ToolCard } from "./tool-card";

const meta: Meta<typeof ToolCard> = { title: "Átomos PRENTER/ToolCard", component: ToolCard };
export default meta;

type S = StoryObj<typeof ToolCard>;

/** Read con output — la card colapsa el cuerpo largo, header `⎿ Read …/path · ✓`. */
export const Lectura: S = {
  args: {
    call: {
      tool_id: "toolu_1",
      name: "Read",
      input: JSON.stringify({ file_path: "internal/adapters/agent/claudecode/conductor.go" }),
      status: "ok",
      output: "package claudecode\n\nimport (\n  \"bufio\"\n  \"context\"\n)\n// … 200 líneas más",
    },
  },
};

/** Bash — arg = el comando, output = stdout. */
export const Comando: S = {
  args: {
    call: {
      tool_id: "toolu_2",
      name: "Bash",
      input: JSON.stringify({ command: "go test ./internal/...", description: "correr tests" }),
      status: "ok",
      output: "ok  github.com/alpacapurpura/dev-studio/internal/usecase  0.009s",
    },
  },
};

/** Edit — cuerpo = diff (+ teal / − rojo) + resumen +N −M. */
export const Edicion: S = {
  args: {
    call: {
      tool_id: "toolu_3",
      name: "Edit",
      input: JSON.stringify({
        file_path: "internal/ports/agent.go",
        old_string: "EventError  AgentEventKind = \"error\"",
        new_string: "EventError      AgentEventKind = \"error\"\n\tEventToolCall   AgentEventKind = \"tool.call\"\n\tEventToolResult AgentEventKind = \"tool.result\"",
      }),
      status: "ok",
    },
  },
};

/** Corriendo — dot teal pulsante mientras la tool no terminó. */
export const Corriendo: S = {
  args: {
    call: { tool_id: "toolu_4", name: "Bash", input: JSON.stringify({ command: "npm run build" }), status: "running" },
  },
};

/** Error — ✕ rojo + output en rojo. */
export const ConError: S = {
  args: {
    call: {
      tool_id: "toolu_5",
      name: "Read",
      input: JSON.stringify({ file_path: "no/existe.go" }),
      status: "error",
      output: "Error: file not found: no/existe.go",
    },
  },
};
