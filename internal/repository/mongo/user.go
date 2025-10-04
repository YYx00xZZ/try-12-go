package mongorepo

import (
	"context"
	"log/slog"

	"github.com/YYx00xZZ/try-12-go/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userDocument struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

// UserRepository provides Mongo-backed access to users.
type UserRepository struct {
	collection *mongodriver.Collection
}

// NewUserRepository instantiates a Mongo user repository.
func NewUserRepository(collection *mongodriver.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

// List fetches up to 10 users ordered by their id.
func (r *UserRepository) List(ctx context.Context) ([]repository.User, error) {
	findOpts := options.Find().SetLimit(10).SetSort(bson.D{{Key: "id", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.D{}, findOpts)
	if err != nil {
		slog.Error("mongo find users failed", slog.Any("err", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make([]repository.User, 0)
	for cursor.Next(ctx) {
		var doc userDocument
		if err := cursor.Decode(&doc); err != nil {
			slog.Error("decode mongo user failed", slog.Any("err", err))
			return nil, err
		}
		users = append(users, repository.User{ID: doc.ID, Name: doc.Name})
	}

	if err := cursor.Err(); err != nil {
		slog.Error("mongo cursor error", slog.Any("err", err))
		return nil, err
	}

	slog.Info("fetched users from mongo", slog.Int("count", len(users)))

	return users, nil
}
