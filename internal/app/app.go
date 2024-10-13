package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/config"
	"github.com/kwa0x2/Settle-Backend/internal/di"
	"github.com/kwa0x2/Settle-Backend/internal/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	MongoDatabase *mongo.Database
	Router        *gin.Engine
}

func NewApp() *App {
	config.LoadEnv()
	mongoDatabase := config.ConnectMongoDB()
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                            // Allow requests from this origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},                                   // Exposed headers to the client
		AllowCredentials: true,                                                         // Allow credentials in requests
	}))

	return &App{
		MongoDatabase: mongoDatabase,
		Router:        router,
	}
}

func (a *App) SetupRoutes() {
	container := di.NewContainer(a.MongoDatabase)

	routes.AuthRoute(a.Router, container.AuthController)
}

func (a *App) Run() error {
	return a.Router.Run(":9090")
}
