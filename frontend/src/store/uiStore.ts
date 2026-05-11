import { create } from "zustand";

interface UIState {
  sidebarOpen: boolean;
  detailPanelOpen: boolean;
  weatherPanelOpen: boolean;
  newsPanelOpen: boolean;
  connectionStatus: "connecting" | "connected" | "disconnected";

  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  setDetailPanelOpen: (open: boolean) => void;
  setConnectionStatus: (status: UIState["connectionStatus"]) => void;
  toggleWeatherPanel: () => void;
  toggleNewsPanel: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: true,
  detailPanelOpen: false,
  weatherPanelOpen: false,
  newsPanelOpen: false,
  connectionStatus: "disconnected",

  toggleSidebar: () => set((s) => ({ sidebarOpen: !s.sidebarOpen })),
  setSidebarOpen: (open) => set({ sidebarOpen: open }),
  setDetailPanelOpen: (open) => set({ detailPanelOpen: open }),
  setConnectionStatus: (status) => set({ connectionStatus: status }),
  toggleWeatherPanel: () => set((s) => ({ weatherPanelOpen: !s.weatherPanelOpen })),
  toggleNewsPanel: () => set((s) => ({ newsPanelOpen: !s.newsPanelOpen })),
}));
