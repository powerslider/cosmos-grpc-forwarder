package server

import (
	"context"
	"github.com/pkg/errors"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/jsonconv"
	"google.golang.org/grpc/metadata"
	"time"

	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
	"google.golang.org/grpc"
)

// NewLoggingInterceptor is a gRPC server interceptor for logging requests, responses and errors.
func NewLoggingInterceptor(logger log.Logger, jsonConverter *jsonconv.JSONConverter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		invoker grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if err := recover(); err != nil {
				reqJSON, err := jsonConverter.Marshal(req)
				if err != nil {
					logger.Error("error: request decoding: ", log.Error(errors.WithStack(err)))
				}

				logger.Panic("panicked gRPC request",
					log.String("method", info.FullMethod),
					log.String("request", string(reqJSON)),
					log.Error(err),
				)
				//panic(err)
			}
		}()

		start := time.Now()
		handlerResp, errResp := invoker(ctx, req)
		duration := time.Since(start)

		reqJSON, err := jsonConverter.Marshal(req)
		if err != nil {
			logger.Error("error: request decoding: ", log.Error(errors.WithStack(err)))
		}

		respJSON, err := jsonConverter.Marshal(handlerResp)
		if err != nil {
			logger.Error("error: response decoding: ", log.Error(errors.WithStack(err)))
		}

		md, _ := metadata.FromIncomingContext(ctx)

		headers, err := jsonConverter.Marshal(md)
		if err != nil {
			logger.Error("error: headers decoding: ", log.Error(errors.WithStack(err)))
		}

		logger.Print("gRPC request",
			log.String("request", string(reqJSON)),
			log.String("response", string(respJSON)),
			log.Error(errResp),
			log.Float64("duration", duration.Seconds()),
			log.String("headers", string(headers)),
		)

		return handlerResp, errResp
	}
}
