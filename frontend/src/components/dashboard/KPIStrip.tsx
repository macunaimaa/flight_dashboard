import { useMemo } from "react";
import type { AircraftState } from "../../api/types";
import { deriveStatus } from "./derive";

interface Props { aircraft: AircraftState[]; }

interface KPI { label: string; value: string; unit?: string; tone?: "default" | "warn" | "alert"; sub?: string; }

function metersToFt(m: number | null): number | null { return m === null ? null : m * 3.28084; }
function msToKts(v: number | null): number | null { return v === null ? null : v * 1.94384; }

export function KPIStrip({ aircraft }: Props) {
  const kpis = useMemo<KPI[]>(() => {
    const airborne = aircraft.filter(a => !a.onGround);
    const alts = airborne.map(a => metersToFt(a.baroAltitude ?? a.geoAltitude)).filter((v): v is number => v !== null);
    const spds = airborne.map(a => msToKts(a.velocity)).filter((v): v is number => v !== null);
    const avgAlt = alts.length ? Math.round(alts.reduce((s, v) => s + v, 0) / alts.length) : 0;
    const avgSpd = spds.length ? Math.round(spds.reduce((s, v) => s + v, 0) / spds.length) : 0;
    const approach = airborne.filter(a => {
      const alt = metersToFt(a.baroAltitude ?? a.geoAltitude);
      const vr  = a.verticalRate;
      return alt !== null && alt < 10000 && vr !== null && vr * 196.85 < -200;
    }).length;
    const emergencies = aircraft.filter(a => deriveStatus(a) === "EMERGENCY").length;
    const ground = aircraft.filter(a => a.onGround).length;

    return [
      { label: "AIRCRAFT TRACKED", value: aircraft.length.toLocaleString(), sub: `${airborne.length} airborne · ${ground} ground` },
      { label: "AVG ALTITUDE",    value: avgAlt.toLocaleString(), unit: "ft" },
      { label: "AVG GROUND SPEED",value: avgSpd.toLocaleString(), unit: "kts" },
      { label: "IN APPROACH",     value: String(approach), sub: "<10k ft · descending" },
      { label: "ALERTS",          value: String(emergencies), tone: emergencies > 0 ? "alert" : "default", sub: emergencies ? "active squawks" : "none" },
    ];
  }, [aircraft]);

  return (
    <div style={{
      display: "grid", gridTemplateColumns: `repeat(${kpis.length}, 1fr)`, gap: 1,
      background: "var(--line)", borderBottom: "1px solid var(--line)",
    }}>
      {kpis.map(k => {
        const accent = k.tone === "alert" ? "var(--red)" : "var(--accent)";
        return (
          <div key={k.label} style={{ background: "var(--surface-0)", padding: "10px 16px", display: "flex", flexDirection: "column", gap: 4 }}>
            <div className="sd-mono" style={{ fontSize: 9.5, letterSpacing: 0.6, color: "var(--text-3)" }}>{k.label}</div>
            <div style={{ display: "flex", alignItems: "baseline", gap: 6 }}>
              <span className="sd-mono" style={{ fontSize: 22, fontWeight: 600, color: accent }}>{k.value}</span>
              {k.unit && <span className="sd-mono" style={{ fontSize: 11, color: "var(--text-3)" }}>{k.unit}</span>}
            </div>
            {k.sub && <div style={{ fontSize: 10.5, color: "var(--text-3)" }}>{k.sub}</div>}
          </div>
        );
      })}
    </div>
  );
}
