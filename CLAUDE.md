# DevStudio — construir y mantener software basado en proceso y arquitectura

Producto standalone (graduado del monorepo `prenter-harness`, 2026-07-04). Norte =
[`VISION.md`](./VISION.md) · registro = [`LEDGER.md`](./LEDGER.md) (fichas `DH-NN`, arranca en
DH-12; DH-01..DH-11 viven en la incubadora `prenter-harness/products/devhub/`, congelada).

**Orden de producto (SSoT, 2026-07-07):** todo lo pendiente = [`BACKLOG.md`](./BACKLOG.md)
(fichas `PB-NN`, prioridad por bandas — la firma Chris) · todo lo funcional =
[`INCREMENTO.md`](./INCREMENTO.md) (`CAP-NN` con evidencia real). Regla: idea nueva → PB antes
de trabajarse; entrega verificada → CAP en el mismo commit. Ritual: «¿qué viene?» → leer banda 🔴.
No confundir planos: PRODUCTO (VISION, arch/, código, BACKLOG/INCREMENTO) vs PROYECTO (LEDGER,
epicas/, .claude/, docs/process/, templates/, scripts/, project.config.yaml, research/).

**Qué es:** aplicación de escritorio multiplataforma donde una organización construye y
mantiene software con proceso, arquitectura y documentación **as code**. Usuarios (CTO ·
developer · devops · product owner) trabajan orquestados sobre uno o varios sistemas; el
repositorio GitHub es el conector.

**Decisiones técnicas vigentes:**
- Conexión con Claude Code = **driver CLI-nativo** (DH-10): spawn del `claude` del usuario vía
  stdin/stdout stream-json. BYO licencia — **sin API de Anthropic**; la app no toca credenciales.
- Binario Go + UI embebida (`go:embed`); instalador = el binario solo (Win/Linux/mac).
- Proceso as code = descriptor I-77 (categorías semánticas fijas · gates · dueños · verbos
  CDEvents); la consola deriva TODO del dato, cero hardcode del ciclo.

**Estado:** repo recién fundado — el código vive aún en el monorepo (congelado) y entra por
port gradual gobernado por la épica [`epicas/experiencia-orquestada/`](./epicas/experiencia-orquestada/)
(carpeta temporal por diseño; specs permanentes → `specs/`).

**Arnés de construcción:** kit dev (plugin del marketplace `alpacapurpura/prenter-marketplace`,
pineado por versión). Evoluciona con el producto — mejoras al arnés se upstreamean al kit
(backflow I-59), jamás fork silencioso.

**Estado (DH-13, 2026-07-06):** F1 (esqueleto de la app) entregado y verificado en vivo. Existe
código real: módulo Go `github.com/alpacapurpura/dev-studio` con arquitectura hexagonal
(`internal/domain` → `internal/ports` ← `internal/adapters` ← `internal/usecase`), driver
CLI-nativo propio (`internal/adapters/agent/claudecode`, subproceso `claude -p --input-format
stream-json --output-format stream-json`, un proceso por sesión), API REST + SSE
(`internal/adapters/transport/{http,sse}`), persistencia liviana en `~/.dev-studio/sessions.json`,
y SPA React 19 + Zustand + Tailwind v4 embebida vía `go:embed` (`web/`, binario único
`cmd/dev-studio`). Verificado con 2 sesiones concurrentes reales (procesos `claude`
independientes, sin cruce). Arquitectura as code propia en [`arch/INDEX.md`](./arch/INDEX.md)
— hoy 6 boundaries, 3 enforced con test (`go test ./arch/fitness/...`), 1 brecha de seguridad
documentada sin ocultar
([`arch/boundaries/sesion-aislada-por-cwd.md`](./arch/boundaries/sesion-aislada-por-cwd.md) —
**CERRADA en DH-16**).

**Estado (DH-14, 2026-07-07):** F0 cerrada — norte FIRMADO en
[`epicas/experiencia-orquestada/NORTE-FIRMADO.md`](./epicas/experiencia-orquestada/NORTE-FIRMADO.md):
journeys ×4 (CTO = hueco declarado, bloqueado por PB-21→22) · 7 principios · **reframe Rol: rol
= arnés instalado desde REGISTRY PROPIO de DevStudio** (solo marketplace, cero roles locales;
el descriptor I-77 y el rol entero viajan EN el arnés; permisos por puesto = con multi-usuario)
· orden de port Proyecto→Rol→Sesión-workspace→Historia/Capability→Proceso · multi-usuario TBD
formal sin bloquear (PB-21, dueño Chris). Seam `project.config.yaml` completo (los 4 slots
firmados: +registry, roster real del repo, 10 estados, wip_caps advisory).

**Estado (DH-15, 2026-07-07):** shell PRENTER ENTREGADO (PB-04 ⊕ PB-03, spec CONGELADA en
[`specs/shell/SPEC.md`](./specs/shell/SPEC.md), 24/24 checks del gate en vivo). Design system
= **PRENTER** (SSoT Claude Design `a98c2e0d`, tokens sincronizados en `specs/shell/tokens/`;
teal único acento, dark-first, Jost/Mulish/JetBrains Mono vendorizadas) + Storybook (RN-9:
componente sin story no entra). El rail F1 murió: rail Repositorios→Workspaces (workspace =
sesión 1:1, `web/src/widgets/repos-rail`) + studio-nav + overlays (jamás tapan rails) + panel
**Cambios git real** (status/diff/side-by-side/log + commit SOLO por pathspec; push/pull/fetch
no existen — boundary `git-solo-lectura-y-commit` enforced). Sesión nace SIEMPRE ligada a
historia (picker, RN-1; historias mock hasta PB-07). Persistencia `~/.dev-studio/state.json`
(migra `sessions.json` F1 solo). CAP-07/08/09 nuevas en INCREMENTO.

**Estado (DH-16, 2026-07-07):** PB-02 ⊕ PB-26 ENTREGADOS (spec `specs/workspace-aislado/SPEC.md`,
14/14 checks contra el binario INSTALADO). Toda sesión de repo nace en su **worktree propio**
(`~/.dev-studio/workspaces/{repo}/{slug}`, branch `wt/{slug}`) + rutas protegidas vedadas en
toda puerta + cierre con modal conservar/borrar (jamás `--force`) → **boundary
`sesion-aislada-por-cwd` ENFORCED** (resta token de capacidad de la API local, check TBD).
**Dogfooding:** la app se instala (`scripts/install.sh` → `~/.local/bin/dev-studio` +
lanzador `.desktop` + icono PRENTER) y se ACTUALIZA desde su footer («Actualizar» = rebuild
del repo local + restart; en el gate se actualizó a sí misma). Las rebanadas se prueban desde
la app instalada. CAP-10/11 nuevas. **Siguiente:** banda 🟡 — PB-05 restos → PB-25 (registry)
→ PB-06 (Roles).

**Git:** trunk-based — `main` única, commit/push directo, tags semver cuando haya releases.
