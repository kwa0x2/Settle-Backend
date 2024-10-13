package di

import "github.com/kwa0x2/Settle-Backend/controller"

type Container struct {
	AuthController controller.IAuthController
}

func NewContainer() *Container {

	return &Container{
		AuthController: controller.NewAuthController(),
	}
}
