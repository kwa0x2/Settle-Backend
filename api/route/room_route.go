package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/api/http/delivery"
	"github.com/kwa0x2/Settle-Backend/api/middleware"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/repository"
	"github.com/kwa0x2/Settle-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRoomRoute(db *mongo.Database, group *gin.RouterGroup, env *bootstrap.Env) {
	rr := repository.NewRoomRepository(db, domain.CollectionRoom)
	mr := repository.NewMessageRepository(db, domain.CollectionMessage)

	rd := &delivery.RoomDelivery{
		RoomUsecase: usecase.NewRoomUsecase(rr, mr),
	}

	group.GET("room", middleware.AuthMiddleware(env.AccessTokenSecret), rd.GetRooms)
}
