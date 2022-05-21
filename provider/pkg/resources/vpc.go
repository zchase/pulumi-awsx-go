// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewVPC(ctx *pulumi.Context, name string, args *VPCArgs, opts ...pulumi.ResourceOption) (*VPCOutput, error) {
	if args == nil {
		args = &VPCArgs{}
	}

	component := &VPCOutput{}
	err := ctx.RegisterComponentResource("awsx-go:ec2:Vpc", name, component, opts...)
	if err != nil {
		return nil, err
	}

	if (len(args.AvailabilityZoneNames) > 0) && (args.NumberOfAvailabilityZones > 0) {
		return nil, fmt.Errorf("Only one of [availabilityZoneNames] and [numberOfAvailabilityZones] can be specified")
	}

	availabilityZones := args.AvailabilityZoneNames
	if len(availabilityZones) == 0 {
		desiredCount := args.NumberOfAvailabilityZones
		if desiredCount == 0 {
			desiredCount = 3
		}

		azs, err := aws.GetAvailabilityZones(ctx, &aws.GetAvailabilityZonesArgs{})
		if err != nil {
			return nil, err
		}

		if len(azs.Names) < desiredCount {
			return nil, fmt.Errorf("The configured region for this provider does not have at least %v Availability Zones. Either specify an explicit list of zones in availabilityZoneNames or choose a region with at least %v AZs.", desiredCount, desiredCount)
		}

		availabilityZones = azs.Names
	}

	allocationIds := args.NatGateways.ElasticIpAllocationIds
	natGatewayStrategy := natGatewayStrategy(args.NatGateways.Strategy)
	if natGatewayStrategy == "" {
		natGatewayStrategy = "OnePerAz"
	}

	err = natGatewayStrategy.ValidateEips(allocationIds, availabilityZones)
	if err != nil {
		return nil, err
	}

	cidrBlock := args.CIDRBlock
	if cidrBlock == "" {
		cidrBlock = "10.0.0.0/16"
	}

	subnetSpecs, err := getSubnetSpecs(name, cidrBlock, availabilityZones, args.SubnetSpecs)
	if err != nil {
		return nil, err
	}

	err = validateSubnets(subnetSpecs)
	if err != nil {
		return nil, err
	}

	err = natGatewayStrategy.ValidateStrategy(subnetSpecs)
	if err != nil {
		return nil, err
	}

	vpcTags := map[string]string{
		"Name": name,
	}
	for tagKey, tagValue := range args.Tags {
		vpcTags[tagKey] = tagValue
	}

	instanceTenancy := args.InstanceTenancy
	if instanceTenancy == "" {
		instanceTenancy = "dedicated"
	}

	vpcArgs := &ec2.VpcArgs{
		CidrBlock:                   pulumi.StringPtr(cidrBlock),
		Tags:                        pulumi.ToStringMap(vpcTags),
		EnableClassiclink:           pulumi.Bool(args.EnableClassiclink),
		EnableClassiclinkDnsSupport: pulumi.Bool(args.EnableClassiclinkDNSSupport),
		EnableDnsHostnames:          pulumi.Bool(args.EnableDNSHostnames),
		EnableDnsSupport:            pulumi.Bool(args.EnableDNSSuport),
		InstanceTenancy:             pulumi.StringPtr(instanceTenancy),
	}

	if args.Ipv6CidrBlock != "" {
		vpcArgs.Ipv6CidrBlock = pulumi.StringPtr(args.Ipv6CidrBlock)
		vpcArgs.Ipv6NetmaskLength = nil
	}

	if args.Ipv4IpamPoolId != "" {
		vpcArgs.Ipv4IpamPoolId = pulumi.StringPtr(args.Ipv4IpamPoolId)
		vpcArgs.Ipv4NetmaskLength = pulumi.IntPtr(args.Ipv4NetmaskLength)
		vpcArgs.CidrBlock = nil
	}

	if args.AssignGeneratedIpv6CidrBlock == true {
		vpcArgs.AssignGeneratedIpv6CidrBlock = pulumi.BoolPtr(args.AssignGeneratedIpv6CidrBlock)
		vpcArgs.Ipv4NetmaskLength = pulumi.IntPtr(args.Ipv6NetmaskLength)
		vpcArgs.Ipv6IpamPoolId = pulumi.StringPtr(args.Ipv6IpamPoolId)
		vpcArgs.Ipv6CidrBlockNetworkBorderGroup = pulumi.StringPtr(args.Ipv6CidrBlockNetworkBorderGroup)
	}

	vpc, err := ec2.NewVpc(ctx, name, vpcArgs, opts...)
	if err != nil {
		return nil, err
	}

	vpcId := vpc.ID()
	vpcChildResourceOptions := []pulumi.ResourceOption{pulumi.Parent(vpc), pulumi.DependsOn([]pulumi.Resource{vpc})}

	igw, err := ec2.NewInternetGateway(ctx, name, &ec2.InternetGatewayArgs{
		VpcId: vpcId,
		Tags: pulumi.ToStringMap(map[string]string{
			"Name": name,
		}),
	}, vpcChildResourceOptions...)
	if err != nil {
		return nil, err
	}

	var vpcEndpoints []*ec2.VpcEndpoint
	var subnets []*ec2.Subnet
	var routeTables []*ec2.RouteTable
	var routeTableAssociations []*ec2.RouteTableAssociation
	var routes []*ec2.Route
	var natGateways []*ec2.NatGateway
	var eips []*ec2.Eip
	var publicSubnetIds []pulumi.IDOutput
	var privateSubnetIds []pulumi.IDOutput
	var isolatedSubnetIds []pulumi.IDOutput

	for _, vpcSubnetSpec := range args.VpcEndpointSpecs {
		vpcEndpoint, err := ec2.NewVpcEndpoint(ctx, vpcSubnetSpec.ServiceName, &ec2.VpcEndpointArgs{
			AutoAccept:        pulumi.BoolPtr(vpcSubnetSpec.AutoAccept),
			Policy:            pulumi.Sprintf("%s", vpcSubnetSpec.Policy),
			PrivateDnsEnabled: pulumi.BoolPtr(vpcSubnetSpec.PrivateDNSEnabled),
			RouteTableIds:     pulumi.ToStringArray(vpcSubnetSpec.RouteTableIds),
			SecurityGroupIds:  pulumi.ToStringArray(vpcSubnetSpec.SecurityGroupIds),
			SubnetIds:         pulumi.ToStringArray(vpcSubnetSpec.SubnetIds),
			Tags:              pulumi.ToStringMap(vpcSubnetSpec.Tags),
			VpcEndpointType:   pulumi.Sprintf("%s", vpcSubnetSpec.VpcEndpointType),
			VpcId:             vpcId,
			ServiceName:       pulumi.Sprintf("%s", vpcSubnetSpec.ServiceName),
		}, vpcChildResourceOptions...)
		if err != nil {
			return nil, err
		}

		vpcEndpoints = append(vpcEndpoints, vpcEndpoint)
	}

	for i, zone := range availabilityZones {
		var specs []subnetSpec
		for _, spec := range subnetSpecs {
			if spec.AzName == zone {
				specs = append(specs, spec)
			}
		}

		sort.SliceStable(specs, compareSubnetSpecs(specs))

		for _, spec := range specs {
			subnet, err := ec2.NewSubnet(ctx, spec.SubnetName, &ec2.SubnetArgs{
				VpcId:               vpcId,
				AvailabilityZone:    pulumi.Sprintf("%s", spec.AzName),
				MapPublicIpOnLaunch: pulumi.BoolPtr(strings.ToLower(spec.Type) == "public"),
				CidrBlock:           pulumi.Sprintf("%s", spec.CidrBlock),
				Tags: pulumi.ToStringMap(map[string]string{
					"Name": spec.SubnetName,
				}),
			}, vpcChildResourceOptions...)
			if err != nil {
				return nil, err
			}

			subnets = append(subnets, subnet)

			if spec.IsPublic() {
				publicSubnetIds = append(publicSubnetIds, subnet.ID())
			}

			if spec.IsPrivate() {
				privateSubnetIds = append(privateSubnetIds, subnet.ID())
			}

			if spec.IsIsolated() {
				isolatedSubnetIds = append(isolatedSubnetIds, subnet.ID())
			}

			routeTable, err := ec2.NewRouteTable(ctx, spec.SubnetName, &ec2.RouteTableArgs{
				VpcId: vpcId,
				Tags: pulumi.ToStringMap(map[string]string{
					"Name": spec.SubnetName,
				}),
			}, pulumi.Parent(subnet), pulumi.DependsOn([]pulumi.Resource{subnet}))
			if err != nil {
				return nil, err
			}

			routeTables = append(routeTables, routeTable)

			routeTableAssoc, err := ec2.NewRouteTableAssociation(ctx, spec.SubnetName, &ec2.RouteTableAssociationArgs{
				RouteTableId: routeTable.ID(),
				SubnetId:     subnet.ID(),
			}, pulumi.Parent(routeTable), pulumi.DependsOn([]pulumi.Resource{routeTable}))
			if err != nil {
				return nil, err
			}

			routeTableAssociations = append(routeTableAssociations, routeTableAssoc)

			createNatGateway, err := natGatewayStrategy.ShouldCreateNatGateway(len(natGateways), i)
			if err != nil {
				return nil, err
			}

			if spec.IsPublic() && createNatGateway {
				createEip := len(allocationIds) == 0

				var natGatewayAllocationIDs pulumi.StringOutput
				if createEip {
					eipName := fmt.Sprintf("%s-%v", name, i+1)
					eip, err := ec2.NewEip(ctx, eipName, &ec2.EipArgs{}, pulumi.Parent(subnet), pulumi.DependsOn([]pulumi.Resource{subnet}))
					if err != nil {
						return nil, err
					}
					eips = append(eips, eip)

					natGatewayAllocationIDs = eip.AllocationId
				} else {
					natGatewayAllocationIDs = pulumi.String(allocationIds[i]).ToStringOutput()
				}

				natGatewayName := fmt.Sprintf("%s-nat-gateway-%v", name, i+1)
				natGateway, err := ec2.NewNatGateway(ctx, natGatewayName, &ec2.NatGatewayArgs{
					SubnetId:     subnet.ID(),
					AllocationId: natGatewayAllocationIDs,
					Tags: pulumi.ToStringMap(map[string]string{
						"Name": spec.SubnetName,
					}),
				}, pulumi.Parent(subnet), pulumi.DependsOn([]pulumi.Resource{subnet}))
				if err != nil {
					return nil, err
				}

				natGateways = append(natGateways, natGateway)
			}

			if spec.IsPublic() {
				route, err := ec2.NewRoute(ctx, spec.SubnetName, &ec2.RouteArgs{
					RouteTableId:         routeTable.ID(),
					GatewayId:            igw.ID(),
					DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
				}, pulumi.Parent(routeTable), pulumi.DependsOn([]pulumi.Resource{routeTable}))
				if err != nil {
					return nil, err
				}
				routes = append(routes, route)
			}

			if spec.IsPrivate() {
				var natGatewayID pulumi.IDOutput
				if natGatewayStrategy.IsSingle() {
					natGatewayID = natGateways[0].ID()
				} else {
					natGatewayID = natGateways[i].ID()
				}

				route, err := ec2.NewRoute(ctx, spec.SubnetName, &ec2.RouteArgs{
					RouteTableId:         routeTable.ID(),
					NatGatewayId:         natGatewayID,
					DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
				}, pulumi.Parent(routeTable), pulumi.DependsOn([]pulumi.Resource{routeTable}))
				if err != nil {
					return nil, err
				}
				routes = append(routes, route)
			}
		}
	}

	component.EIPS = eips
	component.InternetGateway = igw
	component.NatGateways = natGateways
	component.RouteTables = routeTables
	component.RouteTableAssociations = routeTableAssociations
	component.Routes = routes
	component.Subnets = subnets
	component.VPC = vpc
	component.VPCEndpoints = vpcEndpoints
	component.VPCID = vpcId
	component.PublicSubnetIDs = pulumi.ToIDArrayOutput(publicSubnetIds)
	component.PrivateSubnetIDs = pulumi.ToIDArrayOutput(privateSubnetIds)
	component.IsolatedSubnetIDs = pulumi.ToIDArrayOutput(isolatedSubnetIds)

	return component, nil
}
