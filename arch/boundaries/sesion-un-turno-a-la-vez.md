---
regla: sesion-un-turno-a-la-vez
version: 1.0
updated: 2026-07-06
status: enforced
ledger: DH-13
sources:
  - url: https://go.dev/ref/mem
    autoridad: oficial
    revisado: 2026-07-06
  - url: https://html.spec.whatwg.org/multipage/server-sent-events.html
    autoridad: oficial
    revisado: 2026-07-06
enforced_by:
  - fitness/arch_test.go:TestOneTurnAtATime
severity: high
---

## L1 · Principio (estándar de industria)

**Un solo productor por recurso serializado.** El stdin de un proceso `claude` es un recurso con
un único escritor válido a la vez: dos turnos concurrentes escribiendo frames NDJSON al mismo
stdin los intercalan y corrompen el mensaje que el proceso intenta parsear. El patrón correcto es
rechazar (no encolar en silencio) un segundo turno mientras el primero sigue en vuelo.

## L2 · Realización (este árbol Go)

- `usecase.SessionService.Turn` (`session_service.go`) chequea `r.meta.Status ==
  StatusStreaming` bajo el mismo `sync.Mutex` que protege el registro — si es así, devuelve
  `ErrBusy` sin tocar el conductor. ⇐ L1.
- El transporte HTTP (`transport/http/sessions.go`, `sessionTurn`) mapea `ErrBusy` a **409**; el
  store del frontend (`shared/store/sessions-store.ts`) además deshabilita el input mientras
  `status === "streaming"` (guard local, espeja el 409 del server).
- `ccSession.Send` (`conductor.go`) además serializa la escritura real a stdin con su propio
  `sync.Mutex` (`sendMu`) — doble cinturón: el use case nunca debería intentar un segundo `Send`
  concurrente, pero si algo cambia, el adaptador tampoco permite que se intercalen.

## Checklist evaluable

| id | qué chequea | severidad | señal en el mapa | enforcer |
|----|-------------|-----------|-------------------|----------|
| turn-rechaza-si-streaming | `Turn()` devuelve `ErrBusy` si la sesión ya está streaming | error | «dos turnos concurrentes corrompen el stdin del proceso claude» | `arch_test.go:TestOneTurnAtATime` |
| stdin-serializado | `ccSession.Send` serializa bajo `sendMu` | error | «JSON intercalado en stdin, claude no puede parsear el frame» | revisión manual (no hay test de concurrencia real del pipe todavía) |

## Changelog

- 2026-07-06 · v1.0 · Nodo fundacional (DH-13) — nace `enforced`: `TestOneTurnAtATime` corre y
  pasa (`go test ./arch/fitness/...`). El segundo check (stdin serializado) queda sin test
  dedicado — cubierto solo por inspección de código.
