package domain

import "time"

type AlertType string

const (
	AlertAltitudeThreshold AlertType = "altitude_threshold"
	AlertSpeedThreshold    AlertType = "speed_threshold"
	AlertGeofenceEnter     AlertType = "geofence_enter"
	AlertGeofenceExit      AlertType = "geofence_exit"
	AlertSquawkMatch       AlertType = "squawk_match"
)

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

type AlertCondition struct {
	Operator string  `bson:"operator" json:"operator"` // lt, gt, eq
	Value    float64 `bson:"value" json:"value"`
	Unit     string  `bson:"unit" json:"unit"`
}

type Alert struct {
	ID          string         `bson:"_id,omitempty" json:"id"`
	TenantID    string         `bson:"tenant_id" json:"tenantId"`
	Name        string         `bson:"name" json:"name"`
	Type        AlertType      `bson:"type" json:"type"`
	Condition   AlertCondition `bson:"condition" json:"condition"`
	GeofenceID  *string        `bson:"geofence_id,omitempty" json:"geofenceId,omitempty"`
	Severity    Severity       `bson:"severity" json:"severity"`
	Active      bool           `bson:"active" json:"active"`
	CreatedAt   time.Time      `bson:"created_at" json:"createdAt"`
}

type AlertEvent struct {
	ID             string            `bson:"_id,omitempty" json:"id"`
	TenantID       string            `bson:"tenant_id" json:"tenantId"`
	AlertID        string            `bson:"alert_id" json:"alertId"`
	ICAO24         string            `bson:"icao24" json:"icao24"`
	Callsign       string            `bson:"callsign" json:"callsign"`
	TriggeredAt    time.Time         `bson:"triggered_at" json:"triggeredAt"`
	Location       GeoJSONPoint      `bson:"location" json:"location"`
	Snapshot       AircraftLiveState `bson:"snapshot" json:"snapshot"`
	Acknowledged   bool              `bson:"acknowledged" json:"acknowledged"`
	AcknowledgedBy *string           `bson:"acknowledged_by,omitempty" json:"acknowledgedBy,omitempty"`
}
