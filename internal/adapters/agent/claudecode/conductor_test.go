package claudecode

import (
	"strings"
	"testing"

	"github.com/alpacapurpura/dev-studio/internal/ports"
)

func TestBuildArgsExploracionEsPlanMode(t *testing.T) {
	args := strings.Join(buildArgs(ports.SpawnOpts{ReadOnly: true}), " ")
	if !strings.Contains(args, "--permission-mode plan") {
		t.Fatalf("exploración debía llevar --permission-mode plan: %s", args)
	}

	trabajo := strings.Join(buildArgs(ports.SpawnOpts{}), " ")
	if strings.Contains(trabajo, "--permission-mode") {
		t.Fatalf("trabajo NO debía restringir permisos: %s", trabajo)
	}
}

func TestBuildArgsResume(t *testing.T) {
	args := strings.Join(buildArgs(ports.SpawnOpts{Resume: "abc"}), " ")
	if !strings.Contains(args, "--resume abc") {
		t.Fatalf("faltó --resume: %s", args)
	}
}

// TestBuildArgsInyeccionArnes (PB-25): la forma-plugin del rol viaja por flags nativos —
// --plugin-dir + UN solo --append-system-prompt (preámbulo + banda Base embebida; la CLI
// rechaza prompt y file a la vez — bug real cazado en el gate DH-18).
func TestBuildArgsInyeccionArnes(t *testing.T) {
	args := strings.Join(buildArgs(ports.SpawnOpts{
		PluginDirs:   []string{"/caché/dev-full-cycle/0.1.0"},
		SystemPrompt: "Trabajás como Ingeniería.\n\n## Banda Base del arnés\n\n# std-spec",
	}), " ")
	for _, quiere := range []string{
		"--plugin-dir /caché/dev-full-cycle/0.1.0",
		"--append-system-prompt Trabajás como Ingeniería.",
		"Banda Base del arnés",
	} {
		if !strings.Contains(args, quiere) {
			t.Errorf("faltó %q en: %s", quiere, args)
		}
	}
	if strings.Contains(args, "--append-system-prompt-file") {
		t.Fatalf("prompt-file y prompt son excluyentes en la CLI — solo debe ir --append-system-prompt: %s", args)
	}

	// sin arnés = cero flags de inyección (sesión legacy/exploración intacta)
	pelado := strings.Join(buildArgs(ports.SpawnOpts{}), " ")
	if strings.Contains(pelado, "--plugin-dir") || strings.Contains(pelado, "--append-system-prompt") {
		t.Fatalf("spawn sin arnés no debía inyectar: %s", pelado)
	}
}
