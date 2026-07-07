package domain

// Repo es un repositorio git registrado en la app — el contenedor de workspaces/sesiones
// (spec shell §5). La Ruta es la raíz local del repo (contiene .git).
type Repo struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Ruta   string `json:"ruta"`
}

// Historia es la referencia liviana al paquete de trabajo que una sesión liga (RN-1).
// v1: viene del board de ejemplo (mock); el contrato ya es el definitivo.
type Historia struct {
	ID     string `json:"id"`
	Titulo string `json:"titulo"`
}
