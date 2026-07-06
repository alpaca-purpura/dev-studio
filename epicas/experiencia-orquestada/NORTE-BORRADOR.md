# Épica «Experiencia Orquestada» — BORRADOR del norte (F0 pendiente de firma)

> ⚠ NADA de este documento está firmado. F0 = debatir y firmar el norte EN una sesión de este
> repo. Carpeta temporal por diseño (patrón Torre de Control): al cierre de la épica se BORRA;
> lo permanente se promueve a `specs/` y al ledger.

## Cruda del operador (2026-07-04, verbatim)

"Quiero crear una Épica para trabajar la interfaz de esta aplicación e ir trayendo lo que
hemos desarrollado de a pocos conforme a cómo hemos decidido será la experiencia del
developer/usuario. El objetivo de esta aplicación es que el usuario (puede ser CTO, developer,
devops, product owner) trabaje en esta aplicación dentro de una organización de forma
'orquestada', con otros usuarios complementarios trabajando a la par uno o varios sistemas a
la par, apalancándose de un repositorio (Github) como conector entre todos."

## Qué gobierna esta épica

1. **La experiencia manda el port.** Ninguna pieza del monorepo entra "porque existe": entra
   cuando el journey del rol la pide, en la forma que el journey pide.
2. **Inventario de herencia disponible** (monorepo congelado): torre de control 2 ejes ·
   descriptor de proceso + board gobernado (ProcesoProvider, 13 verbos CDEvents) · cockpit de
   delivery (sesiones de agente parametrizadas) · lente arquitectura (`meta.clase`) ·
   multi-workspace · 27 endpoints Go · UI Next v0.6.x.
3. **Decisiones vigentes que la épica NO reabre** (salvo ficha): driver CLI-nativo (DH-10) ·
   descriptor I-77 · Go + UI embebida.

## Candidatos de fases (a firmar/reformar en F0)

- **F0 · Norte firmado**: journey por rol (CTO · developer · devops · PO) + principios de la
  experiencia + modelo de "organización" y del conector GitHub (¿qué sincroniza el repo?) +
  auditoría de herencia (orden del port). Candidato de método: skill `service-design-doing`
  (blueprint frontstage/backstage + inventario de interfaces).
- **F1 · Esqueleto de la app**: shell de escritorio (binario `dev-studio` nuevo, limpio) + driver
  CLI-nativo (la story que DH-10 dejó lista) + primera pantalla del journey.
- **F2+ · Port por rebanadas**: cada rebanada = una pieza de la herencia re-entrando por la
  puerta de la experiencia (spec congelada por fase, checkpoint de forks, verificación en
  browser/app real).
- **FN · Instalable v1**: cross-compile Win/Linux/mac + firma + updater (la deuda que DH-09
  dejó dimensionada).

## Reglas de la casa que aplican

SPEC congelada por fase en `specs/` · forks por AskUserQuestion (checkpoint 60s → recomendada)
· fichas DH-NN por fase · verificación real (browser/app) antes de cerrar fase · una
conversación = este repo (no arrastrar contexto del monorepo: lo que haga falta, se cita por
ficha/spec).
