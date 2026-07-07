---
regla: git-solo-lectura-y-commit
version: 1.0
updated: 2026-07-07
status: enforced
ledger: DH-15
sources:
  - url: https://git-scm.com/docs/git#_high_level_commands_porcelain
    autoridad: oficial
    revisado: 2026-07-07
enforced_by:
  - fitness/arch_test.go:TestGitAdapterSinVerbosProhibidos
  - fitness/arch_test.go:TestGitCommitExigePathspec
severity: high
---

## L1 · Principio (estándar de industria)

**Least authority sobre el estado del usuario.** Una herramienta que opera SOBRE el repositorio
de su usuario no debe poder alterar su relación con el remoto ni reescribir su historia: el
blast radius de un bug pasa de «un commit local de más» a «trabajo ajeno perdido / remoto
contaminado». La app limita su autoridad git al mínimo que la experiencia pide: leer estado y
registrar commits locales de archivos explícitamente elegidos.

## L2 · Realización (este árbol Go)

- Los puertos (`ports/git.go`) SON el contrato completo: `GitInfo` (Status/DiffFile/Log) +
  `GitCommit` (Commit con `paths []string` obligatorio). No existe puerto para push/pull/
  fetch/reset/rebase — un caso de uso no puede pedir lo que el contrato no ofrece. ⇐ L1.
- El adaptador (`adapters/git/cli`) ejecuta el `git` del usuario como subproceso (espíritu BYO
  de DH-10) y solo con los verbos: `branch` · `status` · `diff` · `show` · `log` · `add` ·
  `commit` · `rev-parse` · `init` (tests). El fitness test escanea el CÓDIGO del paquete
  (comentarios excluidos): los strings `push`/`pull`/`fetch`/`reset`/`rebase` no aparecen.
- `Commit(cwd, paths, mensaje)` con `len(paths)==0` → `ErrSinPaths` (RN-5 de la spec shell:
  commitear «todo» sin enumerarlo es imposible por diseño; el transporte además responde 400).
- La UI muestra los botones fetch/pull/push DESHABILITADOS con tooltip honesto (RN-7) — la
  decisión de ofrecerlos es de una rebanada futura (PR-flow, PB-14), no un hueco.

## Checklist evaluable

| id | qué chequea | severidad | señal en el mapa | enforcer |
|----|-------------|-----------|-------------------|----------|
| sin-verbos-prohibidos | el fuente de `adapters/git/cli` no contiene push/pull/fetch/reset/rebase | error | «la app tocó el remoto o reescribió historia del usuario» | `arch_test.go:TestGitAdapterSinVerbosProhibidos` |
| commit-exige-pathspec | `Commit` con paths vacío devuelve error | error | «un commit silencioso barrió archivos no elegidos» | `arch_test.go:TestGitCommitExigePathspec` |

## Changelog

- 2026-07-07 · v1.0 · Nace `enforced` con la rebanada shell (DH-15, spec shell §5/RN-4/RN-5).
