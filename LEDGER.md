# Ledger — DevHub (fichas DH-NN)

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
firmados por AskUserQuestion: nombre **DevHub** (repo `devhub`, binario `devhub` — cierra de
nacimiento la deuda del rename DH-01) · **graduación** con célula del monorepo = fuente
CONGELADA read-only del port gradual (mecánica I-69 adelantada deliberadamente: sin clientes
aún) · repo GitHub **privado en alpacapurpura** desde el día 0.

*Desarrollo:* repo nace limpio con VISION.md (visión ampliada: construir y mantener software
basado en proceso y arquitectura; as-code en las 3 superficies; trabajo orquestado
multi-usuario/multi-rol; GitHub = conector) + este ledger + épica «Experiencia Orquestada»
(borrador de norte, F0 pendiente de firma EN este repo) + kit dev instalado como plugin desde
el marketplace (`alpacapurpura/prenter-marketplace`). Decisiones técnicas heredadas que siguen
vigentes: DH-10 (conexión CC = driver CLI-nativo, BYO licencia, SIN API) · descriptor de
proceso I-77 (contrato de ecosistema — el kit lo shipea, DevHub lo interpreta) · binario Go +
UI embebida. El código go+ui NO se copió: el port es gradual, gobernado por la épica, pieza
por pieza según la experiencia decidida.

*Conecta:* DH-01..DH-11 (la historia en la incubadora; DH-11 = la ficha espejo de esta
graduación en el monorepo) · I-NN de ecosistema (graduación — registrada en
`prenter-harness/tooling/strategy/LEDGER.md`) · I-77 (descriptor) · DH-10 (driver CLI-nativo)
· KIT-06 (marketplace del que este repo consume su arnés).

*Siguiente:* F0 de la épica «Experiencia Orquestada» — norte firmado (journey por rol +
principios de la experiencia) + auditoría de herencia (qué pieza del monorepo entra primero y
en qué forma).

<!-- Próximas: DH-13, DH-14, … -->

## Log

| Fecha | Decisión | Fichas |
|---|---|---|
| 2026-07-04 | Fundación del repo propio: graduación de P2 (nombre DevHub confirmado — binario `devhub`, muere la colisión `cockpit` de nacimiento); visión ampliada = construir y mantener software basado en proceso y arquitectura, trabajo orquestado multi-usuario (CTO·developer·devops·PO), GitHub conector; célula del monorepo congelada como fuente del port gradual; kit dev como plugin del marketplace; épica «Experiencia Orquestada» sembrada (F0 por firmar). | DH-12 |
