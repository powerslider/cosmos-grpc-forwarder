package forwarder_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/gogoproto/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/joho/godotenv"
	pb "github.com/powerslider/cosmos-grpc-forwarder/client/grpc/api/cosmos/forwarder/v1"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/configs"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/client"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/grpc/testrunner"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/jsonconv"
	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"
)

type testClients struct {
	forwarderClient pb.ServiceClient
	originClient    tmservice.ServiceClient
}

func TestServiceHandlerForwardedCalls(t *testing.T) {
	ctx := context.Background()

	_ = godotenv.Load("../../.env.test.dist")

	appConfig := configs.InitializeConfig()

	logger := log.InitializeLogger(appConfig.LogLevel, appConfig.LogFormat)

	jsonConverter := jsonconv.NewJSONConverter()

	grpcClients, closer, err := setupTest(ctx, t, appConfig, logger, jsonConverter)
	defer closer()

	if err != nil {
		t.Error(err)
	}

	verifyGetLatestBlock(ctx, t, grpcClients, jsonConverter)
	verifyGetBlockByHeight(ctx, t, grpcClients, jsonConverter)
	verifyGetSyncing(ctx, t, grpcClients, jsonConverter)
	verifyGetNodeInfo(ctx, t, grpcClients, jsonConverter)
	verifyGetValidatorSetByHeight(ctx, t, grpcClients, jsonConverter)
	//verifyGetLatestValidatorSet(ctx, t, grpcClients, jsonConverter)
	//verifyABCIQuery(ctx, t, grpcClients, jsonConverter)

}

func verifyGetLatestBlock(
	ctx context.Context,
	t *testing.T,
	grpcClients testClients,
	jsonConverter *jsonconv.JSONConverter,
) {
	resp, err := grpcClients.forwarderClient.GetLatestBlock(ctx, &pb.GetLatestBlockRequest{})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	originResp, err := grpcClients.originClient.GetLatestBlock(ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	verifyResponses(t, resp, originResp, jsonConverter)
}

func verifyGetBlockByHeight(
	ctx context.Context,
	t *testing.T,
	grpcClients testClients,
	jsonConverter *jsonconv.JSONConverter,
) {
	resp, err := grpcClients.forwarderClient.GetBlockByHeight(ctx, &pb.GetBlockByHeightRequest{
		Height: 8658239,
	})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	originResp, err := grpcClients.originClient.GetBlockByHeight(ctx, &tmservice.GetBlockByHeightRequest{
		Height: 8658239,
	})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	verifyResponses(t, resp, originResp, jsonConverter)
}

func verifyGetSyncing(
	ctx context.Context,
	t *testing.T,
	grpcClients testClients,
	jsonConverter *jsonconv.JSONConverter,
) {
	resp, err := grpcClients.forwarderClient.GetSyncing(ctx, &pb.GetSyncingRequest{})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	originResp, err := grpcClients.originClient.GetSyncing(ctx, &tmservice.GetSyncingRequest{})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	verifyResponses(t, resp, originResp, jsonConverter)
}

func verifyGetNodeInfo(
	ctx context.Context,
	t *testing.T,
	grpcClients testClients,
	jsonConverter *jsonconv.JSONConverter,
) {
	resp, err := grpcClients.forwarderClient.GetNodeInfo(ctx, &pb.GetNodeInfoRequest{})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	originResp, err := grpcClients.originClient.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	verifyResponses(t, resp, originResp, jsonConverter)
}

func verifyGetValidatorSetByHeight(
	ctx context.Context,
	t *testing.T,
	grpcClients testClients,
	jsonConverter *jsonconv.JSONConverter,
) {
	resp, err := grpcClients.forwarderClient.GetValidatorSetByHeight(ctx, &pb.GetValidatorSetByHeightRequest{
		Height:     8658239,
		Pagination: &query.PageRequest{},
	})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	originResp, err := grpcClients.originClient.GetValidatorSetByHeight(ctx, &tmservice.GetValidatorSetByHeightRequest{
		Height:     8658239,
		Pagination: &query.PageRequest{},
	})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	verifyResponses(t, resp, originResp, jsonConverter)
}

//nolint:unused
func verifyGetLatestValidatorSet(
	ctx context.Context,
	t *testing.T,
	grpcClients testClients,
	jsonConverter *jsonconv.JSONConverter,
) {
	resp, err := grpcClients.forwarderClient.GetLatestValidatorSet(ctx, &pb.GetLatestValidatorSetRequest{
		Pagination: &query.PageRequest{},
	})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	originResp, err := grpcClients.originClient.GetLatestValidatorSet(ctx, &tmservice.GetLatestValidatorSetRequest{
		Pagination: &query.PageRequest{},
	})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	verifyResponses(t, resp, originResp, jsonConverter)
}

//nolint:unused
func verifyABCIQuery(
	ctx context.Context,
	t *testing.T,
	grpcClients testClients,
	jsonConverter *jsonconv.JSONConverter,
) {
	resp, err := grpcClients.forwarderClient.ABCIQuery(ctx, &pb.ABCIQueryRequest{
		Data:   []byte{123},
		Path:   "",
		Height: 8658239,
		Prove:  false,
	})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	originResp, err := grpcClients.originClient.ABCIQuery(ctx, &tmservice.ABCIQueryRequest{
		Data:   []byte{123},
		Path:   "",
		Height: 8658239,
		Prove:  false,
	})
	if err != nil {
		testrunner.HandleUnaryResponseError(err, t)
	}

	verifyResponses(t, resp, originResp, jsonConverter)
}

func setupTest(
	ctx context.Context,
	t *testing.T,
	appConfig *configs.Config,
	logger log.Logger,
	jsonConverter *jsonconv.JSONConverter,
) (testClients, func(), error) {
	testConfig := testrunner.NewDefaultTestConfig(logger, appConfig, jsonConverter)

	conn, closer, err := testrunner.NewUnaryTestSetup(ctx, testConfig)
	if err != nil {
		t.Error(err)
	}

	cosmosConn, err := client.NewDefaultGRPCConn(
		ctx,
		logger,
		jsonConverter,
		appConfig.CosmosSDKGRPCEndpoint,
	)
	if err != nil {
		t.Error(err)
	}

	forwarderClient := pb.NewServiceClient(conn)
	cosmosClient := tmservice.NewServiceClient(cosmosConn)

	return testClients{
		forwarderClient: forwarderClient,
		originClient:    cosmosClient,
	}, closer, nil
}

func verifyResponses(
	t *testing.T, resp proto.Message,
	originResp proto.Message,
	jsonConverter *jsonconv.JSONConverter,
) {
	respJson, errRespMarshal := jsonConverter.Marshal(resp)
	originRespJson, errOriginRespMarshal := jsonConverter.Marshal(originResp)

	errMarshal := errors.Join(errRespMarshal, errOriginRespMarshal)
	if errMarshal != nil {
		t.Error(errMarshal)
	}

	errDiff := compare(respJson, originRespJson)
	if errDiff != nil {
		t.Error(errDiff)
	}
}

func compare(actual, expected any) error {
	diff := cmp.Diff(expected, actual)
	if diff == "" {
		return nil
	}

	return fmt.Errorf("Objects should match.\n%s", diff)
}
