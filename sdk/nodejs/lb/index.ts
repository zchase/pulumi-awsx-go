// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

// Export members:
export * from "./applicationLoadBalancer";
export * from "./networkLoadBalancer";
export * from "./targetGroupAttachment";

// Import resources to register:
import { ApplicationLoadBalancer } from "./applicationLoadBalancer";
import { NetworkLoadBalancer } from "./networkLoadBalancer";
import { TargetGroupAttachment } from "./targetGroupAttachment";

const _module = {
    version: utilities.getVersion(),
    construct: (name: string, type: string, urn: string): pulumi.Resource => {
        switch (type) {
            case "awsx-go:lb:ApplicationLoadBalancer":
                return new ApplicationLoadBalancer(name, <any>undefined, { urn })
            case "awsx-go:lb:NetworkLoadBalancer":
                return new NetworkLoadBalancer(name, <any>undefined, { urn })
            case "awsx-go:lb:TargetGroupAttachment":
                return new TargetGroupAttachment(name, <any>undefined, { urn })
            default:
                throw new Error(`unknown resource type ${type}`);
        }
    },
};
pulumi.runtime.registerResourceModule("awsx-go", "lb", _module)
