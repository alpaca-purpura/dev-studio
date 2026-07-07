package cli_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/adapters/git/cli"
)

func TestCreateWorktreeAislado(t *testing.T) {
	repo := initRepo(t)
	g := cli.New()
	ctx := context.Background()

	// base: el repo necesita HEAD
	write(t, repo, "base.txt", "base\n")
	if _, err := g.Commit(ctx, repo, []string{"base.txt"}, "base"); err != nil {
		t.Fatal(err)
	}

	destino := filepath.Join(t.TempDir(), "workspaces", "demo", "mi-historia")
	path, branch, err := g.CreateWorktree(ctx, repo, destino, "mi-historia")
	if err != nil {
		t.Fatalf("create worktree: %v", err)
	}
	if path != destino {
		t.Errorf("path: esperaba %s, obtuve %s", destino, path)
	}
	if branch != "wt/mi-historia" {
		t.Errorf("branch: esperaba wt/mi-historia, obtuve %s", branch)
	}
	// el worktree es un checkout real, aislado del repo raíz
	if _, err := os.Stat(filepath.Join(path, "base.txt")); err != nil {
		t.Fatalf("el worktree no tiene el contenido base: %v", err)
	}
	st, err := g.Status(ctx, path)
	if err != nil {
		t.Fatalf("status en worktree: %v", err)
	}
	if st.Branch != "wt/mi-historia" {
		t.Errorf("branch del worktree: %s", st.Branch)
	}
	// tocar el worktree NO ensucia el repo raíz
	write(t, path, "solo-wt.txt", "x\n")
	stRoot, _ := g.Status(ctx, repo)
	for _, f := range stRoot.Files {
		if f.Path == "solo-wt.txt" {
			t.Error("el archivo del worktree apareció en el status de la raíz")
		}
	}
	// listado git lo conoce
	out := gitOut(t, repo, "worktree", "list")
	if !strings.Contains(out, destino) {
		t.Errorf("worktree list no muestra %s:\n%s", destino, out)
	}
}

func TestCreateWorktreeColisionDeBranch(t *testing.T) {
	repo := initRepo(t)
	g := cli.New()
	ctx := context.Background()
	write(t, repo, "a.txt", "1\n")
	if _, err := g.Commit(ctx, repo, []string{"a.txt"}, "base"); err != nil {
		t.Fatal(err)
	}

	base := t.TempDir()
	_, b1, err := g.CreateWorktree(ctx, repo, filepath.Join(base, "uno"), "misma")
	if err != nil {
		t.Fatal(err)
	}
	_, b2, err := g.CreateWorktree(ctx, repo, filepath.Join(base, "dos"), "misma")
	if err != nil {
		t.Fatalf("segunda con mismo slug debía resolverse con sufijo: %v", err)
	}
	if b1 == b2 {
		t.Errorf("branches iguales: %s", b1)
	}
	if b2 != "wt/misma-2" {
		t.Errorf("esperaba wt/misma-2, obtuve %s", b2)
	}
}

func TestCreateWorktreeSinCommitsFallaHonesto(t *testing.T) {
	repo := initRepo(t) // sin commits
	g := cli.New()
	_, _, err := g.CreateWorktree(context.Background(), repo, filepath.Join(t.TempDir(), "w"), "x")
	if err == nil {
		t.Fatal("repo sin commits debía fallar con mensaje claro")
	}
	if !strings.Contains(err.Error(), "commit") {
		t.Errorf("mensaje poco claro: %v", err)
	}
}

func TestRemoveWorktree(t *testing.T) {
	repo := initRepo(t)
	g := cli.New()
	ctx := context.Background()
	write(t, repo, "a.txt", "1\n")
	if _, err := g.Commit(ctx, repo, []string{"a.txt"}, "base"); err != nil {
		t.Fatal(err)
	}
	path, _, err := g.CreateWorktree(ctx, repo, filepath.Join(t.TempDir(), "w"), "efimero")
	if err != nil {
		t.Fatal(err)
	}

	// sucio → remove rechazado (RN-3: jamás --force)
	write(t, path, "pendiente.txt", "sin commit\n")
	if err := g.RemoveWorktree(ctx, repo, path); err == nil {
		t.Fatal("remove con cambios sin commit debía rechazarse")
	}

	// limpio → remove ok
	if err := os.Remove(filepath.Join(path, "pendiente.txt")); err != nil {
		t.Fatal(err)
	}
	if err := g.RemoveWorktree(ctx, repo, path); err != nil {
		t.Fatalf("remove limpio: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("el worktree debía desaparecer del disco")
	}
}
