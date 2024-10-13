package service

import (
	"github.com/kwa0x2/Settle-Backend/models"
	"github.com/kwa0x2/Settle-Backend/repository"
)

type IUserService interface {
	Create(user *models.User) error
}

type userService struct {
	UserRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &userService{
		UserRepository: userRepository,
	}
}

func (s *userService) Create(user *models.User) error {
	return s.UserRepository.Create(user)
}
