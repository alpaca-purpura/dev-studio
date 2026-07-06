# DevStudio — Visión de producto

> Norte vivo del producto. Registro de decisiones: [`LEDGER.md`](./LEDGER.md) (fichas `DH-NN`).
> Este repo nació de la graduación de la célula P2 del monorepo `prenter-harness` (2026-07-04);
> la historia DH-01..DH-11 vive allá.

## Identidad

**DevStudio es una aplicación para construir y mantener software basado en PROCESO y ARQUITECTURA.**

No es un tracker de tickets ni un facilitador de sesiones de IA: es la consola donde una
organización ejecuta su ciclo de desarrollo completo — con el proceso, la arquitectura y la
documentación como **dato versionado**, no como prosa suelta:

- **Proceso as code** — descriptor de proceso (contrato I-77): estados, transiciones, gates,
  dueños y verbos (CDEvents) como dato; la consola deriva TODO de ahí, cero hardcode del ciclo.
- **Arquitectura as code** — `arquitectura.yaml` por sistema (`meta.clase`); la app la renderiza
  y la usa como contexto de trabajo.
- **Documentación as code** — visión, specs, decisiones (ledger) viajan en el repo y alimentan
  las sesiones de trabajo.

## La experiencia: trabajo orquestado

El usuario — **CTO, developer, devops o product owner** — trabaja DENTRO de la aplicación, en
una organización, de forma **orquestada**: usuarios complementarios trabajando a la par sobre
uno o varios sistemas simultáneamente. El conector entre todos es el **repositorio (GitHub)**:
el repo es el punto de encuentro — código, proceso, arquitectura y documentación convergen ahí.

Cada rol ve y opera lo suyo (roles de primera clase); el proceso gobierna quién hace qué y
cuándo (gates, dueños, transiciones); la entrega se orquesta: ticket especificado → sesión de
agente parametrizada por el proceso de la empresa + contexto as-code inyectado → gates → el
humano encausa y aprueba.

## Arquitectura técnica (decisiones heredadas vigentes)

- **App de escritorio multiplataforma** (Windows · Linux · macOS): binario Go con UI embebida
  (`go:embed`); instalador = el binario solo.
- **Conexión con Claude Code = driver CLI-nativo** (DH-10): la app spawnea el `claude` que el
  usuario ya tiene instalado, vía stdin/stdout (stream-json). **BYO licencia — sin API de
  Anthropic**, la app jamás toca credenciales. ⚠ Verificar Consumer Terms por escrito antes de
  vender instalables.
- **Arneses de construcción de ESTE repo**: kit dev (plugin, marketplace
  `alpacapurpura/prenter-marketplace`) — se instala pineado y evoluciona con el producto.

## Herencia (port gradual, gobernado por la épica)

El monorepo (`prenter-harness/products/devhub/`, congelado) tiene construido: torre de control
2 ejes · descriptor de proceso + board gobernado · cockpit de delivery (sesiones de agente) ·
lente arquitectura · multi-workspace · 27 endpoints Go + UI Next. **Nada se copia wholesale**:
cada pieza entra cuando la experiencia decidida la pida — épica «Experiencia Orquestada»
([`epicas/experiencia-orquestada/`](./epicas/experiencia-orquestada/)).

## TBD (heredados del grill + nuevos)

- Comprador con nombre · pricing · éxito a 12m.
- Modelo de "organización" y sincronización multi-usuario vía GitHub (¿qué es dato compartido
  vs local?).
- Roles/accesos de primera clase (hoy no hay auth).
