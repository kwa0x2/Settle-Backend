package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/api/route"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.MongoDatabase
	ss := app.SocketServer
	s3 := app.S3Client
	gin := gin.Default()

	route.Setup(env, db, gin, ss, s3)

	gin.Run(env.ServerAddress)
}
