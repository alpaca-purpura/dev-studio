# DevStudio — construir y mantener software basado en proceso y arquitectura

Producto standalone (graduado del monorepo `prenter-harness`, 2026-07-04). Norte =
[`VISION.md`](./VISION.md) · registro = [`LEDGER.md`](./LEDGER.md) (fichas `DH-NN`, arranca en
DH-12; DH-01..DH-11 viven en la incubadora `prenter-harness/products/devhub/`, congelada).

**Orden de producto (SSoT, 2026-07-07):** todo lo pendiente = [`BACKLOG.md`](./BACKLOG.md)
(fichas `PB-NN`, prioridad por bandas — la firma Chris) · todo lo funcional =
[`INCREMENTO.md`](./INCREMENTO.md) (`CAP-NN` con evidencia real) · todo lo decidido/entregado
= [`LEDGER.md`](./LEDGER.md). Regla: idea nueva → PB antes de trabajarse; entrega verificada →
CAP en el mismo commit. Ritual: «¿qué viene?» → leer banda 🔴. **Este archivo describe SOLO el
estado vigente** — cero historial de fichas y cero «siguiente» (anti-sesgo: el pasado vive en
LEDGER, el futuro en BACKLOG). No confundir planos: PRODUCTO (VISION, arch/, código,
BACKLOG/INCREMENTO) vs PROYECTO (LEDGER, epicas/, .claude/, docs/process/, templates/,
scripts/, project.config.yaml, research/).

**Qué es:** aplicación de escritorio multiplataforma donde una organización construye y
mantiene software con proceso, arquitectura y documentación **as code**. Usuarios (CTO ·
developer · devops · product owner) trabajan orquestados sobre uno o varios sistemas; el
repositorio GitHub es el conector.

**Decisiones técnicas vigentes:**
- Conexión con Claude Code = **driver CLI-nativo**: spawn del `claude` del usuario vía
  stdin/stdout stream-json, un proceso por sesión. BYO licencia — **sin API de Anthropic**;
  la app no toca credenciales.
- Binario Go + UI embebida (`go:embed`); instalador = el binario solo (Win/Linux/mac).
- **Rol = arnés instalado desde el registry propio** (estándar ArnesIA `nomenclatura-arnes`
  v1 forma-plugin + marketplace git formato prenter-marketplace) — cero roles locales; la
  curaduría vive en el registry.
- Proceso as code = **se DERIVA del arnés**: spine + `spine.categorias` (enum fijo
  `propuesto·en-progreso·completado·descartado·pausado`, terminalidad derivada) + gates/dueños
  de los contratos por caja. El arnés NO shipea descriptor I-77 aparte (un I-77 materializado
  = export/proyección). La consola deriva TODO del dato, cero hardcode del ciclo.

**Arnés de construcción:** kit dev (plugin del marketplace `alpacapurpura/prenter-marketplace`,
pineado por versión). Evoluciona con el producto — mejoras al arnés se upstreamean al kit
(backflow I-59), jamás fork silencioso.

**Estado vigente** (cómo llegó acá → LEDGER · qué funciona con evidencia → INCREMENTO):

- **Stack:** módulo Go `github.com/alpacapurpura/dev-studio`, arquitectura hexagonal
  (`internal/domain` → `internal/ports` ← `internal/adapters` ← `internal/usecase`) + SPA
  React 19 + Zustand + Tailwind v4 embebida (binario único `cmd/dev-studio`). API REST + SSE
  (`internal/adapters/transport/{http,sse}`). Persistencia `~/.dev-studio/state.json`.
- **Sesiones/workspaces:** toda sesión de repo nace en su worktree propio
  (`~/.dev-studio/workspaces/{repo}/{slug}`) con rutas protegidas vedadas en toda puerta y
  cierre con modal conservar/borrar. Taxonomía de ítems `historia·bug·hotfix·tarea·spike` →
  branch `feature/ bugfix/ hotfix/ chore/ spike/`. Exploración = read-only
  (`--permission-mode plan`). RN-1: toda sesión CON EDICIÓN liga a un paquete de trabajo.
- **Arquitectura as code:** [`arch/INDEX.md`](./arch/INDEX.md), boundaries con fitness tests
  (`go test ./arch/fitness/...`). Claves: `git-solo-lectura-y-commit` (commit SOLO por
  pathspec; push/pull/fetch no existen; v1.1 confina `adapters/registry/gitsync` a
  `~/.dev-studio/registry/`) · `sesion-aislada-por-cwd` (enforced).
- **UI:** design system **PRENTER** (tokens en `specs/shell/tokens/`; dark-first, teal único
  acento, Jost/Mulish/JetBrains Mono vendorizadas) + Storybook (RN-9: componente sin story no
  entra). Shell: rail Repositorios→Workspaces (workspace = sesión 1:1) + studio-nav +
  overlays (jamás tapan rails) + panel **Cambios git real** (status/diff/side-by-side/log).
- **Roles (registry de arneses):** instalación por proyecto = lock as-code
  `.devstudio/arneses.yaml` committeado por pathspec (contrato estable: `registry` +
  `arneses[].{id,version,canal}`, evolución solo aditiva — superficie de auditoría que
  ArnesIA lee) + caché forma-plugin INTACTA en `~/.dev-studio/arneses/` (rehidratable, modelo
  npm). Sesión con rol → spawn `--plugin-dir` + UN solo `--append-system-prompt` (preámbulo
  rol×proceso + banda Base embebida — los dos flags de prompt de la CLI son excluyentes).
  Campo para pintar: `arnes.l0.nombre` → `plugin.json name` → `id` (ídem descripcion).
- **Dogfooding:** la app se instala (`scripts/install.sh` → `~/.local/bin/dev-studio` +
  `.desktop`) y se actualiza desde su Config overlay (rebuild del repo local + restart). Las
  rebanadas se verifican contra la app INSTALADA.
- **Specs:** permanentes por rebanada en `specs/*/SPEC.md` (se CONGELAN y evolucionan por
  changelog). El port desde el monorepo lo gobierna la épica
  [`epicas/experiencia-orquestada/`](./epicas/experiencia-orquestada/) (carpeta temporal por
  diseño; norte firmado en `NORTE-FIRMADO.md`).

**Git:** trunk-based — `main` única, commit/push directo, tags semver cuando haya releases.
