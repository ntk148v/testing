package main

import (
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Pull nginx image
		image, err := docker.NewImage(ctx, "nginxImage", &docker.ImageArgs{
			ImageName: pulumi.String("nginx:latest"),
		})
		if err != nil {
			return err
		}

		// Create nginx container
		_, err = docker.NewContainer(ctx, "nginxContainer", &docker.ContainerArgs{
			Image: image.RepoDigest,
			Name:  pulumi.String("nginx_example"),
			Ports: docker.ContainerPortArray{
				&docker.ContainerPortArgs{
					Internal: pulumi.Int(80),
					External: pulumi.Int(8080),
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
