# Referencia visual — "Making of Claude Code" × AgentsRoom → conversación DevStudio

**Capturado:** 2026-07-09 (Chrome DevTools, viewport 1280×800). Screenshots en esta carpeta.
**Objetivo:** que la conversación de DevStudio se sienta REAL (terminal auténtica) + colgada de la
estética Anthropic, con el doble input de AgentsRoom (chatear por el chat O escribir directo en el
"terminal"). Insumo para los subagentes de investigación.

## A · Anthropic — "The Making of Claude Code"
URL: https://www.anthropic.com/features/making-of-claude-code · screenshots `anthropic-0*.png`

**Dispositivo central = doble modo de lectura:** toggle `>_ Read in Terminal` ↔ `¶ Read as Article`.
El modo terminal ES el efecto que queremos.

Estructura del modo terminal (`anthropic-02-terminal.png`):
- **Window chrome macOS**: 3 traffic-lights reales (#FF5F56 / #FFBD2E / #27CA40), ventana redondeada
  con padding, sobre bg near-black.
- **Header**: avatar pixel-art (Clawd) + `Claude Code / A Time Capsule / 2025–2026`.
- **Chapter tabs** horizontales (scrollables): `I. Origins`(activo, caja resaltada) … `VII. The future`.
- **Cuerpo = chat-log**: cada entrada `● @handle — rol` (bullet + handle coral + rol), luego texto
  indentado. Streaming con línea loader `+ Remembering...` (coral) mientras "escribe".
- **Prompts interactivos de rama** (`take_snapshot`): `Would you like to see this photo?` →
  `❯ 1. Yes, show me` / ` 2. Skip` + hint `1/2 or Enter to select · ↑/↓ to navigate · Esc to go back`.
- **Imágenes embebidas** dentro de sub-ventanas estiladas (`anthropic-03-help.png`: ventana anidada
  "VSCode plugin" con callouts anotados).
- **Input real de comando**: `❯ Type a command or /help` (textbox). `/help` → "Show available commands".
- **Footer**: `/ to type a command` (izq) · `Chapter 1 of 7` (der).

Tokens (computados):
- Window bg `#141413`; superficies anidadas `#1A1918`, `#252321`.
- Mono = **JetBrains Mono** en TODO (DevStudio ya la vendoriza → match exacto).
- Accent coral `#D77757` (handle/loader) + lavanda/periwinkle `#B1B9F9` (botón activo + selección).
- Article mode = anthropicSerif sobre crema `#FAF9F5`.

## B · AgentsRoom
URL: https://agentsroom.dev/es · screenshots `agentsroom-0*.png`

**Demo en vivo** = iframe `/demo/index.html` (app React idéntica al escritorio). NO bootea en Chrome de
automatización (root React vacío + 404, anti-hotlink / gating). La referencia real la dan los mockups
inline del landing (fieles al producto).

Dos patrones clave:
1. **Roster de agentes = app de mensajería** (`agentsroom-06-status.png`): panel derecho tipo lista de
   chats. Grupos `ACTIVE / INACTIVE` (mono caps). Cada fila: avatar circular con glow + nombre bold +
   último mensaje/estado + timestamp der (`just now`, `1m`, `2m ago`, `9m`). Badge de rol pill en mono
   caps de color: `BACKEND DEVELOPER`(verde) `DEVOPS ENGINEER`(verde) `FULL-STACK DEVELOPER`(violeta).
   Fila seleccionada = glow naranja. Dots de no-leído en el avatar. → ESTE es el "chatear por el chat".
2. **Cajas de turno TÚ ↔ AGENTE** (`agentsroom-05-flow.png`): cada feature muestra un intercambio con
   cajas bordeadas etiquetadas `TÚ` (prompt) → `AGENTE` (respuesta). Chips de slash-command
   `/review-pr /tests /commit`. Waveform de voz. Selección sketch/screenshot con rects dashed.
3. **Terminales reales** (`agentsroom-04-terminals.png`): mini-mockup con traffic dots + tabs
   `dev/build/logs` (activo pill verde) + `$ npm run dev / ready on :3000`. "Terminales completos,
   separables, pestañas/divisiones". → ESTE es el "escribir directo en el terminal".

Tokens (Tailwind v4):
- Bg warm near-black `#14110C`; superficies `--color-bg-surface #1b1813`, `-3 #2d2820`.
- Texto crema `#F0EBDF`, secundario `#a8a091`, muted `#6E6757`.
- **Multi-accent color-coded por cluster**: cada familia usa un tono Tailwind (emerald #10b981,
  orange #f97316, violet #8b5cf6, pink, indigo, cyan, amber, lime, sky, red) con bg @0.125 + border
  @0.19–0.31 (chips tintados). Naranja = CTA primario.
- Fonts: sans **Plus Jakarta Sans** (var), display **Space Grotesk** + serif "NewYork", mono =
  system stack (Menlo/Cascadia/Fira — NO vendorizada).

## C · Síntesis para DevStudio (PRENTER: dark-first, teal único, JetBrains Mono/Jost/Mulish vendor)
- La **conversación** = frame terminal estilo Anthropic (window chrome + stream + input slash-cmd +
  prompts de rama) porque el driver ya es CLI-nativo (spawn `claude` stdin/stdout stream-json, 1 proc
  por sesión). El terminal puede ser REAL (passthrough stdin/PTY) → autenticidad literal, no simulada.
- El **doble input** = (a) composer estructurado tipo chat "TÚ/AGENTE" ↔ (b) terminal crudo. Ambos
  escriben al MISMO stdin del proceso. AgentsRoom lo finge; DevStudio lo hace de verdad.
- El **roster multi-agente** (workspaces = sesión 1:1) puede pintarse como lista de mensajería
  (estado activo/esperando/inactivo, badge de rol = arnés instalado, no-leídos) — encaja con el rail
  Repositorios→Workspaces existente.
- Paleta: mantener **teal único PRENTER** (más cerca del emerald de AgentsRoom) + JetBrains Mono
  (match Anthropic). NO adoptar el multi-accent de AgentsRoom (rompe RN de acento único).
- Relación backlog: toca **PB-24** (webview nativa / shell-tonto Tauri) y la conversación. Idea nueva
  → PB antes de trabajar (regla de producto). Esto es investigación/diseño, no build.
