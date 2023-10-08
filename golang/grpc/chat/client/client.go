package main

import (
	"context"
	"log"

	"github.com/ntk148v/testing/golang/grpc/chat"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	defer conn.Close()

	c := chat.NewChatServiceClient(conn)
	response, err := c.SayHello(context.Background(), &chat.Message{Body: "Hello from the client"})
	if err != nil {
		log.Fatalf("error when calling SayHello: %s", err)
	}

	log.Printf("Response from server: %s", response.Body)
}
