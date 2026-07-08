// Package claudecode es el único adaptador que conoce el protocolo stream-json de la CLI `claude`.
// Ningún otro paquete parsea sus frames — boundary: conductor-no-parsea-jsonl.
package claudecode

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// Conductor spawnea procesos `claude` en modo headless (stream-json por stdin/stdout).
// BYO licencia: usa el binario que el usuario ya tiene autenticado, nunca toca credenciales (DH-10).
type Conductor struct {
	bin string
}

// Aserción en tiempo de compilación: Conductor es un AgentPort — boundary
// adaptador-agente-intercambiable (arch/boundaries/adaptador-agente-intercambiable.md).
var _ ports.AgentPort = (*Conductor)(nil)

// New crea un Conductor. bin vacío = resolver "claude" del PATH con fallbacks a las rutas
// de instalación típicas — lanzada desde un .desktop, la app hereda un PATH pelado que no
// incluye ~/.local/bin (bug real cazado en dogfooding, DH-16.1).
func New(bin string) *Conductor {
	if bin == "" {
		bin = resolveClaudeBin()
	}
	return &Conductor{bin: bin}
}

func resolveClaudeBin() string {
	if p, err := exec.LookPath("claude"); err == nil {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "claude"
	}
	for _, cand := range []string{
		filepath.Join(home, ".local", "bin", "claude"),
		filepath.Join(home, "bin", "claude"),
		filepath.Join(home, ".npm-global", "bin", "claude"),
		"/usr/local/bin/claude",
		"/opt/homebrew/bin/claude",
	} {
		if fi, err := os.Stat(cand); err == nil && !fi.IsDir() {
			return cand
		}
	}
	return "claude" // que el error de spawn lo diga claro — la UI ya lo muestra
}

// buildArgs arma la línea del subproceso — pura, testeable. Exploración (PB-27) usa el
// mecanismo NATIVO read-only de la CLI: --permission-mode plan (lee/analiza, no edita).
func buildArgs(opts ports.SpawnOpts) []string {
	args := []string{
		"-p",
		"--input-format", "stream-json",
		"--output-format", "stream-json",
		"--include-partial-messages",
		"--verbose",
	}
	if opts.ReadOnly {
		args = append(args, "--permission-mode", "plan")
	}
	if opts.Resume != "" {
		args = append(args, "--resume", opts.Resume)
	}
	// Arnés del rol (PB-25): forma-plugin por --plugin-dir + un solo system prompt
	// (preámbulo + banda Base embebida — la CLI no acepta prompt y file a la vez).
	for _, d := range opts.PluginDirs {
		args = append(args, "--plugin-dir", d)
	}
	if opts.SystemPrompt != "" {
		args = append(args, "--append-system-prompt", opts.SystemPrompt)
	}
	return args
}

func (c *Conductor) Spawn(ctx context.Context, opts ports.SpawnOpts) (ports.AgentSession, error) {
	cmd := exec.CommandContext(ctx, c.bin, buildArgs(opts)...)
	cmd.Dir = opts.Cwd

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("claudecode: stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("claudecode: stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("claudecode: stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("claudecode: start %s: %w", c.bin, err)
	}

	s := &ccSession{
		cmd:    cmd,
		stdin:  stdin,
		events: make(chan ports.AgentEvent, 64),
	}
	go s.logStderr(stderr)
	go s.pump(stdout)
	return s, nil
}

type ccSession struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	events chan ports.AgentEvent

	sendMu    sync.Mutex
	closeOnce sync.Once
	closeErr  error
}

// userFrame es el frame NDJSON de entrada que espera --input-format stream-json.
type userFrame struct {
	Type    string `json:"type"`
	Message struct {
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
}

func (s *ccSession) Send(ctx context.Context, text string) error {
	s.sendMu.Lock()
	defer s.sendMu.Unlock()

	f := userFrame{Type: "user"}
	f.Message.Role = "user"
	f.Message.Content = []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}{{Type: "text", Text: text}}

	line, err := json.Marshal(f)
	if err != nil {
		return fmt.Errorf("claudecode: marshal turn: %w", err)
	}
	line = append(line, '\n')
	_, err = s.stdin.Write(line)
	return err
}

func (s *ccSession) Events() <-chan ports.AgentEvent { return s.events }

func (s *ccSession) Close() error {
	s.closeOnce.Do(func() {
		_ = s.stdin.Close() // EOF limpio: claude termina el proceso solo
		s.closeErr = s.cmd.Wait()
	})
	return s.closeErr
}

func (s *ccSession) logStderr(r io.Reader) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		slog.Debug("claudecode stderr", "line", sc.Text())
	}
}

// rawFrame es el superset de campos que nos interesan de los frames stream-json de salida.
type rawFrame struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	// system/init
	SessionID string `json:"session_id"`
	Model     string `json:"model"`
	// stream_event
	Event struct {
		Type  string `json:"type"`
		Delta struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"delta"`
	} `json:"event"`
	// result
	Result  string `json:"result"`
	IsError bool   `json:"is_error"`
}

// pump lee stdout NDJSON y traduce cada línea a un ports.AgentEvent normalizado.
// bufio.Reader (no Scanner) porque el frame `init` trae la lista completa de tools/skills
// y excede el cap de 64KB por defecto de bufio.Scanner.
func (s *ccSession) pump(stdout io.Reader) {
	defer close(s.events)
	r := bufio.NewReaderSize(stdout, 1<<20)
	for {
		line, err := r.ReadBytes('\n')
		if len(line) > 0 {
			if ev, ok := translate(line); ok {
				s.events <- ev // send bloqueante a propósito: perder un `result` cuelga la UI en streaming
			}
		}
		if err != nil {
			return
		}
	}
}

func translate(line []byte) (ports.AgentEvent, bool) {
	var f rawFrame
	if err := json.Unmarshal(line, &f); err != nil {
		return ports.AgentEvent{}, false
	}
	switch {
	case f.Type == "system" && f.Subtype == "init":
		return ports.AgentEvent{Kind: ports.EventInit, ClaudeSessionID: f.SessionID, Model: f.Model}, true
	case f.Type == "stream_event" && f.Event.Type == "content_block_delta" && f.Event.Delta.Type == "text_delta":
		return ports.AgentEvent{Kind: ports.EventDelta, Text: f.Event.Delta.Text}, true
	case f.Type == "result":
		if f.IsError {
			return ports.AgentEvent{Kind: ports.EventError, Err: f.Result}, true
		}
		return ports.AgentEvent{Kind: ports.EventResult, Text: f.Result}, true
	default:
		return ports.AgentEvent{}, false // system/status, rate_limit_event, hooks, assistant-completo: ignorados a propósito
	}
}
