// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

// Export members:
export * from "./defaultVpc";
export * from "./getDefaultVpc";
export * from "./vpc";

// Export enums:
export * from "../types/enums/ec2";

// Import resources to register:
import { DefaultVpc } from "./defaultVpc";
import { Vpc } from "./vpc";

const _module = {
    version: utilities.getVersion(),
    construct: (name: string, type: string, urn: string): pulumi.Resource => {
        switch (type) {
            case "awsx-go:ec2:DefaultVpc":
                return new DefaultVpc(name, <any>undefined, { urn })
            case "awsx-go:ec2:Vpc":
                return new Vpc(name, <any>undefined, { urn })
            default:
                throw new Error(`unknown resource type ${type}`);
        }
    },
};
pulumi.runtime.registerResourceModule("awsx-go", "ec2", _module)
