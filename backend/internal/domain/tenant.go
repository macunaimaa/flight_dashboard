package domain

import "time"

type PollRegion struct {
	LaMin float64 `bson:"lamin" json:"lamin"`
	LoMin float64 `bson:"lomin" json:"lomin"`
	LaMax float64 `bson:"lamax" json:"lamax"`
	LoMax float64 `bson:"lomax" json:"lomax"`
}

type TenantConfig struct {
	PollRegions        []PollRegion `bson:"poll_regions" json:"pollRegions"`
	MaxTrackedAircraft int          `bson:"max_tracked_aircraft" json:"maxTrackedAircraft"`
	TrackRetentionDays int          `bson:"track_retention_days" json:"trackRetentionDays"`
}

type Tenant struct {
	TenantID  string       `bson:"tenant_id" json:"tenantId"`
	Name      string       `bson:"name" json:"name"`
	Config    TenantConfig `bson:"config" json:"config"`
	Active    bool         `bson:"active" json:"active"`
	CreatedAt time.Time    `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time    `bson:"updated_at" json:"updatedAt"`
}

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

type User struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	TenantID     string    `bson:"tenant_id" json:"tenantId"`
	Email        string    `bson:"email" json:"email"`
	PasswordHash string    `bson:"password_hash" json:"-"`
	Role         Role      `bson:"role" json:"role"`
	Active       bool      `bson:"active" json:"active"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	LastLogin    time.Time `bson:"last_login" json:"lastLogin"`
}
