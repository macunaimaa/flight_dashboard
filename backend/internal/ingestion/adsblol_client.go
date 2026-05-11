package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

// ADSBLolClient fetches live aircraft data from the free adsb.lol API.
type ADSBLolClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewADSBLolClient() *ADSBLolClient {
	return &ADSBLolClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://api.adsb.lol",
	}
}

type adsbLolResponse struct {
	AC    []adsbLolAircraft `json:"ac"`
	Msg   string            `json:"msg"`
	Now   int64             `json:"now"`
	Total int               `json:"total"`
}

type adsbLolAircraft struct {
	Hex           string      `json:"hex"`
	Flight        *string     `json:"flight"`
	AltBaro       interface{} `json:"alt_baro"` // can be int or "ground"
	AltGeom       *float64    `json:"alt_geom"`
	GS            *float64    `json:"gs"`
	Track         *float64    `json:"track"`
	BaroRate      *float64    `json:"baro_rate"`
	Squawk        *string     `json:"squawk"`
	Category      *string     `json:"category"`
	Lat           *float64    `json:"lat"`
	Lon           *float64    `json:"lon"`
	SPI           int         `json:"spi"`
	Seen          *float64    `json:"seen"`
	SeenPos       *float64    `json:"seen_pos"`
	Registration  *string     `json:"r"`
	AircraftType  *string     `json:"t"`
	TrueHeading   *float64    `json:"true_heading"`
	MagHeading    *float64    `json:"mag_heading"`
	Emergency     *string     `json:"emergency"`
	DbFlags       *int        `json:"dbFlags"`
}

// FetchStates fetches aircraft within the bounding box by converting to
// a center-point + radius query (adsb.lol uses point/radius, not bbox).
func (c *ADSBLolClient) FetchStates(ctx context.Context, bbox domain.BBox) ([]RawStateVector, int64, error) {
	// Convert bbox to center point + radius in nautical miles
	centerLat := (bbox.LaMin + bbox.LaMax) / 2
	centerLon := (bbox.LoMin + bbox.LoMax) / 2
	radiusNM := bboxRadiusNM(bbox)
	if radiusNM > 250 {
		radiusNM = 250 // adsb.lol max is 250 NM
	}

	url := fmt.Sprintf("%s/v2/lat/%.4f/lon/%.4f/dist/%.0f",
		c.baseURL, centerLat, centerLon, radiusNM)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("adsb.lol request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, 0, fmt.Errorf("adsb.lol rate limited (429)")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, 0, fmt.Errorf("adsb.lol status %d: %s", resp.StatusCode, string(body))
	}

	var raw adsbLolResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, 0, fmt.Errorf("decode adsb.lol response: %w", err)
	}

	now := time.Now().Unix()
	vectors := make([]RawStateVector, 0, len(raw.AC))
	for _, ac := range raw.AC {
		v := convertADSBLol(ac, now)
		if v != nil {
			vectors = append(vectors, *v)
		}
	}

	return vectors, now, nil
}

// convertADSBLol converts an adsb.lol aircraft to our RawStateVector format.
func convertADSBLol(ac adsbLolAircraft, now int64) *RawStateVector {
	if ac.Lat == nil || ac.Lon == nil {
		return nil // skip aircraft without position
	}

	v := &RawStateVector{
		ICAO24:        ac.Hex,
		Callsign:      ac.Flight,
		OriginCountry: "", // adsb.lol doesn't provide origin country
		LastContact:   now,
		Latitude:      ac.Lat,
		Longitude:     ac.Lon,
		TrueTrack:     ac.Track,
		Squawk:        ac.Squawk,
		SPI:           ac.SPI == 1,
	}

	// TimePosition
	tp := now
	v.TimePosition = &tp

	// alt_baro can be a number (feet) or "ground"
	switch alt := ac.AltBaro.(type) {
	case float64:
		altMeters := alt * 0.3048 // feet to meters
		v.BaroAltitude = &altMeters
		v.OnGround = false
	case string:
		if alt == "ground" {
			v.OnGround = true
		}
	case nil:
		// no altitude info
	}

	// GeoAltitude (feet to meters)
	if ac.AltGeom != nil {
		geoM := *ac.AltGeom * 0.3048
		v.GeoAltitude = &geoM
	}

	// Ground speed: knots to m/s
	if ac.GS != nil {
		velMS := *ac.GS * 0.514444
		v.Velocity = &velMS
	}

	// Vertical rate: ft/min to m/s
	if ac.BaroRate != nil {
		vrateMS := *ac.BaroRate * 0.00508
		v.VerticalRate = &vrateMS
	}

	return v
}

// bboxRadiusNM computes the radius in nautical miles that covers a bounding box.
func bboxRadiusNM(bbox domain.BBox) float64 {
	dLat := (bbox.LaMax - bbox.LaMin) * math.Pi / 180
	dLon := (bbox.LoMax - bbox.LoMin) * math.Pi / 180
	avgLat := (bbox.LaMin + bbox.LaMax) / 2 * math.Pi / 180

	// Approximate distance in km using equirectangular projection
	x := dLon * math.Cos(avgLat)
	y := dLat
	distKM := math.Sqrt(x*x+y*y) * 6371 / 2 // half-diagonal
	return distKM / 1.852                      // km to nautical miles
}
