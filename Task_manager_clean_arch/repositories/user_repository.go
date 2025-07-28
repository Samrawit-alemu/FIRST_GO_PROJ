package repositories

import (
	"context"
	"taskmanager/domain"
	datamodels "taskmanager/repositories/models" // Aliased to avoid name conflicts

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Count(ctx context.Context) (int64, error)
}

// mongoUserRepository is the concrete implementation.
type mongoUserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository is the constructor.
func NewUserRepository(db *mongo.Database) IUserRepository {
	collection := db.Collection("users")
	// Ensure Username is unique
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}
	_, _ = collection.Indexes().CreateOne(context.Background(), indexModel)
	return &mongoUserRepository{collection: collection}
}

// toBsonUser converts a pure domain.User into a BSON-tagged datamodels.User.
func toBsonUser(user *domain.User) *datamodels.User {
	return &datamodels.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
}

// toDomainUser converts a BSON-tagged datamodels.User into a pure domain.User.
func toDomainUser(user *datamodels.User) *domain.User {
	return &domain.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *domain.User) error {
	// 1. Map from Domain model to BSON model before inserting.
	bsonUser := toBsonUser(user)

	result, err := r.collection.InsertOne(ctx, bsonUser)
	if err != nil {
		return err
	}

	// 2. Set the generated ID back on the original domain object.
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var bsonUser datamodels.User

	// 1. Perform the database find operation.
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&bsonUser)
	if err != nil {
		return nil, err // Can be mongo.ErrNoDocuments
	}

	// 2. Map the result from a BSON model back to a Domain model before returning.
	return toDomainUser(&bsonUser), nil
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var bsonUser datamodels.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bsonUser)
	if err != nil {
		return nil, err
	}
	return toDomainUser(&bsonUser), nil
}

func (r *mongoUserRepository) Update(ctx context.Context, user *domain.User) error {
	bsonUser := toBsonUser(user)
	filter := bson.M{"_id": bsonUser.ID}
	update := bson.M{"$set": bsonUser} // Use $set to update all fields in the BSON model
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoUserRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
