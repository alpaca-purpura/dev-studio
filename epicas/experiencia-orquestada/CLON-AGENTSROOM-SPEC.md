# Spec — Clon de AgentsRoom con espíritu propio (insumo F0 · sin firmar)

> ⚠ Insumo para **F0** de la épica [«Experiencia Orquestada»](./NORTE-BORRADOR.md) — no la
> reemplaza ni la firma. Aporta benchmark competitivo verificado + una propuesta concreta de
> modelo de dominio y de journeys por rol, que es exactamente lo que F0 pide debatir y firmar.
> La mayoría de este documento es propuesta sin ratificar — **excepto §0.5**, que documenta
> decisiones de shell/navegación YA ratificadas por Chris y YA construidas en el mockup. Vive
> acá (carpeta temporal de la épica) y se borra o se promueve a `specs/` cuando F0 cierre.
>
> **SSoT de la experiencia:** este archivo + `mockup-clon-agentsroom.html` (mismo folder) son la
> fuente única del shell. Todo pedido de cambio al mockup se aplica en LOS DOS a la vez — el
> mockup para la forma, §0.5 (o la sección que corresponda) para la regla/decisión detrás. Si un
> cambio de shell no queda escrito acá, no pasó.
>
> **Fuentes:** exploración en vivo de `agentsroom.dev/es` (landing + demo interactivo completo vía
> Chrome DevTools MCP, sin muro de login — capturas en el scratchpad de la sesión) · research de
> `conductor.build` (docs públicas, vía WebSearch) · auditoría propia de este repo y de la célula
> congelada `prenter-harness/products/devhub` (herencia disponible, no portada). `operonapp.dev`
> no tuvo presencia indexada verificable — ver nota §1.

## 0. Encargo original (verbatim, 2026-07-06)

Clonar `agentsroom.dev/es` a nivel de spec, mejorándolo con nuestro espíritu. El *workspace* es
un **proyecto** de desarrollo de software; dentro caben **múltiples sesiones**. Todo sistema
tiene **capabilities**; lo que se implementa son **historias de usuario** que incrementan esas
capabilities. Lo que AgentsRoom llama "agente" nosotros lo llamamos **rol**. La estrella del
producto no es el rol — es **el usuario que orquesta todos los roles**.

> **Nota de terminología (2026-07-06):** el encargo verbatim de arriba usa *workspace* para
> "proyecto de software" — esa palabra queda verbatim a propósito, no se toca. §0.5 abajo formaliza
> el renombre vigente: **Repositorio** = ese mismo concepto (proyecto), **Workspace** pasa a nombrar
> la instancia aislada en disco (antes llamada *Worktree*).

## 0.5 Shell y navegación — YA RATIFICADO Y CONSTRUIDO (2026-07-06)

> A diferencia del resto de esta spec (propuesta, F0 sin firmar), lo que sigue **ya se decidió
> con Chris y ya está construido** en `mockup-clon-agentsroom.html`. Esta sección es la fuente
> de la regla; el mockup es la fuente de la forma. Cambian juntos.
>
> **Sincronizado 2026-07-06:** `mockup-clon-agentsroom.html` ya usa los nombres nuevos
> (`.repos-rail`, `.repo-header`, `.workspace-item`, `data-repo`, chip `multi-repositorio`, etc.).
> El prefijo de branch `wt/` y las rutas `.worktrees/` quedan iguales — convención de infra
> aparte del label de producto.

**Navegación de más alto nivel — un solo sidebar, no tabs.** Reemplaza "Torre / Proyecto /
Backlog / Roles" como 4 tabs paralelas (el inventario original de §7 antes de esta revisión) por
una jerarquía de 2 niveles en un único sidebar izquierdo (`.repos-rail`):

- **Repositorio** = Proyecto de software (ej. `dev-studio`, `harness-studio`).
- **Workspace** = instancia física en disco de ese repositorio (`git worktree`, ej.
  `wt/aislar-sesion-worktree` — el prefijo `wt/` de branch queda igual, es convención de infra
  aparte). **Workspace y Sesión son 1:1** (ya vigente por §5.1) — el sidebar
  y el rail de sesiones dentro de Dev Studio muestran el MISMO conjunto de cosas desde dos ejes
  (por disco / por a-quién-le-debo-atención), sincronizados visualmente (`data-session`
  compartido). Nunca son dos jerarquías independientes ni un workspace "contiene" sesiones
  ajenas.

**Torre: ELIMINADA — no es "F2 o esperar", es que no existe.** Chris pidió explícitamente
sacarla del shell. Esto resuelve la duda abierta de §9 sobre la Torre: no se clona, no hay vista
transversal multi-proyecto hoy. Si en el futuro hace falta un inbox transversal, es una historia
de backlog nueva a proponer (el mockup ya tiene el gancho: story "Panorama transversal de
sesiones entre repositorios" / capability `multi-repositorio`), no una resurrección de la Torre tal
cual estaba.

**Backlog y Configuración: overlays, no tabs.** Al entrar a un workspace se aterriza en
**Dev Studio** — única vista, ya no hay tabs de nivel superior. Su header tiene dos botones,
`📋 Backlog` y `⚙ Configuración`, que abren paneles superpuestos (`.panel-overlay`) sobre el área
de chat + panel lateral. El rail de sesiones (232px) y el sidebar Repositorios/Workspaces quedan
SIEMPRE visibles y clickeables detrás/al lado del overlay — decisión explícita de Chris, no
renegociable sin nueva ratificación.

- **Configuración** por ahora contiene solo la pantalla "Añadir rol" (§7). Es donde futuros
  ajustes de proyecto (Git, extensión, roles predeterminados) van a vivir cuando se construyan.

**Regla de negocio nueva: toda sesión liga sí o sí a un paquete de backlog.** Ya NO existe
"sesión libre sin historia" (corrige §4 Developer, paso 2 — esa opción queda OBSOLETA). Flujo
construido: `+ Nueva sesión` → abre Backlog en modo picker → clic en una historia → pasa a
Configuración con el paquete elegido → `Crear sesión aislada` crea la sesión (y su workspace,
1:1) ya ligada. Entrar a Configuración por `⚙` directo, sin pasar por Backlog, muestra un aviso
y redirige a elegir paquete si se intenta crear sin uno — nunca se permite una sesión huérfana.

## 1. Los tres referentes

| | AgentsRoom.dev | Conductor.build | Operon.app |
|---|---|---|---|
| Unidad organizadora | **Proyecto** = carpeta real en disco, agrupable por zona | **Workspace** = git worktree aislado | No verificable |
| Unidad de trabajo | Sesión de agente (chat + terminal nativo del proveedor) | Workspace = sesión (branch + worktree + chat + terminal + diff, todo en una unidad) | — |
| Multi-agente/proveedor | Sí — Claude, Codex, Antigravity, Mistral Vibe, Grok, OpenCode, Aider | Sí — Claude Code, Codex, Cursor | — |
| Aislamiento de sesión | Ninguno declarado (todas las sesiones comparten el mismo `cwd` del proyecto) | **Git worktree + branch por workspace** — aislamiento real en disco, setup script automático (~10s) | — |
| Diferenciador fuerte | Foco explícito en el orquestador: inbox transversal agrupado por estado, backlog desacoplado del agente, "equipos" de roles encadenados | Diff viewer conversacional (comentás inline, el agente resuelve) + checkpoints/revert por turno (git ref privado) + PR flow integrado con CI/Linear | — |
| Multi-usuario / organización | **No** — mono-usuario, local-first | **No** — mono-usuario, local-first | — |
| Modelo de negocio | Free (3 proyectos) + Pro ilimitado, BYOK | Free hoy, BYO API/suscripción, recién Serie A $22M | — |

`operonapp.dev`: sin presencia indexada verificable en esta pasada (WebSearch no encontró nada,
WebFetch bloqueado por red del sandbox). Si Chris tiene URL/captura directa, vale una vuelta
futura antes de firmar F0 — no bloquea esta spec porque agentsroom + conductor ya cubren el
patrón dominante de la categoría.

**Hallazgo central, el más importante de todo este documento:** ni AgentsRoom ni Conductor
resuelven el "modelo de organización" que nuestra propia `VISION.md` ya declara como norte —
**multi-usuario complementario, GitHub como conector**. Los dos referentes son herramientas
mono-usuario, local-first (un developer, muchos agentes en su propia laptop). Eso es intencional
en su segmento, pero es exactamente el vacío que Dev Studio quiere llenar distinto. No hay nada
que clonar ahí — es terreno propio. Ver §6.

## 2. Reencuadre — de "agentes" a "roles": la estrella es quien orquesta

AgentsRoom mezcla dos cosas bajo la palabra "agente": (a) la plantilla reutilizable (system
prompt + estilo + proveedor/modelo por defecto) y (b) la instancia corriendo. Separamos eso
limpio:

- **Rol** = la plantilla (quién es, qué sabe hacer, qué prompt trae, qué proveedor/modelo usa por
  defecto). Es durable, se versiona, vive en el roster del proyecto o en un catálogo curado.
- **Sesión** = una instancia de un rol *corriendo* en un proyecto: su propio chat, su propio
  terminal, su propio estado, opcionalmente su propio Workspace (git worktree/branch) aislado.

Evidencia concreta de que AgentsRoom ya diseña, aunque no lo nombre así, "el usuario como
estrella" (tomada de la exploración en vivo, no de la landing):

1. La vista por defecto **no** agrupa sesiones por conversación/timestamp — agrupa por
   **REQUIERE ENTRADA → POR REVISAR → ACTIVAS → INACTIVAS**. Es la pregunta de un manager
   ("¿a quién le debo atención ahora?"), no la de alguien chateando con un bot.
2. Badge global "N necesita intervención" (`Cmd+Shift+I`) — un inbox único, transversal a todo
   proyecto/rol/proveedor.
3. El Backlog desacopla "qué hay que hacer" de "quién lo hace": una tarea puede vivir sin rol
   asignado; el sistema auto-sugiere el rol más adecuado o el usuario la arrastra sobre cualquier
   sesión activa.
4. "Equipos" (roles encadenados con feedback loop, ej. Dev→QA) convierten al usuario en diseñador
   de un *proceso*, no en redactor de un prompt para un bot.
5. Revisión de cambios **filtrada por quién los hizo** + commit con "adjuntar la conversación del
   rol" — posiciona al usuario como auditor/integrador, nunca como coautor de un solo modelo.

Estos cinco puntos son el corazón a clonar. Todo lo demás (catálogo de 235 roles, Dynamic Island,
app móvil) es superficie — se puede o no copiar según §7.

## 3. Modelo de dominio propuesto

Dev Studio hoy (`internal/domain/session.go`) tiene un único concepto: `Session` con `Cwd` suelto,
sin contenedor. Es deliberadamente mínimo (F1 = esqueleto). La célula congelada
(`prenter-harness/products/devhub`) ya construyó, y verificó, un modelo de dominio mucho más rico
que responde exactamente a lo que este clon necesita — **la tarea de F2+ no es inventar estos
conceptos, es decidir en qué forma entran por la puerta de la experiencia**.

| Concepto Dev Studio | Definición | Ya existe en la herencia como… | Estado |
|---|---|---|---|
| **Proyecto** (= repositorio) | Carpeta real en disco / repo git. Contiene sesiones, roster de roles, backlog, capabilities, ajustes | `registryProject` / `sysWorkspace` (`registry.go`, `workspace.go`) — clave compuesta `"{proyecto}/{sistema}"` | Existe en herencia, no portado |
| **Sesión** | Instancia de un Rol corriendo en un Proyecto; opcionalmente con un Workspace (git worktree/branch) aislado (ver §5) | `ActiveSession` (lock de filesystem) + `SessionService` (dev-studio, en memoria, ya verificado con 2 procesos concurrentes) | Parcial: el runtime ya existe en dev-studio; falta el contenedor Proyecto y el aislamiento por workspace |
| **Rol** | Plantilla: nombre, system prompt, estilo, proveedor/modelo por defecto | `AgentDefinition` + `functional_areas[]` (`types.ts:496`) | Existe en herencia, no portado |
| **Historia de usuario** | Unidad de backlog; referencia ≥1 Capability que incrementa | `Story` (`types.ts:205`) — 10 estados (`idea→refining→refined→ready→developing→developed→reviewing→done→parked→dropped`) | Existe completo en herencia, no portado |
| **Capability** | Unidad de capacidad del sistema: módulo, status (live/beta/deprecated/sunset), dueño, changelog | `Capability` (`types.ts:311`) | Existe completo en herencia, no portado |
| **Proceso as code** | Estados/transiciones/gates/dueños-por-rol que gobiernan cómo se mueve una historia | `ProcesoDescriptor` (Go `proceso.go:96`, contrato I-77) — `Duenos map[rol][]binding`, gates con autoridad | Existe completo en herencia, no portado — decisión vigente que la épica no reabre |
| **Torre / Inbox orquestador** | Vista transversal de TODAS las sesiones de TODOS los proyectos abiertos, agrupada por estado | `TorreProyecto`/`TorreSistema` (veredictos 2 ejes) — más cercano en espíritu a lo que agentsroom hace con su inbox global | Existe en herencia (otro propósito: salud), se puede re-derivar |
| **Organización / multi-usuario** | Varios humanos complementarios (CTO·dev·devops·PO) trabajando a la par, GitHub como conector | **No existe en ningún lado** (ni herencia, ni agentsroom, ni conductor) | TBD real — terreno propio, ver §6 |

Esto responde directamente al pedido de F0 de "auditoría de herencia (orden del port)": **el orden
natural es Proyecto → Rol → Sesión-con-workspace → Historia/Capability → Proceso-as-code**, porque
cada pieza depende de la anterior para tener sentido en pantalla, y es también el orden en que
agentsroom construye su propio onboarding (proyecto primero, agente después, tarea al final).

### Propuesta concreta para los slots `__FILL_ME__` de `project.config.yaml`

```yaml
domain_modules: [proyecto, sesion, rol, historia, capability, proceso, organizacion(TBD)]

agent_roster:
  builder: [full-stack, frontend, backend, movil, devops]
  auditor: [qa, seguridad, arquitecto-software]
  humano-complementario: [product-owner, ideacion]
  # catálogo extendido opcional, curado — NO el volcado sin filtrar de 235 roles de "The Agency"

value_stream: [idea, refining, refined, ready, developing, developed, reviewing, done]
  # reusa 1:1 el modelo de 10 estados de Story en la herencia — no reinventar

wip_caps:
  developing_por_proyecto: N   # evita que un usuario abra 20 sesiones sin revisar ninguna
  reviewing_por_proyecto: M
  # mapea a procesoEstado.Wip, ya existe como campo en la herencia
```

Estos son propuestas de esta spec, no están firmados — quedan para que Chris los ratifique en F0
junto con el resto.

## 4. Journeys por rol (lo que F0 pide explícitamente)

> Formato: pasos numerados, sin Gherkin todavía (RONDA 1 · funcional-primero). Cada journey asume
> que el Proyecto ya existe.

### Developer

1. Abre Dev Studio → elige un Repositorio y un Workspace en el sidebar (o crea uno nuevo) → entra
   directo al rail de sesiones de ese workspace (ver §0.5 — ya no hay Torre/inbox global).
2. Abre el overlay **Backlog** → toma una historia en `ready` — toda sesión liga sí o sí a una
   historia (§0.5; ya NO se puede crear una sesión libre sin historia).
3. El sistema crea una **Sesión**: asigna/sugiere un Rol (Full-Stack/Backend/Frontend), crea un
   Workspace (git worktree + branch aislados) para esa historia (patrón Conductor, ver §5).
4. Chatea con el rol, corre comandos en el terminal embebido, ve el diff en vivo filtrado por lo
   que ese rol tocó.
5. Cuando el rol termina: la historia pasa `developing → developed`; el developer revisa el diff
   (comentario inline, el rol resuelve — patrón Conductor), y dispara el gate de PR.
6. El proceso as code decide el siguiente dueño del gate (auditor/QA humano o rol QA) — la
   historia pasa a `reviewing`.

### Product Owner

1. Entra al Proyecto → tab **Capabilities** → ve el mapa de capacidades vivas/beta/deprecadas.
2. Crea una historia nueva ligada a una Capability existente (o crea la Capability si es
   territorio nuevo) → historia nace en `idea`.
3. Usa el rol **Ideación/PM** para refinarla conversacionalmente (system prompt orientado a
   Gherkin/criterios de aceptación, no a código) → historia pasa `refining → refined`.
4. Marca `ready` cuando el equipo puede tomarla — el gate de "refined→ready" en el descriptor de
   proceso exige dueño PO explícito.

### DevOps

1. Entra al Proyecto → tab **Comandos** (procesos persistentes: dev server, build, watchers) +
   **Túneles de vista previa** (exponer localhost) + **Conexiones SSH** (remoto).
2. Orquesta la infraestructura que las sesiones de developer necesitan para correr y previsualizar
   sin salir de la app.
3. Es dueño de los gates de "listo para deploy" en el descriptor de proceso — el sistema no
   permite `developed → done` sin su aprobación si el descriptor lo exige.

### CTO

1. **Hueco sin resolver** (la Torre que cubría esto se eliminó del shell, ver §0.5): necesita una
   vista multi-proyecto/multi-sistema completa (no solo "mis sesiones" — todas las de la
   organización). Ninguna pantalla del mockup actual lo cubre; queda pendiente de diseño para
   cuando exista el modelo multi-usuario (§6), no antes.
2. Usa el mapa de **Capabilities** a nivel portafolio (parent_cap, dueños, status) para ver salud
   y deuda técnica del conjunto de sistemas.
3. Revisa drift/arquitectura (lente `meta.clase` de la herencia) antes de aprobar gates
   arquitectónicos.

## 5. Qué tomamos de Conductor.build (diferenciadores a clonar)

1. **Sesión = workspace + branch aislados, no solo un `cwd` compartido.** Esto no es solo un
   feature — **resuelve directamente la brecha de seguridad ya documentada en este repo**
   (`arch/boundaries/sesion-aislada-por-cwd.md`: hoy cada sesión tiene su propio `Cwd` pero sin
   validar que no sea `$HOME`/`~/.ssh`; y todas comparten el mismo repo real, sin aislamiento
   entre sí). Adoptar este patrón —nuestro **Workspace**— convierte una debilidad conocida en un
   feature de producto. Candidato fuerte para **DH-14** (la ficha que el LEDGER ya marca como
   siguiente).
2. **Diff viewer conversacional** — no solo mostrar el diff, dejar comentar inline y que el rol lo
   resuelva. Es el feature más valorado en reviews externas de Conductor.
3. **Checkpoints/revert por turno** (git ref privado, separado del historial real) — red de
   seguridad que ningún competidor mono-usuario más ofrece con esa granularidad.
4. **Setup script por proyecto** (`.conductor/settings.toml` equivalente) — automatiza
   `pnpm install`/migraciones/`.env` al crear una sesión nueva, no fricción manual cada vez.
5. **PR flow integrado con CI** — crear PR, seguir checks, responder comments, avisar cuándo está
   listo para mergear. Encaja 1:1 con nuestro gate de proceso as code (`developed→reviewing→done`).

## 6. Lo que NO clonamos — o clonamos distinto (nuestro espíritu)

- **Multi-usuario real, no mono-usuario local-first.** Ni agentsroom ni conductor resuelven
  "varios humanos complementarios trabajando a la par sobre uno o varios sistemas" — es el norte
  que `VISION.md` ya declara y ninguno de los dos referentes lo tiene. GitHub como conector entre
  usuarios (no solo como destino de un PR) es terreno de diseño propio, no de clon. Sigue TBD —
  **F0 debe decidir qué sincroniza el repo entre usuarios** (¿el backlog vive en GitHub Issues?
  ¿las sesiones son visibles para otros miembros de la organización en tiempo real?).
- **Capability explícita y trazable, no un kanban suelto.** AgentsRoom tiene Backlog/tareas pero
  **no** tiene el concepto de "capacidad del sistema que una historia incrementa" — cada tarea es
  una entidad aislada. Nuestra Historia siempre referencia una Capability; eso da trazabilidad de
  producto (qué construimos incrementó qué) que ningún referente ofrece.
- **Proceso as code gobernando gates por rol**, no un kanban de 4 columnas fijas. El
  `ProcesoDescriptor` (10 estados, transiciones con verbo CDEvents, gates con autoridad y dueño
  por rol) es más riguroso que el Backlog TODO/IN PROGRESS/PENDING/DONE de agentsroom — lo
  mantenemos como decisión vigente, no lo cambiamos por el modelo más simple del referente.
- **Catálogo de roles curado, no 235 sin filtrar.** El catálogo "The Agency" de agentsroom (235
  roles importados de un repo open-source MIT) es volumen sin curaduría — genera ruido de
  decisión ("¿qué rol elijo entre 235?"). Preferimos un roster corto y curado por proyecto
  (`agent_roster` arriba) con la opción de agregar un "Rol personalizado" — como hace agentsroom
  también, pero sin el catálogo gigante por defecto.
- **BYO CLI-nativo (decisión DH-10 vigente), no Agent SDK con API key propia.** Conductor y
  agentsroom en su mayoría asumen BYOK de cada proveedor; nuestra decisión ya tomada (spawnear el
  `claude` que el usuario ya tiene instalado, nunca tocar credenciales) se mantiene — no la
  reabre esta spec.
- **Sin Dynamic Island / compañero móvil por ahora** — features de "no romper el ciclo de
  atención al alejarse del escritorio", coherentes con su producto mono-usuario de escritorio,
  pero no hay señal de que sea prioritario para nuestro público (equipos, no un solo developer
  ansioso por su Mac). Se puede reconsiderar en FN si el modelo multi-usuario lo pide distinto
  (ej. notificaciones de organización, no de un solo dispositivo).

## 7. Inventario de pantallas propuesto (arquitectura de información)

Derivado del inventario real de agentsroom (§1 del research), renombrado a nuestra terminología y
con Capability agregada donde agentsroom no la tiene.

| Pantalla | Origen (agentsroom) | Cambio respecto al original |
|---|---|---|
| ~~Torre~~ (inbox global, multi-proyecto) | Topbar con avatares de todos los agentes de todos los proyectos abiertos | **ELIMINADA** — ver §0.5, ratificado 2026-07-06. No se clona. |
| **Dev Studio** (antes "Vista Proyecto") — se entra eligiendo un Workspace en el sidebar Repositorios→Workspaces (§0.5); sidebar sesiones agrupadas por estado + tabs | Vista principal 3 paneles | Ya no es una tab entre 4 — es la única vista; agrega tab **Capabilities** que agentsroom no tiene (pendiente de construir) |
| **Vista Sesión** — chat + terminal / archivos + cambios + pruebas | Ídem | Agrega diff conversacional + checkpoints (Conductor) |
| **Backlog** (kanban) | Kanban TODO/IN PROGRESS/PENDING/DONE | Sigue siendo kanban de 10 estados de `Story` (cada tarjeta muestra su Capability), pero ahora es **overlay** sobre Dev Studio, no tab (§0.5) |
| **Detalle de Capability** | No existe en agentsroom | Nuevo — status, dueño (rol), `parent_cap`, historias vinculadas |
| **Añadir rol** | Selector rol (3 capas) + motor IA | Vive dentro del overlay **Configuración** (§0.5), no como modal/tab suelto. Igual, catálogo curado (§6) |
| **Detalle de rol** (system prompt) | Ídem | Igual |
| **Ajustes de proyecto** — roles predeterminados, Git, extensión | Ídem | Vivirá dentro del overlay **Configuración** (§0.5), junto a Añadir rol, cuando se construya |
| **Tareas programadas / Comandos / Túneles / SSH** | Ídem | Igual — encajan directo en journey DevOps (§4) |
| **Biblioteca de prompts / skills** | Ídem | Igual, exportable a Claude/Cursor/Codex como ya hacen ellos |
| **Proceso** (visor del descriptor) | No existe en agentsroom | Nuevo — de la herencia (`ProcesoProvider`), visualiza estados/gates/dueños |

## 8. Fases propuestas (diálogo directo con los candidatos de `NORTE-BORRADOR.md`)

- **F0** (sin cerrar): además de journey-por-rol + principios, firmar explícitamente **el orden
  de herencia de §3** (Proyecto → Rol → Sesión-workspace → Historia/Capability → Proceso-as-code) y
  dejar abierta, con dueño y fecha, la pregunta de organización/multi-usuario (§6) — no bloquea F1
  ya entregado, pero si no se firma antes de F2 cada rebanada entra sin norte.
- **F1**: ya entregado (DH-13) — esqueleto + driver CLI-nativo + rail multisesión.
- **DH-14** (ya anotada como siguiente en el LEDGER): cerrar `sesion-aislada-por-cwd` — esta spec
  propone resolverla adoptando **workspace aislado por sesión** (§5.1) en lugar de solo validar
  rutas prohibidas; mata dos pájaros (seguridad + feature de producto) de un tiro.
- **F2+**: port por rebanadas en el orden de §3, cada una con su propia spec congelada y
  verificación real (browser/app), como ya manda la épica.
- **FN**: instalable v1 — sin cambios respecto al candidato ya escrito.

## 9. Dudas abiertas (para Chris, antes de firmar F0)

- [ ] ¿El modelo de organización multi-usuario (§6) es un requisito de F0/F2 temprano, o puede
  vivir como TBD hasta más adelante sin bloquear el port de Proyecto/Rol/Historia/Capability
  (que sí tienen sentido en versión mono-usuario mientras tanto)?
- [ ] ¿Confirmamos workspace aislado por sesión como la forma de cerrar DH-14, en vez de (o además
  de) la validación de rutas protegidas ya planteada?
- [ ] ¿El catálogo de roles curado (§6) arranca con la lista corta propuesta en §3, o Chris quiere
  ver primero el system-prompt completo de algún rol de "The Agency" antes de decidir cuánto
  catálogo externo vale la pena importar?
- [ ] `operonapp.dev` — ¿Chris tiene una URL/captura directa? No se encontró indexado; si hay algo
  puntual que rescatar de ahí, vale una vuelta antes de firmar F0.
- [x] ~~¿La Torre (inbox global) es F2 temprano o puede esperar...?~~ RESUELTO 2026-07-06:
  eliminada del shell, ver §0.5. No se clona.

## 10. Composer y tab Cambios (detalle — ronda 2 de exploración, 2026-07-06)

> Pasada específica sobre dos zonas de `agentsroom.dev/demo/index.html?lang=es` (URL directa al
> workspace, sin landing) a pedido de Chris, con click/hover real en cada elemento — no solo
> inventario visual. Fuente: exploración Chrome DevTools MCP, click-by-click.

### 10.1 Composer — los 7 botones del input de mensaje

Fila de iconos circulares bajo el textarea, agrupados por color/función. Adoptamos los 7,
renombrados a español neutro, con una decisión de producto por grupo:

| # | Original | Grupo/color | Qué hace en agentsroom | Nombre Dev Studio | Decisión |
|---|---|---|---|---|---|
| 1 | Boceto | naranja | Modal tipo Excalidraw completo (formas, color, texto) → "Adjuntar al mensaje" | **Boceto** | Clonar — es genuinamente útil para explicarle algo visual a un rol sin escribir un párrafo |
| 2 | Captura | naranja | Popover con atajo global `Ctrl+Shift+2` + "Capturar ahora" (requiere permiso OS) | **Captura de pantalla** | Clonar el popover; la captura real depende de permisos de escritorio (fuera del alcance del mockup) |
| 3 | HTML `</>` | violeta | Popover: capturar una página como HTML vía extensión de Chrome propia | **Importar página web** | Clonar en espíritu, sin comprometernos a una extensión de Chrome propia en F0 — nota como TBD de alcance |
| 4 | Móvil | violeta | Popover con QR + mockups de share-sheet iOS/Android | **Enviar desde el móvil** | Igual — depende de que exista app móvil (no está en el alcance de F1-F2) |
| 5 | Prompts (biblioteca) | violeta | Modal grande: tabs Proyecto/Global, carpetas, estado vacío "+ Crear el primer prompt" | **Biblioteca de prompts** | Clonar completo — encaja con roles reutilizables (§3), sin dependencia de infraestructura nueva |
| 6 | Modo voz | rosa | Modal: "no disponible en la demo, requiere backend" | **Conversación por voz** | Marcar como no-F1/F2 explícito — mismo patrón: mensaje claro de "todavía no", no ocultar el botón |
| 7 | Dictar | rosa | Igual patrón, mic → transcripción | **Dictar mensaje** | Igual — distinto de "Modo voz" (unidireccional vs bidireccional) |

Detalles de comportamiento a preservar: botón Enviar **deshabilitado** hasta que hay texto; al
escribir aparecen dos acciones — **Enviar** (avión) y **Añadir a la cola** (⇒, encola el mensaje
mientras el rol sigue trabajando en otra cosa — útil para no interrumpir un turno en curso,
coherente con el boundary `TestOneTurnAtATime` que ya tenemos). El composer es redimensionable
(drag del borde superior) y expandible (botón "↗").

Debajo del composer hay una fila secundaria: tab **Terminal** (cerrable) + **"+"** (nueva
terminal) + chip **Comandos** (procesos persistentes: dev server/build/watchers, con detección
automática por IA y toggle "reabrir al iniciar") + chip **SSH remoto** (conexiones a máquinas
externas) + 2 iconos de layout (vista dividida / separar en ventana flotante). Esto mapea 1:1 al
journey DevOps de §4 — no es superficie nueva, es la forma concreta que ya anticipamos ahí.

### 10.2 Tab Cambios — por qué es "más visual" y qué clonar

Confirmado: el tab Cambios de agentsroom tiene dos sub-tabs, **Árbol de trabajo** (lista plana de
archivos, no árbol de carpetas real pese al nombre) e **Historial** (timeline tipo git-log). Lo
que lo hace más visual que un `git diff` de terminal:

1. **Chips de filtro por rol** ("Todos 3" / "Full-Stack 2" / "QA 1") — filtra la lista de archivos
   por qué rol los tocó. Encaja directo con nuestro modelo Sesión↔Rol (§3).
2. **Selección con checkboxes + commit por lote** — el footer muestra "Se hará commit de X/Y
   archivos de [rol]" y tiene un checkbox **"Adjuntar la conversación del rol al commit"** (sube
   la conversación como link no listado en el mensaje de commit) — clonamos esto tal cual, es la
   forma concreta del "commit con contexto de conversación" que ya proponíamos en §5.
3. **Dos niveles de diff**: un "quick look" (click directo en el archivo → modal unificado con
   syntax highlighting, sin números de línea) y un flujo **"Revisar cambios"** completo
   (side-by-side ORIGINAL/MODIFICADO, numeración real por columna, contador "N/3 revisados",
   navegación Siguiente/Finalizar, pantalla de cierre "Listo para hacer commit"). Clonamos ambos
   niveles — el quick look para un vistazo rápido, el flujo completo para revisión seria antes de
   aprobar el trabajo de un rol.
4. **Límite real encontrado, no lo replicamos como si fuera nuestro):** ni el diff simple ni el
   diff side-by-side de agentsroom permiten comentar sobre una línea específica (confirmado
   probando click/hover en varias líneas — no pasa nada). El **diff conversacional** que
   propusimos adoptar de Conductor (§5.2) sigue siendo un diferenciador nuestro genuino frente a
   agentsroom, no algo que ya exista ahí — vale la pena remarcarlo para no perder la ambición al
   portar esto en F2+.

Nombres adaptados: "Árbol de trabajo" → **Cambios pendientes** (es más preciso, no hay árbol de
carpetas real); "Revisar cambios" se mantiene igual; "Commit" / "Commit + Cerrar" →
**Confirmar cambios** / **Confirmar y cerrar sesión**.
