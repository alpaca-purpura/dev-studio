# 03 · Interop multi-CLI — hostear varios CLIs agénticos tras un shell-tonto (brief completo)

> Detalle crudo de la investigación (subagente 3). Fuente para refinar/implementar sin re-investigar.
> Anclado en los ports hexagonales del repo (`internal/ports/agent.go`, `conductor.go`, `arnes.go`).
> UX objetivo (ver `00-visual-reference.md`): pane terminal REAL + composer chat estructurado
> escribiendo al MISMO proceso, roster multi-agente estilo mensajería, prompts de rama (permisos).

## 0. Veredicto rápido
| CLI | ¿Eventos estructurados? | Mejor transporte | ¿Fit UX chat rica? |
|---|---|---|---|
| **Claude Code** | Sí (NDJSON stream-json + envelope de control) | stdio stream-json (ya implementado) | Full |
| **Amp** (Sourcegraph) | Sí — **compatible con Claude Code** stream-json | stdio stream-json (near drop-in) | Full |
| **Cursor** (`cursor-agent`) | Sí (NDJSON stream-json) | stdio stream-json | Full (turnos one-shot-ish) |
| **Gemini CLI** | Sí (NDJSON stream-json) **o** ACP | stdio stream-json; ACP para aprobaciones en vivo | Full (aprobaciones necesitan ACP) |
| **OpenAI Codex** | Sí (`exec --json` JSONL) **o** JSON-RPC app-server | stdio JSON-RPC `codex app-server` | Full (bidi necesita app-server) |
| **OpenCode** | Sí (HTTP REST + SSE, **Go SDK oficial**) | HTTP server `opencode serve` | Full (transporte distinto) |
| **Aider** | **No** — solo texto crudo de terminal | spawn en PTY, diff del worktree vos mismo | **Terminal-only fallback** |

Cinco de seis emiten eventos estructurados → UX de mensajería completa. **Aider es el único raw-text** —
constriñe ese workspace al pane terminal + composer fino "escribir a stdin"; sin tool-cards, sin modales
de permiso, sin roster streaming (infierís busy/idle del estado del proceso + scrapeás el diff del worktree).

## 1. Briefs por proveedor

### 1.1 Claude Code (`claude`) — baseline, ya wired
- **Headless:** `claude -p --input-format stream-json --output-format stream-json --include-partial-messages
  --verbose` (lo que `conductor.go` spawnea hoy). `--permission-mode plan` read-only; `--resume <id>` /
  `--continue`; `--plugin-dir` + un `--append-system-prompt` (inyección del arnés).
- **Output:** NDJSON discriminado por `type`: `system`(`init`), `assistant` (envuelve `BetaMessage`;
  tool calls como bloques `tool_use`), `user` (bloques `tool_result`), `stream_event` (deltas), `result`.
- **Bidireccional:** sí — `--input-format stream-json` mantiene vivo el proceso, lee frames NDJSON
  `{"type":"user","message":{...}}`. Debajo, envelope de control (`control_request`/`control_response`):
  `interrupt`, `set_permission_mode`, handshake `can_use_tool`.
- **Permisos wire-level:** CLI escribe `control_request` con `request.subtype:"can_use_tool"` a **stdout**;
  host responde en **stdin** con `control_response` (eco `request_id`, `{behavior:"allow"|"deny"}`). El
  `canUseTool` del SDK es wrapper de esto. Orden fijo: hooks → deny → ask → mode → allow → `canUseTool`.
- **Resume:** `session_id` (UUID) en `init` y cada `result`; `--resume`/`--continue`/`--fork-session`.
- Docs: https://code.claude.com/docs/en/headless · https://code.claude.com/docs/en/cli-reference ·
  https://code.claude.com/docs/en/agent-sdk/streaming-output · https://code.claude.com/docs/en/agent-sdk/permissions

### 1.2 Amp (Sourcegraph, `amp`) — el gemelo más cercano
- **Headless:** `amp -x "<prompt>"` / `--execute` (también `echo … | amp -x`).
- **Stream:** `--stream-json` (requiere `--execute`) → NDJSON; `--stream-json-thinking` agrega thinking.
  El manual dice que el formato es **"compatible con Claude Code"**: mismo `system`(`init` → `cwd`,
  `session_id`, `tools`, `mcp_servers`), `user`, `assistant` (`text`/`tool_use`/`thinking`), `result`
  (`subtype: success|error_during_execution|error_max_turns`, `duration_ms`, `is_error`, `num_turns`,
  `usage`, `permission_denials`). Subagents vía `parent_tool_use_id`.
- **Bidireccional:** sí — `--stream-json-input` lee un mensaje JSON por línea de stdin
  (`{"type":"user","message":{"role":"user","content":[{"type":"text","text":"…"}]}}`, opcional
  `"steer":true`). Análogo directo de `--input-format stream-json` de Claude.
- **Permisos:** default post-rebuild = **sin prompts de aprobación**; `--dangerously-allow-all` legacy;
  gating custom vía hook Plugin `tool.call`. (Headless Amp da *visibilidad* de tool calls pero no un
  handshake approve/deny out-of-the-box.)
- **Resume:** `amp threads continue <threadId> -x "…" --stream-json`.
- **SDK:** `@ampcode/sdk` (TS+Python) — `execute()` async iterable espejando los tipos del CLI.
- Docs: https://ampcode.com/manual · https://ampcode.com/manual/appendix · https://ampcode.com/news/streaming-json

### 1.3 Cursor CLI (`cursor-agent` / `agent`)
- **Headless:** `agent -p "<prompt>"` / `--print`; `--output-format` solo con `-p`.
- **Stream:** `--output-format {text|json|stream-json}`. `stream-json` = NDJSON: `system`/`init`
  (`apiKeySource`, `cwd`, `session_id`, `model`, `permissionMode`), `assistant` (`message.content[]`),
  tool events (`subtype:"started"|"completed"`, `call_id`, payloads tipados `readToolCall`/`writeToolCall`),
  `result` (`subtype:"success"`, `is_error`, `duration_ms`). `--stream-partial-output` (con `-p`) agrega
  deltas char/chunk. **En fallo el stream puede terminar sin `result`** (exit nonzero + stderr) — manejar
  defensivo.
- **Bidireccional:** no documentado protocolo persistente stdin-JSON; multi-turno vía `-p` repetido +
  `--resume <chatId>` / `--continue` / `create-chat` (devuelve chat ID). Verificar `cursor-agent --help`.
- **Permisos:** **no** auto-aprobado por default. Bypass: `-f`/`--force`/`--yolo`, `--approve-mcps`,
  `--trust` (headless-only), `--sandbox`. Sin canal de intercepción per-call — pre-autorizar. Caveat beta.
- **API server:** el CLI local no tiene server mode; la Cloud Agents API (`api.cursor.com/v1/agents`) es
  para agentes remotos — fuera de scope de un shell-tonto local.
- Docs: https://cursor.com/docs/cli/reference/output-format · https://cursor.com/docs/cli/reference/parameters · https://cursor.com/docs/cli/headless

### 1.4 Gemini CLI (`gemini`)
- **Headless:** `gemini -p "<prompt>"`; `-i` = corre luego cae a interactivo. Exit 0/1/42(input)/53(turn-limit).
- **Stream:** `--output-format {text|json|stream-json}`. `json` = objeto final único. `stream-json` = JSONL:
  `init` (session id, model), `message`, `tool_use` (call+args), `tool_result`, `error`, `result` (stats +
  token usage por modelo). Tools MCP vía mismos `tool_use`/`tool_result`. ⚠ Footgun (issue #9281): un tool
  error no-fatal puede abortar todo con exit 54 *solo* bajo `--output-format json`.
- **Permisos:** `--approval-mode {default|auto_edit|yolo|plan}`, `--sandbox`, + **Policy Engine**
  (`~/.gemini/policies/*.toml`, reglas → `allow|deny|ask_user`, scopable `interactive:false`). **Crucial:
  en headless `ask_user` degrada a `deny`** — no podés interceptar aprobación en vivo por stdio plano;
  solo pre-configurar.
- **Intercepción en vivo → modo ACP:** `--acp` (a.k.a. `--experimental-acp`) expone **JSON-RPC 2.0 sobre
  stdio** (Agent Client Protocol): `initialize`, `newSession`, `loadSession`, `prompt`, `cancel`,
  `setSessionMode` (cambia nivel de aprobación mid-session — el mecanismo de intercepción en vivo), +
  **filesystem proxy** (host gatea reads/writes) + MCP bidireccional. Esta es la superficie "server"
  host-drivable.
- **Resume:** `-r/--resume [latest|<index>|<uuid>]` (headless), `--list-sessions`, `--delete-session`.
- Docs: https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/headless.md ·
  https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/acp-mode.md ·
  https://github.com/google-gemini/gemini-cli/blob/main/docs/reference/policy-engine.md

### 1.5 OpenAI Codex (`codex`)
- **Headless one-shot:** `codex exec "<prompt>"` (mensaje final a stdout, progreso a stderr; default
  sandbox read-only). `codex exec --json` → **JSONL** de `ThreadEvent`: `thread.started`{thread_id},
  `turn.started`, `item.started/updated/completed`{item}, `turn.completed`{usage}, `turn.failed`, `error`.
  `ThreadItem.details`: `agent_message`, `reasoning`, `command_execution`(command, aggregated_output,
  exit_code, status), `file_change`(changes[], status), `mcp_tool_call`, `web_search`, `todo_list`.
  `--output-schema` constriñe el mensaje final. (SSoT schema: `codex-rs/exec/src/exec_events.rs`.)
- **Bidireccional / server:** `codex exec` es one-shot (encadenar con `codex exec resume --last "<next>"`).
  Para canal persistente host-driven usar **`codex app-server`** — JSON-RPC 2.0 sobre stdio (o
  `--listen unix://`): `initialize`, `thread/start|resume|fork|list`, `turn/start|steer|interrupt`,
  `item/*`, `fs/*`, `command/exec`, `process/spawn|writeStdin|resizePty` (¡PTY-capable!). Codegen de tipos
  vía `codex app-server generate-ts|generate-json-schema`. También `codex mcp-server` (dialecto MCP
  estándar). **Ya no hay `codex proto` stdin/stdout crudo** — removido en la reescritura Rust.
- **Aprobaciones/sandbox:** `--sandbox {read-only|workspace-write|danger-full-access}`,
  `--ask-for-approval {untrusted|on-request|never}`, `--dangerously-bypass-approvals-and-sandbox`
  (`--yolo`). **`codex exec` hardcodea `approval_policy: Never`** — headless el único lever es el sandbox
  tier. **Intercepción en vivo solo vía los servers:** app-server manda
  `item/commandExecution/requestApproval` / `item/fileChange/requestApproval` (cliente responde
  `accept|acceptForSession|decline|cancel`); mcp-server manda `execCommandApproval`/`applyPatchApproval`.
- **Resume:** `codex exec resume [<id>|--last]`, `codex resume`, `codex fork`; sesiones en `$CODEX_HOME`.
  TS SDK `@openai/codex-sdk`.
- Docs: https://github.com/openai/codex/blob/main/codex-rs/app-server/README.md ·
  https://github.com/openai/codex/blob/main/codex-rs/exec/src/exec_events.rs · https://learn.chatgpt.com/docs/non-interactive-mode

### 1.6 OpenCode (`opencode`, repo `anomalyco/opencode`, ex `sst/opencode`)
- **Server mode (interfaz estructurada primaria):** `opencode serve [--port 4096] [--hostname 127.0.0.1]`
  → **HTTP REST + SSE**, sobre Hono con **OpenAPI 3.1 en `GET /doc`**. Rutas: `POST /session` (crear,
  `parentID` para fork), `GET/DELETE /session/:id`, `session.prompt` (soporta `format:{type:"json_schema"}`),
  `GET /event` (SSE proyecto), `GET /global/event`, `POST /permission/:requestID/reply`, `/pty`, `/config`.
  Auth `OPENCODE_SERVER_PASSWORD`.
- **Event stream:** SSE `{"type":"<entity>.<action>","properties":{…}}`. Incluye `session.created/updated/
  idle/error`, `message.updated`, `message.part.updated/removed`, `permission.asked/replied`,
  `tool.execute.before/after`, `file.edited`, `command.executed`, `todo.updated`. Mensajes = **Parts**
  discriminadas: `TextPart`, `ReasoningPart`, `FilePart`, **`ToolPart`** (FSM `pending → running →
  completed|error`, campo `raw` acumula input streameado), `SubtaskPart`. SSE por sesión no es contrato
  cerrado — subscribir `/event` y filtrar por `sessionID`.
- **SDKs oficiales (¡3 lenguajes incl. Go!):** `@opencode-ai/sdk` (JS), **`github.com/sst/opencode-sdk-go`**
  (Go — directamente relevante a un host Go), `opencode-ai` (Python). Clientes generados finos sobre `/doc`.
- **Headless one-shot:** `opencode run "<msg>" --format json`, `-c`/`--continue`, `-s <id>`, `--fork`,
  `--auto` (auto-aprueba no-denegados), `--attach <url>`. `opencode export <id>` dumpea sesión JSON.
- **Permisos programáticos:** server emite `permission.asked` (tool, input, session/dir) → host responde
  `POST /permission/:requestID/reply` con `once|always|reject` → `permission.replied`. Default vía
  `opencode.json` `permission` map (`allow|ask|deny`, last-match-wins). **Pinear versión** — nombres han
  churneado (`permission.updated`→`permission.asked`).
- **Provider-agnostic:** 75+ providers + local (Ollama, LM Studio, llama.cpp) + cualquier OpenAI-compatible.
- Docs: https://opencode.ai/docs/server/ · https://opencode.ai/docs/sdk/ · https://opencode.ai/docs/cli/ ·
  https://opencode.ai/docs/permissions/ · https://github.com/sst/opencode-sdk-go

### 1.7 Aider — el outlier raw-text
- **Scripting:** `aider --message "<msg>"` / `-m` (o `--message-file`/`-f`) → un turno, aplica edits a
  disco, sale.
- **No existe output estructurado.** Verificado contra `aider/args.py`: **no hay flag `--json`**.
  `--stream`/`--no-stream` togglea impresión token-a-token *de terminal* (no event stream); `--pretty`
  togglea ANSI. El único "JSON" es analytics PostHog opt-in, no un canal de edits/eventos.
- **No es agente de tool-call:** parsea bloques SEARCH/REPLACE del texto del modelo y los aplica. Los
  edits solo aparecen como diffs impresos + `Applied edit to <path>`. El host debe **diffear el worktree
  git** antes/después de cada corrida para saber qué cambió.
- **Python API** (`from aider.coders import Coder; coder.run("…")`) explícitamente inestable, **in-process
  only** (embeber CPython — no protocolo de subproceso).
- **Permisos:** blanket `--yes-always`; intercepción per-call solo subclaseando `InputOutput` vía la API
  Python in-process. `--auto-commits` (default on) = único registro semi-estructurado.
- **Sin daemon/server/socket.** "Resume" = re-parsear `.aider.chat.history.md` en proceso fresco vía
  `--restore-chat-history`.
- **Veredicto:** tratar como programa de terminal tonto. Adapter = spawn `aider -m … --yes-always` en PTY,
  snapshot del worktree para "files changed" derivado, scrapear stdout para el transcript.
- Docs: https://aider.chat/docs/scripting.html · https://aider.chat/docs/config/options.html · https://github.com/Aider-AI/aider/blob/main/aider/args.py

## 2. Matriz de normalización
| CLI | Invocación (headless) | Transporte | ¿Estructurados? | ¿Tool-calls? | ¿Permisos host-answerable? | Resume | Fallback terminal |
|---|---|---|---|---|---|---|---|
| **Claude Code** | `claude -p --input-format stream-json --output-format stream-json` | **stdio NDJSON** + control | ✅ `system/assistant/user/stream_event/result` | ✅ bloques `tool_use`/`tool_result` | ✅ **en vivo** — `control_request can_use_tool` ⇄ `control_response` | ✅ `session_id`, `--resume`/`--continue`/`--fork-session` | ✅ PTY |
| **Amp** | `amp -x --stream-json [--stream-json-input]` | **stdio NDJSON** (Claude-compat) | ✅ mismo shape | ✅ | ⚠ visible, sin handshake (default sin prompts; gate vía plugin) | ✅ `amp threads continue <id>` | ✅ PTY |
| **Cursor** | `cursor-agent -p --output-format stream-json` | **stdio NDJSON** | ✅ `system/init`, `assistant`, tool `started/completed`, `result` | ✅ tipado | ⚠ pre-autorizar (`-f`/`--yolo`/`--trust`); sin canal en vivo | ✅ `--resume <chatId>`/`--continue`/`create-chat` | ✅ PTY |
| **Gemini** | `gemini -p --output-format stream-json` | **stdio NDJSON**; o **ACP (JSON-RPC/stdio)** | ✅ `init/message/tool_use/tool_result/error/result` | ✅ (incl. MCP) | ❌ plano (`ask_user`→`deny`); ✅ **en vivo vía ACP** `setSessionMode` | ✅ `-r [latest|id]` | ✅ PTY |
| **Codex** | `codex exec --json`; o `codex app-server` | one-shot **stdio JSONL**; o **stdio JSON-RPC** | ✅ `ThreadEvent`/`ThreadItem`; app-server `item/*` | ✅ `command_execution`/`file_change`/`mcp_tool_call` | ❌ `exec` fuerza `never`; ✅ **en vivo vía app-server** `requestApproval` | ✅ `codex exec resume --last`/`thread/resume` | ✅ PTY (app-server `process/spawn`+PTY) |
| **OpenCode** | `opencode serve` (o `run --format json`) | **HTTP REST + SSE** (Go/JS/Py SDK) | ✅ `message.part.updated` (Part union), `session.*` | ✅ `ToolPart` FSM + `tool.execute.*` | ✅ **en vivo** — `permission.asked` ⇄ `POST /permission/:id/reply` | ✅ REST `session` id/list; `-c`/`-s`/`--fork` | ⚠ vía `/pty` endpoint |
| **Aider** | `aider -m "…" --yes-always` | **texto crudo (PTY)** | ❌ ninguno | ❌ (diffs impresos; diff worktree) | ❌ blanket `--yes-always` | ⚠ re-parse history (`--restore-chat-history`) | ✅ **única** opción |

Leyenda: ✅ primera clase · ⚠ parcial/workaround · ❌ ausente.

## 3. Cuatro familias de transporte (+ un caso degenerado)
- **Familia A — stdio NDJSON ("stream-json"):** Claude, Amp, Cursor, Gemini(stream-json), Codex(`exec
  --json`, one-shot). **Es exactamente lo que el adapter `claudecode` ya hace.** Amp es line-for-line
  compatible; los otros = mismo *patrón*, vocabulario `type` distinto. Un `StreamJSONAdapter` +
  traductor de frames por proveedor cubre 4 de 6.
- **Familia B — stdio JSON-RPC 2.0 (bidireccional, más rica):** Codex `app-server`, Gemini `--acp`.
  Proceso persistente; request/response + notifs; **intercepción de aprobación en vivo** como requests
  server→client; filesystem proxy; control de PTY. Acá viven steering mid-turn + modales reales.
- **Familia C — HTTP + SSE:** solo OpenCode — pero trae **Go SDK oficial**, así que para un host Go es
  la integración rica de *menor esfuerzo* (sin parser hand-rolled). Trade-off: manejás un server local +
  puerto en vez de un pipe stdio; la "autenticidad terminal" viene de su `/pty`, no del stdout del proceso.
- **Familia D — terminal crudo:** Aider. Sin canal estructurado; terminal-only.

**Familia B es un estándar emergente: ACP (Agent Client Protocol, agentclientprotocol.com)** — contrato
cross-vendor JSON-RPC-over-stdio para exactamente esto — Zed ya driftea varios agentes por ACP, Gemini
lo habla nativo. El app-server de Codex es un dialecto paralelo. **Recomendación: dar forma al envelope
normalizado de DevStudio alineado-a-ACP** en vez de inventar uno.

## 4. Interfaz de adapter recomendada (Go, extiende los ports existentes)
El repo **ya tiene el contrato correcto** — `internal/ports/agent.go` define `AgentPort.Spawn →
AgentSession{Send, Events, Close}` con `AgentEvent` normalizado. Solo debe crecer de "Claude-shaped" a
"provider-agnostic + rico":

**Qué es Claude-shaped hoy y debe generalizar:**
- `AgentEvent` solo tiene `init/delta/result/error` + campo `Text`. Para tool-cards + modales de rama
  necesita kinds **tool-call, tool-result, reasoning, permission-request**.
- El campo se llama literal `ClaudeSessionID`. Renombrar a `ProviderSessionID`.
- **No hay canal de permisos** ni **passthrough raw-PTY** — ambos requeridos por la UX.

```go
// Descriptor de capacidades que cada adapter anuncia — drive de degradación graceful de UI.
type Capabilities struct {
    StructuredEvents bool // false => Aider => UI terminal-only
    ToolCallEvents   bool
    LivePermissions  bool // host puede responder approve/deny en runtime
    MultiTurnStdin   bool // proceso persistente acepta turnos nuevos (vs re-spawn por turno)
    RawPTY           bool // expone un byte-terminal auténtico
    Resume           bool
}

type EnvelopeKind string
const (
    EvSessionInit   EnvelopeKind = "session.init"      // session id, model, tools
    EvTextDelta     EnvelopeKind = "text.delta"        // texto assistant streaming
    EvReasoning     EnvelopeKind = "reasoning"         // thinking (opcional)
    EvToolCall      EnvelopeKind = "tool.call"         // {name, input, callID, status}
    EvToolResult    EnvelopeKind = "tool.result"       // {callID, output, status}
    EvFileChange    EnvelopeKind = "file.change"       // {path, kind, diff} (puede derivarse)
    EvPermissionReq EnvelopeKind = "permission.request"// {requestID, tool, input} -> necesita reply
    EvTurnResult    EnvelopeKind = "turn.result"       // cost, usage, duration
    EvError         EnvelopeKind = "error"
    EvRaw           EnvelopeKind = "raw"               // bytes de terminal opacos passthrough
)

type Envelope struct {
    Kind      EnvelopeKind
    Text      string
    Tool      *ToolInfo       // name, callID, input, output, status
    Perm      *PermissionReq  // requestID, tool, input, blockedPath
    SessionID string
    Model     string
    Usage     *Usage          // cost, tokens, duration
    Raw       []byte          // para EvRaw
    Err       string
}

type AgentSession interface {
    Send(ctx context.Context, text string) error                            // turno estructurado (composer)
    ReplyPermission(ctx context.Context, requestID string, allow bool) error// responde el prompt de rama
    Interrupt(ctx context.Context) error                                    // cancela turno en vuelo
    Events() <-chan Envelope
    RawTerminal() (io.ReadWriteCloser, bool)                                // PTY auténtico (ok=false => ninguno)
    Close() error
}

type AgentPort interface {
    Capabilities() Capabilities
    Spawn(ctx context.Context, opts SpawnOpts) (AgentSession, error)
}
```

**Obligaciones por adapter (todos satisfacen el mismo `AgentPort`):**
- `StreamJSONAdapter` (Claude/Amp/Cursor/Gemini-stream-json/Codex-exec): reusar el pump de `conductor.go`;
  swap del `translate()` por proveedor. `Send` escribe frame NDJSON user (o re-spawnea con `--resume`
  para Cursor/Codex-exec sin stdin persistente). `ReplyPermission` escribe `control_response` (Claude) o
  no-op (Amp/Cursor — pre-autorizar). `RawTerminal` = tee del PTY del hijo.
- `JSONRPCAdapter` (Codex app-server, Gemini ACP): cliente JSON-RPC 2.0 stdio; mapea `turn/*`+`item/*`
  (o ACP `prompt`/session updates) a envelopes; `ReplyPermission` responde el `requestApproval`/`ask_user`
  server→client; `Interrupt` = `turn/interrupt`/`cancel`.
- `HTTPAdapter` (OpenCode): wrap `github.com/sst/opencode-sdk-go`; `Spawn` = `opencode serve` +
  `session.create`; `Events` = SSE `/event` filtrado por `sessionID` → envelopes; `ReplyPermission` =
  `POST /permission/:id/reply`; `RawTerminal` vía `/pty`.
- `RawPTYAdapter` (Aider): spawn `aider -m … --yes-always` en PTY; `Capabilities{StructuredEvents:false}`;
  `Events` solo `EvRaw` (+ `EvFileChange` derivado de diff del worktree al salir); `RawTerminal` es todo.

La lógica `--permission-mode plan`/`--plugin-dir`/`--append-system-prompt` de `buildArgs` queda — es
asunto privado del adapter Claude. `SpawnOpts` (Cwd, ReadOnly, Resume, PluginDirs, SystemPrompt)
generaliza; cada adapter mapea `ReadOnly` a su sandbox/plan (Claude `--permission-mode plan`, Codex
`--sandbox read-only`, Gemini `--approval-mode plan`, Cursor `--mode plan`, OpenCode config, Aider
`--dry-run`). **Boundary preservado:** `conductor-no-parsea-jsonl` generaliza a "cada adapter dueño de su
protocolo; el usecase solo ve `Envelope`".

## 5. Relación con el registry ArnesIA (cada proveedor como adapter instalable)
Patrón ArnesIA (`internal/ports/arnes.go`): marketplace git de plugins pineados por versión, lock
committeado `.devstudio/arneses.yaml` (`registry` + `arneses[].{id,version,canal}`), caché plugin en
`~/.dev-studio/arneses/` (rehidratable, modelo npm). Mapea a proveedores **con un caveat**:
- **Transfiere limpio:** el patrón *manifest + version-pinning + lock por-proyecto*. Un
  `.devstudio/drivers.yaml` (espejo de `arneses.yaml`) pinea **qué proveedor + versión** usa un
  workspace, y un **manifest de capacidades** por proveedor (el `Capabilities` de arriba) = la superficie
  auditable additive-only que ArnesIA ya lee. La UI lee capacidades del manifest para decidir qué panes
  mostrar (roster streaming + tool-cards vs terminal-only) — sin branching per-proveedor hardcodeado.
- **NO transfiere:** un arnés es **data** (dirs + prompts cargados vía flags nativos — DH-10 "cero API").
  Un adapter de proveedor es **código** (parser Go de protocolo). No podés shipear un parser Codex/OpenCode
  como plugin de marketplace git como un rol. Dos modelos:
  1. **Adapters compilados + registry declarativo** (recomendado near-term): las 4 familias shipean en el
     binario Go; el registry/lock solo *selecciona y configura* per workspace. Bajo esfuerzo, aprovecha
     los ports hexagonales + el Go SDK de OpenCode.
  2. **Sidecars adapter out-of-process** (true "shell-tonto"): cada adapter = proceso externo que traduce
     proveedor↔el envelope normalizado DevStudio — que es esencialmente **lo que ACP ya es**. Adoptar ACP
     como envelope interno → Gemini gratis + seam plug-in limpio donde proveedores nuevos llegan como
     bridges ACP sin recompilar el shell. Más pesado, pero máxima expresión de "el shell solo conoce el
     protocolo normalizado".

El arnés (rol) y el proveedor (engine) = **dos ejes ortogonales instalables** de un workspace: *qué
modelo de colaborador* (arnés) × *qué runtime CLI* (adapter). Extensión natural del lock existente.

## 6. Relación con PB-24 (Tauri shell-tonto vs Go+webview vs Wails)
PB-24 debate la **ventana nativa**. El diseño de adapters es **ortogonal a esa elección pero refuerza
mantener Go como el core del driver**:
- Toda la capa de adapters es Go (`internal/adapters/agent/*`), y la integración rica más rica —
  **el Go SDK oficial de OpenCode** — es Go-nativa. Codex app-server codegen'ea tipos Go; la familia
  stream-json es un pump NDJSON Go que ya tenés.
- En el modelo Tauri de PB-24, el **daemon Go ES el sidecar** — *es* esta capa de adapters. Tauri/Wails
  solo renderizarían el webview + pasarían bytes stdin/PTY. → el trabajo multi-proveedor es **el mismo
  código pase lo que pase** con la ventana, y argumenta contra una reescritura Rust del adapter (tirarías
  el Go SDK + conductor).
- "Shell-tonto" son **dos shells tontos**: la *ventana* (Tauri/Wails renderiza webview) y el *driver* (el
  shell renderiza envelopes normalizados; proveedores = adapters). Componen: ventana tonta → webview →
  UI-tonto ← envelope normalizado ← adapters de proveedor.
- Constraint concreto para PB-24: **autenticidad terminal exige PTY real**, y el transporte de OpenCode es
  HTTP (su PTY es un endpoint `/pty`, no el stdout del proceso). La ventana ganadora debe reenviar bytes
  PTY a un pane estilo xterm.js — stdio-family lo da nativo, OpenCode vía `/pty`, Aider su propio PTY.

## 7. Mínimo común denominador (hace universal el doble chat+terminal)
El piso universal — presente en **todo** proveedor incl. Aider — son exactamente dos cosas:
1. **Un byte-terminal crudo (PTY passthrough).** Todo CLI se spawnea en PTY y su stdout se streamea a un
   pane xterm.js. Solo esto entrega el frame terminal (window chrome + stream vivo + input `❯`) con
   autenticidad literal, para las 7 superficies. El **composer** chat degrada a "escribir texto al stdin
   del proceso" — que también funciona en todos (Aider `-m`/stdin, stream-json aceptan user frame,
   servers aceptan un prompt call).
2. **Una señal gruesa de liveness** — busy / waiting / idle — derivable sin eventos estructurados: *busy*
   mientras el proceso emite output, *idle* al terminar el turno / salir el proceso, *waiting* al bloquear
   en input o permiso. Suficiente para el **status active/waiting/idle del roster** (panel mensajería)
   para todo proveedor, incl. Aider.

Todo lo más rico = **enhancement gated en `Capabilities.StructuredEvents`**:
| Elemento UX (del visual reference) | MCD (todos) | Enhanced (proveedores con eventos) |
|---|---|---|
| Frame terminal + input `❯` | ✅ PTY real | igual |
| Composer chat "TÚ →" | ✅ escribir a stdin | igual |
| Roster active/waiting/idle | ✅ de liveness del proceso | + token/cost/turn de `turn.result` |
| Bubbles assistant | ⚠ scrapear texto PTY | ✅ de `text.delta`/`assistant` |
| Tool-cards (`[Using Bash…]`) | ❌ | ✅ de `tool.call`/`tool.result` |
| Modal de rama (`❯ 1. Yes / 2. Skip`) | ❌ (solo terminal) | ✅ de `permission.request` ⇄ `ReplyPermission` (Claude, OpenCode, Codex-app-server, Gemini-ACP) |
| Badge de rol (arnés) | ✅ | ✅ |

**Construir el terminal-pane + stdin-composer + liveness-roster como base universal; layerear bubbles,
tool-cards y modales de permiso encima, driven por `Capabilities` del manifest.** Aider (y cualquier
raw-text futuro) renderiza terminal-only en un shell idéntico — sin special-casing, solo un flag.

## 8. Huecos concretos en el adapter `claudecode` ACTUAL (accionable)
Incluso para Claude solo, `conductor.go`'s `translate()` hoy emite solo `init/delta/result/error` y
descarta `assistant`/tool frames. Para llegar a la UX objetivo:
1. **Parsear bloques `tool_use` / `tool_result`** (en frames `assistant`/`user`) → `EvToolCall`/
   `EvToolResult`. Hoy ignorados (`default: return false`).
2. **Manejar el envelope `control_request` `can_use_tool`** → `EvPermissionReq`, e implementar
   `ReplyPermission` que escribe el `control_response` en stdin. Hoy no hay canal de permisos — el modal
   de rama no se puede construir sin esto.
3. **Exponer handle raw-PTY** (`RawTerminal`) para el pane terminal auténtico. Hoy stdout se consume solo
   como NDJSON parseado.
4. **Renombrar `ClaudeSessionID` → `ProviderSessionID`** + método `Capabilities()`.

Estos cuatro = el puente de "streamea texto a un pane" a "el shell mensajería, tool-aware,
permission-answering" del visual reference — y definen el port extendido que todo otro proveedor implementa.

## Fuentes (primarias por proveedor)
- **Claude Code:** https://code.claude.com/docs/en/headless · https://code.claude.com/docs/en/cli-reference · https://code.claude.com/docs/en/agent-sdk/streaming-output · https://code.claude.com/docs/en/agent-sdk/permissions
- **Amp:** https://ampcode.com/manual · https://ampcode.com/manual/appendix · https://ampcode.com/news/streaming-json · https://ampcode.com/manual/sdk
- **Cursor:** https://cursor.com/docs/cli/reference/output-format · https://cursor.com/docs/cli/reference/parameters · https://cursor.com/docs/cli/headless
- **Gemini:** https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/headless.md · https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/acp-mode.md · https://github.com/google-gemini/gemini-cli/blob/main/docs/reference/policy-engine.md
- **Codex:** https://github.com/openai/codex/blob/main/codex-rs/app-server/README.md · https://github.com/openai/codex/blob/main/codex-rs/exec/src/exec_events.rs · https://learn.chatgpt.com/docs/non-interactive-mode · https://github.com/openai/codex/blob/main/codex-rs/docs/codex_mcp_interface.md
- **OpenCode:** https://opencode.ai/docs/server/ · https://opencode.ai/docs/sdk/ · https://opencode.ai/docs/cli/ · https://opencode.ai/docs/permissions/ · https://github.com/sst/opencode-sdk-go
- **Aider:** https://aider.chat/docs/scripting.html · https://aider.chat/docs/config/options.html · https://github.com/Aider-AI/aider/blob/main/aider/args.py
- **Cross-vendor:** https://agentclientprotocol.com/get-started/introduction (ACP — shape recomendado del envelope normalizado)
