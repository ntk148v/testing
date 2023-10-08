package chat

import (
	context "context"
	"log"
)

type Server struct{}

// mustEmbedUnimplementedChatServiceServer implements ChatServiceServer.
func (*Server) mustEmbedUnimplementedChatServiceServer() {
	panic("unimplemented")
}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from the client: %s", in.Body)
	return &Message{Body: "Hello from the server"}, nil
}
