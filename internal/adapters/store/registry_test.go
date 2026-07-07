package store_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/adapters/store"
	"github.com/alpacapurpura/dev-studio/internal/domain"
)

func TestMigraSessionsJSONLegacy(t *testing.T) {
	dir := t.TempDir()
	legacy := filepath.Join(dir, "sessions.json")
	statePath := filepath.Join(dir, "state.json")

	legacyJSON := `[{"id":"s1","nombre":"vieja","cwd":"/tmp","status":"idle","conv":[]}]`
	if err := os.WriteFile(legacy, []byte(legacyJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	st := store.NewState(statePath, legacy)
	sessions, err := st.Load(context.Background())
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(sessions) != 1 || sessions[0].ID != "s1" || sessions[0].Nombre != "vieja" {
		t.Fatalf("migración: esperaba la sesión legacy, obtuve %+v", sessions)
	}
	repos, err := st.LoadRepos(context.Background())
	if err != nil || len(repos) != 0 {
		t.Fatalf("repos tras migración: %v %+v", err, repos)
	}
}

func TestPersisteReposYSesionesJuntos(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")
	st := store.NewState(statePath, filepath.Join(dir, "no-legacy.json"))
	ctx := context.Background()

	if err := st.SaveRepos(ctx, []domain.Repo{{ID: "r1", Nombre: "demo", Ruta: "/x"}}); err != nil {
		t.Fatal(err)
	}
	if err := st.Save(ctx, []domain.Session{{ID: "s9", Nombre: "n", Cwd: "/x", RepoID: "r1", Conv: []domain.Turn{}}}); err != nil {
		t.Fatal(err)
	}

	// reabrir desde disco: ambas mitades sobreviven
	st2 := store.NewState(statePath, "")
	repos, err := st2.LoadRepos(ctx)
	if err != nil || len(repos) != 1 || repos[0].ID != "r1" {
		t.Fatalf("repos: %v %+v", err, repos)
	}
	sessions, err := st2.Load(ctx)
	if err != nil || len(sessions) != 1 || sessions[0].RepoID != "r1" {
		t.Fatalf("sessions: %v %+v", err, sessions)
	}
}

func TestPrimerRunSinArchivos(t *testing.T) {
	dir := t.TempDir()
	st := store.NewState(filepath.Join(dir, "state.json"), filepath.Join(dir, "sessions.json"))
	sessions, err := st.Load(context.Background())
	if err != nil || len(sessions) != 0 {
		t.Fatalf("primer run: %v %+v", err, sessions)
	}
}
