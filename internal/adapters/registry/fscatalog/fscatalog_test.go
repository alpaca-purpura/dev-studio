package fscatalog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLeerCatalogoCompleto(t *testing.T) {
	c := New()
	arneses, err := c.Leer(context.Background(), "testdata/marketplace")
	if err != nil {
		t.Fatalf("Leer: %v", err)
	}
	if len(arneses) != 2 {
		t.Fatalf("esperaba 2 arneses, obtuve %d: %+v", len(arneses), arneses)
	}

	dfc := arneses[0]
	if dfc.ID != "dev-full-cycle" {
		t.Fatalf("primer arnés: %+v", dfc)
	}
	if dfc.Rol != "Ingeniería · Desarrollo full-cycle" {
		t.Errorf("rol del manifiesto no leído: %q", dfc.Rol)
	}
	if dfc.Proceso == "" || dfc.Canal != "beta" || dfc.Version != "0.1.0" {
		t.Errorf("meta incompleto: %+v", dfc)
	}
	if len(dfc.Fases) != 4 || dfc.Fases[0] != "spec" {
		t.Errorf("fases no leídas: %v", dfc.Fases)
	}
	// descripción cae a plugin.json (arnes.l0.json no la trae)
	if dfc.Descripcion == "" {
		t.Error("descripción no derivada de plugin.json")
	}
}

func TestLeerDegradaSinManifiesto(t *testing.T) {
	// demo-auditor no tiene arnes.l0.json → degradación honesta a plugin.json:
	// id/nombre/descripcion/version presentes, rol cae al nombre.
	c := New()
	arneses, err := c.Leer(context.Background(), "testdata/marketplace")
	if err != nil {
		t.Fatalf("Leer: %v", err)
	}
	da := arneses[1]
	if da.ID != "demo-auditor" || da.Version != "0.2.0" || da.Rol == "" {
		t.Fatalf("degradación sin manifiesto mal resuelta: %+v", da)
	}
}

// Cadena de fallback canónica BENDECIDA por ArnesIA (interop 2026-07-07, ficha HS-12):
// nombre → plugin.json name → id (ídem descripcion). arnes.l0.nombre es el campo
// canónico para pintar — si el manifiesto lo trae, gana sobre plugin.json.
func TestCadenaFallbackNombreDescripcion(t *testing.T) {
	escribir := func(t *testing.T, dir, rel, contenido string) {
		t.Helper()
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(contenido), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("l0.nombre y l0.descripcion ganan sobre plugin.json", func(t *testing.T) {
		dir := t.TempDir()
		escribir(t, dir, "arnes.l0.json",
			`{"id":"x","nombre":"Nombre L0","descripcion":"Descripción L0","rol":"r"}`)
		escribir(t, dir, ".claude-plugin/plugin.json",
			`{"name":"x","description":"descr plugin","version":"1.0.0"}`)
		a, err := leerArnes(dir)
		if err != nil {
			t.Fatal(err)
		}
		if a.Nombre != "Nombre L0" || a.Descripcion != "Descripción L0" {
			t.Fatalf("l0 no gana la cadena: %+v", a)
		}
	})

	t.Run("sin nombre en ningún lado cae al id", func(t *testing.T) {
		dir := t.TempDir()
		escribir(t, dir, "arnes.l0.json", `{"id":"solo-id","rol":"r"}`)
		a, err := leerArnes(dir)
		if err != nil {
			t.Fatal(err)
		}
		if a.Nombre != "solo-id" {
			t.Fatalf("nombre debía caer al id: %+v", a)
		}
	})
}

func TestLeerRechazaDirSinMarketplace(t *testing.T) {
	c := New()
	if _, err := c.Leer(context.Background(), t.TempDir()); err == nil {
		t.Fatal("dir sin .claude-plugin/marketplace.json debía fallar (error honesto, jamás catálogo vacío inventado)")
	}
}

func TestSourceDir(t *testing.T) {
	c := New()
	dir, err := c.SourceDir("testdata/marketplace", "dev-full-cycle")
	if err != nil {
		t.Fatalf("SourceDir: %v", err)
	}
	if dir != "testdata/marketplace/plugins/dev-full-cycle/0.1.0" {
		t.Fatalf("source dir: %s", dir)
	}
	if _, err := c.SourceDir("testdata/marketplace", "no-existe"); err == nil {
		t.Fatal("arnés inexistente debía fallar")
	}
}
