// Command dev-studio es el binario único (Go + UI embebida) del producto: sirve la API local
// de sesiones y el bundle de la SPA desde el mismo proceso.
package main

import (
	"context"
	"errors"
	"flag"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	claudecode "github.com/alpacapurpura/dev-studio/internal/adapters/agent/claudecode"
	gitcli "github.com/alpacapurpura/dev-studio/internal/adapters/git/cli"
	"github.com/alpacapurpura/dev-studio/internal/adapters/store"
	httptransport "github.com/alpacapurpura/dev-studio/internal/adapters/transport/http"
	"github.com/alpacapurpura/dev-studio/internal/adapters/transport/sse"
	"github.com/alpacapurpura/dev-studio/internal/usecase"
	"github.com/alpacapurpura/dev-studio/web"
)

// inyectadas por -ldflags en scripts/install.sh (PB-26)
var (
	version   = "dev"
	buildDate = ""
)

func main() {
	addr := flag.String("addr", "127.0.0.1:4173", "dirección local donde sirve el shell")
	claudeBin := flag.String("claude-bin", "", "ruta al binario claude (vacío = PATH)")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("dev-studio: home dir: %v", err)
	}
	statePath := filepath.Join(home, ".dev-studio", "state.json")
	legacyPath := filepath.Join(home, ".dev-studio", "sessions.json") // F1 → migra solo

	broker := sse.NewBroker()
	agent := claudecode.New(*claudeBin)
	appState := store.NewState(statePath, legacyPath)
	sessions, err := usecase.NewSessionService(ctx, agent, appState, broker)
	if err != nil {
		log.Fatalf("dev-studio: session service: %v", err)
	}
	repos, err := usecase.NewRepoService(ctx, appState, home)
	if err != nil {
		log.Fatalf("dev-studio: repo service: %v", err)
	}
	gitAdapter := gitcli.New()
	git := usecase.NewGitService(gitAdapter, gitAdapter, gitAdapter)

	api := httptransport.NewRouter(sessions, repos, git, broker, httptransport.Deps{
		Home:          home,
		WorkspacesDir: filepath.Join(home, ".dev-studio", "workspaces"),
		Version:       version,
		BuildDate:     buildDate,
		AppConfigPath: filepath.Join(home, ".dev-studio", "app.json"),
		BinPath:       filepath.Join(home, ".local", "bin", "dev-studio"),
	})
	ui := uiHandler()

	mux := http.NewServeMux()
	mux.Handle("/api/", api)
	mux.Handle("/events", api)
	mux.Handle("/", ui)

	srv := &http.Server{Addr: *addr, Handler: mux}
	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()
		_ = srv.Shutdown(shutCtx)
	}()

	slog.Info("dev-studio escuchando", "addr", *addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("dev-studio: serve: %v", err)
	}
}

// uiHandler sirve la SPA embebida desde web/dist, con fallback a index.html (client-side routing).
func uiHandler() http.Handler {
	sub, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatalf("dev-studio: embed sub: %v", err)
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Stat(sub, r.URL.Path[1:]); err != nil {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL.Path = "/"
			fileServer.ServeHTTP(w, r2)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
