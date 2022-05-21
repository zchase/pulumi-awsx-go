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
	pulumi.Output

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

	subnets := subnetOutput.Ids().ApplyT(func(subnetIDs ec2.LookupSubnetResult) pulumi.AnyOutput {
		return pulumi.All(subnetIDs).ApplyT(func(ids []string) []ec2.LookupSubnetResultOutput {
			var result []ec2.LookupSubnetResultOutput
			for _, id := range ids {
				lookupResult := ec2.LookupSubnetOutput(ctx, ec2.LookupSubnetOutputArgs{
					Id: pulumi.StringPtr(id),
				}, opts...)

				result = append(result, lookupResult)
			}

			return result
		}).(pulumi.AnyOutput)
	}).(pulumi.AnyOutput)

	subnetValues := subnets.ApplyT(func(subs []ec2.LookupSubnetResult) defaultSubnetOutput {
		var publicSubnetIDs []pulumi.StringOutput
		var privateSubnetIDs []pulumi.StringOutput

		for _, s := range subs {
			if s.MapPublicIpOnLaunch {
				publicSubnetIDs = append(publicSubnetIDs, pulumi.String(s.Id).ToStringOutput())
				continue
			}

			privateSubnetIDs = append(privateSubnetIDs, pulumi.String(s.Id).ToStringOutput())
		}

		return defaultSubnetOutput{
			PublicSubnetIDs:  pulumi.ToStringArrayOutput(publicSubnetIDs),
			PrivateSubnetIDs: pulumi.ToStringArrayOutput(privateSubnetIDs),
		}
	}).(defaultSubnetOutput)

	return &DefaultVPCOutput{
		VPCID:            pulumi.String(vpc.Id).ToStringOutput(),
		PublicSubnetIDs:  subnetValues.PublicSubnetIDs,
		PrivateSubnetIDs: subnetValues.PrivateSubnetIDs,
	}, nil
}
