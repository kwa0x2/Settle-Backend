package bootstrap

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	SteamApiKey            string `mapstructure:"STEAM_API_KEY"`
	SteamRedirectUrl       string `mapstructure:"STEAM_REDIRECT_URL"`
	RedirectLoginUrl       string `mapstructure:"REDIRECT_LOGIN_URL"`
	MongoUri               string `mapstructure:"MONGO_URI"`
	MongoDBName            string `mapstructure:"MONGO_DB_NAME"`
	AWSRegion              string `mapstructure:"AWS_REGION"`
	S3BucketName           string `mapstructure:"S3_BUCKET_NAME"`
	AWSAccessKeyID         string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey     string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile("../.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
