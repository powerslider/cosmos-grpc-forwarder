package testrunner

import (
	"context"
	"net"
	"testing"

	"github.com/powerslider/cosmos-grpc-forwarder/pkg/configs"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/forwarder"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/client"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/server"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/jsonconv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const (
	_bufSize        = 1024 * 1024
	_testServerName = "test-server"
	_testServerAddr = "BUFCONN"
)

// UnaryTestConfig handles config params for test runner setup.
type UnaryTestConfig struct {
	Logger             log.Logger
	Config             *configs.Config
	JSONConverter      *jsonconv.JSONConverter
	ClientInterceptors []grpc.UnaryClientInterceptor
	ServerInterceptors []grpc.UnaryServerInterceptor
	ClientOptions      []grpc.DialOption
	ServerOptions      []grpc.ServerOption
}

// HandleUnaryResponseError asserts gRPC error status codes.
func HandleUnaryResponseError(err error, t *testing.T) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			t.Errorf("gRPC Code: %s, response error message: %s", e.Code(), e.Message())
		}

		t.Errorf("raw response error: %s", err.Error())
	}
}

// NewDefaultTestConfig is a constructor function for a test runner config with sane defaults.
func NewDefaultTestConfig(
	logger log.Logger, config *configs.Config, jsonConverter *jsonconv.JSONConverter) UnaryTestConfig {
	return UnaryTestConfig{
		ServerInterceptors: []grpc.UnaryServerInterceptor{
			server.NewLoggingInterceptor(logger, jsonConverter),
		},
		ClientInterceptors: []grpc.UnaryClientInterceptor{
			client.NewLoggingInterceptor(logger, jsonConverter),
		},
		Config:        config,
		Logger:        logger,
		JSONConverter: jsonConverter,
	}
}

// NewUnaryTestSetup runs a gRPC server and provides a gRPC connection for clients.
// After gRPC calls are made in a test, the server performs graceful shutdown.
func NewUnaryTestSetup(
	ctx context.Context,
	config UnaryTestConfig,
) (*grpc.ClientConn, func(), error) {
	lis := bufconn.Listen(_bufSize)
	grpcServer := server.NewGRPCServer(
		_testServerName,
		_testServerAddr,
		lis,
		config.Logger,
		config.ServerInterceptors,
		config.ServerOptions...,
	)

	// TODO: This should be abstracted away in a gRPC service registration function.
	forwarder.InitializeGRPCHandlers(
		ctx,
		config.Config.CosmosSDKGRPCEndpoint,
		grpcServer,
		config.Logger,
		config.JSONConverter,
	)

	errCh := make(chan error)

	go grpcServer.Start(ctx, errCh)

	config.ClientOptions = append(config.ClientOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(getBufDialer(lis)),
	)

	closer := func() {
		lis.Close()
		grpcServer.Shutdown(ctx)
	}

	conn, err := client.NewGRPCConn(ctx, "", config.ClientInterceptors, config.ClientOptions...)

	return conn, closer, err
}

func getBufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return lis.Dial()
	}
}
