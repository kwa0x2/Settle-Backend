package utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"mime/multipart"
	"time"
)

func UploadFile(file multipart.File, fileHeader *multipart.FileHeader, env *bootstrap.Env, s3Client *s3.Client) (string, error) {
	// Generate a unique filename using the current Unix timestamp and the original filename.
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)

	// Prepare the S3 PutObject request parameters.
	params := &s3.PutObjectInput{
		Bucket: aws.String(env.S3BucketName), // Specify the S3 bucket name.
		Key:    aws.String(fileName),         // Set the object key (filename) in the bucket.
		Body:   file,                         // Set the file body to be uploaded.
	}

	// Upload the file to S3 using the PutObject method of the S3 client.
	_, err := s3Client.PutObject(context.TODO(), params)
	if err != nil {
		return "", err // Return an error if the upload fails.
	}

	// Construct the file URL to access the uploaded file.
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", env.S3BucketName, fileName)
	return fileURL, nil // Return the URL of the uploaded file.
}
