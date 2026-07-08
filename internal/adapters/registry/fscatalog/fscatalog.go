// Package fscatalog lee el catálogo de arneses de un marketplace git materializado en
// disco. Formato consumido = el estándar de la fábrica (ArnesIA) YA en producción:
// `.claude-plugin/marketplace.json` (índice de plugins) + por arnés `arnes.l0.json`
// (nomenclatura-arnes v1: meta rol×proceso + fases + spine) y `.claude-plugin/plugin.json`.
// Solo LECTURA — DevStudio consume arneses, jamás los crea ni edita (RN-4).
package fscatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alpacapurpura/dev-studio/internal/domain"
	"github.com/alpacapurpura/dev-studio/internal/ports"
)

type Catalog struct{}

var _ ports.RegistryCatalog = (*Catalog)(nil)

func New() *Catalog { return &Catalog{} }

// marketplaceFile es .claude-plugin/marketplace.json (subset que DevStudio necesita).
type marketplaceFile struct {
	Name    string `json:"name"`
	Plugins []struct {
		Name        string `json:"name"`
		Source      string `json:"source"`
		Description string `json:"description"`
	} `json:"plugins"`
}

// pluginFile es .claude-plugin/plugin.json del arnés.
type pluginFile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// arnesL0 es el subset de arnes.l0.json que DevStudio interpreta en esta rebanada.
// El spine queda EN el archivo del caché para PB-08 (proceso as code) — acá no se parsea.
type arnesL0 struct {
	ID          string   `json:"id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Rol         string   `json:"rol"`
	Proceso     string   `json:"proceso"`
	Canal       string   `json:"canal"`
	Fases       []string `json:"fases"`
}

// Leer devuelve los arneses del marketplace en el orden del índice. Un marketplace sin
// índice = error honesto, jamás un catálogo vacío inventado.
func (c *Catalog) Leer(_ context.Context, dir string) ([]domain.Arnes, error) {
	mp, err := leerMarketplace(dir)
	if err != nil {
		return nil, err
	}
	arneses := make([]domain.Arnes, 0, len(mp.Plugins))
	for _, p := range mp.Plugins {
		src := filepath.Join(dir, filepath.Clean(p.Source))
		a, err := leerArnes(src)
		if err != nil {
			return nil, fmt.Errorf("fscatalog: arnés %q: %w", p.Name, err)
		}
		if a.ID == "" {
			a.ID = p.Name
		}
		if a.Descripcion == "" {
			a.Descripcion = p.Description
		}
		arneses = append(arneses, a)
	}
	return arneses, nil
}

// SourceDir resuelve el directorio forma-plugin de un arnés del índice (para materializar).
func (c *Catalog) SourceDir(dir, arnesID string) (string, error) {
	mp, err := leerMarketplace(dir)
	if err != nil {
		return "", err
	}
	for _, p := range mp.Plugins {
		if p.Name == arnesID {
			return filepath.Join(dir, filepath.Clean(p.Source)), nil
		}
	}
	return "", fmt.Errorf("fscatalog: arnés %q no está en el catálogo", arnesID)
}

func leerMarketplace(dir string) (*marketplaceFile, error) {
	raw, err := os.ReadFile(filepath.Join(dir, ".claude-plugin", "marketplace.json"))
	if err != nil {
		return nil, fmt.Errorf("fscatalog: %s no es un marketplace (.claude-plugin/marketplace.json): %w", dir, err)
	}
	var mp marketplaceFile
	if err := json.Unmarshal(raw, &mp); err != nil {
		return nil, fmt.Errorf("fscatalog: marketplace.json ilegible: %w", err)
	}
	return &mp, nil
}

// LeerArnes arma el domain.Arnes de un dir forma-plugin: arnes.l0.json (manifiesto,
// puede faltar → degradación honesta) + plugin.json (nombre/descr/versión de fallback).
func (c *Catalog) LeerArnes(src string) (domain.Arnes, error) { return leerArnes(src) }

func leerArnes(src string) (domain.Arnes, error) {
	var a domain.Arnes

	var pf pluginFile
	if raw, err := os.ReadFile(filepath.Join(src, ".claude-plugin", "plugin.json")); err == nil {
		_ = json.Unmarshal(raw, &pf)
	}

	if raw, err := os.ReadFile(filepath.Join(src, "arnes.l0.json")); err == nil {
		var l0 arnesL0
		if jerr := json.Unmarshal(raw, &l0); jerr != nil {
			return a, fmt.Errorf("arnes.l0.json ilegible: %w", jerr)
		}
		a = domain.Arnes{
			ID: l0.ID, Nombre: l0.Nombre, Descripcion: l0.Descripcion, Rol: l0.Rol,
			Proceso: l0.Proceso, Canal: l0.Canal, Fases: l0.Fases,
		}
	}

	// Cadena de fallback canónica (bendecida por ArnesIA, interop 2026-07-07 / HS-12):
	// nombre → plugin.json name → id; ídem descripcion. arnes.l0.nombre es el campo
	// canónico para pintar.
	if a.ID == "" {
		a.ID = pf.Name
	}
	if a.Nombre == "" {
		a.Nombre = pf.Name
	}
	if a.Nombre == "" {
		a.Nombre = a.ID
	}
	if a.Rol == "" {
		a.Rol = a.Nombre // degradación: sin manifiesto, el rol es el nombre del plugin
	}
	if a.Descripcion == "" {
		a.Descripcion = pf.Description
	}
	a.Version = pf.Version
	if a.ID == "" {
		return a, fmt.Errorf("sin arnes.l0.json ni plugin.json — no es un arnés (nomenclatura-arnes v1)")
	}
	return a, nil
}
