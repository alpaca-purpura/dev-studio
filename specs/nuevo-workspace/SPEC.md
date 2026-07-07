# SPEC — Flujo «Nuevo Workspace» v2: ubicación → propósito + taxonomía estándar (PB-27)

> Estado: **CONGELADA 2026-07-07** — los 3 forks ratificados por Chris en sesión (ficha DH-17).
> Origen: primer dogfooding (DH-16.1) — Chris: «al crear nuevo workspace, primero preguntar si
> es un nuevo worktree o el actual; luego el usuario debe poder navegar y comprender (sin
> permisos de edición) o trabajar sobre una historia/bugfix/etc existente o crear uno».
> Evoluciona (no rompe): spec shell §4.1 (el picker ya no es la puerta única) y spec
> workspace-aislado §2.1 (branch por tipo) — changelogs anotados allá.

## 1. Taxonomía de paquetes de trabajo (estándar de industria, ratificada)

Convergencia de Jira/Azure (issue types) + Gitflow (branch naming) + Conventional Commits:

| Tipo | Definición | Branch | Commit | Origen estándar |
|---|---|---|---|---|
| `historia` | Funcionalidad nueva — SIEMPRE incrementa una capability (regla de la casa) | `feature/{slug}` | `feat:` | Story (Jira/Scrum) |
| `bug` | Defecto en algo ya entregado | `bugfix/{slug}` | `fix:` | Bug |
| `hotfix` | Corrección urgente sobre producción | `hotfix/{slug}` | `fix:` | Bug crítico (Gitflow) |
| `tarea` | Trabajo técnico sin cara de usuario (refactor/deuda/tooling) | `chore/{slug}` | `chore:`/`refactor:` | Task |
| `spike` | Investigación — el resultado es conocimiento, no merge | `spike/{slug}` | — | Spike (Scrum) |

- Epic = contenedor, jamás una sesión. Release = evento, no trabajo.
- **Muere el prefijo `wt/{slug}`**: la branch del workspace nace con el prefijo de su tipo.
  Workspaces existentes con `wt/…` no se migran (legacy visible, sin drama).
- Exploración EN worktree (comparar una branch sin tocar tu checkout): branch `explore/{slug}`.

## 2. El wizard (modal 2 pasos — reemplaza el salto directo al Backlog)

**Paso 1 — ¿Dónde?**
- ⎇ **Workspace nuevo** — worktree aislado. Habilita los 3 propósitos.
- 📂 **Checkout actual** — raíz del repo, sin aislamiento. Habilita SOLO «Navegar y
  comprender» (las otras opciones visibles pero deshabilitadas con la razón: «trabajo con
  edición siempre en workspace aislado»). El backend lo ENFORCEA (422), no solo la UI.

**Paso 2 — ¿Para qué?**
- 🔍 **Navegar y comprender** → sesión de EXPLORACIÓN: read-only vía `--permission-mode plan`
  (mecanismo nativo de la CLI: lee/analiza/responde, no edita). Sin ítem del backlog. Nombre
  editable (default `explorar`). Chip 🔍 Exploración en header y rail.
- 📋 **Trabajar un ítem existente** → picker del Backlog (flujo de hoy) → Config → crear.
  Las tarjetas muestran su TIPO.
- ➕ **Crear un ítem nuevo** → mini-form: tipo (los 5) + título → Config con el paquete →
  crear. v1: el ítem viaja y persiste CON la sesión (`Session.Historia{ID,Titulo,Tipo}`,
  en state.json); board global de ítems = PB-07 (los creados migran allá).

**RN-1 evolucionada (shell/workspace):** «toda sesión CON PERMISOS DE EDICIÓN liga a un
paquete de trabajo; una sesión de exploración es read-only y no lo necesita». El espíritu
(nada que mute sin trazabilidad) queda intacto.

## 3. Contratos técnicos

- `domain.Historia` gana `Tipo` (`historia|bug|hotfix|tarea|spike`). `domain.Session` gana
  `Modo` (`trabajo|exploracion`; vacío = trabajo legacy).
- `ports.SpawnOpts` gana `ReadOnly bool` → el conductor agrega `--permission-mode plan`.
  (Función `buildArgs` pura, con test.)
- `CreateWorktree(ctx, repoRoot, destino, branch)` — el llamador pasa la branch completa
  (`feature/x`, `explore/y`); colisión → sufijo `-2`… (tests ajustados).
- `POST /api/sessions`: `{repo_id, modo, ubicacion, historia{id,titulo,tipo}, nombre, rol}`.
  Guards del backend: `ubicacion=checkout` + `modo=trabajo` → 422 · `modo=trabajo` sin
  historia → 422 (RN-1) · exploración en checkout → cwd = raíz del repo, SIN worktree,
  read-only (excepción legítima y acotada al boundary sesion-aislada-por-cwd: no muta).

## 4. AC — gate de cierre

- [x] **AC-1** — Wizard: «+ Nuevo Workspace» abre el modal (ya no salta al Backlog); paso 1 →
  paso 2; en «Checkout actual» solo Explorar habilitado (con razón visible).
- [x] **AC-2** — Exploración en checkout: sesión nace en la raíz, chip 🔍, SIN worktree; el
  proceso `claude` corre con `--permission-mode plan`; un turno que pide editar un archivo NO
  lo modifica (verificación real: pedir edición → archivo intacto) y un turno de lectura
  responde.
- [x] **AC-3** — Trabajar ítem existente: picker → sesión con branch `feature/{slug}` (o el
  prefijo del tipo elegido) — `git worktree list` lo confirma.
- [x] **AC-4** — Crear ítem nuevo (tipo bug + título): sesión con branch `bugfix/{slug}`, chip
  del ítem con su tipo, persistido en state.json.
- [x] **AC-5** — Guard: POST con checkout+trabajo → 422; trabajo sin historia → 422.
- [x] **AC-6** — Suite verde (tests worktree ajustados a prefijos + buildArgs + fitness).

## Changelog

- v1 CONGELADA 2026-07-07 — 3 forks ratificados (taxonomía · wizard+RN-1 evolucionada ·
  crear-ítem entra ya).
- v1 VERIFICADA 2026-07-07 — 10/10 checks en vivo contra el binario instalado: wizard con
  guard visual + API (422×2) · exploración con `--permission-mode plan` visto en el proceso
  real, turno de lectura respondido y pedido de edición BLOQUEADO («Plan mode activo — no
  puedo crear archivos», archivo inexistente) · branches spike/… y bugfix/… en worktree list.
  Ficha DH-17.
