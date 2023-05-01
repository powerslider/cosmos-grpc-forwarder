package forwarder

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"

	pb "github.com/powerslider/cosmos-grpc-forwarder/client/grpc/api/cosmos/forwarder/v1"
)

type ServiceHandler struct {
	ServiceGRPCClient tmservice.ServiceClient
	*pb.UnimplementedServiceServer
}

func NewServiceHandler(client tmservice.ServiceClient) *ServiceHandler {
	return &ServiceHandler{
		ServiceGRPCClient:          client,
		UnimplementedServiceServer: &pb.UnimplementedServiceServer{},
	}
}

func (h *ServiceHandler) GetNodeInfo(ctx context.Context, req *pb.GetNodeInfoRequest) (*pb.GetNodeInfoResponse, error) {
	resp, err := h.ServiceGRPCClient.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
	if err != nil {
		return nil, err
	}

	appVersion := resp.GetApplicationVersion()
	return &pb.GetNodeInfoResponse{
		DefaultNodeInfo: resp.DefaultNodeInfo,
		ApplicationVersion: &pb.VersionInfo{
			Name:             appVersion.GetName(),
			AppName:          appVersion.GetAppName(),
			Version:          appVersion.GetVersion(),
			GitCommit:        appVersion.GetGitCommit(),
			BuildTags:        appVersion.GetBuildTags(),
			GoVersion:        appVersion.GetGoVersion(),
			BuildDeps:        remapBuildDeps(appVersion.GetBuildDeps()),
			CosmosSdkVersion: appVersion.GetCosmosSdkVersion(),
		},
	}, nil
}

func (h *ServiceHandler) GetSyncing(ctx context.Context, req *pb.GetSyncingRequest) (*pb.GetSyncingResponse, error) {
	resp, err := h.ServiceGRPCClient.GetSyncing(ctx, &tmservice.GetSyncingRequest{})
	if err != nil {
		return nil, err
	}

	return &pb.GetSyncingResponse{
		Syncing: resp.Syncing,
	}, nil
}

func (h *ServiceHandler) GetLatestBlock(ctx context.Context, req *pb.GetLatestBlockRequest) (*pb.GetLatestBlockResponse, error) {
	resp, err := h.ServiceGRPCClient.GetLatestBlock(ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}

	return &pb.GetLatestBlockResponse{
		BlockId:  resp.GetBlockId(),
		Block:    resp.GetBlock(),
		SdkBlock: remapSDKBlock(resp.GetSdkBlock()),
	}, nil
}

func (h *ServiceHandler) GetBlockByHeight(ctx context.Context, req *pb.GetBlockByHeightRequest) (*pb.GetBlockByHeightResponse, error) {
	resp, err := h.ServiceGRPCClient.GetBlockByHeight(ctx, &tmservice.GetBlockByHeightRequest{
		Height: req.Height,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetBlockByHeightResponse{
		BlockId:  resp.GetBlockId(),
		Block:    resp.GetBlock(),
		SdkBlock: remapSDKBlock(resp.GetSdkBlock()),
	}, nil
}

func (h *ServiceHandler) GetLatestValidatorSet(
	ctx context.Context, req *pb.GetLatestValidatorSetRequest) (*pb.GetLatestValidatorSetResponse, error) {
	resp, err := h.ServiceGRPCClient.GetLatestValidatorSet(ctx, &tmservice.GetLatestValidatorSetRequest{
		Pagination: req.Pagination,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetLatestValidatorSetResponse{
		BlockHeight: resp.GetBlockHeight(),
		Validators:  remapValidators(resp.GetValidators()),
		Pagination:  resp.GetPagination(),
	}, nil
}

func (h *ServiceHandler) GetValidatorSetByHeight(
	ctx context.Context, req *pb.GetValidatorSetByHeightRequest) (*pb.GetValidatorSetByHeightResponse, error) {
	resp, err := h.ServiceGRPCClient.GetValidatorSetByHeight(ctx, &tmservice.GetValidatorSetByHeightRequest{
		Height:     req.Height,
		Pagination: req.Pagination,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetValidatorSetByHeightResponse{
		BlockHeight: resp.GetBlockHeight(),
		Validators:  remapValidators(resp.GetValidators()),
		Pagination:  resp.GetPagination(),
	}, nil
}

func (h *ServiceHandler) ABCIQuery(ctx context.Context, req *pb.ABCIQueryRequest) (*pb.ABCIQueryResponse, error) {
	resp, err := h.ServiceGRPCClient.ABCIQuery(ctx, &tmservice.ABCIQueryRequest{
		Data:   req.Data,
		Path:   req.Path,
		Height: req.Height,
		Prove:  req.Prove,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ABCIQueryResponse{
		Code:      resp.GetCode(),
		Log:       resp.GetLog(),
		Info:      resp.GetInfo(),
		Index:     resp.GetIndex(),
		Key:       resp.GetKey(),
		Value:     resp.GetValue(),
		ProofOps:  remapProofOps(resp.GetProofOps()),
		Height:    resp.GetHeight(),
		Codespace: resp.GetCodespace(),
	}, nil
}

func remapValidators(resp []*tmservice.Validator) []*pb.Validator {
	validators := make([]*pb.Validator, 0)

	for _, v := range resp {
		validators = append(validators, &pb.Validator{
			Address:          v.GetAddress(),
			PubKey:           v.GetPubKey(),
			VotingPower:      v.GetVotingPower(),
			ProposerPriority: v.GetProposerPriority(),
		})
	}

	return validators
}

func remapBuildDeps(resp []*tmservice.Module) []*pb.Module {
	modules := make([]*pb.Module, 0)

	for _, m := range resp {
		modules = append(modules, &pb.Module{
			Path:    m.Path,
			Version: m.Version,
			Sum:     m.Sum,
		})
	}

	return modules
}

func remapProofOps(resp *tmservice.ProofOps) *pb.ProofOps {
	proofOps := make([]pb.ProofOp, 0)
	for _, p := range resp.Ops {
		proofOps = append(proofOps, pb.ProofOp{
			Type: p.GetType(),
			Key:  p.GetKey(),
			Data: p.GetData(),
		})
	}

	return &pb.ProofOps{
		Ops: proofOps,
	}
}

func remapSDKBlock(resp *tmservice.Block) *pb.Block {
	header := resp.GetHeader()

	return &pb.Block{
		Header: pb.Header{
			Version:            header.GetVersion(),
			ChainID:            header.GetChainID(),
			Height:             header.GetHeight(),
			Time:               header.GetTime(),
			LastBlockId:        header.GetLastBlockId(),
			LastCommitHash:     header.GetLastCommitHash(),
			DataHash:           header.GetDataHash(),
			ValidatorsHash:     header.GetValidatorsHash(),
			NextValidatorsHash: header.GetNextValidatorsHash(),
			ConsensusHash:      header.GetConsensusHash(),
			AppHash:            header.GetAppHash(),
			LastResultsHash:    header.GetLastResultsHash(),
			EvidenceHash:       header.GetEvidenceHash(),
			ProposerAddress:    header.GetProposerAddress(),
		},
		Data:       resp.GetData(),
		Evidence:   resp.GetEvidence(),
		LastCommit: resp.GetLastCommit(),
	}
}

//func (h *ServiceHandler) serviceGRPCClient(ctx context.Context) (pb.ServiceClient, error) {
//	grpcConn, err := NewCosmosSDKGRPCConn(ctx, h.Logger, "grpc.osmosis.zone:9090")
//	if err != nil {
//		return nil, err
//	}
//
//	c := tmservice.NewServiceClient(grpcConn)
//
//	return c, nil
//}
