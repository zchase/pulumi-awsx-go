package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const TargetGroupAttachmentIdentifier = "awsx-go:lb:TargetGroupAttachment"

type TargetGroupAttachmentArgs struct {
	Instance       *ec2.Instance    `pulumi:"instance"`
	InstanceID     string           `pulumi:"instanceId"`
	Lambda         *lambda.Function `pulumi:"lambda"`
	LambdaARN      string           `pulumi:"lambdaArn"`
	TargetGroup    *lb.TargetGroup  `pulumi:"targetGroup"`
	TargetGroupARN string           `pulumi:"targetGroupARN"`
}

type TargetGroupAttachment struct {
	pulumi.ResourceState

	LambdaPermission      *lambda.Permission        `pulumi:"lambdaPermission"`
	TargetGroupAttachment *lb.TargetGroupAttachment `pulumi:"targetGroupAttachment"`
}

func NewTargetGroupAttachment(ctx *pulumi.Context, name string, args *TargetGroupAttachmentArgs, opts ...pulumi.ResourceOption) (*TargetGroupAttachment, error) {
	if args == nil {
		args = &TargetGroupAttachmentArgs{}
	}

	component := &TargetGroupAttachment{}
	err := ctx.RegisterComponentResource(TargetGroupAttachmentIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(component))

	if args.TargetGroupARN != "" && args.TargetGroup != nil {
		return nil, fmt.Errorf("Exactly 1 of [targetGroup] or [targetGroupArn] must be provided")
	}

	validAttachmentArgs := 0
	attachmentArgsValidation := []bool{args.InstanceID != "", args.Instance != nil, args.LambdaARN != "", args.Lambda != nil}
	for _, existingArg := range attachmentArgsValidation {
		if existingArg {
			validAttachmentArgs++
		}
	}

	if validAttachmentArgs > 1 {
		return nil, fmt.Errorf("Exactly 1 of [instance], [instanceId], [lambda] or [lambdaArn] must be provided.")
	}

	var targetGroupARN pulumi.StringOutput
	var targetType pulumi.StringPtrOutput
	if args.TargetGroup != nil {
		targetGroupARN = args.TargetGroup.Arn
		targetType = args.TargetGroup.TargetType
	} else {
		if args.TargetGroupARN == "" {
			return nil, fmt.Errorf("Unreachable")
		}

		targetGroup, err := lb.LookupTargetGroup(ctx, &lb.LookupTargetGroupArgs{Arn: &args.TargetGroupARN})
		if err != nil {
			return nil, err
		}

		targetGroupARN = pulumi.String(args.TargetGroupARN).ToStringOutput()
		targetType = pulumi.String(targetGroup.TargetType).ToStringPtrOutput()
	}

	var targetID pulumi.StringOutput
	var availabilityZone pulumi.StringOutput
	if args.Instance != nil {
		targetID = pulumi.All(targetType, args.Instance.ID().ToStringOutput(), args.Instance.PrivateIp).ApplyT(func(t *string, instanceId, privateIP string) string {
			if *t == "instance" {
				return instanceId
			}

			return privateIP
		}).(pulumi.StringOutput)
		availabilityZone = args.Instance.AvailabilityZone
	} else if args.InstanceID != "" {
		instanceOutputs := ec2.LookupInstanceOutput(ctx, ec2.LookupInstanceOutputArgs{
			InstanceId: pulumi.StringPtr(args.InstanceID),
		})

		targetID = pulumi.All(targetType, instanceOutputs.Id(), instanceOutputs.PrivateIp()).ApplyT(func(t *string, instanceId, privateIP string) string {
			if *t == "instance" {
				return instanceId
			}

			return privateIP
		}).(pulumi.StringOutput)
		availabilityZone = instanceOutputs.AvailabilityZone()
	} else if args.Lambda != nil {
		targetID = args.Lambda.Arn
	} else if args.LambdaARN != "" {
		targetID = pulumi.String(args.LambdaARN).ToStringOutput()
	} else {
		return nil, fmt.Errorf("Unreachable condition")
	}

	if args.Lambda != nil || args.LambdaARN != "" {
		lambdaFunc := pulumi.Input(args.Lambda)
		if lambdaFunc == nil {
			lambdaFunc = pulumi.String(args.LambdaARN)
		}

		lambdaPermission, err := lambda.NewPermission(ctx, name, &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Principal: pulumi.String("elasticloadbalancing.amazonaws.com"),
			SourceArn: targetGroupARN,
			Function:  lambdaFunc,
		}, opts...)
		if err != nil {
			return nil, err
		}

		component.LambdaPermission = lambdaPermission
	}

	targetGroupAttachmentOpts := opts
	if component.LambdaPermission != nil {
		targetGroupAttachmentOpts = append(targetGroupAttachmentOpts, pulumi.DependsOn([]pulumi.Resource{component.LambdaPermission}))
	}

	targetGroupAttachment, err := lb.NewTargetGroupAttachment(ctx, name, &lb.TargetGroupAttachmentArgs{
		TargetGroupArn:   targetGroupARN,
		TargetId:         targetID,
		AvailabilityZone: availabilityZone,
	}, targetGroupAttachmentOpts...)
	if err != nil {
		return nil, err
	}
	component.TargetGroupAttachment = targetGroupAttachment

	return component, nil
}
