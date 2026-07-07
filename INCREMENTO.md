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
| CAP-02 | **Sesiones Claude Code multisesión** — driver CLI-nativo propio: un subproceso `claude` por sesión, stream-json stdin/stdout, un turno a la vez + cola de mensajes (Encolar, client-side, RN-8) | funcional ⚠ | DH-13 (2 sesiones ALPHA/BETA) + DH-15: re-verificado con 2 sesiones en repos distintos, turnos concurrentes `ALFA-1`/`BETA-2` respondidos sin cruce (AC-6b) | sin validación de rutas protegidas ni aislamiento entre sesiones — `arch/boundaries/sesion-aislada-por-cwd.md` → PB-02 (la siguiente). No exponer fuera del equipo hasta cerrarla |
| CAP-03 | **Rail Repositorios→Workspaces** — workspace = sesión (1:1, RN-2): repos colapsables, workspace-item con branch + diff ±N + status-dot reales, atajos ⌘1-9, sesiones F1 sin repo bajo «(sin repositorio)» | funcional | DH-15: AC-10-ui (migración visible), branch/±N leídos del git real de cada repo, navegación por atajos y clic ejercida en vivo. El rail F1 (tabs WARP) murió reemplazado | dots `conflict`/`archived` + dropdown ⋯ llegan con PB-02 |
| CAP-04 | **API REST + SSE** — `/api/sessions` + `/api/repos` + `/api/sessions/{id}/git/*` + `/events` (evento `dock` multiplexado) | funcional | DH-13 + DH-15: todo el gate AC corrió sobre esta API en vivo (registro de repos, status/diff/log/commit git) | sin replay `Last-Event-ID` — PB-18 |
| CAP-05 | **Persistencia de estado** — `~/.dev-studio/state.json` (repos + sesiones), escritura atómica, migra `sessions.json` F1 solo y sin pérdida | funcional | DH-15: migración real ejercida (las 2 sesiones F1 del archivo legacy aparecieron bajo «(sin repositorio)» y sobreviven en state.json: 2 repos + 4 sesiones) | sin SQLite (decisión deliberada a este tamaño, no brecha) |
| CAP-06 | **Design system PRENTER** — tokens de marca (teal único acento, dark-first + light) en `theme.css`, fuentes Jost/Mulish/JetBrains Mono VENDORIZADAS en el binario, Storybook con ~20 átomos + organismos del shell (RN-9) | funcional | DH-15: `npm run build-storybook` verde; UI completa renderizando PRENTER en navegador (screenshots); grep `#a8742c` = 0 (AC-9); SSoT sincronizado de Claude Design `a98c2e0d` a `specs/shell/tokens/` | fuentes reales de marca (Coco Gothic/Sansation) pendientes — hoy sustitutos declarados |
| CAP-07 | **Registro de repositorios** — alta por ruta local (valida `.git`), idempotente, eliminación; contenedor de workspaces | funcional | DH-15: AC-1 en vivo — ruta sin `.git` rechazada con error visible; demo-alfa/demo-beta registrados desde la UI | clonar desde URL no existe (rechazo honesto) — idea nueva si se quiere |
| CAP-08 | **Panel Cambios git REAL** — status + quick-look unificado + revisión side-by-side (N/M) + historial (`git log`) + **commit solo por selección explícita** (pathspec, RN-5); push/pull/fetch NO EXISTEN en el adapter (boundary `git-solo-lectura-y-commit`, fitness test) | funcional | DH-15: AC-4 en vivo — 3 archivos reales listados, diff real, commit de 2/3 tildados verificado por `git show --name-only` (el destildado quedó fuera), historial refleja el commit; AC-5 botones sync deshabilitados + fitness verde | filtro por rol necesita PB-02; «adjuntar conversación» → PB-10; «generar mensaje con IA» stub |
| CAP-09 | **Sesión ligada a historia** (RN-1) — flujo picker: `+ Nuevo Workspace` → Backlog picker → historia → Config con paquete → `Crear sesión aislada`; sin paquete no hay sesión (hint + redirección) | funcional ⚠ | DH-15: AC-2/AC-3 en vivo — sesión creada ligada (chip 📋 en header), turno real `SHELL-OK` respondido por `claude`; Config directo sin paquete redirige al picker | historias = datos de ejemplo hasta PB-07; «sesión aislada» = cwd raíz del repo hasta PB-02 (el nombre ya es el contrato) |

## Log de entregas

| Fecha | Qué entró | PB | Ficha |
|---|---|---|---|
| 2026-07-06 | F1 esqueleto: CAP-01..CAP-06 (fundación, previa a este registro) | — (pre-backlog) | DH-13 |
| 2026-07-07 | Shell PRENTER: CAP-01/03/04/05/06 actualizadas + CAP-07/08/09 nuevas — 24/24 checks del gate en vivo | PB-04 (⊕ PB-03) | DH-15 |
