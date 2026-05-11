import { Cartesian3, Math as CesiumMath, Color } from "cesium";
import type { AircraftState } from "../api/types";

export function aircraftPosition(state: AircraftState): Cartesian3 {
  const [lon, lat] = state.location.coordinates;
  const alt = state.baroAltitude ?? 0;
  return Cartesian3.fromDegrees(lon, lat, alt);
}

export function aircraftRotation(state: AircraftState): number {
  return CesiumMath.toRadians(-(state.trueTrack ?? 0));
}

export function aircraftColor(state: AircraftState): Color {
  if (state.onGround) return Color.GRAY;
  if (state.baroAltitude && state.baroAltitude > 10000) return Color.CYAN;
  if (state.baroAltitude && state.baroAltitude > 5000)
    return Color.fromCssColorString("#4FC3F7");
  return Color.fromCssColorString("#81D4FA");
}

export function aircraftLabel(state: AircraftState): string {
  return state.callsign?.trim() || state.icao24.toUpperCase();
}
