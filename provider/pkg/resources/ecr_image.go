package resources

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	do "github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/zchase/pulumi-awsx-go/pkg/utils"
)

const ImageIdentifier = "awsx-go:ecr:Image"

type ImageArgs struct {
	Args          map[string]string   `pulumi:"args"`
	CacheFrom     []string            `pulumi:"cacheFrom"`
	DockerFile    string              `pulumi:"dockerfile"`
	Env           map[string]string   `pulumi:"env"`
	ExtraOptions  []string            `pulumi:"extraOptions"`
	Path          string              `pulumi:"path"`
	RepositoryURL pulumi.StringOutput `pulumi:"repositoryUrl"`
	Target        string              `pulumi:"target"`
}

type Image struct {
	pulumi.ResourceState

	ImageURI pulumi.StringOutput `pulumi:"imageUri"`
}

func NewImage(ctx *pulumi.Context, name string, args *ImageArgs, opts ...pulumi.ResourceOption) (*Image, error) {
	if args == nil {
		args = &ImageArgs{}
	}

	component := &Image{}
	err := ctx.RegisterComponentResource(ImageIdentifier, name, component, opts...)
	if err != nil {
		return nil, err
	}

	imageURI, err := computeImageFromAsset(ctx, name, args, args.RepositoryURL, component)
	if err != nil {
		return nil, err
	}

	component.ImageURI = imageURI

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"imageUri": imageURI,
	}); err != nil {
		return nil, err
	}

	return component, nil
}

type imageRepositoryCreds struct {
	Registry pulumi.StringOutput
	Username pulumi.StringOutput
	Password pulumi.StringOutput
}

func computeImageFromAsset(ctx *pulumi.Context, name string, args *ImageArgs, repositoryURL pulumi.StringOutput, parent pulumi.Resource) (pulumi.StringOutput, error) {
	imageName := getImageName(ctx, args)

	dockerBuild := do.DockerBuildArgs{
		Args:         pulumi.ToStringMap(args.Args),
		Context:      pulumi.String(args.Path),
		Dockerfile:   pulumi.String(args.DockerFile),
		Env:          pulumi.ToStringMap(args.Env),
		ExtraOptions: pulumi.ToStringArray(args.ExtraOptions),
		Target:       pulumi.String(args.Target),
	}

	if len(args.CacheFrom) > 0 {
		dockerBuild.CacheFrom = do.CacheFromPtr(&do.CacheFromArgs{
			Stages: pulumi.ToStringArray(args.CacheFrom),
		})
	}

	repoCreds := repositoryURL.ApplyT(func(repoURL string) (ecr.GetCredentialsResultOutput, error) {
		if repoURL == "" {
			return ecr.GetCredentialsResultOutput{}, fmt.Errorf("RepositoryURL is empty")
		}

		parsedURL, err := url.Parse(fmt.Sprintf("https://%s", repoURL))
		if err != nil {
			return ecr.GetCredentialsResultOutput{}, err
		}

		repoID := strings.Split(parsedURL.Hostname(), ".")[0]

		return ecr.GetCredentialsOutput(ctx, ecr.GetCredentialsOutputArgs{
			RegistryId: pulumi.String(repoID),
		}, pulumi.Parent(parent)), nil
	}).(ecr.GetCredentialsResultOutput)

	imageRepoCreds := repoCreds.ApplyT(func(creds ecr.GetCredentialsResult) (imageRepositoryCreds, error) {
		data, err := base64.StdEncoding.DecodeString(creds.AuthorizationToken)
		if err != nil {
			return imageRepositoryCreds{}, err
		}

		parts := strings.Split(string(data), ":")
		if len(parts) != 2 {
			return imageRepositoryCreds{}, fmt.Errorf("Invalid credentials")
		}

		return imageRepositoryCreds{
			Registry: pulumi.String(creds.ProxyEndpoint).ToStringOutput(),
			Username: pulumi.String(parts[0]).ToStringOutput(),
			Password: pulumi.String(parts[1]).ToStringOutput(),
		}, nil
	}).(pulumi.AnyOutput)

	uniqueImageName, err := do.NewImage(ctx, name, &do.ImageArgs{
		ImageName:      repositoryURL,
		LocalImageName: pulumi.String(imageName),
		Build:          dockerBuild,
		SkipPush:       pulumi.Bool(false),
		Registry: do.ImageRegistryArgs{
			Server: imageRepoCreds.ApplyT(func(x interface{}) pulumi.StringOutput {
				c := x.(imageRepositoryCreds)
				return c.Registry
			}).(pulumi.StringOutput),
			Username: imageRepoCreds.ApplyT(func(x interface{}) pulumi.StringOutput {
				c := x.(imageRepositoryCreds)
				return c.Username
			}).(pulumi.StringOutput),
			Password: imageRepoCreds.ApplyT(func(x interface{}) pulumi.StringOutput {
				c := x.(imageRepositoryCreds)
				return c.Password
			}).(pulumi.StringOutput),
		},
	}, pulumi.Parent(parent))
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	return uniqueImageName.ImageName, nil
}

func getImageName(ctx *pulumi.Context, inputs *ImageArgs) string {
	buildSig := "."
	if inputs.Path != "" {
		buildSig = inputs.Path
	}

	if inputs.DockerFile != "" {
		buildSig += fmt.Sprintf(";dockerfile=%s", inputs.DockerFile)
	}

	for k, v := range inputs.Args {
		buildSig += fmt.Sprintf("`;arg[%s]=%s`", k, v)
	}

	buildSig += ctx.Stack()
	return strings.ToLower(fmt.Sprintf("%s-container", utils.SHA1Hash(buildSig)))
}
