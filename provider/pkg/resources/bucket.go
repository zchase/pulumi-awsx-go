package resources

import "github.com/aws/aws-sdk-go/service/s3"

type ExistingBucketInputs struct {
	ARN  string `pulumi:"arn"`
	Name string `pulumi:"name"`
}

type BucketArgs struct {
	AccelerationStatus                string                               `pulumi:"accelerationStatus"`
	ACL                               string                               `pulumi:"acl"`
	ARN                               string                               `pulumi:"arn"`
	Bucket                            string                               `pulumi:"bucket"`
	BucketPrefix                      string                               `pulumi:"bucketPrefix"`
	CORSRules                         []s3.CORSRule                        `pulumi:"corsRules"`
	ForceDestroy                      bool                                 `pulumi:"forceDestroy"`
	Grants                            []s3.Grant                           `pulumi:"grants"`
	HostedZoneID                      string                               `pulumi:"hostedZoneId"`
	LifecycleRules                    []s3.LifecycleRule                   `pulumi:"lifecycleRules"`
	Loggings                          []s3.LoggingEnabled                  `pulumi:"loggings"`
	ObjectLockConfiguration           s3.ObjectLockConfiguration           `pulumi:"objectLockConfiguration"`
	Policy                            string                               `pulumi:"policy"`
	ReplicationConfiguration          s3.ReplicationConfiguration          `pulumi:"replicationConfiguration"`
	RequestPayer                      string                               `pulumi:"requestPayer"`
	ServerSideEncryptionConfiguration s3.ServerSideEncryptionConfiguration `pulumi:"serverSideEncryptionConfiguration"`
	Tags                              map[string]string                    `pulumi:"tags"`
	Versioning                        s3.VersioningConfiguration           `pulumi:"versioning"`
	Website                           s3.WebsiteConfiguration              `pulumi:"website"`
	WebsiteDomain                     string                               `pulumi:"websiteDomain"`
	WebsiteEndpoint                   string                               `pulumi:"websiteEndpoint"`
}

func requiredBucket(name string, inputs *RequiredBucketInputs, defaults aws)
