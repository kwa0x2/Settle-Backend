package repository

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRoomRepository struct {
	database   *mongo.Database
	collection string
}

func NewUserRoomRepository(db *mongo.Database, collection string) domain.UserRoomRepository {
	return &userRoomRepository{
		database:   db,
		collection: collection,
	}
}

func (urr *userRoomRepository) Create(ctx context.Context, userRoom *domain.UserRoom) error {
	collection := urr.database.Collection(urr.collection)

	_, err := collection.InsertOne(ctx, userRoom)
	return err
}
