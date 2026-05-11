package domain

import "time"

type GeoJSONPoint struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [lon, lat]
}

func NewPoint(lon, lat float64) GeoJSONPoint {
	return GeoJSONPoint{
		Type:        "Point",
		Coordinates: []float64{lon, lat},
	}
}

type AircraftLiveState struct {
	TenantID       string       `bson:"tenant_id" json:"tenantId"`
	ICAO24         string       `bson:"icao24" json:"icao24"`
	Callsign       string       `bson:"callsign" json:"callsign"`
	OriginCountry  string       `bson:"origin_country" json:"originCountry"`
	Location       GeoJSONPoint `bson:"location" json:"location"`
	BaroAltitudeM  *float64     `bson:"baro_altitude_m" json:"baroAltitude"`
	GeoAltitudeM   *float64     `bson:"geo_altitude_m" json:"geoAltitude"`
	VelocityMS     *float64     `bson:"velocity_ms" json:"velocity"`
	TrueTrackDeg   *float64     `bson:"true_track_deg" json:"trueTrack"`
	VerticalRateMS *float64     `bson:"vertical_rate_ms" json:"verticalRate"`
	OnGround       bool         `bson:"on_ground" json:"onGround"`
	Squawk         *string      `bson:"squawk" json:"squawk"`
	SPI            bool         `bson:"spi" json:"spi"`
	PositionSource int          `bson:"position_source" json:"positionSource"`

	ExternalEventKey string    `bson:"external_event_key" json:"-"`
	SourceTimestamp  time.Time `bson:"source_timestamp" json:"sourceTimestamp"`
	LastContact      time.Time `bson:"last_contact" json:"lastContact"`
	IngestedAt       time.Time `bson:"ingested_at" json:"-"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updatedAt"`
}

func (a *AircraftLiveState) Longitude() float64 {
	if len(a.Location.Coordinates) >= 2 {
		return a.Location.Coordinates[0]
	}
	return 0
}

func (a *AircraftLiveState) Latitude() float64 {
	if len(a.Location.Coordinates) >= 2 {
		return a.Location.Coordinates[1]
	}
	return 0
}

type TrackPoint struct {
	TenantID       string       `bson:"tenant_id" json:"tenantId"`
	ICAO24         string       `bson:"icao24" json:"icao24"`
	Location       GeoJSONPoint `bson:"location" json:"location"`
	AltitudeM      *float64     `bson:"altitude_m" json:"altitude"`
	VelocityMS     *float64     `bson:"velocity_ms" json:"velocity"`
	HeadingDeg     *float64     `bson:"heading_deg" json:"heading"`
	VerticalRateMS *float64     `bson:"vertical_rate_ms" json:"verticalRate"`
	OnGround       bool         `bson:"on_ground" json:"onGround"`
	Timestamp      time.Time    `bson:"timestamp" json:"timestamp"`

	ExternalEventKey string `bson:"external_event_key" json:"-"`
}

type BBox struct {
	LaMin float64 `json:"lamin"`
	LoMin float64 `json:"lomin"`
	LaMax float64 `json:"lamax"`
	LoMax float64 `json:"lomax"`
}

type AircraftFilter struct {
	Callsign *string `json:"callsign,omitempty"`
	Country  *string `json:"country,omitempty"`
	OnGround *bool   `json:"onGround,omitempty"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
}
