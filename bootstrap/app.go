package bootstrap

import "go.mongodb.org/mongo-driver/v2/mongo"

type Applicaton struct {
	Env           *Env
	MongoDatabase *mongo.Database
}

func App() Applicaton {
	app := &Applicaton{}
	app.Env = NewEnv()
	app.MongoDatabase = ConnectMongoDB(app.Env)
	return *app
}
