package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, uri string, logger *slog.Logger) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	logger.Info("connected to MongoDB", "uri", uri)
	return client, nil
}

func EnsureIndexes(ctx context.Context, db *mongo.Database, logger *slog.Logger) error {
	indexes := map[string][]mongo.IndexModel{
		"aircraft_live_state": {
			{
				Keys:    bson.D{{Key: "tenant_id", Value: 1}, {Key: "icao24", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "location", Value: "2dsphere"}},
			},
			{
				Keys:    bson.D{{Key: "external_event_key", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "updated_at", Value: 1}},
			},
			{
				Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "callsign", Value: 1}},
			},
		},
		"aircraft_track_points": {
			{
				Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "icao24", Value: 1}, {Key: "timestamp", Value: -1}},
			},
			{
				Keys:    bson.D{{Key: "timestamp", Value: 1}},
				Options: options.Index().SetExpireAfterSeconds(30 * 24 * 3600), // 30 days TTL
			},
			{
				Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "timestamp", Value: 1}},
			},
			{
				Keys:    bson.D{{Key: "external_event_key", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		"users": {
			{
				Keys:    bson.D{{Key: "tenant_id", Value: 1}, {Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		"tenants": {
			{
				Keys:    bson.D{{Key: "tenant_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		"audit_logs": {
			{
				Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "timestamp", Value: -1}},
			},
			{
				Keys:    bson.D{{Key: "timestamp", Value: 1}},
				Options: options.Index().SetExpireAfterSeconds(90 * 24 * 3600), // 90 days TTL
			},
		},
	}

	for collection, models := range indexes {
		coll := db.Collection(collection)
		names, err := coll.Indexes().CreateMany(ctx, models)
		if err != nil {
			return fmt.Errorf("create indexes on %s: %w", collection, err)
		}
		logger.Info("indexes ensured", "collection", collection, "indexes", names)
	}

	return nil
}
