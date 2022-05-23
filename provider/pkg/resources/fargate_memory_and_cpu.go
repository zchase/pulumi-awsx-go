package resources

import (
	"fmt"
	"math"
	"sort"
)

func newRange(start, end int) []float64 {
	result := make([]float64, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, float64(i))
	}
	return result
}

func fargateCost(vcpu, memGB float64) float64 {
	return (0.04048 * vcpu) + (0.004445 * memGB)
}

type cpuMemoryConfig struct {
	Vcpu  float64
	MemGB float64
	Cost  float64
}

func newCpuMemoryConfigs(vcpu float64, memGBs []float64) []cpuMemoryConfig {
	result := make([]cpuMemoryConfig, len(memGBs))
	for _, memGB := range memGBs {
		result = append(result, cpuMemoryConfig{
			Vcpu:  vcpu,
			MemGB: memGB,
			Cost:  fargateCost(vcpu, memGB),
		})
	}
	return result
}

func sortByCost(configs []cpuMemoryConfig) func(x, y int) bool {
	return func(x, y int) bool {
		config1 := configs[x]
		config2 := configs[y]

		if config1.Cost <= config2.Cost {
			return true
		}

		return false
	}
}

func fargateConfigsByPriceAscending() []cpuMemoryConfig {
	allConfigs := []cpuMemoryConfig{}
	allConfigs = append(allConfigs, newCpuMemoryConfigs(0.25, []float64{0.5, 1, 2})...)
	allConfigs = append(allConfigs, newCpuMemoryConfigs(0.5, newRange(1, 4))...)
	allConfigs = append(allConfigs, newCpuMemoryConfigs(1, newRange(2, 8))...)
	allConfigs = append(allConfigs, newCpuMemoryConfigs(2, newRange(4, 16))...)
	allConfigs = append(allConfigs, newCpuMemoryConfigs(4, newRange(8, 30))...)

	sort.SliceStable(allConfigs, sortByCost(allConfigs))
	return allConfigs
}

type fargateContainerMemoryAndCpu struct {
	Cpu               float64
	Memory            float64
	MemoryReservation float64
}

type requestedMemoryAndCpu struct {
	RequestedVCPU float64
	RequestedGB   float64
}

type fargateContainerMemoryAndCpuResult struct {
	Memory string
	Cpu    string
}

func calculateFargateMemoryAndCPU(containers []fargateContainerMemoryAndCpu) (*fargateContainerMemoryAndCpuResult, error) {
	// First, determine how much VCPU/GB that the user is asking for in their containers.
	requested := getRequestedVCPUandMemory(containers)

	// Max CPU that can be requested is only 4.  Don't exceed that.  No need to worry about a
	// min as we're finding the first config that provides *at least* this amount.
	requestedVCPU := math.Min(requested.RequestedVCPU, 4)

	// Max memory that can be requested is only 30.  Don't exceed that.  No need to worry about
	// a min as we're finding the first config that provides *at least* this amount.
	requestedGB := math.Min(requested.RequestedGB, 30)

	// Get all configs that can at least satisfy this pair of cpu/memory needs.
	var config *cpuMemoryConfig
	allConfigs := fargateConfigsByPriceAscending()
	for _, c := range allConfigs {
		if (c.Vcpu >= requestedVCPU) && (c.MemGB >= requestedGB) {
			config = &c
			break
		}
	}

	if config == nil {
		return nil, fmt.Errorf("Could not find fargate config that could satisfy: %v vCPU and %vGB.", requestedVCPU, requestedGB)
	}

	// Want to return docker CPU units, not vCPU values. From AWS:
	//
	// You can determine the number of CPU units that are available per Amazon EC2 instance type by multiplying the
	// number of vCPUs listed for that instance type on the Amazon EC2 Instances detail page by 1,024.
	//
	// We return `memory` in MB units because that appears to be how AWS normalized these internally so this avoids
	// refresh issues.
	return &fargateContainerMemoryAndCpuResult{
		Memory: fmt.Sprintf("%v", config.MemGB*1024),
		Cpu:    fmt.Sprintf("%v", config.Vcpu*1024),
	}, nil
}

func getRequestedVCPUandMemory(containers []fargateContainerMemoryAndCpu) requestedMemoryAndCpu {
	minTaskMemoryMB := float64(0)
	minTaskCPUUnits := float64(0)

	for _, containerDef := range containers {
		if containerDef.MemoryReservation > 0 {
			minTaskMemoryMB += containerDef.MemoryReservation
		} else if containerDef.Memory > 0 {
			minTaskMemoryMB += containerDef.Memory
		}

		minTaskCPUUnits += containerDef.Cpu
	}

	// Convert docker cpu units values into vcpu values.  i.e. 256->.25, 4096->4.
	requestedVCPU := minTaskCPUUnits / 1024

	// Convert memory into GB values.  i.e. 2048MB -> 2GB.
	requestedGB := minTaskMemoryMB / 1024

	return requestedMemoryAndCpu{
		RequestedVCPU: requestedVCPU,
		RequestedGB:   requestedGB,
	}
}
