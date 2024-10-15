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

	gin := gin.Default()

	route.Setup(env, db, gin)

	gin.Run(env.ServerAddress)
}
