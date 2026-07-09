// Package ports declara los contratos (DIP) entre el dominio/casos de uso y los adaptadores concretos.
package ports

import "context"

// AgentEventKind clasifica un evento normalizado emitido por un AgentSession.
type AgentEventKind string

const (
	EventInit   AgentEventKind = "init"
	EventDelta  AgentEventKind = "delta"
	EventResult AgentEventKind = "result"
	EventError  AgentEventKind = "error"

	// Eventos de herramienta — normalizados y provider-agnósticos (análogos a `tool_call` /
	// `tool_call_update` de ACP, decisión B de la épica conversación-terminal-auténtica). El
	// adaptador traduce sus frames crudos a estos; ningún otro paquete conoce `tool_use`/`tool_result`.
	EventToolCall   AgentEventKind = "tool.call"   // el agente invoca una herramienta (con su input)
	EventToolResult AgentEventKind = "tool.result" // el resultado de una herramienta ya invocada
)

// AgentEvent es el evento normalizado que cruza el puerto — ningún adaptador expone su JSON crudo.
type AgentEvent struct {
	Kind            AgentEventKind
	Text            string // delta/result: texto · tool.result: el output de la herramienta
	ClaudeSessionID string
	Model           string
	Err             string

	// Campos de herramienta (tool.call / tool.result). ToolID parea la llamada con su resultado.
	ToolID      string // id de la invocación (tool_use id) — parea call↔result
	ToolName    string // tool.call: nombre de la herramienta ("Bash", "Read", "Edit", …)
	ToolInput   string // tool.call: input de la herramienta como JSON (para render de la card/diff)
	ToolIsError bool   // tool.result: la herramienta terminó en error
}

// SpawnOpts parametriza el arranque de un agente para una sesión.
type SpawnOpts struct {
	Resume   string // ClaudeSessionID a resumir, vacío = sesión nueva
	Cwd      string // working dir aislado de esta sesión
	ReadOnly bool   // sesión de exploración (PB-27): el agente no puede editar archivos

	// Inyección del arnés del rol (PB-25, patrón HS-11 de la fábrica): la forma-plugin
	// intacta se carga con flags nativos de la CLI — jamás se escribe en el árbol del
	// proyecto (METODOLOGIA §9 ArnesIA; DH-10 intacto: archivos + flags, cero API).
	// La banda Base (CLAUDE.md del arnés) viaja EMBEBIDA en SystemPrompt: la CLI exige
	// UNO de --append-system-prompt / --append-system-prompt-file, no ambos (gate DH-18).
	PluginDirs   []string // --plugin-dir por cada uno (el arnés del rol de ESTA sesión)
	SystemPrompt string   // --append-system-prompt (preámbulo del rol + banda Base)
}

// TranscriptItem es un ítem del transcript reconstruido de una sesión (R1.5, ruta A): el
// historial ordenado (texto + tool-cards en su lugar) que sobrevive al salir/reentrar. Cada
// adapter lo reconstruye de SU propia fuente — Claude de su JSONL; el usecase solo ve esto.
type TranscriptItem struct {
	Kind        string `json:"kind"`                     // "user" | "assistant" | "tool.call" | "tool.result"
	Text        string `json:"text,omitempty"`           // user/assistant: el texto · tool.result: el output
	ToolID      string `json:"tool_id,omitempty"`        // tool.call/tool.result: parea la llamada con su resultado
	ToolName    string `json:"tool_name,omitempty"`      // tool.call: nombre de la herramienta
	ToolInput   string `json:"tool_input,omitempty"`     // tool.call: input como JSON (para la card/diff)
	ToolIsError bool   `json:"tool_is_error,omitempty"`  // tool.result
}

// AgentSession es una conversación viva con un agente (un proceso, un canal de eventos).
type AgentSession interface {
	Send(ctx context.Context, text string) error
	Events() <-chan AgentEvent
	Close() error
}

// AgentPort spawnea agentes. Hoy el único adaptador es Claude Code CLI (driver CLI-nativo, DH-10);
// mañana podría entrar otro agente sin tocar el caso de uso.
type AgentPort interface {
	Spawn(ctx context.Context, opts SpawnOpts) (AgentSession, error)

	// History reconstruye el transcript ordenado de una sesión ya existente (R1.5, ruta A) SIN
	// proceso vivo — el adapter lee su propia fuente (Claude: el JSONL de ~/.claude/projects).
	// providerSessionID vacío o sin fuente → (nil, nil): la sesión aún no tiene historial rico.
	History(ctx context.Context, providerSessionID, cwd string) ([]TranscriptItem, error)
}
