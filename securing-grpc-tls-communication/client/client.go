package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/handle-errors/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func getTlsCredentials() (credentials.TransportCredentials, error) {
	clientCert, clientCertErr := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key") // Load client certificate
	if clientCertErr != nil {
		return nil, fmt.Errorf("could not load client key pair: %v", clientCertErr)
	}

	certPool := x509.NewCertPool() // Certificate pool for CA check
	caCert, caCertErr := os.ReadFile("cert/ca.crt")
	if caCertErr != nil {
		return nil, fmt.Errorf("could not read Cert CA: %v", caCertErr)
	}

	if ok := certPool.AppendCertsFromPEM(caCert); !ok { // Adds the CA to the certificate pool
		return nil, errors.New("failed to append CA cert")
	}

	return credentials.NewTLS(&tls.Config{
		ServerName:   "*.microservices.dev",         // Server name used during certificate generation
		Certificates: []tls.Certificate{clientCert}, // Provides the client certificate
		RootCAs:      certPool,                      // Roots the CA for the client to verify the server certificate
	}), nil
}

func main() {
	/*
		This implementation is useful for local development, but in a production environment,
		it is best practice to delegate certificate management to a third party, which we will see
		in detail in chapter 8 when we deploy microservices to a Kubernetes environment
	*/
	tlsCredentials, tlsCredentialsErr := getTlsCredentials()
	if tlsCredentialsErr != nil {
		log.Fatalf("failed to get tls credentials. %v", tlsCredentialsErr)
	}
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials)) // Adds TLS configuration to gRPC dial options
	conn, err := grpc.NewClient("localhost:8080", opts...)
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
		log.Fatalf("Failed to create order. %v", errCreate)
	} else {
		log.Printf("Order %d is created successfully", orderResponse.OrderId)
	}
}
