---
regla: conductor-encapsula-stream-json
version: 1.0
updated: 2026-07-06
status: proposed
ledger: DH-13
sources:
  - url: https://code.claude.com/docs/en/headless
    autoridad: oficial
    revisado: 2026-07-06
  - url: https://martinfowler.com/eaaDev/EventSourcing.html
    autoridad: experto
    revisado: 2026-07-06
enforced_by:
  - fitness/arch_test.go:TestConductorIsOnlyStreamJSONParser  # TBD — hoy solo disciplina de code review
severity: high
---

## L1 · Principio (estándar de industria)

**Un solo parser por protocolo externo.** El formato de frames de `claude --output-format
stream-json` (system/init, stream_event/content_block_delta, assistant, result…) es un detalle
de integración inestable entre versiones de la CLI. Si dos paquetes lo parsean cada uno a su
manera, una actualización de `claude` rompe en dos lugares con dos síntomas distintos.

## L2 · Realización (este árbol Go)

- Solo `internal/adapters/agent/claudecode/conductor.go` (`translate`, `rawFrame`) conoce los
  campos crudos del protocolo. Todo lo que sale de ahí es `ports.AgentEvent` (4 kinds: init,
  delta, result, error) — un evento normalizado, no el JSON crudo. ⇐ L1.
- El caso de uso (`session_service.go`) y el transporte (`transport/http`, `transport/sse`) solo
  ven `AgentEvent`/`dockFrame` — nunca un campo como `stop_reason` o `content_block_delta`.
- Verificado empíricamente contra la CLI real (`claude -p --input-format stream-json
  --output-format stream-json --include-partial-messages --verbose`) antes de escribir
  `translate()` — no se adivinó el schema.

## Checklist evaluable

| id | qué chequea | severidad | señal en el mapa | enforcer |
|----|-------------|-----------|-------------------|----------|
| un-solo-parser | ningún paquete fuera de `claudecode/` decodea un frame `rawFrame`/campos crudos del protocolo | error | «dos parsers del mismo protocolo divergen en la próxima versión de `claude`» | pendiente automatizar (hoy: code review) |
| eventos-normalizados-cruzan-el-puerto | `AgentEvent` nunca expone el JSON crudo, solo `Kind/Text/ClaudeSessionID/Model/Err` | warn | «un campo crudo se filtra al caso de uso, acopla la app a la versión de la CLI» | pendiente automatizar |

## Changelog

- 2026-07-06 · v1.0 · Nodo fundacional (DH-13) — nace `proposed`: la regla ya se respeta en el
  código (un solo paquete parsea), pero no hay todavía un test que la haga fallar en CI si alguien
  la rompe. Candidato: `go-arch-lint` o un `arch_test.go` que falle si `internal/usecase` o
  `internal/adapters/transport/*` importan tipos internos de `claudecode` más allá de `New`.
