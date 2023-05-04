package main

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"

	"google.golang.org/grpc/status"

	"github.com/powerslider/cosmos-grpc-forwarder/pkg/jsonconv"

	"github.com/joho/godotenv"
	pb "github.com/powerslider/cosmos-grpc-forwarder/client/grpc/api/cosmos/forwarder/v1"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/configs"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/client"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
)

func main() {
	var err error

	ctx := context.Background()

	_ = godotenv.Load(".env.dist")

	conf := configs.InitializeConfig()

	logger := log.InitializeLogger(conf.LogLevel, conf.LogFormat)

	jsonConverter := jsonconv.NewJSONConverter()

	forwarderClient := getGPRCClient(ctx, logger, conf, jsonConverter)

	if _, err = forwarderClient.GetLatestBlock(ctx, &pb.GetLatestBlockRequest{}); err != nil {
		handleResponseError(err, logger)
	}

	if _, err = forwarderClient.GetBlockByHeight(ctx, &pb.GetBlockByHeightRequest{
		Height: 8658239,
	}); err != nil {
		handleResponseError(err, logger)
	}

	if _, err = forwarderClient.GetSyncing(ctx, &pb.GetSyncingRequest{}); err != nil {
		handleResponseError(err, logger)
	}

	if _, err = forwarderClient.GetNodeInfo(ctx, &pb.GetNodeInfoRequest{}); err != nil {
		handleResponseError(err, logger)
	}

	if _, err = forwarderClient.GetValidatorSetByHeight(ctx, &pb.GetValidatorSetByHeightRequest{
		Height:     8658239,
		Pagination: &query.PageRequest{},
	}); err != nil {
		handleResponseError(err, logger)
	}
}

func getGPRCClient(
	ctx context.Context,
	logger log.Logger,
	conf *configs.Config,
	jsonConverter *jsonconv.JSONConverter,
) pb.ServiceClient {
	conn, err := client.NewDefaultGRPCConn(
		ctx,
		logger,
		jsonConverter,
		fmt.Sprintf("%s:%d", conf.ServerHost, conf.ServerPort),
	)

	if err != nil {
		logger.Fatal("error: cannot create a gRPC connection: ", log.Error(err))
	}

	return pb.NewServiceClient(conn)
}

func handleResponseError(err error, logger log.Logger) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			logger.Panic(fmt.Sprintf("gRPC Code: %s, response error message: %s", e.Code(), e.Message()))
		}

		logger.Panic("raw response error: ", log.Error(err))
	}
}
