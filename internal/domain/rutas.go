package domain

import (
	"path/filepath"
	"strings"
)

// raícesProtegidas: directorios de sistema donde la app JAMÁS registra un repo ni pone el
// cwd de una sesión (boundary sesion-aislada-por-cwd, spec workspace-aislado §2.3).
var raicesProtegidas = []string{
	"/etc", "/usr", "/bin", "/sbin", "/var", "/boot", "/root", "/proc", "/sys", "/dev",
}

// dotDirsProtegidos: subdirectorios sensibles del home (credenciales/estado de la app).
var dotDirsProtegidos = []string{
	".ssh", ".gnupg", ".aws", ".kube", ".docker", ".dev-studio",
}

// RutaProtegida decide si una ruta está vedada como repo/cwd de sesión. Pura (el home llega
// como argumento) — el dominio no toca el OS.
func RutaProtegida(home, ruta string) bool {
	p := filepath.Clean(ruta)
	if p == "/" || p == filepath.Clean(home) {
		return true
	}
	for _, raiz := range raicesProtegidas {
		if p == raiz || strings.HasPrefix(p, raiz+"/") {
			return true
		}
	}
	h := filepath.Clean(home)
	for _, dot := range dotDirsProtegidos {
		full := filepath.Join(h, dot)
		if p == full || strings.HasPrefix(p, full+"/") {
			return true
		}
	}
	return false
}
