# BACKLOG de producto — DevStudio

> **Plano: PRODUCTO.** SSoT del ORDEN de todo lo que queremos que la aplicación haga y aún no
> hace. El detalle de cada ítem vive en su fuente (spec/épica/boundary — citada por fila); este
> archivo no lo duplica. Lo ya construido y verificado vive en [`INCREMENTO.md`](./INCREMENTO.md)
> (el otro extremo del par). Norte: [`VISION.md`](./VISION.md) · decisiones: [`LEDGER.md`](./LEDGER.md).
>
> Creado 2026-07-07 — auditoría: el trabajo pendiente vivía desparramado en LEDGER §Siguiente,
> `NORTE-BORRADOR.md`, `CLON-AGENTSROOM-SPEC.md` §§5-10, VISION §TBD y `arch/boundaries/`.
> Acá se consolida. **La prioridad de abajo es PROPUESTA sin firmar** — Chris la firma moviendo
> ítems de banda (regla 2).

## Reglas de uso (el contrato para nacer ordenados)

1. **Toda funcionalidad deseada nace acá** como ficha `PB-NN` con su fuente citada. Una idea
   suelta en conversación se registra ANTES de trabajarse — si no está acá, no existe.
2. **La prioridad la firma Chris.** Claude propone banda; mover un ítem de banda es decisión de
   Chris y queda anotada (fecha en la columna Estado). Bandas: 🔴 Ahora · 🟡 Próximo · 🟢 Después ·
   🧱 Deuda declarada · ⚪ TBD/bloqueado.
3. **El detalle vive en la fuente, el orden vive acá.** Este archivo nunca duplica specs.
4. **Toda PB incrementa una capability** (existente `CAP-NN` de INCREMENTO.md, o `nueva`) —
   trazabilidad historia→capability, la misma regla que el producto le impone a sus usuarios.
5. **Entrar a construcción:** la PB se refina por el proceso del arnés (spec + ficha `DH-NN`);
   la fila anota `en-curso · DH-NN`.
6. **Cerrar:** verificación REAL en vivo → la fila pasa a `entregada AAAA-MM-DD` → la capability
   se alta/actualiza en INCREMENTO.md **en el mismo commit**. Sin evidencia no hay entrega.
7. **Ritual de sesión:** "¿qué viene?" / "backlog" → se lee la banda 🔴 y se propone el
   siguiente paso. Al cerrar una sesión que tocó producto, este archivo refleja la realidad.

## 🔴 Ahora

| ID | Ítem | Incrementa | Fuente (detalle) | Estado |
|---|---|---|---|---|
| PB-01 | **F0: firmar el norte de la épica** — journeys por rol, principios, modelo de organización, orden de herencia, dudas abiertas §9. Gate de definición: deuda reconocida en DH-13 (F1 entregó código sin norte firmado) | — (gate) | `epicas/experiencia-orquestada/NORTE-FIRMADO.md` + `CLON-AGENTSROOM-SPEC.md` §8-9 | **entregada 2026-07-07** (gate — sin CAP) · DH-14 |
| PB-04 | **Shell PRENTER** (movida de 🟡, 2026-07-07: «cerremos el shell primero»): rail Repositorios→Workspaces (workspace=sesión 1:1 — el rail F1 murió) + studio-nav + overlays + **panel Cambios git REAL** (commit por pathspec; sin push/pull — boundary nuevo) + tokens PRENTER + Storybook (absorbió PB-03) | CAP-01/03/04/05/06 + CAP-07/08/09 nuevas | `specs/shell/SPEC.md` (CONGELADA, gate AC-1..10) | **entregada 2026-07-07** (24/24 checks en vivo) · DH-15 |
| PB-02 | **Workspace aislado por sesión** (worktree + `wt/{slug}` en `~/.dev-studio/workspaces/`) **+ rutas protegidas** + cierre con modal conservar/borrar — brecha `sesion-aislada-por-cwd` CERRADA (boundary enforced; resta token de capacidad, check TBD) | CAP-02 → **CAP-10** | `specs/workspace-aislado/SPEC.md` (CONGELADA, gate 14/14) | **entregada 2026-07-07** · DH-16 |
| PB-26 | **Instalable dogfooding + updater rebuild-local** — install.sh + .desktop + icono + versión embebida + «Actualizar» (rebuild + restart + reload). La app se actualizó a sí misma en el gate | **CAP-11** | `specs/workspace-aislado/SPEC.md` §3 | **entregada 2026-07-07** · DH-16 |
| PB-03 | **Sistema de diseño propio: Storybook + atomic design** | CAP-06 | premisa operador 2026-07-07 | **entregada 2026-07-07 vía PB-04** (fusionada): Storybook + átomos PRENTER — RN-9 · DH-15 |

| PB-27 | **Flujo «Nuevo Workspace» v2** — wizard ubicación→propósito: exploración read-only (`--permission-mode plan`) · trabajo SIEMPRE en worktree ligado a ítem · taxonomía estándar 5 tipos con branch por prefijo (murió `wt/`) · crear ítem desde el wizard. RN-1 evolucionada | CAP-09 + **CAP-12** | `specs/nuevo-workspace/SPEC.md` (CONGELADA, 10/10) | **entregada 2026-07-07** · DH-17 |

## 🟡 Próximo — shell + port por rebanadas (orden de herencia §3 de la spec clon)

| ID | Ítem | Incrementa | Fuente (detalle) | Estado |
|---|---|---|---|---|
| PB-05 | **Proyecto/Repositorio** como contenedor (rebanada 1 del port) — vaciada por el reframe DH-14: su único resto con journey (roster por proyecto + conexión registry) ES PB-25 | — | spec clon §3 · spec shell §5 | **fusionada en PB-25 2026-07-07** (firma Chris); metadatos declarados sin journey · DH-18 |
| PB-25 | **Registry propio de arneses** (⊕ restos PB-05) — DevStudio ADOPTA el estándar ArnesIA (nomenclatura-arnes v1 + marketplace git en producción): conexión con el git del usuario, instalación por proyecto = lock as-code `.devstudio/arneses.yaml` + caché forma-plugin, inyección por sesión vía `--plugin-dir`+system-prompt (DH-10 intacto). Roster mock MUERTO | **CAP-13** + CAP-02 | `specs/registry-arneses/SPEC.md` (CONGELADA, 7/7 en vivo con el arnés real de ArnesIA) | **entregada 2026-07-07** · DH-18 |
| PB-06 | **Rol** — vista Roles rica conectada al registry: explorar catálogo con detalle/versiones/upgrade (la instalación básica ya vive en Config, DH-18). Depende de PB-25 ✓ | nueva | NORTE-FIRMADO §Modelo Rol/Registry · spec registry-arneses §6 | propuesta |
| PB-07 | **Historia + Capability en la app** — kanban 10 estados como overlay, tarjeta muestra su capability, regla dura sesión↔historia (picker de `+ Nueva sesión`) | nueva | spec clon §0.5 + §3 (herencia `Story`/`Capability`) | propuesta |
| PB-08 | **Proceso as code** — interpretar el proceso DEL ARNÉS (reconciliación spine⟷I-77 CERRADA por interop 2026-07-07, HS-12): spine + `spine.categorias` (enum FIJO I-77 RN-28: propuesto·en-progreso·completado·descartado·pausado; terminalidad derivada ∈{completado,descartado}) + gates/dueños DERIVADOS de contratos por caja (`box.contract.schema.json`: estado «de → a» · gate{tipo} · ruta[]; dueño de estado = caja cuya transición LLEGA · con dueño-caja = de rol, resto = operador). El arnés NO shipea descriptor I-77 aparte (I-77 materializado = export/proyección). Visor Proceso. Nota: semilla `marketplace-arneses` 0.1.0 sin `categorias` → rehidratar del dogfood actualizado o esperar publish fase-5 | nueva | spec clon §3 + §7 · spec registry-arneses changelog v1.1 · decisión vigente I-77 | propuesta |

## 🟢 Después — profundizar la sesión + diferenciadores

| ID | Ítem | Incrementa | Fuente (detalle) | Estado |
|---|---|---|---|---|
| PB-09 | Vista Sesión completa: terminal embebido, composer (7 acciones, cola de mensajes, redimensionable) | CAP-02/03 | spec clon §10.1 | propuesta |
| PB-10 | Tab Cambios: quick-look + revisar side-by-side, chips por rol, commit por lote + adjuntar conversación | nueva | spec clon §10.2 | propuesta |
| PB-11 | Diff conversacional (comentar inline, el rol resuelve) — diferenciador genuino, ningún referente lo tiene | PB-10 | spec clon §5.2 + §10.2.4 | propuesta |
| PB-12 | Checkpoints/revert por turno (git ref privado) | CAP-02 | spec clon §5.3 | propuesta |
| PB-13 | Setup script por proyecto (install/migraciones/.env al crear sesión) | PB-05 | spec clon §5.4 | propuesta |
| PB-14 | PR flow integrado con CI (crear PR, checks, comments) — encaja con gates del proceso | PB-08 | spec clon §5.5 | propuesta |
| PB-15 | Pantalla Detalle de Capability (status, dueño, parent, historias vinculadas) | PB-07 | spec clon §7 | propuesta |
| PB-16 | Biblioteca de prompts (tabs proyecto/global, exportable) | nueva | spec clon §10.1.5 | propuesta |
| PB-17 | DevOps: comandos persistentes + túneles de vista previa + SSH | nueva | spec clon §4 (journey DevOps) + §10.1 | propuesta |

## 🧱 Deuda técnica declarada del producto

| ID | Ítem | Incrementa | Fuente (detalle) | Estado |
|---|---|---|---|---|
| PB-18 | SSE replay por `Last-Event-ID` | CAP-04 | `arch/INDEX.md` §stack (TBD declarado) | propuesta |
| PB-19 | Automatizar checks de boundaries `proposed` (conductor-encapsula-stream-json; go-arch-lint o test import-graph) | — (arch) | `arch/INDEX.md` tabla boundaries | propuesta |
| PB-20 | Segundo adaptador de agente — prueba real de intercambiabilidad de `AgentPort` | CAP-02 | LEDGER DH-13 §Siguiente (timing sin decidir) | propuesta |

## ⚪ TBD / bloqueado por F0

| ID | Ítem | Incrementa | Fuente (detalle) | Estado |
|---|---|---|---|---|
| PB-21 | Organización multi-usuario vía GitHub (qué sincroniza el repo) — terreno propio, ningún referente lo tiene. Incluye permisos por puesto sobre el registry (reframe F0) | nueva | VISION §TBD · spec clon §6 (hallazgo central) · NORTE-FIRMADO §Fuera | **TBD firmado F0 2026-07-07** — dueño Chris; decidir antes de PB-22/PB-23. No bloquea el port |
| PB-22 | Vista CTO multi-proyecto (journey CTO firmado como hueco declarado en F0) | PB-21 | spec clon §4 CTO + §0.5 · NORTE-FIRMADO §Journeys | bloqueada: PB-21 |
| PB-23 | Roles/accesos de primera clase (auth) + permisos por puesto («lo que mi jefe me habilita») | PB-21 | VISION §TBD · NORTE-FIRMADO §Modelo Rol/Registry | bloqueada: PB-21 |
| PB-24 | Instalable v1: **ventana nativa — fork Tauri-2-shell-tonto (patrón arnesia HS-04: daemon Go = sidecar, Rust mínimo, know-how casa ya pagado: mitigaciones WebKitGTK/Mint) vs Wails (un solo toolchain Go, menos maduro — watchlist de arnesia)** + cross-compile + firma + updater. ⚠ Antes de vender: Consumer Terms de Anthropic por escrito | nueva | NORTE-FIRMADO §Fases · VISION §Arquitectura ⚠ · harness-studio HS-04 + research fase3 (2026-07-05) | bloqueada: rebanadas F2+ |

## Trabajo de PROYECTO relacionado (no es backlog de producto — solo punteros)

- ~~Rellenar los 4 slots `__FILL_ME__` de `project.config.yaml`~~ — HECHO 2026-07-07 (DH-14):
  los 4 rellenos con firma (registry en domain_modules · roster ajustado al repo · 10 estados ·
  wip_caps advisory).
- Al cerrar la ÉPICA (no F0 — corregido en DH-14): promover lo permanente a `specs/` (carpeta
  aún no existe) y borrar la carpeta temporal, como manda el patrón. El norte firmado vive en
  la carpeta mientras la épica corra.
