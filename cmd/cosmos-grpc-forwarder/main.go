package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/configs"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/forwarder"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/server"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load(".env.dist")

	conf := configs.InitializeConfig()

	logFormat, err := log.ParseFormat(conf.LogFormat)
	if err != nil {
		panic(err)
	}

	logger := log.New(
		log.WithFormat(logFormat),
		log.AddCaller(),
		log.LogToStdout(),
	)

	grpcServer := server.InitialiazeNewGRPCServer(ctx, conf, logger)

	////grpcConn, err := forwarder.NewCosmosSDKGRPCConn(ctx, logger, conf.CosmosSDKGRPCEndpoint)
	//if err != nil {
	//	panic(err)
	//}
	//defer grpcConn.Close()

	//c := tmservice.NewServiceClient(grpcConn)
	//resp, err := c.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	//fmt.Println(resp)

	//serviceServer := forwarder.NewServiceHandler(tmservice.NewServiceClient(grpcConn))
	//pb.RegisterServiceServer(grpcServer.Instance(), serviceServer)
	forwarder.InitializeGRPCHandlers(ctx, conf.CosmosSDKGRPCEndpoint, grpcServer, logger)

	if err := grpcServer.Run(ctx); err != nil {
		logger.Panic("error starting the gRPC server: ", log.Error(err))
	}

}

func queryState() error {
	//// Create a connection to the gRPC server.
	//grpcConn, err := grpc.Dial(
	//	"grpc.osmosis.zone:9090", // your gRPC server address.
	//	grpc.WithInsecure(),      // The Cosmos SDK doesn't support any transport security mechanism.
	//	// This instantiates a general gRPC codec which handles proto bytes. We pass in a nil interface registry
	//	// if the request/response types contain interface instead of 'nil' you should pass the application specific codec.
	//	grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	//)
	return nil
}
