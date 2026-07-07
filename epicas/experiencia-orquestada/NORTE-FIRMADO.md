# Épica «Experiencia Orquestada» — NORTE FIRMADO (F0 cerrada 2026-07-07 · DH-14)

> Firmado por Chris en la sesión F0 (2026-07-07), fork por fork vía AskUserQuestion — ficha
> `DH-14` en [`LEDGER.md`](../../LEDGER.md). Este archivo ES el borrador renombrado
> (`NORTE-BORRADOR.md` → `NORTE-FIRMADO.md`). Carpeta temporal por diseño (patrón Torre de
> Control): al cierre de la ÉPICA (no de F0) se borra; lo permanente se promueve a `specs/`.
> Este norte gobierna F2+ hasta entonces.

## Cruda del operador (2026-07-04, verbatim)

"Quiero crear una Épica para trabajar la interfaz de esta aplicación e ir trayendo lo que
hemos desarrollado de a pocos conforme a cómo hemos decidido será la experiencia del
developer/usuario. El objetivo de esta aplicación es que el usuario (puede ser CTO, developer,
devops, product owner) trabaje en esta aplicación dentro de una organización de forma
'orquestada', con otros usuarios complementarios trabajando a la par uno o varios sistemas a
la par, apalancándose de un repositorio (Github) como conector entre todos."

## Cruda del reframe (2026-07-07, verbatim — la decisión más grande de F0)

"los 'roles', el 'proceso', en sí vienen de los 'arneses' que estamos creando en ArnesIA. Es
decir, al instalar la aplicación me conectaré a un repositorio de plugins y desde allí en la
vista 'Roles' miraré cuáles hay (Todos son complementarios) y añadiré los que me correspondan
o los que mi jefe me ha permitido de acuerdo a mi puesto. Y esto es por Proyecto y de acuerdo
a eso se pinta el backlog y con esos arneses son los que trabajará el claude code en ese
proyecto específico."

## Qué gobierna esta épica (ratificado, sin cambios)

1. **La experiencia manda el port.** Ninguna pieza del monorepo entra "porque existe": entra
   cuando el journey del rol la pide, en la forma que el journey pide.
2. **Inventario de herencia disponible** (monorepo congelado): torre de control 2 ejes ·
   descriptor de proceso + board gobernado (ProcesoProvider, 13 verbos CDEvents) · cockpit de
   delivery (sesiones de agente parametrizadas) · lente arquitectura (`meta.clase`) ·
   multi-workspace · 27 endpoints Go · UI Next v0.6.x.
3. **Decisiones vigentes que la épica NO reabre** (salvo ficha): driver CLI-nativo (DH-10) ·
   descriptor I-77 · Go + UI embebida · shell y navegación §0.5 (spec clon, ratificado
   2026-07-06).

## Modelo Rol/Registry — FIRMADO (el reframe)

- **Rol = arnés (plugin) instalado desde un registry.** DevStudio no crea ni almacena roles:
  consume e interpreta. Un rol custom = publicarlo al registry de tu organización. Cero roles
  locales — un solo origen, cero drift. Todos los roles son complementarios.
- **Registry PROPIO de DevStudio** (fork firmado — se descartó usar el sistema de plugins
  nativo de Claude Code como mecanismo): formato, repositorio y conexión propios. Nota técnica
  vigente: el payload del arnés igual debe materializarse en el proyecto como artefactos que el
  `claude` spawneado cargue (DH-10 no se toca). Diseño del mecanismo = PB-25.
- **Instalación por Proyecto.** Roster del proyecto = arneses instalados. De ahí se pinta el
  backlog (el descriptor I-77 viaja EN el arnés; DevStudio lo interpreta — ya decidido en
  DH-12, hoy extendido: no solo el proceso, el ROL ENTERO viene del arnés) y con esos arneses
  trabaja el `claude` de cada sesión de ese proyecto.
- **Permisos por puesto** («lo que mi jefe me habilita») = **con multi-usuario** (PB-21→PB-23).
  v1 mono-usuario instala libre del registry.

## Journeys por rol — FIRMADOS

Base: spec clon §4, con el reframe integrado. **Paso previo común a todos:** en Configuración
del proyecto se conecta el registry y se instalan los roles — eso ES el roster del proyecto.

- **Developer** — firmado tal cual (6 pasos §4): sidebar Repositorio→Workspace → overlay
  Backlog → toma historia `ready` (toda sesión liga sí o sí a una historia, §0.5) → sesión con
  rol sugerido/asignado DE LOS INSTALADOS + workspace aislado (worktree + branch, patrón
  Conductor) → chat/terminal/diff filtrado por rol → gates del proceso
  (`developing→developed→reviewing`).
- **Product Owner** — firmado tal cual (4 pasos §4): mapa Capabilities → historia nueva ligada
  a capability (o capability nueva) → refina con rol instalado (Ideación/PM) → `ready` con gate
  de dueño PO explícito. Forma exacta de la vista Capabilities (tab interna vs overlay) se
  decide en PB-15, no acá.
- **DevOps** — firmado tal cual (3 pasos §4): comandos persistentes + túneles de vista previa +
  SSH (forma concreta ya mapeada: §10.1, chips bajo el composer) · dueño de los gates de deploy
  según el descriptor del arnés. Materialización = PB-17.
- **CTO** — firmado como **HUECO DECLARADO**: pasos 2-3 válidos (mapa Capabilities a nivel
  portafolio + lente arquitectura antes de gates arquitectónicos); el paso 1 (vista
  multi-proyecto de la organización) NO tiene pantalla desde que la Torre se eliminó (§0.5) y
  queda bloqueado por PB-21→PB-22 — se diseña recién cuando exista el modelo multi-usuario.

## Principios de la experiencia — FIRMADOS (7)

1. **La estrella es quien orquesta, no el rol.** Toda vista responde «¿a quién le debo atención
   ahora?» — agrupación por estado de atención (REQUIERE ENTRADA → POR REVISAR → ACTIVAS →
   INACTIVAS), nunca por timestamp de chat. *(spec §2)*
2. **La experiencia manda el port.** Ninguna pieza de la herencia entra «porque existe»; entra
   cuando el journey la pide, en la forma que la pide. *(regla 1 de la épica)*
3. **Todo deriva del dato as-code.** Proceso (I-77), arquitectura, documentación y el rol
   entero (arnés del registry): la consola interpreta, cero hardcode del ciclo. *(VISION +
   reframe 2026-07-07)*
4. **Nada huérfano: sesión→historia→capability.** Toda sesión liga historia; toda historia
   incrementa capability. Trazabilidad de punta a punta — el diferenciador que ningún referente
   tiene. *(§0.5 + §6)*
5. **Aislamiento por defecto.** Sesión = workspace propio (worktree + branch); la seguridad es
   feature, no parche. *(§5.1 → PB-02)*
6. **El repositorio es el conector.** Código, proceso, arquitectura y docs convergen en el repo
   (GitHub); lo multi-usuario sincroniza por ahí, no por un backend paralelo. *(VISION)*
7. **Honestidad de superficie.** Lo que no existe se dice («todavía no disponible»), no se
   esconde; brechas declaradas en la UI igual que en INCREMENTO. *(patrón §10.1 + regla de la
   casa)*

## Orden de herencia del port — FIRMADO

**Proyecto → Rol → Sesión-workspace → Historia/Capability → Proceso as code** (spec §3 tal
cual; ratifica el orden del backlog PB-05→06→07→08). Dos notas:

- **PB-02 (workspace aislado) corre ADELANTADA** en banda 🔴 por ser brecha de seguridad — se
  reconoce fuera de secuencia, no reordena la rebanada.
- **La rebanada Rol creció con el reframe**: incluye el registry propio (PB-25 precede o
  acompaña a PB-06).

## Dudas §9 de la spec — resueltas en F0

| Duda | Decisión (2026-07-07) |
|---|---|
| ¿Multi-usuario temprano o TBD? | **TBD formal SIN bloquear el port** — las rebanadas se portan mono-usuario; PB-21 queda ⚪ con dueño Chris, a decidir antes de PB-22 (vista CTO) y PB-23 (auth) |
| ¿Workspace aislado cierra DH-14/PB-02? | **Sí, Y ADEMÁS validación de rutas protegidas** ($HOME, ~/.ssh…) en la misma entrega — el worktree aísla sesiones entre sí pero no impide registrar un proyecto en ruta peligrosa |
| Catálogo de roles curado | **Disuelta por el reframe**: sin catálogo local ni import de «The Agency» — la curaduría vive en el registry |
| operonapp.dev | **Cerrada sin revisar** — agentsroom + conductor cubren el patrón; no condiciona el norte |
| Torre | Ya estaba resuelta (§0.5): eliminada |

## Slots del seam (project.config.yaml) — ratificados

`domain_modules` = spec §3 **+ registry** · `agent_roster` ajustado al repo real
(backend-go/frontend-web — no el aspiracional con móvil) · `value_stream` = **10 estados
completos** (incluye `parked`/`dropped`) · `wip_caps` = **advisory: se visualizan, no
restringen**. Cada slot cita esta firma en el archivo.

## Fases (estado real)

- **F0 · Norte firmado** — CERRADA hoy (esta firma, DH-14).
- **F1 · Esqueleto de la app** — ENTREGADA (DH-13): CAP-01..06.
- **F2+ · Port por rebanadas** — orden firmado arriba; cada rebanada con spec congelada,
  checkpoint de forks y verificación real (browser/app).
- **FN · Instalable v1** — sin cambios (PB-24; ⚠ Consumer Terms antes de vender).

## Qué quedó explícitamente FUERA de esta firma

- Modelo de organización/multi-usuario — qué sincroniza el repo entre usuarios: TBD PB-21,
  dueño Chris.
- Vista CTO multi-proyecto (paso 1 del journey): bloqueada PB-21→PB-22.
- Permisos por puesto: llegan con multi-usuario (PB-23).
- Mecanismo interno del registry propio (formato, distribución, versionado): rebanada PB-25.
- Forma de la vista Capabilities (tab vs overlay): PB-15.
- Números de wip_caps como límite duro: hoy solo visualización.
- Promoción de lo permanente a `specs/`: al cierre de la ÉPICA, no de F0.

## Reglas de la casa que aplican (sin cambios)

SPEC congelada por fase en `specs/` · forks por AskUserQuestion (checkpoint 60s → recomendada)
· fichas DH-NN por fase · verificación real (browser/app) antes de cerrar fase · una
conversación = este repo (no arrastrar contexto del monorepo: lo que haga falta, se cita por
ficha/spec).
