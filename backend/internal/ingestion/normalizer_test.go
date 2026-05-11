package ingestion

import (
	"testing"
	"time"
)

func ptrF64(v float64) *float64 { return &v }
func ptrI64(v int64) *int64     { return &v }
func ptrStr(v string) *string   { return &v }

func TestNormalize_MissingPositionReturnsError(t *testing.T) {
	tests := []struct {
		name string
		raw  RawStateVector
	}{
		{
			name: "missing latitude",
			raw:  RawStateVector{ICAO24: "abc123", Longitude: ptrF64(10), Latitude: nil},
		},
		{
			name: "missing longitude",
			raw:  RawStateVector{ICAO24: "abc123", Longitude: nil, Latitude: ptrF64(10)},
		},
		{
			name: "both nil",
			raw:  RawStateVector{ICAO24: "abc123"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state, err := Normalize(tc.raw, time.Now().Unix())
			if err == nil {
				t.Fatalf("expected error for missing position, got nil (state=%+v)", state)
			}
			if state != nil {
				t.Fatalf("expected nil state, got %+v", state)
			}
		})
	}
}

func TestNormalize_HappyPath(t *testing.T) {
	timePos := int64(1_700_000_000)
	apiTS := int64(1_700_000_500)
	raw := RawStateVector{
		ICAO24:         "  abc123  ",
		Callsign:       ptrStr("  UAL123  "),
		OriginCountry:  "United States",
		TimePosition:   ptrI64(timePos),
		LastContact:    1_700_000_499,
		Longitude:      ptrF64(-122.4194),
		Latitude:       ptrF64(37.7749),
		BaroAltitude:   ptrF64(10000),
		GeoAltitude:    ptrF64(10050),
		Velocity:       ptrF64(250),
		TrueTrack:      ptrF64(90),
		VerticalRate:   ptrF64(-5),
		OnGround:       false,
		Squawk:         ptrStr("1200"),
		SPI:            true,
		PositionSource: 0,
	}

	state, err := Normalize(raw, apiTS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state == nil {
		t.Fatal("expected state, got nil")
	}

	if state.ICAO24 != "abc123" {
		t.Errorf("ICAO24 not trimmed: %q", state.ICAO24)
	}
	if state.Callsign != "UAL123" {
		t.Errorf("Callsign not trimmed: %q", state.Callsign)
	}
	if state.OriginCountry != "United States" {
		t.Errorf("OriginCountry = %q", state.OriginCountry)
	}
	if got := state.Longitude(); got != -122.4194 {
		t.Errorf("Longitude = %v, want -122.4194", got)
	}
	if got := state.Latitude(); got != 37.7749 {
		t.Errorf("Latitude = %v, want 37.7749", got)
	}
	if state.Location.Type != "Point" {
		t.Errorf("Location.Type = %q, want Point", state.Location.Type)
	}
	if !state.SPI {
		t.Error("SPI should be true")
	}
	if state.OnGround {
		t.Error("OnGround should be false")
	}

	// External event key uses TimePosition when present
	wantKey := "abc123:1700000000"
	if state.ExternalEventKey != wantKey {
		t.Errorf("ExternalEventKey = %q, want %q", state.ExternalEventKey, wantKey)
	}

	// Source timestamp uses TimePosition
	if state.SourceTimestamp.Unix() != timePos {
		t.Errorf("SourceTimestamp = %v, want unix %d", state.SourceTimestamp, timePos)
	}

	// LastContact converted from int64
	if state.LastContact.Unix() != 1_700_000_499 {
		t.Errorf("LastContact = %v", state.LastContact)
	}

	// Ingestion timestamps should be very recent
	if time.Since(state.IngestedAt) > time.Second {
		t.Errorf("IngestedAt not recent: %v", state.IngestedAt)
	}
}

func TestNormalize_FallsBackToAPITimestamp(t *testing.T) {
	apiTS := int64(1_700_000_500)
	raw := RawStateVector{
		ICAO24:       "deadbe",
		Longitude:    ptrF64(0),
		Latitude:     ptrF64(0),
		LastContact:  apiTS,
		TimePosition: nil, // forces fallback
	}

	state, err := Normalize(raw, apiTS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantKey := "deadbe:1700000500"
	if state.ExternalEventKey != wantKey {
		t.Errorf("ExternalEventKey = %q, want %q (apiTS fallback)", state.ExternalEventKey, wantKey)
	}
	if state.SourceTimestamp.Unix() != apiTS {
		t.Errorf("SourceTimestamp = %v, want unix %d", state.SourceTimestamp, apiTS)
	}
}

func TestNormalize_NilCallsignBecomesEmpty(t *testing.T) {
	raw := RawStateVector{
		ICAO24:    "abc",
		Longitude: ptrF64(0),
		Latitude:  ptrF64(0),
		Callsign:  nil,
	}
	state, err := Normalize(raw, time.Now().Unix())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Callsign != "" {
		t.Errorf("Callsign = %q, want empty string", state.Callsign)
	}
}

func TestToTrackPoint_CopiesFields(t *testing.T) {
	raw := RawStateVector{
		ICAO24:       "abc123",
		Longitude:    ptrF64(2.3522),
		Latitude:     ptrF64(48.8566),
		BaroAltitude: ptrF64(9000),
		Velocity:     ptrF64(220),
		TrueTrack:    ptrF64(45),
		VerticalRate: ptrF64(0),
		TimePosition: ptrI64(1_700_000_000),
		LastContact:  1_700_000_000,
		OnGround:     false,
	}
	state, err := Normalize(raw, 1_700_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	state.TenantID = "tenant-1"

	tp := ToTrackPoint(state)
	if tp.TenantID != "tenant-1" {
		t.Errorf("TenantID = %q, want tenant-1", tp.TenantID)
	}
	if tp.ICAO24 != "abc123" {
		t.Errorf("ICAO24 = %q", tp.ICAO24)
	}
	if tp.ExternalEventKey != state.ExternalEventKey {
		t.Errorf("ExternalEventKey not carried over: tp=%q state=%q",
			tp.ExternalEventKey, state.ExternalEventKey)
	}
	if tp.AltitudeM == nil || *tp.AltitudeM != 9000 {
		t.Errorf("AltitudeM = %v, want 9000", tp.AltitudeM)
	}
	if tp.HeadingDeg == nil || *tp.HeadingDeg != 45 {
		t.Errorf("HeadingDeg = %v, want 45", tp.HeadingDeg)
	}
	if !tp.Timestamp.Equal(state.SourceTimestamp) {
		t.Errorf("Timestamp = %v, want %v", tp.Timestamp, state.SourceTimestamp)
	}
}
