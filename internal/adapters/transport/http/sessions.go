package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

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

func createSession(svc *usecase.SessionService, repos *usecase.RepoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createSessionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		cwd := req.Cwd
		if req.RepoID != "" {
			repo, ok := repos.Get(req.RepoID)
			if !ok {
				http.Error(w, "repo not found", http.StatusUnprocessableEntity)
				return
			}
			// v1 pre-PB-02: la sesión trabaja sobre la raíz del repo; el workspace
			// aislado (worktree+branch) entra en la rebanada siguiente.
			cwd = repo.Ruta
		}
		if cwd == "" {
			cwd, _ = os.UserHomeDir()
		}
		cwd = filepath.Clean(cwd)
		nombre := req.Nombre
		if nombre == "" {
			nombre = "Nueva sesión"
		}
		sess, err := svc.CreateSession(usecase.NewSession{
			Nombre: nombre, Cwd: cwd, RepoID: req.RepoID, Historia: req.Historia, Rol: req.Rol,
		})
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

func deleteSession(svc *usecase.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.Close(r.PathValue("id")); errors.Is(err, usecase.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
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
