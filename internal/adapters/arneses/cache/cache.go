// Package cache materializa la forma-plugin INTACTA de un arnés al caché local
// (~/.dev-studio/arneses/{id}/{version}). Copia fiel del registry — jamás edita (RN-4):
// editar arneses es trabajo de la fábrica (ArnesIA), no de la app de rol.
package cache

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alpacapurpura/dev-studio/internal/ports"
)

type Cache struct {
	base string // ~/.dev-studio/arneses
}

var _ ports.ArnesCache = (*Cache)(nil)

func New(baseDir string) *Cache { return &Cache{base: baseDir} }

func (c *Cache) Dir(id, version string) string {
	return filepath.Join(c.base, id, version)
}

// Materializar copia el árbol forma-plugin al caché. Idempotente: un destino previo se
// reemplaza completo (jamás merge parcial — el caché refleja EXACTO lo publicado).
func (c *Cache) Materializar(_ context.Context, srcDir, id, version string) (string, error) {
	if id == "" || version == "" {
		return "", fmt.Errorf("cache: id/version vacíos")
	}
	dest := c.Dir(id, version)
	if err := os.RemoveAll(dest); err != nil {
		return "", fmt.Errorf("cache: limpiar destino: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", fmt.Errorf("cache: crear base: %w", err)
	}
	if err := os.CopyFS(dest, os.DirFS(srcDir)); err != nil {
		return "", fmt.Errorf("cache: copiar forma-plugin: %w", err)
	}
	return dest, nil
}
