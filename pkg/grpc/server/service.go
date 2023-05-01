package server

import "google.golang.org/grpc"

// GRPCService contains gRPC service metadata needed for service
// registration when instantiating a new gRPC server.
type GRPCService struct {
	ServiceDesc *grpc.ServiceDesc
	Instance    any
}

// GRPCServiceRegistry is a collection of gRPC service metadata.
type GRPCServiceRegistry []*GRPCService

// NewGRPCServiceRegistry is a constructor function for createing a new GRPCServiceRegistry.
func NewGRPCServiceRegistry() *GRPCServiceRegistry {
	return &GRPCServiceRegistry{}
}

// Register registers a new gRPC service server instance.
func (r *GRPCServiceRegistry) Register(serviceDesc *grpc.ServiceDesc, instance any) {
	*r = append(*r, &GRPCService{
		ServiceDesc: serviceDesc,
		Instance:    instance,
	})
}
