# SPEC — Registry propio de arneses (PB-25 ⊕ restos PB-05)

> Estado: **CONGELADA 2026-07-07** — 6 forks ratificados por Chris en sesión (ficha DH-18);
> los 2 forks de formato/materialización se RE-firmaron tras la revisión de doctrina de
> ArnesIA (`~/Proyectos/harness-studio`, mandato de Chris: «no construyas por construir»).
> Materializa el reframe F0 (DH-14): **rol = arnés instalado desde registry propio** — cero
> roles locales, la curaduría vive en el registry. PB-05 se FUSIONA acá (fork 1): su único
> resto con journey (roster por proyecto + conexión del registry) ES esta rebanada;
> metadatos cosméticos declarados sin journey.

## 0. Decisión madre: DevStudio ADOPTA el estándar ArnesIA (cruce firmado)

ArnesIA es la fábrica de arneses; DevStudio es la «app de rol» de su propio diagrama de
ecosistema («publica → marketplace git → instala → proyecto → ejecuta ← DevHub y apps de
rol»). DevStudio **no inventa formato**: consume el estándar firmado e implementado de la
fábrica y replica su patrón de inyección. Pedidos a ArnesIA = solo aditivos (prompt de
interop entregado a Chris en sesión; nada bloquea esta rebanada).

Fuentes normativas (cross-repo por contrato, jamás por ruta en runtime):
- `harness-studio/arch/contracts/nomenclatura-arnes.md` **v1 FIRMADA** (forma-plugin /
  forma-instalada, manifiesto `arnes.l0.json`, tabla clase→ubicación).
- `harness-studio/METODOLOGIA.md` §9 «Los 3 cuerpos» (inyección por flags, jamás escribir
  la maquinaria en el árbol del proyecto).
- `prenter-marketplace` (formato marketplace git EN PRODUCCIÓN: `.claude-plugin/
  marketplace.json` + `catalogo.json` canales/versiones + carpeta por plugin/versión).
- Arnés real de referencia: `harness-studio/dogfood/dev-full-cycle/` (plugin.json +
  arnes.l0.json con spine + 4 skills con contrato + CLAUDE.md banda Base).

## 1. Los 3 planos del mecanismo

### 1.1 Registry (fuente) — marketplace git de la organización

- **Fuente**: URL git (se clona con el **git del usuario** — BYO credenciales, espíritu
  DH-10) o **ruta local** directa (un registry clonado ES una carpeta; dev-loop). Config
  a nivel APP (norte: «al instalar la aplicación me conectaré»), persistida en
  `state.json` (`registry: {source}`). Clon/actualización en `~/.dev-studio/registry/`.
- **Formato leído**: `catalogo.json` (canales/versiones, estado habilitada/deprecada) +
  `.claude-plugin/marketplace.json` (plugins[] name/source/description) + por arnés:
  `arnes.l0.json` (id · rol · proceso · fases · spine · canal) y `.claude-plugin/
  plugin.json` (name/description/version). Sin catalogo.json se degrada honesto a
  marketplace.json solo (versión única).
- **Sync**: adaptador NUEVO `adapters/registry/gitsync` (clone/pull con el git del
  usuario) **confinado a `~/.dev-studio/registry/`**. El boundary
  `git-solo-lectura-y-commit` NO se toca (su scope = repo DEL usuario,
  `adapters/git/cli`); gana nota v1.1 con la distinción + fitness test nuevo (§5).
- Publicar un arnés custom = push al repo del registry (fuera de la app v1 — lo hace
  ArnesIA/la org).

### 1.2 Instalación (por Proyecto) — modelo npm: lock en repo + caché local

- Instalar arnés en un repo =
  1. copiar la **forma-plugin INTACTA** a caché `~/.dev-studio/arneses/{id}/{version}/`;
  2. escribir/actualizar **`.devstudio/arneses.yaml`** en la raíz del repo (el ROSTER
     as-code: `registry` de origen + lista `{id, version, canal}`);
  3. commit por pathspec del lock: `chore(devstudio): instala arnés {id}@{version}`
     (puerto `GitCommit` existente — jamás `add .`).
- El lock viaja por GitHub → otro usuario/máquina **rehidrata** del registry (multi-usuario
  ready, principio 6). Desinstalar = quitar del lock + commit (el caché se conserva).
- El árbol del proyecto queda LIMPIO: ni `.claude/` fundido ni payload vendorizado
  (METODOLOGIA §9 de la fábrica: la maquinaria no se escribe en el árbol). Cada arnés
  queda auditable como unidad por el loader de ArnesIA.

### 1.3 Ejecución (por Sesión) — inyección por flags nativos (DH-10 intacto)

Sesión de trabajo con rol = arnés instalado → el spawn suma (mismo patrón HS-11 que la
fábrica usa para su propio kit, flags verificados en la CLI):

```
claude -p … --plugin-dir ~/.dev-studio/arneses/{id}/{version}
           --append-system-prompt "<preámbulo del rol>\n\n## Banda Base…\n<CLAUDE.md>"
```

> Ajuste cazado EN el gate: la CLI rechaza `--append-system-prompt` y
> `--append-system-prompt-file` JUNTOS («use only one») — la banda Base (CLAUDE.md del
> arnés) viaja EMBEBIDA en el único `--append-system-prompt`, no como flag aparte.

- **Preámbulo compuesto** (DevStudio lo deriva, cero campo nuevo pedido): rol · proceso ·
  fases del `arnes.l0.json` («Trabajás como {rol} en el proceso {proceso}…»).
- La sesión carga **SOLO el arnés de SU rol** (skills namespaced por plugin) — contexto
  limpio; los demás roles instalados no contaminan.
- `Session.Rol` pasa a ser el **id del arnés**. Exploración (CAP-12) sigue sin rol.
- Sin arnés instalado no hay rol elegible: **muere `MOCK_ROSTER`** (fork 6 — muere en
  PB-25, no en PB-06). Sesión de trabajo puede crearse **sin rol** (elegirlo es opcional
  mientras el roster esté vacío — honestidad: la app no finge roles que no existen).

## 2. Contratos técnicos

- `domain.Arnes{ID, Nombre, Descripcion, Rol, Proceso, Version, Canal, Fases []string}` —
  parseado de `arnes.l0.json` + `plugin.json` (nombre/descripción caen a plugin.json si el
  manifiesto no los trae). El **spine NO se interpreta** en esta rebanada (queda en el
  archivo del caché para PB-08 — reconciliación spine⟷I-77 fichada allá).
- `domain.RegistryEstado{Source, SyncedAt, Arneses []Arnes}` · lock
  `domain.ArnesInstalado{ID, Version, Canal}`.
- **Contrato estable del lock `.devstudio/arneses.yaml`** (declarado 2026-07-07 a pedido
  recíproco de ArnesIA — el lock es superficie de auditoría in situ, detector 3° de
  nomenclatura-arnes v1.1, leído read-only por la fábrica): campos `registry` (source del
  marketplace, top-level) + `arneses[].{id, version, canal}`. **Evolución solo aditiva** —
  esos 4 campos no se renombran ni cambian de semántica; campos nuevos solo SUMAN. Fichado
  en LEDGER (DH-18.2).
- **Cadena de fallback canónica para pintar** (bendecida por ArnesIA, HS-12; schema L0
  reparado — `arnes.l0.nombre` es el campo canónico): `nombre → plugin.json name → id`,
  ídem `descripcion`.
- Ports nuevos: `RegistrySync` (Sync(ctx, source) → dir) · `RegistryCatalog`
  (Leer(ctx, dir) → []Arnes) · `ArnesStore` (lock del repo: Leer/Instalar/Desinstalar) —
  adapters `registry/gitsync`, `registry/fscatalog`, `arneses/lockfile`.
- `ports.SpawnOpts` gana `PluginDirs []string` + `SystemPrompt string` (preámbulo +
  banda Base embebida) → `buildArgs` (pura, test) emite los flags.
- API: `GET /api/registry` (estado+catálogo) · `PUT /api/registry {source}` (conectar+sync)
  · `POST /api/registry/sync` · `GET/POST /api/repos/{id}/arneses` ·
  `DELETE /api/repos/{id}/arneses/{arnesId}`. Guards 422: source inválido/sin clonar ·
  instalar id inexistente en catálogo · sesión con rol no instalado en el repo.
- UI (Config overlay): sección **Registry** (conectar source + estado + Actualizar) +
  roster REAL de instalados agrupado por `proceso` (grupos mock Builder/Auditor/Humano
  mueren con el mock) + instalar/desinstalar desde el catálogo (lista simple — la vista
  Roles rica sigue siendo PB-06). Stories nuevas (RN-9).

## 3. Reglas de negocio

- **RN-1** El registry es la ÚNICA fuente de roles (DH-14): cero creación local. Roster
  vacío se muestra vacío con destino («conectá el registry / instalá desde el catálogo»).
- **RN-2** El lock `.devstudio/arneses.yaml` es el SSoT del roster del proyecto; la app lo
  commitea por pathspec y NUNCA toca otros archivos del repo del usuario.
- **RN-3** `gitsync` opera EXCLUSIVAMENTE bajo `~/.dev-studio/registry/` (fitness). El
  adapter `git/cli` del usuario sigue sin push/pull/fetch (scanner intacto).
- **RN-4** Payload del arnés = solo lectura para DevStudio: se copia INTACTO, jamás se
  edita (editar arneses es trabajo de ArnesIA, no de la app de rol).
- **RN-5** Sesión con rol exige arnés instalado en SU repo (422); el spawn inyecta solo
  los artefactos de ese arnés.

## 4. AC — gate de cierre (verificación REAL contra la app instalada)

- [x] **AC-1** Conectar registry (ruta local a un marketplace git REAL con
  `dev-full-cycle` copiado de harness-studio): Config muestra el catálogo con
  rol/proceso/versión. Source inválido → error visible (422). *(2026-07-07 en vivo:
  catálogo con dev-full-cycle 0.1.0·beta —rol/proceso/fases del arnes.l0.json REAL de
  ArnesIA— + demo-auditor 0.2.0 degradado; `/tmp/no-existe-nada` → 422 con el error real
  de git)*
- [x] **AC-2** Instalar `dev-full-cycle` en un repo: caché con forma-plugin intacta
  (plugin.json + arnes.l0.json + skills + CLAUDE.md, `diff -r` contra el registry = 0) y
  `.devstudio/arneses.yaml` committeado — `git show` lo confirma, SOLO el lock en el
  commit. *(en vivo: DIFF=0; commit `2c90f41 chore(devstudio): instala arnés
  dev-full-cycle@0.1.0`, 1 file changed, 8 insertions)*
- [x] **AC-3** Roster real: el instalado aparece en Config agrupado por proceso (mock
  MUERTO — grep `MOCK_ROSTER` = 0); crear sesión de trabajo con ese rol. *(en vivo desde
  la app instalada: wizard→ítem tarea→Config con «ROSTER DEL PROYECTO · DEMO-GAMMA»
  agrupado por proceso + detalle rol/proceso/fases; sesión `chore/probar-arnes-…` con
  chip del rol; grep=0 — screenshot ac3-config-roster.png)*
- [x] **AC-4** Spawn con arnés: el proceso `claude` de la sesión corre con `--plugin-dir`
  + `--append-system-prompt` visibles en `/proc/{pid}/cmdline`; un turno real confirma
  las skills namespaced del arnés. *(en vivo: cmdline con `--plugin-dir
  ~/.dev-studio/arneses/dev-full-cycle/0.1.0` + preámbulo compuesto + banda Base
  std-spec embebida; el turno respondió EXACTO las 4 skills:
  `dev-full-cycle:builder/releaser/reviewer/spec-writer`)*
- [x] **AC-5** Desinstalar: lock actualizado + commit + roster lo refleja; sesión nueva
  con ese rol → 422. *(en vivo: ciclo instala/desinstala demo-auditor — commits
  `c6f986a`/`aa5d980`, roster queda solo dev-full-cycle, POST con rol desinstalado → 422
  con mensaje accionable)*
- [x] **AC-6** Guards API: instalar id inexistente → 422 · sesión trabajo con rol no
  instalado → 422 · registry sin conectar → roster vacío honesto. *(en vivo: 422×2 +
  `/api/registry` inicial `{"source":"","arneses":null}`)*
- [x] **AC-7** Suite verde: unit (fscatalog sobre fixture marketplace · lockfile ·
  buildArgs · gitsync confinado · ArnesService con dobles · inyección en spawn) + fitness
  (registry-sin-push nuevo + scanner git/cli intacto + hexagonal). *(`go test ./...`
  verde tras el fix del flag; tsc + vite + build-storybook verdes)*

**Bonus del gate (dogfooding real):** la app se actualizó a sí misma DOS veces durante el
gate (POST /api/update → rebuild → syscall.Exec, misma PID) — la primera cazó y arregló
un bug del updater (`npm: orden no encontrada` bajo el PATH pelado del .desktop: nvm no
entra por `bash -lc` → fallbacks explícitos en install.sh); la segunda cazó el conflicto
`--append-system-prompt` vs `-file` (excluyentes) con el proceso claude muriendo al nacer.

## 5. Boundary tocado

`arch/boundaries/git-solo-lectura-y-commit.md` → v1.1: nota de scope (repo del usuario)
+ check nuevo `registry-confinado` (el paquete `adapters/registry` no acepta destinos
fuera de `~/.dev-studio/registry/` — test). El principio L1 (least authority sobre el
estado del usuario) queda MÁS fuerte: el sync vive en un paquete que no puede ver repos
de usuario.

## 6. Fuera de alcance (con destino)

- Vista Roles rica (explorar catálogo con detalle/versiones/upgrade) → **PB-06**.
- Interpretar spine/proceso (pintar backlog del descriptor) + reconciliación I-77 → **PB-08**.
- Publicar arneses desde la app / permisos por puesto → ArnesIA fase 5 / **PB-21→23**.
- Upgrade de versión instalada (hoy: desinstalar+instalar) → PB-06 o ficha nueva.
- ~~Respuestas del prompt interop ArnesIA (categoria en spine, campo nombre, bendición del
  lock) → se incorporan por changelog cuando lleguen; nada bloquea.~~ **RECIBIDAS
  2026-07-07** (ficha HS-12 de la fábrica) — incorporadas en changelog v1.1.

## Changelog

- v1 CONGELADA 2026-07-07 — 6 forks (PB-05 fusión · formato ArnesIA as-is [re-firmado
  tras revisión doctrina] · registry marketplace git + git del usuario · materialización
  lock+caché+flags [re-firmado] · roster mock muere acá · alcance UI en Config).
- v1 VERIFICADA 2026-07-07 — 7/7 ACs en vivo contra el binario instalado, con el arnés
  REAL `dev-full-cycle` de ArnesIA (interop de verdad, no fixture inventado). 2 bugs
  cazados y arreglados EN el gate (updater sin npm en PATH .desktop · flags de system
  prompt excluyentes). Ficha DH-18.
- v1.1 2026-07-07 — **respuestas del interop ArnesIA incorporadas** (5 pedidos
  respondidos, ficha HS-12 de la fábrica): **(1)** publish fase-5 RATIFICADO — contrato
  estable = marketplace.json + forma-plugin intacta con `arnes.l0.json` (evolución solo
  aditiva) + `catalogo.json {marketplace, canales{canal→versión}, versiones[].{version,
  estado: habilitada|deprecada, fuente, fecha}}`; DevStudio puede depender de esos campos
  (PB-06 versiones/upgrade construye sobre esto). **(2)** `spine.categorias` ACEPTADO:
  mapa hermano opcional estado→categoría, enum FIJO = I-77 RN-28
  (`propuesto·en-progreso·completado·descartado·pausado`), terminalidad DERIVADA
  (∈{completado,descartado}) — ya en schema L0 + dogfood dev-full-cycle; la semilla
  `marketplace-arneses` 0.1.0 NO se mutó → para PB-08 rehidratar del dogfood actualizado
  o esperar el próximo publish. **(3)** campo canónico para pintar = `arnes.l0.nombre`
  (schema reparado); cadena de fallback bendecida `nombre → plugin.json name → id`, ídem
  descripcion — fscatalog alineado (leía plugin.json pisando l0.descripcion; faltaba el
  eslabón final →id). **(4)** lock `.devstudio/arneses.yaml` BENDECIDO como superficie de
  auditoría in situ (detector 3° nomenclatura-arnes v1.1, ArnesIA lo lee read-only);
  pedido recíproco cumplido: contrato estable `registry·id·version·canal` declarado en §2
  + fichado LEDGER DH-18.2. **(5)** spine⟷I-77 CONFORME: spine(+categorias) = subconjunto
  navegable canónico; gates/dueños se DERIVAN de los contratos por caja
  (`box.contract.schema.json`: estado «de → a» · gate{tipo} · ruta[]); el arnés NO shipea
  descriptor I-77 aparte — un I-77 materializado es proyección/export del arnés.
  Derivación bendecida: dueño de estado = caja cuya transición LLEGA a él · transiciones
  con dueño-caja = de rol, resto = operador · terminalidad = categoría. PB-08 re-redactada
  con este modelo.
