package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type versionResp struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	Source    string `json:"source"`
}

type appConfig struct {
	Source string `json:"source"`
}

func readAppConfig(path string) appConfig {
	var cfg appConfig
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &cfg)
	}
	return cfg
}

func getVersion(deps Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := readAppConfig(deps.AppConfigPath)
		writeJSON(w, http.StatusOK, versionResp{Version: deps.Version, BuildDate: deps.BuildDate, Source: cfg.Source})
	}
}

// postUpdate — updater dogfooding (spec workspace-aislado §3.2): rebuild del repo local
// (`scripts/install.sh --update`) y restart del proceso vía exec del binario nuevo. RN-4:
// el updater no muta el repo — compila el working tree tal como está (install.sh solo LEE
// la versión con rev-parse). Falla el build → 500 con el output real del compilador.
func postUpdate(deps Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := readAppConfig(deps.AppConfigPath)
		if cfg.Source == "" {
			http.Error(w, "sin source configurado — corré scripts/install.sh del repo primero", http.StatusConflict)
			return
		}
		// bash -lc: PATH del login shell (go/npm), no el PATH mínimo del lanzador .desktop
		cmd := exec.Command("bash", "-lc", "scripts/install.sh --update")
		cmd.Dir = cfg.Source
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			http.Error(w, "el build falló:\n"+out.String(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]string{"detalle": "binario nuevo instalado — reiniciando"})

		// restart: reemplaza el proceso por el binario RECIÉN instalado (mismo addr; la SPA
		// pollea /api/version hasta ver la versión nueva y recarga)
		go func() {
			time.Sleep(400 * time.Millisecond)
			env := os.Environ()
			_ = syscall.Exec(deps.BinPath, os.Args, env)
			// si exec falla (p. ej. corriendo desde `go run` sin instalar), el proceso sigue
			// vivo con la versión vieja — el polling de la SPA lo hace visible.
		}()
	}
}
