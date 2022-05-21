// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

/**
 * Pseudo resource representing the default VPC and associated subnets for an account and region. This does not create any resources. This will be replaced with `getDefaultVpc` in the future.
 */
export class DefaultVpc extends pulumi.ComponentResource {
    /** @internal */
    public static readonly __pulumiType = 'awsx-go:ec2:DefaultVpc';

    /**
     * Returns true if the given object is an instance of DefaultVpc.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is DefaultVpc {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === DefaultVpc.__pulumiType;
    }

    public /*out*/ readonly privateSubnetIds!: pulumi.Output<string[]>;
    public /*out*/ readonly publicSubnetIds!: pulumi.Output<string[]>;
    /**
     * The VPC ID for the default VPC
     */
    public /*out*/ readonly vpcId!: pulumi.Output<string>;

    /**
     * Create a DefaultVpc resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args?: DefaultVpcArgs, opts?: pulumi.ComponentResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            resourceInputs["privateSubnetIds"] = undefined /*out*/;
            resourceInputs["publicSubnetIds"] = undefined /*out*/;
            resourceInputs["vpcId"] = undefined /*out*/;
        } else {
            resourceInputs["privateSubnetIds"] = undefined /*out*/;
            resourceInputs["publicSubnetIds"] = undefined /*out*/;
            resourceInputs["vpcId"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(DefaultVpc.__pulumiType, name, resourceInputs, opts, true /*remote*/);
    }
}

/**
 * The set of arguments for constructing a DefaultVpc resource.
 */
export interface DefaultVpcArgs {
}