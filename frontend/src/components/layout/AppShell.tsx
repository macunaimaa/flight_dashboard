import { useEffect, useRef } from "react";
import { TopBar } from "./TopBar";
import { AircraftList } from "../aircraft/AircraftList";
import { AircraftInfoPanel } from "../aircraft/AircraftInfoPanel";
import { KPIStrip } from "../dashboard/KPIStrip";
import { EventTicker } from "../dashboard/EventTicker";
import { Map2D } from "../map2d/Map2D";
import { CesiumGlobe } from "../../cesium/CesiumViewer";
import { WeatherPanel } from "../weather/WeatherPanel";
import { NewsPanel } from "../news/NewsPanel";
import { useUIStore } from "../../store/uiStore";
import { useAircraftStore } from "../../store/aircraftStore";
import { useWebSocket } from "../../hooks/useWebSocket";

export function AppShell() {
  useWebSocket();

  const mapMode = useUIStore(s => s.mapMode);
  const wxOpen = useUIStore(s => s.weatherPanelOpen);
  const newsOpen = useUIStore(s => s.newsPanelOpen);
  const toggleWx = useUIStore(s => s.toggleWeatherPanel);
  const toggleNews = useUIStore(s => s.toggleNewsPanel);
  const aircraft = useAircraftStore(s => s.aircraftList);
  const selectedIcao24 = useAircraftStore(s => s.selectedIcao24);
  const aircraftMap = useAircraftStore(s => s.aircraft);

  const viewportCenterRef = useRef<{ lat: number; lon: number }>({ lat: 0, lon: 0 });

  if (selectedIcao24) {
    const selected = aircraftMap.get(selectedIcao24);
    if (selected) {
      viewportCenterRef.current = {
        lon: selected.location.coordinates[0],
        lat: selected.location.coordinates[1],
      };
    }
  }

  // Mobile detection
  useEffect(() => {
    const update = () => document.documentElement.dataset.skydeckSize =
      window.innerWidth < 760 ? "mobile" : "desktop";
    update();
    window.addEventListener("resize", update);
    return () => window.removeEventListener("resize", update);
  }, []);

  return (
    <div style={{
      display: "grid",
      gridTemplateRows: "auto auto 1fr auto",
      height: "100vh",
      background: "var(--bg)",
      color: "var(--text)",
    }}>
      <TopBar />
      <KPIStrip aircraft={aircraft} />
      <main style={{ display: "grid", gridTemplateColumns: "300px 1fr 360px", minHeight: 0, position: "relative" }}
        className="sd-main">
        <AircraftList />
        <section style={{ position: "relative", background: "var(--bg)" }}>
          {mapMode === "cesium" ? <CesiumGlobe /> : <Map2D />}
          {wxOpen && (
            <WeatherPanel
              lat={viewportCenterRef.current.lat}
              lon={viewportCenterRef.current.lon}
              onClose={toggleWx}
            />
          )}
          {newsOpen && <NewsPanel onClose={toggleNews} />}
        </section>
        <AircraftInfoPanel />
      </main>
      <EventTicker aircraft={aircraft} />

      {/* Mobile layout overrides */}
      <style>{`
        @media (max-width: 760px) {
          .sd-main { grid-template-columns: 1fr !important; }
          .sd-main > aside:first-child {
            position: fixed; left: 0; top: 0; bottom: 0; width: 84vw; max-width: 320px;
            z-index: 40; transform: translateX(-100%); transition: transform .2s ease;
          }
          .sd-main > aside:first-child[data-open="true"] { transform: translateX(0); }
          .sd-main > aside:last-child {
            position: fixed; left: 0; right: 0; bottom: 32px;
            width: 100%; max-height: 60vh; border-left: none;
            border-top: 1px solid var(--line); z-index: 30;
            border-top-left-radius: 12px; border-top-right-radius: 12px;
          }
        }
      `}</style>
    </div>
  );
}
