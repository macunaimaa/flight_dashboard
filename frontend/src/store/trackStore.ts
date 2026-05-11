import { create } from "zustand";
import type { TrackPoint } from "../api/types";
import { getTrack } from "../api/aircraft";

interface TrackStore {
  points: TrackPoint[];
  loading: boolean;
  error: string | null;
  icao24: string | null;
  fetchTrack: (icao24: string) => Promise<void>;
  clear: () => void;
}

export const useTrackStore = create<TrackStore>((set) => ({
  points: [],
  loading: false,
  error: null,
  icao24: null,

  fetchTrack: async (icao24: string) => {
    set({ loading: true, error: null, icao24, points: [] });
    try {
      const resp = await getTrack(icao24);
      set({ points: resp.data, loading: false });
    } catch (e) {
      set({ error: e instanceof Error ? e.message : "Failed to fetch track", loading: false });
    }
  },

  clear: () => {
    set({ points: [], loading: false, error: null, icao24: null });
  },
}));
