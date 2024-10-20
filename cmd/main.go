package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/api/middleware"
	"github.com/kwa0x2/Settle-Backend/api/route"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	metrics "github.com/kwa0x2/Settle-Backend/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	metrics.RegisterMetrics()
}

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.MongoDatabase
	ss := app.SocketServer
	s3 := app.S3Client
	router := gin.Default()

	router.Use(middleware.MetricsMiddleware())

	route.Setup(env, db, router, ss, s3)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run(env.ServerAddress)
}
