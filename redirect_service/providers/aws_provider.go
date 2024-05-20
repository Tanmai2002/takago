package providers

import (
	"context"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config to Create any client
var finalcfg aws.Config

// S3 Client
var s3_client *s3.Client

// Bucket Name
var bucketName string = "takago-2"

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

func GetFileFromS3(filepat string) (io.ReadCloser, *string, error) {
	finalpath := strings.Join(strings.Split(filepath.Join("dist", filepat), "\\"), "/")
	log.Println("Getting File", finalpath)
	resp, err := s3_client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(finalpath),
	})
	if err != nil {
		return nil, nil, err
	}
	return resp.Body, resp.ContentType, nil

	// s3_client.GetObject(context.TODO(), *s3.GetObjectInput{
	// 	Bucket: &bucketName,
	// 	Key:    &aws.String(filepath)
	// })

}
