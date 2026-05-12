import { useMemo } from "react";
import { useAircraftStore } from "../../store/aircraftStore";
import { useUIStore } from "../../store/uiStore";
import { deriveStatus, deriveAirline } from "../dashboard/derive";
import { FilterChips, filterByChip, chipCounts } from "../dashboard/FilterChips";
import { StatusPill } from "../dashboard/StatusPill";
import { formatAltitude, formatSpeed } from "../../utils/formatting";

export function AircraftList() {
  const aircraft = useAircraftStore(s => s.aircraftList);
  const selectedId = useAircraftStore(s => s.selectedIcao24);
  const select = useAircraftStore(s => s.selectAircraft);
  const chip = useUIStore(s => s.filterChip);
  const setChip = useUIStore(s => s.setFilterChip);

  const filtered = useMemo(() => filterByChip(aircraft, chip), [aircraft, chip]);
  const counts = useMemo(() => chipCounts(aircraft), [aircraft]);

  // Sort: emergencies first, then by altitude descending
  const rows = useMemo(() => [...filtered].sort((a, b) => {
    const sa = deriveStatus(a), sb = deriveStatus(b);
    if (sa === "EMERGENCY" && sb !== "EMERGENCY") return -1;
    if (sb === "EMERGENCY" && sa !== "EMERGENCY") return 1;
    return (b.baroAltitude ?? 0) - (a.baroAltitude ?? 0);
  }), [filtered]);

  return (
    <aside style={{ display: "flex", flexDirection: "column", background: "var(--surface-0)",
      borderRight: "1px solid var(--line)", minWidth: 0 }}>
      <div style={{ padding: "12px 14px 8px", borderBottom: "1px solid var(--line)" }}>
        <div className="sd-mono" style={{ fontSize: 9.5, letterSpacing: 0.6, color: "var(--text-3)" }}>FLEET</div>
        <div style={{ display: "flex", alignItems: "baseline", gap: 8, marginTop: 2 }}>
          <span style={{ fontSize: 18, fontWeight: 600, color: "var(--text)" }}>{filtered.length}</span>
          <span className="sd-mono" style={{ fontSize: 11, color: "var(--text-3)" }}>of {aircraft.length}</span>
        </div>
      </div>
      <FilterChips active={chip} onChange={setChip} counts={counts} />
      <div style={{ flex: 1, overflow: "auto" }}>
        {rows.map(a => {
          const s = deriveStatus(a);
          const isSel = a.icao24 === selectedId;
          return (
            <div key={a.icao24} onClick={() => select(a.icao24)} style={{
              padding: "10px 14px", cursor: "pointer", borderBottom: "1px solid var(--line)",
              background: isSel ? "var(--surface-2)" : "transparent",
              borderLeft: `3px solid ${isSel ? "var(--accent)" : "transparent"}`,
            }}>
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <div className="sd-mono" style={{ fontWeight: 600, color: isSel ? "var(--accent)" : "var(--text)" }}>
                  {a.callsign || a.icao24}
                </div>
                <StatusPill status={s} />
              </div>
              <div style={{ fontSize: 11, color: "var(--text-3)", marginTop: 2, display: "flex", justifyContent: "space-between" }}>
                <span>{deriveAirline(a.callsign)}</span>
                <span className="sd-mono">{a.originCountry}</span>
              </div>
              <div style={{ display: "flex", gap: 12, marginTop: 4, fontSize: 10.5 }} className="sd-mono">
                <span style={{ color: "var(--text-2)" }}>{formatAltitude(a.baroAltitude)}</span>
                <span style={{ color: "var(--text-3)" }}>{formatSpeed(a.velocity)}</span>
                {a.squawk && <span style={{ color: s === "EMERGENCY" ? "var(--red)" : "var(--text-3)", marginLeft: "auto" }}>SQK {a.squawk}</span>}
              </div>
            </div>
          );
        })}
        {rows.length === 0 && (
          <div style={{ padding: 32, textAlign: "center", color: "var(--text-3)", fontSize: 12 }}>
            No aircraft match this filter.
          </div>
        )}
      </div>
    </aside>
  );
}
