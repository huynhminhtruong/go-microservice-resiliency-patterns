package main

import (
	"context"
	"log"
	"time"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/timeout-pattern/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderServiceClient := order.NewOrderServiceClient(conn) // Create a connection to Order Server
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
	log.Println("Creating order...")
	_, errCreate := orderServiceClient.Create(ctx, &order.CreateOrderRequest{
		UserId:    23,
		ProductId: 123,
		Price:     12.3,
	})
	if errCreate != nil {
		log.Printf("Failed to create order. Err: %v", errCreate)
	} else {
		log.Println("Order is created successfully")
	}
}
