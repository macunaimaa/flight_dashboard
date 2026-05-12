import { useMemo, useRef } from "react";
import { useAircraftStore } from "../../store/aircraftStore";
import { deriveStatus } from "../dashboard/derive";
import { AIRPORTS, project } from "./airports";

/* 2D radar overlay — a togglable alternative to Cesium for low-bandwidth
   sessions and mobile. Uses an equirectangular projection on an SVG canvas.
   Continent silhouettes are intentionally schematic (no GeoJSON dependency)
   so the bundle stays light. Swap to TopoJSON if you want geographic fidelity. */

const CONTINENTS = [
  "M 80 200 L 140 180 L 220 200 L 280 240 L 300 300 L 280 360 L 220 400 L 160 420 L 120 380 L 100 340 L 80 280 Z",
  "M 280 460 L 340 460 L 360 520 L 340 600 L 320 680 L 280 740 L 260 700 L 250 620 L 260 540 Z",
  "M 440 200 L 560 220 L 560 280 L 520 320 L 460 320 L 440 280 Z",
  "M 460 340 L 560 340 L 580 420 L 560 520 L 520 620 L 480 640 L 440 560 L 440 440 Z",
  "M 560 200 L 780 180 L 880 240 L 900 320 L 860 380 L 800 420 L 720 420 L 640 380 L 580 320 L 560 260 Z",
  "M 740 440 L 840 440 L 860 500 L 800 540 L 740 520 Z",
  "M 820 620 L 920 620 L 940 680 L 900 720 L 840 720 L 820 680 Z",
];

export function Map2D() {
  const aircraft = useAircraftStore(s => s.aircraftList);
  const selectedId = useAircraftStore(s => s.selectedIcao24);
  const select = useAircraftStore(s => s.selectAircraft);
  const ref = useRef<HTMLDivElement>(null);

  // Project aircraft into pixel space, scaled to canvas
  const items = useMemo(() => aircraft.map(a => {
    const [lon, lat] = a.location.coordinates;
    const [px, py] = project(lon, lat);
    return { a, x: px * 1000, y: py * 500, status: deriveStatus(a) };
  }), [aircraft]);

  const selected = items.find(i => i.a.icao24 === selectedId);

  return (
    <div ref={ref} style={{ position: "absolute", inset: 0, background: "var(--bg)", overflow: "hidden" }}
      onClick={() => select(null)}>
      <svg viewBox="0 0 1000 500" preserveAspectRatio="xMidYMid slice"
        style={{ width: "100%", height: "100%", display: "block" }}>
        <defs>
          <radialGradient id="m2d-vig" cx="50%" cy="50%" r="65%">
            <stop offset="60%" stopColor="rgba(0,0,0,0)" />
            <stop offset="100%" stopColor="rgba(0,0,0,0.55)" />
          </radialGradient>
        </defs>
        {/* grid */}
        <g stroke="var(--line)" strokeWidth="0.4" opacity="0.7">
          {Array.from({ length: 11 }, (_, i) => (
            <line key={`v${i}`} x1={i * 100} y1="0" x2={i * 100} y2="500" />
          ))}
          {Array.from({ length: 6 }, (_, i) => (
            <line key={`h${i}`} x1="0" y1={i * 100} x2="1000" y2={i * 100} />
          ))}
        </g>
        {/* continents (schematic) — scaled down from 1000x800 source */}
        <g transform="scale(1, 0.625)" fill="var(--surface-1)" stroke="var(--line-2)" strokeWidth="0.8">
          {CONTINENTS.map((d, i) => <path key={i} d={d} />)}
        </g>
        {/* airport range rings */}
        {Object.entries(AIRPORTS).slice(0, 8).map(([code, ap]) => {
          const [px, py] = project(ap.coord[0], ap.coord[1]);
          const cx = px * 1000, cy = py * 500;
          return (
            <g key={code} opacity="0.35">
              {[8, 16, 24].map(r => (
                <circle key={r} cx={cx} cy={cy} r={r} fill="none"
                  stroke="var(--text-3)" strokeWidth="0.3" strokeDasharray="2 2" />
              ))}
              <rect x={cx - 1.2} y={cy - 1.2} width="2.4" height="2.4" fill="var(--text-3)" />
            </g>
          );
        })}
        {/* selected trail (synthesized straight line back along heading) */}
        {selected && (() => {
          const h = (selected.a.trueTrack ?? 0) * Math.PI / 180;
          const dx = -Math.sin(h) * 60, dy = Math.cos(h) * 60;
          return (
            <line x1={selected.x + dx} y1={selected.y + dy} x2={selected.x} y2={selected.y}
              stroke="var(--accent)" strokeWidth="1.2" strokeOpacity="0.7" />
          );
        })()}
        {/* aircraft */}
        {items.map(({ a, x, y, status }) => {
          const isSel = a.icao24 === selectedId;
          const color = status === "EMERGENCY" ? "var(--red)"
            : isSel ? "var(--accent)"
            : status === "CLIMBING" ? "var(--mint)"
            : status === "DESCENDING" ? "var(--cyan)" : "var(--text-2)";
          return (
            <g key={a.icao24} transform={`translate(${x} ${y}) rotate(${a.trueTrack ?? 0})`}
              style={{ cursor: "pointer" }}
              onClick={(e) => { e.stopPropagation(); select(a.icao24); }}>
              {isSel && <circle r="6" fill="none" stroke="var(--accent)" strokeWidth="0.8" className="sd-pulse" />}
              <path d="M 0 -3 L 2.4 2.4 L 0 1 L -2.4 2.4 Z" fill={color}
                stroke={isSel ? "var(--bg)" : "transparent"} strokeWidth="0.5" />
            </g>
          );
        })}
        {/* radar sweep */}
        <g className="sd-sweep" style={{ transformOrigin: "500px 250px" }} opacity="0.18">
          <defs>
            <linearGradient id="m2d-sweep" x1="0" y1="0" x2="1" y2="0">
              <stop offset="0%" stopColor="var(--accent)" stopOpacity="0" />
              <stop offset="100%" stopColor="var(--accent)" stopOpacity="0.6" />
            </linearGradient>
          </defs>
          <path d="M 500 250 L 1000 250 A 500 500 0 0 0 920 30 Z" fill="url(#m2d-sweep)" />
        </g>
        <rect x="0" y="0" width="1000" height="500" fill="url(#m2d-vig)" pointerEvents="none" />
      </svg>
    </div>
  );
}
