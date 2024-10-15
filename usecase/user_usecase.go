package usecase

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/models"
	"github.com/kwa0x2/Settle-Backend/repository"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
	"time"
)

type IUserService interface {
	Create(user *models.User) error
	CreateAndJoinRoom(user *models.User, userRoom *models.UserRoom) error
}

type userService struct {
	UserRepository     repository.IUserRepository
	UserRoomRepository repository.IUserRoomRepository
}

func NewUserService(userRepository repository.IUserRepository, userRoomRepository repository.IUserRoomRepository) IUserService {
	return &userService{
		UserRepository:     userRepository,
		UserRoomRepository: userRoomRepository,
	}
}

func (s *userService) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	if err := user.Validate(); err != nil {
		return err
	}
	return s.UserRepository.Create(ctx, user)
}

func (s *userService) CreateAndJoinRoom(user *models.User, userRoom *models.UserRoom) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := s.UserRepository.GetDatabase().Client().StartSession()
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
		if userCreateErr := s.UserRepository.Create(txCtx, user); userCreateErr != nil {
			return nil, userCreateErr
		}

		userRoom.UserID = user.ID
		userRoom.CreatedAt = time.Now().UTC()
		userRoom.UpdatedAt = time.Now().UTC()
		userRoom.Visible = true
		if userRoomValidateErr := userRoom.Validate(); userRoomValidateErr != nil {
			return nil, userRoomValidateErr
		}
		if userRoomCreateErr := s.UserRoomRepository.Create(txCtx, userRoom); userRoomCreateErr != nil {
			return nil, userRoomCreateErr
		}

		return nil, nil
	}, txnOptions)

	return err
}
