package ports

import (
	"context"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

// GitInfo lee el estado git de un cwd — SOLO lectura. Ningún método de este puerto ni de
// GitCommit toca el remoto ni reescribe historia (boundary git-solo-lectura-y-commit, RN-4):
// push/pull/fetch/reset/rebase NO EXISTEN en el contrato.
type GitInfo interface {
	Status(ctx context.Context, cwd string) (domain.GitStatus, error)
	DiffFile(ctx context.Context, cwd, path string) (domain.GitDiff, error)
	Log(ctx context.Context, cwd string, n int) ([]domain.GitLogEntry, error)
}

// GitCommit registra un commit local. La firma exige pathspec explícito (RN-5): commitear
// «todo» sin enumerarlo es imposible por diseño — paths vacío = error del adaptador.
type GitCommit interface {
	Commit(ctx context.Context, cwd string, paths []string, mensaje string) (sha string, err error)
}

// GitWorkspace crea/borra los workspaces aislados por sesión (spec workspace-aislado, PB-02):
// worktree + branch wt/{slug}. Remove sin force — un worktree sucio se conserva, jamás se
// pierde trabajo en silencio (RN-3).
type GitWorkspace interface {
	CreateWorktree(ctx context.Context, repoRoot, destino, slug string) (path, branch string, err error)
	RemoveWorktree(ctx context.Context, repoRoot, path string) error
}
