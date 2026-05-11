package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func cacheKey(tenantID, icao24 string) string {
	return fmt.Sprintf("aircraft:%s:%s", tenantID, icao24)
}

func tenantSetKey(tenantID string) string {
	return fmt.Sprintf("aircraft_set:%s", tenantID)
}

func (c *RedisCache) SetAircraftState(ctx context.Context, tenantID, icao24 string, state *domain.AircraftLiveState, ttl time.Duration) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	pipe := c.client.Pipeline()
	pipe.Set(ctx, cacheKey(tenantID, icao24), data, ttl)
	pipe.SAdd(ctx, tenantSetKey(tenantID), icao24)
	pipe.Expire(ctx, tenantSetKey(tenantID), ttl+30*time.Second)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *RedisCache) GetAircraftState(ctx context.Context, tenantID, icao24 string) (*domain.AircraftLiveState, error) {
	data, err := c.client.Get(ctx, cacheKey(tenantID, icao24)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var state domain.AircraftLiveState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (c *RedisCache) GetAllAircraftStates(ctx context.Context, tenantID string) (map[string]*domain.AircraftLiveState, error) {
	icaos, err := c.client.SMembers(ctx, tenantSetKey(tenantID)).Result()
	if err != nil {
		return nil, err
	}

	if len(icaos) == 0 {
		return make(map[string]*domain.AircraftLiveState), nil
	}

	keys := make([]string, len(icaos))
	for i, icao := range icaos {
		keys[i] = cacheKey(tenantID, icao)
	}

	results, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	states := make(map[string]*domain.AircraftLiveState, len(results))
	for i, val := range results {
		if val == nil {
			continue
		}
		str, ok := val.(string)
		if !ok {
			continue
		}
		var state domain.AircraftLiveState
		if err := json.Unmarshal([]byte(str), &state); err != nil {
			continue
		}
		states[icaos[i]] = &state
	}
	return states, nil
}

func (c *RedisCache) DeleteAircraftState(ctx context.Context, tenantID, icao24 string) error {
	pipe := c.client.Pipeline()
	pipe.Del(ctx, cacheKey(tenantID, icao24))
	pipe.SRem(ctx, tenantSetKey(tenantID), icao24)
	_, err := pipe.Exec(ctx)
	return err
}

func (c *RedisCache) GetTrackedICAOs(ctx context.Context, tenantID string) ([]string, error) {
	return c.client.SMembers(ctx, tenantSetKey(tenantID)).Result()
}
