package bootstrap

import (
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Applicaton struct {
	Env           *Env
	MongoDatabase *mongo.Database
	SocketServer  *socket.Server
}

func App() Applicaton {
	app := &Applicaton{}
	app.Env = NewEnv()
	app.MongoDatabase = ConnectMongoDB(app.Env)
	app.SocketServer = socket.NewServer(nil, nil)
	return *app
}
