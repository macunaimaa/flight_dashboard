package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

type MongoTenantRepo struct {
	tenants *mongo.Collection
	users   *mongo.Collection
}

func NewMongoTenantRepo(db *mongo.Database) *MongoTenantRepo {
	return &MongoTenantRepo{
		tenants: db.Collection("tenants"),
		users:   db.Collection("users"),
	}
}

// TenantRepository

func (r *MongoTenantRepo) GetByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.tenants.FindOne(ctx, bson.M{"tenant_id": tenantID}).Decode(&t)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	}
	return &t, err
}

func (r *MongoTenantRepo) GetAllActive(ctx context.Context) ([]domain.Tenant, error) {
	cursor, err := r.tenants.Find(ctx, bson.M{"active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tenants []domain.Tenant
	if err := cursor.All(ctx, &tenants); err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *MongoTenantRepo) Create(ctx context.Context, tenant *domain.Tenant) error {
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()
	_, err := r.tenants.InsertOne(ctx, tenant)
	return err
}

func (r *MongoTenantRepo) Update(ctx context.Context, tenant *domain.Tenant) error {
	tenant.UpdatedAt = time.Now()
	filter := bson.M{"tenant_id": tenant.TenantID}
	update := bson.M{"$set": tenant}
	_, err := r.tenants.UpdateOne(ctx, filter, update)
	return err
}

// UserRepository

func (r *MongoTenantRepo) GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	var u domain.User
	err := r.users.FindOne(ctx, bson.M{"tenant_id": tenantID, "email": email}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	}
	return &u, err
}

func (r *MongoTenantRepo) GetByUserID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := r.users.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	}
	return &u, err
}

func (r *MongoTenantRepo) CreateUser(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	_, err := r.users.InsertOne(ctx, user)
	return err
}

func (r *MongoTenantRepo) UpdateLastLogin(ctx context.Context, id string, t time.Time) error {
	_, err := r.users.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"last_login": t}},
		options.Update(),
	)
	return err
}
