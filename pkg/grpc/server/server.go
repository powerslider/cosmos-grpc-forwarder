package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/powerslider/cosmos-grpc-forwarder/pkg/log"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server contains all needed parameters for instantiating a new gRPC server.
type Server struct {
	Name           string
	Addr           string
	Listener       net.Listener
	serverInstance *grpc.Server
	logger         log.Logger
}

// NewGRPCServer creates a new instance of server.Server.
func NewGRPCServer(
	name string,
	addr string,
	listener net.Listener,
	logger log.Logger,
	// serviceRegistry *GRPCServiceRegistry,
	interceptors []grpc.UnaryServerInterceptor,
	opts ...grpc.ServerOption,
) *Server {
	// Set unary interceptors server option
	interceptorsOption := grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...))
	opts = append(opts, interceptorsOption)

	s := &Server{
		Name:           name,
		Addr:           addr,
		Listener:       listener,
		logger:         logger,
		serverInstance: grpc.NewServer(opts...),
	}

	// Enable server reflection feature.
	reflection.Register(s.serverInstance)

	return s
}

// Instance return the underlying instance of grpc.Server.
func (s *Server) Instance() *grpc.Server {
	return s.serverInstance
}

// Start starts an instantiated server.Server.
func (s *Server) Start(ctx context.Context, errChan chan error) {
	s.logger.Info(fmt.Sprintf("[Start] %s %s server starting on %s\n", s.Name, "gRPC", s.Addr))

	if err := s.serverInstance.Serve(s.Listener); err != nil {
		errChan <- err
	}
}

// Shutdown performs a graceful shutdown of an instance of server.Server.
func (s *Server) Shutdown(ctx context.Context) {
	s.logger.Info(fmt.Sprintf("[Shutdown] %s gRPC server is shutting down\n", s.Name))

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		if s.serverInstance != nil {
			s.serverInstance.GracefulStop()
		}
		done <- struct{}{}
	}()

	select {
	case <-ctxWithTimeout.Done():
		s.logger.Info("Timed out waiting for server to close.")
	case <-done:
		s.logger.Info("Server gracefully stopped.")
	}
}

// Run manages the gRPC server lifecycle on start and on shutdown.
func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error)

	go s.Start(ctx, errChan)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		s.Shutdown(ctx)

		return nil
	case err := <-errChan:
		return errors.WithStack(err)
	}
}
