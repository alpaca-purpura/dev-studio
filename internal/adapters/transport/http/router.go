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

// Deps agrupa lo que el router necesita además de los servicios (config del host).
type Deps struct {
	Home          string
	WorkspacesDir string // ~/.dev-studio/workspaces
	Version       string
	BuildDate     string
	AppConfigPath string // ~/.dev-studio/app.json (source del updater)
	BinPath       string // ~/.local/bin/dev-studio (destino del rebuild + exec)
}

// NewRouter monta la API local (sesiones + repos + git + registry de arneses +
// versión/update + eventos).
func NewRouter(sessions *usecase.SessionService, repos *usecase.RepoService, git *usecase.GitService, arneses *usecase.ArnesService, broker *sse.Broker, deps Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/sessions", listSessions(sessions))
	mux.HandleFunc("POST /api/sessions", createSession(sessions, repos, git, arneses, deps.Home, deps.WorkspacesDir))
	mux.HandleFunc("GET /api/sessions/{id}", getSession(sessions))
	mux.HandleFunc("PATCH /api/sessions/{id}", patchSession(sessions))
	mux.HandleFunc("DELETE /api/sessions/{id}", deleteSession(sessions, repos, git))
	mux.HandleFunc("POST /api/sessions/{id}/turn", sessionTurn(sessions))
	mux.HandleFunc("GET /api/sessions/{id}/transcript", sessionTranscript(sessions))

	mux.HandleFunc("GET /api/version", getVersion(deps))
	mux.HandleFunc("POST /api/update", postUpdate(deps))

	mux.HandleFunc("GET /api/sessions/{id}/git/status", sessionGitStatus(sessions, git))
	mux.HandleFunc("GET /api/sessions/{id}/git/diff", sessionGitDiff(sessions, git))
	mux.HandleFunc("GET /api/sessions/{id}/git/log", sessionGitLog(sessions, git))
	mux.HandleFunc("POST /api/sessions/{id}/git/commit", sessionGitCommit(sessions, git))

	mux.HandleFunc("GET /api/repos", listRepos(repos))
	mux.HandleFunc("POST /api/repos", registerRepo(repos))
	mux.HandleFunc("DELETE /api/repos/{id}", deleteRepo(repos))
	mux.HandleFunc("GET /api/repos/{id}/git/status", repoGitStatus(repos, git))

	// registry de arneses (PB-25): conexión a nivel app + roster por repo
	mux.HandleFunc("GET /api/registry", getRegistry(arneses))
	mux.HandleFunc("PUT /api/registry", putRegistry(arneses))
	mux.HandleFunc("POST /api/registry/sync", syncRegistry(arneses))
	mux.HandleFunc("GET /api/repos/{id}/arneses", listArnesesInstalados(arneses, repos))
	mux.HandleFunc("POST /api/repos/{id}/arneses", instalarArnes(arneses, repos))
	mux.HandleFunc("DELETE /api/repos/{id}/arneses/{arnesId}", desinstalarArnes(arneses, repos))

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
