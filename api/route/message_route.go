package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/api/http/delivery"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/repository"
	"github.com/kwa0x2/Settle-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewMessageRoute(db *mongo.Database, group *gin.RouterGroup) {
	mr := repository.NewMessageRepository(db, domain.CollectionMessage)
	md := &delivery.MessageDelivery{
		MessageUsecase: usecase.NewMessageUsecase(mr),
	}

	group.POST("message/history", md.GetMessageHistory)
}