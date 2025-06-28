package main

import (
	"log"

	"go-micro.dev/v5"
	"go-micro.dev/v5/broker"
)

func main() {
	// Initialize the service
	service := micro.NewService(micro.Name("example.subscriber"))
	service.Init()

	// Start the broker
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker connect error: %v", err)
	}

	// Subscribe to messages
	_, err := broker.Subscribe("example.topic", func(p broker.Event) error {
		log.Printf("Received message: %s\n", string(p.Message().Body))
		return nil
	})
	if err != nil {
		log.Fatalf("Error subscribing: %v", err)
	}

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
