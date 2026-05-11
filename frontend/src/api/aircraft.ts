import { apiClient } from "./client";
import type { AircraftState, TrackPoint } from "./types";

export async function listAircraft(params?: {
  callsign?: string;
  country?: string;
  onGround?: boolean;
  limit?: number;
  offset?: number;
}) {
  const query = new URLSearchParams();
  if (params?.callsign) query.set("callsign", params.callsign);
  if (params?.country) query.set("country", params.country);
  if (params?.onGround !== undefined)
    query.set("on_ground", String(params.onGround));
  if (params?.limit) query.set("limit", String(params.limit));
  if (params?.offset) query.set("offset", String(params.offset));

  const qs = query.toString();
  return apiClient.get<AircraftState[]>(
    `/api/v1/aircraft${qs ? `?${qs}` : ""}`
  );
}

export async function getAircraft(icao24: string) {
  return apiClient.get<AircraftState>(`/api/v1/aircraft/${icao24}`);
}

export async function getAircraftInBBox(
  lamin: number,
  lomin: number,
  lamax: number,
  lomax: number,
  limit = 1000
) {
  return apiClient.get<AircraftState[]>(
    `/api/v1/aircraft/bbox?lamin=${lamin}&lomin=${lomin}&lamax=${lamax}&lomax=${lomax}&limit=${limit}`
  );
}

export async function getTrack(icao24: string, from?: string, to?: string) {
  const query = new URLSearchParams();
  if (from) query.set("from", from);
  if (to) query.set("to", to);
  const qs = query.toString();
  return apiClient.get<TrackPoint[]>(
    `/api/v1/aircraft/${icao24}/track${qs ? `?${qs}` : ""}`
  );
}
