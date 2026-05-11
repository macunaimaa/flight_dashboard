import { useAuthStore } from "../../store/authStore";
import { useUIStore } from "../../store/uiStore";
import { useAircraftStore } from "../../store/aircraftStore";
import type { CSSProperties } from "react";

export function TopBar() {
  const logout = useAuthStore((s) => s.logout);
  const tenantId = useAuthStore((s) => s.tenantId);
  const connectionStatus = useUIStore((s) => s.connectionStatus);
  const toggleSidebar = useUIStore((s) => s.toggleSidebar);
  const aircraftCount = useAircraftStore((s) => s.aircraftList.length);
  const weatherPanelOpen = useUIStore((s) => s.weatherPanelOpen);
  const newsPanelOpen = useUIStore((s) => s.newsPanelOpen);
  const toggleWeatherPanel = useUIStore((s) => s.toggleWeatherPanel);
  const toggleNewsPanel = useUIStore((s) => s.toggleNewsPanel);

  const statusColor =
    connectionStatus === "connected"
      ? "#22c55e"
      : connectionStatus === "connecting"
        ? "#eab308"
        : "#ef4444";

  return (
    <div style={styles.bar}>
      <div style={styles.left}>
        <button onClick={toggleSidebar} style={styles.menuBtn}>
          &#9776;
        </button>
        <span style={styles.brand}>Aircraft Dashboard</span>
        <span style={styles.tenant}>{tenantId}</span>
      </div>
      <div style={styles.center}>
        <span style={styles.stat}>{aircraftCount} aircraft tracked</span>
      </div>
      <div style={styles.right}>
        <button
          onClick={toggleWeatherPanel}
          style={weatherPanelOpen ? styles.activeBtn : styles.actionBtn}
        >
          Weather
        </button>
        <button
          onClick={toggleNewsPanel}
          style={newsPanelOpen ? styles.activeBtn : styles.actionBtn}
        >
          News
        </button>
        <span style={{ ...styles.status, background: statusColor }} />
        <span style={styles.statusText}>{connectionStatus}</span>
        <button onClick={logout} style={styles.logoutBtn}>
          Logout
        </button>
      </div>
    </div>
  );
}

const styles: Record<string, CSSProperties> = {
  bar: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    height: 48,
    background: "#111827",
    borderBottom: "1px solid #1f2937",
    padding: "0 16px",
    zIndex: 100,
    position: "relative",
  },
  left: { display: "flex", alignItems: "center", gap: 12 },
  center: { display: "flex", alignItems: "center" },
  right: { display: "flex", alignItems: "center", gap: 10 },
  menuBtn: {
    background: "none",
    border: "none",
    color: "#9ca3af",
    fontSize: 18,
    cursor: "pointer",
    padding: "4px 8px",
  },
  brand: { color: "#e0e6ed", fontWeight: 700, fontSize: 14 },
  tenant: {
    color: "#6b7280",
    fontSize: 12,
    background: "#1f2937",
    padding: "2px 8px",
    borderRadius: 4,
  },
  stat: { color: "#9ca3af", fontSize: 13 },
  status: {
    width: 8,
    height: 8,
    borderRadius: "50%",
    display: "inline-block",
  },
  statusText: { color: "#9ca3af", fontSize: 12 },
  actionBtn: {
    background: "none",
    border: "1px solid #374151",
    color: "#9ca3af",
    borderRadius: 4,
    padding: "4px 12px",
    fontSize: 12,
    cursor: "pointer",
  },
  activeBtn: {
    background: "#2563eb33",
    border: "1px solid #2563eb",
    color: "#93c5fd",
    borderRadius: 4,
    padding: "4px 12px",
    fontSize: 12,
    cursor: "pointer",
  },
  logoutBtn: {
    background: "none",
    border: "1px solid #374151",
    color: "#9ca3af",
    borderRadius: 4,
    padding: "4px 12px",
    fontSize: 12,
    cursor: "pointer",
  },
};
