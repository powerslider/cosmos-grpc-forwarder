package forwarder

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	pb "github.com/powerslider/cosmos-grpc-forwarder/client/grpc/api/cosmos/forwarder/v1"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/server"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
)

func InitializeGRPCHandlers(
	ctx context.Context,
	cosmosSDKGRPCEndpoint string,
	grpcServer *server.Server,
	logger log.Logger) {

	grpcConn, err := NewCosmosSDKGRPCConn(ctx, logger, cosmosSDKGRPCEndpoint)
	if err != nil {
		logger.Panic("error: cannot create gRPC connection to Cosmos SDK endpoint: ", log.Error(err))
	}
	// defer grpcConn.Close()

	serviceServer := NewServiceHandler(tmservice.NewServiceClient(grpcConn))
	pb.RegisterServiceServer(grpcServer.Instance(), serviceServer)

}
