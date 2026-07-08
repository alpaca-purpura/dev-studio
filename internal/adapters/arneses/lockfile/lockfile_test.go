package lockfile

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

func TestRoundTrip(t *testing.T) {
	repo := t.TempDir()
	l := New()
	ctx := context.Background()

	// repo sin lock = roster vacío honesto, sin error
	reg, inst, err := l.Leer(ctx, repo)
	if err != nil || reg != "" || len(inst) != 0 {
		t.Fatalf("repo sin lock: reg=%q inst=%v err=%v", reg, inst, err)
	}

	quiere := []domain.ArnesInstalado{
		{ID: "dev-full-cycle", Version: "0.1.0", Canal: "beta"},
		{ID: "demo-auditor", Version: "0.2.0"},
	}
	lockPath, err := l.Escribir(ctx, repo, "git@github.com:org/mp.git", quiere)
	if err != nil {
		t.Fatalf("Escribir: %v", err)
	}
	if lockPath != filepath.Join(repo, ".devstudio", "arneses.yaml") {
		t.Fatalf("lock path: %s", lockPath)
	}

	reg, inst, err = l.Leer(ctx, repo)
	if err != nil {
		t.Fatalf("Leer: %v", err)
	}
	if reg != "git@github.com:org/mp.git" || len(inst) != 2 || inst[0].ID != "dev-full-cycle" || inst[0].Canal != "beta" {
		t.Fatalf("round-trip roto: reg=%q inst=%+v", reg, inst)
	}

	// legible por humanos (as-code): YAML con claves esperadas
	raw, _ := os.ReadFile(lockPath)
	for _, k := range []string{"registry:", "arneses:", "id: dev-full-cycle", "version: 0.1.0"} {
		if !strings.Contains(string(raw), k) {
			t.Errorf("lock sin %q:\n%s", k, raw)
		}
	}
}
