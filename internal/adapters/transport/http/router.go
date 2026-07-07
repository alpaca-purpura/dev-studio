// Package http expone la SessionService por REST + monta el broker SSE. Único paquete que
// conoce net/http — el dominio y el caso de uso no importan este paquete (boundary:
// dominio-independiente-de-transporte).
package http

import (
	"net"
	"net/http"

	"github.com/alpacapurpura/dev-studio/internal/adapters/transport/sse"
	"github.com/alpacapurpura/dev-studio/internal/usecase"
)

// NewRouter monta la API local (sesiones + repos + git + eventos) sobre un http.ServeMux.
func NewRouter(sessions *usecase.SessionService, repos *usecase.RepoService, git *usecase.GitService, broker *sse.Broker) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/sessions", listSessions(sessions))
	mux.HandleFunc("POST /api/sessions", createSession(sessions, repos))
	mux.HandleFunc("GET /api/sessions/{id}", getSession(sessions))
	mux.HandleFunc("PATCH /api/sessions/{id}", patchSession(sessions))
	mux.HandleFunc("DELETE /api/sessions/{id}", deleteSession(sessions))
	mux.HandleFunc("POST /api/sessions/{id}/turn", sessionTurn(sessions))

	mux.HandleFunc("GET /api/sessions/{id}/git/status", sessionGitStatus(sessions, git))
	mux.HandleFunc("GET /api/sessions/{id}/git/diff", sessionGitDiff(sessions, git))
	mux.HandleFunc("GET /api/sessions/{id}/git/log", sessionGitLog(sessions, git))
	mux.HandleFunc("POST /api/sessions/{id}/git/commit", sessionGitCommit(sessions, git))

	mux.HandleFunc("GET /api/repos", listRepos(repos))
	mux.HandleFunc("POST /api/repos", registerRepo(repos))
	mux.HandleFunc("DELETE /api/repos/{id}", deleteRepo(repos))
	mux.HandleFunc("GET /api/repos/{id}/git/status", repoGitStatus(repos, git))

	mux.HandleFunc("GET /events", broker.ServeHTTP)

	return withLocalOnly(mux)
}

// withLocalOnly confina la API a peticiones del propio shell (Host allowlist) — MVP: solo
// localhost/127.0.0.1. Sin token de capacidad todavía (TBD, ver arch/boundaries).
func withLocalOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		if host != "localhost" && host != "127.0.0.1" && host != "" {
			http.Error(w, "forbidden host", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
