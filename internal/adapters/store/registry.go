// Package store persiste el estado de la app en disco (JSON plano, escritura atómica).
// Desde el shell PRENTER (DH-15) el archivo es state.json {repos, sessions}; un
// sessions.json de F1 se migra solo en el primer Load (sin pérdida — spec shell §6).
package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

type stateFile struct {
	Repos    []domain.Repo    `json:"repos"`
	Sessions []domain.Session `json:"sessions"`
}

// State implementa ports.SessionStore + ports.RepoStore sobre UN archivo state.json.
type State struct {
	mu         sync.Mutex
	path       string
	legacyPath string
	loaded     bool
	data       stateFile
}

// NewState crea el store en path (p.ej. ~/.dev-studio/state.json). legacyPath es el
// sessions.json de F1 a migrar si state.json todavía no existe (vacío = sin migración).
func NewState(path, legacyPath string) *State {
	return &State{path: path, legacyPath: legacyPath}
}

// NewRegistry — compat con el nombre F1 (sessions.json directo, sin repos). Deprecated:
// el binario usa NewState; queda para no romper llamadas viejas en tests/tools.
func NewRegistry(path string) *State { return &State{path: path} }

func (s *State) ensureLoadedLocked() error {
	if s.loaded {
		return nil
	}
	data, err := os.ReadFile(s.path)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, &s.data); err != nil {
			return err
		}
	case errors.Is(err, os.ErrNotExist):
		// primer run con state.json: ¿hay un sessions.json legacy de F1 para migrar?
		if s.legacyPath != "" {
			if legacy, lerr := os.ReadFile(s.legacyPath); lerr == nil {
				var sessions []domain.Session
				if jerr := json.Unmarshal(legacy, &sessions); jerr == nil {
					s.data.Sessions = sessions
				}
			}
		}
	default:
		return err
	}
	if s.data.Repos == nil {
		s.data.Repos = []domain.Repo{}
	}
	if s.data.Sessions == nil {
		s.data.Sessions = []domain.Session{}
	}
	s.loaded = true
	return nil
}

func (s *State) persistLocked() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// --- ports.SessionStore ---

func (s *State) Load(_ context.Context) ([]domain.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureLoadedLocked(); err != nil {
		return nil, err
	}
	out := make([]domain.Session, len(s.data.Sessions))
	copy(out, s.data.Sessions)
	return out, nil
}

func (s *State) Save(_ context.Context, sessions []domain.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureLoadedLocked(); err != nil {
		return err
	}
	if sessions == nil {
		sessions = []domain.Session{}
	}
	s.data.Sessions = sessions
	return s.persistLocked()
}

// --- ports.RepoStore ---

func (s *State) LoadRepos(_ context.Context) ([]domain.Repo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureLoadedLocked(); err != nil {
		return nil, err
	}
	out := make([]domain.Repo, len(s.data.Repos))
	copy(out, s.data.Repos)
	return out, nil
}

func (s *State) SaveRepos(_ context.Context, repos []domain.Repo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureLoadedLocked(); err != nil {
		return err
	}
	if repos == nil {
		repos = []domain.Repo{}
	}
	s.data.Repos = repos
	return s.persistLocked()
}
