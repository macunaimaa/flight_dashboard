import { useAircraftStore } from "../../store/aircraftStore";
import {
  formatAltitude,
  formatSpeed,
} from "../../utils/formatting";

export function AircraftList() {
  const aircraftList = useAircraftStore((s) => s.aircraftList);
  const selectedIcao24 = useAircraftStore((s) => s.selectedIcao24);
  const selectAircraft = useAircraftStore((s) => s.selectAircraft);

  return (
    <div style={styles.list}>
      {aircraftList.length === 0 && (
        <div style={styles.empty}>No aircraft in view</div>
      )}
      {aircraftList.map((ac) => (
        <div
          key={ac.icao24}
          style={{
            ...styles.item,
            ...(ac.icao24 === selectedIcao24 ? styles.selected : {}),
          }}
          onClick={() => selectAircraft(ac.icao24)}
        >
          <div style={styles.row}>
            <span style={styles.callsign}>
              {ac.callsign || ac.icao24.toUpperCase()}
            </span>
            <span style={styles.country}>{ac.originCountry}</span>
          </div>
          <div style={styles.row}>
            <span style={styles.detail}>
              {formatAltitude(ac.baroAltitude)}
            </span>
            <span style={styles.detail}>{formatSpeed(ac.velocity)}</span>
            <span
              style={{
                ...styles.badge,
                background: ac.onGround ? "#374151" : "#1e3a5f",
              }}
            >
              {ac.onGround ? "GND" : "AIR"}
            </span>
          </div>
        </div>
      ))}
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  list: {
    flex: 1,
    overflowY: "auto",
    padding: 0,
  },
  empty: {
    color: "#6b7280",
    fontSize: 13,
    textAlign: "center",
    padding: 24,
  },
  item: {
    padding: "8px 16px",
    borderBottom: "1px solid #1f2937",
    cursor: "pointer",
    transition: "background 0.15s",
  },
  selected: {
    background: "#1e3a5f",
  },
  row: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    gap: 8,
  },
  callsign: {
    color: "#e0e6ed",
    fontWeight: 600,
    fontSize: 13,
    fontFamily: "monospace",
  },
  country: {
    color: "#6b7280",
    fontSize: 11,
  },
  detail: {
    color: "#9ca3af",
    fontSize: 11,
  },
  badge: {
    color: "#e0e6ed",
    fontSize: 10,
    padding: "1px 6px",
    borderRadius: 3,
    fontWeight: 600,
  },
};
