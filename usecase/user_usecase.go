package usecase

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
	"time"
)

type userUsecase struct {
	userRepository     domain.UserRepository
	userRoomRepository domain.UserRoomRepository
}

func NewUserUsecase(userRepository domain.UserRepository, userRoomRepository domain.UserRoomRepository) domain.UserUsecase {
	return &userUsecase{
		userRepository:     userRepository,
		userRoomRepository: userRoomRepository,
	}
}

func (uu *userUsecase) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	if err := user.Validate(); err != nil {
		return err
	}
	return uu.userRepository.Create(ctx, user)
}

func (uu *userUsecase) CreateAndJoinRoom(user *domain.User, userRoom *domain.UserRoom) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := uu.userRepository.GetDatabase().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		user.CreatedAt = time.Now().UTC()
		user.UpdatedAt = time.Now().UTC()
		if validateErr := user.Validate(); validateErr != nil {
			return nil, validateErr
		}
		if userCreateErr := uu.userRepository.Create(txCtx, user); userCreateErr != nil {
			return nil, userCreateErr
		}

		userRoom.UserID = user.ID
		userRoom.CreatedAt = time.Now().UTC()
		userRoom.UpdatedAt = time.Now().UTC()
		userRoom.Visible = true
		if userRoomValidateErr := userRoom.Validate(); userRoomValidateErr != nil {
			return nil, userRoomValidateErr
		}
		result, userRoomCreateErr := uu.userRoomRepository.Create(txCtx, userRoom)
		if userRoomCreateErr != nil {
			return nil, userRoomCreateErr
		}

		userRoom.ID = result.InsertedID.(bson.ObjectID)

		return nil, nil
	}, txnOptions)

	return err
}
