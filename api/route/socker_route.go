package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/api/socket/adapter"
	"github.com/kwa0x2/Settle-Backend/api/socket/gateway"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/repository"
	"github.com/kwa0x2/Settle-Backend/usecase"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewSocketRoute(server *socket.Server, router *gin.Engine, db *mongo.Database) {
	mr := repository.NewMessageRepository(db, domain.CollectionMessage)
	mu := usecase.NewMessageUsecase(mr)

	sg := gateway.NewSocketGateway(server, "/chat")
	sa := adapter.NewSocketAdapter(sg, mu)

	sa.HandleConnection()

	router.GET("socket.io/*any", gin.WrapH(server.ServeHandler(nil)))
	router.POST("socket.io/*any", gin.WrapH(server.ServeHandler(nil)))

}
