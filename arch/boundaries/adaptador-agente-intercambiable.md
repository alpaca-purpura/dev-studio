---
regla: adaptador-agente-intercambiable
version: 1.0
updated: 2026-07-06
status: enforced
ledger: DH-13
sources:
  - url: https://alistair.cockburn.us/hexagonal-architecture/
    autoridad: estándar
    revisado: 2026-07-06
  - url: https://code.claude.com/docs/en/headless
    autoridad: oficial
    revisado: 2026-07-06
enforced_by:
  - internal/adapters/agent/claudecode/conductor.go (aserción `var _ ports.AgentPort = (*Conductor)(nil)`)
severity: high
---

## L1 · Principio (estándar de industria)

**Ports & adapters (hexagonal).** El caso de uso depende de una interfaz (`AgentPort`), nunca
de un binario concreto. Hoy el único agente es Claude Code CLI (driver CLI-nativo, BYO licencia
— decisión heredada DH-10 de la incubadora); mañana podría ser otro CLI o el Agent SDK sin tocar
`usecase/`.

## L2 · Realización (este árbol Go+React)

- `internal/ports/agent.go` declara `AgentPort` (`Spawn`) y `AgentSession` (`Send`/`Events`/`Close`)
  — el único contrato que `internal/usecase/session_service.go` conoce.
- `internal/adapters/agent/claudecode/conductor.go` es el ÚNICO paquete que importa `os/exec` con
  el binario `claude`. La composition root (`cmd/dev-studio/main.go`) es el único lugar que
  instancia `claudecode.New(...)` — nada más en el árbol conoce el tipo concreto. ⇐ L1.
- Verificación mínima hoy: aserción de compilador (`var _ ports.AgentPort = (*Conductor)(nil)`,
  `conductor.go:24`) — si `Conductor` deja de satisfacer el puerto, el build rompe. No hay
  segundo adaptador todavía que pruebe la intercambiabilidad real (queda en debates abiertos).

## Checklist evaluable

| id | qué chequea | severidad | señal en el mapa | enforcer |
|----|-------------|-----------|-------------------|----------|
| conductor-implementa-port | `Conductor` satisface `ports.AgentPort` en tiempo de compilación | error | «`usecase` no compila sin el adaptador concreto» | `go build ./...` (aserción en conductor.go) |
| usecase-no-importa-claudecode | `internal/usecase` no importa `internal/adapters/agent/claudecode` | error | «el caso de uso queda pegado a un solo agente» | pendiente: `go list -deps` (ver `dominio-independiente-de-transporte` como plantilla) |

## Changelog

- 2026-07-06 · v1.0 · Nodo fundacional (DH-13) — F1 esqueleto de la app. Nace `enforced` por la
  aserción de compilador; el segundo check (usecase no importa el adaptador) queda pendiente de
  automatizar, hoy es solo revisión manual del import graph.
