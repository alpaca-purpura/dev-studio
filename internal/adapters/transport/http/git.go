package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alpacapurpura/dev-studio/internal/usecase"
)

// sessionCwd resuelve la sesión → su cwd; 404 si no existe.
func sessionCwd(svc *usecase.SessionService, w http.ResponseWriter, r *http.Request) (string, bool) {
	sess, ok := svc.Get(r.PathValue("id"))
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return "", false
	}
	return sess.Cwd, true
}

func sessionGitStatus(sessions *usecase.SessionService, git *usecase.GitService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cwd, ok := sessionCwd(sessions, w, r)
		if !ok {
			return
		}
		st, err := git.Status(r.Context(), cwd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, st)
	}
}

func sessionGitDiff(sessions *usecase.SessionService, git *usecase.GitService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cwd, ok := sessionCwd(sessions, w, r)
		if !ok {
			return
		}
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "bad request: falta path", http.StatusBadRequest)
			return
		}
		d, err := git.DiffFile(r.Context(), cwd, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, d)
	}
}

func sessionGitLog(sessions *usecase.SessionService, git *usecase.GitService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cwd, ok := sessionCwd(sessions, w, r)
		if !ok {
			return
		}
		n, _ := strconv.Atoi(r.URL.Query().Get("n"))
		entries, err := git.Log(r.Context(), cwd, n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, entries)
	}
}

type commitReq struct {
	Paths   []string `json:"paths"`
	Mensaje string   `json:"mensaje"`
}

func sessionGitCommit(sessions *usecase.SessionService, git *usecase.GitService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cwd, ok := sessionCwd(sessions, w, r)
		if !ok {
			return
		}
		var req commitReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if len(req.Paths) == 0 {
			http.Error(w, "bad request: paths vacío (RN-5: commit solo por selección explícita)", http.StatusBadRequest)
			return
		}
		sha, err := git.Commit(r.Context(), cwd, req.Paths, req.Mensaje)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"sha": sha})
	}
}
