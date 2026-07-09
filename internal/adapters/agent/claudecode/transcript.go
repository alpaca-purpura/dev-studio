package claudecode

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alpacapurpura/dev-studio/internal/ports"
)

// jsonlLine es una línea del transcript persistido de Claude Code
// (~/.claude/projects/<cwd-encoded>/<session_id>.jsonl). Mismos bloques de content que el
// stream vivo; envueltos con type/isSidechain. isSidechain=true = turnos internos de un
// subagente (Task) → NO son del transcript principal.
type jsonlLine struct {
	Type        string `json:"type"`
	IsSidechain bool   `json:"isSidechain"`
	Message     struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"` // string suelto o []rawContentBlk
	} `json:"message"`
}

// contentBlocks normaliza el content: string suelto → un bloque text; array → tal cual.
func contentBlocks(raw json.RawMessage) []rawContentBlk {
	if len(raw) == 0 {
		return nil
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return []rawContentBlk{{Type: "text", Text: s}}
	}
	var blocks []rawContentBlk
	if json.Unmarshal(raw, &blocks) == nil {
		return blocks
	}
	return nil
}

// parseTranscript reconstruye el transcript ordenado de las líneas JSONL (en orden de archivo =
// cronológico). Texto consecutivo del mismo rol se fusiona en una burbuja; las tool-cards
// (tool_use → tool.call, tool_result → tool.result) quedan EN SU LUGAR entre las burbujas.
func parseTranscript(lines [][]byte) []ports.TranscriptItem {
	var items []ports.TranscriptItem

	appendText := func(role, text string) {
		if text == "" {
			return
		}
		// fusionar con la burbuja previa si es del mismo rol y no hay una tool en medio
		if n := len(items); n > 0 && items[n-1].Kind == role {
			items[n-1].Text += "\n" + text
			return
		}
		items = append(items, ports.TranscriptItem{Kind: role, Text: text})
	}

	for _, ln := range lines {
		var l jsonlLine
		if json.Unmarshal(ln, &l) != nil {
			continue
		}
		if l.IsSidechain || (l.Type != "user" && l.Type != "assistant") {
			continue // system/summary/attachment/… y ruido de subagente: fuera del transcript
		}
		for _, b := range contentBlocks(l.Message.Content) {
			switch b.Type {
			case "text":
				appendText(l.Type, b.Text) // l.Type = "user" | "assistant"
			case "tool_use":
				items = append(items, ports.TranscriptItem{
					Kind: "tool.call", ToolID: b.ID, ToolName: b.Name, ToolInput: string(b.Input),
				})
			case "tool_result":
				items = append(items, ports.TranscriptItem{
					Kind: "tool.result", ToolID: b.ToolUseID,
					Text: extractContent(b.Content), ToolIsError: b.IsError,
				})
				// thinking/otros: ignorados a propósito
			}
		}
	}
	return items
}

// History (ruta A, R1.5): reconstruye el transcript de una sesión SIN proceso vivo, leyendo el
// JSONL que Claude Code ya persiste para --resume. Boundary conductor-no-parsea-jsonl: este
// paquete es el ÚNICO que conoce el formato. providerSessionID vacío o sin archivo → (nil, nil).
func (c *Conductor) History(_ context.Context, providerSessionID, _ string) ([]ports.TranscriptItem, error) {
	if providerSessionID == "" {
		return nil, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil // sin home no hay dónde buscar — no es error del transcript
	}
	// el session_id es único entre proyectos → glob por nombre de archivo, sin recodificar el cwd.
	matches, _ := filepath.Glob(filepath.Join(home, ".claude", "projects", "*", providerSessionID+".jsonl"))
	if len(matches) == 0 {
		return nil, nil
	}
	lines, err := readLines(matches[0])
	if err != nil {
		return nil, err
	}
	return parseTranscript(lines), nil
}

// readLines lee el JSONL respetando líneas gigantes (el frame init o un tool output largo
// exceden el cap de bufio.Scanner) — mismo criterio que el pump del stream.
func readLines(path string) ([][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 1<<20)
	var out [][]byte
	for {
		line, err := r.ReadBytes('\n')
		if len(line) > 0 {
			cp := make([]byte, len(line))
			copy(cp, line)
			out = append(out, cp)
		}
		if err != nil {
			return out, nil
		}
	}
}
