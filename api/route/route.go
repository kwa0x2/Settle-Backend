package route

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

func Setup(env *bootstrap.Env, db *mongo.Database, router *gin.Engine, server *socket.Server, s3 *s3.Client) {
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 made by kwa -> https://github.com/kwa0x2")
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4725"},                            // Allow requests from this origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},                                   // Exposed headers to the client
		AllowCredentials: true,                                                         // Allow credentials in requests
	}))

	publicRouter := router.Group("/api/v1")

	NewAuthRoute(env, db, publicRouter)
	NewSocketRoute(server, router, db)
	NewMessageRoute(db, publicRouter)
	NewAttachmentRoute(env, db, publicRouter, s3)
}
