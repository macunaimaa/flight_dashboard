import type { CSSProperties } from "react";
import type { FlightStatus } from "./derive";

interface Props { status: FlightStatus; }

const TONE: Record<FlightStatus, { color: string; bg: string; border: string }> = {
  CRUISE:     { color: "var(--text-2)", bg: "var(--surface-2)", border: "var(--line-2)" },
  GROUND:     { color: "var(--text-3)", bg: "var(--surface-2)", border: "var(--line-2)" },
  CLIMBING:   { color: "var(--mint)",   bg: "color-mix(in oklab, var(--mint) 14%, transparent)", border: "color-mix(in oklab, var(--mint) 35%, transparent)" },
  DESCENDING: { color: "var(--cyan)",   bg: "color-mix(in oklab, var(--cyan) 14%, transparent)", border: "color-mix(in oklab, var(--cyan) 35%, transparent)" },
  EMERGENCY:  { color: "var(--red)",    bg: "color-mix(in oklab, var(--red) 14%, transparent)",  border: "color-mix(in oklab, var(--red) 45%, transparent)" },
};

export function StatusPill({ status }: Props) {
  const t = TONE[status];
  const style: CSSProperties = {
    fontSize: 10, fontWeight: 600, letterSpacing: 0.5,
    padding: "3px 8px", borderRadius: 3, display: "inline-flex", alignItems: "center", gap: 4,
    color: t.color, background: t.bg, border: `1px solid ${t.border}`,
  };
  return (
    <span className="sd-mono" style={style}>
      {status === "EMERGENCY" && <span className="sd-blink">●</span>}
      {status}
    </span>
  );
}
