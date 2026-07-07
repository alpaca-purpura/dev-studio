import { useEffect } from "react";
import { SessionRail } from "./widgets/session-rail";
import { ChatPanel } from "./pages/chat/ui/chat-panel";
import { useSessions } from "./shared/store/sessions-store";

export function App() {
  const init = useSessions((s) => s.init);

  useEffect(() => {
    void init();
  }, [init]);

  return (
    <div className="flex h-full">
      <SessionRail />
      <ChatPanel />
    </div>
  );
}
