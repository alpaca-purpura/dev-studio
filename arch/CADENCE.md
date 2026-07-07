# CADENCE — cómo vive el árbol de arquitectura

> La arquitectura as code no se revisa por calendario; se revisa **al cambiar**. Mismo mecanismo
> que el proyecto hermano `harness-studio` (`arch/CADENCE.md`), copiado porque es una práctica
> de la casa, no un detalle de ese repo. Norte: [`../VISION.md`](../VISION.md) · el «por qué»
> firmado: [`../LEDGER.md`](../LEDGER.md).

## Las tres separaciones

- **El «por qué»** → una ficha `DH-NN` del [`../LEDGER.md`](../LEDGER.md). Nunca se duplica acá.
- **La «prueba»** → una check en [`fitness/`](./fitness/) que falla `go test` si el código viola
  la regla. Una arquitectura que no se puede romper en CI es un deseo, no una restricción.
- **El «dibujo»** → [`model/`](./model/), texto que renderiza (Mermaid), nunca sincronizado a mano.

## Anatomía de un boundary node (`boundaries/<regla>.md`)

Igual formato que `harness-studio` — front-matter (`regla/version/status/ledger/sources/
enforced_by/severity`) + `## L1 · Principio` + `## L2 · Realización` + `## Checklist evaluable`
+ `## Changelog`. Ver cualquier nodo existente en [`boundaries/`](./boundaries/) como plantilla.

## El ritual (al cambiar, no por calendario)

1. **Decisión estructural** nace en una conversación → ficha `DH-NN` (o la muta).
2. **¿Toca un boundary?** Si crea/cambia una regla enforçable → nodo nuevo o bump de versión. Si
   es una decisión sin regla verificable (p.ej. «usamos Zustand») → basta la ficha, no todo va acá.
3. **L1 + L2**, con `⇐ L1` explícito en cada afirmación de mapeo. Divergencia → `⚠ divergencia`
   marcada y justificada (nunca silenciosa) — ver `sesion-aislada-por-cwd` como ejemplo real.
4. **Check + enforcer.** Si el check ya corre, `status: enforced` + el test en `enforced_by:`. Si
   no, `status: proposed` y la brecha se documenta explícitamente (no se omite el nodo).
5. **Diagrama, si cambió la topología.** Hoy un solo binario — cuando F2+ agregue un segundo
   componente desplegable, `model/` suma un container diagram (C4 nivel 2).
6. **Bump + changelog.**

## Reglas del árbol (no re-negociar — heredadas de la casa)

- **El LEDGER es el «por qué» canónico.** `arch/` no forkea el diario, lo proyecta.
- **Todo boundary DEBE declarar su `enforced_by:`**, aunque hoy sea "code review, TBD
  automatizar" — eso es honesto; omitir el campo no lo es.
- **Superar, no borrar.** Un boundary superado pasa a `status: superseded` con puntero al
  reemplazo.
- **Honestidad ante todo:** `status: enforced` exige un test que corre y pasa HOY (`go test
  ./arch/fitness/...`). Una brecha real (ver `sesion-aislada-por-cwd`) se documenta con checks
  marcados `❌ TBD`, nunca se pinta verde para que el árbol se vea mejor.

## Hacia dónde va

F1 entregó 5 boundaries, 2 con test real. F2+ (port por rebanadas, gobernado por la épica
«Experiencia Orquestada») sumará nodos a medida que cada rebanada traiga una regla estructural
nueva — nunca por anticipado.
