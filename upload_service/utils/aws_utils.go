package utils

import (
	"context"
	"log"
	"path/filepath"

	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config to Create any client
var finalcfg aws.Config

// S3 Client
var s3_client *s3.Client

// Bucket Name
var bucketName string = "takago-2"

// Intialize AWS Config and Bucket. This function runs when the package is imported
func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithDefaultRegion("ap-south-1"))
	if err != nil {
		panic(err)
	}
	finalcfg = cfg
	log.Println(cfg.Region)
	log.Println(cfg.Credentials.Retrieve(context.Background()))
	client := s3.NewFromConfig(finalcfg)
	s3_client = client
}

// UploadS3File Uploads a file to an S3 bucket.
func UploadS3File(fileName string, filePath string, prefix string) string {

	manager := manager.NewUploader(s3_client)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)

	}

	output, err := manager.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(filepath.ToSlash(filepath.Join(prefix, fileName))),
		Body:   file,
	})
	if err != nil {
		panic(err)
	}
	return output.Location

}
