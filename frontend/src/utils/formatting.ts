export function formatAltitude(meters: number | null): string {
  if (meters === null) return "N/A";
  return `${Math.round(meters * 3.28084).toLocaleString()} ft`;
}

export function formatSpeed(ms: number | null): string {
  if (ms === null) return "N/A";
  const knots = ms * 1.94384;
  return `${Math.round(knots)} kts`;
}

export function formatVerticalRate(ms: number | null): string {
  if (ms === null) return "N/A";
  const fpm = ms * 196.85;
  return `${fpm > 0 ? "+" : ""}${Math.round(fpm)} fpm`;
}

export function formatHeading(deg: number | null): string {
  if (deg === null) return "N/A";
  return `${Math.round(deg)}°`;
}

export function formatCoordinate(value: number, type: "lat" | "lon"): string {
  const abs = Math.abs(value);
  const dir = type === "lat" ? (value >= 0 ? "N" : "S") : value >= 0 ? "E" : "W";
  return `${abs.toFixed(4)}° ${dir}`;
}

export function timeAgo(timestamp: string): string {
  const diff = Date.now() - new Date(timestamp).getTime();
  const seconds = Math.floor(diff / 1000);
  if (seconds < 5) return "just now";
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.floor(seconds / 60);
  return `${minutes}m ago`;
}
