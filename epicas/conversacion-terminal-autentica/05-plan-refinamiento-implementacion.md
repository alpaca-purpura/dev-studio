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

**R0 · Probe del protocolo de control** (spike, 1 archivo, descartable) — ✅ **EJECUTADA 2026-07-09 → [`07-R0-probe-findings.md`](./07-R0-probe-findings.md)** (`claude` v2.1.205 pineada)
- Hecho: se spawneó el `claude` real en stream-json, se disparó un tool gated (Bash) y se logueó
  stdin/stdout crudo + se desensambló el binario. **Hallazgos clave:** (1) `--permission-prompt-tool
  stdio` está MAL (el flag exige un MCP tool, no `"stdio"`); (2) el modo `default` headless AUTO-PERMITE
  los tools (no hay «ask» por defecto); (3) el handshake `initialize` funciona pero el `initialize`
  mínimo NO enruta el permiso; (4) shapes de `init`/`tool_use`/`tool_result`/`text_delta`/`result`
  capturados EN VIVO (R1 des-riesgado); (5) el «ask» es `control_request{subtype:can_use_tool}` ⇄
  `control_response{behavior:allow|deny}` — falta pinear el disparador (primera tarea de R2).

**R1 · Tool-cards** (hueco 1 · `conductor.go`) — 🔨 **IMPLEMENTADA · pendiente Chris-verify in-app (G)**
- Hecho (TDD RED→GREEN): `translate()` → `[]AgentEvent`, parsea `tool_use` (assistant) + `tool_result`
  (user) → `EventToolCall`/`EventToolResult` (nombres ACP-flavored, decisión B). `dockFrame` gana campos
  tool; `consume()` publica `tool.call`/`tool.result` (inline, sin tocar estado del turno). FE: `ToolCall`
  type + `toolCalls` en el store (parea por tool_id, reset por turno) + átomo `ToolCard` (`⎿ Name arg ·
  estado`, diff +/− con `DiffStat`, output colapsable, teal/mono) + story (RN-9) + wiring en `session-view`.
- Tests: `conductor_test.go` (TestTranslateToolUse/ToolResult con frames REALES del probe R0) +
  `session_toolcards_test.go` (consume→frames pareados). `go test ./...` verde + `go vet` + `tsc+vite` build.
- **Live-verify (claude REAL, pipeline aislado):** un turno Read+Edit produjo las 2 tool-cards con datos
  reales — Read (verde, output `1\tline one\n…`), Edit (diff real `old_string`→`new_string`); el
  `tool.result` del Edit vino `is_error:true` («requested permissions to write») = write auto-denegado
  headless (card roja honesta; el write-success es de **R2**, tal como R0 predijo). El chequeo visual en la
  app INSTALADA queda para tu G (dogfooding «Actualizar» → mandar un turno con tools → ver las cards).
- AC: un turno que usa Read+Edit muestra 2 tool-cards con input + output/diff reales. **Datos ✅ probados
  en vivo; render in-app = Chris-verify.**

**R1.5 · Transcript persistente + reconstrucción en re-entry** (ruta A · firmado Chris 2026-07-09) — ✅ **IMPLEMENTADA + verificada en vivo**
- **Problema:** las tool-cards de R1 son **live-only** (mueren al salir/recargar). Lo persistido (`domain.Session.
  Conv []Turn{role,text}`) es **solo texto** → al reentrar, las cards no vuelven a su lugar en el transcript.
- **Ruta A (firmada):** al reentrar/abrir una sesión con `ClaudeSessionID`, el **adapter reconstruye el
  transcript ordenado** (texto + tool.call/tool.result **en su lugar**) desde **su propia fuente** — para
  Claude = su JSONL `~/.claude/projects/<cwd-encoded>/<session_id>.jsonl` (transcript completo que claude
  ya escribe para `--resume`). Simétrico con la rehidratación del resume (`02` §A.5).
- **Diseño (provider-agnóstico, decisión B):**
  - Puerto: `History(ctx) ([]Envelope, error)` en `AgentSession` + flag `Capabilities.Replay` — cada adapter
    reconstruye del suyo (Claude JSONL · OpenCode API history · Aider degrada a texto vía Capabilities).
  - Adapter `claudecode`: **único parser del JSONL** (boundary `conductor-no-parsea-jsonl` INTACTO) — localiza
    el archivo por `~/.claude/projects/*/{ClaudeSessionID}.jsonl` (el session_id es único), parsea líneas
    `user`/`assistant` (content[] → text · tool_use · tool_result) → `[]Envelope` ordenado. Archivo nuevo
    (ej. `transcript.go`), NO toca el pump del stream vivo.
  - Usecase + HTTP: `GET /api/sessions/{id}/transcript` → ítems ordenados `{kind, role, text, tool_*}`.
  - FE: al activar una sesión con `claude_session_id` → fetch transcript → render rico (texto + tool-cards
    en orden) reemplazando el `conv` solo-texto; re-fetch al cerrar un turno (`result`) para plegar las
    cards vivas al historial persistido. `session.transcript` en el store.
- **Verificación REAL (piso HARD):** sesión con turnos que usaron tools → **salir de la sesión / recargar la
  app / reabrir** → las tool-cards reaparecen **en su lugar** en el transcript (reconstruidas del JSONL),
  no solo el texto. Leer el JSONL real + confirmar el orden.
- **Liga con:** R4 (la capability `Replay` vive con `Capabilities`/`ProviderSessionID`/resume — decisión B ya
  los adelantó) · `02` §A.5 (rehidratar tras restart) · `06` P4 (el transcript del mock ya muestra cards inline).
- **Verificado en vivo (2026-07-09, CDP sobre la app real):** turno Read+Edit → **recarga dura (= salir/
  reentrar)** → reactivar la sesión → el transcript se **reconstruyó del JSONL** con las 2 cards EN SU LUGAR
  (Read ✓ + output `1 alpha…`, Edit diff `+// reentry OK` ✕), no solo texto. Screenshot `r15-reentry.png`.
  Entregado: puerto `TranscriptItem`/`History` · adapter `transcript.go` (parser JSONL, único parser) ·
  usecase `Transcript` (History con fallback a conv) · `GET /api/sessions/{id}/transcript` · FE store
  `transcript` + `loadTranscript` (switchTo/init/post-result) + render `foldTranscript` (pliega call+result
  en una card). TDD: `transcript_test.go` (fixture shape real) + `session_toolcards_test.go` (History + fallback).

**R2 · Permisos + modal de rama** (hueco 2 · depende de R0) — ⚠️ **RE-SCOPEADA por R0** (ver `07` §2/§4)
- **PRIMERO (spike A→C, continuación de R0):** pinear cómo el CLI enruta `can_use_tool` al cliente sobre
  stdio + capturar UNA transacción viva (allow Y deny) ANTES de la UI. `--permission-prompt-tool stdio`
  del plan viejo NO sirve (§0 de `07`). Tres caminos (`07` §4): **A** control-protocol `canUseTool`
  (recomendado — todo en el stream, `updatedInput`/`updatedPermissions` = opciones 1/2 del modal);
  **B** `--permission-prompt-tool mcp__…` (MCP tool in-process); **C** PreToolUse hook (lo que el binario
  bendice como fallback). Piso HARD: no construir el modal sin captura viva de `can_use_tool`.
- Después: `conductor.go` parsea `control_request{subtype:can_use_tool}` → `EvPermissionReq`;
  `ReplyPermission(requestID, allow, [message])` escribe `control_response{response:{behavior}}` en stdin.
  Endpoint `POST /api/sessions/{id}/permission` + reducer. UI: modal de rama (`❯ 1. Permitir · 2. Siempre
  · 3. Rechazar`) — «2. Siempre» = `updatedPermissions`; «3. Rechazar» = `behavior:deny` + `message`.
- Fallback de degradación: `--permission-mode acceptEdits` + allow-rules (no bloquea; se pierde el modal).
- AC: un tool gated pausa el turno, muestra el modal, la respuesta permite/rechaza en vivo (leer logs).

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
- `06-especificacion-mockup.md` — cada propuesta del mockup (P1–P6) especificada.
- `07-R0-probe-findings.md` — hallazgos del probe R0 en vivo (shapes verificados + mecanismo de permisos).
- `mock-conversacion.html` — mock visual (Artifact publicado; ver README para URL).
- `*.png` — capturas de referencia (Anthropic making-of, AgentsRoom, render del mock).
