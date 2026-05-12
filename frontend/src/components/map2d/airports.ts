/* A small airport coordinate table used by the 2D radar overlay
   to render range rings + airport markers. Coordinates are
   [longitude, latitude] in degrees. Extend as needed. */

export const AIRPORTS: Record<string, { name: string; coord: [number, number] }> = {
  JFK: { name: "New York JFK",      coord: [-73.78,  40.64] },
  LAX: { name: "Los Angeles",       coord: [-118.41, 33.94] },
  ORD: { name: "Chicago O'Hare",    coord: [-87.90,  41.97] },
  ATL: { name: "Atlanta",           coord: [-84.43,  33.64] },
  DFW: { name: "Dallas/Fort Worth", coord: [-97.04,  32.90] },
  YYZ: { name: "Toronto Pearson",   coord: [-79.63,  43.68] },
  GRU: { name: "São Paulo",         coord: [-46.47, -23.43] },
  LHR: { name: "London Heathrow",   coord: [ -0.46,  51.47] },
  CDG: { name: "Paris CDG",         coord: [  2.55,  49.01] },
  AMS: { name: "Amsterdam",         coord: [  4.76,  52.31] },
  FRA: { name: "Frankfurt",         coord: [  8.57,  50.04] },
  MAD: { name: "Madrid",            coord: [ -3.57,  40.49] },
  DXB: { name: "Dubai",             coord: [ 55.36,  25.25] },
  DOH: { name: "Doha",              coord: [ 51.61,  25.27] },
  IST: { name: "Istanbul",          coord: [ 28.81,  41.28] },
  DEL: { name: "Delhi",             coord: [ 77.10,  28.57] },
  BOM: { name: "Mumbai",            coord: [ 72.87,  19.09] },
  HKG: { name: "Hong Kong",         coord: [113.91,  22.31] },
  PEK: { name: "Beijing",           coord: [116.59,  40.08] },
  PVG: { name: "Shanghai Pudong",   coord: [121.81,  31.14] },
  HND: { name: "Tokyo Haneda",      coord: [139.78,  35.55] },
  NRT: { name: "Tokyo Narita",      coord: [140.39,  35.76] },
  ICN: { name: "Seoul Incheon",     coord: [126.45,  37.46] },
  SIN: { name: "Singapore",         coord: [103.99,   1.36] },
  SYD: { name: "Sydney",            coord: [151.18, -33.94] },
  AKL: { name: "Auckland",          coord: [174.79, -37.01] },
  JNB: { name: "Johannesburg",      coord: [ 28.24, -26.13] },
  CAI: { name: "Cairo",             coord: [ 31.41,  30.11] },
  MEX: { name: "Mexico City",       coord: [-99.07,  19.43] },
};

/** Equirectangular projection — lon/lat → unit square (0..1, 0..1). */
export function project(lon: number, lat: number): [number, number] {
  return [(lon + 180) / 360, (90 - lat) / 180];
}
