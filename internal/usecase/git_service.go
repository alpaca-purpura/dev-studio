package usecase

import (
	"context"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// GitService expone las lecturas git + el commit explícito al transporte, manteniendo la
// dirección hexagonal (transport → usecase → ports). Es deliberadamente delgado: la política
// (solo lectura + commit por pathspec) vive en los PUERTOS — acá no se agrega superficie.
type GitService struct {
	info   ports.GitInfo
	commit ports.GitCommit
}

func NewGitService(info ports.GitInfo, commit ports.GitCommit) *GitService {
	return &GitService{info: info, commit: commit}
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
