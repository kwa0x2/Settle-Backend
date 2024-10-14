package repository

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetDatabase() *mongo.Database
}

type userRepository struct {
	Collection *mongo.Collection
	Database   *mongo.Database
}

func NewUserRepository(database *mongo.Database) IUserRepository {
	return &userRepository{
		Collection: database.Collection("users"),
		Database:   database,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) GetDatabase() *mongo.Database {
	return r.Database
}
