package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// ErrNoEsRepo: la ruta no existe o no contiene un repositorio git.
var ErrNoEsRepo = errors.New("la ruta no es un repositorio git (falta .git)")

// ErrRutaProtegida: ruta vedada como repo/cwd (boundary sesion-aislada-por-cwd, RN-2).
var ErrRutaProtegida = errors.New("ruta protegida — la app no opera ahí")

// ValidarRuta aplica la regla de rutas protegidas del dominio (spec workspace-aislado §2.3).
func ValidarRuta(home, ruta string) error {
	if domain.RutaProtegida(home, ruta) {
		return fmt.Errorf("%w: %s", ErrRutaProtegida, ruta)
	}
	return nil
}

// RepoService registra los repositorios de la app (spec shell §5).
// v1 solo acepta rutas locales; clonar desde URL queda explícitamente fuera (spec §9).
type RepoService struct {
	mu    sync.Mutex
	repos []domain.Repo
	store ports.RepoStore
	ctx   context.Context
	home  string
}

func NewRepoService(ctx context.Context, store ports.RepoStore, home string) (*RepoService, error) {
	repos, err := store.LoadRepos(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo service: load: %w", err)
	}
	return &RepoService{repos: repos, store: store, ctx: ctx, home: home}, nil
}

// Register valida que la ruta sea un repo git local NO protegido y lo registra. Misma ruta
// ya registrada → devuelve el existente (idempotente, no duplica).
func (s *RepoService) Register(ruta string) (domain.Repo, error) {
	abs, err := filepath.Abs(filepath.Clean(ruta))
	if err != nil {
		return domain.Repo{}, err
	}
	if err := ValidarRuta(s.home, abs); err != nil {
		return domain.Repo{}, err
	}
	if fi, err := os.Stat(filepath.Join(abs, ".git")); err != nil || !fi.IsDir() {
		return domain.Repo{}, fmt.Errorf("%w: %s", ErrNoEsRepo, abs)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.repos {
		if r.Ruta == abs {
			return r, nil
		}
	}
	repo := domain.Repo{ID: newID(), Nombre: filepath.Base(abs), Ruta: abs}
	s.repos = append(s.repos, repo)
	return repo, s.store.SaveRepos(s.ctx, s.repos)
}

func (s *RepoService) List() []domain.Repo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.Repo, len(s.repos))
	copy(out, s.repos)
	return out
}

func (s *RepoService) Get(id string) (domain.Repo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.repos {
		if r.ID == id {
			return r, true
		}
	}
	return domain.Repo{}, false
}

func (s *RepoService) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.repos {
		if r.ID == id {
			s.repos = append(s.repos[:i], s.repos[i+1:]...)
			return s.store.SaveRepos(s.ctx, s.repos)
		}
	}
	return ErrNotFound
}
