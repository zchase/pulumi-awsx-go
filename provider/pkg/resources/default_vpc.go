package resources

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

const defaultVPCIdentifier = "awsx-go:ec2:DefaultVpc"

type DefaultVPCArgs struct{}

type DefaultVPC struct {
	pulumi.ResourceState

	VPCID            pulumi.StringOutput      `pulumi:"vpcId"`
	PrivateSubnetIDs pulumi.StringArrayOutput `pulumi:"privateSubnetIds"`
	PublicSubnetIDs  pulumi.StringArrayOutput `pulumi:"publicSubnetIds"`
}

func NewDefaultVPC(ctx *pulumi.Context, name string, args *DefaultVPCArgs, opts ...pulumi.ResourceOption) (*DefaultVPC, error) {
	if args == nil {
		args = &DefaultVPCArgs{}
	}

	component := &DefaultVPC{}
	err := ctx.RegisterComponentResource(defaultVPCIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	defaultVPC, err := getDefaultVPC(ctx)
	if err != nil {
		return nil, err
	}

	component.VPCID = defaultVPC.VPCID
	component.PublicSubnetIDs = defaultVPC.PublicSubnetIDs
	component.PrivateSubnetIDs = defaultVPC.PrivateSubnetIDs

	return component, nil
}
