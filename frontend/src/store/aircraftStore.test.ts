import { beforeEach, describe, expect, it } from "vitest";
import type { AircraftState } from "../api/types";
import { useAircraftStore } from "./aircraftStore";

function makeAC(icao24: string, callsign = ""): AircraftState {
  return {
    icao24,
    callsign,
    originCountry: "US",
    location: { type: "Point", coordinates: [0, 0] },
    baroAltitude: null,
    geoAltitude: null,
    velocity: null,
    trueTrack: null,
    verticalRate: null,
    onGround: false,
    squawk: null,
    sourceTimestamp: "",
    updatedAt: "",
  };
}

describe("aircraftStore", () => {
  beforeEach(() => {
    useAircraftStore.getState().clear();
  });

  it("applyUpdates inserts new aircraft", () => {
    useAircraftStore.getState().applyUpdates([makeAC("abc", "UAL1"), makeAC("def", "DAL2")]);
    const s = useAircraftStore.getState();
    expect(s.aircraft.size).toBe(2);
    expect(s.aircraftList).toHaveLength(2);
  });

  it("applyUpdates upserts (replaces) existing aircraft by icao24", () => {
    useAircraftStore.getState().applyUpdates([makeAC("abc", "OLD")]);
    useAircraftStore.getState().applyUpdates([makeAC("abc", "NEW")]);
    const s = useAircraftStore.getState();
    expect(s.aircraft.size).toBe(1);
    expect(s.aircraft.get("abc")?.callsign).toBe("NEW");
  });

  it("aircraftList is sorted by callsign (falls back to icao24)", () => {
    useAircraftStore
      .getState()
      .applyUpdates([
        makeAC("zzz", "ZULU"),
        makeAC("aaa", "ALPHA"),
        makeAC("mmm"), // no callsign — sorts by icao24
      ]);
    const callsigns = useAircraftStore
      .getState()
      .aircraftList.map((a) => a.callsign || a.icao24);
    expect(callsigns).toEqual(["ALPHA", "mmm", "ZULU"]);
  });

  it("applyRemovals deletes by icao24 and updates aircraftList", () => {
    useAircraftStore
      .getState()
      .applyUpdates([makeAC("a"), makeAC("b"), makeAC("c")]);
    useAircraftStore.getState().applyRemovals(["a", "c"]);
    const s = useAircraftStore.getState();
    expect(s.aircraft.size).toBe(1);
    expect(s.aircraft.has("b")).toBe(true);
    expect(s.aircraftList).toHaveLength(1);
  });

  it("applyRemovals clears selectedIcao24 when the selected aircraft is removed", () => {
    useAircraftStore.getState().applyUpdates([makeAC("abc")]);
    useAircraftStore.getState().selectAircraft("abc");
    useAircraftStore.getState().applyRemovals(["abc"]);
    expect(useAircraftStore.getState().selectedIcao24).toBeNull();
  });

  it("applyRemovals preserves selectedIcao24 when a different aircraft is removed", () => {
    useAircraftStore.getState().applyUpdates([makeAC("abc"), makeAC("def")]);
    useAircraftStore.getState().selectAircraft("abc");
    useAircraftStore.getState().applyRemovals(["def"]);
    expect(useAircraftStore.getState().selectedIcao24).toBe("abc");
  });

  it("selectAircraft updates selectedIcao24", () => {
    useAircraftStore.getState().selectAircraft("xyz");
    expect(useAircraftStore.getState().selectedIcao24).toBe("xyz");
    useAircraftStore.getState().selectAircraft(null);
    expect(useAircraftStore.getState().selectedIcao24).toBeNull();
  });

  it("clear resets everything", () => {
    useAircraftStore.getState().applyUpdates([makeAC("a")]);
    useAircraftStore.getState().selectAircraft("a");
    useAircraftStore.getState().clear();
    const s = useAircraftStore.getState();
    expect(s.aircraft.size).toBe(0);
    expect(s.aircraftList).toEqual([]);
    expect(s.selectedIcao24).toBeNull();
  });
});
