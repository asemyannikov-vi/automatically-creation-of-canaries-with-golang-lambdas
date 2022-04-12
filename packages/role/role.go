package role

import (
	"encoding/json"
	"fmt"
	"project/packages/models"
	"project/packages/policy"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func attachPolicy(client *iam.IAM, session *session.Session, config *aws.Config, service *models.ServiceParameters, policyDocument *models.PolicyDocument) error {
	policy := policy.CreatePolicy(session, config, service, policyDocument)
	_, err := client.AttachRolePolicy(&iam.AttachRolePolicyInput{
		PolicyArn: policy.Policy.Arn,
		RoleName:  aws.String(service.Role.Name),
	})
	if err != nil {
		fmt.Println("Failed to attach a Policy to Role. Error:", err.Error())
		return err
	}
	return nil
}

func CreateRole(session *session.Session, config *aws.Config, service *models.ServiceParameters, policyDocument *models.PolicyDocument) error {
	// Create a IAM Client from just a session.
	client := iam.New(session, config)

	// Builds our role document for IAM.
	roleDocument := CreateLambdaBasicRoleDocument(service)

	bytes, err := json.Marshal(&roleDocument)
	if err != nil {
		fmt.Println("Error marshaling role. Error: ", err.Error())
		return err
	}

	// Create parameters of the IAM Role.
	parameters := iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(string(bytes)),
		RoleName:                 aws.String(service.Role.Name),
	}

	// Create IAM Role.
	role, err := client.CreateRole(&parameters)
	if err != nil {
		fmt.Println("Failed to create an IAM Role. Error:", err.Error())
		return err
	}
	fmt.Println("Role", service.Role.Name, "has successful created.")
	service.Role.ARN = *role.Role.Arn

	// Attach a Policy to IAM Role.
	err = attachPolicy(client, session, config, service, policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func attachLambdaPolicy(client *iam.IAM, session *session.Session, config *aws.Config, service *models.ServiceParameters) error {
	for _, item := range service.Lambdas {
		_, err := client.AttachRolePolicy(&iam.AttachRolePolicyInput{
			PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonRDSFullAccess"),
			RoleName:  aws.String(item.Role.Name),
		})
		if err != nil {
			fmt.Println("Failed to attach a Policy to Role. Error:", err.Error())
			return err
		}

		_, err = client.AttachRolePolicy(&iam.AttachRolePolicyInput{
			PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"),
			RoleName:  aws.String(item.Role.Name),
		})
		if err != nil {
			fmt.Println("Failed to attach a Policy to Role. Error:", err.Error())
			return err
		}
	}
	return nil
}

func CreateLambdaRole(session *session.Session, config *aws.Config, service *models.ServiceParameters) error {
	for _, item := range service.Lambdas {
		// Create a IAM Client from just a session.
		client := iam.New(session, config)

		// Builds our role document for IAM.
		roleDocument := CreateLambdaBasicRoleDocument(service)

		bytes, err := json.Marshal(&roleDocument)
		if err != nil {
			fmt.Println("Error marshaling role. Error: ", err.Error())
			return err
		}

		// Create parameters of the IAM Role.
		parameters := iam.CreateRoleInput{
			AssumeRolePolicyDocument: aws.String(string(bytes)),
			RoleName:                 aws.String(item.Role.Name),
		}

		// Create IAM Role.
		role, err := client.CreateRole(&parameters)
		if err != nil {
			fmt.Println("Failed to create an IAM Role. Error:", err.Error())
			return err
		}
		fmt.Println("Role", item.Role.Name, "has successful created.")
		item.Role.ARN = *role.Role.Arn
		fmt.Println("ARN Role: ", item.Role.ARN, "must be equals", *role.Role.Arn)

		// Attach a Policy to IAM Role.
		err = attachLambdaPolicy(client, session, config, service)
		if err != nil {
			return err
		}
	}
	return nil
}
