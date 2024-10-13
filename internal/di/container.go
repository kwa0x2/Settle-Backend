package di

import (
	"github.com/kwa0x2/Settle-Backend/controller"
	"github.com/kwa0x2/Settle-Backend/repository"
	"github.com/kwa0x2/Settle-Backend/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Container struct {
	AuthController controller.IAuthController
}

func NewContainer(database *mongo.Database) *Container {

	userRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepository)

	userRoomRepository := repository.NewUserRoomRepository(database)
	userRoomService := service.NewUserRoomService(userRoomRepository)

	return &Container{
		AuthController: controller.NewAuthController(userService, userRoomService),
	}
}
