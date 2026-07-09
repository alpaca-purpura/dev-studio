# Épica «Conversación terminal-auténtica» — SSoT (norte + insumos + mockup)

> **Punto de entrada de la ÉPICA** (elevada por Chris 2026-07-09 · patrón Torre de Control, carpeta
> temporal · **F0 CERRADA 2026-07-09 · DH-19**). Objetivo: que la conversación de dev-studio se sienta como un **terminal REAL** (estética
> making-of de Claude Code de Anthropic) con **doble input** (chatear por el chat O escribir directo en
> el terminal, ambos al mismo proceso `claude`) — el efecto de AgentsRoom, pero literal. Backlog:
> **PB-28** (banda 🔴). **Gobierna: [`NORTE-FIRMADO.md`](./NORTE-FIRMADO.md)** (F0 firmada 2026-07-09). Los
> docs 00–05 son insumos de research; el mockup es la SSoT-de-forma. Al cerrar la épica → promover a
> `specs/` + borrar esta carpeta. **F0 CERRADA 2026-07-09 (DH-19):** norte firmado fork-por-fork →
> [`NORTE-FIRMADO.md`](./NORTE-FIRMADO.md). Decisión clave: **B = multi-proveedor desde F1** (Chris
> divergió del recomendado «Claude primero») → el Envelope nace ACP-aligned y Capabilities temprano.

## TL;DR (el veredicto en 5 líneas)
1. La making-of de Anthropic NO es un TTY real — es un transcript React estilizado como terminal.
   dev-studio puede hacerlo **literal** porque ya spawnea `claude` stream-json real (AgentsRoom lo finge).
2. **Solución = Híbrido 3a:** un proceso `claude` stream-json → dos superficies (cards estructuradas +
   xterm.js en modo controlado, deltas reales = typing auténtico; teclas → turno stream-json → mismo stdin).
3. **4 huecos en `internal/adapters/agent/claudecode/conductor.go`:** (1) parsear tool_use/tool_result,
   (2) `control_request:can_use_tool` → modal de rama + `ReplyPermission`, (3) exponer raw-PTY,
   (4) `ProviderSessionID` + `Capabilities()`. **Probe de 1 archivo ANTES de UI** (protocolo sub-doc + version-dep).
4. **Multi-CLI:** 4 familias (stdio NDJSON · stdio JSON-RPC · HTTP+SSE OpenCode con Go SDK · raw Aider).
   Alinear el Envelope con **ACP**. Amp = near drop-in Claude-compat. MCD universal = PTY + liveness.
5. **PRENTER teal único + JetBrains Mono** (ya vendorizada = match Anthropic). Descartar multi-accent de AgentsRoom.

## Orden de lectura para arrancar
0. **`NORTE-FIRMADO.md`** — el norte de la épica (F0 firmada · DH-19): qué construimos, principios, fases F1–F4, decisiones A–E cerradas.
1. **`01-sintesis.md`** — síntesis ejecutiva.
2. **`00-visual-reference.md`** — qué mirar + tokens de ambos referentes + índice de screenshots.
3. **`mock-conversacion.html`** — cómo se vería, en PRENTER. Artifact publicado:
   https://claude.ai/code/artifact/3fa4ae61-17b6-4e43-a6fc-2dde21b68da6 (o abrir el .html local).
3b. **`06-especificacion-mockup.md`** — cada propuesta del mockup especificada (P1–P6): qué es, de qué
   evento stream-json sale, comportamiento, microcopy, tokens, rebanada. Leer junto con el mockup.
4. **`04-mapa-codigo-actual.md`** — dónde tocar (file:line) + GAP summary.
5. **`02-claude-code-protocol.md`** — protocolo stream-json + terminal-en-webview (detalle completo).
6. **`03-multi-cli-interop.md`** — normalización multi-CLI + contrato Go + ACP (detalle completo).
7. **`05-plan-refinamiento-implementacion.md`** — plan accionable: R0 probe → R1 tool-cards → R2 permisos →
   R3 terminal+doble input → R4 Capabilities → R5 multi-CLI. Decisiones pendientes para Chris.

## Cómo continuar (F0 ✅ cerrada · próximo = F1)
> "Épica «Conversación terminal-auténtica». Leé `epicas/conversacion-terminal-autentica/NORTE-FIRMADO.md`
> + este README + `05-plan…md`. F0 ✅ firmada (DH-19): A=probe SÍ · **B=multi-proveedor desde F1** ·
> C=faux 3a · D=DH por rebanada · E=PB-24 no bloquea. Arrancá **F1**: R0 probe + R1 tool-cards
> (Envelope ACP-aligned + Capabilities desde el inicio, por B) → F2 → F3 → F4. Cada rebanada
> R0–R5 entrega valor sola (ver `05-plan…`)."

## Screenshots (en esta carpeta)
- `anthropic-01-hero.png` · `anthropic-02-terminal.png` (el efecto) · `anthropic-03-help.png` (imagen embebida en sub-ventana)
- `agentsroom-04-terminals.png` (terminales) · `agentsroom-05-flow.png` (cajas TÚ↔AGENTE) · `agentsroom-06-status.png` (roster mensajería)
- `agentsroom-01-demo.png` · `agentsroom-02-app.png` · `agentsroom-03-app-full.png` (landing + demo que no bootea)
- `mock-render-check.png` (render del mock)

## Contexto de producto
Refina **PB-09** (terminal embebido — absorbida) y se relaciona con **PB-24** (shell nativo Tauri/Wails)
y **PB-25** (registry ArnesIA). El eje multi-proveedor encaja con **PB-20** (segundo adapter). Ver
`01-sintesis.md` §6 y `03-multi-cli-interop.md` §5–6.
