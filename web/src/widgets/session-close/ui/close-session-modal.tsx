import { useState } from "react";
import { useSessions } from "../../../shared/store/sessions-store";
import { useUi } from "../../../shared/store/ui-store";
import { Button, Chip, ModalShell } from "../../../shared/ui";

/** Modal de cierre (spec workspace-aislado §2.2, fork «preguntar»): conservar o borrar el
 *  workspace. Borrar sucio → git lo rechaza y acá se muestra el aviso (RN-3). */
export function CloseSessionModal() {
  const closeRequestId = useUi((s) => s.closeRequestId);
  const requestClose = useUi((s) => s.requestClose);
  const target = useSessions((s) => s.sessions.find((x) => x.id === closeRequestId));
  const closeSession = useSessions((s) => s.closeSession);
  const [busy, setBusy] = useState(false);
  const [aviso, setAviso] = useState<string | null>(null);
  const [nombreCerrado, setNombreCerrado] = useState("");

  // ojo: al cerrar, la sesión sale del store y `target` muere — el AVISO (workspace
  // conservado) debe sobrevivir a ese unmount, por eso se evalúa ANTES que target.
  if (!closeRequestId) return null;
  if (!target && !aviso) return null;

  const cerrar = async (workspace: "keep" | "remove") => {
    if (busy || !target) return;
    setBusy(true);
    setNombreCerrado(target.nombre);
    try {
      const resp = await closeSession(closeRequestId, workspace);
      if (resp.detalle) {
        setAviso(resp.detalle);
      } else {
        requestClose(null);
      }
    } finally {
      setBusy(false);
    }
  };

  return (
    <ModalShell
      open
      onClose={() => {
        setAviso(null);
        requestClose(null);
      }}
      title={aviso ? "Workspace conservado" : "Cerrar sesión"}
      subtitle={target?.nombre ?? nombreCerrado}
      icon="⎇"
      footer={
        aviso ? (
          <Button
            onClick={() => {
              setAviso(null);
              requestClose(null);
            }}
          >
            Entendido
          </Button>
        ) : target?.workspace ? (
          <>
            <Button variant="ghost" onClick={() => requestClose(null)} disabled={busy}>
              Cancelar
            </Button>
            <Button variant="outline" onClick={() => void cerrar("remove")} disabled={busy}>
              Borrar workspace
            </Button>
            <Button onClick={() => void cerrar("keep")} disabled={busy}>
              Conservar workspace
            </Button>
          </>
        ) : (
          <>
            <Button variant="ghost" onClick={() => requestClose(null)} disabled={busy}>
              Cancelar
            </Button>
            <Button onClick={() => void cerrar("keep")} disabled={busy}>
              Cerrar sesión
            </Button>
          </>
        )
      }
    >
      {aviso ? (
        <p className="text-sm leading-relaxed text-muted-foreground">{aviso}</p>
      ) : target?.workspace ? (
        <div className="space-y-3 text-sm leading-relaxed text-muted-foreground">
          <p>
            Esta sesión tiene su propio workspace aislado. ¿Qué hacemos con él al cerrar?
          </p>
          <div className="flex flex-wrap gap-2">
            <Chip>{target?.branch ?? "wt/…"}</Chip>
            <Chip className="max-w-full truncate">{target?.workspace}</Chip>
          </div>
          <p className="text-xs">
            <b className="text-foreground">Conservar</b>: el worktree y su branch quedan en
            disco (el trabajo sigue ahí). <b className="text-foreground">Borrar</b>: se elimina
            solo si está limpio — con cambios sin commit, git lo rechaza y se conserva (jamás
            se fuerza).
          </p>
        </div>
      ) : (
        <p className="text-sm leading-relaxed text-muted-foreground">
          Sesión sin workspace propio (anterior al aislamiento). Se cierra el proceso y
          desaparece del rail.
        </p>
      )}
    </ModalShell>
  );
}
