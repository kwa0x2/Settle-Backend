package repository

import (
	"context"
	"errors"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (mr *messageRepository) UpdateByID(ctx context.Context, messageID bson.ObjectID, update bson.D) error {
	collection := mr.database.Collection(mr.collection)

	_, err := collection.UpdateByID(ctx, messageID, update)
	if err != nil {
		return err
	}

	return nil
}

func (mr *messageRepository) Find(ctx context.Context, filter bson.D, opts *options.FindOptionsBuilder) ([]domain.Message, error) {
	collection := mr.database.Collection(mr.collection)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []domain.Message

	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (mr *messageRepository) FindOne(ctx context.Context, filter bson.D) (domain.Message, error) {
	collection := mr.database.Collection(mr.collection)

	var message domain.Message
	err := collection.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return message, nil
		}
		return message, err
	}

	return message, nil
}

func (ur *messageRepository) GetDatabase() *mongo.Database {
	return ur.database
}
