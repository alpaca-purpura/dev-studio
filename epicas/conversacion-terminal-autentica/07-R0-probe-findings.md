# 07 · R0 — Hallazgos del probe del protocolo de control (verificado en vivo)

> **Rebanada R0 (spike descartable, ejecutada 2026-07-09).** Salida que `05-plan…` pide: «un `.md`
> con el shape verificado + la versión de `claude`». **Precede a R2** (permisos + modal de rama).
> Método: se spawneó el `claude` REAL en `-p --input-format stream-json --output-format stream-json
> --verbose`, se mandó un turno que dispara un tool gated (`Bash: echo …`) y se logueó **stdin/stdout
> crudo**. También se desensambló el binario (`strings`) para el shape del control-protocol. Scripts +
> logs crudos en el scratchpad de la sesión (`r0_probe.py`, `r0_probe2.py`, `r0_run_*.log`) — descartables.

- **Versión pineada:** `claude` **2.1.205** (Bun-compiled ELF, 257 MB, `~/.local/share/claude/versions/2.1.205`).
- **Modelo del turno:** `claude-opus-4-8[1m]` · `apiKeySource: none` (BYO licencia, DH-10 intacto ✓).
- **Costo del spike:** 2 turnos reales ≈ `$0.23` c/u (cuenta contra el rate-limit `five_hour`, no $ directo).

---

## 0 · TL;DR — el probe corrigió el plan (por eso R0 existía)

1. **`--permission-prompt-tool stdio` (lo que asumían `02`/`05`) está MAL, doble:** (a) el flag NO
   aparece en el help estándar de 2.1.205 y (b) cuando se pasa, **exige el NOMBRE DE UN MCP TOOL**
   (`"... must be an MCP tool"`, `"... not found. Available MCP tools: …"`) — `"stdio"` no es válido.
   El modal de rama de R2 **no se construye sobre ese flag tal cual**.
2. **En modo `default` headless, el `claude` AUTO-PERMITE los tools** (Bash corrió sin pedir permiso,
   con `permissions:{}` en settings del user, `permission_denials:[]`). **No hay «ask» por defecto.**
3. **El handshake de control funciona:** el cliente manda `control_request{subtype:initialize}` y el
   CLI responde `control_response{subtype:success, response:{commands:[…]}}` — **verificado en vivo**.
   Pero el `initialize` mínimo **NO** alcanza para que el CLI enrute el permiso al cliente.
4. **La ruta «ask» existe y es:** `control_request` con `request.subtype = "can_use_tool"` (CLI→cliente)
   ⇄ `control_response{response:{behavior: allow|deny, …}}` (cliente→CLI). Falta pinear el disparador
   exacto que la habilita sobre stdio crudo (§3.3) — **eso es lo primero de R2**.
5. **Todo lo demás de la conversación está DES-RIESGADO y capturado en vivo:** `system/init`,
   `tool_use`, `tool_result`, `text_delta` (typing), `result` (costo/usage). **R1 se puede construir ya.**

---

## 1 · Event shapes verificados EN VIVO (usables tal cual para R1/P4/P6)

Todos capturados del stdout real. JSON recortado (se elide `signature`/`uuid` largos).

### 1.1 · `system/init` (→ P3 chips · P5 command palette · roster)
Campos vivos relevantes: `cwd`, `session_id`, `model`, `permissionMode`, `apiKeySource:"none"`,
`tools:[…]`, `slash_commands:[…]`, `plugins:[…]`, `mcp_servers:[{name,status}]`,
`capabilities:["interrupt_receipt_v1"]`, `claude_code_version:"2.1.205"`, `output_style`.
```jsonc
{"type":"system","subtype":"init","cwd":"…/r0_workdir","session_id":"e4339525-…",
 "tools":["Task","Bash","Edit","Read","Write",…],"model":"claude-opus-4-8[1m]",
 "permissionMode":"default","slash_commands":["design-md",…,"init","mcp","model",…],
 "apiKeySource":"none","capabilities":["interrupt_receipt_v1"],"claude_code_version":"2.1.205", …}
```
> **Nota:** `slash_commands` (init) es la lista para sembrar la palette (P5). El propio protocolo avisa
> que puede cambiar mid-sesión y se re-empuja aparte (ver §3.2 `supportedCommands`).

### 1.2 · `assistant` con `tool_use` (→ P4.c/d tool-cards · hueco #1)
Llega **streameado**: `content_block_start{content_block:{type:"tool_use",id,name,input:{}}}` →
varios `content_block_delta{delta:{type:"input_json_delta",partial_json:"…"}}` → y el frame `assistant`
consolidado:
```jsonc
{"type":"assistant","message":{"model":"claude-opus-4-8","id":"msg_…","role":"assistant",
  "content":[{"type":"tool_use","id":"toolu_01E9xFiJ…","name":"Bash",
              "input":{"command":"echo PROBE_OK_12345","description":"Echo probe string"},
              "caller":{"type":"direct"}}],
  "stop_reason":null,…},"session_id":"e4339525-…","request_id":"req_…"}
```
Campos para la card: `content[].id` (`toolu_…`, pareo con el result), `.name`, `.input` (obj), `.caller.type`.

### 1.3 · `user` con `tool_result` (→ output/diff de la tool-card · hueco #1)
```jsonc
{"type":"user","message":{"role":"user","content":[
   {"tool_use_id":"toolu_01E9xFiJ…","type":"tool_result","content":"PROBE_OK_12345","is_error":false}]},
 "session_id":"e4339525-…",
 "tool_use_result":{"stdout":"PROBE_OK_12345","stderr":"","interrupted":false,"isImage":false,"noOutputExpected":false}}
```
Pareo por `tool_use_id`. `is_error:true` ⇒ card en rojo. `tool_use_result` trae stdout/stderr crudos
(útil para el cuerpo de la card). Para Edit/Write el `input` del `tool_use` trae el diff a renderizar.

### 1.4 · `stream_event / text_delta` (→ typing auténtico · P4.b)
```jsonc
{"type":"stream_event","event":{"type":"content_block_delta","index":0,
  "delta":{"type":"text_delta","text":"ROBE_OK_12345`"}},"session_id":"e4339525-…"}
```
`term.write(delta.text)` ES el typing real (P4.b). Ojo: hay bloques `thinking` (`signature_delta`) —
filtrar por `content_block.type`/`delta.type` antes de escribir al terminal. `message_start` trae
`ttft_ms`; `message_delta` trae `stop_reason` (`tool_use`/`end_turn`).

### 1.5 · `result` (→ P6 footer costo/estado · cierre de turno)
```jsonc
{"type":"result","subtype":"success","is_error":false,"duration_ms":12568,"duration_api_ms":12017,
 "num_turns":2,"result":"Output: `PROBE_OK_12345`","stop_reason":"end_turn","session_id":"e4339525-…",
 "total_cost_usd":0.232508,
 "usage":{"input_tokens":12754,"output_tokens":120,"cache_creation_input_tokens":14950,"cache_read_input_tokens":32476},
 "modelUsage":{"claude-opus-4-8[1m]":{"inputTokens":12754,"outputTokens":120,"costUSD":0.2325,"contextWindow":1000000}},
 "permission_denials":[],"terminal_reason":"completed"}
```
P6 sale entero de acá: `total_cost_usd`, `usage.{input,output}_tokens`, `num_turns`, `duration_ms`.
`permission_denials` (array) lista los tools auto-denegados del turno (hoy vacío — nada se deniega).

### 1.6 · `system/status` + ruido de hooks
- `system/status{status:"requesting"}` entre turnos (señal de liveness → roster «Activo», P2).
- `system/hook_started` / `system/hook_response` — **el SessionStart hook del user (caveman) se cuela
  en el stream** con su payload entero. El conductor debe **filtrar `type:"system"` subtypes de hook**
  (no son del turno). Gotcha real observado en el probe.
- `rate_limit_event{rate_limit_info:{status:"allowed", rateLimitType:"five_hour", …}}` — insumo para un
  aviso de cuota en el footer (P6) / roster.

---

## 2 · El mecanismo de permisos (el corazón de R2) — qué es verdad

### 2.1 · Lo que se probó (2 corridas en vivo)
| Corrida | Flags | Resultado |
|---|---|---|
| `default`, sin handshake | `-p … stream-json` | Bash **auto-permitido**, sin `can_use_tool`. `permission_denials:[]`. |
| `default`, **con** `initialize` | idem + handshake | `initialize` respondido OK; Bash **igual auto-permitido**, sin `can_use_tool`. |

**Conclusión:** ni el modo `default` ni el `initialize` mínimo bastan para que el CLI **pida** permiso
al cliente. El CLI trata al host stream-json como responsable y auto-permite.

### 2.2 · Reglas de supresión de `canUseTool` (extraídas del binario, verbatim)
- `bypassPermissions` → auto-aprueba todo (excepto deny rules). `canUseTool` NO se invoca.
- Entradas *bare* en `allowedTools` → auto-aprueban el tool entero **antes** del callback.
- Allow-rules de settings files → **shadowean** el callback (silenciosamente).
- Recomendación LITERAL del binario: *«To gate every tool call, use a PreToolUse hook instead.»*

### 2.3 · Modos de permiso reales en 2.1.205
`--permission-mode` choices = **`acceptEdits` · `auto` · `bypassPermissions` · `manual` · `dontAsk` · `plan`**
(desaparecieron los viejos; aparecieron `auto`, `manual`, `dontAsk`). `plan` = read-only (ya usado en
exploración PB-27). `manual`/`auto`/`dontAsk` = sin probar aún (candidatos para el «ask» — ver R2).

---

## 3 · Referencia del control-protocol (para el conductor Go)

### 3.1 · Subtypes que el CLI ACEPTA del cliente (control_request cliente→CLI)
`initialize` · `interrupt` · `set_permission_mode` · `set_mcp_permission_mode_override` · `set_model` ·
`set_max_thinking_tokens` · `set_cwd` · `rewind_files` · `get_settings` · `apply_flag_settings` ·
`remote_control` · `mcp_message`/`mcp_toggle`/`mcp_authenticate` · `hook_callback` · `side_question`.
(Desconocido → `"Unsupported control request subtype: …"`.)

### 3.2 · `initialize` — handshake VERIFICADO en vivo
Request (cliente): `{"type":"control_request","request_id":"req-init-1","request":{"subtype":"initialize"}}`
Response (CLI): `{"type":"control_response","response":{"subtype":"success","request_id":"req-init-1","response":{"commands":[{"name","description","argumentHint"},…]}}}`
> El `response.commands` del init = fuente canónica de la palette (P5), más fresca que `slash_commands`
> del `system/init`. El protocolo re-empuja la lista si cambia mid-sesión (`supportedCommands` se captura
> una vez en `initialize`).

### 3.3 · `can_use_tool` — el «ask» (shape del binario; falta captura viva)
Emitido por el CLI cuando un tool necesita aprobación **y el cliente registró el callback**:
```jsonc
// CLI → cliente
{"type":"control_request","request_id":"…","request":{
   "subtype":"can_use_tool","tool_name":"Bash","input":{…},
   "tool_use_id":"toolu_…","title":"…","display_name":"…","description":"…", …}}
// cliente → CLI  (allow / deny)
{"type":"control_response","response":{"subtype":"success","request_id":"…",
   "response":{"behavior":"allow","updatedInput":{…}?,"updatedPermissions":[…]?}}}
{"type":"control_response","response":{"subtype":"success","request_id":"…",
   "response":{"behavior":"deny","message":"por qué"}}}
```
- `updatedInput` = el handler puede **modificar el input del tool** antes de permitir (valida contra el
  schema; si falla → `InputValidationError`). Palanca para «editar y permitir».
- `updatedPermissions` = persistir una allow-rule (mapea a la opción **«2. Permitir siempre»** de P4.e).
- Evento hermano: cuando un tool se **auto-deniega** sin prompt, sale un evento aparte (no `can_use_tool`)
  para que el host pinte la denegación — el «ask» y el «deny short-circuit» son distintos.

### 3.4 · El disparador que falta pinear (primera tarea de R2)
Cómo hacer que el CLI **enrute** `can_use_tool` al cliente sobre stdio. Candidatos (del binario):
- Flag interno `requireCanUseTool` / `permissionPromptToolName` (existen; `--permission-prompt-tool` es
  válido pero **exige un MCP tool**, no `"stdio"`).
- `canUseTool` y `permissionPromptToolName` son **mutuamente excluyentes** («use one or the other»).

---

## 4 · Los 3 caminos para el modal de R2 (recomendación)

| Camino | Cómo | Pro | Contra |
|---|---|---|---|
| **A · Control-protocol `canUseTool`** | Handshake que registra el callback → CLI manda `can_use_tool` por el mismo stdout; conductor responde `control_response`. | Todo en el stream que el conductor YA posee; cero proceso extra; `updatedInput`/`updatedPermissions` gratis (opciones 1/2 de P4.e). | El disparador exacto sobre stdio crudo falta pinear (§3.4); puede requerir replicar el handshake del Agent SDK. |
| **B · `--permission-prompt-tool mcp__…`** | Exponer un **MCP server in-process** con un tool de aprobación; el CLI lo llama para decidir. | Flag documentado, camino soportado. | El conductor debe hablar MCP in-process (`mcp_message`); el tool «no puede requerir interacción» → bloquear internamente contra la UI por un canal lateral. Más pesado. |
| **C · PreToolUse hook** | Registrar un hook que llama al server local de dev-studio, bloquea en la decisión del modal y devuelve allow/deny. | **Recomendado por el propio binario** («gate every tool call»); estable/documentado. | Suma un binario-helper que el hook invoca; arquitectura fuera del stream. |

**Recomendación:** arrancar R2 con un **spike A→C** de 1 archivo (continuación de este probe): intentar A
(pinear el handshake) contra el `claude` instalado; si en ~1 sesión no enruta, caer a **C (PreToolUse
hook)** que es el camino que el binario bendice. **No** construir el modal de rama antes de tener UNA
captura viva de `can_use_tool` (allow Y deny). Fallback de degradación (ya en el plan): `--permission-mode
acceptEdits` + allow-rules — el turno no se bloquea, se pierde el modal.

---

## 5 · Impacto en los huecos de `conductor.go` (revalidado)

| Hueco (de `04`) | Estado tras R0 |
|---|---|
| #1 tool_use/tool_result → tool-cards (R1) | ✅ **Des-riesgado.** Shapes vivos en §1.2/1.3. Construir ya. |
| #2 `can_use_tool` → modal + `ReplyPermission` (R2) | ⚠️ **Re-scopeado.** El flag del plan era erróneo (§0, §2). R2 empieza con el spike A→C (§4) para capturar el shape vivo antes de la UI. |
| #3 raw-PTY (`RawTerminal`) | Sin cambio — escalón futuro (C firmado = faux 3a). Para 3a el terminal se alimenta de `text_delta` (§1.4), no de PTY. |
| #4 `Capabilities()` + `ProviderSessionID` (adelantado a F1 por B) | `init.capabilities:["interrupt_receipt_v1"]` + `apiKeySource` + `permissionMode` + `tools` → insumos reales para el manifest de `Capabilities`. `session_id` sale de `init` y del `result`. |

**Gotchas para el conductor (observados):**
- `pump` ya usa `bufio` 1 MB — bien: el `system/init` real pesa >30 KB.
- **Filtrar los `system/hook_*`** del stream (el SessionStart hook del user se cuela entero).
- El `session_id` es estable dentro de la corrida, pero **leerlo del evento** (init/result), no asumirlo
  (riesgo #12235 de resume documentado en `05`).
- `--output-format stream-json` **exige `--verbose`** en `-p` (el CLI aborta si falta) — ya está en `buildArgs`.
- `input_json_delta` llega en trozos; el `tool_use.input` sólo está completo en el frame `assistant`
  consolidado — parsear ahí, no de los deltas parciales.

---

## 6 · Qué queda descartable vs qué se promueve
- **Descartable:** `r0_probe.py`, `r0_probe2.py`, `r0_run_*.log`, `claude_strings.txt` (scratchpad).
- **Se promueve (este doc):** el shape verificado + la corrección del mecanismo de permisos. Al cierre de
  la épica migra a `specs/driver-multiproveedor/` como referencia del adapter Claude.
