// Package web embebe el bundle compilado de la SPA (npm run build -> dist/) en el binario Go.
package web

import "embed"

//go:embed all:dist
var DistFS embed.FS
