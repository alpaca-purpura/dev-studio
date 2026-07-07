// Package cli implementa los puertos GitInfo/GitCommit ejecutando el `git` del usuario como
// subproceso (mismo espíritu BYO del driver DH-10: cero librerías de red, cero credenciales).
//
// Boundary git-solo-lectura-y-commit (RN-4): este paquete NO contiene los verbos
// push / pull / fetch / reset / rebase — la app jamás mueve el repo del usuario contra un
// remoto ni reescribe historia. El fitness test arch/fitness lo escanea.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

// ErrSinPaths: Commit exige pathspec explícito (RN-5).
var ErrSinPaths = errors.New("git commit: se requiere al menos un path explícito")

// Git ejecuta el binario git del usuario.
type Git struct {
	bin string
}

// New crea el adaptador. bin vacío = "git" del PATH.
func New() *Git { return &Git{bin: "git"} }

func (g *Git) run(ctx context.Context, cwd string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, g.bin, args...)
	cmd.Dir = cwd
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(errb.String()))
	}
	return out.String(), nil
}

// Status implementa ports.GitInfo: branch + archivos tocados + shortstat vs HEAD.
func (g *Git) Status(ctx context.Context, cwd string) (domain.GitStatus, error) {
	branch, err := g.run(ctx, cwd, "branch", "--show-current")
	if err != nil {
		return domain.GitStatus{}, err
	}
	st := domain.GitStatus{Branch: strings.TrimSpace(branch), Files: []domain.GitFile{}}

	porcelain, err := g.run(ctx, cwd, "status", "--porcelain")
	if err != nil {
		return domain.GitStatus{}, err
	}
	for _, line := range strings.Split(porcelain, "\n") {
		if len(line) < 4 {
			continue
		}
		x, y := line[0], line[1]
		path := strings.TrimSpace(line[3:])
		// rename: "R  old -> new" — mostramos el destino
		if i := strings.Index(path, " -> "); i >= 0 {
			path = path[i+4:]
		}
		state := "M"
		switch {
		case x == 'U' || y == 'U', x == 'A' && y == 'A', x == 'D' && y == 'D':
			state = "U" // conflicto (unmerged)
		case x == '?' && y == '?':
			state = "A"
		case x == 'A' || y == 'A':
			state = "A"
		case x == 'D' || y == 'D':
			state = "D"
		case x == 'R' || y == 'R':
			state = "R"
		}
		st.Files = append(st.Files, domain.GitFile{Path: path, State: state})
	}

	// shortstat vs HEAD — en repo sin commits HEAD no existe: add/del quedan 0, no es error.
	if shortstat, err := g.run(ctx, cwd, "diff", "HEAD", "--shortstat"); err == nil {
		st.Add, st.Del = parseShortStat(shortstat)
	}
	return st, nil
}

func parseShortStat(s string) (add, del int) {
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		fields := strings.Fields(part)
		if len(fields) < 2 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		switch {
		case strings.Contains(part, "insertion"):
			add = n
		case strings.Contains(part, "deletion"):
			del = n
		}
	}
	return add, del
}

// DiffFile implementa ports.GitInfo: quick-look unificado + panes original/modificado.
func (g *Git) DiffFile(ctx context.Context, cwd, path string) (domain.GitDiff, error) {
	d := domain.GitDiff{Path: path}

	// modificado = el archivo en el working tree (puede no existir si fue eliminado)
	raw, err := os.ReadFile(filepath.Join(cwd, filepath.Clean(path)))
	if err == nil {
		if bytes.IndexByte(raw, 0) >= 0 {
			d.Binary = true
			return d, nil
		}
		d.Modified = string(raw)
	}

	// original = el archivo en HEAD (falla en archivo nuevo o repo sin commits → vacío)
	if orig, err := g.run(ctx, cwd, "show", "HEAD:"+path); err == nil {
		d.Original = orig
	}

	// unificado: tracked → diff HEAD; untracked → diff --no-index contra /dev/null
	if unified, err := g.run(ctx, cwd, "diff", "HEAD", "--", path); err == nil && unified != "" {
		d.Unified = unified
	} else {
		// untracked (o sin HEAD): git diff --no-index sale con código 1 cuando hay diff —
		// exit 1 acá no es error, es "hay diferencias".
		out, nerr := g.runAllowExit1(ctx, cwd, "diff", "--no-index", "--", os.DevNull, path)
		if nerr == nil {
			d.Unified = out
		}
	}
	return d, nil
}

// runAllowExit1: git diff sale 1 cuando hay diferencias — eso no es un fallo.
func (g *Git) runAllowExit1(ctx context.Context, cwd string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, g.bin, args...)
	cmd.Dir = cwd
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	var exitErr *exec.ExitError
	if err != nil && (!errors.As(err, &exitErr) || exitErr.ExitCode() != 1) {
		return "", err
	}
	return out.String(), nil
}

// Log implementa ports.GitInfo.
func (g *Git) Log(ctx context.Context, cwd string, n int) ([]domain.GitLogEntry, error) {
	if n <= 0 {
		n = 30
	}
	out, err := g.run(ctx, cwd, "log", "-n", strconv.Itoa(n), "--pretty=format:%H%x1f%s%x1f%an%x1f%aI%x1e")
	if err != nil {
		// repo sin commits: log falla — historial vacío, no error para la UI
		if strings.Contains(err.Error(), "does not have any commits") {
			return []domain.GitLogEntry{}, nil
		}
		return nil, err
	}
	entries := []domain.GitLogEntry{}
	for _, rec := range strings.Split(out, "\x1e") {
		rec = strings.TrimSpace(rec)
		if rec == "" {
			continue
		}
		f := strings.Split(rec, "\x1f")
		if len(f) != 4 {
			continue
		}
		entries = append(entries, domain.GitLogEntry{SHA: f[0], Mensaje: f[1], Autor: f[2], Fecha: f[3]})
	}
	return entries, nil
}

// Commit implementa ports.GitCommit: stage + commit SOLO de los paths dados (RN-5).
func (g *Git) Commit(ctx context.Context, cwd string, paths []string, mensaje string) (string, error) {
	if len(paths) == 0 {
		return "", ErrSinPaths
	}
	if strings.TrimSpace(mensaje) == "" {
		mensaje = "cambios desde DevStudio"
	}
	addArgs := append([]string{"add", "--"}, paths...)
	if _, err := g.run(ctx, cwd, addArgs...); err != nil {
		return "", err
	}
	commitArgs := append([]string{"commit", "-m", mensaje, "--"}, paths...)
	if _, err := g.run(ctx, cwd, commitArgs...); err != nil {
		return "", err
	}
	sha, err := g.run(ctx, cwd, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(sha), nil
}
