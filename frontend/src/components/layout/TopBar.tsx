import { useState } from "react";
import { useAircraftStore } from "../../store/aircraftStore";
import { useUIStore } from "../../store/uiStore";
import { useAuthStore } from "../../store/authStore";

export function TopBar() {
  const aircraft = useAircraftStore(s => s.aircraftList);
  const select = useAircraftStore(s => s.selectAircraft);
  const connectionStatus = useUIStore(s => s.connectionStatus);
  const mapMode = useUIStore(s => s.mapMode);
  const setMapMode = useUIStore(s => s.setMapMode);
  const toggleWx = useUIStore(s => s.toggleWeatherPanel);
  const toggleNews = useUIStore(s => s.toggleNewsPanel);
  const logout = useAuthStore(s => s.logout);
  const tenantId = useAuthStore(s => s.tenantId);

  const [q, setQ] = useState("");
  const matches = q.trim().length >= 2
    ? aircraft.filter(a => {
        const s = q.toLowerCase();
        return (a.callsign || "").toLowerCase().includes(s)
          || a.icao24.toLowerCase().includes(s)
          || (a.originCountry || "").toLowerCase().includes(s);
      }).slice(0, 8)
    : [];

  const dot = connectionStatus === "connected" ? "var(--mint)"
    : connectionStatus === "connecting" ? "var(--amber)" : "var(--red)";

  return (
    <header style={{
      height: 52, display: "flex", alignItems: "center", gap: 16, padding: "0 16px",
      background: "var(--surface-0)", borderBottom: "1px solid var(--line)", position: "relative", zIndex: 20,
    }}>
      {/* Brand */}
      <div style={{ display: "flex", alignItems: "center", gap: 10, flex: "0 0 auto" }}>
        <div style={{ width: 28, height: 28, borderRadius: 4, background: "var(--accent)",
          display: "grid", placeItems: "center", color: "var(--bg)", fontWeight: 800, fontSize: 13 }}>S</div>
        <div>
          <div style={{ fontWeight: 700, fontSize: 13, letterSpacing: 0.6 }}>SKYDECK</div>
          <div className="sd-mono" style={{ fontSize: 9, color: "var(--text-3)", letterSpacing: 0.4 }}>
            FLIGHT OPS · {tenantId || "—"}
          </div>
        </div>
      </div>

      {/* Search */}
      <div style={{ flex: 1, maxWidth: 480, position: "relative" }}>
        <input value={q} onChange={e => setQ(e.target.value)} placeholder="Search callsign, ICAO24, country…"
          style={{
            width: "100%", background: "var(--surface-1)", border: "1px solid var(--line-2)",
            color: "var(--text)", padding: "7px 12px", borderRadius: 4, fontSize: 12, outline: "none",
            fontFamily: "IBM Plex Sans",
          }} />
        {matches.length > 0 && (
          <div style={{
            position: "absolute", top: "calc(100% + 4px)", left: 0, right: 0,
            background: "var(--surface-1)", border: "1px solid var(--line-2)", borderRadius: 4,
            boxShadow: "0 12px 32px rgba(0,0,0,0.5)", overflow: "hidden", zIndex: 30,
          }}>
            {matches.map(a => (
              <div key={a.icao24} onClick={() => { select(a.icao24); setQ(""); }}
                style={{ padding: "8px 12px", cursor: "pointer", display: "flex", justifyContent: "space-between", borderBottom: "1px solid var(--line)" }}>
                <span className="sd-mono" style={{ color: "var(--text)" }}>{a.callsign || a.icao24}</span>
                <span className="sd-mono" style={{ color: "var(--text-3)", fontSize: 11 }}>{a.originCountry}</span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Map mode toggle */}
      <div style={{ display: "flex", border: "1px solid var(--line-2)", borderRadius: 4, overflow: "hidden" }}>
        {(["cesium","2d"] as const).map(m => {
          const active = mapMode === m;
          return (
            <button key={m} onClick={() => setMapMode(m)} className="sd-mono" style={{
              padding: "6px 12px", fontSize: 10.5, fontWeight: 600, letterSpacing: 0.5,
              background: active ? "var(--surface-3)" : "transparent",
              color: active ? "var(--accent)" : "var(--text-3)",
              border: "none", cursor: "pointer",
            }}>{m === "cesium" ? "3D" : "RADAR"}</button>
          );
        })}
      </div>

      {/* Status pill */}
      <div className="sd-mono" style={{
        display: "flex", alignItems: "center", gap: 8, padding: "6px 10px",
        border: "1px solid var(--line-2)", borderRadius: 4, fontSize: 10.5, color: "var(--text-2)",
      }}>
        <span style={{ width: 7, height: 7, borderRadius: "50%", background: dot }} className={connectionStatus === "connected" ? "sd-blink" : ""} />
        {connectionStatus.toUpperCase()} · {aircraft.length} CONTACTS
      </div>

      {/* Action buttons */}
      <div style={{ display: "flex", gap: 6 }}>
        <button onClick={toggleWx} className="sd-mono" style={btnStyle}>WX</button>
        <button onClick={toggleNews} className="sd-mono" style={btnStyle}>NEWS</button>
        <button onClick={logout} className="sd-mono" style={{ ...btnStyle, color: "var(--text-3)" }}>↪</button>
      </div>
    </header>
  );
}

const btnStyle: React.CSSProperties = {
  background: "var(--surface-1)", border: "1px solid var(--line-2)",
  color: "var(--text-2)", padding: "6px 12px", borderRadius: 4,
  fontSize: 10.5, fontWeight: 600, letterSpacing: 0.5, cursor: "pointer",
};
