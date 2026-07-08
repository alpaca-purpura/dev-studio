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
)

// AgentEvent es el evento normalizado que cruza el puerto — ningún adaptador expone su JSON crudo.
type AgentEvent struct {
	Kind            AgentEventKind
	Text            string
	ClaudeSessionID string
	Model           string
	Err             string
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
}
