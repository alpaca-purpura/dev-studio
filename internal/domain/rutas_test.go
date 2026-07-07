package domain_test

import (
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/domain"
)

func TestRutaProtegida(t *testing.T) {
	home := "/home/usuaria"
	protegidas := []string{
		"/home/usuaria",           // $HOME exacto
		"/home/usuaria/.ssh",      // credenciales
		"/home/usuaria/.ssh/keys", // dentro de una protegida
		"/home/usuaria/.gnupg",
		"/home/usuaria/.aws",
		"/home/usuaria/.kube",
		"/home/usuaria/.docker",
		"/home/usuaria/.dev-studio", // la app no opera sobre sí misma
		"/", "/etc", "/etc/nginx", "/usr", "/usr/local", "/bin", "/sbin",
		"/var", "/var/log", "/boot", "/root", "/proc", "/sys", "/dev",
	}
	for _, p := range protegidas {
		if !domain.RutaProtegida(home, p) {
			t.Errorf("%s debía ser protegida", p)
		}
	}

	permitidas := []string{
		"/home/usuaria/Proyectos/dev-studio",
		"/home/usuaria/repos/x",
		"/home/usuaria/.config/misrepos/uno", // no está en la denylist
		"/tmp/scratch/demo",
		"/opt/repos/legacy",
	}
	for _, p := range permitidas {
		if domain.RutaProtegida(home, p) {
			t.Errorf("%s debía permitirse", p)
		}
	}
}
