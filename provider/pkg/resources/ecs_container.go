package resources

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type TaskDefinitionContainerDependencyInputs struct {
	Condition     string `pulumi:"condition"`
	ContainerName string `pulumi:"containerName"`
}

type TaskDefinitionKeyValuePairInputs struct {
	Name  string `pulumi:"name"`
	Value string `pulumi:"value"`
}

type TaskDefinitionEnvironmentFileInputs struct {
	Type  string `pulumi:"type"`
	Value string `pulumi:"value"`
}

type TaskDefinitionHostEntryInputs struct {
	Hostname  string `pulumi:"hostname"`
	IPAddress string `pulumi:"ipAddress"`
}

type TaskDefinitionFirelensConfigurationInputs struct {
	Options interface{} `pulumi:"options"`
	Type    string      `pulumi:"type"`
}

type TaskDefinitionHealthCheckInputs struct {
	Command     []string `pulumi:"command"`
	Interval    int      `pulumi:"interval"`
	Retries     int      `pulumi:"retries"`
	StartPeriod int      `pulumi:"startPeriod"`
	Timeout     int      `pulumi:"timeout"`
}

type TaskDefinitionKernelCapabilitiesInputs struct {
	Add  []string `pulumi:"add"`
	Drop []string `pulumi:"drop"`
}

type TaskDefinitionDeviceInputs struct {
	ContainerPath string   `pulumi:"containerPath"`
	HostPath      string   `pulumi:"hostPath"`
	Permissions   []string `pulumi:"permissions"`
}

type TaskDefinitionTmpfsInputs struct {
	ContainerPath string   `pulumi:"containerPath"`
	MountOptions  []string `pulumi:"mountOptions"`
	Size          int      `pulumi:"size"`
}

type TaskDefinitionLinuxParametersInputs struct {
	Capabilities       TaskDefinitionKernelCapabilitiesInputs `pulumi:"capabilities"`
	Devices            []TaskDefinitionDeviceInputs           `pulumi:"devices"`
	InitProcessEnabled bool                                   `pulumi:"initProcessEnabled"`
	MaxSwap            int                                    `pulumi:"maxSwap"`
	SharedMemorySize   int                                    `pulumi:"sharedMemorySize"`
	Swappiness         int                                    `pulumi:"swappiness"`
	TMPFS              TaskDefinitionTmpfsInputs              `pulumi:"tmpfs"`
}

type TaskDefinitionSecretInputs struct {
	Name      string `pulumi:"name"`
	ValueFrom string `pulumi:"valueFrom"`
}

type TaskDefinitionLogConfigurationInputs struct {
	LogDriver     string                       `pulumi:"logDriver"`
	Options       interface{}                  `pulumi:"options"`
	SecretOptions []TaskDefinitionSecretInputs `pulumi:"secretOptions"`
}

type TaskDefinitionMountPointInputs struct {
	ContainerPath string `pulumi:"containerPath"`
	ReadOnly      bool   `pulumi:"readOnly"`
	SourceVolume  string `pulumi:"sourceVolume"`
}

type TaskDefinitionPortMappingInputs struct {
	ContainerPort pulumi.IntPtrInput `pulumi:"containerPort"`
	HostPort      pulumi.IntPtrInput `pulumi:"hostPort"`
	Protocol      pulumi.String      `pulumi:"protocol"`
	TargetGroup   *lb.TargetGroup    `pulumi:"targetGroup"`
}

type TaskDefinitionRepositoryCredentialsInputs struct {
	CredentialsParameter string `pulumi:"credentialsParameter"`
}

type TaskDefinitionResourceRequirementInputs struct {
	Type  string `pulumi:"type"`
	Value string `pulumi:"value"`
}

type TaskDefinitionSystemControlInputs struct {
	Namespace string `pulumi:"namespace"`
	Value     string `pulumi:"value"`
}

type TaskDefinitionUlimitInputs struct {
	HardLimit int    `pulumi:"hardLimit"`
	Name      string `pulumi:"name"`
	SoftLimit int    `pulumi:"softLimit"`
}

type TaskDefinitionVolumeFromInputs struct {
	ReadOnly        bool   `pulumi:"readOnly"`
	SourceContainer string `pulumi:"sourceContainer"`
}

type TaskDefinitionContainerDefinitionInputs struct {
	Command                []string                                  `pulumi:"command"`
	CPU                    int                                       `pulumi:"cpu"`
	DependsOn              []TaskDefinitionContainerDependencyInputs `pulumi:"dependsOn"`
	DisableNetworking      bool                                      `pulumi:"disableNetworking"`
	DnsSearchDomains       []string                                  `pulumi:"dnsSearchDomains"`
	DnsServers             []string                                  `pulumi:"dnsServers"`
	DockerLabels           map[string]string                         `pulumi:"dockerLabels"`
	DockerSecurityOptions  []string                                  `pulumi:"dockerSecurityOptions"`
	EntryPoint             []string                                  `pulumi:"entryPoint"`
	Environment            []TaskDefinitionKeyValuePairInputs        `pulumi:"environment"`
	EnvironmentFiles       []TaskDefinitionEnvironmentFileInputs     `pulumi:"environmentFiles"`
	Essential              bool                                      `pulumi:"essential"`
	ExtraHosts             []TaskDefinitionHostEntryInputs           `pulumi:"extraHosts"`
	FirelensConfiguration  TaskDefinitionFirelensConfigurationInputs `pulumi:"firelensConfiguration"`
	HealthCheck            TaskDefinitionHealthCheckInputs           `pulumi:"healthCheck"`
	Hostname               string                                    `pulumi:"hostname"`
	Image                  string                                    `pulumi:"image"`
	Interactive            bool                                      `pulumi:"interactive"`
	Links                  []string                                  `pulumi:"links"`
	LinuxParameters        TaskDefinitionLinuxParametersInputs       `pulumi:"linuxParameters"`
	LogConfiguration       *TaskDefinitionLogConfigurationInputs     `pulumi:"logConfiguration"`
	Memory                 int                                       `pulumi:"memory"`
	MemoryReservation      int                                       `pulumi:"memoryReservation"`
	MountPoints            []TaskDefinitionMountPointInputs          `pulumi:"mountPoints"`
	Name                   string                                    `pulumi:"name"`
	PortMappings           []TaskDefinitionPortMappingInputs         `pulumi:"portMappings"`
	Privileged             bool                                      `pulumi:"privileged"`
	PseudoTerminal         bool                                      `pulumi:"pseudoTerminal"`
	ReadonlyRootFilesystem bool                                      `pulumi:"readonlyRootFilesystem"`
	RepositoryCredentials  TaskDefinitionRepositoryCredentialsInputs `pulumi:"repositoryCredentials"`
	ResourceRequirements   []TaskDefinitionResourceRequirementInputs `pulumi:"resourceRequirements"`
	Secrets                []TaskDefinitionSecretInputs              `pulumi:"secrets"`
	StartTimeout           int                                       `pulumi:"startTimeout"`
	StopTimeout            int                                       `pulumi:"stopTimeout"`
	SystemControls         []TaskDefinitionSystemControlInputs       `pulumi:"systemControls"`
	Ulimits                []TaskDefinitionUlimitInputs              `pulumi:"ulimits"`
	User                   string                                    `pulumi:"user"`
	VolumesFrom            []TaskDefinitionVolumeFromInputs          `pulumi:"volumesFrom"`
	WorkingDirectory       string                                    `pulumi:" workingDirectory"`
}

func computeContainerDefinitions(parent pulumi.Resource, containers map[string]TaskDefinitionContainerDefinitionInputs, logGroupID *LogGroupID) []TaskDefinitionContainerDefinitionInputs {
	var result []TaskDefinitionContainerDefinitionInputs
	for containerName, container := range containers {
		result = append(result, computeContainerDefinition(parent, containerName, container, logGroupID))
	}
	return result
}

func computeContainerDefinition(parent pulumi.Resource, containerName string, container TaskDefinitionContainerDefinitionInputs, logGroupID *LogGroupID) TaskDefinitionContainerDefinitionInputs {
	var resolvedMappings []TaskDefinitionPortMappingInputs
	for _, mappingInput := range container.PortMappings {
		containerPort := pulumi.All(mappingInput.ContainerPort, mappingInput.TargetGroup.Port, mappingInput.HostPort).ApplyT(func(args []interface{}) *int {
			containerPort := args[0].(*int)
			targetGroupPort := args[1].(*int)
			hostPort := args[2].(*int)

			if containerPort != nil {
				return containerPort
			}

			if targetGroupPort != nil {
				return targetGroupPort
			}

			return hostPort
		}).(pulumi.IntPtrInput)

		hostPort := pulumi.All(mappingInput.TargetGroup.Port, mappingInput.HostPort).ApplyT(func(args []interface{}) *int {
			targetGroupPort := args[0].(*int)
			hostPort := args[1].(*int)

			if targetGroupPort != nil {
				return targetGroupPort
			}

			return hostPort
		}).(pulumi.IntPtrInput)

		resolvedMappings = append(resolvedMappings, TaskDefinitionPortMappingInputs{
			ContainerPort: containerPort,
			HostPort:      hostPort,
			Protocol:      mappingInput.Protocol,
		})
	}

	container.PortMappings = resolvedMappings
	container.Name = containerName

	if container.LogConfiguration == nil && logGroupID != nil {
		container.LogConfiguration = &TaskDefinitionLogConfigurationInputs{
			LogDriver: "awsLogs",
			Options: map[string]pulumi.StringInput{
				"awslogs-group":         logGroupID.LogGroupName,
				"awslogs-region":        logGroupID.LogGroupRegion,
				"awslogs-stream-prefix": pulumi.String(containerName),
			},
		}
	}

	return container
}

type mappedContainerDef struct {
	ContainerName string
	TGArn         pulumi.StringOutput
	TGPort        pulumi.IntPtrOutput
}

func computeLoadBalancers(containers map[string]TaskDefinitionContainerDefinitionInputs) ecs.ServiceLoadBalancerArrayOutput {
	var mappedContainers []mappedContainerDef
	for containerName, containerDefinition := range containers {
		portMappings := containerDefinition.PortMappings

		if len(portMappings) == 0 {
			return ecs.ServiceLoadBalancerArrayOutput{}
		}

		for _, mapping := range portMappings {
			mappedContainers = append(mappedContainers, mappedContainerDef{
				ContainerName: containerName,
				TGArn:         mapping.TargetGroup.Arn,
				TGPort:        mapping.TargetGroup.Port,
			})
		}
	}

	result := ecs.ServiceLoadBalancerArrayOutput{}
	return pulumi.All(result, mappedContainers).ApplyT(func(lbs []ecs.ServiceLoadBalancerOutput, mappedContainers []mappedContainerDef) []ecs.ServiceLoadBalancerOutput {
		for _, containerGroup := range mappedContainers {
			slb := pulumi.All(containerGroup.TGArn, containerGroup.TGPort).ApplyT(func(tgArn string, tgPort int) ecs.ServiceLoadBalancer {
				return ecs.ServiceLoadBalancer{
					ContainerName:  containerGroup.ContainerName,
					TargetGroupArn: &tgArn,
					ContainerPort:  tgPort,
				}
			}).(ecs.ServiceLoadBalancerOutput)

			lbs = append(lbs, slb)
		}

		return lbs
	}).(ecs.ServiceLoadBalancerArrayOutput)
}
