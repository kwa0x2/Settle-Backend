package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Setup(env *bootstrap.Env, db *mongo.Database, gin *gin.Engine) {
	publicRouter := gin.Group("/api/v1")

	NewAuthRoute(env, db, publicRouter)

}
