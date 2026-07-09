# Épica «Conversación terminal-auténtica» — NORTE FIRMADO (F0 cerrada 2026-07-09 · DH-19)

> Firmado por Chris en la sesión F0 (2026-07-09), fork por fork vía AskUserQuestion — ficha
> `DH-19` en [`LEDGER.md`](../../LEDGER.md). Este archivo ES el borrador renombrado
> (`NORTE-BORRADOR.md` → `NORTE-FIRMADO.md`). Carpeta temporal por diseño (patrón Torre de
> Control): al cierre de la ÉPICA (no de F0) se borra; lo permanente se promueve a `specs/`.
> Este norte gobierna F1+ hasta entonces.
>
> **SSoT de la experiencia** mientras la épica corra = esta carpeta + `mock-conversacion.html` +
> `06-especificacion-mockup.md`. Todo cambio de forma se escribe en el mockup Y el `06` a la vez
> (patrón `spec-mapa-funcional.md`). Si un cambio no queda escrito en ambos, no pasó.
>
> **Relación con «Experiencia Orquestada»:** esta épica **aterriza sobre el shell** que aquella
> entregó (rail Repositorios→Workspaces, `specs/shell/SPEC.md`). No lo reemplaza — profundiza la
> **superficie de sesión/conversación** (el hueco PB-09) y el **driver multi-proveedor**.

## Cruda del operador (2026-07-09, verbatim)

"quiero saber cómo se hizo esta web https://www.anthropic.com/features/making-of-claude-code y cómo
podríamos tomar el cómo se ve para la conversación del dev studio, considerando que nos colgaremos de
anthropic, quiero que se sienta lo más real, tal como https://agentsroom.dev/es#demo … en el caso de
agentsroom se muy detallista ya que se puede conversar por el chat de conversación o escribiendo directo
en su «terminal». Quiero dar ese efecto."

"Elevar a épica para atenderla de una vez."

## El norte (qué construimos)

Que la **conversación de dev-studio se sienta como un terminal REAL** —estética *making-of de Claude
Code* de Anthropic— con **doble input**: el usuario chatea por el composer O escribe directo en el
terminal, y **ambos escriben el MISMO stdin del proceso `claude`**. AgentsRoom *finge* ese doble input;
DevStudio lo hace **literal**, porque el driver ya spawnea `claude` stream-json de verdad (BYO
licencia, cero API de Anthropic). El efecto no se simula: se cablea sobre el proceso real.

**Veredicto técnico (de la investigación, ver `01`–`03`): Híbrido 3a** — un proceso stream-json
alimenta DOS superficies desde un solo stream: (a) cards estructuradas (chat, roster, tool-cards, modal
de permisos) y (b) un terminal `xterm.js` en modo controlado (deltas reales = typing auténtico, sin
PTY). Extensible a otros CLIs (Amp, Codex, Gemini, OpenCode, Aider) tras un adapter normalizado.

## Principios / invariantes (RATIFICADOS en F0)

1. **Terminal literal, no fingido** — el terminal se alimenta del stream real; las teclas producen
   turnos reales al mismo stdin. Nada de typewriter fake en vivo.
2. **Doble input = un solo stdin** — composer y terminal convergen en el mismo pipeline (`ccSession.Send`).
3. **Un proceso, dos superficies** — no dos procesos, no PTY para Claude (escalón #3/#4 diferido).
4. **Teal único PRENTER + JetBrains Mono** — se descarta el multi-accent per-cluster de AgentsRoom
   (viola RN de acento único); rol/status por forma+label. Amber = ÚNICA excepción y es semántica
   (atención/«te espera»), no acento de marca. JetBrains Mono ya vendorizada = match Anthropic.
5. **Degradación graceful por `Capabilities`** — la UI decide panes leyendo capacidades del proveedor;
   cero branching per-proveedor hardcodeado. MCD universal = PTY passthrough + liveness busy/wait/idle.
6. **Envelope normalizado alineado a ACP** — no inventar protocolo propio.
7. **Cero API / BYO licencia** — se mantiene el driver CLI-nativo (DH-10 intacto).
8. **Boundary `conductor-no-parsea-jsonl`** — cada adapter dueño de su protocolo; el usecase solo ve `Envelope`.

## Decisiones FIRMADAS en F0 (los forks A–E · vía AskUserQuestion 2026-07-09)

- **A · Probe primero (R0): SÍ** ✅ (recomendado). El protocolo de control `can_use_tool` está
  sub-documentado + version-dependiente (bugs #34046, #12235). Un spike de 1 archivo (descartable)
  valida el envelope REAL contra el `claude` instalado **antes** de construir el modal de rama (R2).
  **Versión pineada: `claude` 2.1.205** (verificada en la máquina de Chris, 2026-07-09).

- **B · Alcance de la primera entrega: MULTI-PROVEEDOR DESDE F1** ✅ **(Chris divergió del recomendado
  «Claude primero»).** Consecuencia firmada — se reordena el slicing:
  - El **`Envelope` nace ACP-aligned desde R1** (no es un refactor tardío en R4). `translate()` emite
    el envelope normalizado desde el primer hueco, no `AgentEvent` Claude-shaped que después se migra.
  - **`Capabilities()` + `ProviderSessionID` (hueco #4) se implementan temprano**, junto a R1/R2 — no
    al final. La UI gatea panes por `Capabilities` desde el arranque.
  - El **segundo adapter (Amp, near drop-in) entra en el alcance FIRME de la épica** (R5 deja de ser
    «fase 2, decisión aparte»): es la prueba de intercambiabilidad de `AgentPort` (absorbe/valida PB-20).
  - Riesgo aceptado: más superficie de golpe y valor visible de la conversación algo más tardío, a
    cambio de no pagar el refactor del Envelope dos veces y validar el seam provider-agnóstico ya.

- **C · Terminal auténtico: FAUX / modo controlado 3a** ✅ (recomendado). xterm.js alimentado de deltas
  reales (typing auténtico) + cards + aprobaciones estructuradas; sin PTY para Claude. PTY real (#3/#4)
  = escalón futuro, solo si aparece necesidad de ANSI/TUI literal (decisión de producto aparte).

- **D · Slicing en el LEDGER: DH POR REBANADA ENTREGADA (R1–R5)** ✅ (recomendado). Cada rebanada
  entrega valor sola → su propia DH + CAP en el commit de entrega. LEDGER granular, calza con la
  disciplina CAP-por-commit. **Esta firma de F0 = DH-19** (aparte de las DH de entrega).

- **E · Relación con PB-24 (shell nativo Tauri/Wails): NO BLOQUEA** ✅ (default recomendado, sin
  objeción de Chris). El adapter es Go pase lo que pase; 3a corre en el HTTP+webview actual. PB-24 se
  decide aparte (sigue bloqueada para rebanadas F2+ de su propio alcance, no de esta épica).

## Fases (FIRMADAS — reordenadas por la decisión B · mapea a las rebanadas R0–R5 de `05-plan…`)

| Fase | Contenido | Rebanadas |
|---|---|---|
| **F0** | ✅ CERRADA — norte firmado + decisiones A–E. Ficha **DH-19** | — |
| **F1** | Probe del protocolo de control + tool-cards + **Envelope ACP-aligned + `Capabilities()` desde el inicio** (por B) | R0, R1, (núcleo R4) |
| **F2** | Permisos (modal de rama) + terminal xterm.js + doble input | R2, R3 |
| **F3** | Cierre de Capabilities + normalización del Envelope (lo que reste de ACP-alignment) | R4 |
| **F4** | Segundo adapter (Amp near-drop-in → prueba de intercambiabilidad, PB-20) + multi-CLI | R5 |

> **Nota de reordenamiento (B):** el trabajo de `Capabilities()`/`Envelope` de R4 se adelanta al núcleo
> de F1 para que Claude nazca ya como «un proveedor más» y el segundo adapter no obligue a un refactor.
> R4 queda como el cierre/pulido de la normalización, no como su nacimiento. Detalle en `05-plan…`.

Cada rebanada entrega valor sola (detalle + AC + archivos a tocar en `05-plan-refinamiento-implementacion.md`).

## Relación con el backlog / otras épicas

- **Absorbe PB-09** (terminal embebido). **Toca:** PB-24 (shell nativo · no bloquea, fork E), PB-25
  (registry ArnesIA = eje proveedor como adapter instalable), PB-20 (segundo adapter =
  intercambiabilidad, ahora en alcance firme por B). Fichada como **PB-28** (banda 🔴 Ahora, firma
  Chris 2026-07-09).
- **Hereda** el shell de «Experiencia Orquestada» (`specs/shell/SPEC.md`).

## Al cierre de la épica

Promover lo permanente a `specs/` (probable: `specs/conversacion-sesion/SPEC.md` +
`specs/driver-multiproveedor/SPEC.md`) y **borrar esta carpeta temporal**. El mockup se promueve como
SSoT-de-forma de la spec correspondiente (patrón `specs/shell` ← mockup de la épica).

## Insumos de esta carpeta (fuentes)

`README.md` (índice) · `00-visual-reference.md` (tokens + screenshots) · `01-sintesis.md` ·
`02-claude-code-protocol.md` · `03-multi-cli-interop.md` · `04-mapa-codigo-actual.md` ·
`05-plan-refinamiento-implementacion.md` · **★ `06-especificacion-mockup.md`** (cada propuesta del
mockup) · **★ `mock-conversacion.html`** (mockup visual · Artifact publicado) · `*.png`.
