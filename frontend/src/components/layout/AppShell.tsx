import { useRef } from "react";
import { TopBar } from "./TopBar";
import { Sidebar } from "./Sidebar";
import { CesiumGlobe } from "../../cesium/CesiumViewer";
import { AircraftInfoPanel } from "../aircraft/AircraftInfoPanel";
import { WeatherPanel } from "../weather/WeatherPanel";
import { NewsPanel } from "../news/NewsPanel";
import { useWebSocket } from "../../hooks/useWebSocket";
import { useUIStore } from "../../store/uiStore";
import { useAircraftStore } from "../../store/aircraftStore";

export function AppShell() {
  useWebSocket();

  const weatherPanelOpen = useUIStore((s) => s.weatherPanelOpen);
  const newsPanelOpen = useUIStore((s) => s.newsPanelOpen);
  const toggleWeatherPanel = useUIStore((s) => s.toggleWeatherPanel);
  const toggleNewsPanel = useUIStore((s) => s.toggleNewsPanel);
  const selectedIcao24 = useAircraftStore((s) => s.selectedIcao24);
  const aircraft = useAircraftStore((s) => s.aircraft);

  const viewportCenterRef = useRef({ lat: 0, lon: 0 });

  if (selectedIcao24) {
    const selected = aircraft.get(selectedIcao24);
    if (selected) {
      viewportCenterRef.current = {
        lon: selected.location.coordinates[0],
        lat: selected.location.coordinates[1],
      };
    }
  }

  return (
    <div style={styles.shell}>
      <TopBar />
      <div style={styles.content}>
        <Sidebar />
        <div style={styles.globe}>
          <CesiumGlobe />
          <AircraftInfoPanel />
          {weatherPanelOpen && (
            <WeatherPanel
              lat={viewportCenterRef.current.lat}
              lon={viewportCenterRef.current.lon}
              onClose={toggleWeatherPanel}
            />
          )}
          {newsPanelOpen && (
            <NewsPanel onClose={toggleNewsPanel} />
          )}
        </div>
      </div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  shell: {
    display: "flex",
    flexDirection: "column",
    height: "100vh",
    width: "100vw",
    overflow: "hidden",
    background: "#0a0e17",
  },
  content: {
    display: "flex",
    flex: 1,
    overflow: "hidden",
  },
  globe: {
    flex: 1,
    position: "relative",
    overflow: "hidden",
  },
};
