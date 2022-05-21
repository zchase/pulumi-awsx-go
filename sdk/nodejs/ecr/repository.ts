// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import { input as inputs, output as outputs, enums } from "../types";
import * as utilities from "../utilities";

import * as pulumiAws from "@pulumi/aws";

/**
 * A [Repository] represents an [aws.ecr.Repository] along with an associated [LifecyclePolicy] controlling how images are retained in the repo.
 *
 * Docker images can be built and pushed to the repo using the [buildAndPushImage] method.  This will call into the `@pulumi/docker/buildAndPushImage` function using this repo as the appropriate destination registry.
 */
export class Repository extends pulumi.ComponentResource {
    /** @internal */
    public static readonly __pulumiType = 'awsx-go:ecr:Repository';

    /**
     * Returns true if the given object is an instance of Repository.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is Repository {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === Repository.__pulumiType;
    }

    /**
     * Underlying repository lifecycle policy
     */
    public readonly lifecyclePolicy!: pulumi.Output<pulumiAws.ecr.LifecyclePolicy | undefined>;
    /**
     * Underlying Repository resource
     */
    public /*out*/ readonly repository!: pulumi.Output<pulumiAws.ecr.Repository>;
    /**
     * The URL of the repository (in the form aws_account_id.dkr.ecr.region.amazonaws.com/repositoryName).
     */
    public /*out*/ readonly url!: pulumi.Output<string>;

    /**
     * Create a Repository resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args?: RepositoryArgs, opts?: pulumi.ComponentResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            resourceInputs["encryptionConfigurations"] = args ? args.encryptionConfigurations : undefined;
            resourceInputs["imageScanningConfiguration"] = args ? args.imageScanningConfiguration : undefined;
            resourceInputs["imageTagMutability"] = args ? args.imageTagMutability : undefined;
            resourceInputs["lifecyclePolicy"] = args ? args.lifecyclePolicy : undefined;
            resourceInputs["name"] = args ? args.name : undefined;
            resourceInputs["tags"] = args ? args.tags : undefined;
            resourceInputs["repository"] = undefined /*out*/;
            resourceInputs["url"] = undefined /*out*/;
        } else {
            resourceInputs["lifecyclePolicy"] = undefined /*out*/;
            resourceInputs["repository"] = undefined /*out*/;
            resourceInputs["url"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(Repository.__pulumiType, name, resourceInputs, opts, true /*remote*/);
    }
}

/**
 * The set of arguments for constructing a Repository resource.
 */
export interface RepositoryArgs {
    /**
     * Encryption configuration for the repository. See below for schema.
     */
    encryptionConfigurations?: pulumi.Input<pulumi.Input<pulumiAws.types.input.ecr.RepositoryEncryptionConfiguration>[]>;
    /**
     * Configuration block that defines image scanning configuration for the repository. By default, image scanning must be manually triggered. See the [ECR User Guide](https://docs.aws.amazon.com/AmazonECR/latest/userguide/image-scanning.html) for more information about image scanning.
     */
    imageScanningConfiguration?: pulumi.Input<pulumiAws.types.input.ecr.RepositoryImageScanningConfiguration>;
    /**
     * The tag mutability setting for the repository. Must be one of: `MUTABLE` or `IMMUTABLE`. Defaults to `MUTABLE`.
     */
    imageTagMutability?: pulumi.Input<string>;
    /**
     * A lifecycle policy consists of one or more rules that determine which images in a repository should be expired. If not provided, this will default to untagged images expiring after 1 day.
     */
    lifecyclePolicy?: inputs.ecr.LifecyclePolicyArgs;
    /**
     * Name of the repository.
     */
    name?: pulumi.Input<string>;
    /**
     * A map of tags to assign to the resource. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
     */
    tags?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
}