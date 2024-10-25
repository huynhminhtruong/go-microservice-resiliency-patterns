package main

import (
	"context"
	"log"
	"time"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/timeout-pattern/order"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cb *gobreaker.CircuitBreaker

func main() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "demo",
		MaxRequests: 3, // Allowed number ò requests for a half-open circuit
		Timeout:     4, // Timeout for an open to half-open transition
		ReadyToTrip: func(counts gobreaker.Counts) bool { // Decides on if the circuit will be open(quyết định xem mạch có mở hay không)
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.1
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Executed on each state change
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderServiceClient(conn)
	for { // Periodic call to simulate end-user requests
		log.Println("Creating order...")
		// Begins wrapping the Order service call with the circuit breaker
		orderResponse, errCreate := cb.Execute(func() (interface{}, error) {
			return orderClient.Create(context.Background(), &order.CreateOrderRequest{
				UserId:    23,
				ProductId: 1234,
				Price:     3.2,
			})
		})

		if errCreate != nil { // Error when state changes
			log.Printf("Failed to create order. Err: %v", errCreate)
		} else { // If the circuit is closed it returns data
			log.Printf("Order %d is created successfully", orderResponse.(*order.CreateOrderResponse).OrderId)
		}
		time.Sleep(1 * time.Second) // Waits for one second to not heat the CPU
	}
}
