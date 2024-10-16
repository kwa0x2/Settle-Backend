package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/api/http/delivery"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/repository"
	"github.com/kwa0x2/Settle-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewAuthRoute(env *bootstrap.Env, db *mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	urr := repository.NewUserRoomRepository(db, domain.CollectionUserRoom)
	ad := &delivery.AuthDelivery{
		UserUsecase: usecase.NewUserUsecase(ur, urr),
		Env:         env,
	}

	group.POST("auth/login", ad.Login)
	group.GET("auth/refresh", ad.RefreshToken)

}
