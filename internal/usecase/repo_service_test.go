package usecase_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/usecase"
)

type memRepoStore struct{ repos []domain.Repo }

func (m *memRepoStore) LoadRepos(context.Context) ([]domain.Repo, error) { return m.repos, nil }
func (m *memRepoStore) SaveRepos(_ context.Context, r []domain.Repo) error {
	m.repos = r
	return nil
}

func gitDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestRegisterValidaRepoGit(t *testing.T) {
	svc, err := usecase.NewRepoService(context.Background(), &memRepoStore{})
	if err != nil {
		t.Fatal(err)
	}

	// ruta sin .git → rechazada
	if _, err := svc.Register(t.TempDir()); err == nil {
		t.Fatal("ruta sin .git debía rechazarse")
	}
	// ruta inexistente → rechazada
	if _, err := svc.Register("/no/existe/xyz"); err == nil {
		t.Fatal("ruta inexistente debía rechazarse")
	}

	dir := gitDir(t)
	repo, err := svc.Register(dir)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if repo.Nombre != filepath.Base(dir) || repo.Ruta != dir || repo.ID == "" {
		t.Errorf("repo mal formado: %+v", repo)
	}

	// duplicado (misma ruta) → devuelve el existente, no crea otro
	again, err := svc.Register(dir)
	if err != nil {
		t.Fatalf("register duplicado: %v", err)
	}
	if again.ID != repo.ID || len(svc.List()) != 1 {
		t.Errorf("duplicado debía devolver el existente: %+v vs %+v", again, repo)
	}
}

func TestRemoveYGet(t *testing.T) {
	store := &memRepoStore{}
	svc, _ := usecase.NewRepoService(context.Background(), store)
	dir := gitDir(t)
	repo, _ := svc.Register(dir)

	got, ok := svc.Get(repo.ID)
	if !ok || got.Ruta != dir {
		t.Fatalf("get: %+v ok=%v", got, ok)
	}
	if err := svc.Remove(repo.ID); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, ok := svc.Get(repo.ID); ok {
		t.Fatal("repo debía desaparecer")
	}
	if len(store.repos) != 0 {
		t.Fatal("remove debía persistirse")
	}
}
