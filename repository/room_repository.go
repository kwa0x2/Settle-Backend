package repository

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type roomRepository struct {
	database   *mongo.Database
	collection string
}

func NewRoomRepository(db *mongo.Database, collection string) domain.RoomRepository {
	return &roomRepository{
		database:   db,
		collection: collection,
	}
}

func (rr *roomRepository) Find(ctx context.Context, filter bson.D, opts *options.FindOptionsBuilder) ([]domain.Room, error) {
	collection := rr.database.Collection(rr.collection)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []domain.Room

	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (rr *roomRepository) UpdateByID(ctx context.Context, roomID bson.ObjectID, update bson.D) error {
	collection := rr.database.Collection(rr.collection)

	_, err := collection.UpdateByID(ctx, roomID, update)
	if err != nil {
		return err
	}

	return nil
}
