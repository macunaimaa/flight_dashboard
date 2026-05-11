import { useEffect, useRef } from "react";
import { useAuthStore } from "../store/authStore";
import { useAircraftStore } from "../store/aircraftStore";
import { useUIStore } from "../store/uiStore";
import { WebSocketManager } from "../ws/connection";
import type { WSIncomingMessage } from "../ws/types";
import type { AircraftState } from "../api/types";

export function useWebSocket() {
  const token = useAuthStore((s) => s.token);
  const applyUpdates = useAircraftStore((s) => s.applyUpdates);
  const applyRemovals = useAircraftStore((s) => s.applyRemovals);
  const setConnectionStatus = useUIStore((s) => s.setConnectionStatus);
  const managerRef = useRef<WebSocketManager | null>(null);

  useEffect(() => {
    if (!token) return;

    const handleMessage = (msg: WSIncomingMessage) => {
      switch (msg.type) {
        case "aircraft.update":
          applyUpdates(msg.payload as AircraftState[]);
          break;
        case "aircraft.remove":
          applyRemovals(msg.payload as string[]);
          break;
        // Phase 2: handle alert.triggered
      }
    };

    const manager = new WebSocketManager(
      token,
      handleMessage,
      setConnectionStatus
    );
    manager.connect();
    managerRef.current = manager;

    return () => {
      manager.disconnect();
      managerRef.current = null;
    };
  }, [token, applyUpdates, applyRemovals, setConnectionStatus]);
}
