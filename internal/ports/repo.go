package ports

import (
	"context"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

// RepoStore persiste el registro de repositorios (junto a las sesiones, en el state de la app).
type RepoStore interface {
	LoadRepos(ctx context.Context) ([]domain.Repo, error)
	SaveRepos(ctx context.Context, repos []domain.Repo) error
}
