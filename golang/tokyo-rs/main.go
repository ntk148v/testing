package main

import (
	"flag"
	"log"
	"time"

	tokyo "github.com/ledongthuc/tokyo_go_sdk"
)

var server = flag.String("server", "wss://combat.sege.dev/socket", "server host")
var userKey = flag.String("key", "thuthu", "user's key")
var userName = flag.String("name", "thuthu", "username")

func main() {
	flag.Parse()
	validateParams()
	log.Printf("Start server: %s, key: %s, name: %s", *server, *userKey, *userName)

	client := tokyo.NewClient(*server, *userKey, *userName)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 300)
		defer ticker.Stop()
		for {
			_ = <-ticker.C
			if !client.ConnReady {
				continue
			}
			otherPlayer, distance, err := client.GetClosestPlayer()
			if err != nil {
				log.Printf("Error when finding user: %v", err)
				client.Throttle(0)
				continue
			}
			client.HeadToPoint(otherPlayer.X, otherPlayer.Y)
			client.Throttle(1)
			if distance < 700 {
				client.Fire()
			}
		}
	}()
	log.Fatal(client.Listen())
}

func validateParams() {
	if server == nil {
		panic("miss server flag")
	}
	if userKey == nil {
		panic("miss key flag")
	}
	if userName == nil {
		panic("miss name flag")
	}
}
