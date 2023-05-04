package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/configs"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/forwarder"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/server"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/jsonconv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env.dist")
	if err != nil {
		panic(err)
	}

	conf := configs.InitializeConfig()

	logger := log.InitializeLogger(conf.LogLevel, conf.LogFormat)

	jsonConverter := jsonconv.NewJSONConverter()

	grpcServer := server.InitialiazeNewGRPCServer(ctx, conf, logger, jsonConverter)

	forwarder.InitializeGRPCHandlers(ctx, conf.CosmosSDKGRPCEndpoint, grpcServer, logger, jsonConverter)

	if err := grpcServer.Run(ctx); err != nil {
		logger.Panic("error starting the gRPC server: ", log.Error(err))
	}
}
