package ports

import (
	"context"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

// SessionStore persiste el registro de sesiones (metadata liviana, no la conversación completa).
type SessionStore interface {
	Load(ctx context.Context) ([]domain.Session, error)
	Save(ctx context.Context, sessions []domain.Session) error
}
