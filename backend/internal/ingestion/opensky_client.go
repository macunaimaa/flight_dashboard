package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

type OpenSkyClient struct {
	httpClient *http.Client
	baseURL    string
	username   string
	password   string
}

func NewOpenSkyClient(username, password string) *OpenSkyClient {
	return &OpenSkyClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://opensky-network.org/api",
		username:   username,
		password:   password,
	}
}

// RawStateVector represents a single aircraft state from the OpenSky API.
// Fields are positional in the JSON response array.
type RawStateVector struct {
	ICAO24         string
	Callsign       *string
	OriginCountry  string
	TimePosition   *int64
	LastContact    int64
	Longitude      *float64
	Latitude       *float64
	BaroAltitude   *float64
	OnGround       bool
	Velocity       *float64
	TrueTrack      *float64
	VerticalRate   *float64
	Sensors        []int
	GeoAltitude    *float64
	Squawk         *string
	SPI            bool
	PositionSource int
}

type openSkyResponse struct {
	Time   int64             `json:"time"`
	States []json.RawMessage `json:"states"`
}

// FetchStates fetches aircraft states within a bounding box from OpenSky Network.
func (c *OpenSkyClient) FetchStates(ctx context.Context, bbox domain.BBox) ([]RawStateVector, int64, error) {
	url := fmt.Sprintf("%s/states/all?lamin=%.4f&lomin=%.4f&lamax=%.4f&lomax=%.4f",
		c.baseURL, bbox.LaMin, bbox.LoMin, bbox.LaMax, bbox.LoMax)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("opensky request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, 0, fmt.Errorf("opensky rate limited (429)")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, 0, fmt.Errorf("opensky status %d: %s", resp.StatusCode, string(body))
	}

	var raw openSkyResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, 0, fmt.Errorf("decode response: %w", err)
	}

	vectors := make([]RawStateVector, 0, len(raw.States))
	for _, stateJSON := range raw.States {
		v, err := parseStateVector(stateJSON)
		if err != nil {
			continue // skip malformed entries
		}
		vectors = append(vectors, *v)
	}

	return vectors, raw.Time, nil
}

// parseStateVector parses a single state vector from positional JSON array.
// OpenSky API format: [icao24, callsign, origin_country, time_position, last_contact,
//
//	longitude, latitude, baro_altitude, on_ground, velocity, true_track,
//	vertical_rate, sensors, geo_altitude, squawk, spi, position_source]
func parseStateVector(data json.RawMessage) (*RawStateVector, error) {
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return nil, err
	}
	if len(arr) < 17 {
		return nil, fmt.Errorf("state vector too short: %d fields", len(arr))
	}

	v := &RawStateVector{}

	json.Unmarshal(arr[0], &v.ICAO24)
	json.Unmarshal(arr[1], &v.Callsign)
	json.Unmarshal(arr[2], &v.OriginCountry)
	json.Unmarshal(arr[3], &v.TimePosition)
	json.Unmarshal(arr[4], &v.LastContact)
	json.Unmarshal(arr[5], &v.Longitude)
	json.Unmarshal(arr[6], &v.Latitude)
	json.Unmarshal(arr[7], &v.BaroAltitude)
	json.Unmarshal(arr[8], &v.OnGround)
	json.Unmarshal(arr[9], &v.Velocity)
	json.Unmarshal(arr[10], &v.TrueTrack)
	json.Unmarshal(arr[11], &v.VerticalRate)
	// arr[12] = sensors (skip)
	json.Unmarshal(arr[13], &v.GeoAltitude)
	json.Unmarshal(arr[14], &v.Squawk)
	json.Unmarshal(arr[15], &v.SPI)
	json.Unmarshal(arr[16], &v.PositionSource)

	return v, nil
}
