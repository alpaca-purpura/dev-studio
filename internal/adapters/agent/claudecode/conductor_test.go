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
