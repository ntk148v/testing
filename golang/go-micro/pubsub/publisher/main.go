package main

import (
	"log"
	"time"

	"go-micro.dev/v5"
	"go-micro.dev/v5/broker"
)

func main() {
	// Initialize the service
	service := micro.NewService(micro.Name("example.publisher"))
	service.Init()

	// Start the broker
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker connect error: %v", err)
	}

	// Publish a message every 5 seconds
	go func() {
		t := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-t.C:
				msg := &broker.Message{
					Header: map[string]string{"id": "1"},
					Body:   []byte("Hello from the Publisher!"),
				}
				if err := broker.Publish("example.topic", msg); err != nil {
					log.Printf("Error publishing: %v", err)
				}
			}
		}
	}()

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
