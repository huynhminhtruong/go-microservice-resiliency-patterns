package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/timeout-pattern/product"
	"google.golang.org/grpc"
)

type server struct {
	product.UnimplementedProductServiceServer
}

func (s *server) Get(ctx context.Context, in *product.GetProductRequest) (*product.GetProductResponse, error) {
	time.Sleep(2 * time.Second)
	return &product.GetProductResponse{Title: "Demo title"}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8081))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	product.RegisterProductServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}
