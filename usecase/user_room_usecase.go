package usecase

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"time"
)

type userRoomUsecase struct {
	userRoomRepository domain.UserRoomRepository
}

func NewUserRoomUsecase(userRoomRepository domain.UserRoomRepository) domain.UserRoomUsecase {
	return &userRoomUsecase{
		userRoomRepository: userRoomRepository,
	}
}

func (uru *userRoomUsecase) Create(userRoom *domain.UserRoom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userRoom.CreatedAt = time.Now().UTC()
	userRoom.UpdatedAt = time.Now().UTC()
	userRoom.Visible = true
	if err := userRoom.Validate(); err != nil {
		return err
	}
	return uru.userRoomRepository.Create(ctx, userRoom)
}
