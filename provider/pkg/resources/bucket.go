package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/zchase/pulumi-awsx-go/provider/pkg/utils"
)

type ExistingBucketInputs struct {
	ARN  string `pulumi:"arn"`
	Name string `pulumi:"name"`
}

type BucketArgs struct {
	AccelerationStatus                string                                             `pulumi:"accelerationStatus"`
	ACL                               string                                             `pulumi:"acl"`
	ARN                               string                                             `pulumi:"arn"`
	Bucket                            string                                             `pulumi:"bucket"`
	BucketPrefix                      string                                             `pulumi:"bucketPrefix"`
	CORSRules                         s3.BucketCorsRuleArray                             `pulumi:"corsRules"`
	ForceDestroy                      bool                                               `pulumi:"forceDestroy"`
	Grants                            s3.BucketGrantArray                                `pulumi:"grants"`
	HostedZoneID                      string                                             `pulumi:"hostedZoneId"`
	LifecycleRules                    s3.BucketLifecycleRuleArray                        `pulumi:"lifecycleRules"`
	Loggings                          s3.BucketLoggingArray                              `pulumi:"loggings"`
	ObjectLockConfiguration           s3.BucketObjectLockConfigurationPtrInput           `pulumi:"objectLockConfiguration"`
	Policy                            string                                             `pulumi:"policy"`
	ReplicationConfiguration          s3.BucketReplicationConfigurationPtrInput          `pulumi:"replicationConfiguration"`
	RequestPayer                      string                                             `pulumi:"requestPayer"`
	ServerSideEncryptionConfiguration s3.BucketServerSideEncryptionConfigurationPtrInput `pulumi:"serverSideEncryptionConfiguration"`
	Tags                              map[string]string                                  `pulumi:"tags"`
	Versioning                        s3.BucketVersioningPtrInput                        `pulumi:"versioning"`
	Website                           s3.BucketWebsitePtrInput                           `pulumi:"website"`
	WebsiteDomain                     string                                             `pulumi:"websiteDomain"`
	WebsiteEndpoint                   string                                             `pulumi:"websiteEndpoint"`
}

type RequiredBucketInputs struct {
	Args     *BucketArgs           `pulumi:"args"`
	Existing *ExistingBucketInputs `pulumi:"existing"`
}

type BucketResultBucketID struct {
	ARN  pulumi.StringOutput
	Name pulumi.StringOutput
}

type BucketResult struct {
	Bucket   *s3.Bucket
	BucketID BucketResultBucketID
}

func requiredBucket(ctx *pulumi.Context, name string, inputs *RequiredBucketInputs, opts ...pulumi.ResourceOption) (*BucketResult, error) {
	if (inputs.Args != nil) && (inputs.Existing != nil) {
		return nil, fmt.Errorf("Can't define bucket args if specifying an existing bucket.")
	}

	existing := inputs.Existing
	if existing != nil {
		if existing.ARN != "" {
			arn := pulumi.String(existing.ARN).ToStringOutput()
			return &BucketResult{
				BucketID: BucketResultBucketID{
					ARN: arn,
					Name: arn.ApplyT(func(arn string) (string, error) {
						resource, err := utils.ParseARN(arn)
						if err != nil {
							return "", err
						}

						return resource.ResourceID, nil
					}).(pulumi.StringOutput),
				},
			}, nil
		}

		if existing.Name != "" {
			name := pulumi.String(existing.Name).ToStringOutput()
			return &BucketResult{
				BucketID: BucketResultBucketID{
					Name: name,
					ARN: name.ApplyT(func(name string) string {
						return fmt.Sprintf("arn:aws:::%s", name)
					}).(pulumi.StringOutput),
				},
			}, nil
		}

		return nil, fmt.Errorf("One of an existing log group name or ARN must be specified")
	}

	t := utils.ReturnValueOrDefault()

	bucketArgs := inputs.Args
	bucket, err := s3.NewBucket(ctx, name, &s3.BucketArgs{
		ForceDestroy:                      pulumi.Bool(true),
		AccelerationStatus:                pulumi.StringPtr(bucketArgs.AccelerationStatus),
		Acl:                               pulumi.StringPtr(bucketArgs.ACL),
		Arn:                               pulumi.StringPtr(bucketArgs.ARN),
		Bucket:                            pulumi.StringPtr(bucketArgs.Bucket),
		BucketPrefix:                      pulumi.StringPtr(bucketArgs.BucketPrefix),
		CorsRules:                         bucketArgs.CORSRules,
		Grants:                            bucketArgs.Grants,
		HostedZoneId:                      pulumi.StringPtr(bucketArgs.HostedZoneID),
		LifecycleRules:                    bucketArgs.LifecycleRules,
		Loggings:                          bucketArgs.Loggings,
		ObjectLockConfiguration:           bucketArgs.ObjectLockConfiguration,
		Policy:                            pulumi.StringPtr(bucketArgs.Policy),
		ReplicationConfiguration:          bucketArgs.ReplicationConfiguration,
		RequestPayer:                      pulumi.StringPtr(bucketArgs.RequestPayer),
		ServerSideEncryptionConfiguration: bucketArgs.ServerSideEncryptionConfiguration,
		Tags:                              pulumi.ToStringMap(bucketArgs.Tags),
		Versioning:                        bucketArgs.Versioning,
		Website:                           bucketArgs.Website,
		WebsiteDomain:                     pulumi.StringPtr(bucketArgs.WebsiteDomain),
		WebsiteEndpoint:                   pulumi.StringPtr(bucketArgs.WebsiteEndpoint),
	}, opts...)
	if err != nil {
		return nil, err
	}

	return &BucketResult{
		Bucket: bucket,
		BucketID: BucketResultBucketID{
			ARN:  bucket.Arn,
			Name: bucket.Bucket,
		},
	}, nil
}
