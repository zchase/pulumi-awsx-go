package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const EC2ServiceIdentifier = "awsx-go:ecs:EC2Service"

type EC2ServiceArgs struct {
	Cluster                         string                                        `pulumi:"cluster"`
	ContinueBeforeSteadyState       bool                                          `pulumi:"continueBeforeSteadyState"`
	DeploymentCircuitBreaker        ecs.ServiceDeploymentCircuitBreakerPtrInput   `pulumi:"deploymentCircuitBreaker"`
	DeploymentController            ecs.ServiceDeploymentControllerPtrInput       `pulumi:"deploymentController"`
	DeploymentMaximumPercent        int                                           `pulumi:"deploymentMaximumPercent"`
	DeploymentMinimumHealthyPercent int                                           `pulumi:"deploymentMinimumHealthyPercent"`
	DesiredCount                    int                                           `pulumi:"desiredCount"`
	EnableEcsManagedTags            bool                                          `pulumi:"enableEcsManagedTags"`
	EnableExecuteCommand            bool                                          `pulumi:"enableExecuteCommand"`
	ForceNewDeployment              bool                                          `pulumi:"forceNewDeployment"`
	HealthCheckGracePeriodSeconds   int                                           `pulumi:"healthCheckGracePeriodSeconds"`
	IAMRole                         string                                        `pulumi:"iamRole"`
	LoadBalancers                   ecs.ServiceLoadBalancerArrayInput             `pulumi:"loadBalancers"`
	Name                            string                                        `pulumi:"name"`
	NetworkConfiguration            ecs.ServiceNetworkConfigurationPtrInput       `pulumi:"networkConfiguration"`
	OrderedPlacementStrategies      ecs.ServiceOrderedPlacementStrategyArrayInput `pulumi:"orderedPlacementStrategies"`
	PlacementConstraints            ecs.ServicePlacementConstraintArrayInput      `pulumi:"placementConstraints"`
	PlatformVersions                string                                        `pulumi:"platformVersions"`
	PropagateTags                   string                                        `pulumi:"propagateTags"`
	SchedulingStrategy              string                                        `pulumi:"schedulingStrategy"`
	ServiceRegistries               ecs.ServiceServiceRegistriesPtrInput          `pulumi:"serviceRegistries"`
	Tags                            map[string]string                             `pulumi:"tags"`
	TaskDefinition                  string                                        `pulumi:"taskDefinition"`
	TaskDefinitionArgs              *EC2TaskDefinitionArgs                        `pulumi:"taskDefinitionArgs"`
}

type EC2Service struct {
	pulumi.ResourceState

	Service        *ecs.Service       `pulumi:"service"`
	TaskDefinition *EC2TaskDefinition `pulumi:"taskDefinition"`
}

func NewEC2Service(ctx *pulumi.Context, name string, args *EC2ServiceArgs, opts ...pulumi.ResourceOption) (*EC2Service, error) {
	var err error
	if args == nil {
		args = &EC2ServiceArgs{}
	}

	component := &EC2Service{}
	err = ctx.RegisterComponentResource(EC2ServiceIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(component))

	if args.TaskDefinition != "" && args.TaskDefinitionArgs != nil {
		return nil, fmt.Errorf("Only one of `taskDefinition` or `taskDefinitionArgs` can be provided.")
	}

	var taskDefinition *EC2TaskDefinition
	taskDefinitionIdentifier := pulumi.String(args.TaskDefinition).ToStringPtrOutput()

	if args.TaskDefinitionArgs != nil {
		taskDefinition, err = NewEC2TaskDefinition(ctx, name, args.TaskDefinitionArgs, opts...)
		if err != nil {
			return nil, err
		}

		component.TaskDefinition = taskDefinition
		taskDefinitionIdentifier = taskDefinition.TaskDefinition.Arn.ToStringPtrOutput()
	}

	service, err := ecs.NewService(ctx, name, &ecs.ServiceArgs{
		Cluster:                         pulumi.StringPtr(args.Cluster),
		DeploymentCircuitBreaker:        args.DeploymentCircuitBreaker,
		DeploymentController:            args.DeploymentController,
		DeploymentMaximumPercent:        pulumi.IntPtr(args.DeploymentMaximumPercent),
		DeploymentMinimumHealthyPercent: pulumi.IntPtr(args.DeploymentMinimumHealthyPercent),
		DesiredCount:                    pulumi.IntPtr(args.DesiredCount),
		EnableEcsManagedTags:            pulumi.BoolPtr(args.EnableEcsManagedTags),
		EnableExecuteCommand:            pulumi.BoolPtr(args.EnableExecuteCommand),
		ForceNewDeployment:              pulumi.BoolPtr(args.ForceNewDeployment),
		HealthCheckGracePeriodSeconds:   pulumi.IntPtr(args.HealthCheckGracePeriodSeconds),
		IamRole:                         pulumi.StringPtr(args.IAMRole),
		LoadBalancers:                   args.LoadBalancers,
		Name:                            pulumi.StringPtr(args.Name),
		NetworkConfiguration:            args.NetworkConfiguration,
		OrderedPlacementStrategies:      args.OrderedPlacementStrategies,
		PlacementConstraints:            args.PlacementConstraints,
		PlatformVersion:                 pulumi.StringPtr(args.PlatformVersions),
		PropagateTags:                   pulumi.StringPtr(args.PropagateTags),
		SchedulingStrategy:              pulumi.StringPtr(args.SchedulingStrategy),
		ServiceRegistries:               args.ServiceRegistries,
		Tags:                            pulumi.ToStringMap(args.Tags),
		TaskDefinition:                  taskDefinitionIdentifier,
	}, opts...)
	if err != nil {
		return nil, err
	}

	component.Service = service

	return component, nil
}
