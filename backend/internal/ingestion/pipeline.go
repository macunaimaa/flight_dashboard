package ingestion

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
	"github.com/macunaimaa/dashboard/backend/internal/repository"
	"github.com/macunaimaa/dashboard/backend/internal/ws"
)

// StatesFetcher abstracts the data source (OpenSky API or mock).
type StatesFetcher interface {
	FetchStates(ctx context.Context, bbox domain.BBox) ([]RawStateVector, int64, error)
}

type Pipeline struct {
	client       StatesFetcher
	aircraftRepo repository.AircraftRepository
	tenantRepo   repository.TenantRepository
	cache        repository.CacheRepository
	hub          *ws.Hub
	logger       *slog.Logger
	stopCh       chan struct{}

	mu            sync.Mutex
	backoffUntil  time.Time
	currentBackoff time.Duration
}

const (
	minBackoff = 30 * time.Second
	maxBackoff = 5 * time.Minute
)

func NewPipeline(
	client StatesFetcher,
	aircraftRepo repository.AircraftRepository,
	tenantRepo repository.TenantRepository,
	cache repository.CacheRepository,
	hub *ws.Hub,
	logger *slog.Logger,
) *Pipeline {
	return &Pipeline{
		client:       client,
		aircraftRepo: aircraftRepo,
		tenantRepo:   tenantRepo,
		cache:        cache,
		hub:          hub,
		logger:       logger.With("component", "ingestion"),
		stopCh:       make(chan struct{}),
		currentBackoff: minBackoff,
	}
}

func (p *Pipeline) Start(interval time.Duration) {
	p.logger.Info("ingestion pipeline started", "interval", interval.String())

	// Wait before first poll to let any previous rate limit cool down
	p.logger.Info("waiting 15s before first poll to avoid rate limit")
	select {
	case <-time.After(15 * time.Second):
	case <-p.stopCh:
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	p.pollAllTenants()

	for {
		select {
		case <-ticker.C:
			// Skip if we're in backoff
			p.mu.Lock()
			if time.Now().Before(p.backoffUntil) {
				remaining := time.Until(p.backoffUntil).Round(time.Second)
				p.mu.Unlock()
				p.logger.Info("skipping poll, in rate limit backoff", "remaining", remaining.String())
				continue
			}
			p.mu.Unlock()

			p.pollAllTenants()
		case <-p.stopCh:
			p.logger.Info("ingestion pipeline stopped")
			return
		}
	}
}

func (p *Pipeline) Stop() {
	close(p.stopCh)
}

func (p *Pipeline) setRateLimited() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.backoffUntil = time.Now().Add(p.currentBackoff)
	p.logger.Warn("rate limited, backing off",
		"backoff", p.currentBackoff.String(),
		"retry_at", p.backoffUntil.Format(time.TimeOnly),
	)
	// Exponential backoff
	p.currentBackoff = p.currentBackoff * 2
	if p.currentBackoff > maxBackoff {
		p.currentBackoff = maxBackoff
	}
}

func (p *Pipeline) resetBackoff() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.currentBackoff = minBackoff
	p.backoffUntil = time.Time{}
}

func (p *Pipeline) pollAllTenants() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tenants, err := p.tenantRepo.GetAllActive(ctx)
	if err != nil {
		p.logger.Error("failed to fetch tenants", "error", err)
		return
	}

	for _, tenant := range tenants {
		p.pollTenant(ctx, &tenant)
	}
}

func (p *Pipeline) pollTenant(ctx context.Context, tenant *domain.Tenant) {
	for _, region := range tenant.Config.PollRegions {
		bbox := domain.BBox{
			LaMin: region.LaMin,
			LoMin: region.LoMin,
			LaMax: region.LaMax,
			LoMax: region.LoMax,
		}

		start := time.Now()

		rawStates, apiTime, err := p.client.FetchStates(ctx, bbox)
		if err != nil {
			if strings.Contains(err.Error(), "429") {
				p.setRateLimited()
				return // Stop polling all regions for this cycle
			}
			p.logger.Warn("opensky fetch failed",
				"tenant_id", tenant.TenantID,
				"error", err,
			)
			continue
		}

		// Success — reset backoff
		p.resetBackoff()

		// Normalize all state vectors
		normalized := make([]domain.AircraftLiveState, 0, len(rawStates))
		for _, raw := range rawStates {
			state, err := Normalize(raw, apiTime)
			if err != nil {
				continue
			}
			state.TenantID = tenant.TenantID
			normalized = append(normalized, *state)
		}

		// Bulk upsert to MongoDB
		if err := p.aircraftRepo.BulkUpsertLiveStates(ctx, normalized); err != nil {
			p.logger.Error("bulk upsert failed",
				"tenant_id", tenant.TenantID,
				"error", err,
			)
		}

		// Compute deltas against Redis cache
		updates, removed := p.computeDeltas(ctx, tenant.TenantID, normalized)

		// Update Redis cache
		for i := range normalized {
			s := &normalized[i]
			_ = p.cache.SetAircraftState(ctx, tenant.TenantID, s.ICAO24, s, 90*time.Second)
		}

		// Insert track points for aircraft with changed positions
		if len(updates) > 0 {
			trackPoints := make([]domain.TrackPoint, 0, len(updates))
			for i := range updates {
				tp := ToTrackPoint(&updates[i])
				trackPoints = append(trackPoints, tp)
			}
			if err := p.aircraftRepo.InsertTrackPoints(ctx, trackPoints); err != nil {
				p.logger.Error("insert track points failed",
					"tenant_id", tenant.TenantID,
					"error", err,
				)
			}
		}

		// Broadcast updates via WebSocket
		if len(updates) > 0 {
			p.hub.BroadcastToTenant(tenant.TenantID, ws.NewAircraftUpdateMsg(updates))
		}
		if len(removed) > 0 {
			p.hub.BroadcastToTenant(tenant.TenantID, ws.NewAircraftRemoveMsg(removed))
			for _, icao := range removed {
				_ = p.cache.DeleteAircraftState(ctx, tenant.TenantID, icao)
			}
			_, _ = p.aircraftRepo.PurgeStaleStates(ctx, tenant.TenantID, time.Now().Add(-90*time.Second))
		}

		duration := time.Since(start)
		p.logger.Info("poll complete",
			"tenant_id", tenant.TenantID,
			"duration_ms", duration.Milliseconds(),
			"raw_count", len(rawStates),
			"normalized", len(normalized),
			"updates", len(updates),
			"removed", len(removed),
		)
	}
}

// computeDeltas compares current states against Redis cache.
func (p *Pipeline) computeDeltas(ctx context.Context, tenantID string, current []domain.AircraftLiveState) ([]domain.AircraftLiveState, []string) {
	cached, err := p.cache.GetAllAircraftStates(ctx, tenantID)
	if err != nil {
		p.logger.Warn("cache read failed, treating all as updates", "error", err)
		return current, nil
	}

	currentICAOs := make(map[string]bool, len(current))
	updates := make([]domain.AircraftLiveState, 0)

	for _, state := range current {
		currentICAOs[state.ICAO24] = true

		prev, exists := cached[state.ICAO24]
		if !exists {
			updates = append(updates, state)
			continue
		}

		if positionChanged(prev, &state) {
			updates = append(updates, state)
		}
	}

	var removed []string
	for icao := range cached {
		if !currentICAOs[icao] {
			removed = append(removed, icao)
		}
	}

	return updates, removed
}

func positionChanged(prev, curr *domain.AircraftLiveState) bool {
	if prev.Longitude() != curr.Longitude() || prev.Latitude() != curr.Latitude() {
		return true
	}
	if !floatPtrEqual(prev.BaroAltitudeM, curr.BaroAltitudeM) {
		return true
	}
	if !floatPtrEqual(prev.VelocityMS, curr.VelocityMS) {
		return true
	}
	if !floatPtrEqual(prev.TrueTrackDeg, curr.TrueTrackDeg) {
		return true
	}
	if prev.OnGround != curr.OnGround {
		return true
	}
	return false
}

func floatPtrEqual(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
