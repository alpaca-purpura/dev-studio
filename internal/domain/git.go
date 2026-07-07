package domain

// GitFile es un archivo tocado en el working tree (status porcelain simplificado).
type GitFile struct {
	Path  string `json:"path"`
	State string `json:"state"` // M modificado · A/? nuevo · D eliminado · R renombrado
}

// GitStatus es la foto del working tree de un cwd.
type GitStatus struct {
	Branch string    `json:"branch"`
	Files  []GitFile `json:"files"`
	Add    int       `json:"add"` // líneas agregadas (diff vs HEAD)
	Del    int       `json:"del"` // líneas eliminadas
}

// GitDiff es el diff de UN archivo, en las dos formas que la UI muestra
// (quick-look unificado y revisión side-by-side).
type GitDiff struct {
	Path     string `json:"path"`
	Unified  string `json:"unified"`
	Original string `json:"original"` // vacío = archivo nuevo
	Modified string `json:"modified"`
	Binary   bool   `json:"binary"`
}

// GitLogEntry es una entrada del historial.
type GitLogEntry struct {
	SHA     string `json:"sha"`
	Mensaje string `json:"mensaje"`
	Autor   string `json:"autor"`
	Fecha   string `json:"fecha"` // ISO-8601
}
