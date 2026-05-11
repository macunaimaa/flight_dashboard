package repository

import (
	"context"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

type AircraftRepository interface {
	UpsertLiveState(ctx context.Context, state *domain.AircraftLiveState) error
	BulkUpsertLiveStates(ctx context.Context, states []domain.AircraftLiveState) error
	GetLiveStates(ctx context.Context, tenantID string, filter domain.AircraftFilter) ([]domain.AircraftLiveState, int64, error)
	GetLiveStateByICAO(ctx context.Context, tenantID string, icao24 string) (*domain.AircraftLiveState, error)
	GetLiveStatesInBBox(ctx context.Context, tenantID string, bbox domain.BBox, limit int) ([]domain.AircraftLiveState, error)
	InsertTrackPoints(ctx context.Context, points []domain.TrackPoint) error
	GetTrackPoints(ctx context.Context, tenantID string, icao24 string, from, to time.Time) ([]domain.TrackPoint, error)
	PurgeStaleStates(ctx context.Context, tenantID string, olderThan time.Time) (int64, error)
}

type TenantRepository interface {
	GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error)
	GetAllActive(ctx context.Context) ([]domain.Tenant, error)
	Create(ctx context.Context, tenant *domain.Tenant) error
	Update(ctx context.Context, tenant *domain.Tenant) error
}

type UserRepository interface {
	GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	UpdateLastLogin(ctx context.Context, id string, t time.Time) error
}

type AuditRepository interface {
	Insert(ctx context.Context, entry *domain.AuditEntry) error
	Query(ctx context.Context, tenantID string, limit, offset int) ([]domain.AuditEntry, error)
}

type CacheRepository interface {
	SetAircraftState(ctx context.Context, tenantID, icao24 string, state *domain.AircraftLiveState, ttl time.Duration) error
	GetAircraftState(ctx context.Context, tenantID, icao24 string) (*domain.AircraftLiveState, error)
	GetAllAircraftStates(ctx context.Context, tenantID string) (map[string]*domain.AircraftLiveState, error)
	DeleteAircraftState(ctx context.Context, tenantID, icao24 string) error
	GetTrackedICAOs(ctx context.Context, tenantID string) ([]string, error)
}
