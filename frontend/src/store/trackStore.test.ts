import { beforeEach, describe, expect, it } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "../test/server";
import { useTrackStore } from "./trackStore";

describe("trackStore", () => {
  beforeEach(() => {
    useTrackStore.getState().clear();
  });

  it("fetchTrack() populates points on success", async () => {
    const points = [
      {
        icao24: "abc",
        location: { type: "Point", coordinates: [0, 0] as [number, number] },
        altitude: 1000,
        velocity: 100,
        heading: 90,
        verticalRate: 0,
        onGround: false,
        timestamp: "2026-05-11T12:00:00Z",
      },
    ];
    server.use(
      http.get("/api/v1/aircraft/abc/track", () =>
        HttpResponse.json({ data: points })
      )
    );

    await useTrackStore.getState().fetchTrack("abc");

    const s = useTrackStore.getState();
    expect(s.points).toHaveLength(1);
    expect(s.points[0].icao24).toBe("abc");
    expect(s.loading).toBe(false);
    expect(s.error).toBeNull();
    expect(s.icao24).toBe("abc");
  });

  it("fetchTrack() sets error on failure and clears loading", async () => {
    server.use(
      http.get("/api/v1/aircraft/missing/track", () =>
        HttpResponse.json({ error: "aircraft not found" }, { status: 404 })
      )
    );

    await useTrackStore.getState().fetchTrack("missing");

    const s = useTrackStore.getState();
    expect(s.loading).toBe(false);
    expect(s.error).toMatch(/not found/i);
    expect(s.points).toEqual([]);
  });

  it("clear() resets the store", () => {
    useTrackStore.setState({
      points: [
        {
          icao24: "x",
          location: { type: "Point", coordinates: [0, 0] },
          altitude: null,
          velocity: null,
          heading: null,
          verticalRate: null,
          onGround: false,
          timestamp: "",
        },
      ],
      loading: true,
      error: "stale",
      icao24: "x",
    });
    useTrackStore.getState().clear();
    const s = useTrackStore.getState();
    expect(s.points).toEqual([]);
    expect(s.loading).toBe(false);
    expect(s.error).toBeNull();
    expect(s.icao24).toBeNull();
  });
});
