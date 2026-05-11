// Use relative URLs so Vite's dev proxy handles routing to the backend.
// In production, set VITE_API_URL / VITE_WS_URL to the real backend origin.
export const API_URL = import.meta.env.VITE_API_URL || "";
export const WS_URL =
  import.meta.env.VITE_WS_URL ||
  `${window.location.protocol === "https:" ? "wss:" : "ws:"}//${window.location.host}`;
export const CESIUM_ION_TOKEN = import.meta.env.VITE_CESIUM_ION_TOKEN || "";
