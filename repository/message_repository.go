package repository

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type messageRepository struct {
	database   *mongo.Database
	collection string
}

func NewMessageRepository(db *mongo.Database, collection string) domain.MessageRepository {
	return &messageRepository{
		database:   db,
		collection: collection,
	}
}

func (mr *messageRepository) Create(ctx context.Context, message *domain.Message) (*mongo.InsertOneResult, error) {
	collection := mr.database.Collection(mr.collection)

	return collection.InsertOne(ctx, message)
}
