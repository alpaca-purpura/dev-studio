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
