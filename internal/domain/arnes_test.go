package domain

import (
	"strings"
	"testing"
)

func TestPreambuloComponeDelMeta(t *testing.T) {
	a := Arnes{
		ID:      "dev-full-cycle",
		Rol:     "Ingeniería · Desarrollo full-cycle",
		Proceso: "desarrollo de software end-to-end (idea → released)",
		Fases:   []string{"spec", "build", "review", "release"},
	}
	p := a.Preambulo()
	for _, quiere := range []string{
		"Ingeniería · Desarrollo full-cycle",
		"desarrollo de software end-to-end",
		"spec → build → review → release",
	} {
		if !strings.Contains(p, quiere) {
			t.Errorf("preámbulo sin %q:\n%s", quiere, p)
		}
	}
}

func TestPreambuloSinProcesoNiFases(t *testing.T) {
	p := Arnes{Rol: "Auditor"}.Preambulo()
	if !strings.Contains(p, "Auditor") || strings.Contains(p, "Fases") {
		t.Errorf("preámbulo mínimo mal compuesto: %s", p)
	}
}
