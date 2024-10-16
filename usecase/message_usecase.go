package usecase

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/domain/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type messageUsecase struct {
	messageRepository domain.MessageRepository
}

func NewMessageUsecase(messageRepository domain.MessageRepository) domain.MessageUsecase {
	return &messageUsecase{
		messageRepository: messageRepository,
	}
}

func (mu *messageUsecase) Create(message *domain.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message.CreatedAt = time.Now().UTC()
	message.UpdatedAt = time.Now().UTC()
	message.ReadStatus = types.Unreaded
	if err := message.Validate(); err != nil {
		return err
	}
	result, err := mu.messageRepository.Create(ctx, message)
	if err != nil {
		return err
	}

	message.ID = primitive.ObjectID(result.InsertedID.(bson.ObjectID))

	return nil
}

func (mu *messageUsecase) SoftDelete(messageID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"room_id": "asdasda"}}
	_, err := mu.messageRepository.UpdateByID(ctx, messageID, update)
	if err != nil {
		return err
	}
	return nil
}
