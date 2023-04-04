package main

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// pull an image
	image, err := client.Pull(ctx, "docker.io/library/nginx:latest")
	if err != nil {
		panic(err)
	}

	fmt.Println(imag)
}
