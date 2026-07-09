# 04 · Mapa del código actual (file:line) — grounding para implementar

> Recon read-only del repo (subagente 1). Dónde tocar. Paths relativos a la raíz del repo.

## 1. Driver Claude Code CLI-nativo
**`internal/adapters/agent/claudecode/conductor.go`** — ÚNICO paquete que parsea stream-json (boundary
`conductor-no-parsea-jsonl`).
- **Resolución bin** L33-60: `New(bin)` → `resolveClaudeBin()`: `exec.LookPath("claude")` + fallbacks
  (`~/.local/bin/claude`, `~/bin`, `~/.npm-global/bin`, `/usr/local/bin`, `/opt/homebrew/bin`).
- **Flags** `buildArgs()` L64-87:
  ```go
  args := []string{"-p", "--input-format", "stream-json", "--output-format", "stream-json", "--include-partial-messages", "--verbose"}
  if opts.ReadOnly { args = append(args, "--permission-mode", "plan") }   // exploración PB-27
  if opts.Resume != "" { args = append(args, "--resume", opts.Resume) }
  for _, d := range opts.PluginDirs { args = append(args, "--plugin-dir", d) }   // arnés PB-25
  if opts.SystemPrompt != "" { args = append(args, "--append-system-prompt", opts.SystemPrompt) }
  ```
- **Spawn** `Spawn()` L89-118: `exec.CommandContext`, `cmd.Dir = opts.Cwd`, wires Stdin/Stdout/StderrPipe,
  `cmd.Start()`, `go s.logStderr(stderr)` + `go s.pump(stdout)`. Un proceso por sesión; events channel buf 64.
- **Write stdin** `Send()` L142-160: marshala `userFrame{Type:"user", Message{Role:"user",
  Content:[{type:"text",text}]}}` como NDJSON + `\n` → `s.stdin.Write`. Mutex-guarded. `Close()` L164-170
  cierra stdin (EOF) + `cmd.Wait()`.
- **Parse** `pump()` L202-216: `bufio.NewReaderSize(stdout, 1<<20)` (1MB, no Scanner — init frame > 64KB)
  `ReadBytes('\n')`. `translate()` L218-236 mapea: `system/init` → `EventInit{ClaudeSessionID, Model}`;
  `stream_event`+`content_block_delta`+`text_delta` → `EventDelta{Text}`; `result` → `EventResult`/`EventError`.
  **Todo lo demás ignorado (`default: return false`) ← acá caen tool_use/tool_result/assistant/control.**

**Contrato de evento** `internal/ports/agent.go`: `AgentEventKind` {init,delta,result,error} L9-14;
`SpawnOpts{Resume,Cwd,ReadOnly,PluginDirs,SystemPrompt}` L26-38; `AgentSession{Send,Events,Close}` L41-45;
`AgentPort.Spawn` L49-51.

**Cadena Eventos → frontend:**
1. `internal/usecase/session_service.go` `consume()` L240-279 draina `live.Events()`, acumula deltas en
   `assembling strings.Builder`, `publish()`ea `dockFrame{session_id, kind, text, status,
   claude_session_id, model}` (L30-37) como JSON vía `pub.Publish("dock", data)` L281-290. `Turn()`
   L182-207 spawn-on-first-turn, enforce un-turno (`ErrBusy` si streaming).
2. **SSE** `internal/adapters/transport/sse/broker.go` — broker multiplexado hand-rolled, un `/events`
   para TODAS las sesiones (cliente filtra por session_id). `Publish()` L31-42 dropea subs lentos (sin
   Last-Event-ID replay — "TBD" = PB-18). `ServeHTTP` L59-85 escribe `event: dock\ndata: ...\n\n`.
3. **HTTP** `internal/adapters/transport/http/router.go`: `POST /api/sessions/{id}/turn` L34 (→ `sessions.go`
   `sessionTurn` L233-252, 202 Accepted, 409 busy); `GET /events` L57. `withLocalOnly` L67-83 confina localhost.
4. `cmd/dev-studio/main.go` wire: `broker` L52, `claudecode.New` L53, mux `/api/`, `/events`, `/` (SPA) L97-100.

> **stdin es server-owned.** El frontend nunca toca claude stdin directo — POSTea texto a `/turn`, y
> `ccSession.Send` es el ÚNICO escritor. No hay segundo input path hoy.

## 2. Frontend conversación (React 19 + Zustand v5 + Tailwind v4, `go:embed`)
- **Embed**: `web/embed.go` L6 `//go:embed all:dist` → `DistFS`; served por `uiHandler()` main.go L117-133
  (SPA fallback index.html).
- **Root**: `web/src/App.tsx` — layout `ReposRail | StudioNav | (SessionView + ChangesPanel + overlays)`.
- **Componente conversación**: `web/src/widgets/session-view/ui/session-view.tsx` (252 líneas):
  - Header L92-114 (rol/cwd/model/branch chips).
  - Conversación L117-145: mapea `session.conv` bubbles (user der / assistant izq, `whitespace-pre-wrap`) +
    bubble live `streamBuffer` L135-144.
  - **Composer** L148-232: `<textarea>` redimensionable (drag L150-154, `composerH` 64-320px),
    Enter=send / Enter-while-streaming=queue L158-164, 7 botones icono stub L170-183, Encolar + Send L190-208.
  - **Fila "Terminal" = STUB** L210-230: label estático `⌨ Terminal` ("llega con PB-09") + botones
    Comandos/SSH que abren modales explainer. **No hay terminal real.**
- **Stores** (Zustand): `web/src/shared/store/sessions-store.ts` — `sessions`, `activeId`,
  `streamBuffer: Record<id,string>`, `queue: Record<id,string[]>`. `sendTurn` L72-94 (bubble optimista +
  POST /turn), `enqueue` L96-103 (cola client-side RN-8), `onDock` L105-176 reducer SSE: `delta` append a
  streamBuffer L122-128, `result`/`error` flush buffer→conv L129-156, dispatch próximo encolado L167-175.
  Otros: `repos-store.ts`, `arneses-store.ts`, `ui-store.ts`.
- **SSE client**: `web/src/shared/api/sse.ts` — `connectDock()` un `new EventSource("/events")`,
  `addEventListener("dock", ...)`.
- **REST client**: `web/src/shared/api/client.ts` — `api.turn(id,text)` L43-46 POST `/api/sessions/{id}/turn`.
- **Types**: `web/src/shared/api/types.ts` — `Session.conv: Turn[]` L29, `DockFrame` L94-101.
- Sin render terminal-like en ningún lado (solo chat-bubble). Átomos en `web/src/shared/ui/*`.

## 3. Tokens PRENTER
`specs/shell/tokens/` — `colors.css`, `fonts.css`, `spacing.css`, `typography.css`.
- **colors.css**: dark-first, teal único acento saturado. `--teal-500: #00b7aa` (L17, brand primary);
  `--dark-bg:#000000` L41, `--dark-surface:#0c1110` L42, `--dark-raised:#141a19` L43,
  `--dark-border:#1f2826` L44, `--dark-border-teal:rgba(0,183,170,0.35)` L45; `--focus-ring:var(--teal-400)` L81.
- **Tema aplicado** `web/src/app/styles/theme.css`: dark `--background:#08090a` L31, `--primary:#1fc6b8`
  (teal-400 sobre dark) L37; light `--background:#f6f8f8`, `--primary:#009d92` L99-105.
- **Fonts**: `fonts.css` marca = Coco Gothic/Sansation, sustituidas por **Jost** (display) + **Mulish**
  (body); mono = **JetBrains Mono** L19-22. Vendorizadas: `web/src/app/fonts/{jost-latin,mulish-latin,
  jetbrains-mono-latin}.woff2`, `@font-face` real en theme.css L7-27; `--font-mono:"JetBrains Mono"…` L95.

## 4. Backlog / capabilities relevantes (verbatim)
- **PB-24** (BACKLOG L81): ventana nativa Tauri-2-shell-tonto (daemon Go = sidecar) vs Wails + cross-compile
  + firma + updater. ⚠ Antes de vender: Consumer Terms Anthropic. Bloqueada: rebanadas F2+.
- **PB-09** (L56): "Vista Sesión completa: terminal embebido, composer (7 acciones, cola de mensajes,
  redimensionable)" → CAP-02/03 · spec clon §10.1.
- **PB-10** (L57): Tab Cambios + adjuntar conversación. **PB-11** (L58): Diff conversacional.
- **PB-19** (L71): boundary `conductor-encapsula-stream-json`.
- **CAP-02** (INCREMENTO L15): "Sesiones Claude Code multisesión — driver CLI-nativo propio: un subproceso
  claude por sesión, stream-json stdin/stdout, un turno a la vez + cola; sesión con ROL inyecta su arnés al
  spawn (--plugin-dir + system prompt compuesto)".
- **SPEC.md** (`specs/shell/SPEC.md`): L120 "Fila terminal: tab Terminal, +, Comandos, SSH | STUB |
  terminal PTY = PB-09"; L292 "Terminal PTY embebida + Comandos + SSH | PB-09 / PB-17".

## 5. xterm.js / node-pty / PTY
**Nada existe.** `web/package.json` deps = solo `clsx, react, react-dom, tailwind-merge, zustand`. `go.mod`
= solo `gopkg.in/yaml.v3` (sin creack/pty). Solo referencias forward-looking (stubs): `web/src/shared/mock/
backlog.ts` L24 `hx-spike-pty` "Spike: evaluar librerías PTY"; `session-view.tsx` L210 "terminal PTY = PB-09".

## GAP SUMMARY
1. EXISTE: driver CLI-nativo completo (`conductor.go`) — un `claude` por sesión, stream-json stdin/stdout,
   flags de arnés, plan-mode read-only.
2. EXISTE: pipeline server-owned — `SessionService.Turn`→`ccSession.Send` (único escritor stdin) → SSE
   broker `/events` → Zustand `onDock` → chat-bubble `SessionView`.
3. EXISTE: tokens PRENTER (teal `#00b7aa`/app `#1fc6b8`, dark `#000`/`#08090a`) + JetBrains Mono/Jost/Mulish
   woff2 vendorizadas, dark-first.
4. FALTA: terminal real — sin xterm.js/node-pty, sin backend PTY; fila "Terminal" = stub (PB-09), sin
   superficie de render terminal.
5. FALTA: doble input — un solo path (textarea → POST /turn → Send server-side); sin segundo canal
   escribiendo el mismo stdin, y stdin no expuesto al cliente.
6. FALTA: decisión de shell nativo (PB-24 Tauri vs Wails, bloqueada F2+) — hoy HTTP localhost + go:embed
   SPA; delta events streamean solo texto (sin ANSI/tool-frame passthrough para render terminal-auténtico).
