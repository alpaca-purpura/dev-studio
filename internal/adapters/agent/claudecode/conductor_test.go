package claudecode

import (
	"strings"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// Frames reales capturados del probe R0 (07-R0-probe-findings.md, claude v2.1.205).
// Sirven de fixture: si la CLI cambia el shape, estos tests fallan y avisan.

const frameToolUseBash = `{"type":"assistant","message":{"model":"claude-opus-4-8","id":"msg_011","role":"assistant","content":[{"type":"tool_use","id":"toolu_01E9xFiJFVQorhZsB1z8NrsQ","name":"Bash","input":{"command":"echo PROBE_OK_12345","description":"Echo probe string"},"caller":{"type":"direct"}}],"stop_reason":null},"session_id":"e4339525","request_id":"req_011"}`

const frameToolResultOK = `{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_01E9xFiJFVQorhZsB1z8NrsQ","type":"tool_result","content":"PROBE_OK_12345","is_error":false}]},"session_id":"e4339525","tool_use_result":{"stdout":"PROBE_OK_12345","stderr":"","interrupted":false}}`

// tool_result con content en forma de array de bloques + error (p. ej. una tool que falla).
const frameToolResultArrErr = `{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_XY","type":"tool_result","content":[{"type":"text","text":"boom: file not found"}],"is_error":true}]}}`

func firstOfKind(t *testing.T, line string, kind ports.AgentEventKind) ports.AgentEvent {
	t.Helper()
	evs, ok := translate([]byte(line))
	if !ok {
		t.Fatalf("translate no emitió nada para %s", kind)
	}
	for _, ev := range evs {
		if ev.Kind == kind {
			return ev
		}
	}
	t.Fatalf("translate no emitió un evento %s; emitió %v", kind, evs)
	return ports.AgentEvent{}
}

// TestTranslateToolUse (hueco #1, R1): el frame `assistant` con un bloque tool_use debe
// producir un EventToolCall con id/nombre/input reales — hoy `default: return false` los tira.
func TestTranslateToolUse(t *testing.T) {
	ev := firstOfKind(t, frameToolUseBash, ports.EventToolCall)
	if ev.ToolID != "toolu_01E9xFiJFVQorhZsB1z8NrsQ" {
		t.Errorf("ToolID = %q", ev.ToolID)
	}
	if ev.ToolName != "Bash" {
		t.Errorf("ToolName = %q, quería Bash", ev.ToolName)
	}
	if !strings.Contains(ev.ToolInput, `"command":"echo PROBE_OK_12345"`) {
		t.Errorf("ToolInput no trae el command real: %s", ev.ToolInput)
	}
}

// TestTranslateToolResult: el frame `user` con tool_result debe parear por tool_use_id y traer
// el output. Cubre content-string (ok) y content-array (error).
func TestTranslateToolResult(t *testing.T) {
	ev := firstOfKind(t, frameToolResultOK, ports.EventToolResult)
	if ev.ToolID != "toolu_01E9xFiJFVQorhZsB1z8NrsQ" {
		t.Errorf("ToolID = %q (no parea con la llamada)", ev.ToolID)
	}
	if ev.Text != "PROBE_OK_12345" {
		t.Errorf("Text (output) = %q", ev.Text)
	}
	if ev.ToolIsError {
		t.Errorf("ToolIsError debía ser false")
	}

	errEv := firstOfKind(t, frameToolResultArrErr, ports.EventToolResult)
	if !errEv.ToolIsError {
		t.Errorf("tool_result con is_error:true debía marcar ToolIsError")
	}
	if !strings.Contains(errEv.Text, "boom: file not found") {
		t.Errorf("content-array no extraído: %q", errEv.Text)
	}
}

// TestTranslateSigueEmitiendoDeltaYResult: la refactor a []AgentEvent no rompe los casos previos.
func TestTranslateSigueEmitiendoDeltaYResult(t *testing.T) {
	delta := `{"type":"stream_event","event":{"type":"content_block_delta","delta":{"type":"text_delta","text":"hola"}}}`
	ev := firstOfKind(t, delta, ports.EventDelta)
	if ev.Text != "hola" {
		t.Errorf("delta Text = %q", ev.Text)
	}

	initF := `{"type":"system","subtype":"init","session_id":"e4339525","model":"claude-opus-4-8"}`
	iv := firstOfKind(t, initF, ports.EventInit)
	if iv.ClaudeSessionID != "e4339525" || iv.Model != "claude-opus-4-8" {
		t.Errorf("init mal parseado: %+v", iv)
	}

	res := `{"type":"result","subtype":"success","is_error":false,"result":"listo"}`
	rv := firstOfKind(t, res, ports.EventResult)
	if rv.Text != "listo" {
		t.Errorf("result Text = %q", rv.Text)
	}

	// ruido ignorado: system/status no emite nada
	if _, ok := translate([]byte(`{"type":"system","subtype":"status","status":"requesting"}`)); ok {
		t.Errorf("system/status no debía emitir evento")
	}
}

func TestBuildArgsExploracionEsPlanMode(t *testing.T) {
	args := strings.Join(buildArgs(ports.SpawnOpts{ReadOnly: true}), " ")
	if !strings.Contains(args, "--permission-mode plan") {
		t.Fatalf("exploración debía llevar --permission-mode plan: %s", args)
	}

	trabajo := strings.Join(buildArgs(ports.SpawnOpts{}), " ")
	if strings.Contains(trabajo, "--permission-mode") {
		t.Fatalf("trabajo NO debía restringir permisos: %s", trabajo)
	}
}

func TestBuildArgsResume(t *testing.T) {
	args := strings.Join(buildArgs(ports.SpawnOpts{Resume: "abc"}), " ")
	if !strings.Contains(args, "--resume abc") {
		t.Fatalf("faltó --resume: %s", args)
	}
}

// TestBuildArgsInyeccionArnes (PB-25): la forma-plugin del rol viaja por flags nativos —
// --plugin-dir + UN solo --append-system-prompt (preámbulo + banda Base embebida; la CLI
// rechaza prompt y file a la vez — bug real cazado en el gate DH-18).
func TestBuildArgsInyeccionArnes(t *testing.T) {
	args := strings.Join(buildArgs(ports.SpawnOpts{
		PluginDirs:   []string{"/caché/dev-full-cycle/0.1.0"},
		SystemPrompt: "Trabajás como Ingeniería.\n\n## Banda Base del arnés\n\n# std-spec",
	}), " ")
	for _, quiere := range []string{
		"--plugin-dir /caché/dev-full-cycle/0.1.0",
		"--append-system-prompt Trabajás como Ingeniería.",
		"Banda Base del arnés",
	} {
		if !strings.Contains(args, quiere) {
			t.Errorf("faltó %q en: %s", quiere, args)
		}
	}
	if strings.Contains(args, "--append-system-prompt-file") {
		t.Fatalf("prompt-file y prompt son excluyentes en la CLI — solo debe ir --append-system-prompt: %s", args)
	}

	// sin arnés = cero flags de inyección (sesión legacy/exploración intacta)
	pelado := strings.Join(buildArgs(ports.SpawnOpts{}), " ")
	if strings.Contains(pelado, "--plugin-dir") || strings.Contains(pelado, "--append-system-prompt") {
		t.Fatalf("spawn sin arnés no debía inyectar: %s", pelado)
	}
}
