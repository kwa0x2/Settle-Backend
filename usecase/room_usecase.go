package usecase

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type roomUsecase struct {
	roomRepository    domain.RoomRepository
	messageRepository domain.MessageRepository
}

func NewRoomUsecase(roomRepository domain.RoomRepository, messageRepository domain.MessageRepository) domain.RoomUsecase {
	return &roomUsecase{
		roomRepository:    roomRepository,
		messageRepository: messageRepository,
	}
}

func (ru *roomUsecase) FindAll() ([]domain.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{"updated_at", 1}})
	result, err := ru.roomRepository.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	for i, room := range result {
		filter := bson.D{{"_id", room.LastMessageID}}
		lastMessage, err := ru.messageRepository.FindOne(ctx, filter)
		if err != nil {
			return nil, err
		}
		result[i].LastMessage = &lastMessage
	}

	return result, err
}
