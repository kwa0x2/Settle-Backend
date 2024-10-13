package repository

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"time"
)

type IUserRepository interface {
	Create(user *models.User) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) IUserRepository {
	return &userRepository{
		collection: database.Collection("users"),
	}
}

func (r *userRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	_, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return err
	}
	return nil
}
