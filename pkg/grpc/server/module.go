package server

import (
	"context"
	"fmt"

	"github.com/powerslider/cosmos-grpc-forwarder/pkg/configs"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"

	"google.golang.org/grpc"
)

// InitialiazeNewGRPCServer oeuoeu
func InitialiazeNewGRPCServer(
	ctx context.Context,
	conf *configs.Config,
	logger log.Logger,
) *Server {
	serverAddress := fmt.Sprintf("%s:%d", conf.ServerHost, conf.ServerPort)

	lis, err := NewListener(serverAddress)
	if err != nil {
		logger.Panic("error: cannot create server listener: ", log.Error(err))
	}

	return NewGRPCServer(
		conf.ServerName,
		serverAddress,
		lis,
		logger,
		[]grpc.UnaryServerInterceptor{
			//NewLoggingInterceptor(logger),
		},
	)
}
