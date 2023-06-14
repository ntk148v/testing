package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func set(ctx context.Context, cli *redis.ClusterClient, key, value string) error {
	return cli.Set(ctx, key, value, 0).Err()
}

func get(ctx context.Context, cli *redis.ClusterClient, key string) (string, error) {
	return cli.Get(ctx, key).Result()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"redis-node-0:6379",
			"redis-node-1:6379",
			"redis-node-2:6379",
			"redis-node-3:6379",
			"redis-node-4:6379",
			"redis-node-5:6379",
			"redis-node-6:6379",
			"redis-node-7:6379",
			"redis-node-8:6379",
		},
		Password: "bitnami", // hardcode!
	})

	for i := 0; i < 1000000000; i++ {
		if err := set(ctx, client, fmt.Sprintf("foo%d", i), strconv.Itoa(i)); err != nil {
			fmt.Println("Something went wrong:", err)
			continue
		}
		if val, err := get(ctx, client, fmt.Sprintf("foo%d", i)); err != nil {
			fmt.Println("Something went wrong:", err)
			continue
		} else {
			fmt.Println(val)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
