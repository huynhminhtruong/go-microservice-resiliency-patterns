package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/huynhminhtruong/go-microservice-resiliency-patterns/handle-errors/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	order.UnimplementedOrderServiceServer
}

func getTlsCredentials() (credentials.TransportCredentials, error) {
	serverCert, serverCertErr := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key") // Load server certificate
	if serverCertErr != nil {
		return nil, fmt.Errorf("could not load server key pairs: %s", serverCertErr)
	}

	certPool := x509.NewCertPool() // Certificate pool for CA check
	caCert, caCertErr := os.ReadFile("cert/ca.crt")
	if caCertErr != nil {
		return nil, fmt.Errorf("could not read CA cert: %s", caCertErr)
	}

	if ok := certPool.AppendCertsFromPEM(caCert); !ok { // Adds ca.crt to pool
		return nil, errors.New("failed to append the CA certs")
	}

	return credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAnyClientCert,      // Client authentication type
		Certificates: []tls.Certificate{serverCert}, // Provides server certificate
		ClientCAs:    certPool,                      // Roots the CA for the server to verify client certificates
	}), nil
}

func (s *server) Create(ctx context.Context, in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	return &order.CreateOrderResponse{OrderId: 1243}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	tlsCredentials, tlsCredentialsErr := getTlsCredentials()
	if tlsCredentialsErr != nil {
		log.Fatal("cannot load server TLS credentials: ", tlsCredentials)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(tlsCredentials)) // Adds TLS configuration to server options

	grpcServer := grpc.NewServer(opts...)
	order.RegisterOrderServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}
