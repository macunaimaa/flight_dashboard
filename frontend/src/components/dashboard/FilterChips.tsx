import type { AircraftState } from "../../api/types";
import { deriveStatus } from "./derive";

export type ChipId = "all" | "longhaul" | "widebody" | "climbing" | "descending" | "delayed" | "emergency";

const CHIPS: { id: ChipId; label: string }[] = [
  { id: "all",        label: "All" },
  { id: "climbing",   label: "Climbing" },
  { id: "descending", label: "Descending" },
  { id: "longhaul",   label: "Long-haul" },
  { id: "widebody",   label: "Wide-body" },
  { id: "emergency",  label: "Alert" },
];

interface Props {
  active: ChipId;
  onChange: (id: ChipId) => void;
  counts: Record<ChipId, number>;
}

export function FilterChips({ active, onChange, counts }: Props) {
  return (
    <div style={{ display: "flex", gap: 6, flexWrap: "wrap", padding: "8px 12px", borderBottom: "1px solid var(--line)" }}>
      {CHIPS.map(c => {
        const isActive = c.id === active;
        const isAlert = c.id === "emergency" && counts[c.id] > 0;
        return (
          <button key={c.id} onClick={() => onChange(c.id)} className="sd-mono" style={{
            fontSize: 10.5, fontWeight: 600, letterSpacing: 0.4,
            padding: "4px 10px", borderRadius: 3,
            background: isActive ? "var(--surface-3)" : "transparent",
            border: `1px solid ${isActive ? "var(--accent)" : "var(--line-2)"}`,
            color: isAlert ? "var(--red)" : isActive ? "var(--accent)" : "var(--text-2)",
            cursor: "pointer", display: "inline-flex", alignItems: "center", gap: 6,
          }}>
            {c.label.toUpperCase()}
            <span style={{ color: "var(--text-3)", fontWeight: 500 }}>{counts[c.id]}</span>
          </button>
        );
      })}
    </div>
  );
}

export function filterByChip(list: AircraftState[], chip: ChipId): AircraftState[] {
  if (chip === "all") return list;
  return list.filter(a => {
    const s = deriveStatus(a);
    if (chip === "climbing")   return s === "CLIMBING";
    if (chip === "descending") return s === "DESCENDING";
    if (chip === "emergency")  return s === "EMERGENCY";
    if (chip === "longhaul") {
      const alt = (a.baroAltitude ?? a.geoAltitude ?? 0) * 3.28084;
      const spd = (a.velocity ?? 0) * 1.94384;
      return alt > 35000 && spd > 450;
    }
    if (chip === "widebody") {
      // We don't have model data; approximate using callsign of major widebody operators
      const cs = (a.callsign || "").toUpperCase().slice(0, 3);
      return ["UAL","DAL","AAL","BAW","DLH","AFR","KLM","UAE","QTR","SIA","CPA","JAL","ANA","KAL","QFA","VIR"].includes(cs);
    }
    return true;
  });
}

export function chipCounts(list: AircraftState[]): Record<ChipId, number> {
  const out: Record<ChipId, number> = {
    all: list.length, climbing: 0, descending: 0, emergency: 0, longhaul: 0, widebody: 0, delayed: 0,
  };
  for (const a of list) {
    const s = deriveStatus(a);
    if (s === "CLIMBING")   out.climbing++;
    if (s === "DESCENDING") out.descending++;
    if (s === "EMERGENCY")  out.emergency++;
    const alt = (a.baroAltitude ?? a.geoAltitude ?? 0) * 3.28084;
    const spd = (a.velocity ?? 0) * 1.94384;
    if (alt > 35000 && spd > 450) out.longhaul++;
    const cs = (a.callsign || "").toUpperCase().slice(0, 3);
    if (["UAL","DAL","AAL","BAW","DLH","AFR","KLM","UAE","QTR","SIA","CPA","JAL","ANA","KAL","QFA","VIR"].includes(cs)) out.widebody++;
  }
  return out;
}
