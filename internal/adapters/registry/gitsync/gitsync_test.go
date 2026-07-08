package gitsync

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// mkMarketplaceRepo crea un repo git local con forma de marketplace mínimo y lo devuelve.
func mkMarketplaceRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".claude-plugin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".claude-plugin", "marketplace.json"), []byte(`{"name":"mp","plugins":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "t@t"},
		{"config", "user.name", "t"},
		{"add", ".claude-plugin/marketplace.json"},
		{"commit", "-m", "init"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return dir
}

func TestSyncRutaLocalDirecta(t *testing.T) {
	// una ruta local con marketplace se usa TAL CUAL — sin git, sin copia
	src := mkMarketplaceRepo(t)
	s := New("git", t.TempDir())
	dir, err := s.Sync(context.Background(), src)
	if err != nil {
		t.Fatalf("Sync ruta local: %v", err)
	}
	if dir != src {
		t.Fatalf("ruta local debía usarse directa: %s != %s", dir, src)
	}
}

func TestSyncClonaURLDentroDelBase(t *testing.T) {
	src := mkMarketplaceRepo(t)
	base := t.TempDir()
	s := New("git", base)

	// file:// fuerza el camino "remoto" (clone), no el de ruta local
	url := "file://" + src
	dir, err := s.Sync(context.Background(), url)
	if err != nil {
		t.Fatalf("Sync clone: %v", err)
	}
	if !strings.HasPrefix(dir, base+string(os.PathSeparator)) {
		t.Fatalf("clone FUERA del base confinado (RN-3): %s no está bajo %s", dir, base)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude-plugin", "marketplace.json")); err != nil {
		t.Fatalf("clone sin contenido del marketplace: %v", err)
	}

	// segundo Sync = actualización (pull ff-only) del MISMO destino, no un clon nuevo
	dir2, err := s.Sync(context.Background(), url)
	if err != nil {
		t.Fatalf("Sync segunda vez (pull): %v", err)
	}
	if dir2 != dir {
		t.Fatalf("re-sync cambió el destino: %s != %s", dir2, dir)
	}
}

func TestSyncRechazaSourceInvalido(t *testing.T) {
	s := New("git", t.TempDir())
	if _, err := s.Sync(context.Background(), filepath.Join(t.TempDir(), "no-existe")); err == nil {
		t.Fatal("ruta local inexistente debía fallar honesto")
	}
	if _, err := s.Sync(context.Background(), t.TempDir()); err == nil {
		t.Fatal("dir local SIN marketplace.json debía fallar (no es un marketplace)")
	}
}
