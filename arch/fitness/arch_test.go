// Package fitness contiene los checks ejecutables que hacen cumplir arch/boundaries/*.md.
// Un boundary sin test acá (o sin enforcer equivalente) no puede decir `status: enforced`.
package fitness

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	gitcli "github.com/alpacapurpura/dev-studio/internal/adapters/git/cli"
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

// TestGitAdapterSinVerbosProhibidos enforça git-solo-lectura-y-commit: el fuente del
// adaptador git no contiene los verbos que tocan el remoto o reescriben historia.
func TestGitAdapterSinVerbosProhibidos(t *testing.T) {
	dir := "../../internal/adapters/git/cli"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("leer %s: %v", dir, err)
	}
	forbidden := []string{"push", "pull", "fetch", "reset", "rebase"}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		src, err := os.ReadFile(dir + "/" + e.Name())
		if err != nil {
			t.Fatal(err)
		}
		// se escanea CÓDIGO, no comentarios (la doc del boundary nombra los verbos a propósito)
		var code strings.Builder
		for _, line := range strings.Split(string(src), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue
			}
			if i := strings.Index(line, "//"); i >= 0 && !strings.Contains(line[:i], `"`) {
				line = line[:i]
			}
			code.WriteString(strings.ToLower(line))
			code.WriteString("\n")
		}
		for _, verb := range forbidden {
			if strings.Contains(code.String(), verb) {
				t.Errorf("%s contiene el verbo prohibido %q — viola git-solo-lectura-y-commit", e.Name(), verb)
			}
		}
	}
}

// TestGitCommitExigePathspec enforça RN-5: commit sin paths explícitos = error, siempre.
func TestGitCommitExigePathspec(t *testing.T) {
	g := gitcli.New()
	if _, err := g.Commit(context.Background(), t.TempDir(), nil, "x"); err == nil {
		t.Fatal("Commit sin paths debía fallar (RN-5)")
	}
}

type memRepoStore struct{ repos []domain.Repo }

func (m *memRepoStore) LoadRepos(context.Context) ([]domain.Repo, error) { return m.repos, nil }
func (m *memRepoStore) SaveRepos(_ context.Context, r []domain.Repo) error {
	m.repos = r
	return nil
}

// TestRutasProtegidasRechazadas enforça sesion-aislada-por-cwd (mitad rutas): ninguna puerta
// acepta $HOME ni raíces de sistema como repo/cwd (spec workspace-aislado §2.3, RN-2).
func TestRutasProtegidasRechazadas(t *testing.T) {
	home := "/home/fitness-test"
	repos, err := usecase.NewRepoService(context.Background(), &memRepoStore{}, home)
	if err != nil {
		t.Fatal(err)
	}
	for _, ruta := range []string{home, "/etc", "/", home + "/.ssh", home + "/.dev-studio"} {
		if _, err := repos.Register(ruta); !errors.Is(err, usecase.ErrRutaProtegida) {
			t.Errorf("Register(%s): esperaba ErrRutaProtegida, obtuve %v", ruta, err)
		}
	}
	if err := usecase.ValidarRuta(home, "/etc/nginx"); !errors.Is(err, usecase.ErrRutaProtegida) {
		t.Errorf("ValidarRuta(/etc/nginx): esperaba ErrRutaProtegida, obtuve %v", err)
	}
	if err := usecase.ValidarRuta(home, home+"/Proyectos/x"); err != nil {
		t.Errorf("ruta legítima rechazada: %v", err)
	}
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
