package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

const listName string = "queue:test"

type Event struct {
	Name string `json:"name"`
}

func (e *Event) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

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
		result, err := client.BLPop(ctx, 10*time.Second, listName).Result()
		if err != nil {
			fmt.Printf("Pop event error: %s\n", err)
			continue
		}
		// Array reply
		// A nil multi-bulk when no element could be popped and the timeout expired.
		// A two-element multi-bulk with the first element being the name of the key where an element was popped and the second element being the value of the popped element.
		if result == nil {
			continue
		}
		eventBytes := []byte(result[1])
		event := Event{}
		json.Unmarshal(eventBytes, &event)
		fmt.Printf("Pop event %s\n", event.Name)
	}
}
