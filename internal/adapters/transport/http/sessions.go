package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/usecase"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func listSessions(svc *usecase.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, svc.List())
	}
}

type createSessionReq struct {
	Nombre   string           `json:"nombre"`
	Cwd      string           `json:"cwd"`
	RepoID   string           `json:"repo_id"`
	Historia *domain.Historia `json:"historia"`
	Rol      string           `json:"rol"`
}

// slugify normaliza el nombre a slug de branch/carpeta (misma regla que la SPA).
func slugify(s string) string {
	var b []rune
	prev := '-'
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b = append(b, r)
			prev = r
		default:
			if prev != '-' {
				b = append(b, '-')
				prev = '-'
			}
		}
	}
	out := strings.Trim(string(b), "-")
	if len(out) > 26 {
		out = strings.Trim(out[:26], "-")
	}
	if out == "" {
		out = "sesion"
	}
	return out
}

func createSession(svc *usecase.SessionService, repos *usecase.RepoService, git *usecase.GitService, home, workspacesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createSessionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		nombre := req.Nombre
		if nombre == "" {
			nombre = "Nueva sesión"
		}

		p := usecase.NewSession{Nombre: nombre, RepoID: req.RepoID, Historia: req.Historia, Rol: req.Rol}
		if req.RepoID != "" {
			// PB-02: toda sesión de repo nace en su PROPIO workspace (worktree + wt/{slug})
			repo, ok := repos.Get(req.RepoID)
			if !ok {
				http.Error(w, "repo not found", http.StatusUnprocessableEntity)
				return
			}
			slug := slugify(nombre)
			destino := filepath.Join(workspacesDir, repo.Nombre, slug)
			path, branch, err := git.CreateWorktree(r.Context(), repo.Ruta, destino, slug)
			if err != nil {
				http.Error(w, "workspace: "+err.Error(), http.StatusUnprocessableEntity)
				return
			}
			p.Cwd, p.Workspace, p.Branch = path, path, branch
		} else {
			cwd := req.Cwd
			if cwd == "" {
				cwd, _ = os.UserHomeDir()
			}
			cwd = filepath.Clean(cwd)
			// RN-2: ninguna sesión en ruta protegida, tampoco por la puerta legacy
			if err := usecase.ValidarRuta(home, cwd); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			p.Cwd = cwd
		}

		sess, err := svc.CreateSession(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, sess)
	}
}

func getSession(svc *usecase.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, ok := svc.Get(r.PathValue("id"))
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, sess)
	}
}

type patchSessionReq struct {
	Nombre *string `json:"nombre,omitempty"`
}

func patchSession(svc *usecase.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req patchSessionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		id := r.PathValue("id")
		if req.Nombre != nil {
			if err := svc.Rename(id, *req.Nombre); errors.Is(err, usecase.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
		}
		sess, _ := svc.Get(id)
		writeJSON(w, http.StatusOK, sess)
	}
}

type deleteSessionResp struct {
	Closed           bool   `json:"closed"`
	WorkspaceRemoved bool   `json:"workspace_removed"`
	Detalle          string `json:"detalle,omitempty"`
}

// deleteSession cierra la sesión; ?workspace=remove intenta borrar su worktree (RN-3: git
// rechaza el remove sucio — la app conserva y avisa, jamás --force).
func deleteSession(svc *usecase.SessionService, repos *usecase.RepoService, git *usecase.GitService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, ok := svc.Get(r.PathValue("id"))
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if err := svc.Close(sess.ID); errors.Is(err, usecase.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		resp := deleteSessionResp{Closed: true}
		if r.URL.Query().Get("workspace") == "remove" && sess.Workspace != "" {
			repoRoot := sess.Workspace // fallback: el remove corre desde el propio worktree
			if repo, ok := repos.Get(sess.RepoID); ok {
				repoRoot = repo.Ruta
			}
			if err := git.RemoveWorktree(r.Context(), repoRoot, sess.Workspace); err != nil {
				resp.Detalle = "el workspace se conservó: " + err.Error()
			} else {
				resp.WorkspaceRemoved = true
			}
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

type turnReq struct {
	Text string `json:"text"`
}

func sessionTurn(svc *usecase.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req turnReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		err := svc.Turn(r.PathValue("id"), req.Text)
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			http.Error(w, "not found", http.StatusNotFound)
		case errors.Is(err, usecase.ErrBusy):
			http.Error(w, "busy", http.StatusConflict)
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusAccepted)
		}
	}
}
