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
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const trailIdentifier = "awsx-go:cloudtrail:Trail"

type RequiredBucketInputs struct {
	Args     BucketArgs           `pulumi:"args"`
	Existing ExistingBucketInputs `pulumi:"existing"`
}

type LogGroupInputs struct {
	KMSKeyID        string            `pulumi:"kmsKeyId"`
	NamePrefix      string            `pulumi:"namePrefix"`
	RetentionInDays int               `pulumi:"retentionInDays"`
	Tags            map[string]string `pulumi:"tags"`
}

type OptionalLogGroupInputs struct {
	Args     LogGroupInputs `pulumi:"args"`
	Enable   bool           `pulumi:"enable"`
	Existing bool           `pulumi:"existing"`
}

type TrailArgs struct {
	AdvancedEventSelector      []cloudtrail.AdvancedEventSelector `pulumi:"advancedEventSelectors"`
	CloudWatchLogsGroup        OptionalLogGroupInputs             `pulumi:"cloudWatchLogsGroup"`
	EnableLogFileValidation    bool                               `pulumi:"enableLogFileValidation"`
	EnableLogging              bool                               `pulumi:"enableLogging"`
	EventSelectors             []cloudtrail.EventSelector         `pulumi:"eventSelectors"`
	IncludeGlobalServiceEvents bool                               `pulumi:"includeGlobalServiceEvents"`
	InsightSelectors           []cloudtrail.InsightSelector       `pulumi:"insightSelectors"`
	IsMultiRegionTrail         bool                               `pulumi:"isMultiRegionTrail"`
	isOrganizationTrail        bool                               `pulumi:"isOrganizationTrail"`
	KMSKeyID                   string                             `pulumi:"kmsKeyId"`
	Name                       string                             `pulumi:"name"`
	S3Bucket                   RequiredBucketInputs               `pulumi:"s3Bucket"`
	S3KeyPrefix                string                             `pulumi:"s3KeyPrefix"`
	SNSTopicName               string                             `pulumi:"snsTopicName"`
	Tags                       map[string]string                  `pulumi:"tags"`
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

	return component, nil
}
