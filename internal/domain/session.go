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

// Session es un frente de trabajo: una conversación viva de Claude Code, 1:1 con un
// workspace (spec shell RN-2). RepoID/Historia/Rol entran con el shell PRENTER (DH-15);
// una sesión previa al shell puede no tener RepoID (se agrupa bajo «(sin repositorio)»).
type Session struct {
	ID              string    `json:"id"`
	Nombre          string    `json:"nombre"`
	Cwd             string    `json:"cwd"`
	Status          Status    `json:"status"`
	ClaudeSessionID string    `json:"claude_session_id,omitempty"`
	Model           string    `json:"model,omitempty"`
	RepoID          string    `json:"repo_id,omitempty"`
	Historia        *Historia `json:"historia,omitempty"`
	Rol             string    `json:"rol,omitempty"`
	Workspace       string    `json:"workspace,omitempty"` // ruta del worktree propio (PB-02); vacío = sin aislamiento
	Branch          string    `json:"branch,omitempty"`    // branch del workspace ({prefijo-tipo}/{slug})
	Modo            string    `json:"modo,omitempty"`      // trabajo (default) | exploracion (read-only, PB-27)
	Conv            []Turn    `json:"conv"`
}
