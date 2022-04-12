package main

import (
	"fmt"
	"project/packages/canary"
	"project/packages/lambda"
	"project/packages/models"
	"project/packages/policy"
	"project/packages/role"
	"project/packages/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	title = "srv-migration"
)

var (
	service = models.ServiceParameters{
		Region:  "us-west-2",
		Service: title,
		Bucket:  title + "-us-west-2",
		Policy:  title + "-policy",
		Role: models.Role{
			Name: title + "-role",
		},
		Canary: models.Canary{
			Name:           title + "-canary",
			Function:       "canary",
			Handler:        "handler",
			RuntimeVersion: "syn-nodejs-puppeteer-3.4",
			Expression:     "rate(5 minutes)",
		},
		Lambdas: []models.Lambda{
			{
				Name:     "get-healthcheck-lambda",
				Filename: "get-healthcheck.zip",
				Runtime:  "go1.x",
				Handler:  "get-healthcheck",
				Role: models.Role{
					Name: title + "-lambda-role",
				},
			},
			{
				Name:     "post-messages-lambda",
				Filename: "post-messages.zip",
				Runtime:  "go1.x",
				Handler:  "post-messages",
				Role: models.Role{
					Name: title + "-lambda-role",
				},
			},
		},
	}
)

func main() {
	config := aws.Config{
		Region: aws.String(service.Region),
	}

	session, err := session.NewSession(&config)
	if err != nil {
		fmt.Println("Failed to create aws-session.")
		return
	}

	err = s3.CreateBucket(session, &config, &service)
	if err != nil {
		fmt.Println("Failed to create a S3Bucket for a Canary. Error:", err.Error())
		return
	}

	err = lambda.CreateLambda(session, &config, &service)
	if err != nil {
		fmt.Println("Failed to prepare Lambda for a Canary. Error:", err.Error())
		return
	}

	canaryPolicyDocument := policy.CreateCanaryPolicyDocument(&service)
	err = role.CreateRole(session, &config, &service, &canaryPolicyDocument)
	if err != nil {
		fmt.Println("Failed to create a IAM Role for a Canary. Error:", err.Error())
		return
	}

	err = canary.CreateCanary(session, &config, &service)
	if err != nil {
		fmt.Println("Failed to create a Canary. Error:", err.Error())
		return
	}
}
