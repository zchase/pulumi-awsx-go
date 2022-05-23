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
	Args          map[string]string `pulumi:"args"`
	CacheFrom     []string          `pulumi:"cacheFrom"`
	DockerFile    string            `pulumi:"dockerfile"`
	Env           map[string]string `pulumi:"env"`
	ExtraOptions  []string          `pulumi:"extraOptions"`
	Path          string            `pulumi:"path"`
	RepositoryURL string            `pulumi:"repositoryUrl"`
	Target        string            `pulumi:"target"`
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

	imageURI, err := computeImageFromAsset(ctx, name, args, component)
	if err != nil {
		return nil, err
	}

	component.ImageURI = imageURI

	return component, nil
}

type imageRepositoryCreds struct {
	pulumi.Output

	Registry pulumi.StringOutput
	Username pulumi.StringOutput
	Password pulumi.StringOutput
}

func computeImageFromAsset(ctx *pulumi.Context, name string, args *ImageArgs, parent pulumi.Resource) (pulumi.StringOutput, error) {
	repoURL, err := url.Parse(fmt.Sprintf("https://%s", args.RepositoryURL))
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	repositoryID := strings.Split(repoURL.Hostname(), ",")[0]
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

	//var imageRepoCreds imageRepositoryCreds
	if repositoryID == "" {
		return pulumi.StringOutput{}, fmt.Errorf("Expected registry ID to be defined during push")
	}

	credentials := ecr.GetCredentialsOutput(ctx, ecr.GetCredentialsOutputArgs{
		RegistryId: pulumi.String(repositoryID),
	}, pulumi.Parent(parent))

	imageRepoCreds := credentials.AuthorizationToken().ApplyT(func(authorizationToken string) (imageRepositoryCreds, error) {
		data, err := base64.StdEncoding.DecodeString(authorizationToken)
		if err != nil {
			return imageRepositoryCreds{}, err
		}

		parts := strings.Split(string(data), ":")
		if len(parts) != 2 {
			return imageRepositoryCreds{}, fmt.Errorf("Invalid credentials")
		}

		return imageRepositoryCreds{
			Registry: credentials.ProxyEndpoint(),
			Username: pulumi.String(parts[0]).ToStringOutput(),
			Password: pulumi.String(parts[1]).ToStringOutput(),
		}, nil
	}).(imageRepositoryCreds)

	uniqueImageName, err := do.NewImage(ctx, name, &do.ImageArgs{
		ImageName: pulumi.String(imageName),
		Build:     dockerBuild,
		SkipPush:  pulumi.Bool(false),
		Registry: do.ImageRegistryArgs{
			Server:   imageRepoCreds.Registry,
			Username: imageRepoCreds.Username,
			Password: imageRepoCreds.Password,
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
	return fmt.Sprintf("%s-container", utils.SHA1Hash(buildSig))
}
