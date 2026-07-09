package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
)

type spyAgent struct{ last ports.SpawnOpts }

func (a *spyAgent) Spawn(_ context.Context, opts ports.SpawnOpts) (ports.AgentSession, error) {
	a.last = opts
	return &nullSession{events: make(chan ports.AgentEvent)}, nil
}

func (a *spyAgent) History(context.Context, string, string) ([]ports.TranscriptItem, error) {
	return nil, nil
}

type nullSession struct{ events chan ports.AgentEvent }

func (n *nullSession) Send(context.Context, string) error { return nil }
func (n *nullSession) Events() <-chan ports.AgentEvent    { return n.events }
func (n *nullSession) Close() error                       { close(n.events); return nil }

type nullStore struct{}

func (nullStore) Load(context.Context) ([]domain.Session, error) { return nil, nil }
func (nullStore) Save(context.Context, []domain.Session) error   { return nil }

// TestSpawnInyectaArnesDelRol (PB-25): una sesión con rol resuelve su arnés al spawn y
// el conductor recibe plugin-dir + preámbulo; sin rol, cero inyección.
func TestSpawnInyectaArnesDelRol(t *testing.T) {
	agent := &spyAgent{}
	svc, err := NewSessionService(context.Background(), agent, nullStore{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	svc.SetArnesResolver(func(repoID, rol string) (ArnesInjection, error) {
		if rol != "dev-full-cycle" {
			return ArnesInjection{}, errors.New("no instalado")
		}
		return ArnesInjection{PluginDir: "/caché/dfc/0.1.0", SystemPrompt: "Trabajás como Ingeniería."}, nil
	})

	sess, err := svc.CreateSession(NewSession{Nombre: "s", Cwd: "/tmp", RepoID: "r1", Rol: "dev-full-cycle"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Turn(sess.ID, "hola"); err != nil {
		t.Fatalf("turn: %v", err)
	}
	if len(agent.last.PluginDirs) != 1 || agent.last.PluginDirs[0] != "/caché/dfc/0.1.0" {
		t.Fatalf("plugin dir no inyectado: %+v", agent.last)
	}
	if agent.last.SystemPrompt == "" {
		t.Fatalf("preámbulo no inyectado: %+v", agent.last)
	}

	// sin rol = sin inyección
	libre, err := svc.CreateSession(NewSession{Nombre: "libre", Cwd: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Turn(libre.ID, "hola"); err != nil {
		t.Fatalf("turn libre: %v", err)
	}
	if len(agent.last.PluginDirs) != 0 || agent.last.SystemPrompt != "" {
		t.Fatalf("sesión sin rol no debía inyectar: %+v", agent.last)
	}

	// rol sin arnés instalado → el turno falla honesto (RN-5)
	rota, err := svc.CreateSession(NewSession{Nombre: "rota", Cwd: "/tmp", RepoID: "r1", Rol: "fantasma"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Turn(rota.ID, "hola"); err == nil {
		t.Fatal("rol no instalado debía fallar el spawn")
	}
}
