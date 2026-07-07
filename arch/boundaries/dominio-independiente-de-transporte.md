---
regla: dominio-independiente-de-transporte
version: 1.0
updated: 2026-07-06
status: enforced
ledger: DH-13
sources:
  - url: https://alistair.cockburn.us/hexagonal-architecture/
    autoridad: estándar
    revisado: 2026-07-06
enforced_by:
  - fitness/arch_test.go:TestDomainNoTransportImport
severity: high
---

## L1 · Principio (estándar de industria)

**El dominio no conoce cómo llega al mundo.** `internal/domain` (el modelo: `Session`, `Turn`) y
`internal/usecase` (la orquestación: `SessionService`) no deben importar `net/http` ni ningún
detalle de transporte — si lo hacen, cambiar de REST a gRPC o agregar un segundo transporte
obliga a tocar la lógica de negocio.

## L2 · Realización (este árbol Go)

- `internal/domain/session.go` y `internal/usecase/session_service.go` solo importan la
  librería estándar (`context`, `sync`, `encoding/json`, …) y `internal/ports` — nunca
  `net/http`. ⇐ L1.
- El único sub-puerto que el caso de uso necesita del transporte realtime es la interfaz mínima
  `usecase.EventPublisher` (`Publish(eventType string, data []byte)`) — el broker SSE concreto
  (`adapters/transport/sse.Broker`) la satisface por duck-typing, sin que `usecase` importe
  `adapters/transport/sse`.
- `internal/adapters/transport/http` es el único paquete (junto al `main.go` de la composition
  root) que importa `net/http`.

## Checklist evaluable

| id | qué chequea | severidad | señal en el mapa | enforcer |
|----|-------------|-----------|-------------------|----------|
| domain-sin-net-http | `internal/domain` no depende de `net/http` (directa o transitivamente) | error | «el modelo de negocio se acopla a REST» | `arch_test.go:TestDomainNoTransportImport` |
| usecase-sin-net-http | `internal/usecase` no depende de `net/http` (directa o transitivamente) | error | «la orquestación se acopla a REST» | `arch_test.go:TestDomainNoTransportImport` |

## Changelog

- 2026-07-06 · v1.0 · Nodo fundacional (DH-13) — nace `enforced`: `TestDomainNoTransportImport`
  corre `go list -deps` sobre `internal/domain` y `internal/usecase` y falla si aparece
  `net/http` en el grafo de imports. Pasa hoy (`go test ./arch/fitness/...`).
