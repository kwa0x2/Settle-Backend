package usecase

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
	"time"
)

type messageUsecase struct {
	messageRepository domain.MessageRepository
	roomRepository    domain.RoomRepository
	userRepository    domain.UserRepository
}

func NewMessageUsecase(messageRepository domain.MessageRepository, roomRepository domain.RoomRepository, userRepository domain.UserRepository) domain.MessageUsecase {
	return &messageUsecase{
		messageRepository: messageRepository,
		roomRepository:    roomRepository,
		userRepository:    userRepository,
	}
}

func (mu *messageUsecase) CreateAndUpdateRoom(message *domain.Message) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := mu.messageRepository.GetDatabase().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		message.CreatedAt = time.Now().UTC()
		message.UpdatedAt = time.Now().UTC()
		message.ReadStatus = 0 //unseen
		if validateErr := message.Validate(); validateErr != nil {
			return nil, validateErr
		}
		result, createErr := mu.messageRepository.Create(ctx, message)
		if createErr != nil {
			return nil, createErr
		}

		message.ID = result.InsertedID.(bson.ObjectID)

		update := bson.D{{"$set", bson.D{{"last_message", message}}}}
		if updateErr := mu.roomRepository.UpdateByID(ctx, message.RoomID, update); updateErr != nil {
			return nil, updateErr
		}

		return nil, nil

	}, txnOptions)

	return err
}

func (mu *messageUsecase) SoftDelete(messageID bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.D{{"$set", bson.D{{"deleted_at", time.Now().UTC()}}}}
	if err := mu.messageRepository.UpdateByID(ctx, messageID, update); err != nil {
		return err
	}
	return nil
}

func (mu *messageUsecase) GetByRoomID(roomID bson.ObjectID, limit, offset int64) ([]domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{"created_at", 1}}).SetLimit(limit).SetSkip(offset)
	filter := bson.D{{"room_id", roomID}}
	result, err := mu.messageRepository.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	//for i, message := range result {
	//	filter := bson.D{{"_id", message.SenderID}}
	//	user, err := mu.userRepository.FindOne(ctx, filter)
	//	if err != nil {
	//		return nil, err
	//	}
	//	result[i].User = &user
	//}

	return result, err
}

func (mu *messageUsecase) GetByID(messageID bson.ObjectID) (domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"_id", messageID}}
	result, err := mu.messageRepository.FindOne(ctx, filter)
	if err != nil {
		return domain.Message{}, err
	}
	return result, nil
}

func (mu *messageUsecase) EditMessage(messageID bson.ObjectID, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.D{{"$set", bson.D{{"content", content}, {"updated_at", time.Now().UTC()}}}}
	err := mu.messageRepository.UpdateByID(ctx, messageID, update)
	if err != nil {
		return err
	}
	return nil
}
