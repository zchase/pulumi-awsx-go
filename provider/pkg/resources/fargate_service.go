package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const FargateServiceIdentifier = "awsx-go:ecs:FargateService"

type FargateServiceArgs struct {
	Cluster                         string                                      `pulumi:"cluster"`
	ContinueBeforeSteadyState       bool                                        `pulumi:"continueBeforeSteadyState"`
	DeploymentCircuitBreaker        ecs.ServiceDeploymentCircuitBreakerPtrInput `pulumi:"deploymentCircuitBreaker"`
	DeploymentController            ecs.ServiceDeploymentControllerPtrInput     `pulumi:"deploymentController"`
	DeploymentMaximumPercent        int                                         `pulumi:"deploymentMaximumPercent"`
	DeploymentMinimumHealthyPercent int                                         `pulumi:"deploymentMinimumHealthyPercent"`
	DesiredCount                    int                                         `pulumi:"desiredCount"`
	EnableEcsManagedTags            bool                                        `pulumi:"enableEcsManagedTags"`
	EnableExecuteCommand            bool                                        `pulumi:"enableExecuteCommand"`
	ForceNewDeployment              bool                                        `pulumi:"forceNewDeployment"`
	HealthCheckGracePeriodSeconds   int                                         `pulumi:"healthCheckGracePeriodSeconds"`
	IAMRole                         string                                      `pulumi:"iamRole"`
	LoadBalancers                   ecs.ServiceLoadBalancerArrayInput           `pulumi:"loadBalancers"`
	Name                            string                                      `pulumi:"name"`
	NetworkConfiguration            ecs.ServiceNetworkConfigurationPtrInput     `pulumi:"networkConfiguration"`
	PlacementConstraints            ecs.ServicePlacementConstraintArrayInput    `pulumi:"placementConstraints"`
	PlatformVersions                string                                      `pulumi:"platformVersions"`
	PropagateTags                   string                                      `pulumi:"propagateTags"`
	SchedulingStrategy              string                                      `pulumi:"schedulingStrategy"`
	ServiceRegistries               ecs.ServiceServiceRegistriesPtrInput        `pulumi:"serviceRegistries"`
	Tags                            map[string]string                           `pulumi:"tags"`
	TaskDefinition                  string                                      `pulumi:"taskDefinition"`
	TaskDefinitionArgs              *FargateTaskDefinitionArgs                  `pulumi:"taskDefinitionArgs"`
}

type FargateService struct {
	pulumi.ResourceState

	Service        *ecs.Service           `pulumi:"service"`
	TaskDefinition *FargateTaskDefinition `pulumi:"taskDefinition"`
}

func NewFargateService(ctx *pulumi.Context, name string, args *FargateServiceArgs, opts ...pulumi.ResourceOption) (*FargateService, error) {
	if args == nil {
		args = &FargateServiceArgs{}
	}

	component := &FargateService{}
	err := ctx.RegisterComponentResource(FargateServiceIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(component))

	if args.TaskDefinition != "" && args.TaskDefinitionArgs != nil {
		return nil, fmt.Errorf("Only one of `taskDefinition` or `taskDefinitionArgs` can be provided.")
	}

	var taskDefinition *FargateTaskDefinition
	taskDefinitionIdentifier := pulumi.String(args.TaskDefinition).ToStringPtrOutput()

	if args.TaskDefinitionArgs != nil {
		taskDefinition, err = NewFargateTaskDefinition(ctx, name, args.TaskDefinitionArgs, opts...)
		if err != nil {
			return nil, err
		}

		component.TaskDefinition = taskDefinition
		taskDefinitionIdentifier = taskDefinition.TaskDefinition.Arn.ToStringPtrOutput()
	}

	if args.DesiredCount == 0 {
		args.DesiredCount = 1
	}

	if args.NetworkConfiguration == nil {
		args.NetworkConfiguration, err = getDefaultNetworkConfiguration(ctx, name, component)
		if err != nil {
			return nil, err
		}
	}

	//args.LoadBalancers
	args.LoadBalancers = pulumi.All(args.LoadBalancers.ToServiceLoadBalancerArrayOutput(), taskDefinition.LoadBalancers).ApplyT(func(argsLBs, tdLBs []ecs.ServiceLoadBalancer) []ecs.ServiceLoadBalancer {
		if len(argsLBs) == 0 {
			return tdLBs
		}
		return argsLBs
	}).(ecs.ServiceLoadBalancerArrayOutput)

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
		LaunchType:                      pulumi.String("FARGATE"),
		LoadBalancers:                   args.LoadBalancers,
		Name:                            pulumi.StringPtr(args.Name),
		NetworkConfiguration:            args.NetworkConfiguration,
		PlacementConstraints:            args.PlacementConstraints,
		PlatformVersion:                 pulumi.StringPtr(args.PlatformVersions),
		PropagateTags:                   pulumi.StringPtr(args.PropagateTags),
		SchedulingStrategy:              pulumi.StringPtr(args.SchedulingStrategy),
		ServiceRegistries:               args.ServiceRegistries,
		Tags:                            pulumi.ToStringMap(args.Tags),
		TaskDefinition:                  taskDefinitionIdentifier,
		WaitForSteadyState:              pulumi.BoolPtr(args.ContinueBeforeSteadyState),
	}, opts...)
	if err != nil {
		return nil, err
	}

	component.Service = service

	return component, nil
}

func getDefaultNetworkConfiguration(ctx *pulumi.Context, name string, parent pulumi.Resource) (*ecs.ServiceNetworkConfigurationArgs, error) {
	defaultVpc, err := getDefaultVPC(ctx, pulumi.Parent(parent))
	if err != nil {
		return nil, err
	}

	sgName := fmt.Sprintf("%s-sg", name)
	sg, err := ec2.NewSecurityGroup(ctx, sgName, &ec2.SecurityGroupArgs{
		VpcId: defaultVpc.VPCID,
		Ingress: ec2.SecurityGroupIngressArray{
			&ec2.SecurityGroupIngressArgs{
				FromPort:       pulumi.Int(0),
				ToPort:         pulumi.Int(0),
				Protocol:       pulumi.String("-1"),
				CidrBlocks:     pulumi.ToStringArray([]string{"0.0.0.0/0"}),
				Ipv6CidrBlocks: pulumi.ToStringArray([]string{"::/0"}),
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			&ec2.SecurityGroupEgressArgs{
				FromPort:       pulumi.Int(0),
				ToPort:         pulumi.Int(65535),
				Protocol:       pulumi.String("tcp"),
				CidrBlocks:     pulumi.ToStringArray([]string{"0.0.0.0/0"}),
				Ipv6CidrBlocks: pulumi.ToStringArray([]string{"::/0"}),
			},
		},
	}, pulumi.Parent(parent))
	if err != nil {
		return nil, err
	}

	return &ecs.ServiceNetworkConfigurationArgs{
		Subnets:        defaultVpc.PublicSubnetIDs,
		AssignPublicIp: pulumi.BoolPtr(true),
		SecurityGroups: pulumi.ToStringArrayOutput([]pulumi.StringOutput{sg.ID().ToStringOutput()}),
	}, nil
}
