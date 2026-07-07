# SPEC — Workspace aislado por sesión + instalable dogfooding (PB-02 ⊕ PB-26)

> Estado: **CONGELADA 2026-07-07** — ratificada por Chris (ficha DH-16). Cambios de alcance =
> nueva ratificación + changelog.
> Segunda spec permanente. Cierra la brecha de seguridad `sesion-aislada-por-cwd` como FEATURE
> (patrón Conductor, forma firmada en F0/DH-14: worktree + validación de rutas en la MISMA
> entrega) y convierte la app en el lugar donde Chris la prueba de acá en adelante
> (instalable local + updater rebuild).
>
> **Fuentes:** `arch/boundaries/sesion-aislada-por-cwd.md` (la brecha) · spec clon §5.1 ·
> NORTE-FIRMADO §Dudas (DH-14) · spec shell `specs/shell/SPEC.md` (el shell que la recibe).
> **Forks firmados 2026-07-07:** worktrees en `~/.dev-studio/workspaces/` (fuera del repo del
> usuario) · al cerrar sesión se PREGUNTA (conservar/borrar) · updater = rebuild local (sin
> red, sin tokens, boundary git intacto).

## 1. Qué se construye

**PB-02:** toda sesión nueva nace en su PROPIO workspace: `git worktree add` + branch
`wt/{slug}` — cero cwd compartido entre sesiones, cero sesión en ruta peligrosa. La brecha
`sesion-aislada-por-cwd` pasa de `proposed` a **`enforced`**.

**PB-26:** DevStudio se instala como app de escritorio (binario + lanzador) y se actualiza
DESDE la app (botón «Actualizar» → rebuild del repo local → restart). Dogfooding: cada
rebanada siguiente se prueba en la app instalada, no en `go run`. (PB-24 sigue siendo el
instalable COMERCIAL: cross-compile + firma + ⚠ Consumer Terms — esto es la versión casa.)

## 2. PB-02 — Workspace aislado

### 2.1 Creación (extiende el flujo picker del shell, §4.1 de la spec shell)

1. `Crear sesión aislada` → el backend crea el worktree ANTES de registrar la sesión:
   `git worktree add ~/.dev-studio/workspaces/{repoNombre}/{slug} -b wt/{slug}` (desde la
   raíz del repo, base = HEAD).
2. `Session.Cwd` = la ruta del worktree · `Session.Workspace` = la misma ruta (campo nuevo:
   distingue sesión aislada de legacy) — el subproceso `claude` y el panel Cambios operan ahí.
3. Colisión de branch/carpeta (`wt/{slug}` ya existe) → sufijo `-2`, `-3`… automático.
4. Repo sin commits (HEAD unborn) → error honesto: «el repositorio necesita al menos un
   commit para crear workspaces».
5. El rail muestra por workspace-item la branch REAL de su sesión (`wt/{slug}`) y el diff ±N
   DE SU worktree — ya no el status agregado del repo. Dot `conflict` se activa si el status
   porcelain reporta conflictos (código `U`).

### 2.2 Cierre (fork firmado: PREGUNTAR)

- Cerrar sesión con workspace propio → modal: **«¿Qué hacemos con el workspace?»**
  - **Conservar** (default): la sesión muere, worktree + branch quedan en disco (el trabajo
    committeado o no sigue ahí; merge/limpieza = PB-14).
  - **Borrar**: `git worktree remove` — si el working tree tiene cambios sin commit, git lo
    rechaza y la app conserva + avisa (NUNCA `--force`: cero pérdida silenciosa).
- Sesiones legacy (sin workspace propio) cierran como hoy, sin modal.

### 2.3 Validación de rutas protegidas (la otra mitad de la brecha)

- **Al registrar repo** y **al crear sesión**: la ruta se rechaza si es protegida:
  `$HOME` exacto · `~/.ssh` · `~/.gnupg` · `~/.aws` · `~/.kube` · `~/.docker` ·
  `~/.dev-studio` (la app no opera sobre sí misma) · raíces de sistema (`/`, `/etc`, `/usr`,
  `/bin`, `/sbin`, `/var`, `/boot`, `/root`, `/proc`, `/sys`, `/dev`) y cualquier ruta DENTRO
  de ellas que no sea home del usuario.
- La regla vive en el dominio (`domain.RutaProtegida(home, ruta)`, pura, testeable); los
  usecases la aplican. Mensaje de error honesto con la ruta y el porqué.
- **Boundary `sesion-aislada-por-cwd` → `enforced`**: fitness tests (registrar `$HOME` →
  rechazo · crear sesión con cwd `/etc` → rechazo · sesión nueva vía repo → su cwd es el
  worktree, distinto de la raíz y de otras sesiones).

### 2.4 Puertos / arquitectura

- `ports.GitWorkspace` (nuevo): `CreateWorktree(ctx, repoRoot, slug) (path, branch, error)` ·
  `RemoveWorktree(ctx, repoRoot, path) error`. Implementado en `adapters/git/cli` con los
  verbos `worktree add` / `worktree remove` / `branch` — los verbos PROHIBIDOS del boundary
  `git-solo-lectura-y-commit` (push/pull/fetch/reset/rebase) siguen sin existir; el fitness
  scanner sigue pasando.
- `DELETE /api/sessions/{id}?workspace=keep|remove` — default `keep`.
- `POST /api/sessions` con `repo_id` ahora crea el worktree (v1 del shell usaba la raíz del
  repo — esa línea muere).

## 3. PB-26 — Instalable dogfooding + updater rebuild-local

### 3.1 Instalador (`scripts/install.sh`)

- Build completo (npm build + `go build` con `-ldflags "-X main.version={git describe/SHA}
  -X main.buildDate={fecha}"`) → binario a `~/.local/bin/dev-studio`.
- Lanzador de escritorio `~/.local/share/applications/dev-studio.desktop` + icono PRENTER
  (SVG teal en `assets/`): abre la app en modo ventana (`google-chrome --app=http://…` si
  existe; si no `xdg-open`). El lanzador levanta el binario si no está corriendo.
- Escribe `~/.dev-studio/app.json` → `{ "source": "<ruta del repo dev-studio>" }` — de ahí
  sabe el updater desde dónde rebuildear.
- Idempotente: correrlo de nuevo = actualizar a mano.

### 3.2 Updater en la app

- `GET /api/version` → `{version, build_date, source}` — la SPA muestra la versión en el
  footer del rail (reemplaza un stub del footer).
- Botón **«Actualizar»** (footer del rail) → `POST /api/update`:
  1. Rebuild desde `app.json.source` (npm build + go build con versión nueva) → binario
     nuevo sobre `~/.local/bin/dev-studio`.
  2. Responde `202` y el proceso se REINICIA (`syscall.Exec` del binario nuevo, mismo addr).
  3. La SPA espera polleando `/api/version` hasta ver la versión cambiar → `location.reload()`.
- Falla el build → `500` con el output del compilador visible en la UI (honestidad: el error
  real, no «algo salió mal»).
- **Sin git**: el updater compila el working tree del repo local tal como está — pull/push
  los sigue haciendo Chris; el boundary git-solo-lectura-y-commit no se toca.

## 4. Reglas de negocio

- **RN-1** — Sesión nueva vía repo = worktree propio SIEMPRE; el cwd de una sesión jamás es
  la raíz compartida del repo (las legacy existentes quedan como están, marcadas sin
  aislamiento).
- **RN-2** — Ruta protegida = rechazo en TODA puerta (registrar repo, crear sesión), con
  mensaje que nombra la ruta.
- **RN-3** — Cerrar jamás borra trabajo en silencio: borrar workspace requiere elección
  explícita, y git rechaza el remove sucio (sin `--force`).
- **RN-4** — El updater no ejecuta ningún verbo git QUE MUTE estado (pull/fetch/push/reset/
  rebase/checkout): compila lo que hay en `source` tal como está. Lecturas puras de versión
  (`rev-parse`, `diff --quiet`) en el script de build sí — de ahí sale el SHA visible.
- **RN-5** — La versión visible en la app es la del binario corriendo (SHA + fecha) — tras
  actualizar, se ve cambiar: esa es la prueba.

## 5. AC — gate de cierre (verificación viva)

- [x] **AC-1** — Crear sesión vía picker: existe `~/.dev-studio/workspaces/{repo}/{slug}`,
  `git worktree list` del repo lo muestra con branch `wt/{slug}`, el header/rail de la sesión
  muestran ese cwd/branch.
- [x] **AC-2** — Dos sesiones del mismo repo: worktrees distintos; un archivo tocado en el
  worktree A NO aparece en el panel Cambios de B; commit desde A no ensucia B.
- [x] **AC-3** — Registrar `$HOME` o `/etc` como repo → rechazo visible; crear sesión con cwd
  protegido vía API → 422; fitness `sesion-aislada-por-cwd` en verde (boundary `enforced`).
- [x] **AC-4** — Cerrar sesión con workspace: modal conservar/borrar; conservar deja el
  worktree; borrar con cambios sin commit → git lo rechaza y la app conserva + avisa; borrar
  limpio → desaparece de `git worktree list`.
- [x] **AC-5** — `scripts/install.sh` instala: binario en `~/.local/bin`, `.desktop` presente,
  la app abre desde el lanzador y muestra su versión (SHA) en el footer.
- [x] **AC-6** — «Actualizar» desde la app instalada: rebuild + restart + la SPA recarga sola
  mostrando la versión NUEVA (evidencia RN-5). El botón muestra el error real si el build
  falla.
- [x] **AC-7** — Suite completa verde: `go test ./...` + fitness (los 2 boundaries git) +
  `tsc` + builds.

## 6. Fuera de alcance (con destino)

| Qué | Destino |
|---|---|
| Merge/PR del branch wt/{slug} desde la UI | PB-14 |
| Dot `archived` + listado de worktrees huérfanos | idea nueva si molesta |
| Setup script por proyecto al crear workspace (install/migraciones/.env) | PB-13 |
| Updater con descarga remota / releases firmadas | PB-24 |
| Multiplataforma del instalador (hoy: Linux) | PB-24 |

## Changelog

- v1 2026-07-07 — borrador post-forks (ubicación WT / cierre pregunta / updater rebuild).
- v1 CONGELADA + VERIFICADA 2026-07-07 — AC-1..AC-7 tildadas: 14/14 checks en vivo contra el
  binario INSTALADO (incl. self-update con fix vivo post-restart). Evidencia: ficha DH-16.
- EVOLUCIÓN 2026-07-07 (PB-27, DH-17): §2.1 — la branch del worktree ya no es `wt/{slug}`:
  lleva el prefijo estándar de su tipo (`feature/ bugfix/ hotfix/ chore/ spike/ explore/`,
  `specs/nuevo-workspace/SPEC.md` §1). `CreateWorktree` recibe la branch completa.
