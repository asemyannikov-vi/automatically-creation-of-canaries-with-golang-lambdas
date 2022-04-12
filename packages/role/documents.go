package role

import "project/packages/models"

func CreateLambdaBasicRoleDocument(service *models.ServiceParameters) models.RoleDocument {
	return models.RoleDocument{
		Version: "2012-10-17",
		Statement: []models.StatementRoleDocumentEntry{
			{
				Effect: "Allow",
				Principal: &models.PrincipalRoleDocumentEntry{
					Service: "lambda.amazonaws.com",
				},
				Action: "sts:AssumeRole",
			},
		},
	}
}
