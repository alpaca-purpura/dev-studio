package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alpacapurpura/dev-studio/internal/usecase"
)

// Handlers del registry de arneses (PB-25): conexión a nivel app + instalación por repo.

type registryReq struct {
	Source string `json:"source"`
}

func getRegistry(arneses *usecase.ArnesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, arneses.Estado())
	}
}

func putRegistry(arneses *usecase.ArnesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registryReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Source == "" {
			http.Error(w, "bad request: falta source (URL git o ruta local del marketplace)", http.StatusBadRequest)
			return
		}
		estado, err := arneses.Conectar(r.Context(), req.Source)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		writeJSON(w, http.StatusOK, estado)
	}
}

func syncRegistry(arneses *usecase.ArnesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		estado, err := arneses.Resync(r.Context())
		switch {
		case errors.Is(err, usecase.ErrRegistryNoConectado):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case err != nil:
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			writeJSON(w, http.StatusOK, estado)
		}
	}
}

func listArnesesInstalados(arneses *usecase.ArnesService, repos *usecase.RepoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, ok := repos.Get(r.PathValue("id"))
		if !ok {
			http.Error(w, "repo not found", http.StatusNotFound)
			return
		}
		instalados, err := arneses.Instalados(r.Context(), repo.Ruta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, instalados)
	}
}

type instalarReq struct {
	ArnesID string `json:"arnes_id"`
}

func instalarArnes(arneses *usecase.ArnesService, repos *usecase.RepoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, ok := repos.Get(r.PathValue("id"))
		if !ok {
			http.Error(w, "repo not found", http.StatusNotFound)
			return
		}
		var req instalarReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ArnesID == "" {
			http.Error(w, "bad request: falta arnes_id", http.StatusBadRequest)
			return
		}
		arnes, err := arneses.Instalar(r.Context(), repo.Ruta, req.ArnesID)
		switch {
		case errors.Is(err, usecase.ErrRegistryNoConectado), errors.Is(err, usecase.ErrArnesNoEncontrado):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			writeJSON(w, http.StatusCreated, arnes)
		}
	}
}

func desinstalarArnes(arneses *usecase.ArnesService, repos *usecase.RepoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, ok := repos.Get(r.PathValue("id"))
		if !ok {
			http.Error(w, "repo not found", http.StatusNotFound)
			return
		}
		err := arneses.Desinstalar(r.Context(), repo.Ruta, r.PathValue("arnesId"))
		switch {
		case errors.Is(err, usecase.ErrArnesNoInstalado):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
