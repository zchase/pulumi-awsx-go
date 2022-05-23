package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const RepositoryIdentifier = "awsx-go:ecr:Repository"

type RepositoryArgs struct {
	EncryptionConfigurations   ecr.RepositoryEncryptionConfigurationArrayInput  `pulumi:"encryptionConfigurations"`
	ImageScanningConfiguration ecr.RepositoryImageScanningConfigurationPtrInput `pulumi:"imageScanningConfiguration"`
	ImageTagMutability         string                                           `pulumi:"imageTagMutability"`
	LifecyclePolicy            lifecyclePolicy                                  `pulumi:"lifecyclePolicy"`
	Name                       string                                           `pulumi:"name"`
	Tags                       map[string]string                                `pulumi:"tags"`
}

type Repository struct {
	pulumi.ResourceState

	URL             pulumi.StringOutput  `pulumi:"url"`
	Repository      *ecr.Repository      `pulumi:"repository"`
	LifecyclePolicy *ecr.LifecyclePolicy `pulumi:"lifecyclePolicy"`
}

func NewRepository(ctx *pulumi.Context, name string, args *RepositoryArgs, opts ...pulumi.ResourceOption) (*Repository, error) {
	if args == nil {
		args = &RepositoryArgs{}
	}

	component := &Repository{}
	err := ctx.RegisterComponentResource(RepositoryIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(component))

	lowercaseName := strings.ToLower(name)

	repository, err := ecr.NewRepository(ctx, lowercaseName, &ecr.RepositoryArgs{
		EncryptionConfigurations:   args.EncryptionConfigurations,
		ImageScanningConfiguration: args.ImageScanningConfiguration,
		ImageTagMutability:         pulumi.String(args.ImageTagMutability),
		Name:                       pulumi.String(args.Name),
		Tags:                       pulumi.ToStringMap(args.Tags),
	}, opts...)
	if err != nil {
		return nil, err
	}

	if !args.LifecyclePolicy.Skip {
		policy, err := buildLifecyclePolicy(args.LifecyclePolicy)
		if err != nil {
			return nil, err
		}

		policyJSON, err := policy.MarshalJSON()
		if err != nil {
			return nil, err
		}

		lifecycle, err := ecr.NewLifecyclePolicy(ctx, lowercaseName, &ecr.LifecyclePolicyArgs{
			Repository: repository.Name,
			Policy:     pulumi.String(policyJSON),
		}, opts...)
		if err != nil {
			return nil, err
		}

		component.LifecyclePolicy = lifecycle
	}

	component.Repository = repository
	component.URL = repository.RepositoryUrl

	return component, nil
}

type lifecyclePolicyRule struct {
	Description           string
	MaximumAgeLimit       int
	MaximumNumberOfImages int
	TagPrefixList         []string
	TagStatus             string
}

type lifecyclePolicy struct {
	Skip  bool
	Rules []lifecyclePolicyRule
}

type lifecyclePolicyDocument struct {
	Rules []policyRule `json:"rules"`
}

func (l lifecyclePolicyDocument) MarshalJSON() (string, error) {
	result, err := json.Marshal(l)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

type policyRule struct {
	RulePriority pulumi.IntInput     `json:"rulePriority"`
	Description  pulumi.StringInput  `json:"description"`
	Selection    policyRuleSelection `json:"selection"`
	Action       policyRuleAction    `json:"action"`
}

type policyRuleSelection struct {
	TagStatus     pulumi.StringInput      `json:"tagStatus"`
	TagPrefixList pulumi.StringArrayInput `json:"tagPrefixList"`
	CountType     pulumi.StringInput      `json:"countType"`
	CountUnit     pulumi.StringPtrInput   `json:"countUnit"`
	CountNumber   pulumi.IntInput         `json:"countNumber"`
}

type policyRuleAction struct {
	Type pulumi.StringInput `pulumi:"type"`
}

func buildLifecyclePolicy(policy lifecyclePolicy) (lifecyclePolicyDocument, error) {
	if len(policy.Rules) == 0 {
		policy.Rules = append(policy.Rules, lifecyclePolicyRule{
			Description:           "remove untagged images",
			TagStatus:             "untagged",
			MaximumNumberOfImages: 1,
		})
	}

	return convertRules(policy.Rules)
}

func convertRules(rules []lifecyclePolicyRule) (lifecyclePolicyDocument, error) {
	result := lifecyclePolicyDocument{}

	var nonAnyRules []lifecyclePolicyRule
	var anyRules []lifecyclePolicyRule
	for _, rule := range rules {
		if rule.TagStatus == "any" {
			anyRules = append(anyRules, rule)
			continue
		}

		nonAnyRules = append(nonAnyRules, rule)
	}

	if len(anyRules) > 2 {
		return result, fmt.Errorf("At most one [selection: \"any\"] rule can be provided.")
	}

	orderedRules := nonAnyRules
	orderedRules = append(orderedRules, anyRules...)

	rulePriority := 1
	for _, rule := range orderedRules {
		convertedRule, err := convertRule(rule, rulePriority)
		if err != nil {
			return result, err
		}

		result.Rules = append(result.Rules, convertedRule)
		rulePriority++
	}

	return result, nil
}

func convertRule(rule lifecyclePolicyRule, rulePriority int) (policyRule, error) {
	if (rule.MaximumAgeLimit + rule.MaximumNumberOfImages) == 0 {
		return policyRule{}, fmt.Errorf("Either [maximumNumberOfImages] or [maximumAgeLimit] must be provided with a rule.")
	}

	ruleSelection := policyRuleSelection{}
	if rule.MaximumNumberOfImages > 0 {
		ruleSelection.CountType = pulumi.String("imageCountMoreThan")
		ruleSelection.CountNumber = pulumi.Int(rule.MaximumNumberOfImages)
	} else if rule.MaximumAgeLimit > 0 {
		ruleSelection.CountType = pulumi.String("sinceImagePushed")
		ruleSelection.CountNumber = pulumi.Int(rule.MaximumAgeLimit)
		ruleSelection.CountUnit = pulumi.String("days")
	}

	if (rule.TagStatus == "any") || (rule.TagStatus == "untagged") {
		ruleSelection.TagStatus = pulumi.String(rule.TagStatus)
	} else {
		if len(rule.TagPrefixList) == 0 {
			return policyRule{}, fmt.Errorf("tagPrefixList cannot be empty.")
		}

		ruleSelection.TagStatus = pulumi.String("tagged")
		ruleSelection.TagPrefixList = pulumi.ToStringArray(rule.TagPrefixList)
	}

	return policyRule{
		RulePriority: pulumi.Int(rulePriority),
		Description:  pulumi.String(rule.Description),
		Selection:    ruleSelection,
		Action: policyRuleAction{
			Type: pulumi.String("expire"),
		},
	}, nil
}
