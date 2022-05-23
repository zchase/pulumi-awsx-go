package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/zchase/pulumi-awsx-go/pkg/utils"
)

type LogGroupID struct {
	ARN            pulumi.StringOutput `pulumi:"arn"`
	LogGroupName   pulumi.StringOutput `pulumi:"logGroupName"`
	LogGroupRegion pulumi.StringOutput `pulumi:"logGroupRegion"`
}

type LogGroupIDOutput struct {
	pulumi.AnyOutput

	ARN            pulumi.StringOutput `pulumi:"arn"`
	LogGroupName   pulumi.StringOutput `pulumi:"logGroupName"`
	LogGroupRegion pulumi.StringOutput `pulumi:"logGroupRegion"`
}

type LogGroupInputs struct {
	KMSKeyID        string            `pulumi:"kmsKeyId"`
	NamePrefix      string            `pulumi:"namePrefix"`
	RetentionInDays int               `pulumi:"retentionInDays"`
	Tags            map[string]string `pulumi:"tags"`
}

type ExistingLogGroupInputs struct {
	ARN    string `pulumi:"arn"`
	Name   string `pulumi:"name"`
	Region string `pulumi:"region"`
}

type OptionalLogGroupInputs struct {
	Args     *LogGroupInputs         `pulumi:"args"`
	Enable   bool                    `pulumi:"enable"`
	Existing *ExistingLogGroupInputs `pulumi:"existing"`
}

type LogGroupArgs struct {
	Args     *LogGroupInputs
	Existing *ExistingLogGroupInputs
}

type LogGroupResult struct {
	LogGroup   *cloudwatch.LogGroup
	LogGroupID pulumi.AnyOutput
}

type DefaultLogGroupInputs struct {
	Args     *LogGroupInputs         `pulumi:"args"`
	Existing *ExistingLogGroupInputs `pulumi:"existing"`
	Skip     bool                    `pulumi:"skip"`
}

func defaultLogGroup(ctx *pulumi.Context, name string, inputs *DefaultLogGroupInputs, opts ...pulumi.ResourceOption) (*LogGroupResult, error) {
	if inputs.Skip {
		return nil, nil
	}
	return requiredLogGroup(ctx, name, &LogGroupArgs{
		Args:     inputs.Args,
		Existing: inputs.Existing,
	}, opts...)
}

func optionalLogGroup(ctx *pulumi.Context, name string, args *OptionalLogGroupInputs, opts ...pulumi.ResourceOption) (*LogGroupResult, error) {
	if args != nil && args.Enable == false {
		return &LogGroupResult{}, nil
	}

	if args == nil {
		args = &OptionalLogGroupInputs{}
	}

	return requiredLogGroup(ctx, name, &LogGroupArgs{
		Args:     args.Args,
		Existing: args.Existing,
	}, opts...)
}

func requiredLogGroup(ctx *pulumi.Context, name string, args *LogGroupArgs, opts ...pulumi.ResourceOption) (*LogGroupResult, error) {
	if args.Args != nil && args.Existing != nil {
		return nil, fmt.Errorf("Can't define log group args if specifying an existing log group name")
	}

	existing := args.Existing
	if existing != nil {
		if existing.ARN != "" {
			arnValue := pulumi.String(existing.ARN).ToStringOutput()
			logGroupID, err := makeLogGroupID(ctx, MakeLogGroupIDArgs{
				ARN: &arnValue,
			})
			if err != nil {
				return nil, err
			}

			return &LogGroupResult{
				LogGroupID: logGroupID,
			}, nil
		}

		if existing.Name != "" {
			// TODO: Figure out how pull the region programatically so it does not need to be supplied.
			region := existing.Region
			if region == "" {
				return nil, fmt.Errorf("Must specify a region")
			}

			nameValue := pulumi.String(existing.Name).ToStringOutput()
			regionValue := pulumi.String(region).ToStringOutput()

			logGroupID, err := makeLogGroupID(ctx, MakeLogGroupIDArgs{
				Name:   &nameValue,
				Region: &regionValue,
			})
			if err != nil {
				return nil, err
			}

			return &LogGroupResult{
				LogGroupID: logGroupID,
			}, nil
		}

		return nil, fmt.Errorf("One of an existing log group name or ARN must be specified")
	}

	logGroup, err := cloudwatch.NewLogGroup(ctx, name, &cloudwatch.LogGroupArgs{}, opts...)
	if err != nil {
		return nil, err
	}

	region, err := aws.GetRegion(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	regionValue := pulumi.String(region.Name).ToStringOutput()

	logGroupID, err := makeLogGroupID(ctx, MakeLogGroupIDArgs{
		Name:   &logGroup.Name,
		Region: &regionValue,
	})
	if err != nil {
		return nil, err
	}

	return &LogGroupResult{
		LogGroup:   logGroup,
		LogGroupID: logGroupID,
	}, nil
}

type MakeLogGroupIDArgs struct {
	ARN    *pulumi.StringOutput
	Name   *pulumi.StringOutput
	Region *pulumi.StringOutput
}

func (m MakeLogGroupIDArgs) Validate() error {
	if m.ARN != nil && m.Name != nil {
		return fmt.Errorf("Only one of an existing log group name or ARN can be specified")
	}

	return nil
}

func idFromARN(arn string) (LogGroupID, error) {
	parts, err := utils.ParseARN(arn)
	if err != nil {
		return LogGroupID{}, err
	}

	if parts.Service != "logs" || parts.ResourceType != "log-group" {
		return LogGroupID{}, fmt.Errorf("Invalid log group ARN")
	}

	return LogGroupID{
		ARN:            pulumi.String(arn).ToStringOutput(),
		LogGroupName:   pulumi.String(parts.ResourceID).ToStringOutput(),
		LogGroupRegion: pulumi.String(parts.Region).ToStringOutput(),
	}, nil
}

func buildLogGroupARN(region, name, accountID string) string {
	return fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", region, accountID, name)
}

type AnyOutput struct{ *pulumi.OutputState }

func makeLogGroupID(ctx *pulumi.Context, args MakeLogGroupIDArgs) (pulumi.AnyOutput, error) {
	err := args.Validate()
	if err != nil {
		return pulumi.AnyOutput{}, err
	}

	if args.ARN != nil {
		return args.ARN.ApplyT(idFromARN).(pulumi.AnyOutput), nil
	}

	callerIdentity, err := aws.GetCallerIdentity(ctx)
	if err != nil {
		return pulumi.AnyOutput{}, err
	}

	return pulumi.All(args.Region, args.Name, pulumi.String(callerIdentity.AccountId)).ApplyT(func(args []interface{}) (LogGroupID, error) {
		region := args[0].(string)
		name := args[1].(string)
		accountId := args[2].(string)

		arn := buildLogGroupARN(region, name, accountId)
		return idFromARN(arn)
	}).(pulumi.AnyOutput), nil
}
