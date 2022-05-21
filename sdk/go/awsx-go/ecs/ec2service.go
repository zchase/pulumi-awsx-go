// Code generated by Pulumi SDK Generator DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package ecs

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Create an ECS Service resource for EC2 with the given unique name, arguments, and options.
// Creates Task definition if `taskDefinitionArgs` is specified.
type EC2Service struct {
	pulumi.ResourceState

	// Underlying ECS Service resource
	Service ecs.ServiceOutput `pulumi:"service"`
	// Underlying EC2 Task definition component resource if created from args
	TaskDefinition ecs.TaskDefinitionOutput `pulumi:"taskDefinition"`
}

// NewEC2Service registers a new resource with the given unique name, arguments, and options.
func NewEC2Service(ctx *pulumi.Context,
	name string, args *EC2ServiceArgs, opts ...pulumi.ResourceOption) (*EC2Service, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.NetworkConfiguration == nil {
		return nil, errors.New("invalid value for required argument 'NetworkConfiguration'")
	}
	var resource EC2Service
	err := ctx.RegisterRemoteComponentResource("awsx-go:ecs:EC2Service", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

type ec2serviceArgs struct {
	// ARN of an ECS cluster.
	Cluster *string `pulumi:"cluster"`
	// If `true`, this provider will not wait for the service to reach a steady state (like [`aws ecs wait services-stable`](https://docs.aws.amazon.com/cli/latest/reference/ecs/wait/services-stable.html)) before continuing. Default `false`.
	ContinueBeforeSteadyState *bool `pulumi:"continueBeforeSteadyState"`
	// Configuration block for deployment circuit breaker. See below.
	DeploymentCircuitBreaker *ecs.ServiceDeploymentCircuitBreaker `pulumi:"deploymentCircuitBreaker"`
	// Configuration block for deployment controller configuration. See below.
	DeploymentController *ecs.ServiceDeploymentController `pulumi:"deploymentController"`
	// Upper limit (as a percentage of the service's desiredCount) of the number of running tasks that can be running in a service during a deployment. Not valid when using the `DAEMON` scheduling strategy.
	DeploymentMaximumPercent *int `pulumi:"deploymentMaximumPercent"`
	// Lower limit (as a percentage of the service's desiredCount) of the number of running tasks that must remain running and healthy in a service during a deployment.
	DeploymentMinimumHealthyPercent *int `pulumi:"deploymentMinimumHealthyPercent"`
	// Number of instances of the task definition to place and keep running. Defaults to 0. Do not specify if using the `DAEMON` scheduling strategy.
	DesiredCount *int `pulumi:"desiredCount"`
	// Specifies whether to enable Amazon ECS managed tags for the tasks within the service.
	EnableEcsManagedTags *bool `pulumi:"enableEcsManagedTags"`
	// Specifies whether to enable Amazon ECS Exec for the tasks within the service.
	EnableExecuteCommand *bool `pulumi:"enableExecuteCommand"`
	// Enable to force a new task deployment of the service. This can be used to update tasks to use a newer Docker image with same image/tag combination (e.g., `myimage:latest`), roll Fargate tasks onto a newer platform version, or immediately deploy `ordered_placement_strategy` and `placement_constraints` updates.
	ForceNewDeployment *bool `pulumi:"forceNewDeployment"`
	// Seconds to ignore failing load balancer health checks on newly instantiated tasks to prevent premature shutdown, up to 2147483647. Only valid for services configured to use load balancers.
	HealthCheckGracePeriodSeconds *int `pulumi:"healthCheckGracePeriodSeconds"`
	// ARN of the IAM role that allows Amazon ECS to make calls to your load balancer on your behalf. This parameter is required if you are using a load balancer with your service, but only if your task definition does not use the `awsvpc` network mode. If using `awsvpc` network mode, do not specify this role. If your account has already created the Amazon ECS service-linked role, that role is used by default for your service unless you specify a role here.
	IamRole *string `pulumi:"iamRole"`
	// Configuration block for load balancers. See below.
	LoadBalancers []ecs.ServiceLoadBalancer `pulumi:"loadBalancers"`
	// Name of the service (up to 255 letters, numbers, hyphens, and underscores)
	Name *string `pulumi:"name"`
	// Network configuration for the service. This parameter is required for task definitions that use the `awsvpc` network mode to receive their own Elastic Network Interface, and it is not supported for other network modes. See below.
	NetworkConfiguration ecs.ServiceNetworkConfiguration `pulumi:"networkConfiguration"`
	// Service level strategy rules that are taken into consideration during task placement. List from top to bottom in order of precedence. Updates to this configuration will take effect next task deployment unless `force_new_deployment` is enabled. The maximum number of `ordered_placement_strategy` blocks is `5`. See below.
	OrderedPlacementStrategies []ecs.ServiceOrderedPlacementStrategy `pulumi:"orderedPlacementStrategies"`
	// Rules that are taken into consideration during task placement. Updates to this configuration will take effect next task deployment unless `force_new_deployment` is enabled. Maximum number of `placement_constraints` is `10`. See below.
	PlacementConstraints []ecs.ServicePlacementConstraint `pulumi:"placementConstraints"`
	// Platform version on which to run your service. Only applicable for `launch_type` set to `FARGATE`. Defaults to `LATEST`. More information about Fargate platform versions can be found in the [AWS ECS User Guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html).
	PlatformVersion *string `pulumi:"platformVersion"`
	// Specifies whether to propagate the tags from the task definition or the service to the tasks. The valid values are `SERVICE` and `TASK_DEFINITION`.
	PropagateTags *string `pulumi:"propagateTags"`
	// Scheduling strategy to use for the service. The valid values are `REPLICA` and `DAEMON`. Defaults to `REPLICA`. Note that [*Tasks using the Fargate launch type or the `CODE_DEPLOY` or `EXTERNAL` deployment controller types don't support the `DAEMON` scheduling strategy*](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html).
	SchedulingStrategy *string `pulumi:"schedulingStrategy"`
	// Service discovery registries for the service. The maximum number of `service_registries` blocks is `1`. See below.
	ServiceRegistries *ecs.ServiceServiceRegistries `pulumi:"serviceRegistries"`
	// Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
	Tags map[string]string `pulumi:"tags"`
	// Family and revision (`family:revision`) or full ARN of the task definition that you want to run in your service. Either [taskDefinition] or [taskDefinitionArgs] must be provided.
	TaskDefinition *string `pulumi:"taskDefinition"`
	// The args of task definition that you want to run in your service. Either [taskDefinition] or [taskDefinitionArgs] must be provided.
	TaskDefinitionArgs *EC2ServiceTaskDefinition `pulumi:"taskDefinitionArgs"`
}

// The set of arguments for constructing a EC2Service resource.
type EC2ServiceArgs struct {
	// ARN of an ECS cluster.
	Cluster pulumi.StringPtrInput
	// If `true`, this provider will not wait for the service to reach a steady state (like [`aws ecs wait services-stable`](https://docs.aws.amazon.com/cli/latest/reference/ecs/wait/services-stable.html)) before continuing. Default `false`.
	ContinueBeforeSteadyState pulumi.BoolPtrInput
	// Configuration block for deployment circuit breaker. See below.
	DeploymentCircuitBreaker ecs.ServiceDeploymentCircuitBreakerPtrInput
	// Configuration block for deployment controller configuration. See below.
	DeploymentController ecs.ServiceDeploymentControllerPtrInput
	// Upper limit (as a percentage of the service's desiredCount) of the number of running tasks that can be running in a service during a deployment. Not valid when using the `DAEMON` scheduling strategy.
	DeploymentMaximumPercent pulumi.IntPtrInput
	// Lower limit (as a percentage of the service's desiredCount) of the number of running tasks that must remain running and healthy in a service during a deployment.
	DeploymentMinimumHealthyPercent pulumi.IntPtrInput
	// Number of instances of the task definition to place and keep running. Defaults to 0. Do not specify if using the `DAEMON` scheduling strategy.
	DesiredCount pulumi.IntPtrInput
	// Specifies whether to enable Amazon ECS managed tags for the tasks within the service.
	EnableEcsManagedTags pulumi.BoolPtrInput
	// Specifies whether to enable Amazon ECS Exec for the tasks within the service.
	EnableExecuteCommand pulumi.BoolPtrInput
	// Enable to force a new task deployment of the service. This can be used to update tasks to use a newer Docker image with same image/tag combination (e.g., `myimage:latest`), roll Fargate tasks onto a newer platform version, or immediately deploy `ordered_placement_strategy` and `placement_constraints` updates.
	ForceNewDeployment pulumi.BoolPtrInput
	// Seconds to ignore failing load balancer health checks on newly instantiated tasks to prevent premature shutdown, up to 2147483647. Only valid for services configured to use load balancers.
	HealthCheckGracePeriodSeconds pulumi.IntPtrInput
	// ARN of the IAM role that allows Amazon ECS to make calls to your load balancer on your behalf. This parameter is required if you are using a load balancer with your service, but only if your task definition does not use the `awsvpc` network mode. If using `awsvpc` network mode, do not specify this role. If your account has already created the Amazon ECS service-linked role, that role is used by default for your service unless you specify a role here.
	IamRole pulumi.StringPtrInput
	// Configuration block for load balancers. See below.
	LoadBalancers ecs.ServiceLoadBalancerArrayInput
	// Name of the service (up to 255 letters, numbers, hyphens, and underscores)
	Name pulumi.StringPtrInput
	// Network configuration for the service. This parameter is required for task definitions that use the `awsvpc` network mode to receive their own Elastic Network Interface, and it is not supported for other network modes. See below.
	NetworkConfiguration ecs.ServiceNetworkConfigurationInput
	// Service level strategy rules that are taken into consideration during task placement. List from top to bottom in order of precedence. Updates to this configuration will take effect next task deployment unless `force_new_deployment` is enabled. The maximum number of `ordered_placement_strategy` blocks is `5`. See below.
	OrderedPlacementStrategies ecs.ServiceOrderedPlacementStrategyArrayInput
	// Rules that are taken into consideration during task placement. Updates to this configuration will take effect next task deployment unless `force_new_deployment` is enabled. Maximum number of `placement_constraints` is `10`. See below.
	PlacementConstraints ecs.ServicePlacementConstraintArrayInput
	// Platform version on which to run your service. Only applicable for `launch_type` set to `FARGATE`. Defaults to `LATEST`. More information about Fargate platform versions can be found in the [AWS ECS User Guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html).
	PlatformVersion pulumi.StringPtrInput
	// Specifies whether to propagate the tags from the task definition or the service to the tasks. The valid values are `SERVICE` and `TASK_DEFINITION`.
	PropagateTags pulumi.StringPtrInput
	// Scheduling strategy to use for the service. The valid values are `REPLICA` and `DAEMON`. Defaults to `REPLICA`. Note that [*Tasks using the Fargate launch type or the `CODE_DEPLOY` or `EXTERNAL` deployment controller types don't support the `DAEMON` scheduling strategy*](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html).
	SchedulingStrategy pulumi.StringPtrInput
	// Service discovery registries for the service. The maximum number of `service_registries` blocks is `1`. See below.
	ServiceRegistries ecs.ServiceServiceRegistriesPtrInput
	// Key-value map of resource tags. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
	Tags pulumi.StringMapInput
	// Family and revision (`family:revision`) or full ARN of the task definition that you want to run in your service. Either [taskDefinition] or [taskDefinitionArgs] must be provided.
	TaskDefinition pulumi.StringPtrInput
	// The args of task definition that you want to run in your service. Either [taskDefinition] or [taskDefinitionArgs] must be provided.
	TaskDefinitionArgs *EC2ServiceTaskDefinitionArgs
}

func (EC2ServiceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*ec2serviceArgs)(nil)).Elem()
}

type EC2ServiceInput interface {
	pulumi.Input

	ToEC2ServiceOutput() EC2ServiceOutput
	ToEC2ServiceOutputWithContext(ctx context.Context) EC2ServiceOutput
}

func (*EC2Service) ElementType() reflect.Type {
	return reflect.TypeOf((**EC2Service)(nil)).Elem()
}

func (i *EC2Service) ToEC2ServiceOutput() EC2ServiceOutput {
	return i.ToEC2ServiceOutputWithContext(context.Background())
}

func (i *EC2Service) ToEC2ServiceOutputWithContext(ctx context.Context) EC2ServiceOutput {
	return pulumi.ToOutputWithContext(ctx, i).(EC2ServiceOutput)
}

// EC2ServiceArrayInput is an input type that accepts EC2ServiceArray and EC2ServiceArrayOutput values.
// You can construct a concrete instance of `EC2ServiceArrayInput` via:
//
//          EC2ServiceArray{ EC2ServiceArgs{...} }
type EC2ServiceArrayInput interface {
	pulumi.Input

	ToEC2ServiceArrayOutput() EC2ServiceArrayOutput
	ToEC2ServiceArrayOutputWithContext(context.Context) EC2ServiceArrayOutput
}

type EC2ServiceArray []EC2ServiceInput

func (EC2ServiceArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*EC2Service)(nil)).Elem()
}

func (i EC2ServiceArray) ToEC2ServiceArrayOutput() EC2ServiceArrayOutput {
	return i.ToEC2ServiceArrayOutputWithContext(context.Background())
}

func (i EC2ServiceArray) ToEC2ServiceArrayOutputWithContext(ctx context.Context) EC2ServiceArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(EC2ServiceArrayOutput)
}

// EC2ServiceMapInput is an input type that accepts EC2ServiceMap and EC2ServiceMapOutput values.
// You can construct a concrete instance of `EC2ServiceMapInput` via:
//
//          EC2ServiceMap{ "key": EC2ServiceArgs{...} }
type EC2ServiceMapInput interface {
	pulumi.Input

	ToEC2ServiceMapOutput() EC2ServiceMapOutput
	ToEC2ServiceMapOutputWithContext(context.Context) EC2ServiceMapOutput
}

type EC2ServiceMap map[string]EC2ServiceInput

func (EC2ServiceMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*EC2Service)(nil)).Elem()
}

func (i EC2ServiceMap) ToEC2ServiceMapOutput() EC2ServiceMapOutput {
	return i.ToEC2ServiceMapOutputWithContext(context.Background())
}

func (i EC2ServiceMap) ToEC2ServiceMapOutputWithContext(ctx context.Context) EC2ServiceMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(EC2ServiceMapOutput)
}

type EC2ServiceOutput struct{ *pulumi.OutputState }

func (EC2ServiceOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**EC2Service)(nil)).Elem()
}

func (o EC2ServiceOutput) ToEC2ServiceOutput() EC2ServiceOutput {
	return o
}

func (o EC2ServiceOutput) ToEC2ServiceOutputWithContext(ctx context.Context) EC2ServiceOutput {
	return o
}

// Underlying ECS Service resource
func (o EC2ServiceOutput) Service() ecs.ServiceOutput {
	return o.ApplyT(func(v *EC2Service) ecs.ServiceOutput { return v.Service }).(ecs.ServiceOutput)
}

// Underlying EC2 Task definition component resource if created from args
func (o EC2ServiceOutput) TaskDefinition() ecs.TaskDefinitionOutput {
	return o.ApplyT(func(v *EC2Service) ecs.TaskDefinitionOutput { return v.TaskDefinition }).(ecs.TaskDefinitionOutput)
}

type EC2ServiceArrayOutput struct{ *pulumi.OutputState }

func (EC2ServiceArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*EC2Service)(nil)).Elem()
}

func (o EC2ServiceArrayOutput) ToEC2ServiceArrayOutput() EC2ServiceArrayOutput {
	return o
}

func (o EC2ServiceArrayOutput) ToEC2ServiceArrayOutputWithContext(ctx context.Context) EC2ServiceArrayOutput {
	return o
}

func (o EC2ServiceArrayOutput) Index(i pulumi.IntInput) EC2ServiceOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *EC2Service {
		return vs[0].([]*EC2Service)[vs[1].(int)]
	}).(EC2ServiceOutput)
}

type EC2ServiceMapOutput struct{ *pulumi.OutputState }

func (EC2ServiceMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*EC2Service)(nil)).Elem()
}

func (o EC2ServiceMapOutput) ToEC2ServiceMapOutput() EC2ServiceMapOutput {
	return o
}

func (o EC2ServiceMapOutput) ToEC2ServiceMapOutputWithContext(ctx context.Context) EC2ServiceMapOutput {
	return o
}

func (o EC2ServiceMapOutput) MapIndex(k pulumi.StringInput) EC2ServiceOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *EC2Service {
		return vs[0].(map[string]*EC2Service)[vs[1].(string)]
	}).(EC2ServiceOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*EC2ServiceInput)(nil)).Elem(), &EC2Service{})
	pulumi.RegisterInputType(reflect.TypeOf((*EC2ServiceArrayInput)(nil)).Elem(), EC2ServiceArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*EC2ServiceMapInput)(nil)).Elem(), EC2ServiceMap{})
	pulumi.RegisterOutputType(EC2ServiceOutput{})
	pulumi.RegisterOutputType(EC2ServiceArrayOutput{})
	pulumi.RegisterOutputType(EC2ServiceMapOutput{})
}
