import {
  Viewer,
  Entity,
  ConstantPositionProperty,
  ConstantProperty,
  HorizontalOrigin,
  VerticalOrigin,
  HeightReference,
  LabelStyle,
  Cartesian2,
  NearFarScalar,
  Color,
} from "cesium";
import type { AircraftState } from "../api/types";
import { aircraftPosition, aircraftRotation, aircraftColor, aircraftLabel } from "./utils";

const AIRCRAFT_SVG = `data:image/svg+xml;base64,${btoa(`
<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32">
  <path d="M16 2 L13 12 L4 15 L13 17 L14 28 L16 25 L18 28 L19 17 L28 15 L19 12 Z" fill="white" stroke="#333" stroke-width="0.5"/>
</svg>
`)}`;

export function createAircraftEntity(
  viewer: Viewer,
  state: AircraftState
): Entity {
  return viewer.entities.add({
    id: `aircraft-${state.icao24}`,
    position: aircraftPosition(state),
    billboard: {
      image: AIRCRAFT_SVG,
      scale: 0.7,
      rotation: aircraftRotation(state),
      horizontalOrigin: HorizontalOrigin.CENTER,
      verticalOrigin: VerticalOrigin.CENTER,
      heightReference: HeightReference.NONE,
      color: aircraftColor(state),
      disableDepthTestDistance: Number.POSITIVE_INFINITY,
    },
    label: {
      text: aircraftLabel(state),
      font: "11px monospace",
      fillColor: Color.WHITE,
      outlineColor: Color.BLACK,
      outlineWidth: 2,
      style: LabelStyle.FILL_AND_OUTLINE,
      pixelOffset: new Cartesian2(0, -22),
      scaleByDistance: new NearFarScalar(5e4, 1.0, 5e6, 0.0),
      showBackground: false,
      disableDepthTestDistance: Number.POSITIVE_INFINITY,
    },
    properties: {
      icao24: state.icao24,
      callsign: state.callsign,
      type: "aircraft",
    } as any,
  });
}

export function updateAircraftEntity(
  entity: Entity,
  state: AircraftState
): void {
  entity.position = new ConstantPositionProperty(aircraftPosition(state));

  if (entity.billboard) {
    entity.billboard.rotation = new ConstantProperty(
      aircraftRotation(state)
    );
    entity.billboard.color = new ConstantProperty(aircraftColor(state));
  }

  if (entity.label) {
    entity.label.text = new ConstantProperty(aircraftLabel(state));
  }
}
