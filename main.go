package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/devnsi/pubsub-direct-push/internal/push"
	"github.com/devnsi/pubsub-direct-push/internal/receive"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	receive.UnimplementedPublisherServer
	pusher *push.Pusher
}

func (s *server) Publish(_ context.Context, req *receive.PublishRequest) (*receive.PublishResponse, error) {
	log.Printf("Push published message: topic=%q messages=%d", req.Topic, len(req.Messages))
	var resp receive.PublishResponse
	for _, msg := range req.Messages {
		msg.MessageId = uuid.New().String() // because it's empty on requests from emulator.
		s.log(msg)
		s.push(msg, req)
		resp.MessageIds = append(resp.MessageIds, msg.MessageId)
	}
	return &resp, nil
}

func (s *server) log(msg *receive.PubsubMessage) {
	log.Printf("-> message received %s", msg.MessageId)
}

func (s *server) push(msg *receive.PubsubMessage, req *receive.PublishRequest) {
	msgPush := &push.Message{
		Message: push.MessageBody{
			Data:       base64.StdEncoding.EncodeToString([]byte(msg.Data)),
			MessageID:  msg.MessageId,
			Attributes: msg.Attributes,
		},
		Subscription: req.Topic,
	}
	if err := s.pusher.Push(msgPush); err != nil {
		log.Printf("failed to push message %s: %v", msg.MessageId, err)
	} else {
		log.Printf("-> message pushed %s", msg.MessageId)
	}
}

func main() {
	port := flag.Int("port", 8085, "gRPC server port")
	endpoint := flag.String("receiver", "http://localhost:8080/messages", "HTTP receiver endpoint")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pusher := push.New(*endpoint)
	grpcServer := grpc.NewServer()
	receive.RegisterPublisherServer(grpcServer, &server{pusher: pusher})

	log.Printf("Starting server on :%d, forwarding to %s", *port, *endpoint)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
