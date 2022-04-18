// Copyright 2022 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		// Blocking get the arrived event
		// streams, err := client.XRead(ctx, &redis.XReadArgs{
		// Streams: []string{streamName, "$"},
		// Block:   0,
		// }).Result()

		// Get from the 1st -> last
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
