// Package fitness contiene los checks ejecutables que hacen cumplir arch/boundaries/*.md.
// Un boundary sin test acá (o sin enforcer equivalente) no puede decir `status: enforced`.
package fitness

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
	"github.com/alpacapurpura/dev-studio/internal/usecase"
)

// --- dobles mínimos para no depender de un `claude` real en el test ---

type noopPublisher struct{}

func (noopPublisher) Publish(string, []byte) {}

type memStore struct{}

func (memStore) Load(context.Context) ([]domain.Session, error) { return nil, nil }
func (memStore) Save(context.Context, []domain.Session) error   { return nil }

// blockingAgent nunca emite eventos: la sesión queda `streaming` para siempre, que es
// exactamente el estado que TestOneTurnAtATime necesita sostener.
type blockingAgent struct{}

func (blockingAgent) Spawn(context.Context, ports.SpawnOpts) (ports.AgentSession, error) {
	return &blockingSession{events: make(chan ports.AgentEvent)}, nil
}

type blockingSession struct{ events chan ports.AgentEvent }

func (b *blockingSession) Send(context.Context, string) error { return nil }
func (b *blockingSession) Events() <-chan ports.AgentEvent    { return b.events }
func (b *blockingSession) Close() error                       { close(b.events); return nil }

// TestOneTurnAtATime enforça sesion-un-turno-a-la-vez: un segundo Turn() mientras el
// primero sigue streaming debe rechazarse con ErrBusy, nunca intercalar dos escrituras.
func TestOneTurnAtATime(t *testing.T) {
	svc, err := usecase.NewSessionService(context.Background(), blockingAgent{}, memStore{}, noopPublisher{})
	if err != nil {
		t.Fatalf("new session service: %v", err)
	}
	sess, err := svc.Create("t", "/tmp")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := svc.Turn(sess.ID, "primero"); err != nil {
		t.Fatalf("primer turno debería aceptarse: %v", err)
	}
	if err := svc.Turn(sess.ID, "segundo"); !errors.Is(err, usecase.ErrBusy) {
		t.Fatalf("segundo turno concurrente: esperaba ErrBusy, obtuve %v", err)
	}
}

// TestDomainNoTransportImport enforça dominio-independiente-de-transporte: ni domain/ ni
// usecase/ pueden depender de net/http (transitivamente) — solo los adapters/transport/* lo hacen.
func TestDomainNoTransportImport(t *testing.T) {
	assertNoDep(t, "./internal/domain/...", "net/http")
	assertNoDep(t, "./internal/usecase/...", "net/http")
}

func assertNoDep(t *testing.T, pkg, forbidden string) {
	t.Helper()
	cmd := exec.Command("go", "list", "-deps", pkg)
	cmd.Dir = "../.." // arch/fitness -> raíz del módulo
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go list %s: %v\n%s", pkg, err, out)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == forbidden {
			t.Fatalf("%s depende de %s (directa o transitivamente) — viola dominio-independiente-de-transporte", pkg, forbidden)
		}
	}
}
