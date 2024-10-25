package main

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/timeout-pattern/order"
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
	return &order.CreateOrderResponse{OrderId: 1234}, err
}
