package domain

import "fmt"

// Arnes es un rol publicado en el registry de la organización, en el formato de la
// fábrica (ArnesIA, nomenclatura-arnes v1): forma-plugin con manifiesto arnes.l0.json
// (meta rol×proceso + spine) + plugin.json. DevStudio lo CONSUME e interpreta — jamás
// lo crea ni lo edita (RN-4 spec registry-arneses; DH-14: cero roles locales).
type Arnes struct {
	ID          string   `json:"id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion,omitempty"`
	Rol         string   `json:"rol"`
	Proceso     string   `json:"proceso,omitempty"`
	Version     string   `json:"version"`
	Canal       string   `json:"canal,omitempty"`
	Fases       []string `json:"fases,omitempty"`
}

// Preambulo compone el system prompt del rol que se inyecta al spawn
// (--append-system-prompt). Cero campo nuevo pedido a la fábrica: se deriva del meta.
func (a Arnes) Preambulo() string {
	p := fmt.Sprintf("Trabajás como %s.", a.Rol)
	if a.Proceso != "" {
		p += fmt.Sprintf(" Tu rol sirve al proceso: %s.", a.Proceso)
	}
	if len(a.Fases) > 0 {
		p += " Fases del proceso: "
		for i, f := range a.Fases {
			if i > 0 {
				p += " → "
			}
			p += f
		}
		p += "."
	}
	p += " Las skills de tu arnés definen tu forma de trabajar — usalas."
	return p
}

// ArnesInstalado es una entrada del lock del proyecto (.devstudio/arneses.yaml): el
// roster as-code que viaja por el repo. El payload NO viaja — se rehidrata del registry.
type ArnesInstalado struct {
	ID      string `json:"id" yaml:"id"`
	Version string `json:"version" yaml:"version"`
	Canal   string `json:"canal,omitempty" yaml:"canal,omitempty"`
}

// RegistryEstado es la vista del registry conectado que consume la UI.
type RegistryEstado struct {
	Source  string  `json:"source"`
	Arneses []Arnes `json:"arneses"`
}
