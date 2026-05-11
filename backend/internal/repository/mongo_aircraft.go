package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

type MongoAircraftRepo struct {
	liveStates  *mongo.Collection
	trackPoints *mongo.Collection
}

func NewMongoAircraftRepo(db *mongo.Database) *MongoAircraftRepo {
	return &MongoAircraftRepo{
		liveStates:  db.Collection("aircraft_live_state"),
		trackPoints: db.Collection("aircraft_track_points"),
	}
}

func (r *MongoAircraftRepo) UpsertLiveState(ctx context.Context, state *domain.AircraftLiveState) error {
	filter := bson.M{
		"tenant_id": state.TenantID,
		"icao24":    state.ICAO24,
	}
	update := bson.M{"$set": state}
	opts := options.Update().SetUpsert(true)

	_, err := r.liveStates.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoAircraftRepo) BulkUpsertLiveStates(ctx context.Context, states []domain.AircraftLiveState) error {
	if len(states) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(states))
	for i := range states {
		s := &states[i]
		s.UpdatedAt = time.Now()
		if s.IngestedAt.IsZero() {
			s.IngestedAt = time.Now()
		}

		filter := bson.M{
			"tenant_id": s.TenantID,
			"icao24":    s.ICAO24,
		}
		update := bson.M{"$set": s}

		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		models = append(models, model)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.liveStates.BulkWrite(ctx, models, opts)
	return err
}

func (r *MongoAircraftRepo) GetLiveStates(ctx context.Context, tenantID string, filter domain.AircraftFilter) ([]domain.AircraftLiveState, int64, error) {
	query := bson.M{"tenant_id": tenantID}

	if filter.Callsign != nil && *filter.Callsign != "" {
		query["callsign"] = bson.M{"$regex": *filter.Callsign, "$options": "i"}
	}
	if filter.Country != nil && *filter.Country != "" {
		query["origin_country"] = bson.M{"$regex": *filter.Country, "$options": "i"}
	}
	if filter.OnGround != nil {
		query["on_ground"] = *filter.OnGround
	}

	total, err := r.liveStates.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "callsign", Value: 1}})

	cursor, err := r.liveStates.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []domain.AircraftLiveState
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *MongoAircraftRepo) GetLiveStateByICAO(ctx context.Context, tenantID string, icao24 string) (*domain.AircraftLiveState, error) {
	filter := bson.M{
		"tenant_id": tenantID,
		"icao24":    icao24,
	}

	var state domain.AircraftLiveState
	err := r.liveStates.FindOne(ctx, filter).Decode(&state)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *MongoAircraftRepo) GetLiveStatesInBBox(ctx context.Context, tenantID string, bbox domain.BBox, limit int) ([]domain.AircraftLiveState, error) {
	if limit <= 0 || limit > 5000 {
		limit = 1000
	}

	filter := bson.M{
		"tenant_id": tenantID,
		"location": bson.M{
			"$geoWithin": bson.M{
				"$geometry": bson.M{
					"type": "Polygon",
					"coordinates": bson.A{
						bson.A{
							bson.A{bbox.LoMin, bbox.LaMin},
							bson.A{bbox.LoMax, bbox.LaMin},
							bson.A{bbox.LoMax, bbox.LaMax},
							bson.A{bbox.LoMin, bbox.LaMax},
							bson.A{bbox.LoMin, bbox.LaMin},
						},
					},
				},
			},
		},
	}

	opts := options.Find().SetLimit(int64(limit))
	cursor, err := r.liveStates.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.AircraftLiveState
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoAircraftRepo) InsertTrackPoints(ctx context.Context, points []domain.TrackPoint) error {
	if len(points) == 0 {
		return nil
	}

	docs := make([]interface{}, len(points))
	for i := range points {
		docs[i] = points[i]
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := r.trackPoints.InsertMany(ctx, docs, opts)
	if mongo.IsDuplicateKeyError(err) {
		return nil // idempotent: duplicates are expected
	}
	return err
}

func (r *MongoAircraftRepo) GetTrackPoints(ctx context.Context, tenantID string, icao24 string, from, to time.Time) ([]domain.TrackPoint, error) {
	filter := bson.M{
		"tenant_id": tenantID,
		"icao24":    icao24,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: 1}}).
		SetLimit(10000)

	cursor, err := r.trackPoints.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.TrackPoint
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoAircraftRepo) PurgeStaleStates(ctx context.Context, tenantID string, olderThan time.Time) (int64, error) {
	filter := bson.M{
		"tenant_id":  tenantID,
		"updated_at": bson.M{"$lt": olderThan},
	}
	result, err := r.liveStates.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
