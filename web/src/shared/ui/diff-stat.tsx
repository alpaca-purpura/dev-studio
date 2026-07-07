/** Contador +N −M de un workspace/archivo — mono, ok/crit. */
export function DiffStat({ add, del }: { add: number; del: number }) {
  if (add === 0 && del === 0) return null;
  return (
    <span className="inline-flex items-center gap-1 font-mono text-xs">
      {add > 0 && <span className="text-ok">+{add}</span>}
      {del > 0 && <span className="text-crit">-{del}</span>}
    </span>
  );
}
