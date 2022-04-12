package policy

import (
	"encoding/json"
	"fmt"
	"project/packages/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func CreatePolicy(session *session.Session, config *aws.Config, service *models.ServiceParameters, policyDocument *models.PolicyDocument) *iam.CreatePolicyOutput {
	// Create a IAM Client from just a session.
	client := iam.New(session, config)

	// Builds our policy document for IAM.

	bytes, err := json.Marshal(&policyDocument)
	if err != nil {
		fmt.Println("Error marshaling policy. Error: ", err.Error())
		return nil
	}

	// Create parameters of the IAM Policy.
	parameters := iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(bytes)),
		PolicyName:     aws.String(service.Policy),
	}

	// Create IAM Policy.
	policy, err := client.CreatePolicy(&parameters)
	if err != nil {
		fmt.Println("Failed to create an IAM Policy. Error:", err.Error())
		return nil
	}
	fmt.Println("Policy", service.Policy, "has successful created.")
	return policy
}
