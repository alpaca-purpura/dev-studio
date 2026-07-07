# INCREMENTO de producto — DevStudio (capabilities funcionales)

> **Plano: PRODUCTO.** Registro de lo YA construido y verificado en vivo — el Product Increment.
> Par de [`BACKLOG.md`](./BACKLOG.md): una PB entregada se registra acá en el MISMO commit que la
> cierra (regla 6 del backlog).
>
> **Regla dura:** nada se lista sin verificación REAL citada (acción ejercida + efecto observado —
> nunca "compila" ni "GET 200"). Las brechas conocidas se declaran en su fila, no se ocultan.

## Capabilities

| ID | Capability | Estado | Evidencia de verificación | Brechas declaradas |
|---|---|---|---|---|
| CAP-01 | **Shell de escritorio PRENTER** — binario Go único (`cmd/dev-studio`) con SPA React 19 embebida (`go:embed`): rail Repositorios→Workspaces + studio-nav (Studio/Backlog/Producto/Roles-pronto/Config) + overlays que jamás tapan los rails (RN-3) + stubs honestos con destino (RN-7) | funcional | DH-15 (2026-07-07): 24/24 checks del gate AC-1..AC-7 ejercidos en navegador real (puppeteer + Chrome, screenshots en sesión); Esc vuelve a Studio; rails clickeables bajo overlay verificado por hit-test | instalable (firma/updater) no existe — PB-24 |
| CAP-02 | **Sesiones Claude Code multisesión** — driver CLI-nativo propio: un subproceso `claude` por sesión, stream-json stdin/stdout, un turno a la vez + cola de mensajes (Encolar, client-side, RN-8) | funcional | DH-13 (2 sesiones ALPHA/BETA) + DH-15 (concurrencia sin cruce) + **DH-16: brecha CERRADA** — boundary `sesion-aislada-por-cwd` enforced (worktree por sesión + rutas protegidas) | resta el token de capacidad de la API local (check TBD del boundary) |
| CAP-03 | **Rail Repositorios→Workspaces** — workspace = sesión (1:1, RN-2): repos colapsables, workspace-item con branch `wt/{slug}` + diff ±N + status-dot DEL WORKTREE DE CADA SESIÓN (incl. `conflict` por unmerged), ✕ de cierre, atajos ⌘1-9 | funcional | DH-15 (rail base) + DH-16: branch/±N por sesión verificados en vivo (wt/… visible en rail y header) | dot `archived` + dropdown ⋯ (duplicar/desde branch) siguen stub |
| CAP-04 | **API REST + SSE** — `/api/sessions` + `/api/repos` + `/api/sessions/{id}/git/*` + `/events` (evento `dock` multiplexado) | funcional | DH-13 + DH-15: todo el gate AC corrió sobre esta API en vivo (registro de repos, status/diff/log/commit git) | sin replay `Last-Event-ID` — PB-18 |
| CAP-05 | **Persistencia de estado** — `~/.dev-studio/state.json` (repos + sesiones), escritura atómica, migra `sessions.json` F1 solo y sin pérdida | funcional | DH-15: migración real ejercida (las 2 sesiones F1 del archivo legacy aparecieron bajo «(sin repositorio)» y sobreviven en state.json: 2 repos + 4 sesiones) | sin SQLite (decisión deliberada a este tamaño, no brecha) |
| CAP-06 | **Design system PRENTER** — tokens de marca (teal único acento, dark-first + light) en `theme.css`, fuentes Jost/Mulish/JetBrains Mono VENDORIZADAS en el binario, Storybook con ~20 átomos + organismos del shell (RN-9) | funcional | DH-15: `npm run build-storybook` verde; UI completa renderizando PRENTER en navegador (screenshots); grep `#a8742c` = 0 (AC-9); SSoT sincronizado de Claude Design `a98c2e0d` a `specs/shell/tokens/` | fuentes reales de marca (Coco Gothic/Sansation) pendientes — hoy sustitutos declarados |
| CAP-07 | **Registro de repositorios** — alta por ruta local (valida `.git`), idempotente, eliminación; contenedor de workspaces | funcional | DH-15: AC-1 en vivo — ruta sin `.git` rechazada con error visible; demo-alfa/demo-beta registrados desde la UI | clonar desde URL no existe (rechazo honesto) — idea nueva si se quiere |
| CAP-08 | **Panel Cambios git REAL** — status + quick-look unificado + revisión side-by-side (N/M) + historial (`git log`) + **commit solo por selección explícita** (pathspec, RN-5); push/pull/fetch NO EXISTEN en el adapter (boundary `git-solo-lectura-y-commit`, fitness test) | funcional | DH-15: AC-4 en vivo — 3 archivos reales listados, diff real, commit de 2/3 tildados verificado por `git show --name-only` (el destildado quedó fuera), historial refleja el commit; AC-5 botones sync deshabilitados + fitness verde | filtro por rol necesita PB-02; «adjuntar conversación» → PB-10; «generar mensaje con IA» stub |
| CAP-09 | **Sesión ligada a historia** (RN-1) — flujo picker: `+ Nuevo Workspace` → Backlog picker → historia → Config con paquete → `Crear sesión aislada`; sin paquete no hay sesión (hint + redirección) | funcional | DH-15: AC-2/AC-3 en vivo (turno real `SHELL-OK`) + DH-16: «aislada» ya es literal — la sesión nace en su worktree | historias = datos de ejemplo hasta PB-07 |
| CAP-10 | **Workspace aislado por sesión** (PB-02) — worktree + branch `wt/{slug}` en `~/.dev-studio/workspaces/{repo}/`, rutas protegidas vedadas en toda puerta, cierre con modal conservar/borrar (borrar sucio → git rechaza, jamás `--force`) — boundary `sesion-aislada-por-cwd` **enforced** | funcional | DH-16: 14/14 checks en vivo — worktree en disco + `git worktree list`, 2 sesiones del mismo repo sin cruce de Cambios, $HOME/etc rechazados (422 + error visible), borrar sucio conservado con el error real de git, borrar limpio desaparece | token de capacidad de la API local TBD (único check restante del boundary) |
| CAP-11 | **Instalable dogfooding + updater** (PB-26) — `install.sh` (binario `~/.local/bin` + lanzador app-mode + `.desktop` + icono PRENTER + versión SHA embebida) + «Actualizar» en la app: rebuild del repo local → restart (`syscall.Exec`) → la SPA recarga sola | funcional | DH-16: instalación real verificada (artefactos en disco) y **la app se actualizó a sí misma en el gate** — un fix commiteado solo en el working tree quedó vivo tras el restart (AC-6, verificado con el 422 post-update) | Linux-only; instalable comercial (cross-compile+firma+Consumer Terms) = PB-24 |

## Log de entregas

| Fecha | Qué entró | PB | Ficha |
|---|---|---|---|
| 2026-07-06 | F1 esqueleto: CAP-01..CAP-06 (fundación, previa a este registro) | — (pre-backlog) | DH-13 |
| 2026-07-07 | Shell PRENTER: CAP-01/03/04/05/06 actualizadas + CAP-07/08/09 nuevas — 24/24 checks del gate en vivo | PB-04 (⊕ PB-03) | DH-15 |
| 2026-07-07 | Workspace aislado + instalable: CAP-10/11 nuevas, CAP-02/03/09 actualizadas (brecha de seguridad fundacional CERRADA) — 14/14 checks contra el binario instalado, self-update incluido | PB-02 ⊕ PB-26 | DH-16 |
