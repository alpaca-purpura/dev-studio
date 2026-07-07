# SPEC — Shell PRENTER (rebanada PB-04 ⊕ PB-03) — v1

> **Primera spec permanente del repo** (`specs/`, patrón «SPEC congelada por fase»). Estado:
> **CONGELADA 2026-07-07** — ratificada por Chris (ficha DH-15). Cambios de alcance = nueva
> ratificación + changelog.
> Construye el shell nuevo de la app componiéndolo de átomos catalogados (PB-03 y PB-04 se
> entregan JUNTAS — fork firmado 2026-07-07). La rebanada siguiente (PB-02, workspace aislado)
> aterriza YA sobre este shell.
>
> **Fuentes (por precedencia):**
> 1. `epicas/experiencia-orquestada/mockup-clon-agentsroom.html` — SSoT de la FORMA (se
>    promueve a `specs/shell/` al cierre de la épica).
> 2. `CLON-AGENTSROOM-SPEC.md` §0.5 (shell ratificado) + §10 (composer/cambios, click-by-click).
> 3. `NORTE-FIRMADO.md` (DH-14) — journeys, principios, reframe rol=arnés-de-registry.
> 4. Design system **PRENTER** — Claude Design `a98c2e0d-db82-43f2-8fd7-e7e05c40fd51`
>    (`tokens/*.css`, sincronizados en `specs/shell/tokens/`) + `~/Proyectos/prenter/marketing/brand-guidelines.md`.
>
> **Forks firmados 2026-07-07 (sesión shell):** paleta = PRENTER teal (mockup ya la usa — los
> swatches son los tokens de marca) · roster Config = placeholder honesto del registry ·
> alcance = ambicioso (incluye panel Cambios REAL) · Storybook nace DENTRO de esta rebanada.

## 1. Qué se construye

El shell definitivo de DevStudio reemplaza la UI de F1 (rail simple de sesiones + área en
blanco): **rail Repositorios→Workspaces + studio-nav + vista de sesión (chat/composer) + panel
lateral con Cambios git REALES + overlays Backlog/Configuración/Producto/Roles**, todo
compuesto de átomos PRENTER catalogados en Storybook. El rail de sesiones F1
(`web/src/widgets/session-rail`) **muere**: workspace = sesión (1:1) y el árbol de la izquierda
es la única lista.

**No-objetivos de esta rebanada** (ver §9): worktree/branch real por sesión (PB-02, la
inmediata siguiente), terminal PTY, board de backlog real, registry de arneses, push/pull,
diff conversacional, multi-proveedor.

## 2. Design system PRENTER (PB-03 dentro de esta rebanada)

- **SSoT visual** = proyecto Claude Design «PRENTER Design System» (`a98c2e0d…`). Copia
  sincronizada y versionada: `specs/shell/tokens/{colors,fonts,spacing,typography}.css`
  (2026-07-07, etags anotados en cada archivo). Cambia allá → se re-sincroniza acá.
- **`web/src/app/styles/theme.css` se reescribe**: mueren los tokens ámbar de harness-studio
  (`--primary #a8742c`); entran los PRENTER mapeados a las MISMAS variables semánticas
  (convención shadcn/Tailwind v4 que ya usa la SPA). El mapping canónico es el `:root` del
  mockup (líneas 7-57): dark default + `data-theme="dark"|"light"` explícitos.

| Variable semántica | Dark (default) | Light | Token PRENTER |
|---|---|---|---|
| `--background` / `--foreground` | `#08090a` / `#fff` | `#f6f8f8` / `#0e1312` | lienzo dark · neutral-50/900 |
| `--card` / `--popover` | `#0c1110` | `#ffffff` | `--dark-surface` · `--surface-card` |
| `--primary` / `--primary-foreground` | `#1fc6b8` / `#04211f` | `#009d92` / `#fff` | teal-400 (sobre dark) · teal-600 (sobre claro) |
| `--secondary` / `--muted` / `--accent` | `#121817` | `#eef1f1` | dark-raised aprox · neutral-100 |
| `--border` / `--input` / `--ring` | `#1f2826` / `#2a3432` / `#1fc6b8` | `#dde3e2` / `#c7cfce` / `#009d92` | `--dark-border` · teal |
| `--ok` / `--warn` / `--crit` (+`-soft`) | `#1f9d6b` / `#c98a16` / `#d0483c` | ídem | `--success/--warning/--danger` |
| `--sidebar*` | `#0c1110` base | `#ffffff` base | superficie + borde dark |
| `--shadow-glow` | `0 8px 24px rgba(0,183,170,.22)` | ídem | `--shadow-brand` |
| radios / spacing | sm4·md8·lg14·xl22·full | ídem | escala PRENTER 1:1 |
| tipografías | Jost (display) · Mulish (sans) · JetBrains Mono | ídem | sustitutos vigentes de Coco Gothic/Sansation |
| escala tipo (app) | xs11 · sm12.5 · base14 · lg16 · xl20 · 2xl26 | ídem | override compacto de aplicación (la escala 64-12 de PRENTER es web/marketing) |

- **Fuentes vendorizadas**: los `.woff2` de Jost/Mulish/JetBrains Mono se descargan al build y
  viajan en el binario (`go:embed`) — la app de escritorio NO depende de Google Fonts en
  runtime. Cuando lleguen los binarios de Coco Gothic/Sansation, se reemplazan.
- **Reglas de marca que el shell respeta**: teal = ÚNICO acento saturado (jamás un segundo) ·
  dark-first, light = respiro · bordes hairline, jerarquía por escala/color de texto, no por
  cajas · mono uppercase tracking ancho para etiquetas técnicas (branch, FIG, atajos) · sin
  emojis en copy de marca (los emoji del mockup en chips/botones son iconografía funcional de
  producto, aceptada) · glow teal solo en CTA primaria.
- **Storybook (PB-03)**: nace en `web/` con TODOS los átomos/moléculas que el shell compone —
  catálogo inicial mínimo: `Button` (primary/outline/ghost) · `IconButton` · `Chip` · `Pill`
  (ok/warn/crit/muted/soft) · `StatusDot` (idle/ready/conflict/archived) · `Avatar` (iniciales
  rol) · `Card` · `Input`/`Textarea`/`Select` · `Toggle` · `Kbd` · `Tab` · `FilterChip` ·
  `DiffStat` (+N −M) · `ReleaseChip` · `StatePill` (10 estados story) · organismos: `RepoBlock`,
  `WorkspaceItem`, `StudioNav`, `Composer`, `PanelOverlay`, `ModalShell`. **Regla dura: el
  shell no maqueta con divs sueltos — si falta un átomo, se crea con su story primero.**
  Dark y light por historia (toolbar de tema).

## 3. Estructura del shell (mapa REAL / STUB por zona)

Leyenda: **REAL** = funcional contra backend en esta rebanada · **STUB** = presente, visual,
con honestidad de superficie (principio 7: banner/tooltip «todavía no») · **N/A** = no se
construye ni se muestra.

### 3.1 Rail Repositorios→Workspaces (`repos-rail`, 224px, colapsable no — fijo en v1)

| Elemento | Estado | Comportamiento v1 |
|---|---|---|
| Botón «Inicio» | STUB | placeholder sin destino (como el mockup lo declara) |
| Registrar repositorio (input ruta) | **REAL** | ruta local absoluta → valida que exista y contenga `.git` → persiste. URL remota: v1 rechaza con mensaje honesto («clonar desde URL llega después») |
| Lista repos + contador + colapso por repo | **REAL** | repos registrados persistidos; chevron colapsa |
| `+ Nuevo Workspace` | **REAL** | dispara el flujo picker (§4.1) |
| Dropdown ⋯ (duplicar/desde branch/worktree) | STUB | menú visible, opciones deshabilitadas — «se define con PB-02» |
| Workspace-item: nombre + branch + status-dot + diff ±N | **REAL** parcial | workspaces = sesiones vivas de la API. v1: branch = branch actual del cwd (real, leída por git), diff ±N real (`git diff --shortstat`), status-dot: `idle` (sin cambios) / `ready` (cambios sin conflicto). `conflict`/`archived` quedan para PB-02+ |
| Atajos ⌘1..⌘9 | **REAL** | seleccionan workspace n |
| Sincronía workspace↔sesión activa | **REAL** | 1:1 — seleccionar workspace = activar sesión (mismo `data-session`) |
| Footer: Archivados / Feedback / Ajustes | STUB | tooltips «todavía no» |

### 3.2 Studio-nav (iconos verticales)

| Ítem | Estado | Nota |
|---|---|---|
| Studio | **REAL** | vista por defecto; cierra overlays |
| Backlog | **REAL** (abre overlay) | contenido del overlay: STUB con datos de ejemplo (§3.4) |
| Producto | **REAL** (abre overlay) | contenido STUB (§3.6) · microcopy en español: «Producto» (el mockup dice "Product" — se corrige, locale `es`) |
| Roles | STUB | badge `pronto` — llega con PB-25/PB-06 |
| Config | **REAL** (abre overlay) | contenido: placeholder honesto del registry (§3.5) |

Regla §0.5 intacta: los overlays tapan chat+panel lateral; el rail de repos y el studio-nav
quedan SIEMPRE visibles y clickeables. Esc cierra y vuelve a Studio.

### 3.3 Vista de sesión (`session-main` + `side-panel`)

| Elemento | Estado | Comportamiento v1 |
|---|---|---|
| Header: avatar rol · nombre workspace · cwd · proveedor | **REAL** parcial | datos de la sesión real; rol = el elegido en el picker (dato local, sin arnés aún); proveedor fijo «Claude Code» |
| Chip historia ligada (📋 …) | **REAL** parcial | muestra la historia mock elegida en el picker; persiste con la sesión |
| Chat (turnos, streaming) | **REAL** | driver CLI-nativo existente (CAP-02) — sin cambios de protocolo |
| Composer: textarea + Enviar | **REAL** | Enter envía; deshabilitado sin texto; Shift+Enter = nueva línea |
| Composer: Encolar (⇒) | **REAL** | cola de mensajes client-side: si hay turno en curso, encola y despacha al terminar (respeta `TestOneTurnAtATime`) |
| Composer: redimensionar/expandir | **REAL** | drag del borde + botón ↗ |
| 7 botones (boceto/captura/html/móvil/prompts/voz/dictar) | STUB | modales del mockup tal cual (boceto con canvas dibujable; adjuntar deshabilitado con nota) |
| Fila terminal: tab Terminal, +, Comandos, SSH, layout | STUB | terminal PTY = PB-09; modales Comandos/SSH como el mockup |
| Side-panel tab Archivos | **REAL** | árbol de archivos TOCADOS (derivado del status git), read-only |
| Side-panel tab **Cambios** | **REAL** | ver §3.7 — el corazón ambicioso de la rebanada |
| Side-panel tab Pruebas | STUB | mensaje honesto (como el mockup) |

### 3.4 Overlay Backlog (STUB con regla real)

Board con los **10 estados firmados** (8 lanes visibles del stream + toggle «mostrar pausadas
+ descartadas» que revela `parked`/`dropped` — armoniza el «8 estados» del mockup con el
value_stream de 10 firmado en F0). Datos: historias de EJEMPLO hardcodeadas + banner honesto
«datos de ejemplo — el board real llega con el port de Historia/Capability (PB-07)». WIP
chips advisory (rojo al exceder, jamás bloquea — firmado F0). **El modo picker SÍ es real**
(§4.1): elegir tarjeta en modo picker liga la historia a la sesión nueva.

### 3.5 Overlay Configuración (placeholder honesto del registry — fork firmado)

- Banner superior: **«Los roles vienen del registry de arneses de tu organización — registry
  propio en construcción (PB-25). Este roster es una vista previa estática.»**
- Roster estático (builder/auditor/humano-complementario como el mockup) + detalle de rol con
  prompt de ejemplo. Selector Proveedor/Modelo/Cuenta: visible, fijo en «Claude Code /
  Personal (BYO CLI)»; Modelo v1 sin efecto (STUB declarado).
- **Se ELIMINA «+ Rol personalizado»** (contradice DH-14: cero roles locales).
- `Crear sesión aislada` = **REAL**: crea la sesión ligada al paquete elegido (§4.1). Nota de
  nombre: v1 crea la sesión con cwd = raíz del repo (el AISLAMIENTO real llega con PB-02); el
  botón conserva el nombre porque el flujo y el contrato ya son los definitivos.

### 3.6 Overlay Producto (STUB)

Mapa del mockup (salud + zonas Estudio/Plataforma/Infraestructura + lens) con datos estáticos
que REFLEJAN el INCREMENTO real al momento del build (CAP-01..06 → cajas/áreas live/planned) +
banner «mapa estático — se deriva del dato real cuando el SYSTEM-MAP as code entre (PB-08+)».

### 3.7 Panel Cambios — REAL (alcance ambicioso firmado)

| Sub-elemento | Estado | Comportamiento v1 |
|---|---|---|
| Header branch + contador cambios | **REAL** | branch actual + `git status` del cwd de la sesión |
| Lista archivos con checkboxes | **REAL** | staged+unstaged+untracked; checkbox = selección para commit |
| Chips filtro por rol | STUB v1 | un solo chip «Todos» (atribución por rol exige PB-02: worktree por sesión; se declara) |
| Quick-look (clic en archivo) | **REAL** | diff unificado del archivo (modal), syntax plano |
| «Revisar cambios» side-by-side | **REAL** | ORIGINAL/MODIFICADO por archivo, navegación Anterior/Siguiente, contador N/M, pantalla final «listo para commit». «✨ Generar con IA» v1 STUB (deshabilitado con nota) |
| Menú ⋯ por archivo (ver/historial/carpeta) | STUB | opciones deshabilitadas |
| Mensaje de commit + **Confirmar cambios** | **REAL** | commit SOLO de los archivos tildados, por pathspec explícito (`git add <paths> && git commit`) — jamás `add .` |
| «Adjuntar la conversación del rol al commit» | STUB v1 | checkbox visible deshabilitado — «llega con PB-10» |
| «Confirmar y cerrar sesión» | **REAL** | commit + cierra la sesión (flujo existente de cierre) |
| Botones fetch / pull / push | STUB | deshabilitados con tooltip — RN-4 (§7) |
| Sub-tab Historial | **REAL** | `git log` real del cwd (mensaje, autor, sha corto, fecha); filtros de búsqueda client-side |

## 4. Flujos (cableado real)

### 4.1 Crear sesión (la regla dura §0.5, ahora funcional)

1. `+ Nuevo Workspace` (en un repo) → abre **Backlog en modo picker** (banner).
2. Clic en una historia → cierra Backlog → abre **Configuración** con chip «📋 Nueva sesión
   para: {historia}».
3. Elegir rol del roster (estático) → `Crear sesión aislada` → `POST /api/sessions` con
   `{repoId, historia:{id,titulo}, rol}` → sesión nueva aparece como workspace-item del repo,
   seleccionada, chat listo.
4. Config directo (⚙) sin paquete → hint + `Crear sesión aislada` redirige al picker. **No
   existe sesión sin historia.**
5. Esc en cualquier overlay → Studio.

### 4.2 Navegar / cerrar

- Seleccionar workspace-item = activar sesión (chat + panel Cambios apuntan a su cwd).
- Cerrar sesión (desde «Confirmar y cerrar sesión» o control del item) = flujo F1 existente
  (mata el proceso `claude`, remueve el item).
- ⌘1..⌘9 saltan entre workspaces del repo activo.

## 5. Arquitectura (siguiendo `arch/INDEX.md` — hexagonal estricta)

**Dominio** (`internal/domain`):
- `Repo{ID, Nombre, Ruta}` — nuevo agregado: repositorio registrado.
- `Session` gana `RepoID`, `Historia{ID, Titulo}` (valor liviano, mock-friendly), `Rol string`.

**Ports** (`internal/ports`):
- `RepoRegistry` — registrar/listar/eliminar repos (persistencia).
- `GitInfo` — `Status(cwd)`, `Diff(cwd, path)`, `DiffPair(cwd, path)` (original/modificado),
  `Log(cwd, n)`, `Branch(cwd)`, `ShortStat(cwd)`.
- `GitCommit` — `Commit(cwd, paths []string, mensaje string)` — pathspec explícito obligatorio
  en la firma (imposible «add .» por diseño).

**Adapters**:
- `internal/adapters/git/cli` — ejecuta el `git` del usuario (mismo espíritu BYO del driver
  DH-10: subproceso, sin librerías de red). Solo los verbos de los ports — **ni push, ni
  pull, ni fetch existen en el adapter** (enforcement por omisión, RN-4).
- `internal/adapters/persist` — extiende `sessions.json` → `~/.dev-studio/state.json`
  (repos + sesiones, escritura atómica como hoy).
- `internal/adapters/transport/http` — nuevos endpoints:
  `POST/GET/DELETE /api/repos` · `GET /api/repos/{id}/git/status|log` ·
  `GET /api/sessions/{id}/git/{status|diff|log}` · `POST /api/sessions/{id}/git/commit`
  (body: `{paths[], mensaje}`) · `POST /api/sessions` extendido (repoId/historia/rol).
- SSE existente (`dock`) sin cambios; refresh de status git = on-demand (al abrir el panel y
  tras cada turno del agente), sin watcher en v1.

**Frontend** (`web/src`, FSD como F1): `app/styles/theme.css` (tokens nuevos) ·
`widgets/repos-rail` · `widgets/studio-nav` · `widgets/session-view` (chat+composer) ·
`widgets/changes-panel` · `widgets/overlays/{backlog,config,producto,roles}` ·
`shared/ui/*` (átomos, cada uno con story). El widget `session-rail` de F1 se elimina.
Storybook 8 como devDependency; `npm run storybook`.

**Boundaries** (`arch/`): los 5 existentes siguen (one-turn-at-a-time y domain-sin-transport
con test). Se agrega **`git-solo-lectura-y-commit`** (nuevo, con fitness test): el paquete
`adapters/git` no contiene los strings/verbos `push|pull|fetch|reset|rebase` — la app jamás
mueve el repo del usuario contra remoto ni reescribe historia. `sesion-aislada-por-cwd` queda
como está (brecha declarada — la cierra PB-02 sobre este shell).

## 6. Migración desde F1

- El rail F1 muere; las sesiones existentes en `sessions.json` migran: sesión sin `RepoID` se
  agrupa bajo un pseudo-repo «(sin repositorio)» derivado de su cwd, hasta que el usuario las
  cierre. Sin pérdida de datos.
- `theme.css` ámbar → PRENTER teal (CAP-06 se actualiza en INCREMENTO al entregar).
- Storybook entra como tooling de `web/` (no toca el binario Go; `go:embed` sigue embebiendo
  solo el build de la SPA).

## 7. Reglas de negocio (RN)

- **RN-1** — No existe sesión sin historia ligada (v1: historia mock del board de ejemplo; el
  contrato ya es el definitivo). Config sin paquete redirige al picker, siempre.
- **RN-2** — Workspace y Sesión son 1:1. Un solo conjunto, dos ejes de vista (§0.5).
- **RN-3** — Overlays jamás tapan el rail de repos ni el studio-nav (§0.5, no renegociable).
- **RN-4** — La app NUNCA ejecuta push/pull/fetch/reset/rebase sobre el repo del usuario. v1
  ni siquiera los ofrece; los botones existen deshabilitados con tooltip honesto. Boundary con
  fitness test (§5).
- **RN-5** — Commit solo por selección explícita del usuario (pathspec); imposible commitear
  «todo» sin tildar todo.
- **RN-6** — Teal PRENTER = único acento saturado en toda la UI. Dark-first; light disponible
  por `data-theme`.
- **RN-7** — Todo elemento presente-pero-no-funcional declara su estado («todavía no», con
  destino: PB-NN) — principio 7 del norte, patrón botones voz del mockup.
- **RN-8** — Un turno a la vez por sesión (boundary existente); Encolar respeta la cola, no
  paraleliza.
- **RN-9** — Todo componente visual del shell existe primero como story en Storybook (átomo o
  composición) — PB-03 es parte del gate de cierre, no un anexo.
- **RN-10** — Microcopy en español neutro (seam `locale: es`); «Product» del mockup se rotula
  «Producto». Los términos Studio/Backlog quedan (nombres propios de producto).

## 8. Criterios de aceptación + verificación REAL (gate de cierre)

Verificación en la app viva (browser real), no «GET 200»:

- [x] **AC-1** — Registrar este mismo repo (`~/Proyectos/dev-studio`) por ruta desde la UI; el
  repo aparece con su contador. Registrar una ruta sin `.git` → error claro y visible.
- [x] **AC-2** — Flujo completo: `+ Nuevo Workspace` → picker → elegir historia → Config con
  chip del paquete → rol → `Crear sesión aislada` → la sesión aparece en el árbol, ligada
  (chip 📋 en el header), y RESPONDE un turno real de `claude`.
- [x] **AC-3** — Intentar crear sesión desde ⚙ sin paquete → hint + redirección al picker
  (RN-1 observada).
- [x] **AC-4** — Con archivos modificados reales en el cwd: panel Cambios lista los archivos;
  quick-look muestra el diff real; side-by-side navega N/M; tildar un subconjunto + mensaje +
  Confirmar → `git log` (terminal) muestra el commit SOLO con esos paths. Historial de la UI
  lo refleja.
- [x] **AC-5** — Botones push/pull/fetch deshabilitados; `grep` del adapter git no contiene
  verbos prohibidos; fitness test `git-solo-lectura-y-commit` en verde.
- [x] **AC-6** — Dos sesiones concurrentes en repos/wd distintos: panel Cambios de cada una
  muestra SU status sin cruce; turnos sin cruce (regresión CAP-02).
- [x] **AC-7** — Overlays Backlog/Config/Producto abren tapando solo chat+panel; rails
  siempre clickeables; Esc vuelve a Studio. Banners de honestidad visibles en
  Backlog/Config/Producto.
- [x] **AC-8** — `npm run storybook` levanta el catálogo con TODOS los átomos de §2 en dark y
  light; la app no contiene componentes visuales sin story (revisión de PR).
- [x] **AC-9** — Toda la UI renderiza tokens PRENTER (teal, Jost/Mulish/JetBrains Mono
  embebidas, sin ámbar residual) — verificación visual + grep `#a8742c` = 0.
- [x] **AC-10** — `go test ./...` + `go test ./arch/fitness/...` + build embebido verdes;
  sesiones F1 preexistentes visibles bajo «(sin repositorio)» tras migrar.

## 9. Fuera de alcance (explícito, con destino)

| Qué | Rebanada |
|---|---|
| Worktree + branch por sesión + validación de rutas protegidas (aislamiento real; dots `conflict`, dropdown ⋯, chips por rol) | **PB-02 — la siguiente, sobre este shell** |
| Terminal PTY embebida + Comandos + SSH | PB-09 / PB-17 |
| Board backlog real (Historia/Capability as code) | PB-07 (+ PB-08 proceso) |
| Registry de arneses + vista Roles funcional | PB-25 → PB-06 |
| Diff conversacional + adjuntar conversación al commit | PB-11 / PB-10 |
| push/pull/fetch desde la UI | sin PB — se decidirá con PR-flow (PB-14) |
| Clonar repo desde URL remota | sin PB — nace como idea nueva si se quiere |
| Segundo proveedor (Codex/OpenCode en el select) | PB-20 |
| Checkpoints/revert por turno | PB-12 |

## 10. Dudas abiertas

- [ ] Ninguna bloqueante al ratificar. (La escala tipográfica compacta de app (§2) es
  decisión de esta spec — si PRENTER formaliza una escala de producto después, se
  re-sincroniza.)

## Changelog

- v1 2026-07-07 — borrador post-forks (paleta/roster/alcance/storybook firmados en sesión).
- v1 CONGELADA 2026-07-07 — ratificada por Chris (DH-15) sin ajustes.
- v1 VERIFICADA 2026-07-07 — AC-1..AC-10 tildadas: 24/24 checks en vivo (puppeteer + Chrome
  contra la app real; commit por pathspec probado con `git show`; concurrencia sin cruce;
  migración F1 real) + suite Go/fitness/build-storybook verdes. Evidencia: ficha DH-15.
- EVOLUCIÓN 2026-07-07 (PB-27, DH-17): §4.1 — «+ Nuevo Workspace» ya no salta directo al
  picker: abre el wizard ubicación→propósito (`specs/nuevo-workspace/SPEC.md`). RN-1
  evolucionada: «toda sesión CON EDICIÓN liga a un paquete»; exploración read-only no lo
  necesita.
