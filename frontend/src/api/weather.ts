export interface WeatherData {
  latitude: number;
  longitude: number;
  current: {
    temperature_2m: number;
    relative_humidity_2m: number;
    wind_speed_10m: number;
    wind_direction_10m: number;
    weather_code: number;
    cloud_cover: number;
    surface_pressure: number;
    is_day: number;
  };
  current_units: {
    temperature_2m: string;
    wind_speed_10m: string;
    surface_pressure: string;
  };
}

export async function fetchWeather(lat: number, lon: number): Promise<WeatherData> {
  const params = new URLSearchParams({
    latitude: lat.toString(),
    longitude: lon.toString(),
    current: "temperature_2m,relative_humidity_2m,wind_speed_10m,wind_direction_10m,weather_code,cloud_cover,surface_pressure,is_day",
  });
  const resp = await fetch(`https://api.open-meteo.com/v1/forecast?${params}`);
  if (!resp.ok) throw new Error(`Weather API error: ${resp.status}`);
  return resp.json();
}

export function weatherCodeToDescription(code: number): string {
  const map: Record<number, string> = {
    0: "Clear sky", 1: "Mainly clear", 2: "Partly cloudy", 3: "Overcast",
    45: "Fog", 48: "Depositing rime fog",
    51: "Light drizzle", 53: "Moderate drizzle", 55: "Dense drizzle",
    61: "Slight rain", 63: "Moderate rain", 65: "Heavy rain",
    71: "Slight snow", 73: "Moderate snow", 75: "Heavy snow",
    77: "Snow grains", 80: "Slight showers", 81: "Moderate showers", 82: "Violent showers",
    85: "Slight snow showers", 86: "Heavy snow showers",
    95: "Thunderstorm", 96: "Thunderstorm + hail", 99: "Thunderstorm + heavy hail",
  };
  return map[code] || "Unknown";
}

export function weatherCodeToIcon(code: number): string {
  if (code === 0) return "\u2600";
  if (code <= 2) return "\u26C5";
  if (code === 3) return "\u2601";
  if (code <= 48) return "\u{1F32B}";
  if (code <= 55) return "\u{1F4A7}";
  if (code <= 65) return "\u{1F327}";
  if (code <= 77) return "\u2744";
  if (code <= 82) return "\u{1F327}";
  return "\u26C8";
}
