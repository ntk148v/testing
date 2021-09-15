package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

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

func randName() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
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

	// Generate event
	for {
		name := randName()
		id, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: os.Getenv("REDIS_STREAM"),
			Values: map[string]interface{}{
				name: &Event{
					Name: name,
				}},
		}).Result()
		if err != nil {
			fmt.Printf("Produce event with name %s error: %s\n", name, err)
			continue
		}
		fmt.Printf("Produce event with name %s id: %s\n", name, id)
		time.Sleep(2 * time.Second)
	}
}
