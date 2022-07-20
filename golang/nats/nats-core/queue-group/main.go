package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
)

// TestMessage is a message
type TestMessage struct {
	ID int `json:"id"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	subject := uuid.NewV4().String()
	// Run pub/sub for 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Create 2 consumers
	for i := 0; i < 2; i++ {
		go sub(ctx, subject)
	}

	// Create a publisher
	go pub(ctx, subject)

	select {
	case <-ctx.Done():
		cancel()
		return
	}
}

func pub(ctx context.Context, subject string) {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL, nats.Name("Publish"))
	if err != nil {
		log.Fatalf("[publisher] unable to connect to nats: %v", err)
	}
	defer nc.Close()

	i := 0
	for {
		select {
		case <-ctx.Done():
			nc.Drain()
			return
		default:
			bytes, err := json.Marshal(&TestMessage{ID: i})
			if err != nil {
				log.Fatalf("could not get bytes from literal TestMessage: %v", err)
			}

			if err := nc.Publish(subject, bytes); err != nil {
				log.Fatalf("[publisher] error publishing: %v", err)
			}

			log.Printf("[publisher] sent %d", i)
			time.Sleep(time.Second)
			i++
		}
	}
}

func sub(ctx context.Context, subject string) {
	id := uuid.NewV4().String()

	nc, err := nats.Connect(nats.DefaultURL, nats.Name(id))
	if err != nil {
		log.Fatalf("[consumer: %s] unable to connect to nats: %v", id, err)
	}
	defer nc.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := nc.QueueSubscribe(subject, "worker", func(msg *nats.Msg) {
				log.Printf("[consumer: %s] got message: %s", id, msg.Data)
			}); err != nil {
				log.Fatalf("[consumer: %s] error consuming message: %v", id, err)
			}
		}
	}
}
