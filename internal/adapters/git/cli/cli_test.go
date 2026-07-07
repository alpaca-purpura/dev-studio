package cli_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/adapters/git/cli"
)

// initRepo crea un repo git real en un tmpdir con identidad de test.
func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init", "-b", "main"},
		{"config", "user.email", "test@dev-studio.local"},
		{"config", "user.name", "Test"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return dir
}

func write(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestStatusListaUntrackedYModificados(t *testing.T) {
	dir := initRepo(t)
	g := cli.New()
	ctx := context.Background()

	write(t, dir, "a.txt", "uno\n")
	st, err := g.Status(ctx, dir)
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if st.Branch != "main" {
		t.Errorf("branch: esperaba main, obtuve %q", st.Branch)
	}
	if len(st.Files) != 1 || st.Files[0].Path != "a.txt" {
		t.Fatalf("esperaba a.txt untracked, obtuve %+v", st.Files)
	}

	// commit inicial y modificar: el status debe reflejar M
	if _, err := g.Commit(ctx, dir, []string{"a.txt"}, "inicial"); err != nil {
		t.Fatalf("commit: %v", err)
	}
	write(t, dir, "a.txt", "uno\ndos\n")
	st, err = g.Status(ctx, dir)
	if err != nil {
		t.Fatalf("status 2: %v", err)
	}
	if len(st.Files) != 1 || st.Files[0].State != "M" {
		t.Fatalf("esperaba a.txt M, obtuve %+v", st.Files)
	}
	if st.Add < 1 {
		t.Errorf("shortstat: esperaba add>=1, obtuve %d", st.Add)
	}
}

func TestCommitSoloPathsSeleccionados(t *testing.T) {
	dir := initRepo(t)
	g := cli.New()
	ctx := context.Background()

	write(t, dir, "elegido.txt", "sí\n")
	write(t, dir, "excluido.txt", "no\n")

	sha, err := g.Commit(ctx, dir, []string{"elegido.txt"}, "solo elegido")
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	if len(sha) < 7 {
		t.Errorf("sha corto: %q", sha)
	}
	// el commit contiene SOLO elegido.txt
	files := gitOut(t, dir, "show", "--name-only", "--pretty=format:", "HEAD")
	if !strings.Contains(files, "elegido.txt") || strings.Contains(files, "excluido.txt") {
		t.Fatalf("commit debía contener solo elegido.txt, contiene: %q", files)
	}
	// excluido sigue untracked
	st, _ := g.Status(ctx, dir)
	if len(st.Files) != 1 || st.Files[0].Path != "excluido.txt" {
		t.Fatalf("excluido.txt debía seguir pendiente, status: %+v", st.Files)
	}
}

func TestCommitSinPathsEsError(t *testing.T) {
	dir := initRepo(t)
	g := cli.New()
	if _, err := g.Commit(context.Background(), dir, nil, "vacío"); err == nil {
		t.Fatal("commit sin paths debía fallar (RN-5)")
	}
}

func TestDiffFileNuevoYModificado(t *testing.T) {
	dir := initRepo(t)
	g := cli.New()
	ctx := context.Background()

	write(t, dir, "n.txt", "línea nueva\n")
	d, err := g.DiffFile(ctx, dir, "n.txt")
	if err != nil {
		t.Fatalf("diff untracked: %v", err)
	}
	if d.Original != "" {
		t.Errorf("archivo nuevo: original debía ser vacío, obtuve %q", d.Original)
	}
	if !strings.Contains(d.Modified, "línea nueva") {
		t.Errorf("modified debía tener el contenido, obtuve %q", d.Modified)
	}
	if !strings.Contains(d.Unified, "+línea nueva") {
		t.Errorf("unified debía marcar +línea nueva, obtuve %q", d.Unified)
	}

	if _, err := g.Commit(ctx, dir, []string{"n.txt"}, "base"); err != nil {
		t.Fatal(err)
	}
	write(t, dir, "n.txt", "línea nueva\notra\n")
	d, err = g.DiffFile(ctx, dir, "n.txt")
	if err != nil {
		t.Fatalf("diff modificado: %v", err)
	}
	if !strings.Contains(d.Original, "línea nueva") || strings.Contains(d.Original, "otra") {
		t.Errorf("original debía ser el contenido de HEAD, obtuve %q", d.Original)
	}
	if !strings.Contains(d.Modified, "otra") {
		t.Errorf("modified debía tener la línea nueva, obtuve %q", d.Modified)
	}
}

func TestLogDevuelveEntradas(t *testing.T) {
	dir := initRepo(t)
	g := cli.New()
	ctx := context.Background()

	write(t, dir, "a.txt", "1\n")
	if _, err := g.Commit(ctx, dir, []string{"a.txt"}, "primero"); err != nil {
		t.Fatal(err)
	}
	write(t, dir, "a.txt", "2\n")
	if _, err := g.Commit(ctx, dir, []string{"a.txt"}, "segundo"); err != nil {
		t.Fatal(err)
	}

	entries, err := g.Log(ctx, dir, 10)
	if err != nil {
		t.Fatalf("log: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("esperaba 2 commits, obtuve %d", len(entries))
	}
	if entries[0].Mensaje != "segundo" || entries[1].Mensaje != "primero" {
		t.Errorf("orden/mensajes: %+v", entries)
	}
	if entries[0].Autor != "Test" || entries[0].SHA == "" || entries[0].Fecha == "" {
		t.Errorf("metadata incompleta: %+v", entries[0])
	}
}

func TestStatusEnRepoVacioNoRevienta(t *testing.T) {
	dir := initRepo(t) // sin commits: HEAD no existe
	g := cli.New()
	st, err := g.Status(context.Background(), dir)
	if err != nil {
		t.Fatalf("status repo vacío: %v", err)
	}
	if st.Branch == "" {
		t.Error("branch vacía en repo recién inicializado")
	}
}
