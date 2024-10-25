package main

import (
	"context"
	"log"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/circuit-breaker-pattern-interceptor/order"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("locahost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderServiceClient(conn)
	log.Println("Creating order...")
	orderResponse, errCreate := orderClient.Create(context.Background(), &order.CreateOrderRequest{
		UserId:    -1,
		ProductId: 0,
		Price:     2,
	})

	if errCreate != nil {
		stat := status.Convert(errCreate) // The client converts an error into an error model
		for _, detail := range stat.Details() {
			switch errType := detail.(type) { // Based on error type prints error descriptions
			case *errdetails.BadRequest:
				for _, violation := range errType.GetFieldViolations() {
					log.Printf("The field %s has invalid value. desc: %v", violation.GetField(), violation.GetDescription())
				}
			}
		}
	} else {
		log.Printf("Order %d is created successfully.", orderResponse.OrderId)
	}
}
