package service

import (
	"context"
	"time"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
	"github.com/macunaimaa/dashboard/backend/internal/repository"
)

type AircraftService struct {
	repo  repository.AircraftRepository
	cache repository.CacheRepository
}

func NewAircraftService(repo repository.AircraftRepository, cache repository.CacheRepository) *AircraftService {
	return &AircraftService{repo: repo, cache: cache}
}

func (s *AircraftService) GetLiveStates(ctx context.Context, tenantID string, filter domain.AircraftFilter) ([]domain.AircraftLiveState, int64, error) {
	return s.repo.GetLiveStates(ctx, tenantID, filter)
}

func (s *AircraftService) GetByICAO(ctx context.Context, tenantID, icao24 string) (*domain.AircraftLiveState, error) {
	// Try cache first
	cached, err := s.cache.GetAircraftState(ctx, tenantID, icao24)
	if err == nil && cached != nil {
		return cached, nil
	}
	return s.repo.GetLiveStateByICAO(ctx, tenantID, icao24)
}

func (s *AircraftService) GetInBBox(ctx context.Context, tenantID string, bbox domain.BBox, limit int) ([]domain.AircraftLiveState, error) {
	return s.repo.GetLiveStatesInBBox(ctx, tenantID, bbox, limit)
}

func (s *AircraftService) GetTrackPoints(ctx context.Context, tenantID, icao24 string, from, to time.Time) ([]domain.TrackPoint, error) {
	return s.repo.GetTrackPoints(ctx, tenantID, icao24, from, to)
}

func (s *AircraftService) BulkUpsert(ctx context.Context, states []domain.AircraftLiveState) error {
	return s.repo.BulkUpsertLiveStates(ctx, states)
}
