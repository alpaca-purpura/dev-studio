package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateWorktree implementa ports.GitWorkspace: checkout aislado + la branch pedida desde
// HEAD del repo (spec nuevo-workspace §3: el llamador pasa la branch completa con su prefijo
// estándar — feature/x, bugfix/x, explore/x…). Colisión → sufijo -2, -3…
func (g *Git) CreateWorktree(ctx context.Context, repoRoot, destino, branchBase string) (string, string, error) {
	// HEAD unborn (repo sin commits) → mensaje honesto antes de que git falle críptico
	if _, err := g.run(ctx, repoRoot, "rev-parse", "--verify", "HEAD"); err != nil {
		return "", "", fmt.Errorf("el repositorio necesita al menos un commit para crear workspaces")
	}
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return "", "", err
	}
	branch := branchBase
	path := destino
	for intento := 2; ; intento++ {
		_, err := g.run(ctx, repoRoot, "worktree", "add", path, "-b", branch)
		if err == nil {
			return path, branch, nil
		}
		msg := err.Error()
		colision := strings.Contains(msg, "already exists") || strings.Contains(msg, "ya existe")
		if !colision || intento > 20 {
			return "", "", err
		}
		branch = fmt.Sprintf("%s-%d", branchBase, intento)
		path = fmt.Sprintf("%s-%d", destino, intento)
	}
}

// RemoveWorktree borra el worktree SOLO si está limpio — git rechaza el remove sucio y acá
// jamás se usa --force (RN-3: cero pérdida silenciosa).
func (g *Git) RemoveWorktree(ctx context.Context, repoRoot, path string) error {
	_, err := g.run(ctx, repoRoot, "worktree", "remove", path)
	return err
}
