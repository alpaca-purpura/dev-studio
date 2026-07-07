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
| CAP-01 | **Shell de escritorio** — binario Go único (`cmd/dev-studio`) con SPA React 19 embebida (`go:embed`) | funcional | DH-13 (2026-07-06): app levantada en vivo, UI servida desde el binario, navegada en navegador real | instalable (firma/updater) no existe — PB-24 |
| CAP-02 | **Sesiones Claude Code multisesión** — driver CLI-nativo propio: un subproceso `claude` por sesión, stream-json stdin/stdout, un turno a la vez | funcional ⚠ | DH-13: 2 sesiones concurrentes reales, procesos `claude` independientes con `claude_session_id` distinto, turnos respondidos sin cruce (`ALPHA`/`BETA`), cierre limpia el proceso | sin validación de rutas protegidas ni aislamiento entre sesiones — `arch/boundaries/sesion-aislada-por-cwd.md` → PB-02. No exponer fuera del equipo hasta cerrarla |
| CAP-03 | **Rail de sesiones** — colapsable 224px↔52px, tabs estilo WARP, crear/cerrar/navegar sesiones | funcional | DH-13: navegación entre tabs preservando conversación, colapso y cierre ejercidos en navegador (Chrome DevTools MCP) | — |
| CAP-04 | **API REST + SSE** — `/api/sessions` + `/events` (evento `dock` multiplexado) | funcional | DH-13: la verificación de CAP-02/03 corrió sobre esta API en vivo | sin replay `Last-Event-ID` — PB-18 |
| CAP-05 | **Persistencia de sesiones** — `~/.dev-studio/sessions.json`, escritura atómica | funcional | DH-13: sesiones sobreviven en el archivo; ejercida como parte del flujo crear/cerrar | sin SQLite (decisión deliberada a este tamaño, no brecha) |
| CAP-06 | **Tokens de diseño** — `web/src/app/styles/theme.css` (Tailwind v4, dark por `data-theme`), copiados 1:1 de harness-studio | funcional | DH-13: estilo verificado visualmente en navegador (rail 1:1 con el hermano) | solo tokens: sin Storybook ni catálogo de átomos propio — PB-03 |

## Log de entregas

| Fecha | Qué entró | PB | Ficha |
|---|---|---|---|
| 2026-07-06 | F1 esqueleto: CAP-01..CAP-06 (fundación, previa a este registro) | — (pre-backlog) | DH-13 |
