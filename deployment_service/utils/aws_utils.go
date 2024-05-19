package utils

import (
	"context"
	"log"
	"os"
	"path/filepath"

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
var uploadPath string = "outputs/"

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

func GetFilesInFolderOfAWS(key string) *s3.ListObjectsV2Output {

	files, err := s3_client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: &key,
	})
	log.Println(files.Prefix)
	if err != nil {
		panic(err)
	}
	return files
}

func DownloadFilesFromAWSToLocal(localdir string, key string) string {
	files := GetFilesInFolderOfAWS(key)
	log.Println("Local DIr:", localdir)
	for _, file := range files.Contents {
		log.Println(*file.Key)
		log.Println(filepath.Join(localdir, filepath.FromSlash((*file.Key)[len(key):])))
		DownloadS3File(*file.Key, filepath.Join(localdir, filepath.FromSlash((*file.Key)[len(key):])))
	}

	return localdir

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

func DownloadS3File(fileName string, filePath string) {
	manager := manager.NewDownloader(s3_client)
	// Check if the directory exists
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// If the directory doesn't exist, create it
		os.MkdirAll(dir, 0755)
	}
	file, err := os.Create(filePath)
	log.Default().Println(filePath)
	if err != nil {
		panic(err)
	}
	print("Downloading File:", fileName)
	_, err = manager.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(fileName),
	})
	if err != nil {
		panic(err)
	}
}
