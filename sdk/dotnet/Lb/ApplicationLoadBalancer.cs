// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.AwsxGo.Lb
{
    /// <summary>
    /// Provides an Application Load Balancer resource with listeners, default target group and default security group.
    /// </summary>
    [AwsxGoResourceType("awsx-go:lb:ApplicationLoadBalancer")]
    public partial class ApplicationLoadBalancer : Pulumi.ComponentResource
    {
        /// <summary>
        /// Default security group, if auto-created
        /// </summary>
        [Output("defaultSecurityGroup")]
        public Output<Pulumi.Aws.Ec2.SecurityGroup?> DefaultSecurityGroup { get; private set; } = null!;

        /// <summary>
        /// Default target group, if auto-created
        /// </summary>
        [Output("defaultTargetGroup")]
        public Output<Pulumi.Aws.LB.TargetGroup> DefaultTargetGroup { get; private set; } = null!;

        /// <summary>
        /// Listeners created as part of this load balancer
        /// </summary>
        [Output("listeners")]
        public Output<ImmutableArray<Pulumi.Aws.LB.Listener>> Listeners { get; private set; } = null!;

        /// <summary>
        /// Underlying Load Balancer resource
        /// </summary>
        [Output("loadBalancer")]
        public Output<Pulumi.Aws.LB.LoadBalancer> LoadBalancer { get; private set; } = null!;

        /// <summary>
        /// Id of the VPC in which this load balancer is operating
        /// </summary>
        [Output("vpcId")]
        public Output<string?> VpcId { get; private set; } = null!;


        /// <summary>
        /// Create a ApplicationLoadBalancer resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public ApplicationLoadBalancer(string name, ApplicationLoadBalancerArgs? args = null, ComponentResourceOptions? options = null)
            : base("awsx-go:lb:ApplicationLoadBalancer", name, args ?? new ApplicationLoadBalancerArgs(), MakeResourceOptions(options, ""), remote: true)
        {
        }

        private static ComponentResourceOptions MakeResourceOptions(ComponentResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new ComponentResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = ComponentResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
    }

    public sealed class ApplicationLoadBalancerArgs : Pulumi.ResourceArgs
    {
        /// <summary>
        /// An Access Logs block. Access Logs documented below.
        /// </summary>
        [Input("accessLogs")]
        public Input<Pulumi.Aws.LB.Inputs.LoadBalancerAccessLogsArgs>? AccessLogs { get; set; }

        /// <summary>
        /// The ID of the customer owned ipv4 pool to use for this load balancer.
        /// </summary>
        [Input("customerOwnedIpv4Pool")]
        public Input<string>? CustomerOwnedIpv4Pool { get; set; }

        /// <summary>
        /// Options for creating a default security group if [securityGroups] not specified.
        /// </summary>
        [Input("defaultSecurityGroup")]
        public Pulumi.AwsxGo.Inputs.DefaultSecurityGroupArgs? DefaultSecurityGroup { get; set; }

        /// <summary>
        /// Options creating a default target group.
        /// </summary>
        [Input("defaultTargetGroup")]
        public Inputs.TargetGroupArgs? DefaultTargetGroup { get; set; }

        /// <summary>
        /// Determines how the load balancer handles requests that might pose a security risk to an application due to HTTP desync. Valid values are `monitor`, `defensive` (default), `strictest`.
        /// </summary>
        [Input("desyncMitigationMode")]
        public Input<string>? DesyncMitigationMode { get; set; }

        /// <summary>
        /// Indicates whether HTTP headers with header fields that are not valid are removed by the load balancer (true) or routed to targets (false). The default is false. Elastic Load Balancing requires that message header names contain only alphanumeric characters and hyphens. Only valid for Load Balancers of type `application`.
        /// </summary>
        [Input("dropInvalidHeaderFields")]
        public Input<bool>? DropInvalidHeaderFields { get; set; }

        /// <summary>
        /// If true, deletion of the load balancer will be disabled via
        /// the AWS API. This will prevent this provider from deleting the load balancer. Defaults to `false`.
        /// </summary>
        [Input("enableDeletionProtection")]
        public Input<bool>? EnableDeletionProtection { get; set; }

        /// <summary>
        /// Indicates whether HTTP/2 is enabled in `application` load balancers. Defaults to `true`.
        /// </summary>
        [Input("enableHttp2")]
        public Input<bool>? EnableHttp2 { get; set; }

        /// <summary>
        /// Indicates whether to allow a WAF-enabled load balancer to route requests to targets if it is unable to forward the request to AWS WAF. Defaults to `false`.
        /// </summary>
        [Input("enableWafFailOpen")]
        public Input<bool>? EnableWafFailOpen { get; set; }

        /// <summary>
        /// The time in seconds that the connection is allowed to be idle. Only valid for Load Balancers of type `application`. Default: 60.
        /// </summary>
        [Input("idleTimeout")]
        public Input<int>? IdleTimeout { get; set; }

        /// <summary>
        /// If true, the LB will be internal.
        /// </summary>
        [Input("internal")]
        public Input<bool>? Internal { get; set; }

        /// <summary>
        /// The type of IP addresses used by the subnets for your load balancer. The possible values are `ipv4` and `dualstack`
        /// </summary>
        [Input("ipAddressType")]
        public Input<string>? IpAddressType { get; set; }

        /// <summary>
        /// A listener to create. Only one of [listener] and [listeners] can be specified.
        /// </summary>
        [Input("listener")]
        public Inputs.ListenerArgs? Listener { get; set; }

        [Input("listeners")]
        private List<Inputs.ListenerArgs>? _listeners;

        /// <summary>
        /// List of listeners to create. Only one of [listener] and [listeners] can be specified.
        /// </summary>
        public List<Inputs.ListenerArgs> Listeners
        {
            get => _listeners ?? (_listeners = new List<Inputs.ListenerArgs>());
            set => _listeners = value;
        }

        /// <summary>
        /// The name of the LB. This name must be unique within your AWS account, can have a maximum of 32 characters,
        /// must contain only alphanumeric characters or hyphens, and must not begin or end with a hyphen. If not specified,
        /// this provider will autogenerate a name beginning with `tf-lb`.
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        /// <summary>
        /// Creates a unique name beginning with the specified prefix. Conflicts with `name`.
        /// </summary>
        [Input("namePrefix")]
        public Input<string>? NamePrefix { get; set; }

        [Input("securityGroups")]
        private InputList<string>? _securityGroups;

        /// <summary>
        /// A list of security group IDs to assign to the LB. Only valid for Load Balancers of type `application`.
        /// </summary>
        public InputList<string> SecurityGroups
        {
            get => _securityGroups ?? (_securityGroups = new InputList<string>());
            set => _securityGroups = value;
        }

        [Input("subnetIds")]
        private InputList<string>? _subnetIds;

        /// <summary>
        /// A list of subnet IDs to attach to the LB. Subnets
        /// cannot be updated for Load Balancers of type `network`. Changing this value
        /// for load balancers of type `network` will force a recreation of the resource.
        /// </summary>
        public InputList<string> SubnetIds
        {
            get => _subnetIds ?? (_subnetIds = new InputList<string>());
            set => _subnetIds = value;
        }

        [Input("subnetMappings")]
        private InputList<Pulumi.Aws.LB.Inputs.LoadBalancerSubnetMappingArgs>? _subnetMappings;

        /// <summary>
        /// A subnet mapping block as documented below.
        /// </summary>
        public InputList<Pulumi.Aws.LB.Inputs.LoadBalancerSubnetMappingArgs> SubnetMappings
        {
            get => _subnetMappings ?? (_subnetMappings = new InputList<Pulumi.Aws.LB.Inputs.LoadBalancerSubnetMappingArgs>());
            set => _subnetMappings = value;
        }

        [Input("subnets")]
        private InputList<Pulumi.Aws.Ec2.Subnet>? _subnets;

        /// <summary>
        /// A list of subnets to attach to the LB. Only one of [subnets], [subnetIds] or [subnetMappings] can be specified
        /// </summary>
        public InputList<Pulumi.Aws.Ec2.Subnet> Subnets
        {
            get => _subnets ?? (_subnets = new InputList<Pulumi.Aws.Ec2.Subnet>());
            set => _subnets = value;
        }

        [Input("tags")]
        private InputMap<string>? _tags;

        /// <summary>
        /// A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
        /// </summary>
        public InputMap<string> Tags
        {
            get => _tags ?? (_tags = new InputMap<string>());
            set => _tags = value;
        }

        public ApplicationLoadBalancerArgs()
        {
        }
    }
}
