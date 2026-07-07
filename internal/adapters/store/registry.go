// Package store persiste el registro de sesiones en disco (JSON plano, escritura atómica).
package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

// Registry implementa ports.SessionStore sobre un archivo JSON.
type Registry struct {
	path string
}

// NewRegistry crea el store en path (p.ej. ~/.dev-studio/sessions.json).
func NewRegistry(path string) *Registry {
	return &Registry{path: path}
}

func (r *Registry) Load(_ context.Context) ([]domain.Session, error) {
	data, err := os.ReadFile(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil // primer run: registro vacío no es error
	}
	if err != nil {
		return nil, err
	}
	var sessions []domain.Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *Registry) Save(_ context.Context, sessions []domain.Session) error {
	if sessions == nil {
		sessions = []domain.Session{}
	}
	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, r.path)
}
