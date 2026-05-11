import { describe, expect, it } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "../test/server";
import { apiClient } from "./client";
import { getAircraft, getAircraftInBBox, getTrack, listAircraft } from "./aircraft";

describe("aircraft api", () => {
  it("listAircraft() builds query string from params and forwards Authorization header", async () => {
    apiClient.setToken("test-token");

    let capturedURL = "";
    let capturedAuth: string | null = null;
    server.use(
      http.get("/api/v1/aircraft", ({ request }) => {
        capturedURL = request.url;
        capturedAuth = request.headers.get("authorization");
        return HttpResponse.json({ data: [] });
      })
    );

    await listAircraft({
      callsign: "UAL",
      country: "US",
      onGround: true,
      limit: 25,
      offset: 50,
    });

    expect(capturedAuth).toBe("Bearer test-token");
    const url = new URL(capturedURL);
    expect(url.searchParams.get("callsign")).toBe("UAL");
    expect(url.searchParams.get("country")).toBe("US");
    expect(url.searchParams.get("on_ground")).toBe("true");
    expect(url.searchParams.get("limit")).toBe("25");
    expect(url.searchParams.get("offset")).toBe("50");

    apiClient.setToken(null);
  });

  it("listAircraft() with no params omits query string", async () => {
    let capturedURL = "";
    server.use(
      http.get("/api/v1/aircraft", ({ request }) => {
        capturedURL = request.url;
        return HttpResponse.json({ data: [] });
      })
    );

    await listAircraft();

    expect(new URL(capturedURL).search).toBe("");
  });

  it("getAircraft() hits the correct path", async () => {
    let capturedPath = "";
    server.use(
      http.get("/api/v1/aircraft/abc123", ({ request }) => {
        capturedPath = new URL(request.url).pathname;
        return HttpResponse.json({
          data: { icao24: "abc123", callsign: "UAL" },
        });
      })
    );

    const resp = await getAircraft("abc123");
    expect(capturedPath).toBe("/api/v1/aircraft/abc123");
    expect(resp.data.icao24).toBe("abc123");
  });

  it("getAircraftInBBox() includes bbox params", async () => {
    let capturedURL = "";
    server.use(
      http.get("/api/v1/aircraft/bbox", ({ request }) => {
        capturedURL = request.url;
        return HttpResponse.json({ data: [] });
      })
    );

    await getAircraftInBBox(10, 20, 30, 40, 500);

    const u = new URL(capturedURL);
    expect(u.searchParams.get("lamin")).toBe("10");
    expect(u.searchParams.get("lomin")).toBe("20");
    expect(u.searchParams.get("lamax")).toBe("30");
    expect(u.searchParams.get("lomax")).toBe("40");
    expect(u.searchParams.get("limit")).toBe("500");
  });

  it("getTrack() forwards from/to when present", async () => {
    let capturedURL = "";
    server.use(
      http.get("/api/v1/aircraft/abc/track", ({ request }) => {
        capturedURL = request.url;
        return HttpResponse.json({ data: [] });
      })
    );

    await getTrack("abc", "2026-01-01T00:00:00Z", "2026-01-02T00:00:00Z");

    const u = new URL(capturedURL);
    expect(u.searchParams.get("from")).toBe("2026-01-01T00:00:00Z");
    expect(u.searchParams.get("to")).toBe("2026-01-02T00:00:00Z");
  });

  it("api client throws on non-2xx with backend error message", async () => {
    server.use(
      http.get("/api/v1/aircraft/missing", () =>
        HttpResponse.json({ error: "aircraft not found" }, { status: 404 })
      )
    );

    await expect(getAircraft("missing")).rejects.toThrow(/aircraft not found/);
  });
});
