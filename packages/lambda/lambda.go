package lambda

import (
	"fmt"
	"project/packages/models"
	"project/packages/role"
	"project/packages/s3"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func prepareLambdaCode(session *session.Session, config *aws.Config, service *models.ServiceParameters) error {
	for _, item := range service.Lambdas {
		err := role.CreateLambdaRole(session, config, service)
		if err != nil {
			fmt.Println("Failed to create a IAM Role for a Lambda. Error:", err.Error())
			return err
		}
		err = s3.SendFileToS3Bucket(session, config, service, service.Region+"/"+service.Canary.Name+"/lambdas/", "./"+item.Filename)
		if err != nil {
			fmt.Println("Failed to send a Lambda Code to the S3 Bucket. Error:", err.Error())
			return err
		}
		fmt.Println("Lambda Code", item.Name, "has successful prepared.")
	}
	return nil
}

func CreateLambda(session *session.Session, config *aws.Config, service *models.ServiceParameters) error {
	for _, item := range service.Lambdas {
		err := prepareLambdaCode(session, config, service)
		if err != nil {
			fmt.Println("Failed to prepare the code of Lambda. Error:", err.Error())
			return err
		}
		time.Sleep(15 * time.Second)

		// Create a IAM Client from just a session.
		client := lambda.New(session, config)

		createCode := &lambda.FunctionCode{
			S3Bucket: &service.Bucket,
			S3Key:    aws.String(service.Region + "/" + service.Canary.Name + "/lambdas/" + item.Filename),
		}

		createArgs := &lambda.CreateFunctionInput{
			Code:         createCode,
			FunctionName: &item.Name,
			Handler:      &item.Handler,
			Role:         &item.Role.ARN,
			Runtime:      &item.Runtime,
		}

		_, err = client.CreateFunction(createArgs)
		if err != nil {
			fmt.Println("Failed to create a Lambda. Error:", err.Error())
			return err
		}
		fmt.Println("Lambda", item.Name, "has successful created.")
	}
	return nil
}
