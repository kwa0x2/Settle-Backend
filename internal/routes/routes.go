package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/controller"
)

func AuthRoute(router *gin.Engine, authController controller.IAuthController) {
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.GET("login/steam", authController.SteamLogin)
		authRoutes.GET("login/steam/callback", authController.SteamCallback)
	}
}
