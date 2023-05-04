package forwarder

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	pb "github.com/powerslider/cosmos-grpc-forwarder/client/grpc/api/cosmos/forwarder/v1"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/client"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/server"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/jsonconv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
)

// InitializeGRPCHandlers registers all gRPC handlers to the gRPC server and wires their dependencies.
func InitializeGRPCHandlers(
	ctx context.Context,
	cosmosSDKGRPCEndpoint string,
	grpcServer *server.Server,
	logger log.Logger,
	jsonConverter *jsonconv.JSONConverter,
) {
	grpcConn, err := client.NewDefaultGRPCConn(ctx, logger, jsonConverter, cosmosSDKGRPCEndpoint)
	if err != nil {
		logger.Panic("error: cannot create gRPC connection to Cosmos SDK endpoint: ", log.Error(err))
	}
	// defer grpcConn.Close()

	serviceServer := NewServiceHandler(tmservice.NewServiceClient(grpcConn))
	pb.RegisterServiceServer(grpcServer.Instance(), serviceServer)
}
