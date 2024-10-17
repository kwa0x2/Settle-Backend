package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	message.ID = result.InsertedID.(bson.ObjectID)

	return nil
}

func (mu *messageUsecase) SoftDelete(messageID bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.D{{"$set", bson.D{{"deleted_at", time.Now().UTC()}}}}
	_, err := mu.messageRepository.UpdateByID(ctx, messageID, update)
	if err != nil {
		return err
	}
	return nil
}

func (mu *messageUsecase) GetByRoomID(roomID uuid.UUID) ([]domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{"created_at", 1}})
	filter := bson.D{{"room_id", roomID.String()}}
	result, err := mu.messageRepository.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return result, err
}
