package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Setup(env *bootstrap.Env, db *mongo.Database, router *gin.Engine, server *socket.Server) {
	publicRouter := router.Group("/api/v1")

	NewAuthRoute(env, db, publicRouter)
	NewSocketRoute(server, router)
}
