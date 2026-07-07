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
