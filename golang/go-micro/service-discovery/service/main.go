package main

import (
	"context"
	"log"

	"go-micro.dev/v5"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *interface{}, resp *interface{}) error {
	log.Println("Service A was called")
	return nil
}

func main() {
	// Create a new Service
	service := micro.NewService(
		micro.Name("serviceA"),
	)

	// Initialize the service
	service.Init()

	// Register handler
	micro.RegisterHandler(service.Server(), new(Greeter))

	// Start the service
	if err := service.Run(); err != nil {
		log.Fatalf("Error running service A: %v", err)
	}
}
