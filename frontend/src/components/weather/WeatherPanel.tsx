import { useState, useEffect, useCallback } from "react";
import { fetchWeather, weatherCodeToDescription, weatherCodeToIcon } from "../../api/weather";
import type { WeatherData } from "../../api/weather";

interface Props {
  lat: number;
  lon: number;
  onClose: () => void;
}

export function WeatherPanel({ lat, lon, onClose }: Props) {
  const [weather, setWeather] = useState<WeatherData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchWeather(lat, lon);
      setWeather(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load weather");
    } finally {
      setLoading(false);
    }
  }, [lat, lon]);

  useEffect(() => {
    load();
  }, [load]);

  const windDir = weather?.current.wind_direction_10m ?? 0;
  const windArrow = getWindArrow(windDir);

  return (
    <div style={styles.panel}>
      <div style={styles.header}>
        <div style={styles.headerLeft}>
          <span style={styles.iconText}>{weather ? weatherCodeToIcon(weather.current.weather_code) : "\u2601"}</span>
          <span style={styles.title}>Weather Intel</span>
        </div>
        <button onClick={onClose} style={styles.closeBtn}>&times;</button>
      </div>

      <div style={styles.coords}>
        {lat.toFixed(2)}{lat >= 0 ? "N" : "S"}, {lon.toFixed(2)}{lon >= 0 ? "E" : "W"}
      </div>

      {loading && <div style={styles.loading}>Fetching weather data...</div>}
      {error && <div style={styles.error}>{error}</div>}

      {weather && !loading && (
        <div style={styles.content}>
          <div style={styles.mainRow}>
            <span style={styles.temp}>{Math.round(weather.current.temperature_2m)}{"\u00B0"}C</span>
            <span style={styles.desc}>{weatherCodeToDescription(weather.current.weather_code)}</span>
          </div>

          <div style={styles.grid}>
            <MetricRow label="WIND" value={`${weather.current.wind_speed_10m} km/h ${windArrow}`} />
            <MetricRow label="WIND DIR" value={`${windDir}{"\u00B0"}`} />
            <MetricRow label="HUMIDITY" value={`${weather.current.relative_humidity_2m}%`} />
            <MetricRow label="PRESSURE" value={`${weather.current.surface_pressure} hPa`} />
            <MetricRow label="CLOUDS" value={`${weather.current.cloud_cover}%`} />
          </div>

          <button onClick={load} style={styles.refreshBtn}>Refresh</button>
        </div>
      )}
    </div>
  );
}

function MetricRow({ label, value }: { label: string; value: string }) {
  return (
    <div style={styles.metricRow}>
      <span style={styles.metricLabel}>{label}</span>
      <span style={styles.metricValue}>{value}</span>
    </div>
  );
}

function getWindArrow(deg: number): string {
  const arrows = ["\u2193", "\u2199", "\u2190", "\u2196", "\u2191", "\u2197", "\u2192", "\u2198"];
  const index = Math.round(deg / 45) % 8;
  return arrows[index];
}

const styles: Record<string, React.CSSProperties> = {
  panel: {
    position: "absolute",
    top: 56,
    right: 0,
    width: 280,
    height: "calc(100% - 56px)",
    background: "#111827",
    borderLeft: "1px solid #1f2937",
    zIndex: 70,
    display: "flex",
    flexDirection: "column",
    overflow: "hidden",
  },
  header: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    padding: "12px 16px",
    borderBottom: "1px solid #1f2937",
  },
  headerLeft: {
    display: "flex",
    alignItems: "center",
    gap: 8,
  },
  iconText: {
    fontSize: 20,
  },
  title: {
    color: "#e0e6ed",
    fontWeight: 700,
    fontSize: 13,
    textTransform: "uppercase" as const,
    letterSpacing: 1,
  },
  closeBtn: {
    background: "none",
    border: "none",
    color: "#6b7280",
    fontSize: 20,
    cursor: "pointer",
    padding: "0 4px",
    lineHeight: 1,
  },
  coords: {
    padding: "8px 16px",
    color: "#6b7280",
    fontSize: 11,
    fontFamily: "monospace",
    borderBottom: "1px solid #1f2937",
  },
  content: {
    flex: 1,
    overflow: "auto",
    padding: "12px 16px",
  },
  mainRow: {
    display: "flex",
    alignItems: "baseline",
    gap: 12,
    marginBottom: 16,
  },
  temp: {
    color: "#e0e6ed",
    fontSize: 32,
    fontWeight: 700,
    fontFamily: "monospace",
  },
  desc: {
    color: "#9ca3af",
    fontSize: 13,
  },
  grid: {
    display: "flex",
    flexDirection: "column" as const,
    gap: 2,
  },
  metricRow: {
    display: "flex",
    justifyContent: "space-between",
    padding: "6px 0",
    borderBottom: "1px solid #1f293733",
  },
  metricLabel: {
    color: "#6b7280",
    fontSize: 11,
    textTransform: "uppercase" as const,
    letterSpacing: 0.5,
  },
  metricValue: {
    color: "#e0e6ed",
    fontSize: 12,
    fontFamily: "monospace",
  },
  loading: {
    color: "#6b7280",
    fontSize: 13,
    textAlign: "center" as const,
    padding: 24,
  },
  error: {
    color: "#ef4444",
    fontSize: 12,
    padding: "12px 16px",
    textAlign: "center" as const,
  },
  refreshBtn: {
    marginTop: 16,
    width: "100%",
    background: "#1f2937",
    border: "1px solid #374151",
    color: "#9ca3af",
    borderRadius: 4,
    padding: "8px",
    fontSize: 12,
    cursor: "pointer",
  },
};
