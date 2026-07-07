---
regla: sesion-aislada-por-cwd
version: 2.0
updated: 2026-07-07
status: enforced
ledger: DH-16
sources:
  - url: https://code.claude.com/docs/en/headless
    autoridad: oficial
    revisado: 2026-07-06
  - url: https://git-scm.com/docs/git-worktree
    autoridad: oficial
    revisado: 2026-07-07
severity: critical
---

## L1 · Principio (estándar de industria)

**Cada sesión de un agente con acceso a filesystem debe correr confinada a su propio directorio
de trabajo**, nunca a un cwd compartido entre sesiones ni a rutas sensibles del sistema
(`$HOME`, `~/.ssh`, `~/.claude`) — el blast radius de un turno que se desvía debe limitarse al
workspace que el usuario abrió en esa pestaña.

## L2 · Realización (este árbol Go) — cerrada en DH-16 (PB-02)

- **Workspace aislado por sesión (patrón Conductor):** `createSession` con `repo_id` crea un
  worktree propio (`ports.GitWorkspace` → `adapters/git/cli.CreateWorktree`: `git worktree add
  ~/.dev-studio/workspaces/{repo}/{slug} -b wt/{slug}`) y la sesión nace con `Cwd` =
  `Workspace` = ese checkout. Dos sesiones del mismo repo jamás comparten working dir; tocar
  el worktree A no aparece en el status de B (test del adapter).
- **Rutas protegidas:** `domain.RutaProtegida(home, ruta)` (pura, testeada) veda `$HOME`
  exacto, `~/.ssh`, `~/.gnupg`, `~/.aws`, `~/.kube`, `~/.docker`, `~/.dev-studio` y las raíces
  de sistema (`/`, `/etc`, `/usr`, …). La aplican `RepoService.Register` y la puerta legacy de
  `createSession` (cwd custom) — `ErrRutaProtegida` → 422 con la ruta nombrada (RN-2).
- **Cierre sin pérdida (RN-3):** `DELETE /api/sessions/{id}?workspace=remove` usa
  `RemoveWorktree` SIN `--force` — un worktree sucio se conserva y la respuesta lo dice.
- Sin token de capacidad ni allowlist de Origin en la API local todavía (`router.go` compara
  `r.Host` contra localhost) — check restante, abajo.

## Checklist evaluable

| id | qué chequea | severidad | señal en el mapa | enforcer |
|----|-------------|-----------|-------------------|----------|
| cwd-por-sesion | cada sesión tiene su propio `Cwd`; vía repo = worktree propio | error | «un turno de la sesión A escribe en el proyecto de la sesión B» | `cli/worktree_test.go:TestCreateWorktreeAislado` |
| cwd-no-es-ruta-protegida | rechazar `$HOME`/`~/.ssh`/raíces de sistema como repo/cwd | critical | «un turno mal encausado lee/escribe credenciales» | `arch_test.go:TestRutasProtegidasRechazadas` + `domain/rutas_test.go` |
| remove-sin-force | borrar workspace jamás fuerza un worktree sucio | error | «cerrar una sesión tiró trabajo sin commit» | `cli/worktree_test.go:TestRemoveWorktree` |
| host-origin-allowlist con token | la API local exige token de capacidad, no solo Host allowlist | high | «otra pestaña/proceso local llama la API sin permiso» | ❌ TBD — único check restante |

## Changelog

- 2026-07-07 · v2.0 · **`proposed` → `enforced`** (DH-16, PB-02): worktree por sesión +
  validación de rutas protegidas + remove sin force, los tres con test corriendo. Queda UN
  check TBD (token de capacidad) — documentado, no oculto.
- 2026-07-06 · v1.0 · Nodo fundacional (DH-13) — nace `proposed` y deliberadamente incompleto.
