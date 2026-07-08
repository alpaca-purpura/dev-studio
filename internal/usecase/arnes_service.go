package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// ErrRegistryNoConectado: operación que exige un registry conectado.
var ErrRegistryNoConectado = errors.New("registry no conectado")

// ErrArnesNoEncontrado: el id no está en el catálogo del registry.
var ErrArnesNoEncontrado = errors.New("arnés no está en el catálogo")

// ErrArnesNoInstalado: el rol pedido no está instalado en el proyecto (RN-5).
var ErrArnesNoInstalado = errors.New("arnés no instalado en el proyecto")

// ArnesInjection es lo que el spawn de una sesión con rol suma al conductor (PB-25).
// SystemPrompt = preámbulo del rol + banda Base (CLAUDE.md del arnés) embebida — la CLI
// exige un solo flag de system prompt.
type ArnesInjection struct {
	PluginDir    string
	SystemPrompt string
}

// ArnesService orquesta el registry de arneses: conexión (a nivel app), catálogo,
// instalación por proyecto (lock as-code + caché) e inyección por sesión.
type ArnesService struct {
	mu      sync.Mutex
	source  string
	dir     string // marketplace materializado (post-sync)
	catalog []domain.Arnes

	sync    ports.RegistrySync
	cat     ports.RegistryCatalog
	lock    ports.ArnesLock
	cache   ports.ArnesCache
	commit  ports.GitCommit
	config  ports.RegistryConfigStore
	baseCtx context.Context
}

func NewArnesService(baseCtx context.Context, rs ports.RegistrySync, cat ports.RegistryCatalog, lock ports.ArnesLock, cache ports.ArnesCache, commit ports.GitCommit, config ports.RegistryConfigStore) (*ArnesService, error) {
	s := &ArnesService{
		sync: rs, cat: cat, lock: lock, cache: cache, commit: commit,
		config: config, baseCtx: baseCtx,
	}
	source, err := config.LoadRegistrySource(baseCtx)
	if err != nil {
		return nil, fmt.Errorf("arnes service: load registry source: %w", err)
	}
	if source != "" {
		s.source = source
		// sync perezoso: si el registry no está disponible al boot, el estado lo dice —
		// la app arranca igual (honestidad > bloqueo).
		if dir, err := rs.Sync(baseCtx, source); err == nil {
			if arneses, err := cat.Leer(baseCtx, dir); err == nil {
				s.dir, s.catalog = dir, arneses
			}
		}
	}
	return s, nil
}

// Conectar fija el source del registry (URL git o ruta local), lo sincroniza y persiste.
func (s *ArnesService) Conectar(ctx context.Context, source string) (domain.RegistryEstado, error) {
	dir, err := s.sync.Sync(ctx, source)
	if err != nil {
		return domain.RegistryEstado{}, err
	}
	arneses, err := s.cat.Leer(ctx, dir)
	if err != nil {
		return domain.RegistryEstado{}, err
	}
	if err := s.config.SaveRegistrySource(ctx, source); err != nil {
		return domain.RegistryEstado{}, err
	}
	s.mu.Lock()
	s.source, s.dir, s.catalog = source, dir, arneses
	s.mu.Unlock()
	return domain.RegistryEstado{Source: source, Arneses: arneses}, nil
}

// Estado devuelve la conexión actual + catálogo (vacío honesto si no hay registry).
func (s *ArnesService) Estado() domain.RegistryEstado {
	s.mu.Lock()
	defer s.mu.Unlock()
	return domain.RegistryEstado{Source: s.source, Arneses: s.catalog}
}

// Resync re-sincroniza el registry conectado (pull) y refresca el catálogo.
func (s *ArnesService) Resync(ctx context.Context) (domain.RegistryEstado, error) {
	s.mu.Lock()
	source := s.source
	s.mu.Unlock()
	if source == "" {
		return domain.RegistryEstado{}, ErrRegistryNoConectado
	}
	return s.Conectar(ctx, source)
}

// Instalados devuelve el roster del proyecto (lock), enriquecido con el meta del
// catálogo cuando el registry lo tiene (si no, degradación honesta al dato del lock).
func (s *ArnesService) Instalados(ctx context.Context, repoRoot string) ([]domain.Arnes, error) {
	_, instalados, err := s.lock.Leer(ctx, repoRoot)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	catalog := s.catalog
	s.mu.Unlock()
	out := make([]domain.Arnes, 0, len(instalados))
	for _, in := range instalados {
		a := domain.Arnes{ID: in.ID, Nombre: in.ID, Rol: in.ID, Version: in.Version, Canal: in.Canal}
		// el meta rico vive en el caché (forma-plugin materializada) — fuente primaria
		if meta, err := s.cat.LeerArnes(s.cache.Dir(in.ID, in.Version)); err == nil {
			a = meta
			a.Version, a.Canal = in.Version, in.Canal
		} else {
			for _, c := range catalog {
				if c.ID == in.ID {
					a = c
					a.Version, a.Canal = in.Version, in.Canal
					break
				}
			}
		}
		out = append(out, a)
	}
	return out, nil
}

// Instalar materializa la forma-plugin al caché, actualiza el lock del repo y lo
// commitea por pathspec (RN-2: único archivo tocado; jamás `git add .`).
func (s *ArnesService) Instalar(ctx context.Context, repoRoot, arnesID string) (domain.Arnes, error) {
	s.mu.Lock()
	source, dir, catalog := s.source, s.dir, s.catalog
	s.mu.Unlock()
	if source == "" || dir == "" {
		return domain.Arnes{}, ErrRegistryNoConectado
	}
	var arnes domain.Arnes
	found := false
	for _, a := range catalog {
		if a.ID == arnesID {
			arnes, found = a, true
			break
		}
	}
	if !found {
		return domain.Arnes{}, fmt.Errorf("%w: %s", ErrArnesNoEncontrado, arnesID)
	}

	src, err := s.cat.SourceDir(dir, arnesID)
	if err != nil {
		return domain.Arnes{}, err
	}
	if _, err := s.cache.Materializar(ctx, src, arnes.ID, arnes.Version); err != nil {
		return domain.Arnes{}, err
	}

	_, instalados, err := s.lock.Leer(ctx, repoRoot)
	if err != nil {
		return domain.Arnes{}, err
	}
	entry := domain.ArnesInstalado{ID: arnes.ID, Version: arnes.Version, Canal: arnes.Canal}
	replaced := false
	for i, in := range instalados {
		if in.ID == arnes.ID {
			instalados[i], replaced = entry, true
			break
		}
	}
	if !replaced {
		instalados = append(instalados, entry)
	}
	if err := s.escribirLockYCommitear(ctx, repoRoot, source, instalados,
		fmt.Sprintf("chore(devstudio): instala arnés %s@%s", arnes.ID, arnes.Version)); err != nil {
		return domain.Arnes{}, err
	}
	return arnes, nil
}

// Desinstalar quita el arnés del lock (el caché se conserva) y commitea.
func (s *ArnesService) Desinstalar(ctx context.Context, repoRoot, arnesID string) error {
	registry, instalados, err := s.lock.Leer(ctx, repoRoot)
	if err != nil {
		return err
	}
	kept := instalados[:0]
	var quitado *domain.ArnesInstalado
	for _, in := range instalados {
		if in.ID == arnesID {
			q := in
			quitado = &q
			continue
		}
		kept = append(kept, in)
	}
	if quitado == nil {
		return fmt.Errorf("%w: %s", ErrArnesNoInstalado, arnesID)
	}
	return s.escribirLockYCommitear(ctx, repoRoot, registry, kept,
		fmt.Sprintf("chore(devstudio): desinstala arnés %s@%s", quitado.ID, quitado.Version))
}

// VerificarInstalado valida que el rol esté en el lock del repo (guard RN-5 del
// transporte al crear sesión — 422 antes de crear nada).
func (s *ArnesService) VerificarInstalado(ctx context.Context, repoRoot, rolID string) error {
	_, instalados, err := s.lock.Leer(ctx, repoRoot)
	if err != nil {
		return err
	}
	for _, in := range instalados {
		if in.ID == rolID {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrArnesNoInstalado, rolID)
}

// Injection resuelve la inyección del rol para una sesión (RN-5): exige el arnés en el
// lock del repo; si el caché falta, REHIDRATA del registry (modelo npm).
func (s *ArnesService) Injection(ctx context.Context, repoRoot, rolID string) (ArnesInjection, error) {
	_, instalados, err := s.lock.Leer(ctx, repoRoot)
	if err != nil {
		return ArnesInjection{}, err
	}
	var in *domain.ArnesInstalado
	for i := range instalados {
		if instalados[i].ID == rolID {
			in = &instalados[i]
			break
		}
	}
	if in == nil {
		return ArnesInjection{}, fmt.Errorf("%w: %s", ErrArnesNoInstalado, rolID)
	}

	dir := s.cache.Dir(in.ID, in.Version)
	if _, err := os.Stat(dir); err != nil {
		// rehidratación: otro usuario clonó el repo con el lock pero sin caché local
		s.mu.Lock()
		mpDir := s.dir
		s.mu.Unlock()
		if mpDir == "" {
			return ArnesInjection{}, fmt.Errorf("%w (y el arnés %s no está en caché)", ErrRegistryNoConectado, in.ID)
		}
		src, err := s.cat.SourceDir(mpDir, in.ID)
		if err != nil {
			return ArnesInjection{}, err
		}
		if dir, err = s.cache.Materializar(ctx, src, in.ID, in.Version); err != nil {
			return ArnesInjection{}, err
		}
	}

	arnes, err := s.cat.LeerArnes(dir)
	if err != nil {
		return ArnesInjection{}, fmt.Errorf("arnés %s en caché ilegible: %w", in.ID, err)
	}
	prompt := arnes.Preambulo()
	if base, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md")); err == nil && len(base) > 0 {
		prompt += "\n\n## Banda Base del arnés\n\n" + string(base)
	}
	return ArnesInjection{PluginDir: dir, SystemPrompt: prompt}, nil
}

func (s *ArnesService) escribirLockYCommitear(ctx context.Context, repoRoot, registry string, instalados []domain.ArnesInstalado, mensaje string) error {
	if _, err := s.lock.Escribir(ctx, repoRoot, registry, instalados); err != nil {
		return err
	}
	lockRel := filepath.Join(".devstudio", "arneses.yaml")
	if _, err := s.commit.Commit(ctx, repoRoot, []string{lockRel}, mensaje); err != nil {
		return fmt.Errorf("commitear lock del roster: %w", err)
	}
	return nil
}
