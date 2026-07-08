// Package lockfile lee/escribe el roster as-code del proyecto: .devstudio/arneses.yaml —
// el ÚNICO archivo que DevStudio escribe en el repo del usuario (RN-2 spec
// registry-arneses). Modelo npm: el lock viaja por GitHub (multi-usuario), el payload se
// rehidrata del registry.
package lockfile

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
)

type Lock struct{}

var _ ports.ArnesLock = (*Lock)(nil)

func New() *Lock { return &Lock{} }

type lockFile struct {
	Registry string                  `yaml:"registry"`
	Arneses  []domain.ArnesInstalado `yaml:"arneses"`
}

func path(repoRoot string) string {
	return filepath.Join(repoRoot, ".devstudio", "arneses.yaml")
}

// Leer devuelve el roster del repo. Sin lock = roster vacío honesto, no error.
func (l *Lock) Leer(_ context.Context, repoRoot string) (string, []domain.ArnesInstalado, error) {
	raw, err := os.ReadFile(path(repoRoot))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, fmt.Errorf("lockfile: leer: %w", err)
	}
	var f lockFile
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return "", nil, fmt.Errorf("lockfile: arneses.yaml ilegible: %w", err)
	}
	return f.Registry, f.Arneses, nil
}

func (l *Lock) Escribir(_ context.Context, repoRoot, registry string, instalados []domain.ArnesInstalado) (string, error) {
	if instalados == nil {
		instalados = []domain.ArnesInstalado{}
	}
	raw, err := yaml.Marshal(lockFile{Registry: registry, Arneses: instalados})
	if err != nil {
		return "", fmt.Errorf("lockfile: marshal: %w", err)
	}
	cabecera := "# Roster de arneses del proyecto — escrito por DevStudio (PB-25).\n" +
		"# rol = arnés instalado desde el registry de la organización (DH-14).\n" +
		"# El payload NO viaja acá: se rehidrata del registry al abrir el proyecto.\n"
	p := path(repoRoot)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", fmt.Errorf("lockfile: crear .devstudio: %w", err)
	}
	if err := os.WriteFile(p, append([]byte(cabecera), raw...), 0o644); err != nil {
		return "", fmt.Errorf("lockfile: escribir: %w", err)
	}
	return p, nil
}
