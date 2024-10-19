package bootstrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

func InitS3(env *Env) *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(env.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AWSAccessKeyID, env.AWSSecretAccessKey, "")))
	if err != nil {
		log.Fatal("failed to load aws cfg")
	}

	return s3.NewFromConfig(cfg)
}
