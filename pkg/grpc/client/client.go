package client

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/jsonconv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
	"google.golang.org/grpc/credentials/insecure"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// NewGRPCConn creates a new instance of a gRPC client connection using a context.
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

// NewDefaultGRPCConn is a constructor function with sane defaults for gRPC options and interceptors.
func NewDefaultGRPCConn(
	ctx context.Context,
	logger log.Logger,
	jsonConverter *jsonconv.JSONConverter,
	serverAddr string,
) (*grpc.ClientConn, error) {
	return NewGRPCConn(
		ctx,
		serverAddr,
		[]grpc.UnaryClientInterceptor{
			NewLoggingInterceptor(logger, jsonConverter),
		},
		// The Cosmos SDK doesn't support any transport security mechanism.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// This instantiates a general gRPC codec which handles proto bytes. We pass in a nil interface registry
		// if the request/response types contain interface instead of 'nil' you should pass the application specific codec.
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
}
