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
	"math"
	"math/big"
	"net"
	"strings"
)

func compareSubnetSpecs(specs []subnetSpec) func(x, y int) bool {
	return func(x, y int) bool {
		spec1 := specs[x]
		spec2 := specs[y]

		if spec1.Type == spec2.Type {
			return true
		}

		if spec1.IsPublic() {
			return true
		}

		if spec1.IsPrivate() && spec2.IsPublic() {
			return false
		}

		if spec1.IsPrivate() && spec2.IsIsolated() {
			return true
		}

		return false
	}
}

func validateSubnets(specs []subnetSpec) error {
	overlappingSubnets, err := getOverlappingSubnets(specs)
	if err != nil {
		return err
	}

	if len(overlappingSubnets) > 0 {
		msgParts := []string{
			"The following subnets overlap with at least one other subnet. Make the CIDR for the VPC larger, reduce the size of the subnets per AZ, or use less Availability Zones:\n\n",
		}
		for i, subnet := range overlappingSubnets {
			msgParts = append(msgParts, fmt.Sprintf("%v. %s: %s\n", i+1, subnet.SubnetName, subnet.CidrBlock))
		}

		return fmt.Errorf(strings.Join(msgParts, ""))
	}

	return nil
}

func doSubnetsOverlap(spec1, spec2 subnetSpec) (bool, error) {
	_, ip1, err := net.ParseCIDR(spec1.CidrBlock)
	if err != nil {
		return false, err
	}

	_, ip2, err := net.ParseCIDR(spec2.CidrBlock)
	if err != nil {
		return false, err
	}

	hasOverlap := ip1.Contains(ip2.IP) || ip2.Contains(ip1.IP)
	return hasOverlap, nil
}

func getOverlappingSubnets(specs []subnetSpec) ([]subnetSpec, error) {
	var result []subnetSpec
	for _, x := range specs {
		hasOverlap := false
		for _, y := range specs {
			if x == y {
				continue
			}

			hasOverlap, err := doSubnetsOverlap(x, y)
			if err != nil {
				return nil, err
			}

			if hasOverlap {
				hasOverlap = true
			}
		}

		if hasOverlap {
			result = append(result, x)
		}
	}

	return result, nil
}

func nextPow2(n int) int {
	if n == 0 {
		return 1
	}

	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16

	return n + 1
}

func ip2BigInt(ip net.IP) *big.Int {
	i := big.NewInt(0)
	i.SetBytes(ip)
	return i
}

func cidrSubnetV4(ipRange string, newBits, netNum int) (string, error) {
	_, ip, err := net.ParseCIDR(ipRange)
	if err != nil {
		return "", fmt.Errorf("Error parsing IP range: %v", err)
	}

	ipSubnetMaskBits, _ := ip.Mask.Size()
	newSubnetMask := ipSubnetMaskBits + newBits
	if newSubnetMask > 32 {
		return "", fmt.Errorf("Requested %v new bits, but only %v are available.", newBits, 32-ipSubnetMaskBits)
	}

	addressBI := ip2BigInt(ip.IP)
	newAddressBase := math.Pow(2, float64(32-newSubnetMask))
	netNumBI := big.NewInt(int64(netNum))

	newAddressBaseBI := big.NewInt(int64(newAddressBase))
	newAddressBI := addressBI.Add(
		addressBI,
		newAddressBaseBI.Mul(newAddressBaseBI, netNumBI),
	)

	newAddress := net.IP(newAddressBI.Bytes())

	return fmt.Sprintf("%s/%v", newAddress.String(), newSubnetMask), nil
}

func generateDefaultSubnets(vpcName, vpcCidr string, azNames, azBases []string) ([]subnetSpec, error) {
	var privateSubnets []subnetSpec
	for i, name := range azNames {
		cidrBlock, err := cidrSubnetV4(azBases[i], 1, 0)
		if err != nil {
			return nil, err
		}

		privateSubnets = append(privateSubnets, subnetSpec{
			AzName:     name,
			Type:       "Private",
			SubnetName: fmt.Sprintf("%s-private-%v", vpcName, i+1),
			CidrBlock:  cidrBlock,
		})
	}

	var publicSubnets []subnetSpec
	for i, name := range azNames {
		splitBase, err := cidrSubnetV4(privateSubnets[i].CidrBlock, 0, 1)
		if err != nil {
			return nil, err
		}

		cidrBlock, err := cidrSubnetV4(splitBase, 1, 0)
		if err != nil {
			return nil, err
		}

		publicSubnets = append(publicSubnets, subnetSpec{
			AzName:     name,
			Type:       "Public",
			SubnetName: fmt.Sprintf("%s-public-%v", vpcName, i+1),
			CidrBlock:  cidrBlock,
		})
	}

	return append(privateSubnets, publicSubnets...), nil
}

func getSubnetSpecs(vpcName, vpcCidr string, azNames []string, subnetInputs []subnetSpecInput) ([]subnetSpec, error) {
	newBitsPerAZ := math.Log2(float64(nextPow2(len(azNames))))

	var azBases []string
	for i := range azNames {
		azBase, err := cidrSubnetV4(vpcCidr, int(newBitsPerAZ), i)
		if err != nil {
			return nil, err
		}
		azBases = append(azBases, azBase)
	}

	if len(subnetInputs) == 0 {
		return generateDefaultSubnets(vpcName, vpcCidr, azNames, azBases)
	}

	_, ip, err := net.ParseCIDR(azBases[0])
	if err != nil {
		return nil, fmt.Errorf("Error parsing IP range for non default VPC: %v", err)
	}

	_, baseSubnetMaskBits := ip.Mask.Size()

	var privateSubnetsIn []subnetSpecInput
	var publicSubnetsIn []subnetSpecInput
	var isolatedSubnetsIn []subnetSpecInput
	for _, subnetIn := range subnetInputs {
		if subnetIn.IsPrivate() {
			privateSubnetsIn = append(privateSubnetsIn, subnetIn)
		}

		if subnetIn.IsPublic() {
			publicSubnetsIn = append(publicSubnetsIn, subnetIn)
		}

		if subnetIn.IsIsolated() {
			isolatedSubnetsIn = append(isolatedSubnetsIn, subnetIn)
		}
	}

	var subnetOuts []subnetSpec

	for i, name := range azNames {
		var privateSubnetsOut []subnetSpec
		var publicSubnetsOut []subnetSpec
		var isolatedSubnetsOut []subnetSpec

		// Private subnets
		for j, privateIn := range privateSubnetsIn {
			newBits := privateIn.CIDRMask - baseSubnetMaskBits

			privateSubnetCidrBlock, err := cidrSubnetV4(azBases[i], newBits, j)
			if err != nil {
				return nil, err
			}

			privateSubnetsOut = append(privateSubnetsOut, subnetSpec{
				AzName:     name,
				CidrBlock:  privateSubnetCidrBlock,
				Type:       "Private",
				SubnetName: fmt.Sprintf("%s-%s-%v", vpcName, privateIn.Name, i+1),
			})
		}

		// Public Subnets
		for j, publicIn := range publicSubnetsIn {
			baseCidr := azBases[i]
			if len(privateSubnetsOut) > 0 {
				baseCidr = privateSubnetsOut[len(privateSubnetsOut)-1].CidrBlock
			}

			_, baseIP, err := net.ParseCIDR(baseCidr)
			if err != nil {
				return nil, err
			}

			_, basePublicSubnetMaskBits := baseIP.Mask.Size()

			splitBase := azBases[i]
			if len(privateSubnetsOut) > 0 {
				splitBase, err = cidrSubnetV4(baseCidr, 0, 1)
				if err != nil {
					return nil, err
				}
			}

			newPublicSubnetBits := publicIn.CIDRMask - basePublicSubnetMaskBits
			publicSubnetCidrBlock, err := cidrSubnetV4(splitBase, newPublicSubnetBits, j)
			if err != nil {
				return nil, err
			}

			publicSubnetsOut = append(publicSubnetsOut, subnetSpec{
				AzName:     name,
				CidrBlock:  publicSubnetCidrBlock,
				Type:       "Public",
				SubnetName: fmt.Sprintf("%s-%s-%v", vpcName, publicIn.Name, i+1),
			})
		}

		// Isolated Subnets
		for j, isolatedIn := range isolatedSubnetsIn {
			baseCidr := azBases[i]
			if len(publicSubnetsOut) > 0 {
				baseCidr = publicSubnetsOut[len(publicSubnetsOut)-1].CidrBlock
			} else if len(privateSubnetsOut) > 0 {
				baseCidr = privateSubnetsOut[len(privateSubnetsOut)-1].CidrBlock
			}

			_, baseIP, err := net.ParseCIDR(baseCidr)
			if err != nil {
				return nil, err
			}

			_, baseIsolatedSubnetMaskBits := baseIP.Mask.Size()

			splitBase := azBases[i]
			if (len(publicSubnetsOut) > 0) || (len(privateSubnetsOut) > 0) {
				splitBase, err = cidrSubnetV4(baseCidr, 0, 1)
				if err != nil {
					return nil, err
				}
			}

			newIsolatedSubnetBits := isolatedIn.CIDRMask - baseIsolatedSubnetMaskBits
			isolatedSubnetCidrBlock, err := cidrSubnetV4(splitBase, newIsolatedSubnetBits, j)
			isolatedSubnetsOut = append(isolatedSubnetsOut, subnetSpec{
				AzName:     name,
				CidrBlock:  isolatedSubnetCidrBlock,
				Type:       "Isolated",
				SubnetName: fmt.Sprintf("%s-%s-%v", vpcName, isolatedIn.Name, i+1),
			})
		}

		subnetOuts = append(subnetOuts, privateSubnetsOut...)
		subnetOuts = append(subnetOuts, publicSubnetsOut...)
		subnetOuts = append(subnetOuts, isolatedSubnetsOut...)
	}

	return subnetOuts, nil
}
