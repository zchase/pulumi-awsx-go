// Code generated by Pulumi SDK Generator DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package ecr

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Arguments for building a docker image
type DockerBuild struct {
	// An optional map of named build-time argument variables to set during the Docker build.  This flag allows you to pass built-time variables that can be accessed like environment variables inside the `RUN` instruction.
	Args map[string]string `pulumi:"args"`
	// Images to consider as cache sources
	CacheFrom []string `pulumi:"cacheFrom"`
	// dockerfile may be used to override the default Dockerfile name and/or location.  By default, it is assumed to be a file named Dockerfile in the root of the build context.
	Dockerfile *string `pulumi:"dockerfile"`
	// Environment variables to set on the invocation of `docker build`, for example to support `DOCKER_BUILDKIT=1 docker build`.
	Env map[string]string `pulumi:"env"`
	// An optional catch-all list of arguments to provide extra CLI options to the docker build command.  For example `['--network', 'host']`.
	ExtraOptions []string `pulumi:"extraOptions"`
	// Path to a directory to use for the Docker build context, usually the directory in which the Dockerfile resides (although dockerfile may be used to choose a custom location independent of this choice). If not specified, the context defaults to the current working directory; if a relative path is used, it is relative to the current working directory that Pulumi is evaluating.
	Path *string `pulumi:"path"`
	// The target of the dockerfile to build
	Target *string `pulumi:"target"`
}

// Simplified lifecycle policy model consisting of one or more rules that determine which images in a repository should be expired. See https://docs.aws.amazon.com/AmazonECR/latest/userguide/lifecycle_policy_examples.html for more details.
type LifecyclePolicy struct {
	// Specifies the rules to determine how images should be retired from this repository. Rules are ordered from lowest priority to highest. If there is a rule with a `selection` value of `any`, then it will have the highest priority.
	Rules []LifecyclePolicyRule `pulumi:"rules"`
	// Skips creation of the policy if set to `true`.
	Skip *bool `pulumi:"skip"`
}

// LifecyclePolicyInput is an input type that accepts LifecyclePolicyArgs and LifecyclePolicyOutput values.
// You can construct a concrete instance of `LifecyclePolicyInput` via:
//
//          LifecyclePolicyArgs{...}
type LifecyclePolicyInput interface {
	pulumi.Input

	ToLifecyclePolicyOutput() LifecyclePolicyOutput
	ToLifecyclePolicyOutputWithContext(context.Context) LifecyclePolicyOutput
}

// Simplified lifecycle policy model consisting of one or more rules that determine which images in a repository should be expired. See https://docs.aws.amazon.com/AmazonECR/latest/userguide/lifecycle_policy_examples.html for more details.
type LifecyclePolicyArgs struct {
	// Specifies the rules to determine how images should be retired from this repository. Rules are ordered from lowest priority to highest. If there is a rule with a `selection` value of `any`, then it will have the highest priority.
	Rules LifecyclePolicyRuleArrayInput `pulumi:"rules"`
	// Skips creation of the policy if set to `true`.
	Skip *bool `pulumi:"skip"`
}

func (LifecyclePolicyArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LifecyclePolicy)(nil)).Elem()
}

func (i LifecyclePolicyArgs) ToLifecyclePolicyOutput() LifecyclePolicyOutput {
	return i.ToLifecyclePolicyOutputWithContext(context.Background())
}

func (i LifecyclePolicyArgs) ToLifecyclePolicyOutputWithContext(ctx context.Context) LifecyclePolicyOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LifecyclePolicyOutput)
}

func (i LifecyclePolicyArgs) ToLifecyclePolicyPtrOutput() LifecyclePolicyPtrOutput {
	return i.ToLifecyclePolicyPtrOutputWithContext(context.Background())
}

func (i LifecyclePolicyArgs) ToLifecyclePolicyPtrOutputWithContext(ctx context.Context) LifecyclePolicyPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LifecyclePolicyOutput).ToLifecyclePolicyPtrOutputWithContext(ctx)
}

// LifecyclePolicyPtrInput is an input type that accepts LifecyclePolicyArgs, LifecyclePolicyPtr and LifecyclePolicyPtrOutput values.
// You can construct a concrete instance of `LifecyclePolicyPtrInput` via:
//
//          LifecyclePolicyArgs{...}
//
//  or:
//
//          nil
type LifecyclePolicyPtrInput interface {
	pulumi.Input

	ToLifecyclePolicyPtrOutput() LifecyclePolicyPtrOutput
	ToLifecyclePolicyPtrOutputWithContext(context.Context) LifecyclePolicyPtrOutput
}

type lifecyclePolicyPtrType LifecyclePolicyArgs

func LifecyclePolicyPtr(v *LifecyclePolicyArgs) LifecyclePolicyPtrInput {
	return (*lifecyclePolicyPtrType)(v)
}

func (*lifecyclePolicyPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**LifecyclePolicy)(nil)).Elem()
}

func (i *lifecyclePolicyPtrType) ToLifecyclePolicyPtrOutput() LifecyclePolicyPtrOutput {
	return i.ToLifecyclePolicyPtrOutputWithContext(context.Background())
}

func (i *lifecyclePolicyPtrType) ToLifecyclePolicyPtrOutputWithContext(ctx context.Context) LifecyclePolicyPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LifecyclePolicyPtrOutput)
}

// Simplified lifecycle policy model consisting of one or more rules that determine which images in a repository should be expired. See https://docs.aws.amazon.com/AmazonECR/latest/userguide/lifecycle_policy_examples.html for more details.
type LifecyclePolicyOutput struct{ *pulumi.OutputState }

func (LifecyclePolicyOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LifecyclePolicy)(nil)).Elem()
}

func (o LifecyclePolicyOutput) ToLifecyclePolicyOutput() LifecyclePolicyOutput {
	return o
}

func (o LifecyclePolicyOutput) ToLifecyclePolicyOutputWithContext(ctx context.Context) LifecyclePolicyOutput {
	return o
}

func (o LifecyclePolicyOutput) ToLifecyclePolicyPtrOutput() LifecyclePolicyPtrOutput {
	return o.ToLifecyclePolicyPtrOutputWithContext(context.Background())
}

func (o LifecyclePolicyOutput) ToLifecyclePolicyPtrOutputWithContext(ctx context.Context) LifecyclePolicyPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v LifecyclePolicy) *LifecyclePolicy {
		return &v
	}).(LifecyclePolicyPtrOutput)
}

// Specifies the rules to determine how images should be retired from this repository. Rules are ordered from lowest priority to highest. If there is a rule with a `selection` value of `any`, then it will have the highest priority.
func (o LifecyclePolicyOutput) Rules() LifecyclePolicyRuleArrayOutput {
	return o.ApplyT(func(v LifecyclePolicy) []LifecyclePolicyRule { return v.Rules }).(LifecyclePolicyRuleArrayOutput)
}

// Skips creation of the policy if set to `true`.
func (o LifecyclePolicyOutput) Skip() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v LifecyclePolicy) *bool { return v.Skip }).(pulumi.BoolPtrOutput)
}

type LifecyclePolicyPtrOutput struct{ *pulumi.OutputState }

func (LifecyclePolicyPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**LifecyclePolicy)(nil)).Elem()
}

func (o LifecyclePolicyPtrOutput) ToLifecyclePolicyPtrOutput() LifecyclePolicyPtrOutput {
	return o
}

func (o LifecyclePolicyPtrOutput) ToLifecyclePolicyPtrOutputWithContext(ctx context.Context) LifecyclePolicyPtrOutput {
	return o
}

func (o LifecyclePolicyPtrOutput) Elem() LifecyclePolicyOutput {
	return o.ApplyT(func(v *LifecyclePolicy) LifecyclePolicy {
		if v != nil {
			return *v
		}
		var ret LifecyclePolicy
		return ret
	}).(LifecyclePolicyOutput)
}

// Specifies the rules to determine how images should be retired from this repository. Rules are ordered from lowest priority to highest. If there is a rule with a `selection` value of `any`, then it will have the highest priority.
func (o LifecyclePolicyPtrOutput) Rules() LifecyclePolicyRuleArrayOutput {
	return o.ApplyT(func(v *LifecyclePolicy) []LifecyclePolicyRule {
		if v == nil {
			return nil
		}
		return v.Rules
	}).(LifecyclePolicyRuleArrayOutput)
}

// Skips creation of the policy if set to `true`.
func (o LifecyclePolicyPtrOutput) Skip() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *LifecyclePolicy) *bool {
		if v == nil {
			return nil
		}
		return v.Skip
	}).(pulumi.BoolPtrOutput)
}

// A lifecycle policy rule that determine which images in a repository should be expired.
type LifecyclePolicyRule struct {
	// Describes the purpose of a rule within a lifecycle policy.
	Description *string `pulumi:"description"`
	// The maximum age limit (in days) for your images. Either [maximumNumberOfImages] or [maximumAgeLimit] must be provided.
	MaximumAgeLimit *float64 `pulumi:"maximumAgeLimit"`
	// The maximum number of images that you want to retain in your repository. Either [maximumNumberOfImages] or [maximumAgeLimit] must be provided.
	MaximumNumberOfImages *float64 `pulumi:"maximumNumberOfImages"`
	// A list of image tag prefixes on which to take action with your lifecycle policy. Only used if you specified "tagStatus": "tagged". For example, if your images are tagged as prod, prod1, prod2, and so on, you would use the tag prefix prod to specify all of them. If you specify multiple tags, only the images with all specified tags are selected.
	TagPrefixList []string `pulumi:"tagPrefixList"`
	// Determines whether the lifecycle policy rule that you are adding specifies a tag for an image. Acceptable options are tagged, untagged, or any. If you specify any, then all images have the rule evaluated against them. If you specify tagged, then you must also specify a tagPrefixList value. If you specify untagged, then you must omit tagPrefixList.
	TagStatus LifecycleTagStatus `pulumi:"tagStatus"`
}

// LifecyclePolicyRuleInput is an input type that accepts LifecyclePolicyRuleArgs and LifecyclePolicyRuleOutput values.
// You can construct a concrete instance of `LifecyclePolicyRuleInput` via:
//
//          LifecyclePolicyRuleArgs{...}
type LifecyclePolicyRuleInput interface {
	pulumi.Input

	ToLifecyclePolicyRuleOutput() LifecyclePolicyRuleOutput
	ToLifecyclePolicyRuleOutputWithContext(context.Context) LifecyclePolicyRuleOutput
}

// A lifecycle policy rule that determine which images in a repository should be expired.
type LifecyclePolicyRuleArgs struct {
	// Describes the purpose of a rule within a lifecycle policy.
	Description pulumi.StringPtrInput `pulumi:"description"`
	// The maximum age limit (in days) for your images. Either [maximumNumberOfImages] or [maximumAgeLimit] must be provided.
	MaximumAgeLimit pulumi.Float64PtrInput `pulumi:"maximumAgeLimit"`
	// The maximum number of images that you want to retain in your repository. Either [maximumNumberOfImages] or [maximumAgeLimit] must be provided.
	MaximumNumberOfImages pulumi.Float64PtrInput `pulumi:"maximumNumberOfImages"`
	// A list of image tag prefixes on which to take action with your lifecycle policy. Only used if you specified "tagStatus": "tagged". For example, if your images are tagged as prod, prod1, prod2, and so on, you would use the tag prefix prod to specify all of them. If you specify multiple tags, only the images with all specified tags are selected.
	TagPrefixList pulumi.StringArrayInput `pulumi:"tagPrefixList"`
	// Determines whether the lifecycle policy rule that you are adding specifies a tag for an image. Acceptable options are tagged, untagged, or any. If you specify any, then all images have the rule evaluated against them. If you specify tagged, then you must also specify a tagPrefixList value. If you specify untagged, then you must omit tagPrefixList.
	TagStatus LifecycleTagStatusInput `pulumi:"tagStatus"`
}

func (LifecyclePolicyRuleArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LifecyclePolicyRule)(nil)).Elem()
}

func (i LifecyclePolicyRuleArgs) ToLifecyclePolicyRuleOutput() LifecyclePolicyRuleOutput {
	return i.ToLifecyclePolicyRuleOutputWithContext(context.Background())
}

func (i LifecyclePolicyRuleArgs) ToLifecyclePolicyRuleOutputWithContext(ctx context.Context) LifecyclePolicyRuleOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LifecyclePolicyRuleOutput)
}

// LifecyclePolicyRuleArrayInput is an input type that accepts LifecyclePolicyRuleArray and LifecyclePolicyRuleArrayOutput values.
// You can construct a concrete instance of `LifecyclePolicyRuleArrayInput` via:
//
//          LifecyclePolicyRuleArray{ LifecyclePolicyRuleArgs{...} }
type LifecyclePolicyRuleArrayInput interface {
	pulumi.Input

	ToLifecyclePolicyRuleArrayOutput() LifecyclePolicyRuleArrayOutput
	ToLifecyclePolicyRuleArrayOutputWithContext(context.Context) LifecyclePolicyRuleArrayOutput
}

type LifecyclePolicyRuleArray []LifecyclePolicyRuleInput

func (LifecyclePolicyRuleArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]LifecyclePolicyRule)(nil)).Elem()
}

func (i LifecyclePolicyRuleArray) ToLifecyclePolicyRuleArrayOutput() LifecyclePolicyRuleArrayOutput {
	return i.ToLifecyclePolicyRuleArrayOutputWithContext(context.Background())
}

func (i LifecyclePolicyRuleArray) ToLifecyclePolicyRuleArrayOutputWithContext(ctx context.Context) LifecyclePolicyRuleArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LifecyclePolicyRuleArrayOutput)
}

// A lifecycle policy rule that determine which images in a repository should be expired.
type LifecyclePolicyRuleOutput struct{ *pulumi.OutputState }

func (LifecyclePolicyRuleOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LifecyclePolicyRule)(nil)).Elem()
}

func (o LifecyclePolicyRuleOutput) ToLifecyclePolicyRuleOutput() LifecyclePolicyRuleOutput {
	return o
}

func (o LifecyclePolicyRuleOutput) ToLifecyclePolicyRuleOutputWithContext(ctx context.Context) LifecyclePolicyRuleOutput {
	return o
}

// Describes the purpose of a rule within a lifecycle policy.
func (o LifecyclePolicyRuleOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LifecyclePolicyRule) *string { return v.Description }).(pulumi.StringPtrOutput)
}

// The maximum age limit (in days) for your images. Either [maximumNumberOfImages] or [maximumAgeLimit] must be provided.
func (o LifecyclePolicyRuleOutput) MaximumAgeLimit() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v LifecyclePolicyRule) *float64 { return v.MaximumAgeLimit }).(pulumi.Float64PtrOutput)
}

// The maximum number of images that you want to retain in your repository. Either [maximumNumberOfImages] or [maximumAgeLimit] must be provided.
func (o LifecyclePolicyRuleOutput) MaximumNumberOfImages() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v LifecyclePolicyRule) *float64 { return v.MaximumNumberOfImages }).(pulumi.Float64PtrOutput)
}

// A list of image tag prefixes on which to take action with your lifecycle policy. Only used if you specified "tagStatus": "tagged". For example, if your images are tagged as prod, prod1, prod2, and so on, you would use the tag prefix prod to specify all of them. If you specify multiple tags, only the images with all specified tags are selected.
func (o LifecyclePolicyRuleOutput) TagPrefixList() pulumi.StringArrayOutput {
	return o.ApplyT(func(v LifecyclePolicyRule) []string { return v.TagPrefixList }).(pulumi.StringArrayOutput)
}

// Determines whether the lifecycle policy rule that you are adding specifies a tag for an image. Acceptable options are tagged, untagged, or any. If you specify any, then all images have the rule evaluated against them. If you specify tagged, then you must also specify a tagPrefixList value. If you specify untagged, then you must omit tagPrefixList.
func (o LifecyclePolicyRuleOutput) TagStatus() LifecycleTagStatusOutput {
	return o.ApplyT(func(v LifecyclePolicyRule) LifecycleTagStatus { return v.TagStatus }).(LifecycleTagStatusOutput)
}

type LifecyclePolicyRuleArrayOutput struct{ *pulumi.OutputState }

func (LifecyclePolicyRuleArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]LifecyclePolicyRule)(nil)).Elem()
}

func (o LifecyclePolicyRuleArrayOutput) ToLifecyclePolicyRuleArrayOutput() LifecyclePolicyRuleArrayOutput {
	return o
}

func (o LifecyclePolicyRuleArrayOutput) ToLifecyclePolicyRuleArrayOutputWithContext(ctx context.Context) LifecyclePolicyRuleArrayOutput {
	return o
}

func (o LifecyclePolicyRuleArrayOutput) Index(i pulumi.IntInput) LifecyclePolicyRuleOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) LifecyclePolicyRule {
		return vs[0].([]LifecyclePolicyRule)[vs[1].(int)]
	}).(LifecyclePolicyRuleOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*LifecyclePolicyInput)(nil)).Elem(), LifecyclePolicyArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*LifecyclePolicyPtrInput)(nil)).Elem(), LifecyclePolicyArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*LifecyclePolicyRuleInput)(nil)).Elem(), LifecyclePolicyRuleArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*LifecyclePolicyRuleArrayInput)(nil)).Elem(), LifecyclePolicyRuleArray{})
	pulumi.RegisterOutputType(LifecyclePolicyOutput{})
	pulumi.RegisterOutputType(LifecyclePolicyPtrOutput{})
	pulumi.RegisterOutputType(LifecyclePolicyRuleOutput{})
	pulumi.RegisterOutputType(LifecyclePolicyRuleArrayOutput{})
}
