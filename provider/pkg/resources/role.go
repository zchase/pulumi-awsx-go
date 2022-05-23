package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/zchase/pulumi-awsx-go/pkg/utils"
)

type RolePolicyAttachments struct {
	PolicyARNS []string
}

func defaultExecutionRolePolicyARNs() []string {
	return []string{
		"arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
	}
}

func defaultRoleAssumeRolePolicy(ctx *pulumi.Context) (*iam.GetPolicyDocumentResult, error) {
	args := &iam.GetPolicyDocumentArgs{
		Version: pulumi.StringRef("2012-10-17"),
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Actions: []string{"sts:AssumeRole"},
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type:        "Service",
						Identifiers: []string{"ecs-tasks.amazonaws.com"},
					},
				},
				Effect: pulumi.StringRef("Allow"),
				Sid:    pulumi.StringRef(""),
			},
		},
	}

	return iam.GetPolicyDocument(ctx, args)
}

type RoleWithPolicyInputs struct {
	Description         string                         `pulumi:"description"`
	ForceDetachPolicies bool                           `pulumi:"forceDetachPolicies"`
	InlinePolicies      iam.RoleInlinePolicyArrayInput `pulumi:"inlinePolicies"`
	ManagedPolicyArns   []string                       `pulumi:"managedPolicyArns"`
	MaxSessionDuration  int                            `pulumi:"maxSessionDuration"`
	Name                string                         `pulumi:"string"`
	NamePrefix          string                         `pulumi:"namePrefix"`
	Path                string                         `pulumi:"path"`
	PermissionsBoundary string                         `pulumi:"permissionsBoundary"`
	PolicyARNs          []string                       `pulumi:"policyArns"`
	Tags                map[string]string              `pulumi:"tags"`
}

type DefaultRoleWithPolicyInputs struct {
	Args    *RoleWithPolicyInputs `pulumi:"args"`
	RoleARN string                `pulumi:"roleArn"`
	Skip    bool                  `pulumi:"skip"`
}

type defaultRoleWithPoliciesResult struct {
	RoleARN  pulumi.StringOutput
	Role     *iam.Role
	Policies []*iam.RolePolicyAttachment
}

func defaultRoleWithPolicies(ctx *pulumi.Context, name string, inputs DefaultRoleWithPolicyInputs, assumeRolePolicy string, opts ...pulumi.ResourceOption) (*defaultRoleWithPoliciesResult, error) {
	if inputs.RoleARN != "" && inputs.Args != nil {
		return nil, fmt.Errorf("Can't define role args if specified an existing role ARN")
	}

	if inputs.Skip {
		return nil, nil
	}

	if inputs.RoleARN != "" {
		return &defaultRoleWithPoliciesResult{
			RoleARN: pulumi.String(inputs.RoleARN).ToStringOutput(),
		}, nil
	}

	args := inputs.Args

	if args.MaxSessionDuration == 0 {
		args.MaxSessionDuration = 3600
	}

	var roleName pulumi.StringPtrInput
	var roleNamePrefix pulumi.StringPtrInput

	if args.Name != "" {
		roleName = pulumi.StringPtr(args.Name)
	}

	if args.NamePrefix != "" {
		roleName = nil
		roleNamePrefix = pulumi.StringPtr(args.NamePrefix)
	}

	role, err := iam.NewRole(ctx, name, &iam.RoleArgs{
		AssumeRolePolicy:    pulumi.String(assumeRolePolicy),
		Description:         pulumi.StringPtr(args.Description),
		ForceDetachPolicies: pulumi.BoolPtr(args.ForceDetachPolicies),
		InlinePolicies:      args.InlinePolicies,
		ManagedPolicyArns:   pulumi.ToStringArray(args.ManagedPolicyArns),
		MaxSessionDuration:  pulumi.IntPtr(args.MaxSessionDuration),
		Name:                roleName,
		NamePrefix:          roleNamePrefix,
		Path:                pulumi.StringPtr(args.Path),
		PermissionsBoundary: pulumi.StringPtr(args.PermissionsBoundary),
		Tags:                pulumi.ToStringMap(args.Tags),
	}, opts...)
	if err != nil {
		return nil, err
	}

	var policies []*iam.RolePolicyAttachment
	for _, policyARN := range inputs.Args.PolicyARNs {
		policyAttachmentName := fmt.Sprintf("%s-%s", name, utils.SHA1Hash(policyARN))
		policy, err := iam.NewRolePolicyAttachment(ctx, policyAttachmentName, &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: pulumi.String(policyARN),
		}, opts...)
		if err != nil {
			return nil, err
		}

		policies = append(policies, policy)
	}

	return &defaultRoleWithPoliciesResult{
		Role:     role,
		Policies: policies,
		RoleARN:  role.Arn,
	}, nil
}
