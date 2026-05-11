import type { AircraftState } from "../api/types";

export interface WSAircraftUpdate {
  type: "aircraft.update";
  payload: AircraftState[];
}

export interface WSAircraftRemove {
  type: "aircraft.remove";
  payload: string[];
}

export interface WSAlertTriggered {
  type: "alert.triggered";
  payload: {
    alertId: string;
    alertName: string;
    severity: string;
    icao24: string;
    callsign: string;
    triggeredAt: string;
  };
}

export interface WSSystemInfo {
  type: "system.info";
  payload: { message: string };
}

export type WSIncomingMessage =
  | WSAircraftUpdate
  | WSAircraftRemove
  | WSAlertTriggered
  | WSSystemInfo;
