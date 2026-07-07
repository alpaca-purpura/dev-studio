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
| PB-01 | **F0: firmar el norte de la épica** — journeys por rol, principios, modelo de organización, orden de herencia, dudas abiertas §9. Gate de definición: deuda reconocida en DH-13 (F1 entregó código sin norte firmado) | — (gate) | `epicas/experiencia-orquestada/NORTE-BORRADOR.md` + `CLON-AGENTSROOM-SPEC.md` §8-9 | propuesta |
| PB-02 | **Workspace aislado por sesión** (git worktree + branch) — cierra la brecha de seguridad `sesion-aislada-por-cwd` convirtiéndola en feature (patrón Conductor) | CAP-02 | `arch/boundaries/sesion-aislada-por-cwd.md` · spec clon §5.1 · LEDGER DH-13 §Siguiente (candidata DH-14) | propuesta |
| PB-03 | **Sistema de diseño propio: Storybook + atomic design** — hoy solo hay tokens copiados (`theme.css`); toda UI que viene debe componerse de átomos catalogados. Fundacional: precede a las rebanadas de UI | CAP-06 | premisa operador 2026-07-07 · `project.config.yaml § design_system_ref` | propuesta |

## 🟡 Próximo — shell + port por rebanadas (orden de herencia §3 de la spec clon)

| ID | Ítem | Incrementa | Fuente (detalle) | Estado |
|---|---|---|---|---|
| PB-04 | **Shell nuevo**: sidebar Repositorios→Workspaces + overlays Backlog/Configuración — pasar el mockup ratificado (§0.5) a la app real | nueva (navegación) | spec clon §0.5 + `mockup-clon-agentsroom.html` (SSoT del shell) | propuesta |
| PB-05 | **Proyecto/Repositorio** como contenedor (rebanada 1 del port) | nueva | spec clon §3 (herencia `registryProject`) | propuesta |
| PB-06 | **Rol** — plantillas (prompt + proveedor/modelo), roster corto curado | nueva | spec clon §3 + §6 (no catálogo de 235) | propuesta |
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
| PB-21 | Organización multi-usuario vía GitHub (qué sincroniza el repo) — terreno propio, ningún referente lo tiene | nueva | VISION §TBD · spec clon §6 (hallazgo central) | bloqueada: F0 decide |
| PB-22 | Vista CTO multi-proyecto (hueco reconocido tras eliminar la Torre) | PB-21 | spec clon §4 CTO + §0.5 | bloqueada: PB-21 |
| PB-23 | Roles/accesos de primera clase (auth) | PB-21 | VISION §TBD | bloqueada: PB-21 |
| PB-24 | Instalable v1: cross-compile Win/Linux/mac + firma + updater. ⚠ Antes de vender: verificar Consumer Terms de Anthropic por escrito | nueva | NORTE-BORRADOR §FN · VISION §Arquitectura ⚠ | bloqueada: rebanadas F2+ |

## Trabajo de PROYECTO relacionado (no es backlog de producto — solo punteros)

- Rellenar los 4 slots `__FILL_ME__` de `project.config.yaml` (`domain_modules`, `agent_roster`,
  `value_stream`, `wip_caps`) — propuesta concreta ya escrita en spec clon §3, falta ratificar.
- Al cerrar F0: promover lo permanente de la épica a `specs/` (carpeta aún no existe) y borrar
  la carpeta temporal, como manda el patrón.
