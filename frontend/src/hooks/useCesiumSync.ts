import { useEffect, useRef, type RefObject } from "react";
import {
  Cartesian3,
  Color,
  ScreenSpaceEventType,
  type Viewer,
  type Entity,
} from "cesium";
import { useAircraftStore } from "../store/aircraftStore";
import { useTrackStore } from "../store/trackStore";
import {
  createAircraftEntity,
  updateAircraftEntity,
} from "../cesium/AircraftLayer";

function altitudeToColor(altM: number | null): Color {
  if (altM === null) return Color.GRAY;
  const ft = altM * 3.28084;
  if (ft < 5000) return Color.fromCssColorString("#22c55e");
  if (ft < 15000) return Color.fromCssColorString("#4ade80");
  if (ft < 25000) return Color.fromCssColorString("#facc15");
  if (ft < 35000) return Color.fromCssColorString("#38bdf8");
  return Color.fromCssColorString("#818cf8");
}

export function useCesiumSync(viewerRef: RefObject<Viewer | null>) {
  const entityMapRef = useRef(new Map<string, Entity>());
  const trailEntityRef = useRef<Entity | null>(null);

  useEffect(() => {
    const unsub = useAircraftStore.subscribe((state) => {
      const viewer = viewerRef.current;
      if (!viewer || viewer.isDestroyed()) return;

      const aircraft = state.aircraft;
      const entityMap = entityMapRef.current;

      for (const [icao24, aircraftState] of aircraft) {
        const existing = entityMap.get(icao24);
        if (existing) {
          updateAircraftEntity(existing, aircraftState);
        } else {
          const entity = createAircraftEntity(viewer, aircraftState);
          entityMap.set(icao24, entity);
        }
      }

      for (const [icao24, entity] of entityMap) {
        if (!aircraft.has(icao24)) {
          viewer.entities.remove(entity);
          entityMap.delete(icao24);
        }
      }
    });

    return () => {
      unsub();
      const viewer = viewerRef.current;
      if (viewer && !viewer.isDestroyed()) {
        for (const entity of entityMapRef.current.values()) {
          viewer.entities.remove(entity);
        }
      }
      entityMapRef.current.clear();
    };
  }, [viewerRef]);

  useEffect(() => {
    const viewer = viewerRef.current;
    if (!viewer || viewer.isDestroyed()) return;

    const handler = viewer.screenSpaceEventHandler;

    handler.setInputAction((movement: { position: any }) => {
      const picked = viewer.scene.pick(movement.position);
      if (picked?.id?.properties?.type?.getValue() === "aircraft") {
        const icao24 = picked.id.properties.icao24.getValue();
        useAircraftStore.getState().selectAircraft(icao24);
      } else {
        useAircraftStore.getState().selectAircraft(null);
      }
    }, ScreenSpaceEventType.LEFT_CLICK);

    return () => {
      handler.removeInputAction(ScreenSpaceEventType.LEFT_CLICK);
    };
  }, [viewerRef]);

  useEffect(() => {
    let prev: string | null = useAircraftStore.getState().selectedIcao24;
    const unsub = useAircraftStore.subscribe((state) => {
      if (state.selectedIcao24 === prev) return;
      prev = state.selectedIcao24;
      if (state.selectedIcao24) {
        useTrackStore.getState().fetchTrack(state.selectedIcao24);
      } else {
        useTrackStore.getState().clear();
      }
    });

    return () => {
      unsub();
    };
  }, []);

  useEffect(() => {
    let prevLen = useTrackStore.getState().points.length;
    const unsub = useTrackStore.subscribe((state) => {
      if (state.points.length === prevLen) return;
      prevLen = state.points.length;

      const viewer = viewerRef.current;
      if (!viewer || viewer.isDestroyed()) return;

      if (trailEntityRef.current) {
        viewer.entities.remove(trailEntityRef.current);
        trailEntityRef.current = null;
      }

      const points = state.points;
      if (points.length < 2) return;

      const positions = points.map((p: { location: { coordinates: [number, number] }; altitude: number | null }) =>
        Cartesian3.fromDegrees(
          p.location.coordinates[0],
          p.location.coordinates[1],
          p.altitude ?? 0
        )
      );

      const avgAlt =
        points.reduce((sum: number, p: { altitude: number | null }) => sum + (p.altitude ?? 0), 0) / points.length;
      const trailColor = altitudeToColor(avgAlt).withAlpha(0.8);

      const entity = viewer.entities.add({
        id: "aircraft-trail",
        polyline: {
          positions: positions,
          width: 2,
          material: trailColor,
          clampToGround: false,
        },
      });

      trailEntityRef.current = entity;
    });

    return () => {
      unsub();
      const viewer = viewerRef.current;
      if (viewer && !viewer.isDestroyed() && trailEntityRef.current) {
        viewer.entities.remove(trailEntityRef.current);
        trailEntityRef.current = null;
      }
    };
  }, [viewerRef]);
}
