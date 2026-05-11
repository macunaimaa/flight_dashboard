import { useEffect, useRef } from "react";
import {
  Viewer,
  Ion,
  Color,
  createWorldTerrainAsync,
  SceneMode,
} from "cesium";
import "cesium/Build/Cesium/Widgets/widgets.css";
import { CESIUM_ION_TOKEN } from "../utils/constants";
import { useCesiumSync } from "../hooks/useCesiumSync";

if (CESIUM_ION_TOKEN) {
  Ion.defaultAccessToken = CESIUM_ION_TOKEN;
}

export function CesiumGlobe() {
  const containerRef = useRef<HTMLDivElement>(null);
  const viewerRef = useRef<Viewer | null>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    const viewer = new Viewer(containerRef.current, {
      animation: false,
      timeline: false,
      baseLayerPicker: false,
      sceneModePicker: false,
      navigationHelpButton: false,
      geocoder: false,
      homeButton: false,
      fullscreenButton: false,
      infoBox: false,
      selectionIndicator: false,
      sceneMode: SceneMode.SCENE3D,
      requestRenderMode: true,
      maximumRenderTimeChange: Infinity,
    });

    viewer.scene.highDynamicRange = false;
    viewer.scene.fog.enabled = true;
    viewer.scene.globe.enableLighting = false;
    viewer.scene.screenSpaceCameraController.enableTilt = true;
    viewer.scene.screenSpaceCameraController.enableZoom = true;

    createWorldTerrainAsync()
      .then((terrain) => {
        if (!viewer.isDestroyed()) {
          viewer.terrainProvider = terrain;
        }
      })
      .catch(() => {});

    viewer.scene.globe.baseColor = Color.fromCssColorString("#0a0e17");
    viewer.scene.backgroundColor = Color.fromCssColorString("#0a0e17");
    if (viewer.scene.skyBox) viewer.scene.skyBox.show = true;
    if (viewer.scene.sun) viewer.scene.sun.show = false;
    if (viewer.scene.moon) viewer.scene.moon.show = false;

    viewerRef.current = viewer;

    return () => {
      if (!viewer.isDestroyed()) {
        viewer.destroy();
      }
      viewerRef.current = null;
    };
  }, []);

  useCesiumSync(viewerRef);

  return (
    <div
      ref={containerRef}
      style={{
        width: "100%",
        height: "100%",
        position: "absolute",
        top: 0,
        left: 0,
      }}
    />
  );
}
