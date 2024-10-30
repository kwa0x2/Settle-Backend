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

func NewMessageRoute(db *mongo.Database, group *gin.RouterGroup, env *bootstrap.Env) {
	mr := repository.NewMessageRepository(db, domain.CollectionMessage)
	rr := repository.NewRoomRepository(db, domain.CollectionRoom)
	ur := repository.NewUserRepository(db, domain.CollectionUser)

	md := &delivery.MessageDelivery{
		MessageUsecase: usecase.NewMessageUsecase(mr, rr, ur),
	}

	group.POST("message/history", middleware.AuthMiddleware(env.AccessTokenSecret), md.GetMessageHistory)
}
