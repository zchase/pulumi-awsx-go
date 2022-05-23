import * as pulumi from "@pulumi/pulumi";
import * as awsx from "@pulumi/awsx-go";

const trail = new awsx.cloudtrail.Trail("nodejs-trail", {
    enableLogging: true,
})
