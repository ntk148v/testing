package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rvflash/tcp"
)

func main() {
	bye := make(chan os.Signal, 1)
	signal.Notify(bye, os.Interrupt, syscall.SIGTERM)

	// Creates a server with a logger and a recover on panic as middlewares.
	r := tcp.Default()
	r.ACK(func(c *tcp.Context) {
		// New message received
		// Gets the request body
		buf, err := c.ReadAll()
		if err != nil {
			c.Error(err)
			return
		}
		// Writes something as response
		c.String(string(buf))
	})
	go func() {
		err := r.Run(":9090")
		if err != nil {
			log.Printf("server: %q\n", err)
		}
	}()

	<-bye
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	err := r.Shutdown(ctx)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}
