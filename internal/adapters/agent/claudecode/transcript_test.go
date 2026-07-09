package claudecode

import (
	"strings"
	"testing"
)

// Fixture con el shape REAL del JSONL de Claude Code (verificado en ~/.claude/projects, 2026-07-09):
// líneas {type,isSidechain,message:{role,content:[...]}}; content string o array; sidechain = ruido.
var transcriptFixture = [][]byte{
	[]byte(`{"type":"user","uuid":"u1","message":{"role":"user","content":"leé hello.txt y editalo"}}`), // content STRING
	[]byte(`{"type":"assistant","uuid":"a1","message":{"role":"assistant","content":[{"type":"text","text":"Voy a leerlo."}]}}`),
	[]byte(`{"type":"assistant","uuid":"a1b","message":{"role":"assistant","content":[{"type":"text","text":"Sigo."}]}}`), // consecutivo → fusiona
	[]byte(`{"type":"assistant","uuid":"a2","message":{"role":"assistant","content":[{"type":"tool_use","id":"toolu_1","name":"Read","input":{"file_path":"hello.txt"}}]}}`),
	[]byte(`{"type":"user","uuid":"u2","message":{"role":"user","content":[{"tool_use_id":"toolu_1","type":"tool_result","content":"line one\nline two","is_error":false}]}}`),
	[]byte(`{"type":"assistant","uuid":"sc","isSidechain":true,"message":{"role":"assistant","content":[{"type":"text","text":"RUIDO subagente"}]}}`), // skip
	[]byte(`{"type":"summary","summary":"x"}`), // skip
	[]byte(`{"type":"assistant","uuid":"a3","message":{"role":"assistant","content":[{"type":"tool_use","id":"toolu_2","name":"Edit","input":{"file_path":"hello.txt","old_string":"a","new_string":"b"}}]}}`),
	[]byte(`{"type":"user","uuid":"u3","message":{"role":"user","content":[{"tool_use_id":"toolu_2","type":"tool_result","content":[{"type":"text","text":"permiso denegado"}],"is_error":true}]}}`),
	[]byte(`{"type":"assistant","uuid":"a4","message":{"role":"assistant","content":[{"type":"text","text":"Listo."}]}}`),
}

// TestParseTranscript (R1.5, ruta A): reconstrucción ordenada — texto+tool-cards en su lugar,
// texto consecutivo fusionado, sidechain/summary descartados, content-string y content-array.
func TestParseTranscript(t *testing.T) {
	items := parseTranscript(transcriptFixture)

	type want struct {
		kind, text, toolName string
		isErr                bool
	}
	expected := []want{
		{kind: "user", text: "leé hello.txt y editalo"},
		{kind: "assistant", text: "Voy a leerlo.\nSigo."}, // fusionado
		{kind: "tool.call", toolName: "Read"},
		{kind: "tool.result", text: "line one\nline two"},
		{kind: "tool.call", toolName: "Edit"},
		{kind: "tool.result", text: "permiso denegado", isErr: true}, // content-array + error
		{kind: "assistant", text: "Listo."},
	}

	if len(items) != len(expected) {
		t.Fatalf("esperaba %d ítems, hubo %d: %+v", len(expected), len(items), items)
	}
	for i, e := range expected {
		got := items[i]
		if got.Kind != e.kind {
			t.Errorf("ítem %d: Kind=%q, quería %q", i, got.Kind, e.kind)
		}
		if e.text != "" && got.Text != e.text {
			t.Errorf("ítem %d: Text=%q, quería %q", i, got.Text, e.text)
		}
		if e.toolName != "" && got.ToolName != e.toolName {
			t.Errorf("ítem %d: ToolName=%q, quería %q", i, got.ToolName, e.toolName)
		}
		if got.Kind == "tool.result" && got.ToolIsError != e.isErr {
			t.Errorf("ítem %d: ToolIsError=%v, quería %v", i, got.ToolIsError, e.isErr)
		}
	}

	// la tool.call Read debe conservar su input real (para arg/diff en la card)
	if !strings.Contains(items[2].ToolInput, `"file_path":"hello.txt"`) {
		t.Errorf("tool.call Read sin input real: %q", items[2].ToolInput)
	}
	// pareo call↔result por ToolID
	if items[2].ToolID != items[3].ToolID || items[2].ToolID != "toolu_1" {
		t.Errorf("Read call/result no parean: %q vs %q", items[2].ToolID, items[3].ToolID)
	}
}
