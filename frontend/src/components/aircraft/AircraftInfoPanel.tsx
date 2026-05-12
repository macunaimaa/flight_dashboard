import { useAircraftStore } from "../../store/aircraftStore";
import { deriveStatus, deriveAirline, deriveRegistration } from "../dashboard/derive";
import { StatusPill } from "../dashboard/StatusPill";
import { AltitudeProfile } from "../dashboard/AltitudeProfile";
import { HeadingDial } from "../dashboard/HeadingDial";
import {
  formatAltitude, formatSpeed, formatVerticalRate, formatHeading, formatCoordinate, timeAgo,
} from "../../utils/formatting";

export function AircraftInfoPanel() {
  const aircraft = useAircraftStore(s => s.aircraft);
  const selectedId = useAircraftStore(s => s.selectedIcao24);
  const select = useAircraftStore(s => s.selectAircraft);
  const ac = selectedId ? aircraft.get(selectedId) : undefined;

  if (!ac) {
    return (
      <aside style={panelStyle}>
        <div style={{ padding: 24, color: "var(--text-3)", fontSize: 12 }}>
          <div className="sd-mono" style={{ fontSize: 9.5, letterSpacing: 0.6, marginBottom: 6 }}>INSPECTOR</div>
          Select an aircraft to inspect.
        </div>
      </aside>
    );
  }

  const status = deriveStatus(ac);
  const airline = deriveAirline(ac.callsign);
  const reg = deriveRegistration(ac.icao24);
  const [lon, lat] = ac.location.coordinates;

  return (
    <aside style={panelStyle}>
      <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between",
        padding: "12px 16px", borderBottom: "1px solid var(--line)" }}>
        <div className="sd-mono" style={{ fontSize: 9.5, letterSpacing: 0.6, color: "var(--text-3)" }}>INSPECTOR</div>
        <button onClick={() => select(null)} style={{
          background: "transparent", border: "none", color: "var(--text-3)", cursor: "pointer", fontSize: 18, padding: 0,
        }}>×</button>
      </div>

      <div style={{ overflow: "auto", flex: 1 }}>
        {/* Header */}
        <div style={{ padding: "14px 16px", borderBottom: "1px solid var(--line)" }}>
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
            <div>
              <div className="sd-mono" style={{ fontSize: 22, fontWeight: 700, color: "var(--text)" }}>{ac.callsign || ac.icao24}</div>
              <div style={{ fontSize: 12, color: "var(--text-2)", marginTop: 2 }}>{airline}</div>
              <div className="sd-mono" style={{ fontSize: 10.5, color: "var(--text-3)", marginTop: 2 }}>
                {ac.icao24.toUpperCase()} · REG {reg} · {ac.originCountry}
              </div>
            </div>
            <StatusPill status={status} />
          </div>
        </div>

        {/* Photo slot */}
        <div style={{ height: 120, margin: 16, borderRadius: 4, background:
          "repeating-linear-gradient(135deg, var(--surface-1) 0 10px, var(--surface-2) 10px 20px)",
          border: "1px solid var(--line-2)", display: "grid", placeItems: "center" }}>
          <span className="sd-mono" style={{ fontSize: 10, color: "var(--text-3)", letterSpacing: 0.6 }}>
            AIRCRAFT PHOTO · drop image
          </span>
        </div>

        {/* Altitude profile */}
        <div style={{ padding: "0 16px 14px" }}>
          <SectionLabel>ALTITUDE PROFILE</SectionLabel>
          <AltitudeProfile aircraft={ac} />
        </div>

        {/* Instrument tiles + heading dial */}
        <div style={{ padding: "0 16px 14px", display: "grid", gridTemplateColumns: "1fr 1fr 96px", gap: 8 }}>
          <Tile label="ALT" value={formatAltitude(ac.baroAltitude)} />
          <Tile label="SPD" value={formatSpeed(ac.velocity)} />
          <HeadingDial heading={ac.trueTrack} />
          <Tile label="VS" value={formatVerticalRate(ac.verticalRate)} tone={
            (ac.verticalRate ?? 0) * 196.85 > 200 ? "mint"
              : (ac.verticalRate ?? 0) * 196.85 < -200 ? "cyan" : undefined
          } />
          <Tile label="HDG" value={formatHeading(ac.trueTrack)} />
          <Tile label="SQK" value={ac.squawk || "—"} tone={status === "EMERGENCY" ? "red" : undefined} />
        </div>

        {/* Position */}
        <div style={{ padding: "0 16px 14px" }}>
          <SectionLabel>POSITION</SectionLabel>
          <div className="sd-mono" style={{ fontSize: 11, color: "var(--text-2)", display: "grid", gap: 4 }}>
            <div>{formatCoordinate(lat, "lat")}</div>
            <div>{formatCoordinate(lon, "lon")}</div>
            <div style={{ color: "var(--text-3)" }}>Last update {timeAgo(ac.updatedAt)}</div>
          </div>
        </div>

        {/* Actions */}
        <div style={{ padding: "0 16px 16px", display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 6 }}>
          {["TRACK", "ALERT", "SHARE"].map(l => (
            <button key={l} className="sd-mono" style={{
              padding: "8px 0", fontSize: 11, fontWeight: 600, letterSpacing: 0.5,
              background: l === "TRACK" ? "var(--accent)" : "var(--surface-2)",
              color: l === "TRACK" ? "var(--bg)" : "var(--text-2)",
              border: l === "TRACK" ? "none" : "1px solid var(--line-2)",
              borderRadius: 3, cursor: "pointer",
            }}>{l}</button>
          ))}
        </div>
      </div>
    </aside>
  );
}

const panelStyle: React.CSSProperties = {
  width: 360, display: "flex", flexDirection: "column",
  background: "var(--surface-0)", borderLeft: "1px solid var(--line)",
};

function SectionLabel({ children }: { children: React.ReactNode }) {
  return <div className="sd-mono" style={{
    fontSize: 9.5, letterSpacing: 0.6, color: "var(--text-3)", marginBottom: 6,
  }}>{children}</div>;
}

function Tile({ label, value, tone }: { label: string; value: string; tone?: "mint" | "cyan" | "red" }) {
  const color = tone === "mint" ? "var(--mint)" : tone === "cyan" ? "var(--cyan)" : tone === "red" ? "var(--red)" : "var(--text)";
  return (
    <div style={{ padding: 10, background: "var(--surface-1)", border: "1px solid var(--line)", borderRadius: 4 }}>
      <div className="sd-mono" style={{ fontSize: 9, letterSpacing: 0.5, color: "var(--text-3)" }}>{label}</div>
      <div className="sd-mono" style={{ fontSize: 14, fontWeight: 600, color, marginTop: 2 }}>{value}</div>
    </div>
  );
}
