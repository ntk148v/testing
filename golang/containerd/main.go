package main

import (
	"context"
	"fmt"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("default"))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	alpineImage := "docker.io/library/alpine:latest"
	var image containerd.Image

	// list all images
	images, err := client.ListImages(ctx)
	if err != nil {
		panic(err)
	}

	var exist bool
	for _, i := range images {
		if alpineImage == i.Name() {
			exist = true
			image = i
			break
		}
	}

	if !exist {
		// pull an image

		image, err = client.Pull(ctx, alpineImage, containerd.WithPullUnpack)
		if err != nil {
			panic(err)
		}

		fmt.Println("Pulled image:", image.Name())
	}

	// start a container
	alpine, err := client.NewContainer(ctx, "alpine",
		containerd.WithNewSnapshot("alpine-rootfs", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)))
	if err != nil {
		panic(err)
	}

	fmt.Println("Start container:", alpine.ID())
	defer alpine.Delete(ctx)

	// create a new task
	task, err := alpine.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		panic(err)
	}
	defer task.Delete(ctx)

	// the task is now running and has a pid that can be used to setup networking
	// or other runtime settings outside of containerd
	fmt.Printf("Start task with id %d inside container\n", task.Pid())

	// start the process inside the container
	if err = task.Start(ctx); err != nil {
		panic(err)
	}

	// wait for the task to exit and get the exit status
	// status, err := task.Wait(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// <-status
	task.Kill(ctx, syscall.SIGKILL)
}
