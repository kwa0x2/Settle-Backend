package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/api/socket/adapter"
	"github.com/kwa0x2/Settle-Backend/api/socket/gateway"
	"github.com/zishang520/socket.io/socket"
)

func NewSocketRoute(server *socket.Server, router *gin.Engine) {
	sg := gateway.NewSocketGateway(server, "/chat")
	sa := adapter.NewSocketAdapter(sg)

	sa.HandleConnection()

	router.GET("socket.io/*any", gin.WrapH(server.ServeHandler(nil)))
	router.POST("socket.io/*any", gin.WrapH(server.ServeHandler(nil)))

}
