package handler

import (
	"context"
	"sync"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

// --- AircraftRepository fake ---

type fakeAircraftRepo struct {
	mu sync.Mutex

	listStates []domain.AircraftLiveState
	listTotal  int64
	listErr    error
	listGotFilter domain.AircraftFilter

	byICAO      *domain.AircraftLiveState
	byICAOErr   error
	byICAOCalls []struct{ tenant, icao string }

	bboxStates []domain.AircraftLiveState
	bboxErr    error
	bboxCalls  []struct {
		tenant string
		bbox   domain.BBox
		limit  int
	}

	trackPoints []domain.TrackPoint
	trackErr    error
	trackCalls  []struct {
		tenant, icao string
		from, to     time.Time
	}
}

func (r *fakeAircraftRepo) UpsertLiveState(ctx context.Context, state *domain.AircraftLiveState) error {
	return nil
}
func (r *fakeAircraftRepo) BulkUpsertLiveStates(ctx context.Context, states []domain.AircraftLiveState) error {
	return nil
}
func (r *fakeAircraftRepo) GetLiveStates(ctx context.Context, tenantID string, filter domain.AircraftFilter) ([]domain.AircraftLiveState, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listGotFilter = filter
	return r.listStates, r.listTotal, r.listErr
}
func (r *fakeAircraftRepo) GetLiveStateByICAO(ctx context.Context, tenantID, icao24 string) (*domain.AircraftLiveState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byICAOCalls = append(r.byICAOCalls, struct{ tenant, icao string }{tenantID, icao24})
	return r.byICAO, r.byICAOErr
}
func (r *fakeAircraftRepo) GetLiveStatesInBBox(ctx context.Context, tenantID string, bbox domain.BBox, limit int) ([]domain.AircraftLiveState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bboxCalls = append(r.bboxCalls, struct {
		tenant string
		bbox   domain.BBox
		limit  int
	}{tenantID, bbox, limit})
	return r.bboxStates, r.bboxErr
}
func (r *fakeAircraftRepo) InsertTrackPoints(ctx context.Context, points []domain.TrackPoint) error {
	return nil
}
func (r *fakeAircraftRepo) GetTrackPoints(ctx context.Context, tenantID, icao24 string, from, to time.Time) ([]domain.TrackPoint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trackCalls = append(r.trackCalls, struct {
		tenant, icao string
		from, to     time.Time
	}{tenantID, icao24, from, to})
	return r.trackPoints, r.trackErr
}
func (r *fakeAircraftRepo) PurgeStaleStates(ctx context.Context, tenantID string, olderThan time.Time) (int64, error) {
	return 0, nil
}

// --- CacheRepository fake ---

type fakeCacheRepo struct {
	mu       sync.Mutex
	state    *domain.AircraftLiveState
	stateErr error
}

func (c *fakeCacheRepo) SetAircraftState(ctx context.Context, tenantID, icao24 string, state *domain.AircraftLiveState, ttl time.Duration) error {
	return nil
}
func (c *fakeCacheRepo) GetAircraftState(ctx context.Context, tenantID, icao24 string) (*domain.AircraftLiveState, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state, c.stateErr
}
func (c *fakeCacheRepo) GetAllAircraftStates(ctx context.Context, tenantID string) (map[string]*domain.AircraftLiveState, error) {
	return nil, nil
}
func (c *fakeCacheRepo) DeleteAircraftState(ctx context.Context, tenantID, icao24 string) error {
	return nil
}
func (c *fakeCacheRepo) GetTrackedICAOs(ctx context.Context, tenantID string) ([]string, error) {
	return nil, nil
}

// --- AuditRepository fake (audit svc spawns a goroutine — must be race-safe) ---

type fakeAuditRepo struct {
	mu      sync.Mutex
	entries []*domain.AuditEntry
}

func (a *fakeAuditRepo) Insert(ctx context.Context, entry *domain.AuditEntry) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, entry)
	return nil
}
func (a *fakeAuditRepo) Query(ctx context.Context, tenantID string, limit, offset int) ([]domain.AuditEntry, error) {
	return nil, nil
}

// --- AuthUserRepo fake ---

type fakeAuthUserRepo struct {
	mu sync.Mutex

	user    *domain.User
	userErr error

	updateCalls []struct {
		id string
		t  time.Time
	}
}

func (u *fakeAuthUserRepo) GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.user, u.userErr
}
func (u *fakeAuthUserRepo) UpdateLastLogin(ctx context.Context, id string, t time.Time) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.updateCalls = append(u.updateCalls, struct {
		id string
		t  time.Time
	}{id, t})
	return nil
}
