# DevHub — construir y mantener software basado en proceso y arquitectura

Producto standalone (graduado del monorepo `prenter-harness`, 2026-07-04). Norte =
[`VISION.md`](./VISION.md) · registro = [`LEDGER.md`](./LEDGER.md) (fichas `DH-NN`, arranca en
DH-12; DH-01..DH-11 viven en la incubadora `prenter-harness/products/devhub/`, congelada).

**Qué es:** aplicación de escritorio multiplataforma donde una organización construye y
mantiene software con proceso, arquitectura y documentación **as code**. Usuarios (CTO ·
developer · devops · product owner) trabajan orquestados sobre uno o varios sistemas; el
repositorio GitHub es el conector.

**Decisiones técnicas vigentes:**
- Conexión con Claude Code = **driver CLI-nativo** (DH-10): spawn del `claude` del usuario vía
  stdin/stdout stream-json. BYO licencia — **sin API de Anthropic**; la app no toca credenciales.
- Binario Go + UI embebida (`go:embed`); instalador = el binario solo (Win/Linux/mac).
- Proceso as code = descriptor I-77 (categorías semánticas fijas · gates · dueños · verbos
  CDEvents); la consola deriva TODO del dato, cero hardcode del ciclo.

**Estado:** repo recién fundado — el código vive aún en el monorepo (congelado) y entra por
port gradual gobernado por la épica [`epicas/experiencia-orquestada/`](./epicas/experiencia-orquestada/)
(carpeta temporal por diseño; specs permanentes → `specs/`).

**Arnés de construcción:** kit dev (plugin del marketplace `alpacapurpura/prenter-marketplace`,
pineado por versión). Evoluciona con el producto — mejoras al arnés se upstreamean al kit
(backflow I-59), jamás fork silencioso.

**Git:** trunk-based — `main` única, commit/push directo, tags semver cuando haya releases.
