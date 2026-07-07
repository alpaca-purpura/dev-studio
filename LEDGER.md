# Ledger — DevStudio (fichas DH-NN)

> Registro de decisiones de ESTE producto. Mismo formato/disciplina que la casa prenter-harness.
> **Continuidad:** DH-01..DH-11 viven en `prenter-harness/products/devhub/LEDGER.md` (la
> incubadora, congelada en la graduación). Este repo arranca en **DH-12**.

## Fichas

### DH-12 · Fundación del repo propio — graduación de P2 con visión ampliada — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-04):* "quiero crear una Épica para trabajar la interfaz de esta
aplicación e ir trayendo lo que hemos desarrollado de a pocos conforme a cómo hemos decidido
será la experiencia del developer/usuario. El objetivo es que el usuario (puede ser CTO,
developer, devops, product owner) trabaje en esta aplicación dentro de una organización de
forma 'orquestada', con otros usuarios complementarios trabajando a la par uno o varios
sistemas a la par, apalancándose de un repositorio (GitHub) como conector entre todos. Esto ya
creció más de lo que buscábamos y será un producto propio con su propio repositorio… una
aplicación cuyo objetivo es construir y mantener software basado en proceso y arquitectura,
donde usamos arquitectura as code, proceso as code, software documentation as code, etc. Este
nuevo repositorio debe nacer limpio, con la nueva visión, y como arneses de construcción vamos
a usar los del KIT DEV (plugin), pero lo iremos evolucionando conforme avanzamos." Forks
firmados por AskUserQuestion: nombre **DevStudio** (repo `dev-studio`, binario `dev-studio` — cierra de
nacimiento la deuda del rename DH-01) · **graduación** con célula del monorepo = fuente
CONGELADA read-only del port gradual (mecánica I-69 adelantada deliberadamente: sin clientes
aún) · repo GitHub **privado en alpacapurpura** desde el día 0.

*Desarrollo:* repo nace limpio con VISION.md (visión ampliada: construir y mantener software
basado en proceso y arquitectura; as-code en las 3 superficies; trabajo orquestado
multi-usuario/multi-rol; GitHub = conector) + este ledger + épica «Experiencia Orquestada»
(borrador de norte, F0 pendiente de firma EN este repo) + kit dev instalado como plugin desde
el marketplace (`alpacapurpura/prenter-marketplace`). Decisiones técnicas heredadas que siguen
vigentes: DH-10 (conexión CC = driver CLI-nativo, BYO licencia, SIN API) · descriptor de
proceso I-77 (contrato de ecosistema — el kit lo shipea, DevStudio lo interpreta) · binario Go +
UI embebida. El código go+ui NO se copió: el port es gradual, gobernado por la épica, pieza
por pieza según la experiencia decidida.

*Conecta:* DH-01..DH-11 (la historia en la incubadora; DH-11 = la ficha espejo de esta
graduación en el monorepo) · I-NN de ecosistema (graduación — registrada en
`prenter-harness/tooling/strategy/LEDGER.md`) · I-77 (descriptor) · DH-10 (driver CLI-nativo)
· KIT-06 (marketplace del que este repo consume su arnés).

*Siguiente:* F0 de la épica «Experiencia Orquestada» — norte firmado (journey por rol +
principios de la experiencia) + auditoría de herencia (qué pieza del monorepo entra primero y
en qué forma).

### DH-13 · F1 esqueleto de la app — driver CLI-nativo + rail multisesión — `decidida` · `vig:vigente`

*Cruda (operador, 2026-07-06):* "revisa todos los aspectos técnicos de harness-studio y hagamos
un hermano de la arquitectura para construir la aplicación de escritorio que se conecta con
claude code, manejando múltiples sesiones claude code entre las que puedo navegar al cambiar de
tab de mi aplicación y crear o cerrar tmb desde mi aplicación. Copia también las buenas prácticas
de creación de visión, arquitectura as code, etc. Copia el estilo gráfico del storybook a la
fecha y cómo se ha manejado visualmente las pestañas de sesión como columna colapsable a la
izquierda. Todo lo demás dejalo en blanco. Revisa la arquitectura as code de lo correspondiente
y crea una para nuestro proyecto adaptándola. Levanta la aplicación y verifica que esté
funcionando."

*Desarrollo:* auditoría técnica de `harness-studio` (hermano de arquitectura, mismo operador) en
tres ejes: (1) `arch/` + VISION/LEDGER/METODOLOGIA/CLAUDE.md como patrón documental — L1↔L2 +
checklist evaluable por boundary, fichas `Cruda/Desarrollo/Conecta/Siguiente`; (2) el driver
CLI-nativo: `internal/adapters/agent/claudecode` (subproceso `claude -p --input-format
stream-json --output-format stream-json`, un proceso por sesión, `SessionService` en memoria +
persistencia JSON liviana, SSE multiplexado); (3) el estilo gráfico Storybook (tokens Tailwind v4
`theme.css`, paleta `--primary #a8742c`/dark `#d9a35b`) y el widget `SessionRail` (rail
izquierdo colapsable 224px↔52px, tabs estilo WARP, acento `shadow-[inset_3px_0_0_var(--primary)]`
en la activa). Se construyó el **hermano real** en este repo (no fork, reimplementación propia
contra el protocolo verificado con una corrida real de `claude`): módulo Go
`github.com/alpacapurpura/dev-studio`, arquitectura hexagonal (`domain/ports/usecase/adapters`),
API REST + SSE (`/api/sessions`, `/events`), SPA React 19 + Zustand embebida vía `go:embed`,
`SessionRail` copiado 1:1 en estilo y simplificado en dominio (sin `arnes`/`salud`/`view` — eso
es de ArnesIA). Arquitectura as code propia en `arch/` (5 boundaries, 2 con test real corriendo:
`TestOneTurnAtATime`, `TestDomainNoTransportImport`) — deliberadamente más chica que la del
hermano y con una brecha de seguridad documentada explícitamente (`sesion-aislada-por-cwd`: sin
validación de rutas protegidas todavía). Verificación end-to-end real (no simulada): 2 sesiones
creadas desde el navegador (Chrome DevTools MCP), cada una spawneando su propio proceso `claude`
con `claude_session_id` distinto, turnos concurrentes respondidos correctamente y sin
cruzarse (`ALPHA`/`BETA`), navegación entre tabs preservando conversación, colapso del rail,
cierre de sesión limpiando el proceso.

*Conecta:* DH-10 (driver CLI-nativo, decisión heredada que esta ficha materializa por primera
vez en código propio) · DH-12 (fundación del repo) · épica
[`epicas/experiencia-orquestada/NORTE-BORRADOR.md`](./epicas/experiencia-orquestada/NORTE-BORRADOR.md)
(candidato de fase F1 «esqueleto de la app») · `harness-studio` (proyecto hermano, misma
metodología, fuente del patrón `arch/` y del estilo `SessionRail`) · `arch/INDEX.md` (el árbol
nuevo) · `arch/boundaries/sesion-aislada-por-cwd.md` (brecha de seguridad declarada, no oculta).

*Siguiente:* **DH-14 = cerrar la brecha de `sesion-aislada-por-cwd`** (validación de rutas
protegidas antes de exponer la app fuera de este equipo) + F0 de la épica (norte firmado,
todavía pendiente pese a que F1 ya entregó código) + decidir si el segundo adaptador de agente
(prueba real de intercambiabilidad) entra antes o después del port por rebanadas.

<!-- Próximas: DH-14, DH-15, … -->

## Log

| Fecha | Decisión | Fichas |
|---|---|---|
| 2026-07-04 | Fundación del repo propio: graduación de P2 (nombre DevStudio confirmado — binario `dev-studio`, muere la colisión `cockpit` de nacimiento); visión ampliada = construir y mantener software basado en proceso y arquitectura, trabajo orquestado multi-usuario (CTO·developer·devops·PO), GitHub conector; célula del monorepo congelada como fuente del port gradual; kit dev como plugin del marketplace; épica «Experiencia Orquestada» sembrada (F0 por firmar). | DH-12 |
| 2026-07-06 | F1 esqueleto de la app entregado: driver CLI-nativo propio (subproceso `claude` + stream-json, multisesión real verificada con 2 procesos concurrentes) + rail de sesiones colapsable (estilo Storybook de `harness-studio` copiado 1:1, dominio simplificado) + arquitectura as code propia (`arch/`, 5 boundaries, 2 enforced con test) con brecha de seguridad documentada sin ocultar. | DH-13 |
