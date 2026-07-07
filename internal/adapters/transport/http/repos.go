package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alpacapurpura/dev-studio/internal/usecase"
)

func listRepos(svc *usecase.RepoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, svc.List())
	}
}

type registerRepoReq struct {
	Ruta string `json:"ruta"`
}

func registerRepo(svc *usecase.RepoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRepoReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Ruta == "" {
			http.Error(w, "bad request: falta ruta", http.StatusBadRequest)
			return
		}
		repo, err := svc.Register(req.Ruta)
		if errors.Is(err, usecase.ErrNoEsRepo) || errors.Is(err, usecase.ErrRutaProtegida) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, repo)
	}
}

func deleteRepo(svc *usecase.RepoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.Remove(r.PathValue("id")); errors.Is(err, usecase.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// repoGitStatus sirve el status del repo raíz (para el rail: branch + ±N por workspace
// cuando la sesión comparte el cwd del repo — v1, pre PB-02).
func repoGitStatus(repos *usecase.RepoService, git *usecase.GitService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, ok := repos.Get(r.PathValue("id"))
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		st, err := git.Status(r.Context(), repo.Ruta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, st)
	}
}
