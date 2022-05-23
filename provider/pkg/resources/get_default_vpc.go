package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DefaultVPCOutput struct {
	VPCID            pulumi.StringOutput
	PrivateSubnetIDs pulumi.StringArrayOutput
	PublicSubnetIDs  pulumi.StringArrayOutput
}

type defaultSubnetOutput struct {
	PrivateSubnetIDs pulumi.StringArrayOutput
	PublicSubnetIDs  pulumi.StringArrayOutput
}

func getDefaultVPC(ctx *pulumi.Context, opts ...pulumi.InvokeOption) (*DefaultVPCOutput, error) {
	vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{
		Default: pulumi.BoolRef(true),
	}, opts...)
	if err != nil {
		return nil, err
	}

	if vpc == nil {
		return nil, fmt.Errorf("unable to find default VPC for this region and account")
	}

	subnetOutput := ec2.GetSubnetsOutput(ctx, ec2.GetSubnetsOutputArgs{
		Filters: ec2.GetSubnetsFilterArray{
			&ec2.GetSubnetsFilterArgs{
				Name:   pulumi.String("vpc-id"),
				Values: pulumi.ToStringArray([]string{vpc.Id}),
			},
		},
	}, opts...)

	subnets := subnetOutput.Ids().ApplyT(func(subnetIDs []string) []ec2.LookupSubnetResultOutput {
		var result []ec2.LookupSubnetResultOutput
		for _, id := range subnetIDs {
			lookupResult := ec2.LookupSubnetOutput(ctx, ec2.LookupSubnetOutputArgs{
				Id: pulumi.StringPtr(id),
			}, opts...)

			result = append(result, lookupResult)
		}

		return result
	}).(pulumi.AnyOutput)

	subnetValues := subnets.ApplyT(func(subs interface{}) defaultSubnetOutput {
		var publicSubnetIDs []pulumi.StringOutput
		var privateSubnetIDs []pulumi.StringOutput

		switch subs := subs.(type) {
		case []ec2.LookupSubnetResult:
			for _, s := range subs {
				fmt.Println(s.Id)

				if s.MapPublicIpOnLaunch {
					publicSubnetIDs = append(publicSubnetIDs, pulumi.String(s.Id).ToStringOutput())
					continue
				}

				privateSubnetIDs = append(privateSubnetIDs, pulumi.String(s.Id).ToStringOutput())
			}
		}

		return defaultSubnetOutput{
			PublicSubnetIDs:  pulumi.ToStringArrayOutput(publicSubnetIDs),
			PrivateSubnetIDs: pulumi.ToStringArrayOutput(privateSubnetIDs),
		}
	}).(pulumi.AnyOutput)

	return &DefaultVPCOutput{
		VPCID: pulumi.String(vpc.Id).ToStringOutput(),
		PublicSubnetIDs: subnetValues.ApplyT(func(v interface{}) pulumi.StringArrayOutput {
			so := v.(defaultSubnetOutput)
			return so.PublicSubnetIDs
		}).(pulumi.StringArrayOutput),
		PrivateSubnetIDs: subnetValues.ApplyT(func(v interface{}) pulumi.StringArrayOutput {
			so := v.(defaultSubnetOutput)
			return so.PrivateSubnetIDs
		}).(pulumi.StringArrayOutput),
	}, nil
}
