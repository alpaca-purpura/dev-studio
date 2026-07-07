/** Taxonomía estándar de paquetes de trabajo (spec nuevo-workspace §1). */
export type TipoItem = "historia" | "bug" | "hotfix" | "tarea" | "spike";

export const TIPO_META: Record<TipoItem, { label: string; icon: string; branch: string; commit: string }> = {
  historia: { label: "Historia", icon: "✦", branch: "feature/…", commit: "feat:" },
  bug: { label: "Bug", icon: "🐞", branch: "bugfix/…", commit: "fix:" },
  hotfix: { label: "Hotfix", icon: "🔥", branch: "hotfix/…", commit: "fix:" },
  tarea: { label: "Tarea", icon: "🔧", branch: "chore/…", commit: "chore:" },
  spike: { label: "Spike", icon: "🧪", branch: "spike/…", commit: "—" },
};

export const TIPOS: TipoItem[] = ["historia", "bug", "hotfix", "tarea", "spike"];
