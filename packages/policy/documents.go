package policy

import (
	"project/packages/models"
)

func CreateCanaryPolicyDocument(service *models.ServiceParameters) models.PolicyDocument {
	return models.PolicyDocument{
		Version: "2012-10-17",
		Statement: []models.StatementPolicyDocumentEntry{
			{
				Effect: "Allow",
				Action: []string{
					"s3:PutObject",
					"s3:GetObject",
				},
				Resource: "arn:aws:s3:::" + service.Bucket + "/" + service.Region + "/" + service.Canary.Name + "/*",
			},
			{
				Effect: "Allow",
				Action: []string{
					"s3:GetBucketLocation",
				},
				Resource: "arn:aws:s3:::" + service.Bucket,
			},
			{
				Effect: "Allow",
				Action: []string{
					"logs:CreateLogStream",
					"logs:PutLogEvents",
					"logs:CreateLogGroup",
				},
				Resource: "arn:aws:logs:" + service.Region + ":069377782092:log-group:/aws/lambda/lambda-" + service.Canary.Name + "-*",
			},
			{
				Effect: "Allow",
				Action: []string{
					"s3:ListAllMyBuckets",
					"xray:PutTraceSegments",
				},
				Resource: "*",
			},
			{
				Effect: "Allow",
				Action: []string{
					"lambda:InvokeAsync",
					"lambda:InvokeFunction",
				},
				Resource: "*",
			},
			{
				Effect: "Allow",
				Action: []string{
					"cloudwatch:PutMetricData",
				},
				Resource: "*",
				Condition: &models.ConditionPolicyDocumentEntry{
					StringEquals: models.StringEqualsPolicyDocumentEntry{
						CloudwatchNamespace: "CloudWatchSynthetics",
					},
				},
			},
		},
	}
}

func CreateLambdaBasicExecutionPolicyDocument(service *models.ServiceParameters) models.PolicyDocument {
	return models.PolicyDocument{
		Version: "2012-10-17",
		Statement: []models.StatementPolicyDocumentEntry{
			{
				Effect: "Allow",
				Action: []string{
					"logs:CreateLogStream",
					"logs:PutLogEvents",
					"logs:CreateLogGroup",
				},
				Resource: "arn:aws:logs:" + service.Region + ":069377782092:log-group:/aws/lambda/lambda-" + service.Canary.Name + "-*",
			},
		},
	}
}
