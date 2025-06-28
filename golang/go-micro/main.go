package main

import (
	"fmt"
	"log"
	"net/http"

	"go-micro.dev/v5/web"
)

func main() {
	// Create a new service
	service := web.NewService(
		web.Name("helloworld"),
		web.Address(":8080"),
	)

	// Initialize the service
	service.Init()

	// Set up a route and handler
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	// Assign the handler to the service
	service.Handle("/", http.DefaultServeMux)

	// Start the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
