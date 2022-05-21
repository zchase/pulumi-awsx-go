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
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type vpcEndpointSpecsInput struct {
	AutoAccept        bool              `pulumi:"autoAccept"`
	Policy            string            `pulumi:"policy"`
	PrivateDNSEnabled bool              `pulumi:"privateDnsEnabled"`
	RouteTableIds     []string          `pulumi:"routeTableIds"`
	SecurityGroupIds  []string          `pulumi:"securityGroupIds"`
	ServiceName       string            `pulumi:"serviceName"`
	SubnetIds         []string          `pulumi:"subnetIds"`
	Tags              map[string]string `pulumi:"tags"`
	VpcEndpointType   string            `pulumi:"vpcEndpointType"`
}

type subnetSpecInput struct {
	CIDRMask int    `pulumi:"cidrMask"`
	Name     string `pulumi:"name"`
	Type     string `pulumi:"type"`
}

func (s subnetSpecInput) IsPublic() bool {
	return strings.ToLower(s.Type) == "public"
}

func (s subnetSpecInput) IsPrivate() bool {
	return strings.ToLower(s.Type) == "private"
}

func (s subnetSpecInput) IsIsolated() bool {
	return strings.ToLower(s.Type) == "isolated"
}

type subnetSpec struct {
	CidrBlock  string
	Type       string
	AzName     string
	SubnetName string
}

func (s subnetSpec) IsPublic() bool {
	return strings.ToLower(s.Type) == "public"
}

func (s subnetSpec) IsPrivate() bool {
	return strings.ToLower(s.Type) == "private"
}

func (s subnetSpec) IsIsolated() bool {
	return strings.ToLower(s.Type) == "isolated"
}

type natGatewayInput struct {
	ElasticIpAllocationIds []string `pulumi:"elasticIpAllocationIds"`
	Strategy               string   `pulumi:"strategy"`
}

type VPCArgs struct {
	AssignGeneratedIpv6CidrBlock    bool                    `pulumi:"assignGeneratedIpv6CidrBlock"`
	AvailabilityZoneNames           []string                `pulumi:"availabilityZoneNames"`
	CIDRBlock                       string                  `pulumi:"cidrBlock"`
	EnableClassiclink               bool                    `pulumi:"enableClassiclink"`
	EnableClassiclinkDNSSupport     bool                    `pulumi:"enableClassiclinkDnsSupport"`
	EnableDNSHostnames              bool                    `pulumi:"enableDnsHostnames"`
	EnableDNSSuport                 bool                    `pulumi:"enableDnsSupport"`
	InstanceTenancy                 string                  `pulumi:"instanceTenancy"`
	Ipv4IpamPoolId                  string                  `pulumi:"ipv4IpamPoolId"`
	Ipv4NetmaskLength               int                     `pulumi:"ipv4NetmaskLength"`
	Ipv6CidrBlock                   string                  `pulumi:"ipv6CidrBlock"`
	Ipv6CidrBlockNetworkBorderGroup string                  `pulumi:"ipv6CidrBlockNetworkBorderGroup"`
	Ipv6IpamPoolId                  string                  `pulumi:"ipv6IpamPoolId"`
	Ipv6NetmaskLength               int                     `pulumi:"ipv6NetmaskLength"`
	NatGateways                     natGatewayInput         `pulumi:"natGateways"`
	NumberOfAvailabilityZones       int                     `pulumi:"numberOfAvailabilityZones"`
	SubnetSpecs                     []subnetSpecInput       `pulumi:"subnetSpecs"`
	Tags                            map[string]string       `pulumi:"tags"`
	VpcEndpointSpecs                []vpcEndpointSpecsInput `pulumi:"vpcEndpointSpecs"`
}

type VPCOutput struct {
	pulumi.ResourceState

	EIPS                   []*ec2.Eip                   `pulumi:"eips"`
	InternetGateway        *ec2.InternetGateway         `pulumi:"internetGateway"`
	NatGateways            []*ec2.NatGateway            `pulumi:"natGateways"`
	RouteTableAssociations []*ec2.RouteTableAssociation `pulumi:"routeTableAssociations"`
	RouteTables            []*ec2.RouteTable            `pulumi:"routeTables"`
	Routes                 []*ec2.Route                 `pulumi:"routes"`
	Subnets                []*ec2.Subnet                `pulumi:"subnets"`
	VPC                    *ec2.Vpc                     `pulumi:"vpc"`
	VPCEndpoints           []*ec2.VpcEndpoint           `pulumi:"vpcEndpoints"`
	VPCID                  pulumi.IDOutput              `pulumi:"vpcId"`
	PublicSubnetIDs        pulumi.IDArrayOutput         `pulumi:"publicSubnetIds"`
	PrivateSubnetIDs       pulumi.IDArrayOutput         `pulumi:"privateSubnetIds"`
	IsolatedSubnetIDs      pulumi.IDArrayOutput         `pulumi:"isolatedSubnetIds"`
}
