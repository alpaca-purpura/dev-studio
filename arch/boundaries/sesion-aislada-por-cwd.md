---
regla: sesion-aislada-por-cwd
version: 1.0
updated: 2026-07-06
status: proposed
ledger: DH-13
sources:
  - url: https://code.claude.com/docs/en/headless
    autoridad: oficial
    revisado: 2026-07-06
severity: critical
---

## L1 · Principio (estándar de industria)

**Cada sesión de un agente con acceso a filesystem debe correr confinada a su propio directorio
de trabajo**, nunca a un cwd compartido entre sesiones ni a rutas sensibles del sistema
(`$HOME`, `~/.ssh`, `~/.claude`) — el blast radius de un turno que se desvía debe limitarse al
proyecto que el usuario abrió en esa pestaña.

## L2 · Realización (este árbol Go) — ⚠ divergencia parcial, honesta

- `internal/domain.Session.Cwd` es un campo por sesión (no hay cwd global); `Conductor.Spawn`
  hace `cmd.Dir = opts.Cwd` (`conductor.go`) — cada proceso `claude` arranca en su propio
  directorio. ⇐ L1 (confinamiento básico).
- **Lo que NO está implementado todavía** (a diferencia del proyecto hermano, que sí lo tiene):
  no hay validación de que `Cwd` sea una ruta absoluta, exista, o no contenga/esté-contenida-en
  rutas protegidas (`$HOME`, `~/.ssh`, `~/.gnupg`, `~/.claude`). Hoy `createSession`
  (`transport/http/sessions.go`) solo hace `filepath.Clean` y, si `cwd` viene vacío, usa
  `os.UserHomeDir()` — que es exactamente una de las rutas que debería estar prohibida.
- Sin token de capacidad ni allowlist de Host/Origin en la API local (`router.go` solo compara
  `r.Host` contra `localhost`/`127.0.0.1`) — más liviano que el equivalente del hermano.

## Checklist evaluable

| id | qué chequea | severidad | señal en el mapa | enforcer |
|----|-------------|-----------|-------------------|----------|
| cwd-por-sesion | cada sesión tiene su propio `Cwd`, ninguna comparte working dir | error | «un turno de la sesión A escribe en el proyecto de la sesión B» | ✅ ya cumplido (campo por sesión) |
| cwd-no-es-ruta-protegida | rechazar `$HOME`/`~/.ssh`/`~/.claude`/`/` como cwd de una sesión | critical | «un turno mal encausado lee/escribe credenciales o config de Claude Code» | ❌ TBD — sin test, sin validación en `createSession` |
| host-origin-allowlist con token | la API local exige token de capacidad, no solo Host allowlist | high | «otra pestaña del navegador o proceso local llama la API sin permiso» | ❌ TBD |

## Changelog

- 2026-07-06 · v1.0 · Nodo fundacional (DH-13) — nace `proposed` y **deliberadamente incompleto**:
  documenta tanto lo que F1 ya resuelve (aislamiento básico por cwd) como la brecha real frente al
  proyecto hermano (sin registry de paths protegidos, sin token). Regla de honestidad (METODOLOGIA
  heredada): mejor un check marcado `❌ TBD` que un `enforced` que miente.
