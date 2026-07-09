# 02 · Claude Code stream-json + terminal-en-webview (brief técnico completo)

> Detalle crudo de la investigación (subagente 2). Fuente para refinar/implementar sin re-investigar.
> **Honestidad:** el schema de stdin y el envelope de control (A.3/A.4) NO están cubiertos por ninguna
> página oficial única — son consistentes across SDK source + docs de comunidad, pero **version-sensibles**.
> Verificar contra el `claude` instalado con un probe ANTES de construir UI.

## A) Protocolo stream-json (headless / bidireccional)

### A.1 — Invocación persistente bidireccional
La sesión persistente es real y es lo que usa el **Agent SDK** por debajo ("Streaming Input Mode… un
proceso long-lived que toma input del usuario, maneja interrupciones, expone permission requests y
maneja session management"). El SDK es un wrapper que spawnea esta forma del CLI — dev-studio puede
spawnear el mismo CLI directo, **sin SDK / sin API de Anthropic** (BYO-licencia intacto).

```bash
claude \
  --input-format  stream-json \
  --output-format stream-json \
  --verbose \
  --permission-prompt-tool stdio     # rutea aprobaciones de tool por el canal de control stdio
  # opcional: --model … --append-system-prompt … --plugin-dir … --resume <id> / --session-id <uuid>
  # opcional: --include-partial-messages  (deltas token-level para el typing en vivo)
```

Hechos de flags (con caveats honestos):
- `--output-format stream-json` → NDJSON, un evento por línea, emitido conforme ocurre. En modo `-p`
  print **requiere `--verbose`**. ([headless docs](https://code.claude.com/docs/en/headless))
- `--input-format stream-json` → "el único mecanismo CLI para comunicación **bidireccional**
  programática". Su schema de stdin es **oficialmente indocumentado** (issue abierto
  [#24594](https://github.com/anthropics/claude-code/issues/24594)). Lo de abajo está corroborado
  across SDK source + reverse-engineering de comunidad, no una página oficial.
- `--include-partial-messages` → agrega deltas `stream_event` token-level (A.2). Opcional; solo para
  el typing en vivo con deltas reales.
- `-p / --print`: diseñado para **one-shot**. Para sesión persistente, **omitir el prompt posicional**
  y apoyarse en `--input-format stream-json` para mantener el proceso interactivo-pero-headless. Las
  fuentes difieren en si `-p` debe estar; el Go SDK protocol doc + el SKILL de comunidad spawnean
  **sin** `-p`. Tratar `-p` acá como "no pasar prompt posicional".
- `--bare` (nuevo): salta auto-discovery de hooks/skills/plugins/MCP/CLAUDE.md. **NO usar `--bare`** —
  queremos el CLAUDE.md del proyecto + el plugin del arnés cargados. El rol carga vía `--plugin-dir`
  + `--append-system-prompt`.

### A.2 — Tipos de mensaje de salida (stdout NDJSON)
Cada línea = `{"type": …, "session_id": …, "uuid": …, …}`. El envelope envuelve un **objeto Message
completo de la API Anthropic** bajo la key `message` para assistant/user. **Crítico para el render:
`tool_use` y `tool_result` son *bloques* de content dentro de mensajes `assistant`/`user` — no eventos
top-level.**

| `type` | Shape / campos clave | Notas |
|---|---|---|
| `system` / `subtype:"init"` | `session_id`, `model`, `cwd`, `tools[]`, `mcp_servers[{name,status}]`, `permissionMode`, `apiKeySource`, `plugins[{name,path}]`, `plugin_errors[]`, `slash_commands?` | **Primer** evento. Sembrar el roster de tools / header del agente + el command palette. |
| `assistant` | `{type:"assistant", session_id, parent_tool_use_id, message:{ id, role, model, content:[…bloques…], stop_reason, usage }}` | bloques: `{type:"text",text}`, `{type:"thinking",…}`, `{type:"tool_use",id:"toolu_…",name,input}`. |
| `user` | `{type:"user", session_id, message:{ role:"user", content:[…] }}` | Acá llegan los bloques **`{type:"tool_result",tool_use_id,content}`** (output del tool devuelto a Claude). Con `--replay-user-messages` también eco de tus turnos. |
| `stream_event` | `{type:"stream_event", uuid, session_id, parent_tool_use_id, event:{…SSE Anthropic crudo…}, ttft_ms?}` | Solo con `--include-partial-messages`. `event.type` ∈ `message_start`/`content_block_start`/`content_block_delta`/`content_block_stop`/`message_delta`/`message_stop`. Texto en `event.delta.type=="text_delta" → event.delta.text`; input de tool streaming en `input_json_delta → partial_json`. |
| `system` / `subtype:"api_retry"` | `attempt`, `max_retries`, `retry_delay_ms`, `error_status`, `error` | Toast de retry. |
| `system` / `subtype:"plugin_install"` | `status ∈ started/installed/failed/completed`, `name?`, `error?` | Solo con `CLAUDE_CODE_SYNC_PLUGIN_INSTALL`; precede a `init`. |
| `system` / `subtype:"compact_boundary"` | marca compactación de historia | Divider sutil. |
| `result` | `subtype ∈ success/error_max_turns/error_during_execution/error_max_budget_usd`; `result` (texto final), `is_error`, `duration_ms`, `duration_api_ms`, `num_turns`, `total_cost_usd`, `usage{input_tokens,output_tokens,cache_*}`, `permission_denials[]`, `structured_output?`, `session_id` | Emitido al **final de cada turno** (no del proceso — el proceso sigue vivo para el próximo turno de stdin). |

Orden con parciales: `stream_event(message_start) → content_block_start/delta*/stop → message_delta →
message_stop → assistant(completo) → …tool ejecuta… → user(tool_result) → …próximo turno… → result`.
([streaming-output docs](https://code.claude.com/docs/en/agent-sdk/streaming-output) · shapes exactos
en [takopi cheatsheet](https://takopi.dev/reference/runners/claude/stream-json-cheatsheet/) · Anthropic
reconoce que la lista de tipos está sub-documentada: [#24612](https://github.com/anthropics/claude-code/issues/24612))

### A.3 — Enviar turnos de usuario por stdin
NDJSON, un objeto por línea, terminado en `\n`, stdin abierto. Shape (= `SDKUserMessage` del SDK):
```json
{"type":"user","message":{"role":"user","content":"fix the failing test in auth_test.go"}}
```
`content` puede ser **string** *o* **array de bloques** (texto + imágenes):
```json
{"type":"user","message":{"role":"user","content":[
  {"type":"text","text":"Review this diagram"},
  {"type":"image","source":{"type":"base64","media_type":"image/png","data":"<b64>"}}
]},"parent_tool_use_id":null}
```
- Fin de turno = el newline; el próximo `result` cierra el turno. Encolar más turnos cuando sea (se
  procesan secuencialmente).
- **Slash commands van embebidos en `content`**: `"…then run /commit"`. `claude` expande `/skill` y
  comandos custom antes de correr. Los de diálogo interactivo (`/login`, UI de `/resume`) **no** están
  headless; `/config key=value` sí.
- **Gotcha de framing** (de una UI custom real, [bswen](https://docs.bswen.com/blog/2026-03-21-stream-json-custom-ui-claude-code/)):
  los objetos JSON de stdout **pueden abarcar varios read chunks** — bufferear y splitear en `\n`,
  retener el fragment incompleto del final, `JSON.parse` por línea, fallback a raw si falla parse.
  Cap scrollback (~5000 líneas). No spawnear un proceso por mensaje.

### A.4 — Permission prompts / aprobaciones de tool (protocolo de control)
La parte load-bearing para UX IDE-grade. Las aprobaciones van por un **canal de control bidireccional
JSON-RPC-ish sobre los mismos pipes stdio**, habilitado por `--permission-prompt-tool stdio`. Tanto
`control_request` como `control_response` fluyen en ambas direcciones, matcheados por `request_id`:

```json
// CLI → host  (permiso necesario)
{"type":"control_request","request_id":"req_abc",
 "request":{"subtype":"can_use_tool","tool_name":"Bash",
   "input":{"command":"git add -A"},"tool_use_id":"toolu_xyz",
   "decision_reason":"Command not in allowlist"}}

// host → CLI  ALLOW  (eco de request_id; updatedInput puede reescribir la llamada)
{"type":"control_response","response":{"subtype":"success","request_id":"req_abc",
   "response":{"behavior":"allow","updatedInput":{"command":"git add -A"}}}}

// host → CLI  DENY  (incluir message → se le pasa a Claude como razón)
{"type":"control_response","response":{"subtype":"success","request_id":"req_abc",
   "response":{"behavior":"deny","message":"User declined"}}}
```
Host→CLI también podés mandar control_requests: `interrupt` (cancela el turno — el interrupt del SDK
mapea acá), `set_permission_mode`, y un handshake `initialize`. El CLI **bloquea** en un `can_use_tool`
pendiente (comunidad reporta ~60s default) → el modal debe responder rápido o el turno se cuelga.

**Caveats honestos del protocolo de control:**
- Ninguna página oficial documenta el envelope on-the-wire; el equivalente oficial es el callback
  **`canUseTool`** del SDK ([permissions docs](https://code.claude.com/docs/en/agent-sdk/permissions)).
  Los nombres de campo de arriba están corroborados por el Go SDK protocol doc + spec de comunidad.
- Las fuentes difieren en detalles (`control_request` vs `sdk_control_request`; `subtype:"can_use_tool"`
  vs `"permission"`). Tratar como hipótesis de trabajo y **validar contra tu `claude` pineado** con el
  probe antes de construir UI.
- **Bug version-dependiente a testear:** [#34046](https://github.com/anthropics/claude-code/issues/34046)
  — algunas versiones **no** emitían el `can_use_tool` bajo `--permission-prompt-tool stdio`. Pinear/
  verificar versión.
- Alternativa si el canal de control falla: `--permission-mode acceptEdits` (auto-aprueba writes + fs
  cmds comunes; el resto gated por `--allowedTools`/`permissions.allow`) o `dontAsk`. Cambia la UX de
  aprobación-en-vivo por determinismo.

### A.5 — Slash commands & session resume
- Slash: embebidos en `content` (A.3). `/help`, custom, skills expanden; los de diálogo interactivo no.
- Resume: `--resume <session_id>` (específico) o `--continue` (más reciente en cwd); `--session-id
  <uuid>` fija un id elegido en la **primera** corrida (job addressable/idempotente) → después
  `--resume`. Estado persiste como JSONL en `~/.claude/projects/<encoded-cwd>/<session_id>.jsonl`;
  lookup scoped al proyecto + sus worktrees (importa para las sesiones per-worktree de dev-studio).
  **Gotcha:** el resume puede acuñar un session_id **nuevo** en algunos contextos
  ([#12235](https://github.com/anthropics/claude-code/issues/12235)) — leer siempre el id del evento
  `init`/`result`, no asumirlo estable. En sesión persistente el resume se usa para **rehidratar** tras
  restart de la app.

## B) Renderizar la experiencia DUAL

### Opción 1 — Transcript estructurado (componentes React del stream-json)
Parsear eventos → render de turnos chat-log, tool-cards, retry toasts, footer de costo.
- **Pros:** control total del look "Making of" (bullet `● @handle — rol`, cuerpo indentado, cards de
  rama, tool-cards con inputs/outputs, roster/status). Visualización rica de tools (diff colapsable,
  paths, cmd Bash + output). Copy/scroll/search "gratis" como DOM. **ES la técnica de Anthropic** — el
  faux terminal es un transcript renderizado → Opción 1 *es* el norte visual.
- **Cons:** no es una grilla de caracteres — sin fidelidad ANSI/TUI cruda (spinners, arte de cursor,
  box-drawing) salvo que también renderices ANSI.

### Opción 2 — PTY real → xterm.js
PTY en Go (`github.com/creack/pty`: `pty.Start(exec.Command("claude"))` → `*os.File` master), stream de
bytes master por WebSocket a `xterm.js`, teclas de vuelta (frames binarios), resize como control JSON
(`{"type":"resize","cols","rows"}` → `pty.Setsize`). Receta estándar de web-terminal.
- **Pros:** autenticidad literal — el TUI real de `claude`, ANSI real, cursor real, spinners reales.
- **Cons decisivos para dev-studio:** (1) **perdés eventos estructurados** — el PTY lleva ANSI
  renderizado, no JSON → sin tool-cards/roster/cost/modales salvo que TAMBIÉN corras stream-json
  (Opción 3b). (2) **Contradice el driver committeado** (pipe stream-json, un proceso por sesión — no un
  TTY). (3) Los permission prompts del TUI se vuelven `y/n` crudo que hay que screen-scrapear (el hack
  bswen literalmente hace `pty.write('y\n')`). (4) PTY cross-platform en Windows necesita ConPTY.

### Opción 3 — Híbrido (RECOMENDADO, variante 3a)
**Recomendado = 3a: driver stream-json + faux terminal fiel.** Un proceso `claude` en stream-json (A.1).
Renderizar el *mismo* stream de dos formas desde una sola fuente de verdad:
- **Superficie estructurada** (chat cards, tool roster, modales de aprobación, footer de costo) ←
  parsear `assistant`/`user`/`result`/eventos de control.
- **Superficie terminal** ← `xterm.js` en **modo controlado (detached)** — *no* pegado a un PTY.
  `term.write()` del texto assistant streameado (de `stream_event` text deltas → cadencia de typing
  instantánea) + tus propias líneas de prompt faux (`❯ `, prompts de rama, loader `+ Remembering…` en
  teal). Teclas: interceptar `term.onData`, eco local, en Enter empaquetar la línea como turno
  stream-json (A.3). Slash commands tipeados acá pasan en `content`. Podés emitir un set chico de ANSI
  (color, bold, acentos teal) vos mismo para la textura "terminal real" sin TTY real.

Esto da: el feel literal de *tipear-en-un-terminal* + el doble input (composer O terminal → mismo
stdin) + cards estructuradas + aprobaciones programáticas — **desde un proceso, cero PTY**. Matchea el
driver Y el norte visual (faux por naturaleza). AgentsRoom *finge* el doble input; dev-studio lo hace
real porque ambos inputs escriben el mismo `stdin`.

**¿Un proceso alimenta ambas?** Sí — un proceso stream-json alimenta ambas superficies (renderizás los
eventos dos veces). NO necesitás PTY ni segundo proceso para la experiencia recomendada.

**3b (solo si después querés ANSI/TUI literal):** correr un **segundo** `claude` en PTY real como tab
"terminal crudo", en paralelo al driver stream-json. Doble costo/memoria/complejidad, dos sesiones que
reconciliar. Diferir; NO default.

### Frontend específico (todo dentro del teal único PRENTER)
- **xterm.js:** `@xterm/xterm` + addons `@xterm/addon-fit` (tamaño; `fit()` en resize/mount),
  `@xterm/addon-webgl` (renderer GPU, gran win para streaming; fallback canvas si no hay WebGL2),
  `@xterm/addon-web-links`. Theme: `background:#08090a`, superficies `#0c1110`/`#141a19`,
  `fontFamily:"JetBrains Mono"` (ya vendorizada → match exacto Anthropic), `fontSize:13`, cursor blink.
  Cap `scrollback:5000`.
- **Typing/streaming:** ya tenés cadencia de token real con `--include-partial-messages` `text_delta` —
  `term.write(delta.text)` conforme llegan *es* el stream auténtico (nada de typewriter fake). Para las
  cards estructuradas, espejar los mismos deltas a un text node. El typewriter scripted (estilo "Making
  of") solo para replay de un transcript cerrado.
- **Command palette / slash autocomplete:** sembrar de `system/init` (`slash_commands`/`tools`) → refleja
  el arnés + comandos custom instalados REALES. Overlay arriba del composer/terminal; al elegir, insertar
  `/cmd `. Teal único: fila activa/seleccionada = highlight teal (mapear el periwinkle de selección de
  Anthropic → teal), coral del `●`/loader → teal. **Descartar el multi-accent per-cluster de AgentsRoom**
  (viola RN de acento único); rol/status por forma+label, no arcoíris.
- **Window-chrome macOS:** shell React estático (ventana near-black redondeada, 3 traffic-lights, header
  con avatar estilo-Clawd + título, strip de tabs/chapters si querés el framing "reader"). CSS puro.
- **Prompts de rama** (`❯ 1. Yes / 2. Skip` + hint): renderizar desde la capa de aprobación/control
  (A.4) como filas seleccionables estilo terminal — acá el `can_use_tool` se vuelve un prompt
  interactivo hermoso en vez de un `y/n` crudo.

## C) Escalera de autenticidad (simulado → literalmente-real)
| # | Enfoque | Fidelidad | Esfuerzo | Fit dev-studio |
|---|---|---|---|---|
| 0 | React chat-log estilizado, canned/sin stream | Cosmético | Muy bajo | ✗ no real |
| 1 | **Faux terminal = React de eventos stream-json REALES** (técnica exacta "Making of") | Alta *visual* (sin grilla ANSI) | Bajo–med | ✓ matchea norte; feel "typing-en-terminal" débil |
| 2 | **xterm.js modo controlado, driven por stream-json** — deltas reales `write()`, teclas → turnos stream-json, ANSI-accents self-emitidos ← **RECOMENDADO (3a)** | Alta visual **+ feel real de typing + cards + aprobaciones en vivo**, un proceso, sin PTY | Med | ★ mejor autenticidad-por-esfuerzo; honra el driver |
| 3 | **PTY real → xterm.js** (TUI `claude` literal) | TTY/ANSI literal | Med–alto | △ pierde eventos; contradice driver; ConPTY Windows; `y/n` scrapeado |
| 4 | **PTY real + proceso stream-json paralelo** | Máxima | Alto | △ 2 sesiones, 2× costo/memoria, reconciliación — diferir |

**Elegir #2 (Híbrido 3a).** Máximo peldaño que (a) reusa el driver stream-json verbatim, (b) reproduce
el look faux-terminal + cadencia de typing *auténticamente* (deltas reales), (c) entrega el doble input
real que AgentsRoom finge, (d) convierte `can_use_tool` en UI de prompt de rama de primera clase. Subir
a #3/#4 solo si aparece necesidad concreta de ANSI/TUI literal — y aún así #2 como default.

**Antes de construir UI sobre el protocolo de control: 1 probe** que spawnea
`claude --input-format stream-json --output-format stream-json --verbose --permission-prompt-tool stdio`,
manda un turno que dispara un tool gated, y loguea stdin/stdout crudo — el envelope exacto es
version-dependiente y sub-documentado (A.4).

## Fuentes
**Oficial:** [headless](https://code.claude.com/docs/en/headless) ·
[streaming-output / StreamEvent](https://code.claude.com/docs/en/agent-sdk/streaming-output) ·
[streaming-vs-single-mode](https://code.claude.com/docs/en/agent-sdk/streaming-vs-single-mode) ·
[permissions / canUseTool](https://code.claude.com/docs/en/agent-sdk/permissions) ·
[sessions/resume](https://code.claude.com/docs/en/sessions) ·
[streaming events API](https://platform.claude.com/docs/en/build-with-claude/streaming)
**Comunidad / reverse-engineering:**
[Roasbeef Go SDK cli-protocol.md](https://github.com/Roasbeef/claude-agent-sdk-go/blob/main/docs/cli-protocol.md) ·
[claude-cli-agent-protocol SKILL](https://raw.githubusercontent.com/NeverSight/skills_feed/refs/heads/main/data/skills-md/bohdan-shulha/skills/claude-cli-agent-protocol/SKILL.md) ·
[takopi cheatsheet](https://takopi.dev/reference/runners/claude/stream-json-cheatsheet/) ·
[Background Claude stream-json](https://backgroundclaude.com/blog/stream-json) ·
[bswen custom UI](https://docs.bswen.com/blog/2026-03-21-stream-json-custom-ui-claude-code/) ·
[Avasdream wrapping the CLI](https://avasdream.com/blog/claude-cli-agentic-wrapper)
**Issues abiertos (caveats):** [#24594](https://github.com/anthropics/claude-code/issues/24594) ·
[#24612](https://github.com/anthropics/claude-code/issues/24612) ·
[#34046](https://github.com/anthropics/claude-code/issues/34046) ·
[#12235](https://github.com/anthropics/claude-code/issues/12235)
**Frontend:** [xterm.js](https://xtermjs.org/) · [repo+addons](https://github.com/xtermjs/xterm.js/) ·
[react-xtermjs](https://www.qovery.com/blog/react-xtermjs-a-react-library-to-build-terminals) ·
[creack/pty](https://github.com/creack/pty) ·
[TypeIt streaming](https://macarthur.me/posts/streaming-text-with-typeit/)
