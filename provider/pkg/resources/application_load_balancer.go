package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const ApplicationLoadBalancerIdentifier = "awsx-go:lb:ApplicationLoadBalancer"

type SecurityGroupInputs struct {
	Description         string                             `pulumi:"description"`
	Egress              ec2.SecurityGroupEgressArrayInput  `pulumi:"egress"`
	Ingress             ec2.SecurityGroupIngressArrayInput `pulumi:"ingress"`
	Name                string                             `pulumi:"string"`
	NamePrefix          string                             `pulumi:"namePrefix"`
	RevokeRulesOnDelete bool                               `pulumi:"revokeRulesOnDelete"`
	Tags                map[string]string                  `pulumi:"tags"`
	VpcID               string                             `pulumi:"vpcId"`
}

type DefaultSecurityGroupInputs struct {
	Args            *SecurityGroupInputs `pulumi:"args"`
	SecurityGroupID string               `pulumi:"securityGroupId"`
	Skip            bool                 `pulumi:"skip"`
}

type TargetGroupInputs struct {
	ConnectionTermination          bool                              `pulumi:"connectionTermination"`
	DeRegistrationDelay            int                               `pulumi:"deregistrationDelay"`
	HealthCheck                    lb.TargetGroupHealthCheckPtrInput `pulumi:"healthCheck"`
	LambdaMultiValueHeadersEnabled bool                              `pulumi:"lambdaMultiValueHeadersEnabled"`
	LoadBalancingAlgorithmType     string                            `pulumi:"loadBalancingAlgorithmType"`
	Name                           string                            `pulumi:"name"`
	NamePrefix                     string                            `pulumi:"namePrefix"`
	Port                           int                               `pulumi:"port"`
	PreserveClientIp               string                            `pulumi:"preserveClientIp"`
	Protocol                       string                            `pulumi:"protocol"`
	ProtocolVersion                string                            `pulumi:"protocolVersion"`
	ProxyProtocolV2                bool                              `pulumi:"proxyProtocolV2"`
	SlowStart                      int                               `pulumi:"slowStart"`
	Stickiness                     lb.TargetGroupStickinessPtrInput  `pulumi:"stickiness"`
	Tags                           map[string]string                 `pulumi:"tags"`
	TargetType                     string                            `pulumi:"targetType"`
	VpcID                          string                            `pulumi:"vpcId"`
}

type ListenerInputs struct {
	ALPNPolicy     string                             `pulumi:"alpnPolicy"`
	CertificateARN string                             `pulumi:"certificateArn"`
	DefaultActions lb.ListenerDefaultActionArrayInput `pulumi:"defaultActions"`
	Port           int                                `pulumi:"port"`
	Protocol       string                             `pulumi:"protocol"`
	SSLPolicy      string                             `pulumi:"sslPolicy"`
	Tags           map[string]string                  `pulumi:"tags"`
}

type ApplicationLoadBalancerArgs struct {
	AccessLogs               lb.LoadBalancerAccessLogsPtrInput `pulumi:"accessLogs"`
	CustomerOwnedIpv4Pool    string                            `pulumi:"customerOwnedIpv4Pool"`
	DefaultSecurityGroup     DefaultSecurityGroupInputs        `pulumi:"defaultSecurityGroup"`
	DefaultTargetGroup       TargetGroupInputs                 `pulumi:"defaultTargetGroup"`
	DesyncMitigationMode     string                            `pulumi:"desyncMitigationMode"`
	DropInvalidHeaderFields  bool                              `pulumi:"dropInvalidHeaderFields"`
	EnableDeletionProtection bool                              `pulumi:"enableDeletionProtection"`
	EnableHttp2              bool                              `pulumi:"enableHttp2"`
	EnableWafFailOpen        bool                              `pulumi:"enableWafFailOpen"`
	IdleTimeout              int                               `pulumi:"idleTimeout"`
	Internal                 bool                              `pulumi:"internal"`
	IPAddressType            string                            `pulumi:"ipAddressType"`
	Listener                 *ListenerInputs                   `pulumi:"listener"`
	Listeners                []*ListenerInputs                 `pulumi:"listeners"`
	Name                     string                            `pulumi:"name"`
	NamePrefix               string                            `pulumi:"namePrefix"`
	SecurityGroups           []string                          `pulumi:"securityGroups"`
	SubnetIDs                []string                          `pulumi:"subnetIds"`
	SubnetMappings           []lb.LoadBalancerSubnetMapping    `pulumi:"subnetMappings"`
	Subnets                  []ec2.Subnet                      `pulumi:"subnets"`
	Tags                     map[string]string                 `pulumi:"tags"`
}

type ApplicationLoadBalancer struct {
	pulumi.ResourceState

	DefaultSecurityGroup *ec2.SecurityGroup  `pulumi:"defaultSecurityGroup"`
	DefaultTargetGroup   *lb.TargetGroup     `pulumi:"defaultTargetGroup"`
	Listeners            []*lb.Listener      `pulumi:"listeners"`
	LoadBalancer         *lb.LoadBalancer    `pulumi:"loadBalancer"`
	VpcID                pulumi.StringOutput `pulumi:"vpcId"`
}

func NewApplicationLoadBalancer(ctx *pulumi.Context, name string, args *ApplicationLoadBalancerArgs, opts ...pulumi.ResourceOption) (*ApplicationLoadBalancer, error) {
	if args == nil {
		args = &ApplicationLoadBalancerArgs{}
	}

	opts = append(opts, pulumi.Aliases([]pulumi.Alias{
		{Name: pulumi.String("awsx:x:elasticloadbalancingv2:ApplicationLoadBalancer")},
	}))

	component := &ApplicationLoadBalancer{}
	err := ctx.RegisterComponentResource(ApplicationLoadBalancerIdentifier, name, component, opts...)
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

	var securityGroups pulumi.StringArrayOutput
	if (len(args.SecurityGroups) == 0) && !args.DefaultSecurityGroup.Skip {
		defaultSecurityGroup := args.DefaultSecurityGroup
		if (defaultSecurityGroup.Args != nil) && (defaultSecurityGroup.SecurityGroupID != "") {
			return nil, fmt.Errorf("Only one of [defaultSecurityGroup] [args] or [securityGroupId] can be specified")
		}

		var securityGroupIDs []pulumi.StringOutput
		securityGroupID := defaultSecurityGroup.SecurityGroupID
		securityGroupIDs = []pulumi.StringOutput{pulumi.String(securityGroupID).ToStringOutput()}
		if securityGroupID == "" {
			sgArgs := &ec2.SecurityGroupArgs{
				VpcId: component.VpcID,
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
			}

			if defaultSecurityGroup.Args != nil {
				dSgArgs := defaultSecurityGroup.Args
				sgArgs = &ec2.SecurityGroupArgs{
					VpcId:               pulumi.String(dSgArgs.VpcID),
					Description:         pulumi.String(dSgArgs.Description),
					Ingress:             dSgArgs.Ingress,
					Egress:              dSgArgs.Egress,
					Name:                pulumi.String(dSgArgs.Name),
					NamePrefix:          pulumi.String(dSgArgs.NamePrefix),
					RevokeRulesOnDelete: pulumi.BoolPtr(dSgArgs.RevokeRulesOnDelete),
					Tags:                pulumi.ToStringMap(dSgArgs.Tags),
				}
			}

			securityGroup, err := ec2.NewSecurityGroup(ctx, name, sgArgs, opts...)
			if err != nil {
				return nil, err
			}

			securityGroupIDs = []pulumi.StringOutput{securityGroup.ID().ToStringOutput()}
		}
		securityGroups = pulumi.ToStringArrayOutput(securityGroupIDs)
	}

	var lbNamePrefix pulumi.StringPtrInput
	var lbName pulumi.StringPtrInput
	if args.Name != "" {
		lbName = pulumi.StringPtr(args.Name)
	}

	if args.NamePrefix != "" {
		lbNamePrefix = pulumi.StringPtr(args.NamePrefix)
		lbName = nil
	}

	loadBalancerType := "application"
	lbArgs := &lb.LoadBalancerArgs{
		AccessLogs:               args.AccessLogs,
		CustomerOwnedIpv4Pool:    pulumi.String(args.CustomerOwnedIpv4Pool),
		DropInvalidHeaderFields:  pulumi.BoolPtr(args.DropInvalidHeaderFields),
		EnableDeletionProtection: pulumi.BoolPtr(args.EnableDeletionProtection),
		EnableHttp2:              pulumi.BoolPtr(args.EnableHttp2),
		EnableWafFailOpen:        pulumi.BoolPtr(args.EnableWafFailOpen),
		IdleTimeout:              pulumi.IntPtr(args.IdleTimeout),
		Internal:                 pulumi.BoolPtr(args.Internal),
		LoadBalancerType:         pulumi.String(loadBalancerType),
		Name:                     lbName,
		NamePrefix:               lbNamePrefix,
		SecurityGroups:           securityGroups,
		Subnets:                  subnetIDs,
		Tags:                     pulumi.ToStringMap(args.Tags),
	}

	if args.IPAddressType != "" {
		lbArgs.IpAddressType = pulumi.StringPtr(args.IPAddressType)
	}

	if args.DesyncMitigationMode != "" {
		lbArgs.DesyncMitigationMode = pulumi.StringPtr(args.DesyncMitigationMode)
	}

	loadBalancer, err := lb.NewLoadBalancer(ctx, name, lbArgs, opts...)
	if err != nil {
		return nil, err
	}
	component.LoadBalancer = loadBalancer

	defaultProtocol := getDefaultProtocol(args)

	var tgNamePrefix pulumi.StringPtrInput
	var tgName pulumi.StringPtrInput
	if args.DefaultTargetGroup.Name != "" {
		tgName = pulumi.StringPtr(args.DefaultTargetGroup.Name)
	}

	if args.DefaultTargetGroup.NamePrefix != "" {
		tgNamePrefix = pulumi.StringPtr(args.DefaultTargetGroup.NamePrefix)
		tgName = nil
	}

	tgArgs := &lb.TargetGroupArgs{
		VpcId:                          component.VpcID,
		TargetType:                     pulumi.StringPtr("ip"),
		Port:                           pulumi.IntPtr(defaultProtocol.Port),
		Protocol:                       pulumi.StringPtr(defaultProtocol.Protocol),
		ConnectionTermination:          pulumi.BoolPtr(args.DefaultTargetGroup.ConnectionTermination),
		DeregistrationDelay:            pulumi.IntPtr(args.DefaultTargetGroup.DeRegistrationDelay),
		HealthCheck:                    args.DefaultTargetGroup.HealthCheck,
		LambdaMultiValueHeadersEnabled: pulumi.BoolPtr(args.DefaultTargetGroup.LambdaMultiValueHeadersEnabled),
		Name:                           tgName,
		NamePrefix:                     tgNamePrefix,
		PreserveClientIp:               pulumi.StringPtr(args.DefaultTargetGroup.PreserveClientIp),
		SlowStart:                      pulumi.IntPtr(args.DefaultTargetGroup.SlowStart),
		Stickiness:                     args.DefaultTargetGroup.Stickiness,
		Tags:                           pulumi.ToStringMap(args.DefaultTargetGroup.Tags),
	}

	if args.DefaultTargetGroup.VpcID != "" {
		tgArgs.VpcId = pulumi.String(args.DefaultTargetGroup.VpcID)
	}

	if args.DefaultTargetGroup.TargetType != "" {
		tgArgs.TargetType = pulumi.String(args.DefaultTargetGroup.TargetType)
	}

	if args.DefaultTargetGroup.Port > 0 {
		tgArgs.Port = pulumi.IntPtr(args.DefaultTargetGroup.Port)
	}

	if args.DefaultTargetGroup.Protocol != "" {
		tgArgs.Protocol = pulumi.String(args.DefaultTargetGroup.Protocol)
	}

	if args.DefaultTargetGroup.ProtocolVersion != "" {
		tgArgs.ProtocolVersion = pulumi.StringPtr(args.DefaultTargetGroup.ProtocolVersion)
	}

	if args.DefaultTargetGroup.LoadBalancingAlgorithmType != "" {
		tgArgs.LoadBalancingAlgorithmType = pulumi.String(args.DefaultTargetGroup.LoadBalancingAlgorithmType)
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
		listenerProtocol := getListenerProtocol(listener)

		port := listener.Port
		if port == 0 {
			port = listenerProtocol.Port
		}

		protocol := listener.Protocol
		if protocol == "" {
			protocol = listenerProtocol.Protocol
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
			Port:     pulumi.IntPtr(defaultProtocol.Port),
			Protocol: pulumi.String(defaultProtocol.Protocol),
		}, opts...)
		if err != nil {
			return nil, err
		}

		component.Listeners = append(component.Listeners, listenerResource)
	}

	return component, nil
}

type protocolResult struct {
	Port     int
	Protocol string
}

func getDefaultProtocol(args *ApplicationLoadBalancerArgs) *protocolResult {
	listener := args.Listener
	if len(args.Listeners) > 0 {
		listener = args.Listeners[0]
	}
	return getListenerProtocol(listener)
}

func getListenerProtocol(listener *ListenerInputs) *protocolResult {
	port := 80
	protocol := "HTTP"
	if listener != nil && (listener.Port > 0 || listener.Protocol != "") {
		port = listener.Port
		protocol = listener.Protocol

		if port == 0 {
			port = protocolToPort(protocol)
		}

		if protocol == "" {
			protocol = portToProtocol(port)
		}
	}

	return &protocolResult{
		Port:     port,
		Protocol: protocol,
	}
}

func portToProtocol(port int) string {
	if port == 443 {
		return "HTTPS"
	}

	return "HTTP"
}

func protocolToPort(protocol string) int {
	if protocol == "HTTPS" {
		return 443
	}

	return 80
}
