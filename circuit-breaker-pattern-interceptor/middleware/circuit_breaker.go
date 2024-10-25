package middleware

import (
	"context"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

func CircuitBreakerClientInterceptor(cb *gobreaker.CircuitBreaker) grpc.UnaryClientInterceptor { // A function with a UnaryClientInterceptor type
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		_, cbErr := cb.Execute(func() (interface{}, error) { // Circuit breaker wrapper
			err := invoker(ctx, method, req, reply, cc, opts...) // The gRPC call is resumed within a circuit breaker wrapper
			if err != nil {
				return nil, err
			}

			return nil, nil
		})
		return cbErr
	}
}
