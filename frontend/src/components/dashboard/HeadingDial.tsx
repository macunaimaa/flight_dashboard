interface Props { heading: number | null; size?: number; }

const CARDINALS = [
  { label: "N", deg: 0 }, { label: "E", deg: 90 },
  { label: "S", deg: 180 }, { label: "W", deg: 270 },
];

export function HeadingDial({ heading, size = 96 }: Props) {
  const h = heading ?? 0;
  const r = size / 2;
  return (
    <div style={{ width: size, height: size, position: "relative", flex: "0 0 auto" }}>
      <svg viewBox={`0 0 ${size} ${size}`} style={{ width: size, height: size }}>
        <circle cx={r} cy={r} r={r - 2} fill="var(--surface-1)" stroke="var(--line-2)" strokeWidth="1" />
        {/* tick marks every 30° */}
        {Array.from({ length: 12 }, (_, i) => {
          const a = (i * 30 * Math.PI) / 180;
          const x1 = r + Math.sin(a) * (r - 4);
          const y1 = r - Math.cos(a) * (r - 4);
          const x2 = r + Math.sin(a) * (r - (i % 3 === 0 ? 10 : 7));
          const y2 = r - Math.cos(a) * (r - (i % 3 === 0 ? 10 : 7));
          return <line key={i} x1={x1} y1={y1} x2={x2} y2={y2} stroke="var(--text-3)" strokeWidth="1" />;
        })}
        {CARDINALS.map(c => {
          const a = (c.deg * Math.PI) / 180;
          const x = r + Math.sin(a) * (r - 18);
          const y = r - Math.cos(a) * (r - 18) + 3;
          return (
            <text key={c.label} x={x} y={y} textAnchor="middle" fontSize="9" fontWeight="600"
              fill={c.label === "N" ? "var(--accent)" : "var(--text-3)"} fontFamily="IBM Plex Mono">{c.label}</text>
          );
        })}
        {/* aircraft pointer */}
        <g transform={`rotate(${h} ${r} ${r})`}>
          <path d={`M ${r} ${r - (r - 14)} L ${r + 6} ${r + 4} L ${r} ${r} L ${r - 6} ${r + 4} Z`}
            fill="var(--accent)" stroke="var(--bg)" strokeWidth="0.6" />
        </g>
        <circle cx={r} cy={r} r="2" fill="var(--text)" />
      </svg>
      <div className="sd-mono" style={{
        position: "absolute", inset: 0, display: "flex", alignItems: "flex-end",
        justifyContent: "center", paddingBottom: 6,
        fontSize: 11, color: "var(--text-2)", pointerEvents: "none",
      }}>{Math.round(h)}°</div>
    </div>
  );
}
