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
	"strings"
)

const (
	natGatewayStrategyOnePerAZ = "oneperaz"
	natGatewayStrategySingle   = "single"
	natGatewayStrategyNone     = "none"
)

type natGatewayStrategy string

func (n natGatewayStrategy) ToLower() string {
	return strings.ToLower(string(n))
}

func (n natGatewayStrategy) IsNone() bool {
	return n.ToLower() == natGatewayStrategyNone
}

func (n natGatewayStrategy) IsSingle() bool {
	return n.ToLower() == natGatewayStrategySingle
}

func (n natGatewayStrategy) IsOnePerAZ() bool {
	return n.ToLower() == natGatewayStrategyOnePerAZ
}

func (n natGatewayStrategy) isValidStrategyValue() error {
	switch n.ToLower() {
	case natGatewayStrategyOnePerAZ, natGatewayStrategySingle, natGatewayStrategyNone:
		return nil
	default:
		return fmt.Errorf("Unknown NAT Gateway strategy %s", n)
	}
}

func (n natGatewayStrategy) ValidateStrategy(subnets []subnetSpec) error {
	err := n.isValidStrategyValue()
	if err != nil {
		return err
	}

	hasStrategy := n.IsSingle() || n.IsOnePerAZ()
	hasPublicSubnets := false
	hasPrivateSubnets := false

	for _, subnet := range subnets {
		if subnet.IsPublic() {
			hasPublicSubnets = true
		}

		if subnet.IsPrivate() {
			if n.IsNone() {
				return fmt.Errorf("If private subnets are specified, NAT Gateway strategy cannot be 'None'.")
			}

			hasPrivateSubnets = true
		}
	}

	if hasStrategy && (!hasPublicSubnets || !hasPrivateSubnets) {
		return fmt.Errorf("If NAT Gateway strategy is 'OnePerAz' or 'Single', both private and public subnets must be declared. The private subnet creates the need for a NAT Gateway, and the public subnet is required to host the NAT Gateway resource.")
	}

	return nil
}

func (n natGatewayStrategy) ValidateEips(eips, availabilityZones []string) error {
	err := n.isValidStrategyValue()
	if err != nil {
		return err
	}

	eipsLength := len(eips)
	availabilityZonesLength := len(availabilityZones)

	if n.IsOnePerAZ() {
		if (eipsLength > 0) && (eipsLength != availabilityZonesLength) {
			return fmt.Errorf("The number of Elastic IPs, if specified, must match the number of availability zones for the VPC (%v) when NAT Gateway strategy is '%s'", len(availabilityZones), n)
		}
	}

	if n.IsSingle() {
		if eipsLength > 1 {
			return fmt.Errorf("Exactly one Elastic IP may be specified when NAT Gateway strategy is '%s'.", n)
		}
	}

	if n.IsNone() {
		if eipsLength > 0 {
			return fmt.Errorf("Elastic IP allocation IDs cannot be specified when NAT Gateway strategy is %s.", n)
		}
	}

	return nil
}

func (n natGatewayStrategy) ShouldCreateNatGateway(numGateways, azIndex int) (bool, error) {
	err := n.isValidStrategyValue()
	if err != nil {
		return false, err
	}

	if n.IsOnePerAZ() {
		return numGateways < (azIndex + 1), nil
	}

	if n.IsSingle() {
		return numGateways < 1, nil
	}

	if n.IsNone() {
		return false, nil
	}

	return false, nil
}
