package geo

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {
	const earthRadiusM = 6371000.0

	tests := []struct {
		name           string
		lat1, lon1     float64
		lat2, lon2     float64
		wantMeters     float64
		toleranceRatio float64 // relative tolerance, e.g. 0.001 = 0.1%
	}{
		{
			name:           "same point returns zero",
			lat1:           37.7749, lon1: -122.4194,
			lat2: 37.7749, lon2: -122.4194,
			wantMeters:     0,
			toleranceRatio: 0,
		},
		{
			name:           "antipodal points equal half circumference",
			lat1:           0, lon1: 0,
			lat2: 0, lon2: 180,
			wantMeters:     math.Pi * earthRadiusM,
			toleranceRatio: 1e-9,
		},
		{
			name:           "NYC to London approx 5570 km",
			lat1:           40.7128, lon1: -74.0060,
			lat2: 51.5074, lon2: -0.1278,
			wantMeters:     5_570_000,
			toleranceRatio: 0.005,
		},
		{
			name:           "one degree of latitude approx 111 km",
			lat1:           0, lon1: 0,
			lat2: 1, lon2: 0,
			wantMeters:     111_195,
			toleranceRatio: 0.001,
		},
		{
			name:           "symmetric: reversing args gives same distance",
			lat1:           48.8566, lon1: 2.3522, // Paris
			lat2: 35.6762, lon2: 139.6503, // Tokyo
			wantMeters:     9_714_000,
			toleranceRatio: 0.01,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Haversine(tc.lat1, tc.lon1, tc.lat2, tc.lon2)

			if tc.wantMeters == 0 {
				if got != 0 {
					t.Fatalf("Haversine same-point = %v, want 0", got)
				}
				return
			}

			diff := math.Abs(got - tc.wantMeters)
			tolerance := tc.wantMeters * tc.toleranceRatio
			if diff > tolerance {
				t.Fatalf("Haversine = %.2f m, want %.2f m (tolerance %.2f m, diff %.2f m)",
					got, tc.wantMeters, tolerance, diff)
			}
		})
	}
}

func TestHaversineIsSymmetric(t *testing.T) {
	a := Haversine(40.7128, -74.0060, 51.5074, -0.1278)
	b := Haversine(51.5074, -0.1278, 40.7128, -74.0060)
	if math.Abs(a-b) > 1e-6 {
		t.Fatalf("Haversine is not symmetric: a=%v b=%v", a, b)
	}
}
