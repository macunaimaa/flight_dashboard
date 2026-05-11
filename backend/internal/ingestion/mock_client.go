package ingestion

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

// MockClient generates realistic simulated aircraft data for development.
type MockClient struct {
	aircraft []mockAircraft
	rng      *rand.Rand
}

type mockAircraft struct {
	icao24    string
	callsign  string
	country   string
	lat       float64
	lon       float64
	altitude  float64
	velocity  float64
	heading   float64
	vrate     float64
	onGround  bool
}

var mockFlights = []mockAircraft{
	{"a1b2c3", "GLO1234", "Brazil", -23.4356, -46.4731, 10000, 240, 45, 2.5, false},
	{"d4e5f6", "TAM5678", "Brazil", -22.9068, -43.1729, 8500, 220, 180, -1.2, false},
	{"789abc", "AZU9012", "Brazil", -23.0077, -44.3192, 11000, 250, 270, 0, false},
	{"def012", "AAL345", "United States", -22.8103, -47.0626, 12000, 260, 90, 1.8, false},
	{"345678", "UAE772", "United Arab Emirates", -23.6273, -46.6566, 9800, 230, 315, -0.5, false},
	{"9abcde", "DLH501", "Germany", -24.0083, -46.2372, 11500, 255, 135, 0.3, false},
	{"f01234", "AFR447", "France", -21.7856, -48.1794, 10200, 235, 200, -2.0, false},
	{"567890", "BAW247", "United Kingdom", -22.3361, -49.0706, 9200, 215, 60, 1.5, false},
	{"abcdef", "KLM886", "Netherlands", -23.5015, -47.4526, 10800, 245, 150, 0, false},
	{"112233", "TAP903", "Portugal", -24.5218, -47.1278, 8000, 200, 340, -3.0, false},
	{"445566", "GOL7890", "Brazil", -23.0000, -45.5000, 0, 0, 120, 0, true},
	{"778899", "LAT456", "Brazil", -22.5000, -43.5000, 3000, 180, 225, 5.0, false},
	{"aabb01", "AVA321", "Colombia", -21.5000, -49.5000, 11200, 248, 170, 0.2, false},
	{"ccdd02", "CPA880", "China", -24.8000, -48.0000, 10500, 252, 80, -0.8, false},
	{"eeff03", "QTR774", "Qatar", -23.2000, -44.8000, 9600, 238, 295, 1.0, false},
}

func NewMockClient() *MockClient {
	aircraft := make([]mockAircraft, len(mockFlights))
	copy(aircraft, mockFlights)
	return &MockClient{
		aircraft: aircraft,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// FetchStates returns simulated aircraft with slightly updated positions each call.
func (m *MockClient) FetchStates(_ context.Context, _ domain.BBox) ([]RawStateVector, int64, error) {
	now := time.Now().Unix()
	vectors := make([]RawStateVector, 0, len(m.aircraft))

	for i := range m.aircraft {
		ac := &m.aircraft[i]

		if !ac.onGround {
			// Move aircraft along heading
			headingRad := ac.heading * math.Pi / 180
			speed := ac.velocity / 111000.0 * 30 // degrees per 30s at given speed
			ac.lat += math.Cos(headingRad) * speed
			ac.lon += math.Sin(headingRad) * speed / math.Cos(ac.lat*math.Pi/180)

			// Small random perturbations
			ac.heading += (m.rng.Float64() - 0.5) * 5
			ac.velocity += (m.rng.Float64() - 0.5) * 10
			ac.altitude += ac.vrate * 30
			ac.vrate += (m.rng.Float64() - 0.5) * 0.5

			// Keep in reasonable bounds
			if ac.heading < 0 {
				ac.heading += 360
			}
			if ac.heading >= 360 {
				ac.heading -= 360
			}
			if ac.altitude < 1000 {
				ac.altitude = 1000
				ac.vrate = math.Abs(ac.vrate)
			}
			if ac.altitude > 13000 {
				ac.altitude = 13000
				ac.vrate = -math.Abs(ac.vrate)
			}
			if ac.velocity < 150 {
				ac.velocity = 150
			}
			if ac.velocity > 280 {
				ac.velocity = 280
			}

			// Wrap around region
			if ac.lat < -26 || ac.lat > -19 || ac.lon < -51 || ac.lon > -42 {
				ac.lat = -23.0 + (m.rng.Float64()-0.5)*4
				ac.lon = -46.5 + (m.rng.Float64()-0.5)*6
				ac.heading = m.rng.Float64() * 360
			}
		}

		callsign := ac.callsign
		lat := ac.lat
		lon := ac.lon
		alt := ac.altitude
		vel := ac.velocity
		track := ac.heading
		vrate := ac.vrate
		tp := now

		vectors = append(vectors, RawStateVector{
			ICAO24:        ac.icao24,
			Callsign:      &callsign,
			OriginCountry: ac.country,
			TimePosition:  &tp,
			LastContact:   now,
			Latitude:      &lat,
			Longitude:     &lon,
			BaroAltitude:  &alt,
			OnGround:      ac.onGround,
			Velocity:      &vel,
			TrueTrack:     &track,
			VerticalRate:  &vrate,
			GeoAltitude:   &alt,
			PositionSource: 0,
		})
	}

	return vectors, now, nil
}
