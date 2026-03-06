package main

import (
  "context"
  "flag"
  "fmt"
  "log"
  "net"
  "sync/atomic"

  pb "github.com/devnsi/pubsub-direct-push/internal/handler"
  "google.golang.org/grpc"
)

type server struct {
  pb.UnimplementedPublisherServer
  nextID uint64
}

func (s *server) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
  // Log topic and count
  log.Printf("Publish called: topic=%q messages=%d", req.Topic, len(req.Messages))

  var resp pb.PublishResponse
  for _, msg := range req.Messages {
    // Assign or echo an ID
    idNum := atomic.AddUint64(&s.nextID, 1)
    msgID := fmt.Sprintf("%d", idNum)
    resp.MessageIds = append(resp.MessageIds, msgID)

    // Decode data (if you want) and log
    log.Printf(
      "  → message_id=%s client_id=%q attributes=%v data=%q publish_time=%v",
      msgID, msg.MessageId, msg.Attributes, string(msg.Data), msg.PublishTime,
    )
  }
  return &resp, nil
}

func main() {
  port := flag.Int("port", 50051, "gRPC server port")
  flag.Parse()

  lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  grpcServer := grpc.NewServer()
  pb.RegisterPublisherServer(grpcServer, &server{})

  log.Printf("Starting Pub/Sub gRPC stub on :%d …", *port)
  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalf("gRPC server error: %v", err)
  }
}
