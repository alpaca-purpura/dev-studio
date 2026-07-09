# Ledger — DevStudio (fichas DH-NN)

> Registro de decisiones de ESTE producto. Mismo formato/disciplina que la casa prenter-harness.
> **Continuidad:** DH-01..DH-11 viven en `prenter-harness/products/devhub/LEDGER.md` (la
> incubadora, congelada en la graduación). Este repo arranca en **DH-12**.

## Fichas

### DH-12 · Fundación del repo propio — graduación de P2 con visión ampliada — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-04):* "quiero crear una Épica para trabajar la interfaz de esta
aplicación e ir trayendo lo que hemos desarrollado de a pocos conforme a cómo hemos decidido
será la experiencia del developer/usuario. El objetivo es que el usuario (puede ser CTO,
developer, devops, product owner) trabaje en esta aplicación dentro de una organización de
forma 'orquestada', con otros usuarios complementarios trabajando a la par uno o varios
sistemas a la par, apalancándose de un repositorio (GitHub) como conector entre todos. Esto ya
creció más de lo que buscábamos y será un producto propio con su propio repositorio… una
aplicación cuyo objetivo es construir y mantener software basado en proceso y arquitectura,
donde usamos arquitectura as code, proceso as code, software documentation as code, etc. Este
nuevo repositorio debe nacer limpio, con la nueva visión, y como arneses de construcción vamos
a usar los del KIT DEV (plugin), pero lo iremos evolucionando conforme avanzamos." Forks
firmados por AskUserQuestion: nombre **DevStudio** (repo `dev-studio`, binario `dev-studio` — cierra de
nacimiento la deuda del rename DH-01) · **graduación** con célula del monorepo = fuente
CONGELADA read-only del port gradual (mecánica I-69 adelantada deliberadamente: sin clientes
aún) · repo GitHub **privado en alpacapurpura** desde el día 0.

*Desarrollo:* repo nace limpio con VISION.md (visión ampliada: construir y mantener software
basado en proceso y arquitectura; as-code en las 3 superficies; trabajo orquestado
multi-usuario/multi-rol; GitHub = conector) + este ledger + épica «Experiencia Orquestada»
(borrador de norte, F0 pendiente de firma EN este repo) + kit dev instalado como plugin desde
el marketplace (`alpacapurpura/prenter-marketplace`). Decisiones técnicas heredadas que siguen
vigentes: DH-10 (conexión CC = driver CLI-nativo, BYO licencia, SIN API) · descriptor de
proceso I-77 (contrato de ecosistema — el kit lo shipea, DevStudio lo interpreta) · binario Go +
UI embebida. El código go+ui NO se copió: el port es gradual, gobernado por la épica, pieza
por pieza según la experiencia decidida.

*Conecta:* DH-01..DH-11 (la historia en la incubadora; DH-11 = la ficha espejo de esta
graduación en el monorepo) · I-NN de ecosistema (graduación — registrada en
`prenter-harness/tooling/strategy/LEDGER.md`) · I-77 (descriptor) · DH-10 (driver CLI-nativo)
· KIT-06 (marketplace del que este repo consume su arnés).

*Siguiente:* F0 de la épica «Experiencia Orquestada» — norte firmado (journey por rol +
principios de la experiencia) + auditoría de herencia (qué pieza del monorepo entra primero y
en qué forma).

### DH-13 · F1 esqueleto de la app — driver CLI-nativo + rail multisesión — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-06):* "revisa todos los aspectos técnicos de harness-studio y hagamos
un hermano de la arquitectura para construir la aplicación de escritorio que se conecta con
claude code, manejando múltiples sesiones claude code entre las que puedo navegar al cambiar de
tab de mi aplicación y crear o cerrar tmb desde mi aplicación. Copia también las buenas prácticas
de creación de visión, arquitectura as code, etc. Copia el estilo gráfico del storybook a la
fecha y cómo se ha manejado visualmente las pestañas de sesión como columna colapsable a la
izquierda. Todo lo demás dejalo en blanco. Revisa la arquitectura as code de lo correspondiente
y crea una para nuestro proyecto adaptándola. Levanta la aplicación y verifica que esté
funcionando."

*Desarrollo:* auditoría técnica de `harness-studio` (hermano de arquitectura, mismo operador) en
tres ejes: (1) `arch/` + VISION/LEDGER/METODOLOGIA/CLAUDE.md como patrón documental — L1↔L2 +
checklist evaluable por boundary, fichas `Cruda/Desarrollo/Conecta/Siguiente`; (2) el driver
CLI-nativo: `internal/adapters/agent/claudecode` (subproceso `claude -p --input-format
stream-json --output-format stream-json`, un proceso por sesión, `SessionService` en memoria +
persistencia JSON liviana, SSE multiplexado); (3) el estilo gráfico Storybook (tokens Tailwind v4
`theme.css`, paleta `--primary #a8742c`/dark `#d9a35b`) y el widget `SessionRail` (rail
izquierdo colapsable 224px↔52px, tabs estilo WARP, acento `shadow-[inset_3px_0_0_var(--primary)]`
en la activa). Se construyó el **hermano real** en este repo (no fork, reimplementación propia
contra el protocolo verificado con una corrida real de `claude`): módulo Go
`github.com/alpacapurpura/dev-studio`, arquitectura hexagonal (`domain/ports/usecase/adapters`),
API REST + SSE (`/api/sessions`, `/events`), SPA React 19 + Zustand embebida vía `go:embed`,
`SessionRail` copiado 1:1 en estilo y simplificado en dominio (sin `arnes`/`salud`/`view` — eso
es de ArnesIA). Arquitectura as code propia en `arch/` (5 boundaries, 2 con test real corriendo:
`TestOneTurnAtATime`, `TestDomainNoTransportImport`) — deliberadamente más chica que la del
hermano y con una brecha de seguridad documentada explícitamente (`sesion-aislada-por-cwd`: sin
validación de rutas protegidas todavía). Verificación end-to-end real (no simulada): 2 sesiones
creadas desde el navegador (Chrome DevTools MCP), cada una spawneando su propio proceso `claude`
con `claude_session_id` distinto, turnos concurrentes respondidos correctamente y sin
cruzarse (`ALPHA`/`BETA`), navegación entre tabs preservando conversación, colapso del rail,
cierre de sesión limpiando el proceso.

*Conecta:* DH-10 (driver CLI-nativo, decisión heredada que esta ficha materializa por primera
vez en código propio) · DH-12 (fundación del repo) · épica
[`epicas/experiencia-orquestada/NORTE-FIRMADO.md`](./epicas/experiencia-orquestada/NORTE-FIRMADO.md)
(entonces `NORTE-BORRADOR.md`; candidato de fase F1 «esqueleto de la app») · `harness-studio` (proyecto hermano, misma
metodología, fuente del patrón `arch/` y del estilo `SessionRail`) · `arch/INDEX.md` (el árbol
nuevo) · `arch/boundaries/sesion-aislada-por-cwd.md` (brecha de seguridad declarada, no oculta).

*Siguiente:* **DH-14 = cerrar la brecha de `sesion-aislada-por-cwd`** (validación de rutas
protegidas antes de exponer la app fuera de este equipo) + F0 de la épica (norte firmado,
todavía pendiente pese a que F1 ya entregó código) + decidir si el segundo adaptador de agente
(prueba real de intercambiabilidad) entra antes o después del port por rebanadas.

### DH-14 · F0 norte firmado de la épica — rol = arnés instalado desde registry propio — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-07):* "los 'roles', el 'proceso', en sí vienen de los 'arneses' que
estamos creando en ArnesIA. Es decir, al instalar la aplicación me conectaré a un repositorio
de plugins y desde allí en la vista 'Roles' miraré cuáles hay (Todos son complementarios) y
añadiré los que me correspondan o los que mi jefe me ha permitido de acuerdo a mi puesto. Y
esto es por Proyecto y de acuerdo a eso se pinta el backlog y con esos arneses son los que
trabajará el claude code en ese proyecto específico."

*Desarrollo:* sesión F0 de la épica «Experiencia Orquestada» (= PB-01, gate en banda 🔴) —
norte debatido y firmado bloque por bloque vía AskUserQuestion. Firmas: **(1) journeys ×4**
(spec clon §4): Developer/PO/DevOps tal cual con el reframe integrado; CTO como HUECO DECLARADO
(paso 1 —vista multi-proyecto— bloqueado por PB-21→PB-22 desde que la Torre se eliminó).
**(2) 7 principios de la experiencia** (la estrella es quien orquesta · la experiencia manda el
port · todo deriva del dato as-code · nada huérfano sesión→historia→capability · aislamiento
por defecto · el repositorio es el conector · honestidad de superficie). **(3) El reframe** (la
decisión más grande): Rol = arnés instalado desde marketplace — SOLO marketplace, cero roles
locales; mecanismo = **registry PROPIO de DevStudio** (fork firmado contra la recomendación de
usar el sistema de plugins nativo de CC); permisos por puesto = con multi-usuario. Supersede la
definición de Rol de spec §3 («plantilla local + catálogo curado»). **(4) Orden de herencia**:
Proyecto → Rol → Sesión-workspace → Historia/Capability → Proceso (spec §3 tal cual; PB-02
sigue adelantada en 🔴 por seguridad). **(5) Dudas §9**: multi-usuario TBD formal sin bloquear
(PB-21 dueño Chris) · PB-02 = worktree + validación de rutas protegidas en la misma entrega ·
catálogo disuelto por el reframe · operonapp cerrada sin revisar. **(6) Slots del seam**:
domain_modules spec §3 + `registry` · agent_roster ajustado al repo real (backend-go /
frontend-web) · value_stream 10 estados completos (incluye parked/dropped) · wip_caps advisory
(se visualizan, no restringen). Artefactos del cierre: `NORTE-BORRADOR.md` →
`NORTE-FIRMADO.md` · slots `__FILL_ME__` rellenos en `project.config.yaml` · spec §9 marcada
resuelta · BACKLOG: PB-01 entregada (gate, sin CAP), PB-06 redefinida por el reframe, **PB-25
nueva** (registry propio), PB-21 desbloqueada a TBD-con-dueño.

*Conecta:* DH-12 (I-77: «el kit lo shipea, DevStudio lo interpreta» — hoy extendido al rol
entero) · DH-13 (F1 entregada; la deuda «F0 sin firmar» que esta ficha salda) · DH-10 (driver
CLI-nativo intacto: el arnés se materializa como artefactos que el `claude` spawneado carga) ·
[`epicas/experiencia-orquestada/NORTE-FIRMADO.md`](./epicas/experiencia-orquestada/NORTE-FIRMADO.md)
(el norte) · `CLON-AGENTSROOM-SPEC.md` §9 (dudas resueltas) · `project.config.yaml` (seam
completo) · BACKLOG PB-01/PB-06/PB-21/PB-25.

*Siguiente:* banda 🔴 del BACKLOG: PB-02 (workspace aislado + validación de rutas — forma
firmada hoy) y PB-03 (Storybook + atomic design). Diseñar el registry propio (PB-25) antes o
junto a la rebanada Rol (PB-06). Promover lo permanente de la épica a `specs/` recién al cierre
de la ÉPICA.

### DH-15 · Shell primero — spec congelable del shell PRENTER (PB-04 ⊕ PB-03) — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-07):* "necesito que antes de continuar sobre la construcción de la
aplicación vayamos cerrando cosas, como por ejemplo el shell, tengo un mockup html que quiero
que revises, entiendas y hagamos el spec y con ello construyamos siguiendo la arquitectura y
allí ya revisar los workspace aislados, pero ya en la aplicación." Y sobre la paleta: "eso lo
hemos extraído del system design que creamos para prenter que está en ~/Proyectos/prenter …
si te conectas a Claude Design puedes extraerlo, el system design se llama PRENTER allí."

*Desarrollo:* revisión completa del mockup (`mockup-clon-agentsroom.html`, 1634 líneas: fuente
+ 5 estados renderizados con Chrome headless aislado — el MCP compartido estaba lockeado por
otra sesión viva, HB-73 respetado). Hallazgos: el mockup EVOLUCIONÓ sobre §0.5 — `studio-nav`
de 5 ítems (Studio/Backlog/Product/Roles-pronto/Config) reemplaza los 2 botones de header;
overlays nuevos Product (mapa por zonas Estudio/Plataforma/Infraestructura) y Roles; el rail
de sesiones F1 se FUSIONA con el árbol Repositorios→Workspaces (workspace=sesión 1:1, un solo
conjunto en dos ejes). Confirmado: la paleta del mockup ES el design system **PRENTER** (los
swatches del boceto son los tokens de marca) — tokens extraídos del proyecto Claude Design
«PRENTER Design System» (`a98c2e0d-db82-43f2-8fd7-e7e05c40fd51`, `tokens/*.css`) y
sincronizados a `specs/shell/tokens/` con etags anotados; `brand-guidelines.md` de
`~/Proyectos/prenter` como fuente de reglas de marca (teal único acento, dark-first, Jost/
Mulish/JetBrains Mono como sustitutos vigentes). **Forks firmados:** (1) paleta = PRENTER
(theme.css ámbar muere; CAP-06 se actualizará al entregar); (2) Config/roster = placeholder
honesto del registry (banner PB-25, muere «+ Rol personalizado»); (3) alcance = AMBICIOSO:
además del shell cableado, **panel Cambios git REAL** (status/diff/quick-look/side-by-side/
log + commit por pathspec explícito; push/pull/fetch NO EXISTEN en el adapter — boundary
nuevo `git-solo-lectura-y-commit` con fitness test); (4) Storybook DENTRO de la rebanada
(PB-03 fusionada: átomo sin story no entra al shell, RN-9). Artefactos: **`specs/shell/SPEC.md`
v1** (primera spec permanente del repo — estructura real/stub por zona, flujos picker→crear
sesión, arquitectura hexagonal con ports `RepoRegistry`/`GitInfo`/`GitCommit`, RN-1..10,
AC-1..10 con verificación real, fuera-de-alcance con destino por PB) + BACKLOG rebandeado
(PB-04→🔴 en-curso absorbe PB-03 · PB-02 después del shell «ya en la aplicación» · PB-05
parcialmente adelantada).

*Conecta:* DH-14 (norte firmado que esta rebanada materializa; reframe rol=arnés → roster
placeholder) · DH-13 (F1 cuyo rail muere en este shell) · DH-10 (BYO CLI — el adapter git
sigue el mismo espíritu: subproceso del `git` del usuario) · PB-04/PB-03/PB-02 ·
`specs/shell/SPEC.md` · `specs/shell/tokens/` · marca PRENTER (`~/Proyectos/prenter/marketing/
brand-guidelines.md` + Claude Design `a98c2e0d`).

*Entrega (mismo día, 2026-07-07):* spec ratificada (CONGELADA) y rebanada CONSTRUIDA +
VERIFICADA: tokens PRENTER en `theme.css` (dark default + light) con Jost/Mulish/JetBrains
Mono vendorizadas (3 woff2 en el binario, sin Google Fonts en runtime) · backend TDD (adapter
`git/cli` con 6 tests reales sobre repos git de fixture; `RepoService` con validación `.git`;
store `state.json` con migración automática de `sessions.json`; endpoints repos+git) ·
boundary `git-solo-lectura-y-commit` ENFORCED (2 fitness tests: scan de verbos prohibidos en
código + pathspec obligatorio) · Storybook 10 con átomos + organismos (RN-9) · shell completo
(repos-rail, studio-nav, session-view con Encolar real, changes-panel, 4 overlays). El rail F1
y el chat-panel viejos ELIMINADOS. **Verificación real: 24/24 checks** del gate (AC-1..AC-7
vía puppeteer + Chrome contra la app viva: registro de repos con rechazo de ruta inválida,
picker→sesión ligada→turno `SHELL-OK` real, commit por pathspec verificado con `git show`
—el archivo destildado quedó fuera—, 2 sesiones en repos distintos con turnos concurrentes
`ALFA-1`/`BETA-2` sin cruce, rails clickeables bajo overlays por hit-test, Esc→Studio) +
AC-8/9/10 (build-storybook verde · grep ámbar=0 · `go test ./...` + fitness verdes · migración
real: 2 sesiones F1 visibles bajo «(sin repositorio)»). Bugs cazados en verificación: selector
zustand fabricando `[]` por snapshot (loop React #185) y falso-positivo del script por
`text-transform: uppercase` — ambos documentados en el fix. INCREMENTO: CAP-01/03/04/05/06
actualizadas + **CAP-07/08/09 nuevas**.

*Siguiente:* **PB-02 — workspace aislado por sesión (worktree + branch + validación de rutas
protegidas), ya sobre este shell** (forma firmada en F0; el rail ya pinta branch/status/dots
para recibirla). Después: banda 🟡 (PB-05 restos → PB-25 registry → PB-06 Roles).

### DH-16 · Workspace aislado por sesión + instalable dogfooding (PB-02 ⊕ PB-26) — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-07):* "dale, arranquemos PB-02 y crea el instalador y el mecanismo
para instalarlo y actualizar la aplicación y verla y probar todo en adelante desde el mismo
app."

*Desarrollo:* spec `specs/workspace-aislado/SPEC.md` (segunda permanente) ratificada con 3
forks: worktrees en `~/.dev-studio/workspaces/{repo}/{slug}` (fuera del repo del usuario,
patrón Conductor) · cierre de sesión PREGUNTA (modal conservar/borrar) · updater = REBUILD
LOCAL (sin red, sin tokens, sin git — compila el working tree del repo fuente). **PB-02:**
puerto `GitWorkspace` (CreateWorktree/RemoveWorktree, TDD con 4 tests sobre repos reales:
aislamiento, colisión de slug → sufijo, HEAD unborn → error honesto, remove sucio rechazado)
· toda sesión de repo nace en su worktree con branch `wt/{slug}` (`Session.Workspace/Branch`)
· `domain.RutaProtegida` (pura: $HOME, ~/.ssh/.gnupg/.aws/.kube/.docker/.dev-studio, raíces
de sistema) aplicada en TODA puerta (Register + createSession legacy → 422) · el rail pinta
branch/±N/dot del worktree DE CADA SESIÓN (dot `conflict` parseando unmerged) · modal de
cierre con RN-3 (borrar sucio → git rechaza, la app conserva y muestra el error REAL —
bug de unmount del aviso cazado y arreglado EN el gate). **Boundary
`sesion-aislada-por-cwd`: `proposed` → `enforced`** (v2.0 — 3 de 4 checks con test; resta
token de capacidad, documentado). **PB-26:** `scripts/install.sh` (SPA + binario con
`-ldflags` versión SHA+dirty/fecha → `~/.local/bin/dev-studio` + lanzador `dev-studio-open`
(app-mode Chrome) + `.desktop` + icono PRENTER SVG + `app.json` con el source) · `GET
/api/version` + `POST /api/update` (rebuild vía `install.sh --update` con `bash -lc`, 500
con output del compilador si falla, restart por `syscall.Exec` del binario nuevo) · footer
del rail: versión visible + «Actualizar» con poll de `version@build_date` → `location.reload()`.
**Verificación real: 14/14 checks** en vivo contra el binario INSTALADO (puppeteer):
worktree en disco + `git worktree list` + branch en rail · 2 sesiones mismo repo sin cruce
de Cambios · $HOME rechazado visible · modal: borrar sucio conservado con aviso + borrar
limpio desaparece · **la app se ACTUALIZÓ A SÍ MISMA** (AC-6: rebuild + restart + un fix
Go commiteado solo en el working tree quedó VIVO tras el update — 422 verificado post-restart).
Suite completa verde (5 paquetes + fitness 6 tests).

*Conecta:* DH-14 (forma worktree+rutas firmada en F0) · DH-15 (el shell que recibe esto;
CAP-09 pierde su ⚠) · DH-13/DH-10 (driver intacto — el subproceso `claude` ahora nace en el
worktree) · `arch/boundaries/sesion-aislada-por-cwd.md` v2.0 · `arch/boundaries/git-solo-lectura-y-commit.md`
(worktree add/remove NO tocan los verbos prohibidos — scanner sigue verde) · PB-24 (el
instalable comercial sigue pendiente; esto es la versión casa) · BACKLOG PB-02/PB-26.

*Siguiente:* dogfooding real — Chris abre DevStudio desde el lanzador y prueba las rebanadas
desde la app instalada («Actualizar» tras cada entrega). Banda 🟡: PB-05 restos → PB-25
(registry) → PB-06 (Roles). Deuda visible del boundary: token de capacidad de la API local
(check TBD). La brecha de seguridad fundacional quedó CERRADA — la app ya puede verla gente
fuera del equipo (con criterio).

### DH-17 · Primer feedback de dogfooding: 3 fixes + flujo «Nuevo Workspace» v2 con taxonomía estándar — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-07, tras probar la app instalada):* "1. Cuando ingreso por primera
vez me sale este mensaje [bienvenida de Chrome]… 2. el colapsamiento de repositorios… no está
funcionando… Cuando hago click en Nuevo Workspace me lleva a backlog, no debería hacer nada.
3. cuando escribo algo en studio y mando para que me responda claude code no me responde." Y
al proponer el flujo: "al momento de crear nuevo workspace, primero debo preguntar si es un
nuevo worktree o el worktree actual, luego el usuario debe poder: navegar y comprender (podrá
usar claude code pero no tendrá permisos de edición), luego podrá trabajar sobre una historia
de usuario, bugfix, etc existente (ayudame investigando el estándar en la industria que todo
developer maneje) o crear uno."

*Desarrollo:* **Hotfix DH-16.1** (commit propio): (1) lanzador con `--no-first-run
--no-default-browser-check` + PATH completo; (2) colapso universal del rail (incl. «(sin
repositorio)») + reveal al activar sesión oculta; (3) claude mudo = el .desktop hereda PATH
sin `~/.local/bin` → resolución del binario con fallbacks + la UI muestra el error real del
turno (antes quedaba «streaming» muda) — verificado con turno `FALLBACK-OK` bajo PATH pelado.
**PB-27** (spec `specs/nuevo-workspace/SPEC.md`, 3 forks ratificados): taxonomía estándar de
paquetes de trabajo = convergencia Jira/Scrum issue-types + Gitflow branch-naming +
Conventional Commits → 5 tipos (`historia·bug·hotfix·tarea·spike`) con branch por prefijo
(`feature/ bugfix/ hotfix/ chore/ spike/` — muere `wt/`); wizard 2 pasos (¿dónde? worktree
nuevo vs checkout actual → ¿para qué? explorar / ítem existente / ítem nuevo) con guard
«checkout = solo explorar» enforced en backend (422); **sesiones de EXPLORACIÓN read-only**
vía `--permission-mode plan` (mecanismo nativo de la CLI) — RN-1 evolucionada: «toda sesión
CON EDICIÓN liga a un paquete»; crear ítem desde el wizard (tipo+título, persiste con la
sesión hasta PB-07). **Verificación: 10/10 checks** contra el binario instalado — la joya:
pedido de edición en exploración BLOQUEADO por plan mode («no puedo crear archivos», archivo
inexistente) con el flag visto en el proceso real; branches `spike/…` y `bugfix/…` confirmadas
en `git worktree list`; guards API 422×2.

*Conecta:* DH-16 (la entrega que este feedback ejercitó — el loop dogfooding FUNCIONA: instalar
→ probar → feedback → fix → Actualizar) · DH-14 (RN-1 evoluciona sin perder el espíritu) ·
specs shell §4.1 y workspace-aislado §2.1 (changelogs anotados) · PB-27 · PB-07 (los ítems
creados migran al board real).

*Siguiente:* Chris aprieta «Actualizar» y prueba el wizard. Banda 🟡: PB-05 restos → PB-25
(registry) → PB-06 (Roles). Nota: los ítems tipo `historia` deberán exigir capability cuando
exista el board real (PB-07) — la regla de la casa no se negocia, hoy no hay dónde elegirla.

### DH-18 · Registry de arneses — DevStudio adopta el estándar ArnesIA (PB-25 ⊕ PB-05) — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-07):* "Solo te pido revises y converses con ~/Proyectos/harness-studio
porque aquí se crean los arneses, el formato que necesites debes pasarselo de alguna forma que
pueda constuirlos así, o ve como los construye según su doctrina y te adaptas, pero no
construyas por construir." Y ante el análisis: "quiero saber cuál es mejor, debido a que
arnesia esta construyendose tmb… si nos adaptamos o modificamos arnesia o un cruce de ambos."

*Desarrollo:* sesión banda 🟡. **(1) Auditoría PB-05:** vaciada por el reframe DH-14 — su único
resto con journey (roster por proyecto + conexión del registry) ES PB-25 → **fusionada** (fork
firmado); metadatos declarados sin journey. **(2) Revisión profunda de ArnesIA** (mandato «no
construyas por construir»): la fábrica YA firmó e implementó su formato —
`arch/contracts/nomenclatura-arnes.md` v1 (forma-plugin: plugin.json + `arnes.l0.json` con meta
rol×proceso + spine + skills con contrato + CLAUDE.md banda Base; loader round-trip 13/13) —,
su VISION ya dibuja el eslabón exacto («publica → marketplace git → instala → proyecto →
ejecuta ← apps de rol»: DevStudio ES la app de rol), el marketplace git corre en producción
(prenter-marketplace: marketplace.json + catalogo.json) y su METODOLOGIA §9 manda inyección
por flags, jamás escribir la maquinaria en el árbol del proyecto. Dos forks firmados en sesión
quedaron SUPERSEDIDOS el mismo día y se RE-firmaron: **formato = ArnesIA as-is** (muere el
arnes.yaml propio) y **materialización = CRUCE**: lock as-code `.devstudio/arneses.yaml`
committeado (roster viaja por GitHub, modelo npm: payload rehidratable) + caché forma-plugin
INTACTA en `~/.dev-studio/arneses/` + inyección por sesión con los MISMOS flags que la fábrica
usa para su kit (`--plugin-dir` + system prompt). Prompt de interop entregado a Chris para la
sesión ArnesIA (ratificar formato del publish fase-5 · `categoria` semántica en el spine
[puente a I-77, para PB-08] · campo `nombre` · postura multi-arnés/lock). **(3) Construcción
TDD:** dominio `Arnes`/lock · puertos RegistrySync/RegistryCatalog/ArnesLock/ArnesCache ·
adapters `registry/gitsync` (clone/pull confinado POR CONSTRUCCIÓN a `~/.dev-studio/registry/`
— boundary git-solo-lectura-y-commit v1.1 con fitness nuevo `registry-sin-push`) +
`registry/fscatalog` + `arneses/lockfile` (YAML) + `arneses/cache` · `ArnesService` (conectar/
instalar/desinstalar/inyección con rehidratación) · commit del lock por pathspec · guard RN-5
(422) · SpawnOpts+buildArgs · Config overlay v2 (sección Registry + roster real por proceso +
catálogo instalar/desinstalar — **MOCK_ROSTER MUERTO**, grep=0) + stories RN-9.
**(4) Verificación: 7/7 checks** contra el binario instalado **con el arnés REAL
`dev-full-cycle` de ArnesIA** (interop de verdad): catálogo con rol/proceso/fases del
manifiesto real · caché `diff -r`=0 · commit solo-lock (`git show`) · sesión con rol desde el
wizard → flags en `/proc/{pid}/cmdline` → **el turno respondió EXACTO las 4 skills namespaced**
(`dev-full-cycle:builder/releaser/reviewer/spec-writer`) · ciclo desinstalar + 422s. **La app
se actualizó a sí misma 2 veces EN el gate** y cada vuelta cazó un bug real: (a) updater muerto
bajo el PATH pelado del .desktop (`npm` de nvm no entra por `bash -lc` → fallbacks explícitos
en install.sh); (b) `--append-system-prompt` y `--append-system-prompt-file` son EXCLUYENTES
en la CLI (proceso claude defunct al nacer) → la banda Base viaja EMBEBIDA en el único prompt.
**(5) Feedback dogfooding en la misma sesión (DH-18.1):** globo «No se puede actualizar
Chrome» dentro de la ventana app → flag `OutdatedBuildDetector` deshabilitado en el lanzador
(la causa raíz es el Chrome del sistema; Tauri sigue siendo PB-24) · **Config se muda a su
zona**: sale del studio-nav (zona por sesión) al footer del rail de repositorios (zona app,
patrón ArnesIA) y versión+«Actualizar» viven DENTRO del overlay (§ Aplicación) — spec shell
§3.2 evolucionada con changelog; verificado en vivo tras otro self-update.

*Conecta:* DH-14 (el reframe que esto materializa — rol = arnés instalado, cero roles locales)
· DH-10 (intacto: archivos en disco + flags nativos, cero API) · DH-16/17 (el loop dogfooding
que cazó los 2 bugs) · ArnesIA HS-10/HS-11 (nomenclatura v1 + los 3 cuerpos — los contratos
cross-repo que esta ficha consume) · I-77 (spine⟷descriptor: reconciliación diferida a PB-08,
pedido `categoria` en el prompt interop) · `specs/registry-arneses/SPEC.md` · BACKLOG
PB-05/PB-25/PB-06 · INCREMENTO CAP-13/CAP-02.

*Siguiente:* Chris pega el prompt interop en la sesión ArnesIA (las respuestas entran por
changelog de la spec — nada bloquea). Banda 🟡: PB-06 (vista Roles rica) → PB-07 (Historia/
Capability). Registry duradero RESUELTO en el mismo cierre: `~/Proyectos/marketplace-arneses`
(repo git semilla con dev-full-cycle de ArnesIA + README que declara «los arneses nacen en
ArnesIA, acá solo se publican»); la app quedó CONECTADA a él (persistido en state.json) y el
lock de demo-gamma re-apunta. Cuando ArnesIA entregue su publish fase-5, alimenta ese repo (o
su URL en GitHub) sin tocar DevStudio.

**(6) Interop RESPONDIDO (DH-18.2, mismo día — ficha HS-12 de la fábrica, 4 firmas + 1
reparación):** los 5 pedidos volvieron y entraron por changelog v1.1 de la spec (mecanismo
firmado arriba). (1) Publish fase-5 RATIFICADO con contrato estable (marketplace.json +
forma-plugin + catalogo.json canales/versiones — evolución solo aditiva; base firme para
PB-06 upgrade). (2) `spine.categorias` ACEPTADO: enum FIJO = I-77 RN-28, terminalidad
derivada — ya en schema L0 + dogfood; semilla 0.1.0 sin mutar → PB-08 rehidrata del dogfood
actualizado o espera publish. (3) `arnes.l0.nombre` = campo canónico (ArnesIA REPARÓ su
schema: §2 lo nombraba pero lo rechazaba); cadena de fallback bendecida
`nombre → name → id` ídem descripcion → fscatalog ALINEADO con TDD (leía plugin.json
pisando `l0.descripcion` + faltaba el eslabón final →id; 2 tests nuevos, suite verde).
(4) Lock `.devstudio/arneses.yaml` BENDECIDO como superficie de auditoría in situ (detector
3° de nomenclatura-arnes v1.1, la fábrica lo lee read-only); **pedido recíproco CUMPLIDO:
contrato estable declarado** — campos `registry` + `arneses[].{id, version, canal}`, solo
aditivo (spec §2). (5) spine⟷I-77 CONFORME con nuestra lectura: spine(+categorias) =
subconjunto navegable canónico; gates/dueños DERIVADOS de contratos por caja; el arnés NO
shipea descriptor I-77 (materializado = export/proyección); derivación bendecida (dueño de
estado = caja cuya transición llega · dueño-caja = rol, resto = operador · terminalidad =
categoría). **PB-08 re-redactada con este modelo — el pendiente externo de DH-18 queda
CERRADO; nada quedó bloqueado.**

### DH-19 · F0 norte firmado — épica «Conversación terminal-auténtica» (Híbrido 3a) — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-09):* "quiero saber cómo se hizo esta web
https://www.anthropic.com/features/making-of-claude-code y cómo podríamos tomar el cómo se ve para la
conversación del dev studio, considerando que nos colgaremos de anthropic, quiero que se sienta lo más
real, tal como https://agentsroom.dev/es#demo … en el caso de agentsroom se muy detallista ya que se
puede conversar por el chat de conversación o escribiendo directo en su «terminal». Quiero dar ese
efecto." + "Elevar a épica para atenderla de una vez."

*Desarrollo:* sesión F0 de la épica «Conversación terminal-auténtica» (= PB-28, banda 🔴). Sobre el
research previo (docs `00`–`06` + `mock-conversacion.html` + Artifact) se firmó el norte **fork por
fork vía AskUserQuestion**. **Norte:** que la conversación se sienta como un terminal REAL (estética
making-of de Claude Code) con **doble input** — chat ↔ terminal, ambos al MISMO stdin del `claude` real;
AgentsRoom lo FINGE, DevStudio lo hace LITERAL (el driver ya spawnea stream-json). **Veredicto técnico
ratificado = Híbrido 3a:** un proceso stream-json → dos superficies (cards estructuradas + xterm.js en
modo controlado, deltas reales = typing auténtico, SIN PTY). **8 invariantes ratificados** (terminal
literal · doble input = un solo stdin · un proceso dos superficies · teal único + amber solo semántico ·
degradación por `Capabilities` · Envelope ACP-aligned · cero API/BYO licencia DH-10 intacto · boundary
`conductor-no-parsea-jsonl`). **Decisiones A–E firmadas:** **A** probe R0 primero = SÍ (protocolo de
control `can_use_tool` sub-doc + version-dependiente, bugs #34046/#12235; **`claude` v2.1.205 pineada**).
**B** = **MULTI-PROVEEDOR DESDE F1** (Chris DIVERGIÓ del recomendado «Claude primero»): consecuencia
firmada = el `Envelope` nace ACP-aligned en R1 (no refactor tardío en R4), `Capabilities()` +
`ProviderSessionID` (hueco #4) se adelantan al núcleo de F1, y el segundo adapter Amp (near drop-in)
entra al alcance FIRME (R5 deja de ser «fase 2 aparte») como prueba de intercambiabilidad de `AgentPort`
(absorbe PB-20). **C** terminal = faux/modo-controlado 3a (PTY real = escalón futuro). **D** slicing
LEDGER = **DH por rebanada entregada R1–R5** (cada una + CAP en su commit); esta firma de F0 = DH-19
aparte. **E** PB-24 (shell nativo Tauri/Wails) NO bloquea (3a corre en el HTTP+webview actual; se decide
aparte). **Fases reordenadas por B:** F1 (R0 probe + R1 tool-cards + núcleo Envelope/Capabilities) → F2
(permisos modal de rama + terminal xterm.js + doble input) → F3 (cierre normalización) → F4 (2do adapter
+ multi-CLI). Artefactos del cierre: `NORTE-BORRADOR.md` → `NORTE-FIRMADO.md` · `05-plan…` con decisiones
CERRADAS + §Reordenamiento-por-B · README actualizado · BACKLOG PB-28 = «F0 firmada».

*Conecta:* DH-10 (driver CLI-nativo intacto — el terminal auténtico se cablea sobre el proceso real, cero
API) · DH-13 (el `conductor.go` que esta épica extiende: los 4 huecos) · DH-14/DH-18 (el registry de
arneses = eje «qué colaborador», ortogonal al eje «qué runtime CLI» que esta épica abre; PB-25) ·
[`epicas/conversacion-terminal-autentica/NORTE-FIRMADO.md`](./epicas/conversacion-terminal-autentica/NORTE-FIRMADO.md)
(el norte) + `06-especificacion-mockup.md` + `mock-conversacion.html` (SSoT de forma) · BACKLOG PB-28
(absorbe PB-09; toca PB-20/PB-24/PB-25) · «Experiencia Orquestada» NORTE-FIRMADO (el shell que hereda).

*Siguiente:* **F1** — R0 probe del protocolo de control `can_use_tool` (spike 1 archivo descartable:
spawn `claude` con `--permission-prompt-tool stdio`, disparar un tool gated, loguear stdin/stdout crudo →
`.md` con el shape verificado) ANTES de R2. Después R1 (tool-cards + Envelope ACP-aligned + Capabilities
desde el inicio, por B). Cada rebanada = su DH + CAP al entregar. Promover a `specs/` recién al cierre de
la ÉPICA (probable `specs/conversacion-sesion/` + `specs/driver-multiproveedor/`).

<!-- Próximas: DH-20, DH-21, … -->

## Log

| Fecha | Decisión | Fichas |
|---|---|---|
| 2026-07-04 | Fundación del repo propio: graduación de P2 (nombre DevStudio confirmado — binario `dev-studio`, muere la colisión `cockpit` de nacimiento); visión ampliada = construir y mantener software basado en proceso y arquitectura, trabajo orquestado multi-usuario (CTO·developer·devops·PO), GitHub conector; célula del monorepo congelada como fuente del port gradual; kit dev como plugin del marketplace; épica «Experiencia Orquestada» sembrada (F0 por firmar). | DH-12 |
| 2026-07-06 | F1 esqueleto de la app entregado: driver CLI-nativo propio (subproceso `claude` + stream-json, multisesión real verificada con 2 procesos concurrentes) + rail de sesiones colapsable (estilo Storybook de `harness-studio` copiado 1:1, dominio simplificado) + arquitectura as code propia (`arch/`, 5 boundaries, 2 enforced con test) con brecha de seguridad documentada sin ocultar. | DH-13 |
| 2026-07-07 | F0 norte FIRMADO (cierra PB-01 y la deuda de DH-13): journeys ×4 (CTO = hueco declarado) + 7 principios + reframe Rol = arnés instalado desde REGISTRY PROPIO (solo marketplace, cero roles locales; permisos por puesto con multi-usuario) + orden de herencia Proyecto→Rol→Sesión-workspace→Historia/Capability→Proceso + dudas §9 resueltas (multi-usuario TBD sin bloquear; PB-02 = worktree + validación rutas) + seam completo (value_stream 10 estados, wip_caps advisory). | DH-14 |
| 2026-07-07 | Shell primero: mockup revisado (evolución sobre §0.5: studio-nav 5 ítems, overlays Producto/Roles, workspace=sesión 1:1 — rail F1 muere) + paleta confirmada = design system PRENTER (tokens sincronizados de Claude Design a `specs/shell/tokens/`) + spec v1 del shell escrita (`specs/shell/SPEC.md`, primera spec permanente: alcance ambicioso con panel Cambios git real sin push/pull, roster placeholder del registry, Storybook dentro — PB-03 fusionada en PB-04; PB-02 pasa a construirse DESPUÉS del shell, ya en la app). | DH-15 |
| 2026-07-07 | Shell PRENTER ENTREGADO (PB-04 ⊕ PB-03 → entregadas): tokens+fuentes vendorizadas, backend git TDD, boundary `git-solo-lectura-y-commit` enforced, Storybook (RN-9), rail Repositorios→Workspaces reemplaza al rail F1, panel Cambios real con commit por pathspec, flujo picker→sesión ligada. 24/24 checks del gate en vivo (turno real, commit selectivo probado con `git show`, concurrencia sin cruce, migración F1 real). CAP-07/08/09 nuevas. Siguiente: PB-02 sobre este shell. | DH-15 |
| 2026-07-07 | Workspace aislado + instalable ENTREGADOS (PB-02 ⊕ PB-26): toda sesión nace en su worktree `wt/{slug}` en `~/.dev-studio/workspaces/` + rutas protegidas en toda puerta + modal cierre conservar/borrar (sin --force jamás) → **boundary `sesion-aislada-por-cwd` ENFORCED (la brecha fundacional cerrada)**. Instalable dogfooding: install.sh + .desktop + icono + versión embebida + botón «Actualizar» (rebuild local + restart). 14/14 checks en vivo contra el binario instalado — incluida la app actualizándose a sí misma con un fix vivo post-restart. CAP-10/11 nuevas. | DH-16 |
| 2026-07-07 | Primer feedback de dogfooding: 3 fixes (bienvenida Chrome · colapso universal del rail · claude mudo por PATH del .desktop — resolución con fallbacks + error visible) + **PB-27 «Nuevo Workspace» v2**: taxonomía estándar 5 tipos con branch por prefijo (feature/bugfix/hotfix/chore/spike — muere wt/), wizard ubicación→propósito con guards backend, **exploración read-only** (`--permission-mode plan`, edición BLOQUEADA verificada en vivo), crear ítem desde el wizard. RN-1 evolucionada: «toda sesión con edición liga a un paquete». 10/10 checks. CAP-12 nueva. | DH-17 |
| 2026-07-07 | Registry de arneses ENTREGADO (PB-25 ⊕ PB-05 fusionada): DevStudio ADOPTA el estándar ArnesIA (nomenclatura-arnes v1 + marketplace git en producción — «no construyas por construir») con materialización CRUCE: lock as-code committeado + caché forma-plugin + inyección por sesión `--plugin-dir`+system-prompt (patrón HS-11 de la fábrica; DH-10 intacto). Roster mock MUERTO; guard RN-5. Boundary git-solo-lectura-y-commit v1.1 (registry confinado, fitness nuevo). 7/7 checks en vivo con el arnés REAL dev-full-cycle de ArnesIA — las 4 skills namespaced respondidas por el claude de la sesión. 2 bugs cazados EN el gate (updater PATH .desktop · flags de prompt excluyentes). CAP-13 nueva. Prompt interop → sesión ArnesIA. | DH-18 |
| 2026-07-07 | Interop ArnesIA RESPONDIDO (HS-12 → DH-18.2): publish fase-5 ratificado con contrato estable · `spine.categorias` aceptado (enum I-77 RN-28, terminalidad derivada) · `arnes.l0.nombre` canónico + cadena fallback bendecida (fscatalog alineado TDD) · lock bendecido como detector 3° + **contrato estable recíproco declarado** (`registry·id·version·canal`, solo aditivo, spec §2) · spine⟷I-77 CONFORME: gates/dueños se derivan de contratos por caja, el arnés no shipea I-77 (export/proyección) → PB-08 re-redactada. Pendiente externo de DH-18 CERRADO. | DH-18 |
| 2026-07-09 | F0 norte FIRMADO — épica «Conversación terminal-auténtica» (Híbrido 3a, PB-28 banda 🔴): doble input literal (chat ↔ terminal, mismo stdin del `claude` real) sobre el driver DH-10. 8 invariantes ratificados (terminal literal · teal único + amber semántico · Envelope ACP-aligned · degradación por Capabilities). Decisiones A–E: A probe R0 primero SÍ (`claude` v2.1.205 pineada) · **B multi-proveedor DESDE F1 (Chris divergió): Envelope nace ACP-aligned + Capabilities temprano + 2do adapter Amp en alcance firme = PB-20)** · C terminal faux/3a (PTY diferido) · D DH por rebanada R1–R5 · E PB-24 no bloquea. NORTE-BORRADOR→FIRMADO. Siguiente: F1 = R0 probe + R1 tool-cards. | DH-19 |
