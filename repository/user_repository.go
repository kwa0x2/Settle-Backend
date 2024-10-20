package repository

import (
	"context"
	"errors"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepository struct {
	database   *mongo.Database
	collection string
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	collection := ur.database.Collection(ur.collection)

	_, err := collection.InsertOne(ctx, user)
	return err
}

func (ur *userRepository) FindOne(ctx context.Context, filter bson.D) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	var user domain.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, nil
		}
		return user, err
	}

	return user, nil
}

func (ur *userRepository) GetDatabase() *mongo.Database {
	return ur.database
}
