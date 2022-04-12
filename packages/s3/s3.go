package s3

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"project/packages/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createDirectoryInBucket(path string, client *s3.S3, service *models.ServiceParameters) error {
	_, err := client.PutObject(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Bucket: &service.Bucket,
		Key:    aws.String(path),
	})
	if err != nil {
		fmt.Println("Failed to create a directory '"+path+"' into a S3 Bucket an object. Error:", err.Error())
		return err
	}
	return nil
}

func createTree(client *s3.S3, service *models.ServiceParameters) error {
	err := createDirectoryInBucket(service.Region+"/", client, service)
	if err != nil {
		return err
	}

	err = createDirectoryInBucket(service.Region+"/"+service.Canary.Name+"/source/", client, service)
	if err != nil {
		return err
	}

	err = createDirectoryInBucket(service.Region+"/"+service.Canary.Name+"/logs/", client, service)
	if err != nil {
		return err
	}

	err = createDirectoryInBucket(service.Region+"/"+service.Canary.Name+"/lambdas/", client, service)
	if err != nil {
		return err
	}

	return nil
}

func CreateBucket(session *session.Session, config *aws.Config, service *models.ServiceParameters) error {
	// Create a S3 Bucket from just a session.
	client := s3.New(session, config)

	// Create parameters of the S3 Bucket.
	parameters := s3.CreateBucketInput{
		Bucket: &service.Bucket,
	}

	// Create S3 Bucket.
	_, err := client.CreateBucket(&parameters)
	if err != nil {
		fmt.Println("Failed to create a S3 Bucket. Error:", err.Error())
		return err
	}

	// Create Tree in S3 Bucket.
	err = createTree(client, service)
	if err != nil {
		fmt.Println("Failed to create a tree in S3 Bucket. Error:", err.Error())
		return err
	}

	fmt.Println("S3Bucket", service.Bucket, "has successful created.")
	return nil
}

func SendFileToS3Bucket(session *session.Session, config *aws.Config, service *models.ServiceParameters, pathToDirectory string, pathToFile string) error {
	// Create a S3 Bucket from just a session.
	client := s3.New(session, config)

	// Get the fileName from Path
	fileName := filepath.Base(pathToFile)

	// Open the file from the file path
	upFile, err := os.Open(pathToFile)
	if err != nil {
		return err
	}
	defer upFile.Close()

	// Get the file info
	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = client.PutObject(&s3.PutObjectInput{
		ACL:                aws.String("public-read"),
		Bucket:             &service.Bucket,
		Key:                aws.String(pathToDirectory + "/" + fileName),
		Body:               bytes.NewReader(fileBuffer),
		ContentLength:      aws.Int64(fileSize),
		ContentType:        aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition: aws.String("attachment"),
	})

	if err != nil {
		return err
	}
	return nil
}
