# arch/ — arquitectura y diseño técnico as code de DevStudio

> **La arquitectura vive acá, no en un doc que se pudre.** Cada regla estructural (boundary) es
> un nodo con **L1** (principio con fuente) ↔ **L2** (realización en este árbol Go+React) + una
> **tabla de checks ejecutables**. El «por qué» vive en [`../LEDGER.md`](../LEDGER.md) (ficha
> `DH-13`, sin duplicar); la «prueba» vive en [`fitness/`](./fitness/); el «dibujo» en
> [`model/`](./model/). Patrón copiado del proyecto hermano `harness-studio` (mismo operador,
> mismo estilo `arch/`), adaptado al scope real de DevStudio — ver nota de honestidad abajo.
> Norte: [`../VISION.md`](../VISION.md).

## Cómo se relaciona con el resto del repo

- [`../VISION.md`](../VISION.md) — identidad + decisiones técnicas vigentes (driver CLI-nativo,
  binario Go + UI embebida, proceso/arquitectura/documentación as code). Esta capa aterriza las
  decisiones de arquitectura, no las de producto.
- [`../LEDGER.md`](../LEDGER.md) — el diario firmado. Cada boundary cita `ledger: DH-13` (la
  ficha que entregó el esqueleto F1 de la épica «Experiencia Orquestada»).
- [`../epicas/experiencia-orquestada/NORTE-FIRMADO.md`](../epicas/experiencia-orquestada/NORTE-FIRMADO.md)
  — el norte FIRMADO (F0 cerrada 2026-07-07, DH-14) que gobierna QUÉ entra y en qué orden;
  `arch/` documenta CÓMO se construyó lo que ya entró (F1: esqueleto + driver CLI-nativo).

## Las dos capas (en cada boundary node)

1. **L1 · Principio** — el patrón como estándar de industria (hexagonal, ports&adapters,
   event-normalization…), con fuente y fecha de revisión.
2. **L2 · Realización** — cómo se mapea en ESTE árbol Go+React: qué paquete es dominio/puerto/
   adaptador/transporte, y divergencias marcadas explícitamente (nunca silenciosas).

## El árbol — boundary nodes

| Nodo | Regla | Estado | Checks | Enforcer |
|------|-------|--------|--------|----------|
| [`boundaries/adaptador-agente-intercambiable.md`](./boundaries/adaptador-agente-intercambiable.md) | Claude Code = un adaptador tras `AgentPort` | 🌳 enforced | 2 | aserción de compilador |
| [`boundaries/conductor-encapsula-stream-json.md`](./boundaries/conductor-encapsula-stream-json.md) | Un solo paquete parsea el protocolo `stream-json` | 🌱 proposed | 2 | code review (TBD automatizar) |
| [`boundaries/sesion-un-turno-a-la-vez.md`](./boundaries/sesion-un-turno-a-la-vez.md) | Un turno concurrente por sesión, rechazo explícito (no cola silenciosa) | 🌳 enforced | 2 | `fitness/arch_test.go:TestOneTurnAtATime` |
| [`boundaries/dominio-independiente-de-transporte.md`](./boundaries/dominio-independiente-de-transporte.md) | `domain`/`usecase` no importan `net/http` | 🌳 enforced | 2 | `fitness/arch_test.go:TestDomainNoTransportImport` |
| [`boundaries/sesion-aislada-por-cwd.md`](./boundaries/sesion-aislada-por-cwd.md) | Cada sesión confinada a su cwd — **parcial, brecha documentada** | 🌱 proposed | 3 (1 ✅, 2 ❌ TBD) | ninguno todavía (validación de path pendiente) |
| [`boundaries/git-solo-lectura-y-commit.md`](./boundaries/git-solo-lectura-y-commit.md) | La app nunca push/pull/fetch/reset/rebase; commit solo por pathspec explícito | 🌳 enforced | 2 | `fitness/arch_test.go:TestGitAdapterSinVerbosProhibidos` + `TestGitCommitExigePathspec` |

Leyenda: 🌱 proposed (declarado, chequeo manual o pendiente) · 🌳 enforced (código + test
corriendo). **Total: 6 boundaries · 13 checks · 3 enforced con test real** (`go test
./arch/fitness/...` pasa hoy).

> **Honestidad (heredada de la casa, vía `harness-studio`/METODOLOGIA):** este árbol es
> deliberadamente MÁS CHICO que el del proyecto hermano — F1 es un esqueleto, no la fábrica
> completa. `sesion-aislada-por-cwd` documenta a propósito una brecha de seguridad real (sin
> validación de rutas protegidas) en vez de omitirla u ocultarla detrás de un `enforced` falso.

## Subdirectorios

- [`model/`](./model/) — diagramas as code. Hoy: `system-context.mmd` (C4 nivel 1, Mermaid,
  renderiza nativo en GitHub). Sin container diagram todavía (F1 es un solo binario, un solo
  proceso — el C4 nivel 2 no aporta hasta que haya más de un componente desplegable).
- [`fitness/`](./fitness/) — `arch_test.go`: los 2 checks que hoy corren en `go test
  ./arch/fitness/...`. Sin `go-arch-lint` todavía (herramienta externa, no instalada — el
  import-graph check se hace hoy con `go list -deps` a mano, ver `dominio-independiente-de-transporte`).
- `contracts/`, `conventions/`, `decisions/` — **no existen todavía**. Se crean cuando F2+ (port
  por rebanadas) los necesite; no se anticipan vacíos (regla: no fabricar estructura sin uso real).

## El stack (lo que F1 realmente usa)

| Capa | Pick | Nota |
|---|---|---|
| Backend | binario Go único `dev-studio` (`cmd/dev-studio`) | hexagonal: `domain` → `usecase` → `ports` ← `adapters` |
| Conexión CC | subproceso `claude` + `stream-json` stdin/stdout | mismo patrón conductor que `harness-studio`, implementación propia |
| Persistencia | JSON plano (`~/.dev-studio/state.json`: repos + sesiones; migra `sessions.json` F1 solo), escritura atómica | sin SQLite todavía — no hace falta índice a este tamaño |
| Git del usuario | subproceso `git` (adapter `adapters/git/cli`) — status/diff/log/commit, nada más | boundary `git-solo-lectura-y-commit` (DH-15) |
| Transporte realtime | SSE (`/events`, un solo tipo de evento `dock`) | sin replay por `Last-Event-ID` todavía (TBD) |
| Frontend | Vite + React 19 SPA, `go:embed` | Zustand (store de sesiones), sin router (una sola página) |
| Estilo | Tailwind v4 + tokens copiados de `harness-studio` (`theme.css`) | dark por `data-theme`, mismo mecanismo |
| UI copiada 1:1 | `SessionRail` (rail colapsable, tabs estilo WARP) | simplificado: sin `arnes`/`salud`/`view` — eso es dominio de ArnesIA, no de DevStudio |

## Cómo crece

Igual que el hermano: la arquitectura se revisa **al cambiar**, no en cadencia fija. Detalle del
ritual en [`CADENCE.md`](./CADENCE.md).
