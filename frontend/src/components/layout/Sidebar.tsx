import { useUIStore } from "../../store/uiStore";
import { AircraftList } from "../aircraft/AircraftList";

export function Sidebar() {
  const sidebarOpen = useUIStore((s) => s.sidebarOpen);

  if (!sidebarOpen) return null;

  return (
    <div style={styles.sidebar}>
      <div style={styles.header}>
        <h3 style={styles.title}>Aircraft</h3>
      </div>
      <AircraftList />
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  sidebar: {
    width: 320,
    height: "100%",
    background: "#111827",
    borderRight: "1px solid #1f2937",
    display: "flex",
    flexDirection: "column",
    position: "relative",
    zIndex: 50,
  },
  header: {
    padding: "12px 16px",
    borderBottom: "1px solid #1f2937",
  },
  title: {
    color: "#e0e6ed",
    fontSize: 14,
    fontWeight: 600,
    margin: 0,
  },
};
