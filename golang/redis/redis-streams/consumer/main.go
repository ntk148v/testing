package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var wg sync.WaitGroup

const streamName string = "test"

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
		wg.Wait()
		cancel()
		client.Close()
	}()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	for {
		streams, err := client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{streamName, "0"},
		}).Result()

		if err != nil {
			fmt.Printf("Consume event error: %s\n", err)
			continue
		}

		for _, stream := range streams[0].Messages {
			for _, value := range stream.Values {
				var event Event
				event.UnmarshalBinary([]byte(value.(string)))
				fmt.Printf("Consume event ID %s: %+v\n", stream.ID, event)
			}
			if _, err = client.XDel(ctx, streamName, stream.ID).Result(); err != nil {
				fmt.Printf("Delete event ID %s error: %s\n", stream.ID, err)
			}
		}
	}
}
