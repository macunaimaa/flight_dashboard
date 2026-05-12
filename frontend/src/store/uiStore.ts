import { create } from "zustand";
import type { ChipId } from "../components/dashboard/FilterChips";

export type MapMode = "cesium" | "2d";
export type ConnectionStatus = "connected" | "connecting" | "disconnected";

interface UIState {
  weatherPanelOpen: boolean;
  newsPanelOpen: boolean;
  sidebarOpen: boolean;
  connectionStatus: ConnectionStatus;
  mapMode: MapMode;
  filterChip: ChipId;
  toggleWeatherPanel: () => void;
  toggleNewsPanel: () => void;
  toggleSidebar: () => void;
  setConnectionStatus: (s: ConnectionStatus) => void;
  setMapMode: (m: MapMode) => void;
  setFilterChip: (c: ChipId) => void;
}

export const useUIStore = create<UIState>((set) => ({
  weatherPanelOpen: false,
  newsPanelOpen: false,
  sidebarOpen: true,
  connectionStatus: "connecting",
  mapMode: "cesium",
  filterChip: "all",
  toggleWeatherPanel: () => set((s) => ({ weatherPanelOpen: !s.weatherPanelOpen })),
  toggleNewsPanel: () => set((s) => ({ newsPanelOpen: !s.newsPanelOpen })),
  toggleSidebar: () => set((s) => ({ sidebarOpen: !s.sidebarOpen })),
  setConnectionStatus: (s) => set({ connectionStatus: s }),
  setMapMode: (m) => set({ mapMode: m }),
  setFilterChip: (c) => set({ filterChip: c }),
}));
