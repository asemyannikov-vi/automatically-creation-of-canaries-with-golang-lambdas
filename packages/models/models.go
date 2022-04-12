package models

// RoleDocument is our definition of our role to be uploaded to IAM.
type ServiceParameters struct {
	Region  string
	Service string
	Bucket  string
	Policy  string
	Role    Role
	Lambdas []Lambda
	Canary  Canary
}

type Role struct {
	Name string
	ARN  string
}

type Lambda struct {
	Name       string
	Runtime    string
	Expression string
	Handler    string
	Filename   string
	Role       Role
}

type Canary struct {
	Name           string
	RuntimeVersion string
	Expression     string
	Function       string
	Handler        string
}

// PolicyDocument is our definition of our policies to be uploaded to IAM.
type PolicyDocument struct {
	Version   string
	Statement []StatementPolicyDocumentEntry
}

// StatementEntry will dictate what this policy will allow or not allow.
type StatementPolicyDocumentEntry struct {
	Effect    string                        `json:"Effect"`
	Action    []string                      `json:"Action"`
	Resource  string                        `json:"Resource"`
	Condition *ConditionPolicyDocumentEntry `json:"Condition,omitempty"`
}

// ConditionEntry will dictate what condition for this policy.
type ConditionPolicyDocumentEntry struct {
	StringEquals StringEqualsPolicyDocumentEntry `json:"StringEquals"`
}

// StringEqualsEntry will dictate what namespaces describe condition of the policy.
type StringEqualsPolicyDocumentEntry struct {
	CloudwatchNamespace string `json:"cloudwatch:namespace"`
}

// RoleDocument is our definition of our role to be uploaded to IAM.
type RoleDocument struct {
	Version   string
	Statement []StatementRoleDocumentEntry
}

// StatementEntry will dictate what this role will allow or not allow.
type StatementRoleDocumentEntry struct {
	Effect    string                      `json:"Effect"`
	Principal *PrincipalRoleDocumentEntry `json:"Principal"`
	Action    string                      `json:"Action"`
}

// PrincipalEntry will dictate what principals for this role.
type PrincipalRoleDocumentEntry struct {
	Service string `json:"Service"`
}
