// A simple service:
// - Launch an HTTP server with, at least, a handler registered at endpoint GET /health. This handler should return status code OK ( 200 ).
// - Register the service on Consul.
// - De-register the service on shutdown. Only if I’m running one instance of my service.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
)

// startHTTPServer starts a simple HTTP server with a health check endpoint
func startHTTPServer(port int) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	log.Printf("Starting HTTP server on :%d...\n", port)
	// launch server and passing nil as
	// mux value to use default multiplexer (mux)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		// exit program because
		// server will be launched as
		// goroutine.
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// registerServiceWithConsul registers service with Consul
func registerServiceWithConsul(serviceID, serviceName, serviceAddress string, servicePort int, consulClient *api.Client) error {
	// define the service registration
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceAddress,
		Port:    servicePort,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", serviceAddress, servicePort),
			Interval: "10s",
			Timeout:  "5s",
		},
	}
	// Register the service
	if err := consulClient.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	return nil
}

func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func main() {
	// Service configuration
	serviceID := "example-service-1"
	serviceName := "example-service"
	servicePort := 8080
	serviceAddress := getLocalIP().String()
	consulAddress := "http://localhost:8500"

	config := &api.Config{
		Address: consulAddress, // Consul address
	}

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	// Start a simple HTTP server
	go startHTTPServer(servicePort)

	// Register the service with Consul
	err = registerServiceWithConsul(
		serviceID, serviceName, serviceAddress, servicePort, client)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	fmt.Printf("Service %s is running on %s:%d and registered with Consul.\n", serviceName, serviceAddress, servicePort)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigs

	// Deregister the service upon shutdown
	err = client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		log.Fatalf("Failed to deregister service: %v", err)
	}
	fmt.Println("Service deregistered from Consul")
}
