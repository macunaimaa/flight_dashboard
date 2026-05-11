package domain

import "time"

type GeoJSONPolygon struct {
	Type        string        `bson:"type" json:"type"`
	Coordinates [][][]float64 `bson:"coordinates" json:"coordinates"`
}

type Geofence struct {
	ID           string         `bson:"_id,omitempty" json:"id"`
	TenantID     string         `bson:"tenant_id" json:"tenantId"`
	Name         string         `bson:"name" json:"name"`
	Geometry     GeoJSONPolygon `bson:"geometry" json:"geometry"`
	AltitudeMinM float64        `bson:"altitude_min_m" json:"altitudeMinM"`
	AltitudeMaxM float64        `bson:"altitude_max_m" json:"altitudeMaxM"`
	Active       bool           `bson:"active" json:"active"`
	CreatedAt    time.Time      `bson:"created_at" json:"createdAt"`
}
