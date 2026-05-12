/* Derive enthusiast/spotter fields the backend doesn't supply.
   Frontend-only — replace with backend-enriched values when available. */

import type { AircraftState } from "../../api/types";

export type FlightStatus = "CRUISE" | "CLIMBING" | "DESCENDING" | "EMERGENCY" | "GROUND";

const EMERGENCY_SQUAWKS = new Set(["7500", "7600", "7700"]);

const AIRLINE_PREFIX: Record<string, string> = {
  BA: "British Airways", UA: "United", AA: "American", DL: "Delta", AC: "Air Canada",
  LH: "Lufthansa", AF: "Air France", KL: "KLM", IB: "Iberia", LX: "Swiss",
  EK: "Emirates", QR: "Qatar", EY: "Etihad", TK: "Turkish", SQ: "Singapore",
  CX: "Cathay Pacific", JL: "JAL", NH: "ANA", KE: "Korean Air", QF: "Qantas",
  AM: "Aeromexico", TP: "TAP", VS: "Virgin Atlantic", AY: "Finnair", TG: "Thai",
  EZY: "easyJet", RYR: "Ryanair",
};

export function deriveStatus(ac: AircraftState): FlightStatus {
  if (ac.squawk && EMERGENCY_SQUAWKS.has(ac.squawk)) return "EMERGENCY";
  if (ac.onGround) return "GROUND";
  const vr = ac.verticalRate; // m/s
  if (vr === null) return "CRUISE";
  const fpm = vr * 196.85;
  if (fpm > 200) return "CLIMBING";
  if (fpm < -200) return "DESCENDING";
  return "CRUISE";
}

export function deriveAirline(callsign: string): string {
  const cs = (callsign || "").trim().toUpperCase();
  // Try 3-letter ICAO prefix then 2-letter IATA
  for (const len of [3, 2]) {
    const p = cs.slice(0, len);
    if (AIRLINE_PREFIX[p]) return AIRLINE_PREFIX[p];
  }
  return "—";
}

/** Synthetic registration placeholder. Replace with backend value when available. */
export function deriveRegistration(icao24: string): string {
  // Real reg requires a lookup (hexdb.io etc). Show ICAO24 hex as fallback.
  return icao24.toUpperCase();
}

export function statusColor(s: FlightStatus): string {
  switch (s) {
    case "EMERGENCY":  return "var(--red)";
    case "CLIMBING":   return "var(--mint)";
    case "DESCENDING": return "var(--cyan)";
    case "GROUND":     return "var(--text-4)";
    default:           return "var(--text-3)";
  }
}

/** Long-haul heuristic: aircraft model includes wide-body type code. */
const WIDEBODY = /A3[5678]0|B77|B78|B74/i;
export function isWidebody(model: string | undefined): boolean {
  return !!model && WIDEBODY.test(model);
}
