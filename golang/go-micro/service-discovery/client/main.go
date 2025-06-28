package main

import (
	"context"
	"log"

	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
)

func main() {
	// create a new service
	service := micro.NewService(micro.Name("serviceB"))
	service.Init()

	// wequest message
	req := service.Client().NewRequest("serviceA", "Greeter.Hello", &map[string]interface{}{}, client.WithContentType("application/json"))
	rsp := &map[string]interface{}{}

	// call the service
	if err := service.Client().Call(context.Background(), req, rsp); err != nil {
		log.Fatalf("Error calling service A: %v", err)
	}

	log.Println("Successfully called service A")
}
