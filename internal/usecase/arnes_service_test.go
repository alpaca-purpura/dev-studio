package usecase

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

// --- dobles de puertos ---

type fakeSync struct{ dir string }

func (f *fakeSync) Sync(_ context.Context, source string) (string, error) {
	if source == "roto" {
		return "", errors.New("source inválido")
	}
	return f.dir, nil
}

type fakeCatalog struct {
	arneses []domain.Arnes
	sources map[string]string       // id -> dir forma-plugin
	metas   map[string]domain.Arnes // dir -> meta (LeerArnes)
}

func (f *fakeCatalog) Leer(context.Context, string) ([]domain.Arnes, error) { return f.arneses, nil }
func (f *fakeCatalog) SourceDir(_, id string) (string, error) {
	d, ok := f.sources[id]
	if !ok {
		return "", errors.New("no está")
	}
	return d, nil
}
func (f *fakeCatalog) LeerArnes(dir string) (domain.Arnes, error) {
	m, ok := f.metas[dir]
	if !ok {
		return domain.Arnes{}, errors.New("dir sin arnés")
	}
	return m, nil
}

type memLock struct {
	registry string
	inst     []domain.ArnesInstalado
	path     string
}

func (m *memLock) Leer(context.Context, string) (string, []domain.ArnesInstalado, error) {
	return m.registry, append([]domain.ArnesInstalado{}, m.inst...), nil
}
func (m *memLock) Escribir(_ context.Context, repoRoot, registry string, inst []domain.ArnesInstalado) (string, error) {
	m.registry, m.inst = registry, inst
	m.path = filepath.Join(repoRoot, ".devstudio", "arneses.yaml")
	return m.path, nil
}

type fakeCache struct{ base string }

func (f *fakeCache) Dir(id, version string) string { return filepath.Join(f.base, id, version) }
func (f *fakeCache) Materializar(_ context.Context, _, id, version string) (string, error) {
	d := f.Dir(id, version)
	if err := os.MkdirAll(d, 0o755); err != nil {
		return "", err
	}
	return d, nil
}

type spyCommit struct {
	paths    []string
	mensajes []string
}

func (s *spyCommit) Commit(_ context.Context, _ string, paths []string, mensaje string) (string, error) {
	s.paths = append(s.paths, paths...)
	s.mensajes = append(s.mensajes, mensaje)
	return "abc123", nil
}

type memConfig struct{ source string }

func (m *memConfig) LoadRegistrySource(context.Context) (string, error) { return m.source, nil }
func (m *memConfig) SaveRegistrySource(_ context.Context, s string) error {
	m.source = s
	return nil
}

func newTestArnesService(t *testing.T) (*ArnesService, *spyCommit, *memLock, *fakeCache) {
	t.Helper()
	dfc := domain.Arnes{ID: "dev-full-cycle", Nombre: "dev-full-cycle", Rol: "Ingeniería", Proceso: "e2e", Version: "0.1.0", Canal: "beta"}
	cacheBase := t.TempDir()
	cache := &fakeCache{base: cacheBase}
	cat := &fakeCatalog{
		arneses: []domain.Arnes{dfc},
		sources: map[string]string{"dev-full-cycle": "/mp/plugins/dev-full-cycle/0.1.0"},
		metas: map[string]domain.Arnes{
			filepath.Join(cacheBase, "dev-full-cycle", "0.1.0"): dfc,
		},
	}
	lock := &memLock{}
	commit := &spyCommit{}
	svc, err := NewArnesService(context.Background(), &fakeSync{dir: "/mp"}, cat, lock, cache, commit, &memConfig{})
	if err != nil {
		t.Fatal(err)
	}
	return svc, commit, lock, cache
}

func TestConectarYEstado(t *testing.T) {
	svc, _, _, _ := newTestArnesService(t)
	est, err := svc.Conectar(context.Background(), "/mp")
	if err != nil {
		t.Fatalf("Conectar: %v", err)
	}
	if est.Source != "/mp" || len(est.Arneses) != 1 {
		t.Fatalf("estado: %+v", est)
	}
	if _, err := svc.Conectar(context.Background(), "roto"); err == nil {
		t.Fatal("source roto debía fallar")
	}
	// el fallo NO pisa la conexión previa
	if svc.Estado().Source != "/mp" {
		t.Fatalf("conexión previa perdida: %+v", svc.Estado())
	}
}

func TestInstalarActualizaLockYCommitea(t *testing.T) {
	svc, commit, lock, _ := newTestArnesService(t)
	if _, err := svc.Conectar(context.Background(), "/mp"); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	a, err := svc.Instalar(context.Background(), repo, "dev-full-cycle")
	if err != nil {
		t.Fatalf("Instalar: %v", err)
	}
	if a.ID != "dev-full-cycle" {
		t.Fatalf("arnés: %+v", a)
	}
	if len(lock.inst) != 1 || lock.inst[0].Version != "0.1.0" || lock.registry != "/mp" {
		t.Fatalf("lock: %+v reg=%q", lock.inst, lock.registry)
	}
	// commit SOLO del lock, por pathspec (RN-2)
	if len(commit.paths) != 1 || commit.paths[0] != filepath.Join(".devstudio", "arneses.yaml") {
		t.Fatalf("pathspec del commit: %v", commit.paths)
	}
	if !strings.Contains(commit.mensajes[0], "instala arnés dev-full-cycle@0.1.0") {
		t.Fatalf("mensaje: %v", commit.mensajes)
	}

	// re-instalar = idempotente (reemplaza la entrada, no duplica)
	if _, err := svc.Instalar(context.Background(), repo, "dev-full-cycle"); err != nil {
		t.Fatal(err)
	}
	if len(lock.inst) != 1 {
		t.Fatalf("instalación duplicada: %+v", lock.inst)
	}

	// id fuera del catálogo → error tipado (422 en transporte)
	if _, err := svc.Instalar(context.Background(), repo, "nope"); !errors.Is(err, ErrArnesNoEncontrado) {
		t.Fatalf("esperaba ErrArnesNoEncontrado: %v", err)
	}
}

func TestInstalarSinRegistryFalla(t *testing.T) {
	svc, _, _, _ := newTestArnesService(t)
	if _, err := svc.Instalar(context.Background(), t.TempDir(), "dev-full-cycle"); !errors.Is(err, ErrRegistryNoConectado) {
		t.Fatalf("esperaba ErrRegistryNoConectado: %v", err)
	}
}

func TestDesinstalar(t *testing.T) {
	svc, commit, lock, _ := newTestArnesService(t)
	if _, err := svc.Conectar(context.Background(), "/mp"); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	if _, err := svc.Instalar(context.Background(), repo, "dev-full-cycle"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Desinstalar(context.Background(), repo, "dev-full-cycle"); err != nil {
		t.Fatalf("Desinstalar: %v", err)
	}
	if len(lock.inst) != 0 {
		t.Fatalf("lock no vació: %+v", lock.inst)
	}
	if !strings.Contains(commit.mensajes[len(commit.mensajes)-1], "desinstala arnés dev-full-cycle@0.1.0") {
		t.Fatalf("mensaje: %v", commit.mensajes)
	}
	if err := svc.Desinstalar(context.Background(), repo, "dev-full-cycle"); !errors.Is(err, ErrArnesNoInstalado) {
		t.Fatalf("desinstalar lo no instalado: %v", err)
	}
}

func TestInjectionExigeInstalado(t *testing.T) {
	svc, _, _, _ := newTestArnesService(t)
	if _, err := svc.Conectar(context.Background(), "/mp"); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()

	// RN-5: rol no instalado → error tipado
	if _, err := svc.Injection(context.Background(), repo, "dev-full-cycle"); !errors.Is(err, ErrArnesNoInstalado) {
		t.Fatalf("esperaba ErrArnesNoInstalado: %v", err)
	}

	if _, err := svc.Instalar(context.Background(), repo, "dev-full-cycle"); err != nil {
		t.Fatal(err)
	}
	inj, err := svc.Injection(context.Background(), repo, "dev-full-cycle")
	if err != nil {
		t.Fatalf("Injection: %v", err)
	}
	if inj.PluginDir == "" || !strings.Contains(inj.SystemPrompt, "Ingeniería") {
		t.Fatalf("inyección incompleta: %+v", inj)
	}
	// sin CLAUDE.md en el caché fake → banda Base ausente del prompt, honesto
	if strings.Contains(inj.SystemPrompt, "Banda Base") {
		t.Fatalf("banda Base fantasma sin CLAUDE.md: %+v", inj)
	}
}

func TestInjectionEmbebeBandaBase(t *testing.T) {
	svc, _, _, cache := newTestArnesService(t)
	if _, err := svc.Conectar(context.Background(), "/mp"); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	if _, err := svc.Instalar(context.Background(), repo, "dev-full-cycle"); err != nil {
		t.Fatal(err)
	}
	// el arnés trae CLAUDE.md (banda Base) → viaja EMBEBIDO en el system prompt
	if err := os.WriteFile(filepath.Join(cache.Dir("dev-full-cycle", "0.1.0"), "CLAUDE.md"), []byte("# std-spec\nregla base"), 0o644); err != nil {
		t.Fatal(err)
	}
	inj, err := svc.Injection(context.Background(), repo, "dev-full-cycle")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(inj.SystemPrompt, "Banda Base") || !strings.Contains(inj.SystemPrompt, "regla base") {
		t.Fatalf("banda Base no embebida: %q", inj.SystemPrompt)
	}
}

func TestInjectionRehidrataCacheAusente(t *testing.T) {
	svc, _, lock, cache := newTestArnesService(t)
	if _, err := svc.Conectar(context.Background(), "/mp"); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	// lock traído por git (otro usuario instaló) — caché local vacío
	lock.inst = []domain.ArnesInstalado{{ID: "dev-full-cycle", Version: "0.1.0", Canal: "beta"}}
	if _, err := os.Stat(cache.Dir("dev-full-cycle", "0.1.0")); err == nil {
		t.Fatal("precondición: caché debía estar vacío")
	}
	inj, err := svc.Injection(context.Background(), repo, "dev-full-cycle")
	if err != nil {
		t.Fatalf("Injection con rehidratación: %v", err)
	}
	if inj.PluginDir != cache.Dir("dev-full-cycle", "0.1.0") {
		t.Fatalf("no rehidrató al caché: %+v", inj)
	}
}
