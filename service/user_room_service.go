package service

import (
	"github.com/kwa0x2/Settle-Backend/models"
	"github.com/kwa0x2/Settle-Backend/repository"
	"time"
)

type IUserRoomService interface {
	Create(userRoom *models.UserRoom) error
}

type userRoomService struct {
	UserRoomRepository repository.IUserRoomRepository
}

func NewUserRoomService(userRoomRepository repository.IUserRoomRepository) IUserRoomService {
	return &userRoomService{
		UserRoomRepository: userRoomRepository,
	}
}

func (s *userRoomService) Create(userRoom *models.UserRoom) error {
	userRoom.CreatedAt = time.Now().UTC()
	userRoom.UpdatedAt = time.Now().UTC()
	userRoom.Visible = true
	return s.UserRoomRepository.Create(userRoom)
}
