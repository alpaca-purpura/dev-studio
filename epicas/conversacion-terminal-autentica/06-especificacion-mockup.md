# 06 · Especificación del mockup — cada propuesta, elemento por elemento

> **Lectura obligatoria junto con `mock-conversacion.html`.** El mockup es la SSoT de la FORMA; este doc
> es la regla/decisión DETRÁS de cada elemento (patrón de `spec-mapa-funcional.md`: "el mockup para la
> forma, la sección para la regla"). Cada propuesta se especifica: **qué es · de qué evento stream-json
> sale · comportamiento/estados · microcopy · tokens PRENTER · rebanada** (R# de `05-plan…`). El mockup
> es exploratorio (CSS inline); al refinar, la forma se recompone de átomos reales PRENTER/Storybook
> ("lo que veo = lo que se programa"). Nada acá es todavía contrato congelado — es la propuesta a firmar.

Tokens PRENTER usados (de `04-mapa-codigo-actual.md` §3): ground `#08090a` · window `#0c1110` · raised
`#141a19` · border `#1f2826` · **teal** `#1fc6b8` (único acento de marca) · texto `#e6ecea`/`#8a9995`/
`#5a6a66` · **amber `#e0b341`** (SOLO estado semántico «te espera / atención», no es acento de marca) ·
red `#e5484d` (error) · **JetBrains Mono** en toda la superficie.

Mapa evento→UI (referencia rápida, de `02` §A.2 / `03` §8): `system/init` → header+roster+palette ·
`stream_event/text_delta` → typing en transcript+terminal · `assistant.tool_use` → tool-card (R1) ·
`user.tool_result` → output de la tool-card (R1) · `control_request:can_use_tool` → modal de rama (R2) ·
`result` → footer costo + estado idle (roster).

---

## P1 · Window chrome (marco tipo terminal macOS)
- **Qué:** ventana redondeada near-black con 3 traffic-lights (rojo/amarillo/verde), título
  `dev-studio · workspace {slug}`, y a la derecha un **toggle de modo** `>_ Terminal` / `¶ Chat`.
- **Por qué:** es el gesto que "cuelga" la conversación de la estética making-of de Anthropic; encuadra
  todo como un terminal auténtico. El toggle es el mismo dispositivo `Read in Terminal / Read as Article`
  de Anthropic, adaptado a "vista terminal cruda" vs "vista chat estructurada" del MISMO stream.
- **Datos:** estático + `{slug}`/repo del workspace activo (ya existe en el store de sesiones).
- **Comportamiento:** el toggle cambia el ÉNFASIS de render (terminal-primero vs cards-primero), no el
  proceso — las dos vistas leen el mismo stream. Traffic-lights decorativos (no cierran la app; opcional:
  el rojo = cerrar workspace con el modal conservar/borrar que ya existe, CAP-10).
- **Microcopy:** `>_ Terminal` · `¶ Chat`.
- **Estilo:** chrome `linear-gradient` sutil + `border-bottom #1f2826`; toggle activo = `teal-faint` bg + texto teal.
- **Rebanada:** R3 (shell del terminal) — el chrome como shell React CSS puro.

## P2 · Roster de workspaces (lista estilo mensajería)
- **Qué:** rail izquierdo = lista de sesiones/agentes como una app de chat. Header `WORKSPACES {n}`.
  Grupos por estado: **Activo · Te espera · Inactivo**. Cada fila: avatar circular con glow + **dot de
  estado**, nombre (arnés), sub-línea (última actividad/estado), **badge de rol** (= el arnés instalado),
  timestamp derecho, **badge de no-leídos**.
- **Por qué:** es el "chatear por el chat" de AgentsRoom — los agentes como contactos con no-leídos y
  estado en vivo. Encaja con el rail Repositorios→Workspaces existente (workspace = sesión 1:1).
- **Datos:** lista de sesiones del store · estado derivado de **liveness** (MCD, `03` §7): *activo*
  mientras emite output, *te espera* al bloquear en un `control_request` (permiso) o input, *inactivo*
  al `result`/idle · badge de rol = `arnes.l0.nombre` (de `--plugin-dir`, PB-25) · no-leídos = mensajes
  de agente sin ver (patrón notificaciones).
- **Estados de fila:** normal · **seleccionada** (glow teal, `teal-faint` + borde `teal-line`) ·
  **atención/te-espera** (stripe + borde + texto **amber** — un agente pide permiso) · inactiva (avatar
  atenuado, dot gris).
- **Microcopy (ejemplos del mock):** `Activo` / `Te espera` / `Inactivo` · sub `editando conductor.go…`
  · `pide permiso: go test…` · `exploración read-only` · badges `Backend` `Frontend` `Plan` · ts
  `ahora`/`12s`/`4m`.
- **Estilo:** avatar gradient teal (o violeta para distinguir agentes; **el violeta es identidad de
  avatar, NO acento de UI** — se mantiene teal único para acentos). Dot: teal=activo, amber=espera,
  gris=inactivo.
- **Rebanada:** base = R3/R4 (liveness + Capabilities). No-leídos + badges de rol se enriquecen con
  PB-25 (registry) y el eje notificaciones.

## P3 · Header de contexto (chips)
- **Qué:** fila de chips bajo el chrome: **rol** (teal, `arnesia-dev/backend`), **model** (`opus-4.8`),
  **branch** (`feature/terminal-embebido`), **cwd** (`~/…/workspaces/dev-studio/{slug}`).
- **Por qué:** contexto siempre visible de con QUÉ estás hablando y DÓNDE. Reemplaza el header actual de
  `session-view.tsx` (rol/cwd/model/branch ya existen ahí).
- **Datos:** de `system/init` (`model`, `cwd`, `permissionMode`, `tools`) + el arnés del spawn + branch
  del worktree. El chip rol es teal (acento); el resto neutro.
- **Rebanada:** ya existe parcialmente (R1 lo alimenta de `init` en vez de props).

## P4 · Transcript (stream chat-log)
Contenedor scrolleable, `scrollback` cap ~5000 (de `02` §A.3). Sub-elementos:

### P4.a · Turno del usuario (TÚ)
- **Qué:** bullet `❯` + `Tú` + cuerpo del prompt. Alineado como en el mock (izquierda, atenuado).
- **Datos:** eco optimista al enviar (ya lo hace `sendTurn` en el store) + confirmación por el frame
  `user` si se usa `--replay-user-messages`.
- **Microcopy:** el texto del usuario, con `/comando` resaltado en teal.

### P4.b · Turno del agente
- **Qué:** `● {arnés} — arnés {Rol} · {model}` (bullet teal) + cuerpo que **streamea** token a token.
- **Datos:** `stream_event/text_delta` → `term.write`/append (typing REAL, no simulado, `02` §B).
  Header del turno de `assistant.message.model` + arnés del spawn.
- **Estilo:** bullet `●` teal (mapea el coral de Anthropic → teal), handle bold, rol atenuado.

### P4.c · Tool-card (Read / lectura)
- **Qué:** tarjeta `⎿ {Tool} {arg} … {status}` + cuerpo con output (líneas numeradas).
- **Por qué:** el hueco #1 (`conductor.go` hoy `default: return false`). Sin esto no hay visualización
  de herramienta = el diferenciador del look IDE-grade.
- **Datos:** `assistant.tool_use{name,input,id}` (header) → `user.tool_result{tool_use_id,content}`
  (output). Ver `02` §A.2, `03` §8.
- **Estados:** `running` (spinner) → `ok` / `error` (rojo). Colapsable (output largo).
- **Microcopy:** `Read internal/adapters/agent/claudecode/conductor.go · OK`.
- **Rebanada:** **R1**.

### P4.d · Tool-card (Edit / escritura, con diff)
- **Qué:** igual que P4.c pero cuerpo = **diff** (líneas `+` teal / `−` rojo) + resumen `+14 −1`.
- **Datos:** `tool_use(Edit/Write)` input → diff renderizado. Para proveedores sin diff estructurado se
  deriva del worktree (`03`: Aider).
- **Rebanada:** **R1** (render) · enlaza con PB-10/PB-11 (tab Cambios + diff conversacional).

### P4.e · Prompt de rama (modal de permiso)
- **Qué:** bloque **amber** con la pregunta + comando en chip + **opciones numeradas** seleccionables
  (`❯ 1. Permitir una vez` / `2. Permitir siempre en este workspace` / `3. Rechazar y decir por qué`) +
  hint (`1–3 o Enter · ↑/↓ · Esc · vía control_request:can_use_tool`).
- **Por qué:** el hueco #2. Convierte el `can_use_tool` crudo en el prompt de rama hermoso de la
  making-of (`❯ 1. Yes / 2. Skip`). Es el corazón del "se siente real": el agente PIDE y vos decidís en vivo.
- **Datos:** `control_request:can_use_tool{tool_name,input,tool_use_id,decision_reason}` → render;
  la elección → `ReplyPermission(requestID, allow, [message])` que escribe `control_response` (`02` §A.4).
  Mientras está abierto, el turno **bloquea** (roster pasa a «Te espera», P2) — responder rápido (~60s timeout).
- **Estados:** opción enfocada (`❯`, caja amber) · navegación teclado · «Rechazar» abre input de razón
  (se pasa como `message` del deny). Fallback si el canal falla: `--permission-mode acceptEdits` (`02` §A.4).
- **Microcopy:** las 3 opciones + hint. La opción 2 es scoped al workspace (allow-rule persistida).
- **Estilo:** **amber = semántico (atención), no acento de marca.** Único lugar donde aparece amber.
- **Rebanada:** **R2** (depende de **R0** probe — protocolo version-dependiente).

### P4.f · Loader de streaming
- **Qué:** `Pensando •••` (dots animados, teal) mientras el agente procesa entre turnos/tools.
- **Datos:** entre `message_start` y `result`, o mientras una tool corre. Respeta `prefers-reduced-motion`.
- **Microcopy:** `Pensando` (mapea `+ Remembering…` de Anthropic).
- **Rebanada:** R1/R3.

## P5 · Doble input (dock con tabs Terminal | Chat)
- **Qué:** pestañas `>_ Terminal` / `¶ Chat`; el input activo abajo. **Terminal** = línea `❯` con caret
  parpadeante, acepta comandos (`/…`) o texto libre. **Chat** = composer redimensionable (el actual, 7
  acciones + cola).
- **Por qué:** EL diferenciador — el doble input que AgentsRoom finge y DevStudio hace **literal**.
  Ambos escriben el **mismo stdin** del proceso `claude` (`ccSession.Send`, único pipeline).
- **Datos/comportamiento:**
  - Terminal: `xterm.js` modo controlado (sin PTY). `term.write` de los deltas; `term.onData` intercepta
    teclas, en Enter empaqueta la línea como turno stream-json → `POST /turn` (mismo pipeline que el chat).
  - Slash-command palette sembrado de `system/init` (`slash_commands`/`tools`) = refleja el arnés real.
  - Chat: composer existente (`session-view.tsx` L148-232) — Enter=enviar, Enter-mientras-streamea=encolar (RN-8).
- **Microcopy:** placeholder `/test acepta comandos o escribe libre…`.
- **Estilo:** tab activa teal; caret teal (`cb` blink, reduced-motion off).
- **Rebanada:** **R3** (terminal + segundo input path). Es donde nace el "doble input = un solo stdin".

## P6 · Footer de estado
- **Qué:** línea inferior: `⏻ conectado · stream-json · {n} turno · {tok} tok · ${cost} · {dur}s`.
- **Por qué:** transparencia del motor real (lo que Anthropic muestra como cost/usage) + señal de conexión.
- **Datos:** de `result` (`total_cost_usd`, `usage{input/output_tokens}`, `num_turns`, `duration_ms`) +
  estado del proceso (conectado/reconectando — liga con SSE replay PB-18).
- **Microcopy:** `conectado` (teal) · `stream-json` · `1 turno` · `2 471 tok` · `$0.04` · `3.2s`.
- **Rebanada:** R1 (parsear `result` completo; hoy solo se usa parcialmente).

---

## Matriz elemento → evento → rebanada (resumen)
| Elemento | Evento stream-json fuente | Rebanada |
|---|---|---|
| P1 Window chrome + toggle | (estático + workspace) | R3 |
| P2 Roster mensajería | liveness + `init` + arnés (PB-25) | R3/R4 |
| P3 Chips de contexto | `system/init` | R1 |
| P4.a Turno usuario | eco optimista / `user` | (existe) |
| P4.b Turno agente (typing) | `stream_event/text_delta` | R1/R3 |
| P4.c/d Tool-cards (+diff) | `assistant.tool_use` → `user.tool_result` | **R1** |
| P4.e Prompt de rama | `control_request:can_use_tool` ⇄ `ReplyPermission` | **R2** (dep R0) |
| P4.f Loader | entre `message_start`…`result` | R1/R3 |
| P5 Doble input | teclas → turno stream-json (mismo stdin) | **R3** |
| P6 Footer costo | `result` (cost/usage/duration) | R1 |

## Qué NO adopta el mockup (decisiones negativas)
- ❌ Multi-accent per-cluster de AgentsRoom (emerald/orange/violet/…) → **teal único**; rol/status por
  forma+label. Amber es la ÚNICA excepción y es semántica (atención), no acento.
- ❌ PTY real para Claude → **modo controlado 3a** (PTY = escalón futuro / otros proveedores).
- ❌ Typewriter fake en vivo → **deltas reales** (typewriter solo para replay de transcript cerrado).
