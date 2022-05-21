# coding=utf-8
# *** WARNING: this file was generated by Pulumi SDK Generator. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from .. import _utilities

__all__ = ['ImageArgs', 'Image']

@pulumi.input_type
class ImageArgs:
    def __init__(__self__, *,
                 repository_url: pulumi.Input[str],
                 args: Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]] = None,
                 cache_from: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 dockerfile: Optional[pulumi.Input[str]] = None,
                 env: Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]] = None,
                 extra_options: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 path: Optional[pulumi.Input[str]] = None,
                 target: Optional[pulumi.Input[str]] = None):
        """
        The set of arguments for constructing a Image resource.
        :param pulumi.Input[str] repository_url: Url of the repository
        :param pulumi.Input[Mapping[str, pulumi.Input[str]]] args: An optional map of named build-time argument variables to set during the Docker build.  This flag allows you to pass built-time variables that can be accessed like environment variables inside the `RUN` instruction.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] cache_from: Images to consider as cache sources
        :param pulumi.Input[str] dockerfile: dockerfile may be used to override the default Dockerfile name and/or location.  By default, it is assumed to be a file named Dockerfile in the root of the build context.
        :param pulumi.Input[Mapping[str, pulumi.Input[str]]] env: Environment variables to set on the invocation of `docker build`, for example to support `DOCKER_BUILDKIT=1 docker build`.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] extra_options: An optional catch-all list of arguments to provide extra CLI options to the docker build command.  For example `['--network', 'host']`.
        :param pulumi.Input[str] path: Path to a directory to use for the Docker build context, usually the directory in which the Dockerfile resides (although dockerfile may be used to choose a custom location independent of this choice). If not specified, the context defaults to the current working directory; if a relative path is used, it is relative to the current working directory that Pulumi is evaluating.
        :param pulumi.Input[str] target: The target of the dockerfile to build
        """
        pulumi.set(__self__, "repository_url", repository_url)
        if args is not None:
            pulumi.set(__self__, "args", args)
        if cache_from is not None:
            pulumi.set(__self__, "cache_from", cache_from)
        if dockerfile is not None:
            pulumi.set(__self__, "dockerfile", dockerfile)
        if env is not None:
            pulumi.set(__self__, "env", env)
        if extra_options is not None:
            pulumi.set(__self__, "extra_options", extra_options)
        if path is not None:
            pulumi.set(__self__, "path", path)
        if target is not None:
            pulumi.set(__self__, "target", target)

    @property
    @pulumi.getter(name="repositoryUrl")
    def repository_url(self) -> pulumi.Input[str]:
        """
        Url of the repository
        """
        return pulumi.get(self, "repository_url")

    @repository_url.setter
    def repository_url(self, value: pulumi.Input[str]):
        pulumi.set(self, "repository_url", value)

    @property
    @pulumi.getter
    def args(self) -> Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]]:
        """
        An optional map of named build-time argument variables to set during the Docker build.  This flag allows you to pass built-time variables that can be accessed like environment variables inside the `RUN` instruction.
        """
        return pulumi.get(self, "args")

    @args.setter
    def args(self, value: Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]]):
        pulumi.set(self, "args", value)

    @property
    @pulumi.getter(name="cacheFrom")
    def cache_from(self) -> Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]:
        """
        Images to consider as cache sources
        """
        return pulumi.get(self, "cache_from")

    @cache_from.setter
    def cache_from(self, value: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]):
        pulumi.set(self, "cache_from", value)

    @property
    @pulumi.getter
    def dockerfile(self) -> Optional[pulumi.Input[str]]:
        """
        dockerfile may be used to override the default Dockerfile name and/or location.  By default, it is assumed to be a file named Dockerfile in the root of the build context.
        """
        return pulumi.get(self, "dockerfile")

    @dockerfile.setter
    def dockerfile(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "dockerfile", value)

    @property
    @pulumi.getter
    def env(self) -> Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]]:
        """
        Environment variables to set on the invocation of `docker build`, for example to support `DOCKER_BUILDKIT=1 docker build`.
        """
        return pulumi.get(self, "env")

    @env.setter
    def env(self, value: Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]]):
        pulumi.set(self, "env", value)

    @property
    @pulumi.getter(name="extraOptions")
    def extra_options(self) -> Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]:
        """
        An optional catch-all list of arguments to provide extra CLI options to the docker build command.  For example `['--network', 'host']`.
        """
        return pulumi.get(self, "extra_options")

    @extra_options.setter
    def extra_options(self, value: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]]):
        pulumi.set(self, "extra_options", value)

    @property
    @pulumi.getter
    def path(self) -> Optional[pulumi.Input[str]]:
        """
        Path to a directory to use for the Docker build context, usually the directory in which the Dockerfile resides (although dockerfile may be used to choose a custom location independent of this choice). If not specified, the context defaults to the current working directory; if a relative path is used, it is relative to the current working directory that Pulumi is evaluating.
        """
        return pulumi.get(self, "path")

    @path.setter
    def path(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "path", value)

    @property
    @pulumi.getter
    def target(self) -> Optional[pulumi.Input[str]]:
        """
        The target of the dockerfile to build
        """
        return pulumi.get(self, "target")

    @target.setter
    def target(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "target", value)


class Image(pulumi.ComponentResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 args: Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]] = None,
                 cache_from: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 dockerfile: Optional[pulumi.Input[str]] = None,
                 env: Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]] = None,
                 extra_options: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 path: Optional[pulumi.Input[str]] = None,
                 repository_url: Optional[pulumi.Input[str]] = None,
                 target: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        """
        Builds a docker image and pushes to the ECR repository

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[Mapping[str, pulumi.Input[str]]] args: An optional map of named build-time argument variables to set during the Docker build.  This flag allows you to pass built-time variables that can be accessed like environment variables inside the `RUN` instruction.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] cache_from: Images to consider as cache sources
        :param pulumi.Input[str] dockerfile: dockerfile may be used to override the default Dockerfile name and/or location.  By default, it is assumed to be a file named Dockerfile in the root of the build context.
        :param pulumi.Input[Mapping[str, pulumi.Input[str]]] env: Environment variables to set on the invocation of `docker build`, for example to support `DOCKER_BUILDKIT=1 docker build`.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] extra_options: An optional catch-all list of arguments to provide extra CLI options to the docker build command.  For example `['--network', 'host']`.
        :param pulumi.Input[str] path: Path to a directory to use for the Docker build context, usually the directory in which the Dockerfile resides (although dockerfile may be used to choose a custom location independent of this choice). If not specified, the context defaults to the current working directory; if a relative path is used, it is relative to the current working directory that Pulumi is evaluating.
        :param pulumi.Input[str] repository_url: Url of the repository
        :param pulumi.Input[str] target: The target of the dockerfile to build
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: ImageArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Builds a docker image and pushes to the ECR repository

        :param str resource_name: The name of the resource.
        :param ImageArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ImageArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 args: Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]] = None,
                 cache_from: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 dockerfile: Optional[pulumi.Input[str]] = None,
                 env: Optional[pulumi.Input[Mapping[str, pulumi.Input[str]]]] = None,
                 extra_options: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 path: Optional[pulumi.Input[str]] = None,
                 repository_url: Optional[pulumi.Input[str]] = None,
                 target: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        if opts is None:
            opts = pulumi.ResourceOptions()
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.version is None:
            opts.version = _utilities.get_version()
        if opts.id is not None:
            raise ValueError('ComponentResource classes do not support opts.id')
        else:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ImageArgs.__new__(ImageArgs)

            __props__.__dict__["args"] = args
            __props__.__dict__["cache_from"] = cache_from
            __props__.__dict__["dockerfile"] = dockerfile
            __props__.__dict__["env"] = env
            __props__.__dict__["extra_options"] = extra_options
            __props__.__dict__["path"] = path
            if repository_url is None and not opts.urn:
                raise TypeError("Missing required property 'repository_url'")
            __props__.__dict__["repository_url"] = repository_url
            __props__.__dict__["target"] = target
            __props__.__dict__["image_uri"] = None
        super(Image, __self__).__init__(
            'awsx-go:ecr:Image',
            resource_name,
            __props__,
            opts,
            remote=True)

    @property
    @pulumi.getter(name="imageUri")
    def image_uri(self) -> pulumi.Output[str]:
        """
        Unique identifier of the pushed image
        """
        return pulumi.get(self, "image_uri")

