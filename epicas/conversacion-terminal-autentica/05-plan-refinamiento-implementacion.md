# 05 · Plan de refinamiento + implementación (accionable)

> Cómo arrancar una conversación nueva y ejecutar. Anclado en los docs 00–04 + el mock. Este plan
> propone; el refinamiento por el proceso del arnés (spec + DH-NN) decide el slicing final.

## Norte
La conversación de dev-studio se siente como un **terminal real** (estética making-of de Anthropic) con
**doble input** (chat ↔ terminal, mismo stdin) — colgada de Claude Code hoy, extensible a otros CLIs
mañana. Veredicto = **Híbrido 3a** (`02-claude-code-protocol.md` §B/§C).

## Decisiones ya tomadas por la investigación (no re-abrir sin motivo)
1. **Híbrido 3a**: un proceso `claude` stream-json → dos superficies. Sin PTY para Claude; xterm.js en
   modo controlado, deltas reales = typing auténtico. (PTY real = escalón #3/#4, diferido.)
2. **Teal único PRENTER**. Descartar el multi-accent per-cluster de AgentsRoom. Rol/status por forma+label.
3. **JetBrains Mono** en toda la superficie terminal (ya vendorizada = match Anthropic).
4. **Envelope normalizado alineado a ACP** para no inventar protocolo (para el eje multi-CLI).
5. **Adapters compilados + registry declarativo** near-term (no sidecars ACP todavía).

## Decisiones CERRADAS en F0 (firma Chris 2026-07-09 · DH-19 · ver `NORTE-FIRMADO.md`)
- **A. Probe primero: SÍ** ✅ — R0 es un spike descartable de 1 archivo que valida el envelope real de
  `can_use_tool` contra el `claude` instalado (**v2.1.205 pineada**) antes de tocar UI. Precede a R2.
- **B. Alcance: MULTI-PROVEEDOR DESDE F1** ✅ **(Chris divergió del recomendado «Claude primero»).**
  Consecuencia en el slicing (ver §Reordenamiento por B): el `Envelope` nace ACP-aligned en R1 (no
  refactor en R4); `Capabilities()` + `ProviderSessionID` (hueco #4) se implementan temprano; el segundo
  adapter (Amp) entra en el alcance FIRME (R5 deja de ser «fase 2 aparte») como prueba de PB-20.
- **C. Terminal: FAUX / modo controlado 3a** ✅ — sin PTY para Claude; PTY real = escalón futuro.
- **D. Slicing LEDGER: DH por rebanada entregada (R1–R5)** ✅ — cada rebanada su DH + CAP en el commit
  de entrega. La firma de F0 = DH-19 aparte.
- **E. PB-24 (shell nativo): NO BLOQUEA** ✅ — 3a corre en el HTTP+webview actual; PB-24 se decide aparte.

## Reordenamiento por la decisión B (multi-proveedor desde F1)
El slicing R0–R5 de abajo se escribió asumiendo «Claude primero». Con B firmado, dos ajustes de orden
(el CONTENIDO de cada rebanada no cambia, solo CUÁNDO se hace el núcleo provider-agnóstico):
1. **R1 emite `Envelope` normalizado (ACP-aligned) desde el vamos** — `translate()` NO produce
   `AgentEvent` Claude-shaped para migrarlo luego; nace en la forma final. El puente `dockFrame`→`onDock`
   transporta el Envelope, no el shape viejo.
2. **`Capabilities()` + `ProviderSessionID` (R4/hueco #4) se adelantan al núcleo de F1** — la UI gatea
   panes por `Capabilities` desde R1/R2. R4 queda como el CIERRE/pulido de la normalización, no su
   nacimiento. El segundo adapter (R5 · Amp) sube a alcance firme de la épica.

## Rebanadas propuestas (cada una entrega valor sola · TDD, RN del proyecto)
> **Piso HARD**: verificación REAL en vivo (ejercer la acción real + leer logs + confirmar efecto), no
> "HTTP 200". Bug fix = regression test RED primero. (test-design-doctrine.md, tdd-mandatory.md.)

**R0 · Probe del protocolo de control** (spike, 1 archivo, descartable)
- Spawn `claude --input-format stream-json --output-format stream-json --verbose --permission-prompt-tool
  stdio`, mandar un turno que dispare un tool gated (ej. Bash fuera de allowlist), loguear stdin/stdout
  crudo. Documentar el envelope REAL (`control_request`/`control_response` shape + `request_id`) + la
  versión de `claude`. Salida: un `.md` con el shape verificado. **Precede a R2.**

**R1 · Tool-cards** (hueco 1 · `conductor.go`)
- `translate()`: parsear bloques `tool_use` (en `assistant`) + `tool_result` (en `user`) → nuevos
  `AgentEventKind` (tool.call / tool.result). Propagar por `dockFrame` → `onDock` → nuevas cards en
  `session-view.tsx`. Render estilo mock (tarjeta con `⎿ Read path · OK`, diff add/del).
- AC: un turno que usa Read+Edit muestra 2 tool-cards con input + output/diff reales.

**R2 · Permisos + modal de rama** (hueco 2 · depende de R0)
- Spawn con `--permission-prompt-tool stdio`. `conductor.go`: parsear `control_request:can_use_tool` →
  `EvPermissionReq`; implementar `ReplyPermission(requestID, allow, [message])` que escribe
  `control_response` en stdin. Nuevo endpoint `POST /api/sessions/{id}/permission` + reducer. UI: modal de
  rama estilo terminal (`❯ 1. Permitir · 2. Siempre · 3. Rechazar`) como el mock.
- Fallback si el canal falla: `--permission-mode acceptEdits` + allow-rules (degrada UX, no bloquea).
- AC: un tool gated pausa el turno, muestra el modal, la respuesta del usuario permite/rechaza en vivo
  (leer logs: el proceso continúa/aborta según la elección).

**R3 · Terminal xterm.js modo controlado + doble input** (hueco 3 + segundo input path)
- FE: agregar `@xterm/xterm` + addons `fit`/`webgl`/`web-links`. Pane terminal que `term.write()` los
  deltas reales (typing auténtico) + líneas de prompt faux. `term.onData` → en Enter, empaquetar línea
  como turno → **mismo pipeline que el composer** (POST /turn). Tabs Chat|Terminal (como el mock).
- BE: exponer `RawTerminal()` en el port (aunque para 3a el terminal se alimenta de deltas, no de PTY —
  RawTerminal queda para el escalón PTY futuro / otros proveedores).
- Command palette sembrado de `system/init` (`slash_commands`/`tools`).
- AC: escribir en el terminal Y en el chat produce turnos al MISMO proceso; el stream se ve como typing
  en el terminal. Window-chrome macOS + JetBrains Mono + teal.

**R4 · Capabilities + normalización** (hueco 4 · habilita multi-CLI)
- Renombrar `ClaudeSessionID`→`ProviderSessionID`; agregar `Capabilities()` al `AgentPort`; refactor
  `AgentEvent` → `Envelope` alineado-ACP (`02`/`03`). El adapter Claude declara sus capacidades. UI lee
  capabilities para decidir panes (gating graceful).
- AC: la UI renderiza idéntica leyendo `Capabilities` de un manifest; ningún branching per-proveedor
  hardcodeado.

**R5+ · Segundo adapter (multi-CLI)** — fase 2, decisión aparte
- Empezar por **Amp** (`StreamJSONAdapter`, near drop-in Claude-compatible) → prueba real de
  intercambiabilidad (PB-20). Después Familia B (Codex app-server / Gemini ACP) y C (OpenCode Go SDK).
  Aider = `RawPTYAdapter` terminal-only. MCD universal (`03` §7) = PTY passthrough + liveness busy/wait/idle.

## Archivos a tocar (referencia rápida — de `04-mapa-codigo-actual.md`)
- BE: `internal/adapters/agent/claudecode/conductor.go` (`translate`, `Send`, nuevo `ReplyPermission`,
  `RawTerminal`), `internal/ports/agent.go` (Envelope/Capabilities/AgentSession), `internal/usecase/
  session_service.go` (`consume`/`publish` nuevos kinds), `internal/adapters/transport/http/router.go` +
  `sessions.go` (endpoint permission), `internal/adapters/transport/sse/broker.go`.
- FE: `web/package.json` (+xterm), `web/src/widgets/session-view/ui/session-view.tsx` (tool-cards, modal,
  terminal pane, tabs), `web/src/shared/store/sessions-store.ts` (`onDock` nuevos kinds, replyPermission),
  `web/src/shared/api/{sse.ts,client.ts,types.ts}`.
- Boundary: respetar `conductor-no-parsea-jsonl` (cada adapter dueño de su protocolo; usecase solo ve
  Envelope). Ver PB-19 (automatizar el check).

## Riesgos / caveats a no olvidar
- Protocolo de control **version-dependiente** → R0 probe obligatorio antes de R2. Pinear versión de `claude`.
- El resume puede acuñar session_id nuevo (#12235) → leer id del evento, no asumir.
- Objetos JSON de stdout pueden abarcar varios read chunks → bufferear/splitear en `\n` (ya lo hace el
  pump con `ReadBytes`, pero cuidar al agregar frames grandes).
- SSE broker dropea subs lentos + sin replay (PB-18) → cuidar con streams largos de tool output.
- Windows: si se va a PTY real (escalón futuro) → ConPTY.

## Entregables de este research (todos en esta carpeta)
- `00-visual-reference.md` — tokens + inventario de interacción + screenshots (anthropic-*.png, agentsroom-*.png).
- `01-sintesis.md` — síntesis ejecutiva.
- `02-claude-code-protocol.md` — protocolo stream-json + terminal-en-webview (completo).
- `03-multi-cli-interop.md` — normalización multi-CLI + contrato Go + ACP (completo).
- `04-mapa-codigo-actual.md` — file:line del código actual + GAP summary.
- `05-plan-refinamiento-implementacion.md` — este plan.
- `mock-conversacion.html` — mock visual (Artifact publicado; ver README para URL).
- `*.png` — capturas de referencia (Anthropic making-of, AgentsRoom, render del mock).
