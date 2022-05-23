package resources

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/zchase/pulumi-awsx-go/pkg/utils"
)

const EC2TaskDefinitionIdentifier = "awsx-go:ecs:EC2TaskDefinition"

type EC2TaskDefinitionArgs struct {
	Containers            map[string]TaskDefinitionContainerDefinitionInputs `pulumi:"containers"`
	CPU                   string                                             `pulumi:"cpu"`
	EphemeralStorage      ecs.TaskDefinitionEphemeralStoragePtrInput         `pulumi:"ephemeralStorage"`
	ExecutionRole         DefaultRoleWithPolicyInputs                        `pulumi:"executionRole"`
	Family                string                                             `pulumi:"family"`
	InferenceAccelerators ecs.TaskDefinitionInferenceAcceleratorArrayInput   `pulumi:"inferenceAccelerators"`
	IPCMode               string                                             `pulumi:"ipcMode"`
	LogGroup              DefaultLogGroupInputs                              `pulumi:"logGroup"`
	Memory                string                                             `pulumi:"memory"`
	NetworkMode           string                                             `pulumi:"networkMode"`
	PIDMode               string                                             `pulumi:"pidMode"`
	PlacementConstraints  ecs.TaskDefinitionPlacementConstraintArrayInput    `pulumi:"placementConstraints"`
	ProxyConfiguration    ecs.TaskDefinitionProxyConfigurationPtrInput       `pulumi:"proxyConfiguration"`
	RuntimePlatform       ecs.TaskDefinitionRuntimePlatformPtrInput          `pulumi:"runtimePlatform"`
	SkipDestroy           bool                                               `pulumi:"skipDestroy"`
	Tags                  map[string]string                                  `pulumi:"tags"`
	TaskRole              DefaultRoleWithPolicyInputs                        `pulumi:"taskRole"`
	Volumes               ecs.TaskDefinitionVolumeArrayInput                 `pulumi:"volumes"`
}

type EC2TaskDefinition struct {
	pulumi.ResourceState

	ExecutionRole  *iam.Role                          `pulumi:"executionRole"`
	LoadBalancers  ecs.ServiceLoadBalancerArrayOutput `pulumi:"loadBalancers"`
	LogGroup       *cloudwatch.LogGroup               `pulumi:"logGroup"`
	TaskDefinition *ecs.TaskDefinition                `pulumi:"taskDefinition"`
	TaskRole       *iam.Role                          `pulumi:"taskRole"`
}

func NewEC2TaskDefinition(ctx *pulumi.Context, name string, args *EC2TaskDefinitionArgs, opts ...pulumi.ResourceOption) (*EC2TaskDefinition, error) {
	if args == nil {
		args = &EC2TaskDefinitionArgs{}
	}

	component := &EC2TaskDefinition{}
	err := ctx.RegisterComponentResource(EC2TaskDefinitionIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(component))

	dLogGroup, err := defaultLogGroup(ctx, name, &args.LogGroup, opts...)
	if err != nil {
		return nil, err
	}

	component.LogGroup = dLogGroup.LogGroup

	taskRoleName := fmt.Sprintf("%s-task", name)

	if len(args.TaskRole.Args.PolicyARNs) == 0 {
		args.TaskRole.Args.PolicyARNs = defaultExecutionRolePolicyARNs()
	}

	defaultPolicyDoc, err := defaultRoleAssumeRolePolicy(ctx)
	if err != nil {
		return nil, err
	}

	taskRole, err := defaultRoleWithPolicies(ctx, taskRoleName, args.TaskRole, defaultPolicyDoc.Json, opts...)
	if err != nil {
		return nil, err
	}
	component.TaskRole = taskRole.Role

	executionRoleName := fmt.Sprintf("%s-execution", name)

	if len(args.ExecutionRole.Args.PolicyARNs) == 0 {
		args.ExecutionRole.Args.PolicyARNs = defaultExecutionRolePolicyARNs()
	}

	executionRole, err := defaultRoleWithPolicies(ctx, executionRoleName, args.ExecutionRole, defaultPolicyDoc.Json, opts...)
	if err != nil {
		return nil, err
	}
	component.ExecutionRole = executionRole.Role

	containerDefinitions := computeContainerDefinitions(component, args.Containers, &dLogGroup.LogGroupID)

	computeLoadBalancers(ctx, args.Containers).ApplyT(func(x interface{}) interface{} {
		lbs := x.(ecs.ServiceLoadBalancerArrayOutput)
		component.LoadBalancers = lbs
		return nil
	})

	taskDefinitionArgs, err := buildTaskDefinitionArgs(ctx, name, args, containerDefinitions, taskRole.RoleARN, executionRole.RoleARN)
	if err != nil {
		return nil, err
	}

	taskDefinition, err := ecs.NewTaskDefinition(ctx, name, taskDefinitionArgs, opts...)
	if err != nil {
		return nil, err
	}

	component.TaskDefinition = taskDefinition

	return component, nil
}

func buildTaskDefinitionArgs(ctx *pulumi.Context, name string, args *EC2TaskDefinitionArgs, containerDefinitions []TaskDefinitionContainerDefinitionInputs, taskRoleARN, executionRoleARN pulumi.StringOutput) (*ecs.TaskDefinitionArgs, error) {
	containerDefJSON, err := json.Marshal(containerDefinitions)
	if err != nil {
		return nil, err
	}

	containerDefHash := utils.SHA1Hash(fmt.Sprintf("%s%s", ctx.Stack(), string(containerDefJSON)))
	defaultFamily := fmt.Sprintf("%s-%s", name, containerDefHash)
	if args.Family == "" {
		args.Family = defaultFamily
	}

	return &ecs.TaskDefinitionArgs{
		ContainerDefinitions:  pulumi.String(string(containerDefJSON)),
		Cpu:                   pulumi.StringPtr(args.CPU),
		EphemeralStorage:      args.EphemeralStorage,
		ExecutionRoleArn:      executionRoleARN,
		Family:                pulumi.String(args.Family),
		InferenceAccelerators: args.InferenceAccelerators,
		IpcMode:               pulumi.StringPtr(args.IPCMode),
		Memory:                pulumi.StringPtr(args.Memory),
		NetworkMode:           pulumi.StringPtr(args.NetworkMode),
		PidMode:               pulumi.StringPtr(args.PIDMode),
		PlacementConstraints:  args.PlacementConstraints,
		ProxyConfiguration:    args.ProxyConfiguration,
		RuntimePlatform:       args.RuntimePlatform,
		SkipDestroy:           pulumi.BoolPtr(args.SkipDestroy),
		Tags:                  pulumi.ToStringMap(args.Tags),
		TaskRoleArn:           taskRoleARN,
		Volumes:               args.Volumes,
	}, nil
}
