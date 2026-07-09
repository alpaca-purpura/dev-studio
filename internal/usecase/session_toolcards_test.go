package usecase

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// scriptAgent emite una secuencia fija de eventos y cierra el canal (un turno completo).
// hist = transcript que devuelve History() (R1.5).
type scriptAgent struct {
	evs  []ports.AgentEvent
	hist []ports.TranscriptItem
}

func (a *scriptAgent) Spawn(context.Context, ports.SpawnOpts) (ports.AgentSession, error) {
	ch := make(chan ports.AgentEvent, len(a.evs)+1)
	for _, e := range a.evs {
		ch <- e
	}
	close(ch)
	return &scriptSession{events: ch}, nil
}

func (a *scriptAgent) History(context.Context, string, string) ([]ports.TranscriptItem, error) {
	return a.hist, nil
}

type scriptSession struct{ events chan ports.AgentEvent }

func (s *scriptSession) Send(context.Context, string) error { return nil }
func (s *scriptSession) Events() <-chan ports.AgentEvent    { return s.events }
func (s *scriptSession) Close() error                       { return nil }

// capturePub guarda cada dockFrame publicado (deserializado del JSON que va al SSE).
type capturePub struct {
	mu     sync.Mutex
	frames []dockFrame
}

func (c *capturePub) Publish(_ string, data []byte) {
	var f dockFrame
	if json.Unmarshal(data, &f) != nil {
		return
	}
	c.mu.Lock()
	c.frames = append(c.frames, f)
	c.mu.Unlock()
}

func (c *capturePub) byKind(kind string) []dockFrame {
	c.mu.Lock()
	defer c.mu.Unlock()
	var out []dockFrame
	for _, f := range c.frames {
		if f.Kind == kind {
			out = append(out, f)
		}
	}
	return out
}

// TestConsumePublicaToolCards (R1): un turno que usa una herramienta debe publicar los frames
// tool.call y tool.result al canal realtime, pareados por tool_id, con input/output reales —
// sin tocar el estado de streaming ni el buffer de texto (las cards son inline, no burbujas).
func TestConsumePublicaToolCards(t *testing.T) {
	pub := &capturePub{}
	agent := &scriptAgent{evs: []ports.AgentEvent{
		{Kind: ports.EventInit, ClaudeSessionID: "cc1", Model: "claude-opus-4-8"},
		{Kind: ports.EventToolCall, ToolID: "toolu_1", ToolName: "Read", ToolInput: `{"file_path":"conductor.go"}`},
		{Kind: ports.EventToolResult, ToolID: "toolu_1", Text: "package claudecode", ToolIsError: false},
		{Kind: ports.EventDelta, Text: "Leí el archivo."},
		{Kind: ports.EventResult, Text: "Leí el archivo."},
	}}
	svc, err := NewSessionService(context.Background(), agent, nullStore{}, pub)
	if err != nil {
		t.Fatal(err)
	}
	sess, err := svc.CreateSession(NewSession{Nombre: "s", Cwd: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Turn(sess.ID, "leé el conductor"); err != nil {
		t.Fatal(err)
	}

	// consume corre async; esperar a que llegue el result final del turno.
	waitFor(t, func() bool { return len(pub.byKind("result")) > 0 })

	calls := pub.byKind("tool.call")
	if len(calls) != 1 {
		t.Fatalf("esperaba 1 tool.call, hubo %d: %+v", len(calls), calls)
	}
	if calls[0].ToolName != "Read" || calls[0].ToolID != "toolu_1" {
		t.Errorf("tool.call mal: %+v", calls[0])
	}
	if calls[0].ToolInput != `{"file_path":"conductor.go"}` {
		t.Errorf("tool.call input perdido: %q", calls[0].ToolInput)
	}

	results := pub.byKind("tool.result")
	if len(results) != 1 {
		t.Fatalf("esperaba 1 tool.result, hubo %d", len(results))
	}
	if results[0].ToolID != "toolu_1" {
		t.Errorf("tool.result no parea con la llamada: %+v", results[0])
	}
	if results[0].Text != "package claudecode" {
		t.Errorf("tool.result output perdido: %q", results[0].Text)
	}

	// la sesión debe quedar idle tras el result — las cards no interfieren con el ciclo del turno.
	if got, _ := svc.Get(sess.ID); got.Status != "idle" {
		t.Errorf("status tras el turno = %q, quería idle", got.Status)
	}
}

// runClosedTurn corre un turno cuyo agente cierra el canal (EventInit + EventResult), y espera a
// que la sesión quede idle con su ClaudeSessionID puesto (camino real del consume).
func runClosedTurn(t *testing.T, agent ports.AgentPort, text string) (*SessionService, string) {
	t.Helper()
	svc, err := NewSessionService(context.Background(), agent, nullStore{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	sess, _ := svc.CreateSession(NewSession{Nombre: "s", Cwd: "/tmp"})
	if err := svc.Turn(sess.ID, text); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { g, _ := svc.Get(sess.ID); return g.ClaudeSessionID != "" && g.Status == "idle" })
	return svc, sess.ID
}

// TestTranscriptPathHistory (R1.5): una sesión que ya corrió (tiene ClaudeSessionID) reconstruye
// su transcript ordenado (texto + tool-cards en su lugar) desde History del adapter.
func TestTranscriptPathHistory(t *testing.T) {
	hist := []ports.TranscriptItem{
		{Kind: "user", Text: "leé y editá"},
		{Kind: "tool.call", ToolID: "t1", ToolName: "Read", ToolInput: `{"file_path":"x"}`},
		{Kind: "tool.result", ToolID: "t1", Text: "contenido"},
		{Kind: "assistant", Text: "listo"},
	}
	agent := &scriptAgent{
		evs:  []ports.AgentEvent{{Kind: ports.EventInit, ClaudeSessionID: "cc-xyz"}, {Kind: ports.EventResult, Text: "listo"}},
		hist: hist,
	}
	svc, id := runClosedTurn(t, agent, "hola")
	got, err := svc.Transcript(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(hist) {
		t.Fatalf("esperaba %d ítems del History, hubo %d: %+v", len(hist), len(got), got)
	}
	if got[1].Kind != "tool.call" || got[1].ToolName != "Read" || got[2].Kind != "tool.result" {
		t.Errorf("orden/tool-cards mal reconstruidos: %+v", got)
	}
}

// TestTranscriptFallbackConv (R1.5): si el adapter no tiene historial rico (History vacío), el
// transcript cae al conv persistido solo-texto — mejor que una pantalla en blanco.
func TestTranscriptFallbackConv(t *testing.T) {
	agent := &scriptAgent{ // hist nil → History devuelve vacío
		evs: []ports.AgentEvent{{Kind: ports.EventInit, ClaudeSessionID: "cc-a"}, {Kind: ports.EventResult, Text: "hey"}},
	}
	svc, id := runClosedTurn(t, agent, "hola")
	got, err := svc.Transcript(id)
	if err != nil {
		t.Fatal(err)
	}
	// conv = [user "hola", assistant "hey"]
	if len(got) != 2 || got[0].Kind != "user" || got[0].Text != "hola" || got[1].Kind != "assistant" || got[1].Text != "hey" {
		t.Fatalf("fallback conv esperado [user hola, assistant hey], hubo %+v", got)
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("timeout esperando la condición (¿consume no publicó?)")
}
