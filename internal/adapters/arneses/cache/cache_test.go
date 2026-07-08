package cache

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestMaterializarCopiaIntacta(t *testing.T) {
	src := t.TempDir()
	// mini forma-plugin
	for _, f := range []struct{ p, c string }{
		{".claude-plugin/plugin.json", `{"name":"x","version":"1.0.0"}`},
		{"arnes.l0.json", `{"id":"x"}`},
		{"skills/a/SKILL.md", "---\nname: a\n---\n"},
		{"CLAUDE.md", "# base\n"},
	} {
		p := filepath.Join(src, f.p)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(f.c), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	base := t.TempDir()
	c := New(base)
	dir, err := c.Materializar(context.Background(), src, "x", "1.0.0")
	if err != nil {
		t.Fatalf("Materializar: %v", err)
	}
	if dir != filepath.Join(base, "x", "1.0.0") || dir != c.Dir("x", "1.0.0") {
		t.Fatalf("dir del caché: %s", dir)
	}
	for _, p := range []string{".claude-plugin/plugin.json", "arnes.l0.json", "skills/a/SKILL.md", "CLAUDE.md"} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("payload incompleto, falta %s: %v", p, err)
		}
	}

	// re-materializar = idempotente (reemplaza, no acumula)
	if _, err := c.Materializar(context.Background(), src, "x", "1.0.0"); err != nil {
		t.Fatalf("re-materializar: %v", err)
	}
}
