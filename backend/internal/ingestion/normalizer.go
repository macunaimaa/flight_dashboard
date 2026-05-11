package ingestion

import (
	"fmt"
	"strings"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

// Normalize converts a RawStateVector into a domain AircraftLiveState.
// Returns nil if the vector lacks position data (lat/lon).
func Normalize(raw RawStateVector, apiTimestamp int64) (*domain.AircraftLiveState, error) {
	// Skip aircraft without position
	if raw.Latitude == nil || raw.Longitude == nil {
		return nil, fmt.Errorf("no position data for %s", raw.ICAO24)
	}

	icao24 := strings.TrimSpace(raw.ICAO24)

	// Build external event key for idempotency
	var eventKey string
	if raw.TimePosition != nil {
		eventKey = fmt.Sprintf("%s:%d", icao24, *raw.TimePosition)
	} else {
		eventKey = fmt.Sprintf("%s:%d", icao24, apiTimestamp)
	}

	// Determine source timestamp
	var sourceTS time.Time
	if raw.TimePosition != nil {
		sourceTS = time.Unix(*raw.TimePosition, 0)
	} else {
		sourceTS = time.Unix(apiTimestamp, 0)
	}

	callsign := ""
	if raw.Callsign != nil {
		callsign = strings.TrimSpace(*raw.Callsign)
	}

	state := &domain.AircraftLiveState{
		ICAO24:         icao24,
		Callsign:       callsign,
		OriginCountry:  raw.OriginCountry,
		Location:       domain.NewPoint(*raw.Longitude, *raw.Latitude),
		BaroAltitudeM:  raw.BaroAltitude,
		GeoAltitudeM:   raw.GeoAltitude,
		VelocityMS:     raw.Velocity,
		TrueTrackDeg:   raw.TrueTrack,
		VerticalRateMS:  raw.VerticalRate,
		OnGround:       raw.OnGround,
		Squawk:         raw.Squawk,
		SPI:            raw.SPI,
		PositionSource: raw.PositionSource,

		ExternalEventKey: eventKey,
		SourceTimestamp:  sourceTS,
		LastContact:      time.Unix(raw.LastContact, 0),
		IngestedAt:       time.Now(),
		UpdatedAt:        time.Now(),
	}

	return state, nil
}

// ToTrackPoint converts an AircraftLiveState to a TrackPoint for historical storage.
func ToTrackPoint(state *domain.AircraftLiveState) domain.TrackPoint {
	return domain.TrackPoint{
		TenantID:         state.TenantID,
		ICAO24:           state.ICAO24,
		Location:         state.Location,
		AltitudeM:        state.BaroAltitudeM,
		VelocityMS:       state.VelocityMS,
		HeadingDeg:       state.TrueTrackDeg,
		VerticalRateMS:   state.VerticalRateMS,
		OnGround:         state.OnGround,
		Timestamp:        state.SourceTimestamp,
		ExternalEventKey: state.ExternalEventKey,
	}
}
