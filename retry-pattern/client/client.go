package main

import (
	"context"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/retry-pattern/shipping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor( // Enables interceptor for Unary connections
		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted), // Retries only once those codes are returned
		grpc_retry.WithMax(5), // Retrues five times at max
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)), // Uses a one-second timeout for each retry
	)))
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("locahost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect shipping service. Err: %v", err)
	}

	defer conn.Close()

	shippingServiceClient := shipping.NewShippingServiceClient(conn)
	log.Println("Creating shipping...")
	_, errCreate := shippingServiceClient.Create(context.Background(), &shipping.CreateShippingRequest{
		UserId:  23,
		OrderId: 2344,
	})
	if errCreate != nil {
		log.Printf("Failed to create shipping. Err: %v", errCreate) // May return ContextDeadLineExceeded
	} else {
		log.Println("Shipping is created successfully")
	}
}
