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
