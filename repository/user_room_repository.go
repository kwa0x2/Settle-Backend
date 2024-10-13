package repository

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IUserRoomRepository interface {
	Create(userRoom *models.UserRoom) error
}

type userRoomRepository struct {
	Collection *mongo.Collection
}

func NewUserRoomRepository(database *mongo.Database) IUserRoomRepository {
	return &userRoomRepository{
		Collection: database.Collection("user_rooms"),
	}
}

func (r *userRoomRepository) Create(userRoom *models.UserRoom) error {
	_, err := r.Collection.InsertOne(context.TODO(), userRoom)
	return err
}
