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
	"github.com/alpacapurpura/dev-studio/internal/adapters/store"
	httptransport "github.com/alpacapurpura/dev-studio/internal/adapters/transport/http"
	"github.com/alpacapurpura/dev-studio/internal/adapters/transport/sse"
	"github.com/alpacapurpura/dev-studio/internal/usecase"
	"github.com/alpacapurpura/dev-studio/web"
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
	registryPath := filepath.Join(home, ".dev-studio", "sessions.json")

	broker := sse.NewBroker()
	agent := claudecode.New(*claudeBin)
	sessions, err := usecase.NewSessionService(ctx, agent, store.NewRegistry(registryPath), broker)
	if err != nil {
		log.Fatalf("dev-studio: session service: %v", err)
	}

	api := httptransport.NewRouter(sessions, broker)
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
