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

const listName string = "queue:test"

func randName() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

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
		name := randName()
		if _, err := client.RPush(ctx, listName, &Event{Name: name}).Result(); err != nil {
			fmt.Printf("Push event error: %s\n", err)
			continue
		}
		fmt.Printf("Push event %s\n", name)
		time.Sleep(2 * time.Second)
	}
}
