package client

import (
	"context"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// NewGRPCConn creates a new instance of a gRPC client connection using a context
func NewGRPCConn(
	ctx context.Context,
	serverAddr string,
	interceptors []grpc.UnaryClientInterceptor,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	// Set unary interceptors client option.
	interceptorsOption := grpc.WithUnaryInterceptor(grpcmiddleware.ChainUnaryClient(interceptors...))
	opts = append(opts, interceptorsOption)

	// Set client blocking option which waits for the server connection before returning.
	//opts = append(opts, grpc.WithBlock())

	conn, err := grpc.DialContext(ctx, serverAddr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
