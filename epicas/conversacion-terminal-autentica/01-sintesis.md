# Síntesis — conversación DevStudio "terminal auténtica × Anthropic × AgentsRoom"

**Fecha:** 2026-07-09. **Insumos:** `00-visual-reference.md` (screenshots) + 3 investigaciones
(mapa código actual · Claude Code terminal-en-webview · interop multi-CLI). **Naturaleza:**
investigación/diseño. NO se tocó código. Relaciona con **PB-09** (terminal embebido) y **PB-24**
(shell nativo). Idea/redefinición → PB antes de construir (regla de producto · M2: solo `/pm-*`
edita BACKLOG).

---

## 1 · ¿Cómo se hizo la web de Anthropic? (respuesta corta)

**El "terminal" de la making-of NO es un TTY real — es un transcript renderizado con React,
estilizado como terminal.** Un `<div>` con window-chrome macOS que "escribe" un guion. Por eso es
100% replicable en web. Piezas:
- Toggle de modo `>_ Read in Terminal` ↔ `¶ Read as Article` (mismo contenido, dos vistas).
- Chat-log: `● @handle — rol` + cuerpo indentado; loader `+ Remembering...`; efecto typing.
- Input de slash-command real (`/help`) + **prompts de rama** (`❯ 1. Yes / 2. Skip`).
- Imágenes en sub-ventanas anotadas. JetBrains Mono, coral+lavanda sobre `#141413`.

**Consecuencia para nosotros:** lo que Anthropic finge (terminal renderizado), DevStudio lo hace
**de verdad** — ya spawnea `claude` real con stream-json. La estética es la misma; el motor es
auténtico. Eso ES el "lo más real" que pediste.

## 2 · ¿Cómo se logra en DevStudio? Veredicto = **Híbrido 3a**

**Un solo proceso `claude` stream-json → renderizar el MISMO stream de dos formas:**

1. **Superficie estructurada** (cards de chat, roster, tool-cards, modal de aprobación, footer de
   costo) ← parsear eventos `assistant`/`user`/`result`/control.
2. **Superficie terminal** = `xterm.js` en **modo controlado (SIN PTY)** ← `term.write()` de los
   deltas de token reales (typing auténtico, no simulado). Teclas en el terminal → interceptar
   `onData` → empaquetar la línea como turno stream-json → **al MISMO stdin**.

**Doble input real:** composer (chat) ↔ terminal, ambos escriben el mismo `stdin` del proceso.
AgentsRoom lo FINGE; DevStudio lo hace real. Sin PTY, sin segundo proceso, sin SDK, BYO-licencia
intacta.

### Escalera de autenticidad (elegir #2)
| # | Enfoque | Fidelidad | Fit |
|---|---|---|---|
| 1 | React chat-log de eventos reales (= técnica exacta de Anthropic) | Alta visual | ✓ pero sin "typing en terminal" |
| **2** | **xterm.js modo controlado + deltas reales + teclas→stream-json** | **Alta visual + typing real + cards + aprobaciones** | ★ **mejor autenticidad/esfuerzo, respeta el driver** |
| 3 | PTY real → xterm.js (TUI `claude` literal) | ANSI/TUI literal | △ pierde eventos estructurados; contradice driver; ConPTY Windows |
| 4 | PTY real + proceso stream-json paralelo | Máxima | △ 2 sesiones, 2× costo — diferir |

## 3 · Los 4 huecos del adapter actual (`conductor.go`) — el puente accionable

Hoy `translate()` emite solo `init/delta/result/error` y **descarta `assistant`/tool frames**. Para
llegar a la UX objetivo:
1. **Parsear bloques `tool_use` / `tool_result`** (en frames `assistant`/`user`) → `EvToolCall`/
   `EvToolResult`. Hoy ignorados (`default: return false`). → sin esto no hay tool-cards.
2. **Manejar el envelope de control `control_request:can_use_tool`** → `EvPermissionReq` +
   implementar `ReplyPermission` que escribe `control_response` en stdin. → sin esto no hay modal de
   rama (`❯ 1. Yes / 2. Skip`). ⚠ Sub-documentado + version-dependiente → **primero un probe de 1
   archivo** que dispare un tool gated y loguee stdin/stdout crudo. Bugs conocidos: #34046 (no emite
   `can_use_tool`), #12235 (session_id cambia al resume). Fallback: `--permission-mode acceptEdits`
   + allow-rules.
3. **Exponer handle raw-PTY** (`RawTerminal`) para la pane terminal auténtica.
4. **Renombrar `ClaudeSessionID` → `ProviderSessionID`** + método `Capabilities()` para que el mismo
   UI sirva a todos los providers.

Respeta el boundary `conductor-no-parsea-jsonl` (cada adapter dueño de su protocolo; el usecase solo
ve `Envelope`).

## 4 · Multi-CLI: 4 familias de transporte + Aider outlier

Seis CLIs colapsan en 4 patrones. El seam `ports.AgentPort` ya existe — solo crece de "Claude-shaped"
a "provider-agnostic".

| CLI | Transporte | Eventos estructurados | Tool-calls | Permisos en vivo | Fit UX |
|---|---|---|---|---|---|
| **Claude Code** | stdio NDJSON + control | ✅ | ✅ | ✅ `can_use_tool` | Full (ya wired) |
| **Amp** | stdio NDJSON (**Claude-compatible**) | ✅ mismo shape | ✅ | ⚠ visible, sin handshake | Full (near drop-in) |
| **Cursor** | stdio NDJSON | ✅ | ✅ | ⚠ pre-autorizar | Full (turnos re-spawn) |
| **Gemini** | stdio NDJSON **o ACP** | ✅ | ✅ | ❌ plano / ✅ vía ACP | Full (approvals=ACP) |
| **Codex** | `exec --json` **o `app-server` JSON-RPC** | ✅ | ✅ | ❌ exec / ✅ app-server | Full (bidi=app-server) |
| **OpenCode** | HTTP REST + SSE (**Go SDK oficial**) | ✅ | ✅ ToolPart FSM | ✅ `permission.asked` | Full (otro transporte) |
| **Aider** | raw text (PTY) | ❌ | ❌ (diff worktree) | ❌ `--yes-always` | **Terminal-only** |

**Familias:** A=stdio NDJSON (reusa el pump actual; Amp/Cursor/Gemini/Codex-exec) · B=stdio JSON-RPC
(Codex app-server, Gemini ACP — approvals en vivo + steering) · C=HTTP+SSE (OpenCode, Go SDK) ·
D=raw terminal (Aider). **Recomendación: alinear el Envelope interno con ACP**
(agentclientprotocol.com — estándar cross-vendor, Zed ya lo usa, Gemini lo habla nativo) en vez de
inventar uno.

### Mínimo común denominador (hace universal el doble input)
Presente en **todo** provider incl. Aider: **(1) PTY passthrough** (byte terminal + composer=escribir
stdin) + **(2) señal de liveness busy/waiting/idle** (del estado del proceso, sin eventos). Eso basta
para el frame terminal + roster de mensajería. Todo lo rico (bubbles, tool-cards, modal de rama) =
**enhancement gated en `Capabilities.StructuredEvents`**. Aider degrada a terminal-only en un shell
idéntico — sin special-casing, solo un flag.

## 5 · Contrato Go propuesto (extiende `internal/ports/agent.go`)

```go
type Capabilities struct { StructuredEvents, ToolCallEvents, LivePermissions, MultiTurnStdin, RawPTY, Resume bool }
type AgentSession interface {
    Send(ctx, text) error                              // turno estructurado (composer chat)
    ReplyPermission(ctx, requestID, allow) error       // responde el prompt de rama
    Interrupt(ctx) error                               // cancela turno en vuelo
    Events() <-chan Envelope                            // init/text.delta/tool.call/tool.result/permission.request/turn.result/error/raw
    RawTerminal() (io.ReadWriteCloser, bool)           // PTY auténtico (ok=false ⇒ ninguno)
    Close() error
}
type AgentPort interface { Capabilities() Capabilities; Spawn(ctx, SpawnOpts) (AgentSession, error) }
```
`SpawnOpts` (Cwd/ReadOnly/Resume/PluginDirs/SystemPrompt) generaliza; cada adapter mapea `ReadOnly` a
su sandbox (Claude `--permission-mode plan`, Codex `--sandbox read-only`, Gemini `--approval-mode
plan`, OpenCode config, Aider `--dry-run`).

## 6 · Relación ArnesIA + PB-24
- **Dos ejes ortogonales instalables:** *qué colaborador* (arnés/rol, DATA — dirs+prompts vía flags,
  cero API) × *qué runtime CLI* (provider adapter, CÓDIGO — parser Go). El lock `.devstudio/
  arneses.yaml` gana un gemelo `.devstudio/drivers.yaml` (pinea provider+versión+`Capabilities`
  manifest que el UI lee para decidir panes).
- Un arnés NO se puede shipear como plugin git para un parser Codex/OpenCode. Modelos: (1)
  **adapters compilados + registry declarativo** (recomendado near-term; aprovecha el Go SDK de
  OpenCode) o (2) **sidecars out-of-process = ACP** (máxima expresión "shell-tonto", providers nuevos
  = bridges ACP sin recompilar).
- **PB-24:** el adapter es Go pase lo que pase (Tauri/Wails/Go-webview). En el modelo Tauri, el
  **daemon Go ES el sidecar** = esta capa de adapters. Argumenta contra reescribir en Rust (tirarías
  el Go SDK + conductor). "Shell-tonto" son DOS shells tontos: la *ventana* (Tauri/Wails renderiza
  webview) y el *driver* (renderiza envelopes normalizados). Componen: ventana tonta → webview →
  UI-tonto ← envelope ← adapters. Constraint: la autenticidad terminal exige PTY real → la ventana
  ganadora debe reenviar bytes PTY a un pane xterm.js (stdio nativo, OpenCode vía `/pty`, Aider PTY
  propio).

## 7 · Frontend (dentro de PRENTER teal único)
- `@xterm/xterm` + addons `fit` (tamaño), `webgl` (GPU, clave para streaming), `web-links`. Theme
  `background:#08090a`, superficies `#0c1110`/`#141a19`, `fontFamily:"JetBrains Mono"` (ya
  vendorizada = match Anthropic exacto), `scrollback:5000`.
- **Typing real:** `--include-partial-messages` da `text_delta` → `term.write(delta.text)` ES el
  stream auténtico (nada de typewriter fake en vivo). Typewriter solo para replay de transcript
  cerrado.
- Command palette sembrado de `system/init` (`slash_commands`/`tools`) = refleja el arnés instalado
  real.
- Window-chrome macOS = shell React CSS puro.
- **Mantener teal único PRENTER.** DESCARTAR el multi-accent por-cluster de AgentsRoom (viola RN de
  acento único): estado/rol por forma+label, no por arcoíris. Mapear la selección periwinkle de
  Anthropic → teal; el `●`/loader coral → teal.

## 8 · Próximos pasos sugeridos (para que decida Chris)
1. **Probe de 1 archivo** del control-protocol `can_use_tool` contra el `claude` instalado (valida
   envelope + versión) — ANTES de tocar UI.
2. Fichar en BACKLOG (vía `/pm-*`): ¿refinar PB-09 hacia "conversación terminal-auténtica Híbrido
   3a" + nuevo PB "adapter multi-provider (Capabilities + ACP-aligned Envelope)"? — decisión de
   producto, no la tomo yo.
3. Los 4 huecos de `conductor.go` (§3) son el MVP incremental: tool-cards → permisos → raw-PTY →
   Capabilities. Cada uno entrega valor solo.
4. Mock visual (Artifact) para validar el look con teal PRENTER antes de construir.
