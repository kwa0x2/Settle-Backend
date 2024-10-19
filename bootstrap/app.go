package bootstrap

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Applicaton struct {
	Env           *Env
	MongoDatabase *mongo.Database
	SocketServer  *socket.Server
	S3Client      *s3.Client
}

func App() Applicaton {
	app := &Applicaton{}
	app.Env = NewEnv()
	app.MongoDatabase = ConnectMongoDB(app.Env)
	app.SocketServer = socket.NewServer(nil, nil)
	app.S3Client = InitS3(app.Env)
	return *app
}
