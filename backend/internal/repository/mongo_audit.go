package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

type MongoAuditRepo struct {
	coll *mongo.Collection
}

func NewMongoAuditRepo(db *mongo.Database) *MongoAuditRepo {
	return &MongoAuditRepo{
		coll: db.Collection("audit_logs"),
	}
}

func (r *MongoAuditRepo) Insert(ctx context.Context, entry *domain.AuditEntry) error {
	_, err := r.coll.InsertOne(ctx, entry)
	return err
}

func (r *MongoAuditRepo) Query(ctx context.Context, tenantID string, limit, offset int) ([]domain.AuditEntry, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.coll.Find(ctx, bson.M{"tenant_id": tenantID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.AuditEntry
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
