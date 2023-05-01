package server

import (
	"context"

	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
	"google.golang.org/grpc"
)

func NewLoggingInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		return nil, nil
	}
}
