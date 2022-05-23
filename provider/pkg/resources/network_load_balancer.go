package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const NetworkLoadBalancerIdentifier = "awsx-go:lb:NetworkLoadBalancer"

type NetworkLoadBalancerArgs struct {
	AccessLogs                   lb.LoadBalancerAccessLogsPtrInput `pulumi:"accessLogs"`
	CustomerOwnedIpv4Pool        string                            `pulumi:"customerOwnedIpv4Pool"`
	DefaultTargetGroup           TargetGroupInputs                 `pulumi:"defaultTargetGroup"`
	DesyncMitigationMode         string                            `pulumi:"desyncMitigationMode"`
	DropInvalidHeaderFields      bool                              `pulumi:"dropInvalidHeaderFields"`
	EnableDeletionProtection     bool                              `pulumi:"enableDeletionProtection"`
	EnableCrossZoneLoadBalancing bool                              `pulumi:"enableCrossZoneLoadBalancing"`
	EnableWafFailOpen            bool                              `pulumi:"enableWafFailOpen"`
	IdleTimeout                  int                               `pulumi:"idleTimeout"`
	Internal                     bool                              `pulumi:"internal"`
	IPAddressType                string                            `pulumi:"ipAddressType"`
	Listener                     *ListenerInputs                   `pulumi:"listener"`
	Listeners                    []*ListenerInputs                 `pulumi:"listeners"`
	Name                         string                            `pulumi:"name"`
	NamePrefix                   string                            `pulumi:"namePrefix"`
	SubnetIDs                    []string                          `pulumi:"subnetIds"`
	SubnetMappings               []lb.LoadBalancerSubnetMapping    `pulumi:"subnetMappings"`
	Subnets                      []ec2.Subnet                      `pulumi:"subnets"`
	Tags                         map[string]string                 `pulumi:"tags"`
}

type NetworkLoadBalancer struct {
	pulumi.ResourceState

	DefaultTargetGroup *lb.TargetGroup     `pulumi:"defaultTargetGroup"`
	Listeners          []*lb.Listener      `pulumi:"listeners"`
	LoadBalancer       *lb.LoadBalancer    `pulumi:"loadBalancer"`
	VpcID              pulumi.StringOutput `pulumi:"vpcId"`
}

func NewNetworkLoadBalancer(ctx *pulumi.Context, name string, args *NetworkLoadBalancerArgs, opts ...pulumi.ResourceOption) (*NetworkLoadBalancer, error) {
	if args == nil {
		args = &NetworkLoadBalancerArgs{}
	}

	opts = append(opts, pulumi.Aliases([]pulumi.Alias{
		{Name: pulumi.String("awsx:x:elasticloadbalancingv2:NetworkLoadBalancer")},
	}))

	component := &NetworkLoadBalancer{}
	err := ctx.RegisterComponentResource(NetworkLoadBalancerIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(component))

	subnetArgsSizes := []int{len(args.SubnetIDs), len(args.Subnets), len(args.SubnetMappings)}
	subnetArgsSet := 0
	for _, argSize := range subnetArgsSizes {
		if argSize > 0 {
			subnetArgsSet++
		}
	}

	if subnetArgsSet > 1 {
		return nil, fmt.Errorf("Only one of [subnets], [subnetIds] or [subnetMappings] can be specified")
	}

	var subnetIDs pulumi.StringArrayOutput
	subnetMappings := lb.LoadBalancerSubnetMappingArrayOutput{}
	if len(args.Subnets) > 0 {
		component.VpcID = args.Subnets[0].VpcId

		var sIds []pulumi.StringOutput
		for _, subnet := range args.Subnets {
			sIds = append(sIds, subnet.ID().ToStringOutput())
		}
		subnetIDs = pulumi.ToStringArrayOutput(sIds)
	} else if len(args.SubnetIDs) > 0 {
		subnetIDs = pulumi.ToStringArray(args.SubnetIDs).ToStringArrayOutput()

		firstSubnet, err := ec2.LookupSubnet(ctx, &ec2.LookupSubnetArgs{Id: pulumi.StringRef(args.SubnetIDs[0])})
		if err != nil {
			return nil, err
		}

		component.VpcID = pulumi.String(firstSubnet.VpcId).ToStringOutput()
	} else if len(args.SubnetMappings) > 0 {
		subnetMappings = subnetMappings.ApplyT(func(sm []lb.LoadBalancerSubnetMapping) []lb.LoadBalancerSubnetMapping {
			return args.SubnetMappings
		}).(lb.LoadBalancerSubnetMappingArrayOutput)

		firstSubnet, err := ec2.LookupSubnet(ctx, &ec2.LookupSubnetArgs{Id: pulumi.StringRef(args.SubnetMappings[0].SubnetId)})
		if err != nil {
			return nil, err
		}
		component.VpcID = pulumi.String(firstSubnet.VpcId).ToStringOutput()
	} else {
		defaultVPC, err := getDefaultVPC(ctx, pulumi.Parent(component))
		if err != nil {
			return nil, err
		}
		component.VpcID = defaultVPC.VPCID
		subnetIDs = defaultVPC.PublicSubnetIDs
	}

	if args.Listener != nil && len(args.Listeners) > 0 {
		return nil, fmt.Errorf("Only one of [listener] and [listeners] can be specified")
	}

	loadBalancerType := "network"
	lbArgs := &lb.LoadBalancerArgs{
		AccessLogs:               args.AccessLogs,
		CustomerOwnedIpv4Pool:    pulumi.String(args.CustomerOwnedIpv4Pool),
		DesyncMitigationMode:     pulumi.StringPtr(args.DesyncMitigationMode),
		DropInvalidHeaderFields:  pulumi.BoolPtr(args.DropInvalidHeaderFields),
		EnableDeletionProtection: pulumi.BoolPtr(args.EnableDeletionProtection),
		EnableHttp2:              pulumi.BoolPtr(false),
		EnableWafFailOpen:        pulumi.BoolPtr(args.EnableWafFailOpen),
		IdleTimeout:              pulumi.IntPtr(0),
		Internal:                 pulumi.BoolPtr(args.Internal),
		IpAddressType:            pulumi.StringPtr(args.IPAddressType),
		LoadBalancerType:         pulumi.String(loadBalancerType),
		Name:                     pulumi.StringPtr(args.Name),
		NamePrefix:               pulumi.StringPtr(args.NamePrefix),
		Subnets:                  subnetIDs,
		Tags:                     pulumi.ToStringMap(args.Tags),
	}

	loadBalancer, err := lb.NewLoadBalancer(ctx, name, lbArgs, opts...)
	if err != nil {
		return nil, err
	}
	component.LoadBalancer = loadBalancer

	tgArgs := &lb.TargetGroupArgs{
		VpcId:                          component.VpcID,
		TargetType:                     pulumi.StringPtr(args.DefaultTargetGroup.TargetType),
		Port:                           pulumi.IntPtr(80),
		Protocol:                       pulumi.StringPtr("TCP"),
		ConnectionTermination:          pulumi.BoolPtr(args.DefaultTargetGroup.ConnectionTermination),
		DeregistrationDelay:            pulumi.IntPtr(args.DefaultTargetGroup.DeRegistrationDelay),
		HealthCheck:                    args.DefaultTargetGroup.HealthCheck,
		LambdaMultiValueHeadersEnabled: pulumi.BoolPtr(args.DefaultTargetGroup.LambdaMultiValueHeadersEnabled),
		LoadBalancingAlgorithmType:     pulumi.String(args.DefaultTargetGroup.LoadBalancingAlgorithmType),
		Name:                           pulumi.String(args.DefaultTargetGroup.Name),
		NamePrefix:                     pulumi.String(args.DefaultTargetGroup.NamePrefix),
		PreserveClientIp:               pulumi.StringPtr(args.DefaultTargetGroup.PreserveClientIp),
		ProtocolVersion:                pulumi.StringPtr(args.DefaultTargetGroup.ProtocolVersion),
		SlowStart:                      pulumi.IntPtr(args.DefaultTargetGroup.SlowStart),
		Stickiness:                     args.DefaultTargetGroup.Stickiness.ToTargetGroupStickinessOutput(),
		Tags:                           pulumi.ToStringMap(args.DefaultTargetGroup.Tags),
	}

	if args.DefaultTargetGroup.VpcID != "" {
		tgArgs.VpcId = pulumi.String(args.DefaultTargetGroup.VpcID)
	}

	if args.DefaultTargetGroup.Port > 0 {
		tgArgs.Port = pulumi.IntPtr(args.DefaultTargetGroup.Port)
	}

	if args.DefaultTargetGroup.Protocol != "" {
		tgArgs.Protocol = pulumi.String(args.DefaultTargetGroup.Protocol)
	}

	targetGroup, err := lb.NewTargetGroup(ctx, name, tgArgs, opts...)
	if err != nil {
		return nil, err
	}
	component.DefaultTargetGroup = targetGroup

	listenersToCreate := args.Listeners
	if args.Listener != nil {
		listenersToCreate = append(listenersToCreate, args.Listener)
	}

	for i, listener := range listenersToCreate {
		port := listener.Port
		if port == 0 {
			port = 80
		}

		protocol := listener.Protocol
		if protocol == "" {
			protocol = "TCP"
		}

		listenerResource, err := lb.NewListener(ctx, fmt.Sprintf("%s-%v", name, i), &lb.ListenerArgs{
			LoadBalancerArn: component.LoadBalancer.Arn,
			DefaultActions: lb.ListenerDefaultActionArray{
				&lb.ListenerDefaultActionArgs{
					Type:           pulumi.String("forward"),
					TargetGroupArn: component.DefaultTargetGroup.Arn,
				},
			},
			Port:           pulumi.IntPtr(port),
			Protocol:       pulumi.String(protocol),
			AlpnPolicy:     pulumi.String(listener.ALPNPolicy),
			CertificateArn: pulumi.StringPtr(listener.CertificateARN),
			SslPolicy:      pulumi.StringPtr(listener.SSLPolicy),
			Tags:           pulumi.ToStringMap(listener.Tags),
		}, opts...)
		if err != nil {
			return nil, err
		}

		component.Listeners = append(component.Listeners, listenerResource)
	}

	if len(component.Listeners) == 0 {
		listenerResource, err := lb.NewListener(ctx, fmt.Sprintf("%s-%v", name, 0), &lb.ListenerArgs{
			LoadBalancerArn: component.LoadBalancer.Arn,
			DefaultActions: lb.ListenerDefaultActionArray{
				&lb.ListenerDefaultActionArgs{
					Type:           pulumi.String("forward"),
					TargetGroupArn: component.DefaultTargetGroup.Arn,
				},
			},
			Port:     pulumi.IntPtr(80),
			Protocol: pulumi.String("TCP"),
		}, opts...)
		if err != nil {
			return nil, err
		}

		component.Listeners = append(component.Listeners, listenerResource)
	}

	return component, nil
}
