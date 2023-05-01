package client

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"

	"time"

	"google.golang.org/grpc"
)

func NewLoggingInterceptor(logger log.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		errResp := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start)

		reqJSON, err := json.MarshalIndent(req, "", "  ")
		if err != nil {
			logger.Error("error: request decoding: ", log.Error(errors.WithStack(err)))
		}

		respJSON, err := json.MarshalIndent(reply, "", "  ")
		if err != nil {
			logger.Error("error: request decoding: ", log.Error(errors.WithStack(err)))
		}

		logger.Print("outgoing gRPC request",
			log.String("method", method),
			log.String("request", string(reqJSON)),
			log.String("response", string(respJSON)),
			log.Error(errResp),
			log.Float64("duration", duration.Seconds()),
		)

		return errResp
	}
}
