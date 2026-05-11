export interface AircraftState {
  icao24: string;
  callsign: string;
  originCountry: string;
  location: {
    type: string;
    coordinates: [number, number]; // [lon, lat]
  };
  baroAltitude: number | null;
  geoAltitude: number | null;
  velocity: number | null;
  trueTrack: number | null;
  verticalRate: number | null;
  onGround: boolean;
  squawk: string | null;
  sourceTimestamp: string;
  updatedAt: string;
}

export interface APIResponse<T> {
  data: T;
  meta?: {
    total: number;
    limit: number;
    offset: number;
    timestamp: string;
  };
  error?: string;
}

export interface LoginRequest {
  tenantId: string;
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  expiresAt: string;
  tenantId: string;
  role: string;
}

export interface TrackPoint {
  icao24: string;
  location: {
    type: string;
    coordinates: [number, number];
  };
  altitude: number | null;
  velocity: number | null;
  heading: number | null;
  verticalRate: number | null;
  onGround: boolean;
  timestamp: string;
}

export interface WSMessage {
  type: string;
  payload: unknown;
}
