// Package domain contiene el modelo de negocio de DevStudio, sin dependencias de transporte ni de shell.
package domain

// Status es el estado de una sesión de trabajo.
type Status string

const (
	StatusIdle      Status = "idle"
	StatusStreaming Status = "streaming"
)

// Turn es un turno de la conversación (para replay instantáneo en UI).
// La fuente de verdad de la conversación completa vive en el JSONL propio de Claude Code.
type Turn struct {
	Role string `json:"role"` // user | assistant
	Text string `json:"text"`
}

// Session es un frente de trabajo: una conversación viva de Claude Code, N:1 con un directorio.
type Session struct {
	ID              string `json:"id"`
	Nombre          string `json:"nombre"`
	Cwd             string `json:"cwd"`
	Status          Status `json:"status"`
	ClaudeSessionID string `json:"claude_session_id,omitempty"`
	Model           string `json:"model,omitempty"`
	Conv            []Turn `json:"conv"`
}
