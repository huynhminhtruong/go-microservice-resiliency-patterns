package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/retry-pattern/shipping"
	"google.golang.org/grpc"
)

type server struct {
	shipping.UnimplementedShippingServiceServer
}

func (s *server) Create(ctx context.Context, in *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	time.Sleep(2 * time.Second)
	return &shipping.CreateShippingResponse{ShippingId: 1234}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	shipping.RegisterShippingServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}
