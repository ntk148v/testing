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
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	streamName    string = "test"
	consumerGroup string = "test"
	consumerName  string = "test"
)

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

	// Create consumer group
	if _, err := client.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Result(); err != nil {
		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			panic(err)
		}
	}

	go consumerEvents(ctx, client)
	consumeAutoClaim(ctx, client)
}

func consumerEvents(ctx context.Context, client *redis.Client) {
	for {
		streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Streams:  []string{streamName, ">"},
			Group:    consumerGroup,
			Consumer: consumerName,
		}).Result()

		if err != nil {
			fmt.Printf("Consume event error: %s\n", err)
			continue
		}

		for _, message := range streams[0].Messages {
			for _, value := range message.Values {
				var event Event
				event.UnmarshalBinary([]byte(value.(string)))
				fmt.Printf("Consume event ID %s: %+v\n", message.ID, event)
			}
			// Consumer gets message, but not ACK -> message -> [PENDING LIST]
			// Use XAUTOCLAIM to scan PENDING list then claim
			// https://redis.io/commands/xautoclaim
			// if _, err = client.XAck(ctx, streamName, consumerGroup, message.ID).Result(); err != nil {
			// 	fmt.Printf("ACK event ID %s error: %s\n", message.ID, err)
			// }
		}
	}
}

func consumeAutoClaim(ctx context.Context, client *redis.Client) {
	ticker := time.Tick(time.Second * 10)
	for {
		select {
		case <-ticker:
			messages, _, err := client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
				Stream:   streamName,
				Group:    consumerGroup,
				Consumer: consumerName,
				MinIdle:  10 * time.Second,
				Start:    "0",
			}).Result()

			if err != nil {
				panic(err)
			}

			for _, message := range messages {
				for _, value := range message.Values {
					var event Event
					event.UnmarshalBinary([]byte(value.(string)))
					fmt.Printf("Consume pending event ID %s: %+v\n", message.ID, event)
				}
				if _, err = client.XAck(ctx, streamName, consumerGroup, message.ID).Result(); err != nil {
					fmt.Printf("ACK event ID %s error: %s\n", message.ID, err)
				}
			}
		}
	}
}
