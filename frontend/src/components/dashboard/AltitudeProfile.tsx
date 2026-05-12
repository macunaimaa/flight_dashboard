import type { AircraftState } from "../../api/types";

interface Props { aircraft: AircraftState; width?: number; height?: number; }

/* Altitude profile sparkline.
   With no historical track in scope, we render a synthesized profile anchored
   on the current altitude + vertical rate so it reads as a "where it is in the
   flight" indicator. Wire to backend track history once available. */
export function AltitudeProfile({ aircraft, width = 320, height = 64 }: Props) {
  const altFt = (aircraft.baroAltitude ?? aircraft.geoAltitude ?? 0) * 3.28084;
  const vrFpm = (aircraft.verticalRate ?? 0) * 196.85;
  const maxAlt = 42000;

  // Build 40 points: ascending profile, cruise plateau biased toward current,
  // then either continued cruise or descent depending on vertical rate sign.
  const N = 40;
  const pts: { x: number; y: number }[] = [];
  for (let i = 0; i < N; i++) {
    const t = i / (N - 1);
    let a: number;
    if (t < 0.18)       a = altFt * (t / 0.18) * 0.95;
    else if (t < 0.78)  a = altFt * (0.95 + 0.05 * Math.sin(t * 20));
    else {
      // tail: extrapolate using vertical rate sign
      const tail = (t - 0.78) / 0.22;
      a = vrFpm < -200 ? altFt * (1 - tail * 0.7) : altFt * (1 - tail * 0.1);
    }
    pts.push({ x: t * width, y: height - (a / maxAlt) * height });
  }

  const d = pts.map((p, i) => `${i === 0 ? "M" : "L"} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`).join(" ");
  const dFill = `${d} L ${width} ${height} L 0 ${height} Z`;
  const cursorX = pts[Math.floor(N * 0.5)].x;
  const cursorY = pts[Math.floor(N * 0.5)].y;

  return (
    <svg viewBox={`0 0 ${width} ${height}`} style={{ width: "100%", height, display: "block" }}>
      <defs>
        <linearGradient id="altGrad" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor="var(--accent)" stopOpacity="0.35" />
          <stop offset="100%" stopColor="var(--accent)" stopOpacity="0" />
        </linearGradient>
      </defs>
      {/* altitude grid: FL100/200/300/400 */}
      {[10000, 20000, 30000, 40000].map(fl => {
        const y = height - (fl / maxAlt) * height;
        return (
          <g key={fl}>
            <line x1="0" y1={y} x2={width} y2={y} stroke="var(--line)" strokeDasharray="2 3" strokeWidth="0.5" />
            <text x={width - 4} y={y - 2} textAnchor="end" fontSize="8" fill="var(--text-4)" fontFamily="IBM Plex Mono">FL{fl / 100}</text>
          </g>
        );
      })}
      <path d={dFill} fill="url(#altGrad)" />
      <path d={d} fill="none" stroke="var(--accent)" strokeWidth="1.5" />
      <line x1={cursorX} y1={0} x2={cursorX} y2={height} stroke="var(--accent)" strokeWidth="0.6" strokeDasharray="2 2" opacity="0.6" />
      <circle cx={cursorX} cy={cursorY} r="3" fill="var(--bg)" stroke="var(--accent)" strokeWidth="1.5" />
    </svg>
  );
}
