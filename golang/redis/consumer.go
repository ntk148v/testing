package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Package struct {

	// Activate time in ISO-8601 format
	// Example: 2020-06-19T17:36:46.271Z
	// Format: date-time
	ActivatedAt string `json:"activatedAt,omitempty"`

	// Expiration time in ISO-8601 format
	// Example: 2020-06-19T17:36:46.271Z
	// Required: true
	// Format: date-time
	ExpiredAt string `json:"expiredAt"`

	// A set of limitation
	// Required: true
	Metadata map[string]int64 `json:"metadata"`

	// SOC Service name
	// Required: true
	Service string `json:"service"`
}

type Event struct {
	TenantID string  `json:"tenant_id"`
	Package  Package `json:"package"`
}

func (p *Event) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Event) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

const nsmQueue = "nsm"

func main() {
	// Create new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0, // use default DB
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		client.Close()
	}()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	for {
		result, err := client.BLPop(ctx, 10*time.Second, nsmQueue).Result()
		if err != nil && err != redis.Nil {
			fmt.Println("Pop nsm package error:", err)
			continue
		}

		// Array reply
		// A nil multi-bulk when no element could be popped and the timeout expired.
		// A two-element multi-bulk with the first element being the name of the key
		// where an element was popped and the second element being the value of the popped element.
		if result == nil {
			continue
		}
		eventBytes := []byte(result[1])
		var event Event
		event.UnmarshalBinary(eventBytes)
		fmt.Printf("Pop nsm event %+v\n", event)
	}
}
