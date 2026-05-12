import { useMemo } from "react";
import type { AircraftState } from "../../api/types";
import { deriveStatus } from "./derive";
import { formatAltitude, timeAgo } from "../../utils/formatting";

interface Event { id: string; kind: "DEP" | "ARR" | "SQK" | "WX" | "CLB" | "DSC"; text: string; time: string; tone?: "alert"; }

interface Props { aircraft: AircraftState[]; weatherAlertCount?: number; }

/* Derives a synthetic event stream from current AircraftState observations.
   In production, replace with backend-pushed events (ws topic: /events). */
export function EventTicker({ aircraft, weatherAlertCount = 0 }: Props) {
  const events = useMemo<Event[]>(() => {
    const ev: Event[] = [];
    for (const a of aircraft.slice(0, 80)) {
      const s = deriveStatus(a);
      if (s === "EMERGENCY") {
        ev.push({ id: `sqk-${a.icao24}`, kind: "SQK", tone: "alert",
          text: `${a.callsign || a.icao24} squawking ${a.squawk}`, time: a.updatedAt });
      } else if (s === "CLIMBING" && (a.verticalRate ?? 0) * 196.85 > 1500) {
        ev.push({ id: `clb-${a.icao24}`, kind: "CLB",
          text: `${a.callsign} climbing through ${formatAltitude(a.baroAltitude)}`, time: a.updatedAt });
      } else if (s === "DESCENDING" && (a.baroAltitude ?? 99999) * 3.28084 < 8000) {
        ev.push({ id: `dsc-${a.icao24}`, kind: "ARR",
          text: `${a.callsign} on final approach`, time: a.updatedAt });
      }
    }
    if (weatherAlertCount > 0) {
      ev.unshift({ id: "wx-0", kind: "WX", tone: "alert",
        text: `${weatherAlertCount} active weather advisory · expand panel for details`, time: new Date().toISOString() });
    }
    return ev.slice(0, 24);
  }, [aircraft, weatherAlertCount]);

  if (!events.length) {
    return (
      <div style={{ height: 32, display: "flex", alignItems: "center", padding: "0 16px",
        background: "var(--surface-0)", borderTop: "1px solid var(--line)", color: "var(--text-3)", fontSize: 11 }}>
        <span className="sd-mono">● LIVE</span>
        <span style={{ marginLeft: 12 }}>No notable events</span>
      </div>
    );
  }

  return (
    <div style={{ height: 32, display: "flex", alignItems: "center", background: "var(--surface-0)",
      borderTop: "1px solid var(--line)", overflow: "hidden", position: "relative" }}>
      <div className="sd-mono sd-blink" style={{ flex: "0 0 auto", padding: "0 14px", color: "var(--accent)",
        fontSize: 10.5, fontWeight: 700, letterSpacing: 0.6, borderRight: "1px solid var(--line)" }}>● LIVE</div>
      <div style={{ flex: 1, overflow: "hidden", position: "relative", maskImage: "linear-gradient(90deg, transparent 0, #000 30px, #000 calc(100% - 60px), transparent 100%)" }}>
        <div className="sd-ticker-track" style={{ display: "inline-flex", gap: 28, whiteSpace: "nowrap", paddingLeft: 20 }}>
          {[...events, ...events].map((e, i) => (
            <span key={`${e.id}-${i}`} style={{ display: "inline-flex", alignItems: "center", gap: 8, fontSize: 11.5 }}>
              <span className="sd-mono" style={{
                fontSize: 9.5, fontWeight: 700, padding: "2px 6px", borderRadius: 2,
                background: e.tone === "alert" ? "color-mix(in oklab, var(--red) 18%, transparent)" : "var(--surface-2)",
                color: e.tone === "alert" ? "var(--red)" : "var(--text-2)",
              }}>{e.kind}</span>
              <span style={{ color: "var(--text)" }}>{e.text}</span>
              <span className="sd-mono" style={{ color: "var(--text-4)", fontSize: 10 }}>{timeAgo(e.time)}</span>
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}
