package ports

import (
	"context"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

// RegistrySync obtiene una copia local del marketplace git de la organización. El
// adaptador está CONFINADO a ~/.dev-studio/registry/ — jamás opera sobre repos del
// usuario (nota v1.1 del boundary git-solo-lectura-y-commit; RN-3 spec registry-arneses).
// source = URL git (se clona/actualiza con el git del usuario, BYO credenciales) o ruta
// local directa (se usa tal cual, sin git).
type RegistrySync interface {
	Sync(ctx context.Context, source string) (dir string, err error)
}

// RegistryCatalog lee el catálogo de arneses de un marketplace ya materializado en disco
// (formato prenter-marketplace: .claude-plugin/marketplace.json + por arnés
// arnes.l0.json/plugin.json — nomenclatura-arnes v1 de la fábrica).
type RegistryCatalog interface {
	Leer(ctx context.Context, dir string) ([]domain.Arnes, error)
	// SourceDir resuelve el dir forma-plugin de un arnés del índice (para materializar).
	SourceDir(dir, arnesID string) (string, error)
	// LeerArnes lee el meta de UN dir forma-plugin (arnes.l0.json + plugin.json) — lo usa
	// la inyección para componer el preámbulo desde el caché.
	LeerArnes(dir string) (domain.Arnes, error)
}

// RegistryConfigStore persiste la conexión al registry (a nivel app — norte DH-14:
// «al instalar la aplicación me conectaré a un repositorio de plugins»).
type RegistryConfigStore interface {
	LoadRegistrySource(ctx context.Context) (string, error)
	SaveRegistrySource(ctx context.Context, source string) error
}

// ArnesLock lee/escribe el roster as-code del proyecto (.devstudio/arneses.yaml) — el
// ÚNICO archivo que DevStudio escribe en el repo del usuario (RN-2).
type ArnesLock interface {
	Leer(ctx context.Context, repoRoot string) (registry string, instalados []domain.ArnesInstalado, err error)
	Escribir(ctx context.Context, repoRoot, registry string, instalados []domain.ArnesInstalado) (lockPath string, err error)
}

// ArnesCache materializa la forma-plugin INTACTA de un arnés al caché local
// (~/.dev-studio/arneses/{id}/{version}) desde el marketplace sincronizado. Solo copia,
// jamás edita (RN-4).
type ArnesCache interface {
	Materializar(ctx context.Context, srcDir, id, version string) (dir string, err error)
	Dir(id, version string) string
}
