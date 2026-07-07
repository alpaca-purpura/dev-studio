package domain

// Repo es un repositorio git registrado en la app — el contenedor de workspaces/sesiones
// (spec shell §5). La Ruta es la raíz local del repo (contiene .git).
type Repo struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Ruta   string `json:"ruta"`
}

// TipoItem clasifica un paquete de trabajo según el estándar de industria (spec
// nuevo-workspace §1: Jira issue types ∩ Gitflow branch naming ∩ Conventional Commits).
type TipoItem string

const (
	TipoHistoria TipoItem = "historia" // feature/… · feat:
	TipoBug      TipoItem = "bug"      // bugfix/…  · fix:
	TipoHotfix   TipoItem = "hotfix"   // hotfix/…  · fix:
	TipoTarea    TipoItem = "tarea"    // chore/…   · chore:/refactor:
	TipoSpike    TipoItem = "spike"    // spike/…   · investigación
)

// BranchPrefix devuelve el prefijo de branch estándar del tipo (default: feature).
func (t TipoItem) BranchPrefix() string {
	switch t {
	case TipoBug:
		return "bugfix"
	case TipoHotfix:
		return "hotfix"
	case TipoTarea:
		return "chore"
	case TipoSpike:
		return "spike"
	default:
		return "feature"
	}
}

// Historia es la referencia liviana al paquete de trabajo que una sesión CON EDICIÓN liga
// (RN-1 evolucionada). v1: viaja y persiste con la sesión; board global de ítems = PB-07.
type Historia struct {
	ID     string   `json:"id"`
	Titulo string   `json:"titulo"`
	Tipo   TipoItem `json:"tipo,omitempty"` // vacío = historia (legacy)
}
