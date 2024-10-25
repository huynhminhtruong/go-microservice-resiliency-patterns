package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/circuit-breaker-pattern-simple/order"
	"google.golang.org/grpc"
)

type server struct {
	order.UnimplementedOrderServiceServer
}

func (s *server) Create(ctx context.Context, in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var err error
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	r := rand.New(src)
	if r.Intn(2) == 1 {
		err = errors.New("create order error")
	}
	return &order.CreateOrderResponse{OrderId: 1243}, err
}

func main() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	order.RegisterOrderServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}
