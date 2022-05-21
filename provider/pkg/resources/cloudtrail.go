// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudtrail"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const trailIdentifier = "awsx-go:cloudtrail:Trail"

type TrailArgs struct {
	AdvancedEventSelector      cloudtrail.TrailAdvancedEventSelectorArray `pulumi:"advancedEventSelectors"`
	CloudWatchLogsGroup        *OptionalLogGroupInputs                    `pulumi:"cloudWatchLogsGroup"`
	EnableLogFileValidation    bool                                       `pulumi:"enableLogFileValidation"`
	EnableLogging              bool                                       `pulumi:"enableLogging"`
	EventSelectors             cloudtrail.TrailEventSelectorArray         `pulumi:"eventSelectors"`
	IncludeGlobalServiceEvents bool                                       `pulumi:"includeGlobalServiceEvents"`
	InsightSelectors           cloudtrail.TrailInsightSelectorArray       `pulumi:"insightSelectors"`
	IsMultiRegionTrail         bool                                       `pulumi:"isMultiRegionTrail"`
	isOrganizationTrail        bool                                       `pulumi:"isOrganizationTrail"`
	KMSKeyID                   string                                     `pulumi:"kmsKeyId"`
	Name                       string                                     `pulumi:"name"`
	S3Bucket                   RequiredBucketInputs                       `pulumi:"s3Bucket"`
	S3KeyPrefix                string                                     `pulumi:"s3KeyPrefix"`
	SNSTopicName               string                                     `pulumi:"snsTopicName"`
	Tags                       map[string]string                          `pulumi:"tags"`
}

type Trail struct {
	pulumi.ResourceState

	Bucket   *s3.Bucket           `pulumi:"Bucket"`
	LogGroup *cloudwatch.LogGroup `pulumi:"logGroup"`
	Trail    *cloudtrail.Trail    `pulumi:"trail"`
}

func NewTrail(ctx *pulumi.Context, name string, args *TrailArgs, opts ...pulumi.ResourceOption) (*Trail, error) {
	if args == nil {
		args = &TrailArgs{}
	}

	component := &Trail{}
	err := ctx.RegisterComponentResource(trailIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(component))

	bucket, err := requiredBucket(ctx, name, &args.S3Bucket, opts...)
	if err != nil {
		return nil, err
	}

	policy, err := createBucketCloudTrailPolicy(ctx, name, bucket.BucketID, bucket.Bucket, component)
	if err != nil {
		return nil, err
	}

	logGroup, err := optionalLogGroup(ctx, name, args.CloudWatchLogsGroup, opts...)
	if err != nil {
		return nil, err
	}

	trail, err := cloudtrail.NewTrail(ctx, name, &cloudtrail.TrailArgs{
		S3BucketName: bucket.Bucket.Bucket,
		CloudWatchLogsGroupArn: logGroup.LogGroupID.ApplyT(func(logGropuID LogGroupID) string {
			return fmt.Sprintf("%s:*", logGropuID.ARN)
		}).(pulumi.StringPtrInput),

		AdvancedEventSelectors:     args.AdvancedEventSelector,
		EnableLogFileValidation:    pulumi.BoolPtr(args.EnableLogFileValidation),
		EnableLogging:              pulumi.BoolPtr(args.EnableLogging),
		EventSelectors:             args.EventSelectors,
		IncludeGlobalServiceEvents: pulumi.BoolPtr(args.IncludeGlobalServiceEvents),
		InsightSelectors:           args.InsightSelectors,
		IsMultiRegionTrail:         pulumi.BoolPtr(args.IsMultiRegionTrail),
		IsOrganizationTrail:        pulumi.BoolPtr(args.isOrganizationTrail),
		KmsKeyId:                   pulumi.String(args.KMSKeyID),
		Name:                       pulumi.String(args.Name),
		S3KeyPrefix:                pulumi.String(args.S3KeyPrefix),
		SnsTopicName:               pulumi.String(args.SNSTopicName),
		Tags:                       pulumi.ToStringMap(args.Tags),
	}, pulumi.DependsOn([]pulumi.Resource{policy}), opts)
	if err != nil {
		return nil, err
	}

	component.Bucket = bucket.Bucket
	component.LogGroup = logGroup.LogGroup
	component.Trail = trail

	return component, nil
}

func createBucketCloudTrailPolicy(ctx *pulumi.Context, name string, bucketID BucketResultBucketID, bucket *s3.Bucket, parent pulumi.Resource) (*s3.BucketPolicy, error) {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}
	if bucket != nil {
		opts = append(opts, pulumi.DependsOn([]pulumi.Resource{bucket}))
	}

	return s3.NewBucketPolicy(ctx, name, &s3.BucketPolicyArgs{
		Bucket: bucketID.Name,
		Policy: bucketID.ARN.ApplyT(func(arn string) string {
			policy := defaultCloudTrailPolicy(arn)
			return policy.Json
		}),
	}, opts...)
}

func defaultCloudTrailPolicy(bucketARN string) iam.GetPolicyDocumentResult {
	return iam.GetPolicyDocumentResult{
		Version: pulumi.StringRef("2012-10-17"),
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Sid:       pulumi.StringRef("AWSCloudTrailAclCheck"),
				Effect:    pulumi.StringRef("Allow"),
				Actions:   []string{"s3:GetBucketAcl"},
				Resources: []string{bucketARN},
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type:        "Service",
						Identifiers: []string{"cloudtrail.amazonaws.com"},
					},
				},
			},
			{
				Sid:       pulumi.StringRef("AWSCloudTrailWrite"),
				Effect:    pulumi.StringRef("Allow"),
				Actions:   []string{"s3:PutObject"},
				Resources: []string{fmt.Sprintf("%s/*", bucketARN)},
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type:        "Service",
						Identifiers: []string{"cloudtrail.amazonaws.com"},
					},
				},
				Conditions: []iam.GetPolicyDocumentStatementCondition{
					{
						Test:     "StringEquals",
						Variable: "s3:x-amz-acl",
						Values:   []string{"bucket-owner-full-control"},
					},
				},
			},
		},
	}
}
