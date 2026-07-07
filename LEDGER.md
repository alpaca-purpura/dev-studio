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

<!-- Próximas: DH-16, DH-17, … -->

## Log

| Fecha | Decisión | Fichas |
|---|---|---|
| 2026-07-04 | Fundación del repo propio: graduación de P2 (nombre DevStudio confirmado — binario `dev-studio`, muere la colisión `cockpit` de nacimiento); visión ampliada = construir y mantener software basado en proceso y arquitectura, trabajo orquestado multi-usuario (CTO·developer·devops·PO), GitHub conector; célula del monorepo congelada como fuente del port gradual; kit dev como plugin del marketplace; épica «Experiencia Orquestada» sembrada (F0 por firmar). | DH-12 |
| 2026-07-06 | F1 esqueleto de la app entregado: driver CLI-nativo propio (subproceso `claude` + stream-json, multisesión real verificada con 2 procesos concurrentes) + rail de sesiones colapsable (estilo Storybook de `harness-studio` copiado 1:1, dominio simplificado) + arquitectura as code propia (`arch/`, 5 boundaries, 2 enforced con test) con brecha de seguridad documentada sin ocultar. | DH-13 |
| 2026-07-07 | F0 norte FIRMADO (cierra PB-01 y la deuda de DH-13): journeys ×4 (CTO = hueco declarado) + 7 principios + reframe Rol = arnés instalado desde REGISTRY PROPIO (solo marketplace, cero roles locales; permisos por puesto con multi-usuario) + orden de herencia Proyecto→Rol→Sesión-workspace→Historia/Capability→Proceso + dudas §9 resueltas (multi-usuario TBD sin bloquear; PB-02 = worktree + validación rutas) + seam completo (value_stream 10 estados, wip_caps advisory). | DH-14 |
| 2026-07-07 | Shell primero: mockup revisado (evolución sobre §0.5: studio-nav 5 ítems, overlays Producto/Roles, workspace=sesión 1:1 — rail F1 muere) + paleta confirmada = design system PRENTER (tokens sincronizados de Claude Design a `specs/shell/tokens/`) + spec v1 del shell escrita (`specs/shell/SPEC.md`, primera spec permanente: alcance ambicioso con panel Cambios git real sin push/pull, roster placeholder del registry, Storybook dentro — PB-03 fusionada en PB-04; PB-02 pasa a construirse DESPUÉS del shell, ya en la app). | DH-15 |
| 2026-07-07 | Shell PRENTER ENTREGADO (PB-04 ⊕ PB-03 → entregadas): tokens+fuentes vendorizadas, backend git TDD, boundary `git-solo-lectura-y-commit` enforced, Storybook (RN-9), rail Repositorios→Workspaces reemplaza al rail F1, panel Cambios real con commit por pathspec, flujo picker→sesión ligada. 24/24 checks del gate en vivo (turno real, commit selectivo probado con `git show`, concurrencia sin cruce, migración F1 real). CAP-07/08/09 nuevas. Siguiente: PB-02 sobre este shell. | DH-15 |
