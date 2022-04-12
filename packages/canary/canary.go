package canary

import (
	"fmt"
	"project/packages/models"
	"project/packages/s3"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/synthetics"
)

func prepareCanaryCode(session *session.Session, config *aws.Config, service *models.ServiceParameters) error {
	err := s3.SendFileToS3Bucket(session, config, service, service.Region+"/"+service.Canary.Name+"/source/", "./"+service.Canary.Function+".zip")
	if err != nil {
		fmt.Println("Failed to send a Canary Code to S3 Bucket. Error:", err.Error())
		return err
	}
	fmt.Println("Canary Code", service.Canary.Name, "has successful prepared.")
	return err
}

func CreateCanary(session *session.Session, config *aws.Config, service *models.ServiceParameters) error {
	err := prepareCanaryCode(session, config, service)
	if err != nil {
		fmt.Println("Failed to prepare a Canary Code. Error:", err.Error())
		return err
	}
	time.Sleep(15 * time.Second)

	// Create a Synthetics Client with additional configuration.
	client := synthetics.New(session, config)

	// Create parameters of the Canary.
	fmt.Println("ArtifactS3Location:", "s3://"+service.Bucket+"/"+service.Region+"/"+service.Canary.Name+"/logs/")
	fmt.Println("S3Key:", "s3://"+service.Bucket+"/"+service.Region+"/"+service.Canary.Name+"/source/"+service.Canary.Function+".zip")
	parameters := synthetics.CreateCanaryInput{
		ArtifactS3Location: aws.String("s3://" + service.Bucket + "/" + service.Region + "/" + service.Canary.Name + "/logs/"), // aws.String("s3://migration-us-west-2/us-west-2/migration-canary/logs/"),
		Code: &synthetics.CanaryCodeInput{
			Handler:  aws.String(service.Canary.Function + ".handler"),                                                       // aws.String("canary.handler"),
			S3Bucket: &service.Bucket,                                                                                        // aws.String("migration-us-west-2"),
			S3Key:    aws.String(service.Region + "/" + service.Canary.Name + "/source/" + service.Canary.Function + ".zip"), // aws.String("us-west-2/migration-canary/source/canary.zip")
		},
		ExecutionRoleArn: aws.String(service.Role.ARN), // aws.String("arn:aws:iam::069377782092:role/migration-role")
		Name:             aws.String(service.Canary.Name),
		RuntimeVersion:   aws.String(service.Canary.RuntimeVersion),
		Schedule: &synthetics.CanaryScheduleInput{
			DurationInSeconds: aws.Int64(60),
			Expression:        aws.String(service.Canary.Expression),
		},
	}

	// Create a Canary.
	_, err = client.CreateCanary(&parameters)
	if err != nil {
		fmt.Println("Failed to create a Canary. Error:", err.Error())
		return err
	}

	fmt.Println("Canary", service.Canary.Name, "has successful created.")
	return nil
}
