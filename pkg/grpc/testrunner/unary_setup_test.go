package testrunner_test

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/configs"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/jsonconv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"

	"github.com/pkg/errors"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/testrunner"
)

func TestClientConnection(t *testing.T) {
	ctx := context.Background()

	_ = godotenv.Load("../../../.env.test.dist")

	conf := configs.InitializeConfig()

	logger := log.InitializeLogger(conf.LogLevel, conf.LogFormat)

	jsonConverter := jsonconv.NewJSONConverter()

	config := testrunner.NewDefaultTestConfig(logger, conf, jsonConverter)

	conn, closer, err := testrunner.NewUnaryTestSetup(ctx, config)

	if err != nil {
		t.Error(err)
	}

	if conn == nil {
		t.Error(errors.New("connection cannot be nil"))
	}

	closer()
}
