import { create } from "zustand";
import type { AircraftState } from "../api/types";

interface AircraftStore {
  aircraft: Map<string, AircraftState>;
  selectedIcao24: string | null;
  // Derived state stored explicitly to avoid re-render loops
  aircraftList: AircraftState[];

  applyUpdates: (updates: AircraftState[]) => void;
  applyRemovals: (icao24s: string[]) => void;
  selectAircraft: (icao24: string | null) => void;
  clear: () => void;
}

function buildSortedList(map: Map<string, AircraftState>): AircraftState[] {
  return Array.from(map.values()).sort((a, b) =>
    (a.callsign || a.icao24).localeCompare(b.callsign || b.icao24)
  );
}

export const useAircraftStore = create<AircraftStore>((set) => ({
  aircraft: new Map(),
  selectedIcao24: null,
  aircraftList: [],

  applyUpdates: (updates: AircraftState[]) => {
    set((state) => {
      const newMap = new Map(state.aircraft);
      for (const update of updates) {
        newMap.set(update.icao24, update);
      }
      return { aircraft: newMap, aircraftList: buildSortedList(newMap) };
    });
  },

  applyRemovals: (icao24s: string[]) => {
    set((state) => {
      const newMap = new Map(state.aircraft);
      for (const icao of icao24s) {
        newMap.delete(icao);
      }
      const selectedIcao24 =
        state.selectedIcao24 && icao24s.includes(state.selectedIcao24)
          ? null
          : state.selectedIcao24;
      return {
        aircraft: newMap,
        aircraftList: buildSortedList(newMap),
        selectedIcao24,
      };
    });
  },

  selectAircraft: (icao24: string | null) => {
    set({ selectedIcao24: icao24 });
  },

  clear: () => {
    set({ aircraft: new Map(), selectedIcao24: null, aircraftList: [] });
  },
}));
