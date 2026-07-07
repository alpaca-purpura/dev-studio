package usecase

import (
	"context"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// GitService expone las lecturas git + el commit explícito + los workspaces aislados al
// transporte, manteniendo la dirección hexagonal (transport → usecase → ports). Es
// deliberadamente delgado: la política (solo lectura + commit por pathspec + worktree sin
// force) vive en los PUERTOS — acá no se agrega superficie.
type GitService struct {
	info   ports.GitInfo
	commit ports.GitCommit
	ws     ports.GitWorkspace
}

func NewGitService(info ports.GitInfo, commit ports.GitCommit, ws ports.GitWorkspace) *GitService {
	return &GitService{info: info, commit: commit, ws: ws}
}

func (g *GitService) CreateWorktree(ctx context.Context, repoRoot, destino, slug string) (string, string, error) {
	return g.ws.CreateWorktree(ctx, repoRoot, destino, slug)
}

func (g *GitService) RemoveWorktree(ctx context.Context, repoRoot, path string) error {
	return g.ws.RemoveWorktree(ctx, repoRoot, path)
}

func (g *GitService) Status(ctx context.Context, cwd string) (domain.GitStatus, error) {
	return g.info.Status(ctx, cwd)
}

func (g *GitService) DiffFile(ctx context.Context, cwd, path string) (domain.GitDiff, error) {
	return g.info.DiffFile(ctx, cwd, path)
}

func (g *GitService) Log(ctx context.Context, cwd string, n int) ([]domain.GitLogEntry, error) {
	return g.info.Log(ctx, cwd, n)
}

func (g *GitService) Commit(ctx context.Context, cwd string, paths []string, mensaje string) (string, error) {
	return g.commit.Commit(ctx, cwd, paths, mensaje)
}
