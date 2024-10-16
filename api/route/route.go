package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

func Setup(env *bootstrap.Env, db *mongo.Database, router *gin.Engine, server *socket.Server) {
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 made by kwa -> https://github.com/kwa0x2")
	})

	publicRouter := router.Group("/api/v1")

	NewAuthRoute(env, db, publicRouter)
	NewSocketRoute(server, router, db)
}
