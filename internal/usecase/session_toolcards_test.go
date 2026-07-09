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
type scriptAgent struct{ evs []ports.AgentEvent }

func (a *scriptAgent) Spawn(context.Context, ports.SpawnOpts) (ports.AgentSession, error) {
	ch := make(chan ports.AgentEvent, len(a.evs)+1)
	for _, e := range a.evs {
		ch <- e
	}
	close(ch)
	return &scriptSession{events: ch}, nil
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
