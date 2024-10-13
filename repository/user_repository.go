package repository

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IUserRepository interface {
	Create(user *models.User) error
}

type userRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) IUserRepository {
	return &userRepository{
		Collection: database.Collection("users"),
	}
}

func (r *userRepository) Create(user *models.User) error {
	_, err := r.Collection.InsertOne(context.TODO(), user)
	return err
}
