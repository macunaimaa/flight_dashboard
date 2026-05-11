import { useAircraftStore } from "../../store/aircraftStore";
import {
  formatAltitude,
  formatSpeed,
  formatVerticalRate,
  formatHeading,
  formatCoordinate,
  timeAgo,
} from "../../utils/formatting";

export function AircraftInfoPanel() {
  const aircraft = useAircraftStore((s) =>
    s.selectedIcao24 ? s.aircraft.get(s.selectedIcao24) : undefined
  );
  const selectAircraft = useAircraftStore((s) => s.selectAircraft);

  if (!aircraft) return null;

  const [lon, lat] = aircraft.location.coordinates;

  return (
    <div style={styles.panel}>
      <div style={styles.header}>
        <div>
          <h3 style={styles.callsign}>
            {aircraft.callsign || aircraft.icao24.toUpperCase()}
          </h3>
          <span style={styles.icao}>{aircraft.icao24.toUpperCase()}</span>
        </div>
        <button onClick={() => selectAircraft(null)} style={styles.closeBtn}>
          &times;
        </button>
      </div>

      <div style={styles.grid}>
        <InfoRow label="Country" value={aircraft.originCountry} />
        <InfoRow label="Altitude" value={formatAltitude(aircraft.baroAltitude)} />
        <InfoRow label="Speed" value={formatSpeed(aircraft.velocity)} />
        <InfoRow label="Heading" value={formatHeading(aircraft.trueTrack)} />
        <InfoRow label="V/S" value={formatVerticalRate(aircraft.verticalRate)} />
        <InfoRow label="Lat" value={formatCoordinate(lat, "lat")} />
        <InfoRow label="Lon" value={formatCoordinate(lon, "lon")} />
        <InfoRow
          label="Squawk"
          value={aircraft.squawk || "N/A"}
        />
        <InfoRow
          label="Status"
          value={aircraft.onGround ? "On Ground" : "Airborne"}
        />
        <InfoRow
          label="Updated"
          value={timeAgo(aircraft.updatedAt)}
        />
      </div>
    </div>
  );
}

function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div style={styles.row}>
      <span style={styles.label}>{label}</span>
      <span style={styles.value}>{value}</span>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  panel: {
    position: "absolute",
    bottom: 16,
    right: 16,
    width: 300,
    background: "#111827ee",
    borderRadius: 8,
    border: "1px solid #1f2937",
    backdropFilter: "blur(12px)",
    zIndex: 60,
    overflow: "hidden",
  },
  header: {
    display: "flex",
    alignItems: "flex-start",
    justifyContent: "space-between",
    padding: "12px 16px 8px",
    borderBottom: "1px solid #1f2937",
  },
  callsign: {
    color: "#e0e6ed",
    fontSize: 16,
    fontWeight: 700,
    margin: 0,
    fontFamily: "monospace",
  },
  icao: {
    color: "#6b7280",
    fontSize: 11,
  },
  closeBtn: {
    background: "none",
    border: "none",
    color: "#6b7280",
    fontSize: 20,
    cursor: "pointer",
    padding: "0 4px",
    lineHeight: 1,
  },
  grid: {
    padding: "8px 16px 12px",
  },
  row: {
    display: "flex",
    justifyContent: "space-between",
    padding: "3px 0",
  },
  label: {
    color: "#6b7280",
    fontSize: 12,
  },
  value: {
    color: "#e0e6ed",
    fontSize: 12,
    fontFamily: "monospace",
  },
};
