// Package gitsync materializa el marketplace git de la organización en un caché local.
// CONFINADO POR CONSTRUCCIÓN a su baseDir (~/.dev-studio/registry/): el destino se deriva
// SIEMPRE dentro del base — este paquete no puede tocar repos del usuario (RN-3; nota
// v1.1 del boundary git-solo-lectura-y-commit: aquel scope = repo DEL usuario, este =
// estado propio de la app). Usa el git DEL usuario (BYO credenciales, espíritu DH-10).
// Verbos: clone + pull --ff-only. push/reset/rebase no existen acá tampoco (fitness).
package gitsync

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alpacapurpura/dev-studio/internal/ports"
)

type Sync struct {
	git  string
	base string // ~/.dev-studio/registry — TODO destino vive acá adentro
}

var _ ports.RegistrySync = (*Sync)(nil)

func New(git, baseDir string) *Sync { return &Sync{git: git, base: baseDir} }

// Sync resuelve el source a un dir local con el marketplace:
//   - ruta local existente → se valida y se usa TAL CUAL (dev-loop; sin git);
//   - URL git → clone al base confinado (o pull --ff-only si ya existe).
func (s *Sync) Sync(ctx context.Context, source string) (string, error) {
	if source == "" {
		return "", fmt.Errorf("gitsync: source vacío")
	}
	if esRutaLocal(source) {
		if err := validarMarketplace(source); err != nil {
			return "", err
		}
		return source, nil
	}

	dest := filepath.Join(s.base, destSlug(source))
	if _, err := os.Stat(filepath.Join(dest, ".git")); err == nil {
		if out, err := s.run(ctx, dest, "pull", "--ff-only"); err != nil {
			return "", fmt.Errorf("gitsync: actualizar registry: %w\n%s", err, out)
		}
	} else {
		if err := os.MkdirAll(s.base, 0o755); err != nil {
			return "", fmt.Errorf("gitsync: crear base: %w", err)
		}
		if out, err := s.run(ctx, s.base, "clone", source, dest); err != nil {
			return "", fmt.Errorf("gitsync: clonar registry: %w\n%s", err, out)
		}
	}
	if err := validarMarketplace(dest); err != nil {
		return "", err
	}
	return dest, nil
}

func (s *Sync) run(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, s.git, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// esRutaLocal: un path absoluto/relativo existente en disco (no una URL git).
func esRutaLocal(source string) bool {
	if strings.Contains(source, "://") || strings.HasPrefix(source, "git@") {
		return false
	}
	fi, err := os.Stat(source)
	return err == nil && fi.IsDir()
}

func validarMarketplace(dir string) error {
	if _, err := os.Stat(filepath.Join(dir, ".claude-plugin", "marketplace.json")); err != nil {
		return fmt.Errorf("gitsync: %s no es un marketplace de arneses (falta .claude-plugin/marketplace.json)", dir)
	}
	return nil
}

// destSlug deriva un nombre estable y legible del source para el dir de clone.
func destSlug(source string) string {
	nombre := strings.TrimSuffix(filepath.Base(source), ".git")
	nombre = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			return r
		case r >= 'A' && r <= 'Z':
			return r + ('a' - 'A')
		default:
			return '-'
		}
	}, nombre)
	h := sha256.Sum256([]byte(source))
	return nombre + "-" + hex.EncodeToString(h[:4])
}
