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
| PB-04 | **Shell PRENTER** (movida de 🟡, 2026-07-07: «cerremos el shell primero»): rail Repositorios→Workspaces (workspace=sesión 1:1 — el rail F1 muere) + studio-nav (Studio/Backlog/Producto/Roles-pronto/Config) + overlays + **panel Cambios git REAL** (status/diff/log/commit por pathspec; sin push/pull — boundary nuevo) + tokens PRENTER + Storybook (absorbe PB-03). Alcance real/stub congelado en la spec | CAP-01/03/06 + nueva (navegación) | **`specs/shell/SPEC.md` (v1, por ratificar)** · mockup (SSoT forma) · spec clon §0.5/§10 · DS PRENTER (Claude Design `a98c2e0d`) | **en-curso · DH-15** |
| PB-02 | **Workspace aislado por sesión** (git worktree + branch) **+ validación de rutas protegidas** ($HOME, ~/.ssh…) en la misma entrega — forma firmada en F0. Cierra la brecha `sesion-aislada-por-cwd` convirtiéndola en feature (patrón Conductor). **Orden 2026-07-07: entra DESPUÉS del shell (PB-04), ya en la aplicación** | CAP-02 | `arch/boundaries/sesion-aislada-por-cwd.md` · spec clon §5.1 · NORTE-FIRMADO §Dudas (DH-14) · spec shell §9 | propuesta |
| PB-03 | **Sistema de diseño propio: Storybook + atomic design** | CAP-06 | premisa operador 2026-07-07 | **fusionada en PB-04** (2026-07-07): Storybook nace CON el shell — RN-9 de `specs/shell/SPEC.md`; tokens = PRENTER |

## 🟡 Próximo — shell + port por rebanadas (orden de herencia §3 de la spec clon)

| ID | Ítem | Incrementa | Fuente (detalle) | Estado |
|---|---|---|---|---|
| PB-05 | **Proyecto/Repositorio** como contenedor (rebanada 1 del port) — **parcialmente adelantada por PB-04** (registro de repos por ruta local + persistencia); esta PB queda para lo que falte (metadatos, ajustes por proyecto) | nueva | spec clon §3 (herencia `registryProject`) · spec shell §5 | propuesta |
| PB-25 | **Registry propio de arneses** — formato + repositorio de plugins + conexión de la app (decisión F0: registry PROPIO, no el sistema de plugins nativo de CC; el payload se materializa en el proyecto como artefactos que el `claude` spawneado carga — DH-10 intacto) | nueva | NORTE-FIRMADO §Modelo Rol/Registry (DH-14) | propuesta |
| PB-06 | **Rol** — vista Roles conectada al registry propio: navegar arneses, instalar por proyecto (rol = arnés instalado; cero roles locales, sin catálogo en la app). Redefinida en F0 — depende de PB-25 | nueva | NORTE-FIRMADO §Modelo Rol/Registry · spec clon §3 (nota ★) | propuesta |
| PB-07 | **Historia + Capability en la app** — kanban 10 estados como overlay, tarjeta muestra su capability, regla dura sesión↔historia (picker de `+ Nueva sesión`) | nueva | spec clon §0.5 + §3 (herencia `Story`/`Capability`) | propuesta |
| PB-08 | **Proceso as code** — interpretar descriptor I-77 (estados/gates/dueños) + visor Proceso | nueva | spec clon §3 + §7 · decisión vigente I-77 | propuesta |

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
| PB-24 | Instalable v1: cross-compile Win/Linux/mac + firma + updater. ⚠ Antes de vender: verificar Consumer Terms de Anthropic por escrito | nueva | NORTE-FIRMADO §Fases · VISION §Arquitectura ⚠ | bloqueada: rebanadas F2+ |

## Trabajo de PROYECTO relacionado (no es backlog de producto — solo punteros)

- ~~Rellenar los 4 slots `__FILL_ME__` de `project.config.yaml`~~ — HECHO 2026-07-07 (DH-14):
  los 4 rellenos con firma (registry en domain_modules · roster ajustado al repo · 10 estados ·
  wip_caps advisory).
- Al cerrar la ÉPICA (no F0 — corregido en DH-14): promover lo permanente a `specs/` (carpeta
  aún no existe) y borrar la carpeta temporal, como manda el patrón. El norte firmado vive en
  la carpeta mientras la épica corra.
