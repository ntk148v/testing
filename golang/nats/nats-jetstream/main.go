package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"time"

	"github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
)

// TestMessage is a message that can help test timings on jetstream
type TestMessage struct {
	ID          int       `json:"id"`
	PublishTime time.Time `json:"publish_time"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	stream := uuid.NewV4().String()
	subject := stream

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("unable to connect to nats: %v", err)
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatalf("error getting jetstream: %v", err)
	}

	info, err := js.StreamInfo(stream)
	if err == nil {
		log.Fatalf("stream already exists: %v", info)
	}

	// Run pub/sub for 1 minute
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      stream,
		Subjects:  []string{subject},
		Retention: nats.WorkQueuePolicy,
	}, nats.Context(ctx))
	if err != nil {
		log.Fatalf("can't add stream: %v", err)
	}

	results := make(chan int64)

	var (
		totalTime     int64
		totalMessages int64
	)

	// Create 2 consumers
	for i := 0; i < 2; i++ {
		go sub(ctx, subject, results)
	}

	// Create a publisher
	go pub(ctx, subject, js)

	for {
		select {
		case <-ctx.Done():
			cancel()
			log.Printf("sent %d messages with average time of %f", totalMessages,
				math.Round(float64(totalTime/totalMessages)))
			js.DeleteStream(stream)
			return
		case microsec := <-results:
			totalTime += microsec
			totalMessages++
		}

	}
}

func pub(ctx context.Context, subject string, js nats.JetStream) {
	i := 0

	for {
		start := time.Now()
		bytes, err := json.Marshal(&TestMessage{
			ID:          i,
			PublishTime: start,
		})
		if err != nil {
			log.Fatalf("could not get bytes from literal TestMessage: %v", err)
		}

		_, err = js.Publish(subject, bytes, nats.Context(ctx))
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return
			}

			log.Fatalf("error publishing: %v", err)
		}

		log.Printf("[publisher] sent %d, publish time microsec: %d", i, time.Since(start).Microseconds())
		time.Sleep(1 * time.Second)
		i++
	}
}

func sub(ctx context.Context, subject string, results chan int64) {
	id := uuid.NewV4().String()

	nc, err := nats.Connect(nats.DefaultURL, nats.Name(id))
	if err != nil {
		log.Fatalf("[consumer: %s] unable to connect to nats: %v", id, err)
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("[consumer: %s] error getting jetstream: %v", id, err)
	}

	sub, err := js.PullSubscribe(subject, "group")
	if err != nil {
		log.Fatalf("[consumer: %s] error pulling subscribe")
	}

	for {
		msgs, err := sub.Fetch(1, nats.Context(ctx))
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}

			log.Printf("[consumer: %s] error consuming, sleeping for a second: %v", id, err)
			time.Sleep(1 * time.Second)
			continue
		}

		msg := msgs[0]

		var tMsg *TestMessage

		err = json.Unmarshal(msg.Data, &tMsg)
		if err != nil {
			log.Printf("[consumer: %s] error unmarshaling, sleeping for a second: %v", id, err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Get publish time for logging
		tm := time.Since(tMsg.PublishTime).Microseconds()
		results <- tm
		log.Printf("[consumer: %s] received msg (%d) after waiting microsec: %d", id, tMsg.ID, tm)

		err = msg.Ack(nats.Context(ctx))
		if err != nil {
			log.Printf("[consumer: %s] error acking message: %v", id, err)
		}
	}
}
