import { describe, expect, it, vi } from "vitest";
import {
  formatAltitude,
  formatCoordinate,
  formatHeading,
  formatSpeed,
  formatVerticalRate,
  timeAgo,
} from "./formatting";

describe("formatAltitude", () => {
  it("returns N/A when input is null", () => {
    expect(formatAltitude(null)).toBe("N/A");
  });
  it("converts meters to feet and rounds", () => {
    // 1000m -> 3280.84ft -> 3,281 ft (locale separator may vary)
    const out = formatAltitude(1000);
    expect(out.endsWith(" ft")).toBe(true);
    expect(out.replace(/[,\s]/g, "")).toContain("3281ft");
  });
});

describe("formatSpeed", () => {
  it("returns N/A when input is null", () => {
    expect(formatSpeed(null)).toBe("N/A");
  });
  it("converts m/s to knots and rounds", () => {
    // 100 m/s -> 194.384 kts -> 194
    expect(formatSpeed(100)).toBe("194 kts");
  });
});

describe("formatVerticalRate", () => {
  it("returns N/A when null", () => {
    expect(formatVerticalRate(null)).toBe("N/A");
  });
  it("prefixes positive values with +", () => {
    // 5 m/s -> 984.25 fpm -> +984
    expect(formatVerticalRate(5)).toBe("+984 fpm");
  });
  it("does not prefix negative values", () => {
    expect(formatVerticalRate(-5)).toBe("-984 fpm");
  });
});

describe("formatHeading", () => {
  it("returns N/A when null", () => {
    expect(formatHeading(null)).toBe("N/A");
  });
  it("rounds and adds degree symbol", () => {
    expect(formatHeading(123.7)).toBe("124°");
  });
});

describe("formatCoordinate", () => {
  it("formats positive latitude as N", () => {
    expect(formatCoordinate(40.7128, "lat")).toBe("40.7128° N");
  });
  it("formats negative latitude as S", () => {
    expect(formatCoordinate(-33.8688, "lat")).toBe("33.8688° S");
  });
  it("formats positive longitude as E", () => {
    expect(formatCoordinate(151.2093, "lon")).toBe("151.2093° E");
  });
  it("formats negative longitude as W", () => {
    expect(formatCoordinate(-74.006, "lon")).toBe("74.0060° W");
  });
});

describe("timeAgo", () => {
  it("returns 'just now' for very recent timestamps", () => {
    vi.useFakeTimers();
    const now = new Date("2026-05-11T12:00:00Z");
    vi.setSystemTime(now);
    expect(timeAgo(new Date(now.getTime() - 1000).toISOString())).toBe("just now");
    vi.useRealTimers();
  });
  it("returns seconds for under-minute", () => {
    vi.useFakeTimers();
    const now = new Date("2026-05-11T12:00:00Z");
    vi.setSystemTime(now);
    expect(timeAgo(new Date(now.getTime() - 30_000).toISOString())).toBe("30s ago");
    vi.useRealTimers();
  });
  it("returns minutes for over-minute", () => {
    vi.useFakeTimers();
    const now = new Date("2026-05-11T12:00:00Z");
    vi.setSystemTime(now);
    expect(timeAgo(new Date(now.getTime() - 5 * 60_000).toISOString())).toBe("5m ago");
    vi.useRealTimers();
  });
});
